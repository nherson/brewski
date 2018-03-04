package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/nherson/brewski/temperature"
)

func main() {
	parseFlags()
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	logger.Info("brewski starting",
		zap.String("component", "init"),
	)

	pollingInterval := GetPollingInterval()

	// Set the sysfs device path
	temperature.SetOnewireSysfsDir(GetOnewireSysfsDir())

	// Create a logger for temperature sensors
	tempLogger := logger.With(zap.String("component", "temperature"))

	var temperatureSensors []temperature.Poller
	for _, id := range GetDS18B20IDs() {
		// Create a ChainCallback to record temperature readings
		callbackChain := temperature.NewChainCallback()
		// Add a logging callback to the callback chain... it's always a good
		// idea to log what is going on!
		callbackChain.RegisterCallback(
			temperature.NewLoggingCallback(tempLogger),
		)
		// If influxDB is enabled, create a callback to send sensor readings to it and
		// register the callback into the chain
		if UseInfluxDB() {
			tags := make(map[string]string)
			tags["device"] = id
			database := GetInfluxDBDatabase()
			endpoint := GetInfluxDBEndpoint()
			influxDBCallback, err := temperature.NewInfluxDBCallback(endpoint, database, tags)
			if err != nil {
				logger.Fatal(fmt.Sprintf("error initializing influxdb callback: %s", err.Error()),
					zap.String("component", "init"),
				)
			}
			callbackChain.RegisterCallback(influxDBCallback)
		}
		sensor := temperature.NewDS18B20Sensor(
			id,
			pollingInterval,
			tempLogger.With(zap.String("device", id)),
		)
		sensor.SetCallback(callbackChain)
		temperatureSensors = append(temperatureSensors, sensor)
		sensor.Start()
	}
	waitForExit(logger)
}

// adapted from https://gobyexample.com/signals
func waitForExit(logger *zap.Logger) {
	stop := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-stop
		fmt.Println()
		logger.Info("received stop, shutting down brewski")
		done <- true
	}()
	<-done

	os.Exit(0)
}

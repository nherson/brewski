package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/nherson/brewski/device"
	"github.com/nherson/brewski/handlers"
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
	device.SetOnewireSysfsDir(GetOnewireSysfsDir())

	// Create a logger for temperature sensors
	tempLogger := logger.With(zap.String("component", "temperature"))

	var temperatureSensors []device.Poller
	for _, id := range GetDS18B20IDs() {
		// Create a ChainCallback to record temperature readings
		callbackChain := handlers.NewChainCallback()
		// Add a logging callback to the callback chain... it's always a good
		// idea to log what is going on!
		callbackChain.RegisterCallback(
			handlers.NewLoggingCallback(tempLogger),
		)
		// If influxDB is enabled, create a callback to send sensor readings to it and
		// register the callback into the chain
		if UseInfluxDB() {
			database := GetInfluxDBDatabase()
			endpoint := GetInfluxDBEndpoint()
			influxDBCallback, err := handlers.NewInfluxDBCallback(endpoint, database)
			if err != nil {
				logger.Fatal(fmt.Sprintf("error initializing influxdb callback: %s", err.Error()),
					zap.String("component", "init"),
				)
			}
			callbackChain.RegisterCallback(influxDBCallback)
		}
		ds18b20Reader := device.NewDS18B20(id)
		deviceSensor := device.NewSensor(
			ds18b20Reader,
			pollingInterval,
			tempLogger.With(zap.String("device", id)),
		)
		deviceSensor.SetCallback(callbackChain)
		temperatureSensors = append(temperatureSensors, deviceSensor)
		deviceSensor.Start()
	}
	// Just try making a tilt device and tracking it with basic logging
	tilt, err := device.NewTiltHydrometer()
	tiltChain := handlers.NewChainCallback()
	tiltChain.RegisterCallback(
		handlers.NewLoggingCallback(logger),
	)
	tiltSensor := device.NewSensor(
		tilt,
		pollingInterval,
		logger.With(zap.String("device", "tilt")),
	)
	tiltSensor.SetCallback(tiltChain)
	tiltSensor.Start()
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

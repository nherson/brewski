package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/BurntSushi/toml"

	"github.com/nherson/brewski/config"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "./config.toml", "Path to brewski config file")
}

func main() {
	flag.Parse()

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	mainLogger := logger.With(zap.String("component", "init"))
	mainLogger.Info("brewski starting")

	mainLogger.Info("reading config file at " + configPath)
	b, err := ioutil.ReadFile(configPath)
	if err != nil {
		mainLogger.Fatal("could not read config file", zap.Error(err))
	}

	var conf config.Config
	if _, err := toml.Decode(string(b), &conf); err != nil {
		// handle error
		mainLogger.Fatal("could not parse config file", zap.Error(err))
	}

	// TODO: make this take a logger
	pollers, err := conf.Generate()
	if err != nil {
		mainLogger.Fatal("could not generate sensors", zap.Error(err))
	}
	for _, p := range pollers {
		p.Start()
	}

	waitForExit(mainLogger)
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

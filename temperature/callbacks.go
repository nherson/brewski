package temperature

// This file contains callbacks which can do different operations with
// the read temperature data

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/influxdata/influxdb/client/v2"
	"go.uber.org/zap"
)

// SensorCallback is an interface which a sensor can use to pass data off
// to for arbitrary processing.
type SensorCallback interface {
	Handle(float32, float32) error
}

// LoggingCallback includes structured logging for data
type LoggingCallback struct {
	logger *zap.Logger
}

// NewLoggingCallback returns a LoggingCallback
func NewLoggingCallback(l *zap.Logger) *LoggingCallback {
	return &LoggingCallback{
		logger: l,
	}
}

// Handle is a function that adheres to the TemperaturePoller interface
// Use this function after generating the struct to pass into the
// SetCallback function to use this method as a callback for temperature readings
// The intermediate struct is because this callback is stateful due to the logger
// needing to be predefined (this makes the logger more flexible in config)
func (l *LoggingCallback) Handle(c, f float32) error {
	l.logger.Info("temperature successfully read",
		zap.Float32("fahrenheit", f),
		zap.Float32("celsius", c),
	)
	return nil
}

// StdoutCallback is a basic printing callback handler
type StdoutCallback struct{}

// NewStdoutCallback returns a new stdout callback...
// it just prints things to stdout.
func NewStdoutCallback() *StdoutCallback {
	return &StdoutCallback{}
}

// Handle is a simple callback that just prints to stdout
func (sc *StdoutCallback) Handle(f, c float32) error {
	fmt.Printf("Temperature read! C=%f F=%f\n", c, f)
	return nil
}

// ChainCallback can be used to include multiple sub-callback
// implementations into a single callback handler
type ChainCallback struct {
	callbacks []SensorCallback
}

// NewChainCallback returns an empty ChainCallback
func NewChainCallback() *ChainCallback {
	return &ChainCallback{
		callbacks: make([]SensorCallback, 0),
	}
}

// RegisterCallback adds a callback function to the list of callbacks in the chain
func (cc *ChainCallback) RegisterCallback(scb SensorCallback) {
	cc.callbacks = append(cc.callbacks, scb)
}

// Handle iterates through all registered callbacks and executes them
func (cc *ChainCallback) Handle(c, f float32) error {
	var errList *multierror.Error
	for _, cb := range cc.callbacks {
		err := cb.Handle(c, f)
		if err != nil {
			multierror.Append(errList, err)
		}
	}
	return errList.ErrorOrNil()
}

// InfluxDBCallback sends sensor data to a specified InfluxDB endpoint
type InfluxDBCallback struct {
	c    client.Client
	tags map[string]string
	db   string
}

// NewInfluxDBCallback returns an InfluxDBCallback that can be used to send sensor data to the configured
// influxDB endpoint.  The passed in tags will be included in every callback invocation
func NewInfluxDBCallback(addr string, database string, tags map[string]string) (*InfluxDBCallback, error) {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: "http://localhost:8086",
	})
	if err != nil {
		return nil, err
	}
	return &InfluxDBCallback{
		c:    c,
		tags: tags,
		db:   database,
	}, nil
}

// Handle sends sensor data to the configured InfluxDB database
func (icb *InfluxDBCallback) Handle(c, f float32) error {
	// create a new batch of points
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  icb.db,
		Precision: "s",
	})
	if err != nil {
		return err
	}
	// Turn the measured temperatures into InfluxDB fields
	fields := map[string]interface{}{
		"celsius":    c,
		"fahrenheit": f,
	}
	pt, err := client.NewPoint("temperature", icb.tags, fields, time.Now())
	if err != nil {
		return err
	}
	bp.AddPoint(pt)
	err = icb.c.Write(bp)
	if err != nil {
		return err
	}
	return nil
}

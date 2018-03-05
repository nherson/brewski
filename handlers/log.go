package handlers

import (
	"github.com/nherson/brewski/measurement"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

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

// Handle is a function that adheres to the Callback interface
// Use this function after generating the struct to pass into the
// SetCallback function to use this method as a callback for temperature readings
// The intermediate struct is because this callback is stateful due to the logger
// needing to be predefined (this makes the logger more flexible in config)
func (l *LoggingCallback) Handle(s measurement.Sample) error {
	sampleFields := []zapcore.Field{}
	sampleFields = append(sampleFields, zap.String("device", s.DeviceName()))
	for _, d := range s.Datapoints() {
		sampleFields = append(sampleFields, zap.Float32(d.Name(), d.Value()))
	}

	l.logger.Info("device successfully read",
		sampleFields...,
	)
	return nil
}

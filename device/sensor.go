package device

import (
	"time"

	"github.com/nherson/brewski/handlers"
	"go.uber.org/zap"
)

// Poller interface for a sensor, which can periodically poll
// and submit temperature data
type Poller interface {
	Start()
	Stop()
	SetCallback(handlers.Callback)
}

// Sensor is a simple implementation of a TemperaturePoller
// It sleeps for
type Sensor struct {
	reader   Reader
	logger   *zap.Logger
	interval time.Duration
	control  chan bool
	callback handlers.Callback
}

// NewSensor creates a new polling sensor for a given device reader.
// The device reader will provide measurement samples, and this sensor
// will act as a harness to process retrieved datapoints.
// By default, the sensor will start with a simple stdout handler
// until provided something more specific using the SetCallback method
func NewSensor(r Reader, i time.Duration, l *zap.Logger) *Sensor {
	return &Sensor{
		reader:   r,
		logger:   l,
		interval: i,
		control:  make(chan bool, 1),
		callback: handlers.NewStdoutCallback(),
	}
}

// SetCallback assigns a callback function for the sensor
// for when polling is complete
func (s *Sensor) SetCallback(scb handlers.Callback) {
	s.callback = scb
}

// Start begins the polling process. Given a polling interval, the sensor will
// sleep for the specified amount of time, poll the sensor for a temperature
// reading, and submit the results into the callback function
// This method is non-blocking and will spin off a go routine and return immediately.
func (s *Sensor) Start() {
	// Repeatedly sleep and poll, forever...
	go func() {
		intervalTicker := time.NewTicker(s.interval).C
		for {
			select {
			case <-intervalTicker:
				// Read the temperature from the device
				sample, err := s.reader.Read()
				if err != nil {
					s.logger.Error("error reading data from device",
						zap.String("device", s.reader.Name()),
						zap.String("error", err.Error()),
					)
				}
				// Process the temperature reading
				err = s.callback.Handle(sample)
				if err != nil {
					s.logger.Error("error recording temperature reading",
						zap.String("error", err.Error()),
					)
				}
			case <-s.control:
				return
			}
		}
	}()
}

// Stop tells the sensor to stop periodically polling for data
func (s *Sensor) Stop() {
	s.control <- true
}

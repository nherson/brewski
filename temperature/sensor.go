package temperature

import (
	"time"

	"go.uber.org/zap"
)

// Poller interface for a sensor, which can periodically poll
// and submit temperature data
type Poller interface {
	Start()
	Stop()
	SetCallback(SensorCallback)
}

// Sensor is a simple implementation of a TemperaturePoller
// It sleeps for
type Sensor struct {
	tr       Reader
	logger   *zap.Logger
	interval time.Duration
	control  chan bool
	callback SensorCallback
}

// NewDS18B20Sensor creates a new polling sensor for a DS18B20 temperature probe
// with the given device ID using onewire protocol via w1-gpio kernel module
// using sysfs for reading data.
func NewDS18B20Sensor(deviceID string, i time.Duration, l *zap.Logger) *Sensor {
	ds18b20 := NewDS18B20(deviceID)
	return &Sensor{
		tr:       ds18b20,
		logger:   l,
		interval: i,
		control:  make(chan bool, 1),
		callback: NewStdoutCallback(),
	}
}

// SetCallback assigns a callback function for the sensor
// for when polling is complete
func (s *Sensor) SetCallback(scb SensorCallback) {
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
				c, f, err := s.tr.ReadTemperature()
				if err != nil {
					s.logger.Error("error reading temperature",
						zap.String("error", err.Error()),
					)
				}
				// Process the temperature reading
				err = s.callback.Handle(c, f)
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

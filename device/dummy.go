package device

// Contains a dummy implementation of a device for demonstration purposes
// ... and also some simple testing of wiring things together

import (
	"math/rand"
	"time"

	"github.com/nherson/brewski/measurement"
)

func init() {
	rand.Seed(time.Now().Unix())
}

// DummyDevice is a random number picker
type DummyDevice struct {
	name           string
	possibleValues []float32
}

// NewDummyDevice returns a new dummy device that returns one of the provided
// dummy values at random. If no dummy values are provided, uses a preset list
// of dummy values
func NewDummyDevice(name string, possibleValues ...float32) *DummyDevice {
	// have some defaults ready
	if len(possibleValues) == 0 {
		possibleValues = []float32{0.0, 1.1, 2.2, 3.3, 4.4, 5.5}
	}
	return &DummyDevice{
		name:           name,
		possibleValues: possibleValues,
	}
}

// Read returns one of the dummy stored in the dummy devices possible values,
// selecting one at random to return.
func (d *DummyDevice) Read() ([]measurement.Sample, error) {
	t := time.Now()
	sample := measurement.NewDeviceSample(d.Name())
	value := d.possibleValues[rand.Intn(len(d.possibleValues))]
	sample.AddDatapoint("random", value, t)
	return []measurement.Sample{sample}, nil
}

// Name returns the formatted name of the dummy device
func (d *DummyDevice) Name() string {
	return d.name
}

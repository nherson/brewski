package device

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDS18B20Read(t *testing.T) {
	SetOnewireSysfsDir("../testdata/temperature/ds18b20")

	probe := NewDS18B20("28-0123456789abcd", "28-0123456789abcd")

	// There is a dummy file to 'read' data from in this repo.
	// Read data from there and make sure it looks ok
	samples, err := probe.Read()
	assert.Equal(t, 1, len(samples))

	assert.Nil(t, err)

	sample := samples[0]

	assert.NotNil(t, sample)
	assert.Equal(t, 2, len(sample.Datapoints()))
	fahrenheitDatapointExists := false
	celsiusDatapointExists := false
	for _, datapoint := range sample.Datapoints() {
		switch datapoint.Name() {
		case "celsius":
			assert.Equal(t, float32(21.375), datapoint.Value())
			celsiusDatapointExists = true
		case "fahrenheit":
			assert.Equal(t, float32(70.475), datapoint.Value())
			fahrenheitDatapointExists = true
		}
	}
	assert.True(t, celsiusDatapointExists)
	assert.True(t, fahrenheitDatapointExists)
}

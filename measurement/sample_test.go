package measurement

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Simple test for Sample code.
func TestSample(t *testing.T) {
	assert.Nil(t, nil)

	ts1 := time.Now()
	ts2 := ts1.Add(time.Second * 10)

	d1Name := "field1"
	d1Value := float32(1.234)
	d1Time := ts1

	d2Name := "field2"
	d2Value := float32(2.345)
	d2Time := ts2

	// Create a sample
	sample := NewDeviceSample("testDevice")
	sample.AddDatapoint(d1Name, d1Value, d1Time)
	sample.AddDatapoint(d2Name, d2Value, d2Time)

	// Test that the sample looks ok
	assert.Equal(t, "testDevice", sample.DeviceName())
	assert.Equal(t, d1Name, sample.Datapoints()[0].Name())
	assert.Equal(t, d1Value, sample.Datapoints()[0].Value())
	assert.Equal(t, d2Name, sample.Datapoints()[1].Name())
	assert.Equal(t, d2Value, sample.Datapoints()[1].Value())

}

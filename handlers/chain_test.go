package handlers

import (
	"testing"
	"time"

	"github.com/nherson/brewski/measurement"
	"github.com/stretchr/testify/assert"
)

// mockCallback just stashes samples passed to it
type mockCallback struct {
	samples []measurement.Sample
}

func newMockCallback() *mockCallback {
	return &mockCallback{
		samples: []measurement.Sample{},
	}
}

func (mc *mockCallback) Handle(s measurement.Sample) error {
	mc.samples = append(mc.samples, s)
	return nil
}

func (mc *mockCallback) ReceivedSamples() []measurement.Sample {
	return mc.samples
}

func TestChainCallback(t *testing.T) {
	// create a few mocked callbacks
	mock1 := newMockCallback()
	mock2 := newMockCallback()
	mock3 := newMockCallback()

	// Create a chained callback using the 3 mocked callbacks
	callbackChain := NewChainCallback()
	callbackChain.RegisterCallback(mock1)
	callbackChain.RegisterCallback(mock2)
	callbackChain.RegisterCallback(mock3)

	// create a few device samples
	deviceName := "testDevice1"
	sample1 := measurement.NewDeviceSample(deviceName)
	sample1.AddDatapoint("foobar", 1.234, time.Now())
	sample1.AddDatapoint("garply", 2.345, time.Now())
	sample2 := measurement.NewDeviceSample(deviceName)
	sample2.AddDatapoint("foobar", 3.456, time.Now())
	sample2.AddDatapoint("garply", 4.567, time.Now())

	// Send in the samples to the callback chain
	callbackChain.Handle(sample1)
	callbackChain.Handle(sample2)

	// Make sure the samples look good on the other side
	mock1Samples := mock1.ReceivedSamples()
	assert.Equal(t, "testDevice1", mock1Samples[0].DeviceName())
	assert.Equal(t, "foobar", mock1Samples[0].Datapoints()[0].Name())
	assert.Equal(t, float32(1.234), mock1Samples[0].Datapoints()[0].Value())
	assert.Equal(t, "garply", mock1Samples[0].Datapoints()[1].Name())
	assert.Equal(t, float32(2.345), mock1Samples[0].Datapoints()[1].Value())
	assert.Equal(t, "testDevice1", mock1Samples[1].DeviceName())
	assert.Equal(t, "foobar", mock1Samples[1].Datapoints()[0].Name())
	assert.Equal(t, float32(3.456), mock1Samples[1].Datapoints()[0].Value())
	assert.Equal(t, "garply", mock1Samples[1].Datapoints()[1].Name())
	assert.Equal(t, float32(4.567), mock1Samples[1].Datapoints()[1].Value())

	mock2Samples := mock1.ReceivedSamples()
	assert.Equal(t, "testDevice1", mock2Samples[0].DeviceName())
	assert.Equal(t, "foobar", mock2Samples[0].Datapoints()[0].Name())
	assert.Equal(t, float32(1.234), mock2Samples[0].Datapoints()[0].Value())
	assert.Equal(t, "garply", mock2Samples[0].Datapoints()[1].Name())
	assert.Equal(t, float32(2.345), mock2Samples[0].Datapoints()[1].Value())
	assert.Equal(t, "testDevice1", mock2Samples[1].DeviceName())
	assert.Equal(t, "foobar", mock2Samples[1].Datapoints()[0].Name())
	assert.Equal(t, float32(3.456), mock2Samples[1].Datapoints()[0].Value())
	assert.Equal(t, "garply", mock2Samples[1].Datapoints()[1].Name())
	assert.Equal(t, float32(4.567), mock2Samples[1].Datapoints()[1].Value())

	mock3Samples := mock1.ReceivedSamples()
	assert.Equal(t, "testDevice1", mock3Samples[0].DeviceName())
	assert.Equal(t, "foobar", mock3Samples[0].Datapoints()[0].Name())
	assert.Equal(t, float32(1.234), mock3Samples[0].Datapoints()[0].Value())
	assert.Equal(t, "garply", mock3Samples[0].Datapoints()[1].Name())
	assert.Equal(t, float32(2.345), mock3Samples[0].Datapoints()[1].Value())
	assert.Equal(t, "testDevice1", mock3Samples[1].DeviceName())
	assert.Equal(t, "foobar", mock3Samples[1].Datapoints()[0].Name())
	assert.Equal(t, float32(3.456), mock3Samples[1].Datapoints()[0].Value())
	assert.Equal(t, "garply", mock3Samples[1].Datapoints()[1].Name())
	assert.Equal(t, float32(4.567), mock3Samples[1].Datapoints()[1].Value())
}

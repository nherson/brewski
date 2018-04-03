package device

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/nherson/brewski/measurement"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type mockReader struct {
	samples           []measurement.Sample
	sampleReturnCount int
	index             int
}

func newMockReader(samples []measurement.Sample, sampleReturnCount int) *mockReader {
	return &mockReader{
		samples:           samples,
		index:             0,
		sampleReturnCount: sampleReturnCount,
	}
}

func (mr *mockReader) Name() string {
	return "mockReader"
}

func (mr *mockReader) Read() ([]measurement.Sample, error) {
	toReturn := []measurement.Sample{}
	if len(mr.samples)-mr.index < mr.sampleReturnCount {
		return nil, fmt.Errorf("not enough samples to return")
	}
	for i := 0; i < mr.sampleReturnCount; i++ {
		toReturn = append(toReturn, mr.samples[i+mr.index])
	}
	mr.index += mr.sampleReturnCount
	return toReturn, nil
}

type mockCallback struct {
	t               *testing.T
	expectedSamples []measurement.Sample
	index           int
}

func newMockCallback(t *testing.T, samples []measurement.Sample) *mockCallback {
	return &mockCallback{
		t:               t,
		expectedSamples: samples,
		index:           0,
	}
}

func (mc *mockCallback) Handle(actualSample measurement.Sample) error {
	if mc.index == len(mc.expectedSamples) {
		assert.Fail(mc.t, "ran out of expected samples, but Handle() has been called")
	}
	expectedSample := mc.expectedSamples[mc.index]
	// test sample device name equality
	assert.Equal(mc.t, expectedSample.DeviceName(), actualSample.DeviceName())
	// test sample tag equality
	for expectedKey, expectedValue := range expectedSample.Tags() {
		actualValue, found := actualSample.Tags()[expectedKey]
		if !found {
			assert.Fail(mc.t, fmt.Sprintf("missing expected key '%s'", expectedKey))
		}
		assert.Equal(mc.t, expectedValue, actualValue)
	}
	// test sample datapoint equality
	assert.Equal(mc.t, len(expectedSample.Datapoints()), len(actualSample.Datapoints()))
	for i, expectedDatapoint := range expectedSample.Datapoints() {
		assert.Equal(mc.t, expectedDatapoint.Name(), actualSample.Datapoints()[i].Name())
		assert.Equal(mc.t, expectedDatapoint.Value(), actualSample.Datapoints()[i].Value())
		assert.Equal(mc.t, expectedDatapoint.Time(), actualSample.Datapoints()[i].Time())
	}
	mc.index++
	return nil
}

// taken from https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
func randomString(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")
	s := make([]rune, n)
	for i := range s {
		s[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(s)
}

func generateTestSamples(deviceName string, numSamples int) []measurement.Sample {
	t := time.Now()
	samples := []measurement.Sample{}
	for i := 0; i < numSamples; i++ {
		sampleTime := t.Add(time.Duration(i) * time.Second)
		sample := measurement.NewDeviceSample(deviceName)
		// add a random number of tags
		for j := 0; j < rand.Intn(5); j++ {
			sample.AddTag(randomString(5), randomString(5))
		}
		// add a random number of datapoints
		for k := 0; k < rand.Intn(5); k++ {
			sample.AddDatapoint(randomString(5), rand.Float32(), sampleTime)
		}
		samples = append(samples, sample)
	}
	return samples
}

// The callback for this test makes sure that the sample received is what is expected.
// This is done by pre-generating a batch of samples for the "device" to emit, which is
// handed to both the "device" and the mock callback handler.  When the handler receives
// a sample, it makes sure all the fields are the same as what is expected in its personal
// copy of the sample.
// This is essentially an integration test for the Sensor implementation.
func TestSensor(t *testing.T) {
	assert.True(t, true)
	testSamples := generateTestSamples("testDevice", 100)
	mr := newMockReader(testSamples, 5)
	mc := newMockCallback(t, testSamples)
	logger, _ := zap.NewProduction()
	sensor := NewSensor(mr, time.Millisecond, logger)
	sensor.SetCallback(mc)
	sensor.Start()
	time.Sleep(5 * time.Second)
	sensor.Stop()
}

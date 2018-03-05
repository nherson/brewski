package measurement

import "time"

// Sample holds information about datapoints retrieved from a device
type Sample interface {
	DeviceName() string
	Datapoints() []Datapoint
	AddDatapoint(string, float32, time.Time)
}

// DeviceSample implements the Sample interface to hold data from a
// device sensor reading
type DeviceSample struct {
	datapoints []Datapoint
	deviceName string
}

// NewDeviceSample returns a device sample that contains no
// datapoints (yet)
func NewDeviceSample(deviceName string) *DeviceSample {
	return &DeviceSample{
		deviceName: deviceName,
		datapoints: []Datapoint{},
	}
}

// DeviceName returns the name of the device which this sample belongs to
func (ds *DeviceSample) DeviceName() string {
	return ds.deviceName
}

// Datapoints returns the datapoints that have been added to this sample
func (ds *DeviceSample) Datapoints() []Datapoint {
	return ds.datapoints
}

// AddDatapoint adds a datapoint to this DeviceSample
func (ds *DeviceSample) AddDatapoint(name string, value float32, time time.Time) {
	d := newDatapoint(name, value, time)
	ds.datapoints = append(ds.datapoints, d)
}

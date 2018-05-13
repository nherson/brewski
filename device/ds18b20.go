package device

// Contains a Sensor implementation for the DS18B20 temperature probe

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/nherson/brewski/measurement"
)

// OnewireSysfsDir is the directory where onewire sysfs device interfaces
// are mounted on the filesystem. Below is a default but it can be changed
// using the SetOnewireSysfsDir function
var OnewireSysfsDir = "/sys/bus/w1/devices"

// SetOnewireSysfsDir changes a global var indicating where on the OS
// the sysfs onewire interface can be found.
func SetOnewireSysfsDir(dir string) {
	OnewireSysfsDir = dir
}

// DS18B20 represents a temperature sensor
type DS18B20 struct {
	name string
	ID   string
}

// NewDS18B20 creates a new sensor struct. Once created, the sensor can be
// launched using Start(), where it will perioically poll data from the sensor
// and feed it into a callback function
func NewDS18B20(name string, deviceID string) *DS18B20 {
	return &DS18B20{
		name: name,
		ID:   deviceID,
	}
}

// ReadTemperature reads data from the sensor and saves the temperature reading (celsius)
// into the sensors lastReadTemperatureC field
// if there is a problem reading data an error is returned
func (d *DS18B20) Read() ([]measurement.Sample, error) {
	dataFile := d.getSysfsPath()
	b, err := ioutil.ReadFile(dataFile)
	if err != nil {
		return nil, err
	}
	// simple sanity checks
	// 1. make sure the output is exactly 2 lines
	lines := strings.Split(string(b), "\n")
	if len(lines) != 2 {
		return nil, fmt.Errorf("unexpected number of lines in sensor output: %s", string(len(lines)))
	}
	// 2. Make sure the device is ready to read (first line ends in 'YES')
	if !d.isReady(lines[0]) {
		return nil, fmt.Errorf("device not ready")
	}
	return d.parseTemperature(lines[1])
}

// Returns a celsius and fahrenheit reading given the data line for a sensor
// reading. Returns an error if there is an issue parsing the line
func (d *DS18B20) parseTemperature(dataLine string) ([]measurement.Sample, error) {
	// Take the last field and split on '=' (we expect a format `t=$temp`)
	fields := strings.Split(dataLine, " ")
	thermReading := strings.Split(fields[len(fields)-1], "=")

	// Expect the above split to return 2 fields (the `t` and the `$temp`)
	if len(thermReading) != 2 {
		return nil, fmt.Errorf("unknown error reading temperature from sensor")
	}
	// attempt to parse into a float
	// the actual data value is always a signed int, but lets parse right
	// into a float because we will immediately be dividing by 1000
	rawReading, err := strconv.ParseFloat(thermReading[1], 32)
	if err != nil {
		return nil, fmt.Errorf("error parsing temperature: %s", err.Error())
	}
	t := time.Now()
	// create a sample using the retrieved data
	sample := measurement.NewDeviceSample(d.Name())
	// Add the device ID to the sample as a tag
	sample.AddTag("id", d.ID)
	// convert from millicelsius(?) to celsius
	c := float32(rawReading) / 1000
	// add the celsius reading to the sample
	sample.AddDatapoint("celsius", c, t)
	// convert from celsius to fahrenheit
	f := celsiusToFahrenheit(c)
	sample.AddDatapoint("fahrenheit", f, t)

	// return the results
	// the Reader interface expects an array of samples, but this device is very simple,
	// so just return the single sample alone in an array... it just has the temp
	return []measurement.Sample{sample}, nil

}

// Name returns the formatted name of the sensor, in the format
// DS18B20-<device id>
func (d *DS18B20) Name() string {
	return d.name
}

// checks if the reading indicates that the sensor data
// is good to use. returns true if so, otherwise false
func (d *DS18B20) isReady(s string) bool {
	fields := strings.Split(s, " ")
	n := len(fields)
	return fields[n-1] == "YES"
}

// utility to convert celsium to fahrenheit
func celsiusToFahrenheit(c float32) float32 {
	return (c * 9 / 5) + 32
}

// utility to get a filepath to read sensor data
// uses ONEWIRE_SYSFS_DIR and the sensor's ID
// to build a filepath to read sensor data from
func (d *DS18B20) getSysfsPath() string {
	return filepath.Join(OnewireSysfsDir, d.ID, "w1_slave")
}

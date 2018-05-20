package config

import (
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/BurntSushi/toml"
	"github.com/nherson/brewski/device"
	"github.com/nherson/brewski/outputs"
)

// Contains structs and definitions to turn a TOML config file
// into a ready-to-roll brewski orchestration

// Config is a top level struct that contains all configuration data.
// The TOML file will be parsed into this struct
type Config struct {
	Global  GlobalConfig  `toml:"global"`
	Devices DevicesConfig `toml:"devices"`
	Outputs OutputsConfig `toml:"outputs"`
}

// ParseConfig returns a Config struct generated from the received bytes,
// otherwise returns an error
func ParseConfig(b []byte) (*Config, error) {
	var conf Config
	if _, err := toml.Decode(string(b), &conf); err != nil {
		return nil, err
	}
	return &conf, nil
}

// GlobalConfig contains configuration that is either a global
// default for devices and outputs, or is some other globally
// applied configuration
type GlobalConfig struct {
	PollingInterval duration `toml:"polling-interval"`
	OnesireSysfsDir string   `toml:"onewire-sysfs-dir"`
}

// DevicesConfig holds configuration data for each device being setup for use
type DevicesConfig struct {
	DS18B20s     map[string]*DS18B20Config     `toml:"ds18b20"`
	Tilts        map[string]*TiltConfig        `toml:"tilt"`
	DummyDevices map[string]*DummyDeviceConfig `toml:"dummy-device"`
}

// AllDeviceConfigs returns a mapping from device names to their configuration
func (d *DevicesConfig) AllDeviceConfigs() (map[string]DeviceConfig, error) {
	deviceConfigs := make(map[string]DeviceConfig)
	for name, deviceConfig := range d.DS18B20s {
		if _, found := deviceConfigs[name]; found {
			return nil, fmt.Errorf("duplicate device declared '%s'", name)
		}
		deviceConfigs[name] = deviceConfig
	}
	for name, deviceConfig := range d.Tilts {
		if _, found := deviceConfigs[name]; found {
			return nil, fmt.Errorf("duplicate device declared '%s'", name)
		}
		deviceConfigs[name] = deviceConfig
	}
	for name, deviceConfig := range d.DummyDevices {
		if _, found := deviceConfigs[name]; found {
			return nil, fmt.Errorf("duplicate device declared '%s'", name)
		}
		deviceConfigs[name] = deviceConfig
	}
	return deviceConfigs, nil
}

// OutputsConfig holds configuration data for each output being setup for use
type OutputsConfig struct {
	Logs      map[string]*LogConfig      `toml:"log"`
	Influxdbs map[string]*InfluxdbConfig `toml:"influxdb"`
}

// AllOutputConfigs returns a mapping between an output name and its OutputConfig
func (d *OutputsConfig) AllOutputConfigs() (map[string]OutputConfig, error) {
	outputConfigs := make(map[string]OutputConfig)
	for name, outputConfig := range d.Logs {
		if _, found := outputConfigs[name]; found {
			return nil, fmt.Errorf("duplicate output declared '%s'", name)
		}
		outputConfigs[name] = outputConfig
	}
	for name, outputConfig := range d.Influxdbs {
		if _, found := outputConfigs[name]; found {
			return nil, fmt.Errorf("duplicate output declared '%s'", name)
		}
		outputConfigs[name] = outputConfig
	}
	return outputConfigs, nil
}

// OutputConfig is some configuration for an output
// that can be used to generate a corresponding outputs.Callback
type OutputConfig interface {
	GenerateOutput() (outputs.Callback, error)
}

// DeviceConfig is some configuration for a device that
// can be used to generate a corresponding device.Reader
type DeviceConfig interface {
	GenerateDevice(string) (device.Reader, error)
	OutputNames() []string
}

// DEVICE CONFIG STRUCTS

// DS18B20Config holds configuration data about a ds18b20 temperature sensor
type DS18B20Config struct {
	ID      string   `toml:"id"`
	Outputs []string `toml:"outputs"`
}

// GenerateDevice creates a DS18B20 device from a given configuration
// Validates that an ID is given
func (c *DS18B20Config) GenerateDevice(name string) (device.Reader, error) {
	if c.ID == "" {
		return nil, fmt.Errorf("ds18b20 id cannot be empty")
	}
	return device.NewDS18B20(name, c.ID), nil
}

// OutputNames returns the names of the outputs configured for this
// DS18B20 device
func (c *DS18B20Config) OutputNames() []string {
	return c.Outputs
}

// TiltConfig holds configuration data about a fleet of Tilt Hydrometers (all colors)
type TiltConfig struct {
	TemperatureCalibration float32  `toml:"temperature_calibration"` // defaults to 0
	GravityCalibration     float32  `toml:"gravity_calibration"`     // defaults to 0
	Outputs                []string `toml:"outputs"`
}

// GenerateDevice creates a TiltHydrometer device from a given configuration
func (c *TiltConfig) GenerateDevice(name string) (device.Reader, error) {
	return device.NewTiltHydrometer(name, c.GravityCalibration, c.TemperatureCalibration)
}

// OutputNames returns the names of the outputs configured for this
// tilt hydrometer configuration
func (c *TiltConfig) OutputNames() []string {
	return c.Outputs
}

// DummyDeviceConfig holds configuration data about a DummyDevice
type DummyDeviceConfig struct {
	PossibleValues []float32 `toml:"possible-values"`
	Outputs        []string  `toml:"outputs"`
}

// GenerateDevice creates a DummyDevice from a configuration
func (c *DummyDeviceConfig) GenerateDevice(name string) (device.Reader, error) {
	return device.NewDummyDevice(name, c.PossibleValues...), nil
}

// OutputNames returns the names of the outputs configured for this
// dummy device configuration
func (c *DummyDeviceConfig) OutputNames() []string {
	return c.Outputs
}

// OUTPUT CONFIG STRUCTS

// LogConfig holds configuration data about a logger (using zap)
type LogConfig struct {
}

// GenerateOutput creates a LoggingCallback output from a given configuration
func (c *LogConfig) GenerateOutput() (outputs.Callback, error) {
	l, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	return outputs.NewLoggingCallback(l), nil
}

// InfluxdbConfig holds configuration data for an influxdb database
type InfluxdbConfig struct {
	Address  string `toml:"address"`
	Database string `toml:"database"`
}

// GenerateOutput creates an InfluxdbCallback output from a given configuration
func (c *InfluxdbConfig) GenerateOutput() (outputs.Callback, error) {
	if c.Address == "" {
		return nil, fmt.Errorf("address must be provided for influxdb output")
	} else if c.Database == "" {
		return nil, fmt.Errorf("database must be provided for influxdb output")
	}
	return outputs.NewInfluxDBCallback(c.Address, c.Database)
}

// HELPERS

// from the README at https://github.com/BurntSushi/toml
type duration struct {
	time.Duration
}

// UnmarshalText allows decoding duration strings into time.Duration objects
func (d *duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

package config

import (
	"fmt"

	"github.com/nherson/brewski/device"
	"github.com/nherson/brewski/outputs"
	"go.uber.org/zap"
)

// Generate uses a parsed config to create a list of Sensors
// which will be the main set of input-output pipelines used
// to collect data and operate on it.
func (c *Config) Generate() ([]device.Poller, error) {
	// Get some global config options
	pollingInterval := c.Global.PollingInterval.Duration

	// Get raw device configs for looking up device<-->outputs mappings
	deviceConfigs, err := c.Devices.AllDeviceConfigs()
	if err != nil {
		return nil, err
	}

	outputConfigs, err := c.Outputs.AllOutputConfigs()
	if err != nil {
		return nil, err
	}

	// A place to store outputs that have already been generated
	generatedOutputs := make(map[string]outputs.Callback)

	pollers := []device.Poller{}

	// Iterate over each device, generating (or pulling from cache) all
	// configured outputs to associate with it. Generate one sensor per device
	for deviceName, deviceConfig := range deviceConfigs {
		// Generate the device
		d, err := deviceConfig.GenerateDevice(deviceName)
		if err != nil {
			return nil, err
		}

		// Create a sensor harness for the device
		// TODO: better logging handling...
		sensorLogger, err := zap.NewProduction()
		if err != nil {
			return nil, err
		}
		sensor := device.NewSensor(d, pollingInterval, sensorLogger)

		// create a chained callback for the sensor
		callbackChain := outputs.NewChainCallback()
		for _, outputName := range deviceConfig.OutputNames() {
			var output outputs.Callback
			var outputConfig OutputConfig
			var found bool
			var err error
			output, found = generatedOutputs[outputName]
			// See if the output has already been generated and cached
			if !found {
				// If not found, generate it
				outputConfig, found = outputConfigs[outputName]
				// If the requested output doesn't exist, that's a config error
				if !found {
					return nil, fmt.Errorf("output '%s' does not exist for device '%s'", outputName, deviceName)
				}
				// Generate the output
				output, err = outputConfig.GenerateOutput()
				if err != nil {
					return nil, err
				}
				// Cache generated output for later
				generatedOutputs[outputName] = output
			}
			// Register the output in the callback chain
			callbackChain.RegisterCallback(output)
		}
		// Assign the callback chain to the sensor
		sensor.SetCallback(callbackChain)
		// Append sensor to list of returned sensors
		pollers = append(pollers, sensor)
	}
	return pollers, nil
}

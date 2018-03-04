package main

import (
	"flag"
	"fmt"
	"time"
)

// Command line flags here
var ds18b20s stringSet
var onewireSysfsDir string
var pollingInterval time.Duration
var influxDBEndpoint string
var influxDBDatabase string
var useInfluxDB bool

type stringSet []string

func (ss *stringSet) Set(v string) error {
	*ss = append(*ss, v)
	return nil
}

func (ss *stringSet) String() string {
	return fmt.Sprint(*ss)
}

func init() {
	flag.Var(&ds18b20s, "ds18b20", "device ID of ds18b20 temperature sensor to probe; can be specified multiple times")
	flag.StringVar(&onewireSysfsDir, "onewire-sysfs-dir", "/sys/bus/w1/devices", "Directory to read onewire devices via sysfs")
	flag.DurationVar(&pollingInterval, "polling-interval", time.Second*5, "Interval to wait between device reads")
	flag.StringVar(&influxDBEndpoint, "influxdb-endpoint", "http://localhost:8086", "Endpoint to send influxdb metrics to")
	flag.StringVar(&influxDBDatabase, "influxdb-database", "brewski", "Endpoint to send influxdb metrics to")
	flag.BoolVar(&useInfluxDB, "influxdb-enabled", false, "whether or not to send metrics to influxdb")
}

// Parse the command lines flags and go through and set defaults
// where necessary
func parseFlags() {
	flag.Parse()
}

// functions to pull out the parsed flag data in a clean set of calls

// GetDS18B20IDs returns a list of DS18B20 temperature probe IDs configured for use
func GetDS18B20IDs() []string {
	return []string(ds18b20s)
}

// GetOnewireSysfsDir returns the kernel sysfs dir to use to read onewire
// temperature data
func GetOnewireSysfsDir() string {
	return onewireSysfsDir
}

// GetPollingInterval returns the interval to wait between probe polling instances
func GetPollingInterval() time.Duration {
	return pollingInterval
}

// GetInfluxDBEndpoint returns the URL at which to send sensor data
func GetInfluxDBEndpoint() string {
	return influxDBEndpoint
}

// GetInfluxDBDatabase returns the database where sensor data should be stored
func GetInfluxDBDatabase() string {
	return influxDBDatabase
}

// UseInfluxDB indicates whether or not sensor samples should be sent to an InfluxDB datastore
func UseInfluxDB() bool {
	return useInfluxDB
}

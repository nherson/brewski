package config

import (
	"fmt"
	"testing"

	"github.com/nherson/brewski/outputs"
	"github.com/stretchr/testify/assert"
)

// A big integration test for various bad configs
func TestBadConfig(t *testing.T) {

	var badConfigText string
	var c *Config
	var err error

	// A simple sanity test against badly formatted TOML
	badConfigText = `
		[global]
		[[[][][]][][]]
		{}{}{}
		123456
		\\\\|||||/////
	`
	c, err = ParseConfig([]byte(badConfigText))
	assert.Nil(t, c)
	assert.NotNil(t, err)

	// test duplicate device names within the same device family
	badConfigText = `
		[devices.dummy-device.foobar]
		possible-values = [2.0,3.0,4.0,5.0]

		[devices.dummy-device.foobar]
`
	c, err = ParseConfig([]byte(badConfigText))
	fmt.Println(err.Error())
	assert.Nil(t, c)
	assert.NotNil(t, err)

	// test duplicate names across device types
	badConfigText = `
	[devices.dummy-device.foobar]
	possible-values = [2.0,3.0,4.0,5.0]

	[devices.ds18b20.foobar]
	`
	c, err = ParseConfig([]byte(badConfigText))
	assert.Nil(t, err)
	_, err = c.Generate()
	assert.NotNil(t, err)

	// test duplicate outputs
	badConfigText = `
	[outputs.log.foobar]

	[outputs.log.foobar]
	`
	c, err = ParseConfig([]byte(badConfigText))
	assert.Nil(t, c)
	assert.NotNil(t, err)

	// test duplicate outputs across output types
	badConfigText = `
	[outputs.log.foobar]

	[outputs.influxdb.foobar]
	`
	c, err = ParseConfig([]byte(badConfigText))
	assert.Nil(t, err)
	_, err = c.Generate()
	assert.NotNil(t, err)
}

func TestDS18B20Config(t *testing.T) {
	var err error
	goodConfig := &DS18B20Config{
		ID: "abc123",
	}
	d1, err := goodConfig.GenerateDevice("name")
	assert.Nil(t, err)
	assert.Equal(t, "name", d1.Name())

	badConfig := &DS18B20Config{}
	d2, err := badConfig.GenerateDevice("name2")
	assert.Nil(t, d2)
	assert.NotNil(t, err)
}

func TestInfluxDBConfig(t *testing.T) {
	var err error
	var o outputs.Callback
	goodConfig := &InfluxdbConfig{
		Address:  "http://example.com:5555/endpoint",
		Database: "mycooldatabase",
	}
	o, err = goodConfig.GenerateOutput()
	assert.Nil(t, err)
	assert.NotNil(t, o)

	// No database provided
	badConfig1 := &InfluxdbConfig{
		Address: "http://example.com:5555/endpoint",
	}
	o, err = badConfig1.GenerateOutput()
	assert.Nil(t, o)
	assert.NotNil(t, err)

	// No address provided
	badConfig2 := &InfluxdbConfig{
		Database: "mycooldatabase",
	}
	o, err = badConfig2.GenerateOutput()
	assert.Nil(t, o)
	assert.NotNil(t, err)
}

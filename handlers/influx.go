package handlers

import (
	"time"

	client "github.com/influxdata/influxdb/client/v2"
	"github.com/nherson/brewski/measurement"
)

// InfluxDBCallback sends sensor data to a specified InfluxDB endpoint
type InfluxDBCallback struct {
	c  client.Client
	db string
}

// NewInfluxDBCallback returns an InfluxDBCallback that can be used to send sensor data to the configured
// influxDB endpoint.  The passed in tags will be included in every callback invocation
func NewInfluxDBCallback(addr string, database string) (*InfluxDBCallback, error) {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: "http://localhost:8086",
	})
	if err != nil {
		return nil, err
	}
	return &InfluxDBCallback{
		c:  c,
		db: database,
	}, nil
}

// TODO: integrate measurement.Sample.DeviceName() into the influxDB tags,
// so it doesnt need to be passed into the Callback constructor

// Handle sends sensor data to the configured InfluxDB database
func (icb *InfluxDBCallback) Handle(s measurement.Sample) error {
	// create a new batch of points
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  icb.db,
		Precision: "s",
	})
	if err != nil {
		return err
	}
	// Turn the measured sample into InfluxDB fields
	fields := map[string]interface{}{}
	for _, d := range s.Datapoints() {
		fields[d.Name()] = d.Value()
	}
	pt, err := client.NewPoint(s.DeviceName(), s.Tags(), fields, time.Now())
	if err != nil {
		return err
	}
	bp.AddPoint(pt)
	err = icb.c.Write(bp)
	if err != nil {
		return err
	}
	return nil
}

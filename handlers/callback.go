package handlers

import "github.com/nherson/brewski/measurement"

// Callback is an interface which allows taking arbitrary datapoints
// and doing some processing on them
type Callback interface {
	Handle(measurement.Sample) error
}

package device

import "github.com/nherson/brewski/measurement"

// Reader is an interface that can read temperature data from
// an arbitrary source (currently DS1820B implemented via sysfs)
type Reader interface {
	Read() (measurement.Sample, error)
	Name() string
}

package measurement

import "time"

// Datapoint holds a discrete datapoint with a single label
// indicating what the value represents.  To be used embedded
// in Sample
type Datapoint interface {
	Name() string
	Value() float32
	Time() time.Time
}

type datapoint struct {
	name  string
	value float32
	time  time.Time
}

// NewDatapoint returns a datapoint for the given data
func newDatapoint(n string, v float32, t time.Time) *datapoint {
	return &datapoint{
		name:  n,
		value: v,
		time:  t,
	}
}

func (d *datapoint) Name() string {
	return d.name
}

func (d *datapoint) Value() float32 {
	return d.value
}

func (d *datapoint) Time() time.Time {
	return d.time
}

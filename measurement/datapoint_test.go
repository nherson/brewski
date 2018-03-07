package measurement

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Simple test for Datapoint code.
func TestDatapoint(t *testing.T) {
	assert.Nil(t, nil)

	ts := time.Now()
	d := newDatapoint("fieldName", 1.234, ts)

	assert.Equal(t, "fieldName", d.Name())
	assert.Equal(t, float32(1.234), d.Value())
	assert.Equal(t, ts, d.Time())
}

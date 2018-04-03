package device

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecentData(t *testing.T) {
	data := newRecentData()
	for _, d := range data {
		assert.False(t, d.seen)
	}
	assert.Nil(t, nil)
	// add a few datapoints
	zero := float32(0)
	data["red"].addData(zero, zero)
	data["red"].addData(zero, zero)
	data["red"].addData(zero, zero)
	data["yellow"].addData(zero, zero)
	data["yellow"].addData(zero, zero)
	data["yellow"].addData(zero, zero)
	data["yellow"].addData(zero, zero)
	data["pink"].addData(zero, zero)

	assert.Equal(t, 3, data["red"].count)
	assert.Equal(t, 4, data["yellow"].count)
	assert.Equal(t, 1, data["pink"].count)
	assert.Equal(t, 0, data["black"].count)
	assert.Equal(t, 0, data["purple"].count)
}

func TestAverageData(t *testing.T) {
	ad := &averagedData{
		gravity:     float32(1),
		temperature: float32(0),
		count:       0,
		seen:        false,
	}
	ad.addData(float32(0), float32(4.5))
	assert.True(t, ad.seen)
	assert.Equal(t, float32(0), ad.gravity)
	assert.Equal(t, float32(4.5), ad.temperature)

	ad.addData(float32(0.9), float32(9))
	assert.Equal(t, float32(0.45), ad.gravity)
	assert.Equal(t, float32(6.75), ad.temperature)

	ad.addData(float32(2.43), float32(7.5))
	assert.Equal(t, float32(1.11), ad.gravity)
	assert.Equal(t, float32(7), ad.temperature)

	assert.Equal(t, 3, ad.count)
}

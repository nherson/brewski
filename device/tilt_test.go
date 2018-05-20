package device

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"testing"

	"github.com/go-ble/ble"
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

// returns one value at a time when requested, until there are none left
type mockBluetooth struct {
	tempToReturn    uint16
	gravityToReturn uint16
}

func (mb *mockBluetooth) GetAdvertisements() []ble.Advertisement {
	return []ble.Advertisement{
		newMockAdvertisement(mb.tempToReturn, mb.gravityToReturn),
	}
}

type mockAdvertisement struct {
	md []byte
}

func newMockAdvertisement(temp, gravity uint16) *mockAdvertisement {
	var buffer bytes.Buffer
	pre := "4C000215A495BB10C5B14B44B5121370F02D74DE"
	post := "C7"
	preBytes, err := hex.DecodeString(pre)
	if err != nil {
		panic(err)
	}
	postBytes, err := hex.DecodeString(post)
	if err != nil {
		panic(err)
	}
	buffer.Write(preBytes)
	binary.Write(&buffer, binary.BigEndian, temp)
	binary.Write(&buffer, binary.BigEndian, gravity)
	buffer.Write(postBytes)
	return &mockAdvertisement{
		md: buffer.Bytes(),
	}
}

func (ma *mockAdvertisement) ManufacturerData() []byte {
	return ma.md
}
func (ma *mockAdvertisement) LocalName() string              { return "" }
func (ma *mockAdvertisement) ServiceData() []ble.ServiceData { return nil }
func (ma *mockAdvertisement) Services() []ble.UUID           { return nil }
func (ma *mockAdvertisement) OverflowService() []ble.UUID    { return nil }
func (ma *mockAdvertisement) TxPowerLevel() int              { return 0 }
func (ma *mockAdvertisement) Connectable() bool              { return false }
func (ma *mockAdvertisement) SolicitedService() []ble.UUID   { return nil }
func (ma *mockAdvertisement) RSSI() int                      { return 0 }
func (ma *mockAdvertisement) Addr() ble.Addr                 { return nil }

func TestCalibration(t *testing.T) {
	mb := &mockBluetooth{
		gravityToReturn: uint16(1040),
		tempToReturn:    uint16(67),
	}

	tilt := &TiltHydrometer{
		name:               "test-tilt",
		bluetooth:          mb,
		data:               newRecentData(),
		gravityCalibration: 0.003,
		tempCalibration:    -2,
	}

	samples, err := tilt.Read()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(samples))
	datapoints := samples[0].Datapoints()
	assert.Equal(t, 2, len(datapoints))
	// temp calibrated down 2 points
	assert.Equal(t, float32(65), datapoints[0].Value())
	assert.Equal(t, float32(1.043), datapoints[1].Value())

}

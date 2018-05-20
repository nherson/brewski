package device

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/go-ble/ble"
	"github.com/go-ble/ble/examples/lib/dev"

	"github.com/nherson/brewski/measurement"
)

// Mapping between hardcoded Tilt Hydrometer UUIDs
// and the color that the UUID represents
var colorUUIDMap = map[string]string{
	"A495BB10C5B14B44B5121370F02D74DE": "red",
	"A495BB20C5B14B44B5121370F02D74DE": "green",
	"A495BB30C5B14B44B5121370F02D74DE": "black",
	"A495BB40C5B14B44B5121370F02D74DE": "purple",
	"A495BB50C5B14B44B5121370F02D74DE": "orange",
	"A495BB60C5B14B44B5121370F02D74DE": "blue",
	"A495BB70C5B14B44B5121370F02D74DE": "yellow",
	"A495BB80C5B14B44B5121370F02D74DE": "pink",
}

// TiltHydrometer tracks hydrometer and temperature data from
// Tilt Hydometers.  This single device instance will track every
// device color
type TiltHydrometer struct {
	name               string
	bluetooth          BluetoothScanner
	data               recentData
	tempCalibration    float32
	gravityCalibration float32
}

// NewTiltHydrometer returns a new device capable of reading from a Tilt Hydrometer
func NewTiltHydrometer(name string, gravityCalibration, tempCalibration float32) (*TiltHydrometer, error) {
	// This will immediately start scanning and holding on to discovered advertisements
	b, err := newBluetoothScanner()
	if err != nil {
		return nil, err
	}
	return &TiltHydrometer{
		name:               name,
		bluetooth:          b,
		data:               newRecentData(),
		gravityCalibration: gravityCalibration,
		tempCalibration:    tempCalibration,
	}, nil
}

type recentData map[string]*averagedData

func newRecentData() recentData {
	rd := make(recentData)
	for _, color := range colorUUIDMap {
		rd[color] = &averagedData{
			gravity:     float32(1),
			temperature: float32(0),
			count:       0,
			seen:        false,
		}
	}
	return rd
}

// just reset all our counts to zero, but hold on to the last known values in case
// no new data is received before another read request (so the last value will be returned)
func (rd recentData) clearRecentData() {
	for _, avgData := range rd {
		avgData.count = 0
	}
}

type averagedData struct {
	gravity     float32
	temperature float32
	count       int
	seen        bool
}

// function to incorporate a new data point into this tilt's recent data set, using an averaging method
func (ad *averagedData) addData(gravity, temperature float32) {
	if ad.count == 0 {
		ad.gravity = gravity
		ad.temperature = temperature
		ad.count = 1
		ad.seen = true
	} else {
		ad.gravity = (ad.gravity*float32(ad.count) + gravity) / float32(ad.count+1)
		ad.temperature = (ad.temperature*float32(ad.count) + temperature) / float32(ad.count+1)
		ad.count++
	}
}

// Read reads any data that has been advertised by the tilt since last Read()
// and returns that. If no advertisements have been made by any tilt devices since
// last call to Read, this will return an empty sample with no Datapoints
func (th *TiltHydrometer) Read() ([]measurement.Sample, error) {
	t := time.Now()
	advertisements := th.bluetooth.GetAdvertisements()
	samples := []measurement.Sample{}
	// Process all the newly received advertisements in preparation
	// for returning the most accurate readings possible.
	// If we have received multiple advertisements for the same device
	// since the last Read(), the datapoints will be averaged together
	for _, a := range advertisements {
		var color string
		var isTilt bool
		if color, isTilt = isTiltHydrometer(a); !isTilt {
			continue
		}
		// parse out values for temp and gravity
		temp, gravity := parseTiltData(a)
		th.data[color].addData(gravity, temp)
	}
	// Collect all the tilt data gathered and report it into the sample to return
	for color, data := range th.data {
		// don't report devices we have never even seen
		if !data.seen {
			continue
		}
		sample := measurement.NewDeviceSample("tilt")
		sample.AddTag("color", color)
		sample.AddDatapoint("temperature", data.temperature+th.tempCalibration, t)
		sample.AddDatapoint("gravity", data.gravity+th.gravityCalibration, t)
		samples = append(samples, sample)
	}
	// Clear the recent data counts to prepare for the next read window
	th.data.clearRecentData()
	return samples, nil
}

func parseTiltData(a ble.Advertisement) (float32, float32) {
	md := a.ManufacturerData()
	tempBytes := md[20:22]
	gravityBytes := md[22:24]
	temp := float32(binary.BigEndian.Uint16(tempBytes))
	gravity := float32(binary.BigEndian.Uint16(gravityBytes)) / 1000
	return temp, gravity

}

// determines if the bluetooth advertisement belongs to a tilt hydrometer
// returns the color if so, otherwise returns false for the bool
func isTiltHydrometer(a ble.Advertisement) (string, bool) {
	md := a.ManufacturerData()
	if len(md) != 25 {
		return "", false
	}
	var color string
	// Parse out the constant preamble and the UUID
	discoveredPreamble := md[0:4]
	requiredPreamble := []byte{0x4c, 0x00, 0x02, 0x15}
	uuid := md[4:20]
	if bytes.Compare(discoveredPreamble, requiredPreamble) != 0 {
		return "", false
	}
	discoveredUUID := strings.ToUpper(hex.EncodeToString(uuid))
	// Do a lookup to determine the color of the found device.
	// If the UUID is not in the map, then the device must not
	// be a Tilt Hydrometer
	color, found := colorUUIDMap[discoveredUUID]
	if !found {
		// unknown UUID, so skip
		return "", false
	}
	return color, true
}

// Name returns the name of this device (tilt)
func (th *TiltHydrometer) Name() string {
	return th.name
}

// BluetoothScanner is an interface to read data from a Bluetooth LE device
type BluetoothScanner interface {
	GetAdvertisements() []ble.Advertisement
}

type bluetoothScanner struct {
	advertisements []ble.Advertisement
	lock           *sync.Mutex
}

func newBluetoothScanner() (*bluetoothScanner, error) {
	// Setup bluetooth scanning device
	d, err := dev.NewDevice("default")
	if err != nil {
		return nil, err
	}
	ble.SetDefaultDevice(d)
	// Scan for devices with no timeout
	bs := &bluetoothScanner{
		lock:           &sync.Mutex{},
		advertisements: []ble.Advertisement{},
	}
	// the bluetoothScanner will blindly start gobbling up discovered advertisements
	// and cache them away for consumption by callers of GetAdvertisements
	go func() {
		for {
			// Just keep doing short term scans
			// Use a refreshing scan with a timeout in an infinite loop because there are issues
			// with long running calls to Scan() no longer seeing data
			ctx := ble.WithSigHandler(context.WithTimeout(context.Background(), 1*time.Minute))
			ble.Scan(ctx, true, bs.advHandler, nil)
		}
	}()
	return bs, nil
}

func (bs *bluetoothScanner) GetAdvertisements() []ble.Advertisement {
	bs.lock.Lock()
	defer bs.lock.Unlock()
	// make the array to return
	toReturn := make([]ble.Advertisement, len(bs.advertisements))
	copy(toReturn, bs.advertisements)
	bs.advertisements = []ble.Advertisement{}
	return toReturn
}

// this will be a noop since we can handle all this logic
func (bs *bluetoothScanner) advHandler(a ble.Advertisement) {
	// First, verify that this could possibly be a Tilt
	// By ensuring that the message payload is 25 bytes (iBeacon), and
	// the manufacturer data preamble is the right hex string (0x4c000215)

	// Add the advertisement to the pool
	bs.lock.Lock()
	bs.advertisements = append(bs.advertisements, a)
	bs.lock.Unlock()
}

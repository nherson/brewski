package main

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"strings"

	"github.com/go-ble/ble"
	"github.com/go-ble/ble/examples/lib/dev"
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

func main() {
	fmt.Println("Scanning for Tilt Hydrometer devices, use ctrl-C to exit...")

	// Setup bluetooth scanning device
	d, err := dev.NewDevice("default")
	if err != nil {
		log.Fatalf("can't new device : %s", err)
	}
	ble.SetDefaultDevice(d)

	// Scan for devices with no timeout
	ble.Scan(context.Background(), true, advHandler, nil)
}

// this will be a noop since we can handle all this logic
func advHandler(a ble.Advertisement) {
	// First, verify that this could possibly be a Tilt
	// By ensuring that the message payload is 25 bytes (iBeacon), and
	// the manufacturer data preamble is the right hex string (0x4c000215)
	md := a.ManufacturerData()
	if len(md) != 25 {
		return
	}
	var color string
	// Parse out the constant preamble and the UUID
	discoveredPreamble := md[0:4]
	requiredPreamble := []byte{0x4c, 0x00, 0x02, 0x15}
	uuid := md[4:20]
	if bytes.Compare(discoveredPreamble, requiredPreamble) != 0 {
		fmt.Println(discoveredPreamble)
		fmt.Println(requiredPreamble)
		fmt.Println("preamble mismatch")
	}
	discoveredUUID := strings.ToUpper(hex.EncodeToString(uuid))
	// Do a lookup to determine the color of the found device.
	// If the UUID is not in the map, then the device must not
	// be a Tilt Hydrometer
	color, found := colorUUIDMap[discoveredUUID]
	if !found {
		// unknown UUID, so skip
		return
	}
	fmt.Printf("Found %s Tilt Hydrometer with address: %s\n", color, a.Addr().String())
}

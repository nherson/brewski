Brewski
===

Brewski is a command-line tool written in Go to monitor your homebrewing process. It will track various metrics regarding the status of your homebrew and allow for handling that data in various ways using configurable callbacks (logging, timeseries database, etc)

Motivation
---

I am a novice homebrewer. At the time of initial implementation of brewski, I only have 3 homebrews under my belt. However, I could already see that keeping track of simple variables like temperature and how it changes over time is nearly impossible without digital tools. A little sticker on the side of my brewing vessel with an analog temperature display wasn't cutting it, and I didn't like not knowing details about my beer's status throughout the day.

There are devices and services that allow you to track things like temperature, gravity, CO2 activity, etc but they tend to be proprietary and require your devices to be networked, and your data to be exported. Brewski aims to make adding device interfaces (development ongoing) and sensor reading handlers (also development ongoing) very simple. The goal of brewski is to make any brewing hardware, like sensors and probes, compatible with any computing hardware and/or software, like raspberry pis, timeseries databases, custom dashboaring, alerting pipelines, etc.

Architecture
---

Brewski is split between devices that can read data, and outputs that can act on that data. The two are linked using the `Sample` interface in the `measurement` package.  The `device/*` packages (seperated by category, but currently only `temperature` exists) implement the `device.Reader` interface to read data from the device and return a `Sample`.  On the other side, the `handlers` package has an interface called `Callback` which takes a `Sample` and does some arbitrary processing of the data within it.

The `Reader` implementations and `Callback` implementations are linked together with the `device.Poller` interface, which is a harness that glues together a `Reader` with a `Callback` to do some long-running, presumably periodic, processing of the device's data stream. This interface has a simple implementation in place called `Sensor` that just reads a `Sample` from the `Reader` at a specified interval and passes that `Sample` over to the registered `Callback` for handling.

The result is a library of devices and outputs that allow the user to link together any device to any data processing logic. Send temperature data to InfluxDB, send a text or e-mail when the gravity reading hits a target, etc. As long as the implementations exist, it should be easy to wire any device to any output.

Current Devices Supported
---
* DS18B20 (a very cheap temperature probe using onewire protocol, connected via sysfs)
* Tilt Hydrometer (all colors)

Device Support Wishlist
---
 * BrewNanny
 * Some digital pH meters

Output Methods Supported
---
* Logging (using zap)
* InfluxDB

Output Methods Wishlist
---
* Sending data to external process's `stdin` using defined text protocol/format (ala statsite stream commands)
* Fridge temperature controllers (future)

Developing
---
Checkout the project with `go get -u github.com/nherson/brewski`. This project uses Glide for dependency management. Use `glide install` to setup your vendor directory.

Additional Tools
---
There is a script to discover Tilt Hydrometers in `cmd/tilt-finder`.  Use `go run cmd/tilt-finder/main.go` or `go build ./cmd/tilt-finder && ./tilt-finder` to run. The script will listen for Bluetooth LE advertisements. When a Tilt is found, it will print the color discovered. Remember that Tilt Hydrometers standing straight up go into an idle mode and stop advertising data!

Todo
---
* Example setting up raspberry pi with brewski, influxdb, and grafana
* Configuration file support
* Rewrite executable into something more versatile and well thought out

Brewski
===

Brewski is a command-line tool written in Go to monitor your homebrewing process. It will track various metrics regarding the status of your homebrew and allow for handling that data in various ways using configurable callbacks (logging, timeseries database, etc)

Motivation
---

I am a novice homebrewer. At the time of initial implementation of brewski, I only have 3 homebrews under my belt. However, I could already see that keeping track of simple variables like temperature and how it changes over time is nearly impossible without digital tools. A little sticker on the side of my brewing vessel with an analog temperature display wasn't cutting it, and I didn't like not knowing details about my beer's status throughout the day.

There are devices and services that allow you to track things like temperature, gravity, CO2 activity, etc but they tend to be proprietary and require your devices to be networked, and your data to be exported. Brewski aims to make adding device interfaces (development ongoing) and sensor reading handlers (also development ongoing) very simple. The goal of brewski is to make any brewing hardware, like sensors and probes, compatible with any computing hardware and/or software, like raspberry pis, timeseries databases, custom dashboaring, alerting pipelines, etc.

Current Devices Supported
---
* Temperature
  * DS18B20 (a very cheap temperature probe using onewire protocol, connected via sysfs)

Device Support Wishlist
---
 * Tilt Hydrometer (I need to buy one)

Output Methods Supported
---
* Logging (using zap)
* InfluxDB

Output Methods Wishlist
---
* Sending data to external process's `stdin` using defined text protocol/format (ala statsite stream commands)

Developing
---
Checkout the project with `go get -u github.com/nherson/brewski`. This project uses Glide for dependency management. Use `glide install` to setup your vendor directory.

Todo
---
* Tests for existing codebase
* Example setting up raspberry pi with brewski, influxdb, and grafana
* Configuration file support

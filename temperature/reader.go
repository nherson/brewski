package temperature

// Reader is an interface that can read temperature data from
// an arbitrary source (currently DS1820B implemented via sysfs)
type Reader interface {
	ReadTemperature() (float32, float32, error)
}

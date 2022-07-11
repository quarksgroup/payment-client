package driver

// Driver identifies the payment platform driver
type Driver uint

// Supported drivers
const (
	DriverUnknown Driver = iota
	DriverFDI
	DriverAirtel
)

func (d Driver) String() string {
	switch d {
	case DriverFDI:
		return "fdi"
	case DriverAirtel:
		return "airtel"
	default:
		return "unknown"
	}
}

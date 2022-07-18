package driver

// Driver identifies the payment platform driver
type Driver uint

// Supported drivers
const (
	DriverUnknown Driver = iota
	DriverAirtel
)

func (d Driver) String() string {
	switch d {
	case DriverAirtel:
		return "airtel"
	default:
		return "unknown"
	}
}

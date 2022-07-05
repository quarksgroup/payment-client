package mtn

// Driver identifies the payment platform driver
type Driver uint

// Supported drivers
const (
	DriverUnknown Driver = iota
	DriverFDI
	//DriverMTN
)

func (d Driver) String() string {
	switch d {
	case DriverFDI:
		return "fdi"
	default:
		return "unknown"
	}
}

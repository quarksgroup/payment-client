package airtel

import (
	"context"
)

// Number represents information about a phone number details
type Number struct {
	Phone     string
	FirstName string
	LastName  string
	Status    bool
	HasPin    bool
}

// CheckNumber returns information about a phone number details
type CheckNumber interface {
	// Check takes a phone number and returns information about it and returns the error if it occurs
	Check(context.Context, string) (*Number, *Response, error)
}

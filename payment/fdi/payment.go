package fdi

import "context"

// Provider is the supported payment channel
type Provider string

// Payment represents as single transaction
type Payment struct {
	ID       string
	Amount   float64
	Wallet   string
	Provider Provider
}

// Status reports the status of a requested transaction
type Status struct {
	Ref   string
	GRef  string
	State string
}

// PaymentsService pulls amd pushes funds to and from an funds/payment provider.
type PaymentsService interface {
	// Pull initializes a pull of funds into our wallet
	Pull(context.Context, *Payment) (*Status, *Response, error)

	// Push initializes a push of funds to an external wallet transaction.
	Push(context.Context, *Payment) (*Status, *Response, error)
}

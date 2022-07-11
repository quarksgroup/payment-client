package airtel

import "context"

// Payment represents as single transaction
type Payment struct {
	ID     string
	Ref    string
	Amount int64
	Phone  string
}

// Status reports the status of a requested transaction
type Status struct {
	Ref          string `json:"ref"`
	Status       bool   `json:"status"`
	Message      string `json:"message"`
	ResponseCode string `json:"response_code"`
}

// PaymentService implements the PaymentService interface for push and pull requests
type PaymentService interface {
	// Push initializes a push of funds to an external wallet transaction.
	Push(context.Context, *Payment) (*Status, *Response, error)
	// Pull initializes a Pull of funds to an external wallet transaction.
	Pull(context.Context, *Payment) (*Status, *Response, error)
}

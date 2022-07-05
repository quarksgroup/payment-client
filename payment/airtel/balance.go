package airtel

import "context"

// Balance ...
type Balance struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
	Status   string `json:"status"`
}

// AccountService ...
type AccountService interface {
	// Balance returns balanceInfo about it.
	Balance(context.Context) (*Balance, *Response, error)
}

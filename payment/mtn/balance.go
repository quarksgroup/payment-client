package mtn

import "context"

// TODO: impelement balance service

// Account  ...
type Account struct {
	Currency         string  `json:"currency"`
	BalanceAvailable float64 `json:"balanceAvailable"`
	BalanceActual    float64 `json:"balanceActual"`
}

// Data  ...
type Data struct {
	Date     string    `json:"date"`
	Accounts []Account `json:"accounts"`
}

// Balance ...
type Balance struct {
	Status string `json:"status"`
	Data   Data   `json:"data"`
}

// BalanceService ...
type BalanceService interface {
	// Balance returns balanceInfo about it.
	Balance(context.Context) (*Balance, *Response, error)
}

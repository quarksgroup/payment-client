package payment

import "context"

// TODO: impelement balance service

// Account  ...
type Account struct {
	Currency         string `json:"currency"`
	BalanceAvailable uint64 `json:"balanceAvailable"`
	BalanceActual    uint64 `json:"balanceActual"`
}

// Data  ...
type Data struct {
	Date     string    `json:"date"`
	Accounts []Account `json:"accounts"`
}

// Balance ...
type Balances struct {
	Status string `json:"status"`
	Data   Data   `json:"data"`
}

// BalanceService ...
type BalanceService interface {
	// Balance returns balanceInfo about it.
	Balance(context.Context) (*Balances, *Response, error)
}

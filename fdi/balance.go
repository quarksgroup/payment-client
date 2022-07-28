package fdi

import (
	"context"

	"github.com/quarksgroup/payment-client/client"
)

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

func (c *Client) Balance(ctx context.Context) (*Balance, *client.Response, error) {

	endpoint := "balance/now"
	out := new(balanceResponse)
	res, err := c.do(ctx, "GET", endpoint, nil, out, true)
	return convertBalance(out), res, err

}

type balanceResponse struct {
	Status string `json:"status"`
	Data   struct {
		Date     string `json:"date"`
		Accounts []struct {
			Currency         string  `json:"currency"`
			BalanceAvailable float64 `json:"balanceAvailable"`
			BalanceActual    float64 `json:"balanceActual"`
		} `json:"accounts"`
	} `json:"data"`
}

func convertBalance(res *balanceResponse) *Balance {

	data := &Data{
		Date:     res.Data.Date,
		Accounts: make([]Account, 0),
	}

	for _, item := range res.Data.Accounts {
		acc := &Account{
			Currency:         item.Currency,
			BalanceAvailable: item.BalanceAvailable,
			BalanceActual:    item.BalanceActual,
		}
		data.Accounts = append(data.Accounts, *acc)
	}

	return &Balance{
		Status: res.Status,
		Data:   *data,
	}
}

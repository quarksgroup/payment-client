package client

import (
	"context"

	"github.com/quarksgroup/payment-client/fdi"
)

func (c *Client) Balance(ctx context.Context) (*fdi.Balance, *fdi.Response, error) {
	endpoint := "balance/now"
	out := new(balanceResponse)
	res, err := c.do(ctx, "GET", endpoint, nil, out)
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

func convertBalance(res *balanceResponse) *fdi.Balance {

	data := &fdi.Data{
		Date:     res.Data.Date,
		Accounts: make([]fdi.Account, 0),
	}

	for _, item := range res.Data.Accounts {
		acc := &fdi.Account{
			Currency:         item.Currency,
			BalanceAvailable: item.BalanceAvailable,
			BalanceActual:    item.BalanceActual,
		}
		data.Accounts = append(data.Accounts, *acc)
	}

	return &fdi.Balance{
		Status: res.Status,
		Data:   *data,
	}
}

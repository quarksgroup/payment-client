package fdi

import (
	"context"

	"github.com/quarksgroup/payment-client/fdi"
)

type BalanceService struct {
	client *wrapper
}

func (s *BalanceService) Balance(ctx context.Context) (*fdi.Balance, *fdi.Response, error) {
	endpoint := "balance/now"
	out := new(balanceResponse)
	res, err := s.client.do(ctx, "GET", endpoint, nil, out)
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
		account := &fdi.Account{
			Currency:         item.Currency,
			BalanceAvailable: item.BalanceAvailable,
			BalanceActual:    item.BalanceActual,
		}
		data.Accounts = append(data.Accounts, *account)
	}

	return &fdi.Balance{
		Status: res.Status,
		Data:   *data,
	}
}

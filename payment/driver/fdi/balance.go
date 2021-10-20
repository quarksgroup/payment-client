package fdi

import (
	"context"

	"github.com/quarksgroup/payment-client/payment"
)

type BalanceService struct {
	client *wrapper
}

func (s *BalanceService) Balance(ctx context.Context) (*payment.Balance, *payment.Response, error) {

	endpoint := "/balance/now"
	out := new(balanceResponse)
	res, err := s.client.do(ctx, "GET", endpoint, nil, out)
	return convertBalance(out), res, err

}

type balanceResponse struct {
	Status string `json:"status"`
	Data   struct {
		Date     string `json:"date"`
		Accounts []struct {
			Currency         string `json:"currency"`
			BalanceAvailable uint64 `json:"balanceAvailable"`
			BalanceActual    uint64 `json:"balanceActual"`
		} `json:"accounts"`
	} `json:"data"`
}

func convertBalance(res *balanceResponse) *payment.Balance {

	data := &payment.Data{
		Date:     res.Data.Date,
		Accounts: make([]payment.Account, 0),
	}

	for _, item := range res.Data.Accounts {
		account := &payment.Account{
			Currency:         item.Currency,
			BalanceAvailable: item.BalanceAvailable,
			BalanceActual:    item.BalanceActual,
		}
		data.Accounts = append(data.Accounts, *account)
	}

	return &payment.Balance{
		Status: res.Status,
		Data:   *data,
	}
}

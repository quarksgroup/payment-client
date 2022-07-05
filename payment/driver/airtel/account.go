package airtel

import (
	"context"

	"github.com/quarksgroup/payment-client/payment/airtel"
)

type accountService struct {
	client *wrapper
}

func (a *accountService) Balance(ctx context.Context) (*airtel.Balance, *airtel.Response, error) {
	endpoint := "standard/v1/users/balance"
	out := new(balanceResponse)
	res, err := a.client.do(ctx, "GET", endpoint, nil, out, nil)
	return convertBalance(out), res, err

}

type balanceResponse struct {
	Data struct {
		Balance       string `json:"balance"`
		Currency      string `json:"currency"`
		AccountStatus string `json:"account_status"`
	} `json:"data"`
	Status struct {
		Code        string `json:"code"`
		Success     bool   `json:"success"`
		Result_code string `json:"result_code"`
		Message     string `json:"message"`
	} `json:"status"`
}

func convertBalance(res *balanceResponse) *airtel.Balance {
	return &airtel.Balance{
		Amount:   res.Data.Balance,
		Currency: res.Data.Currency,
		Status:   res.Data.AccountStatus,
	}
}

//account.go contain the implementation of airtel protal account information
package client

import (
	"context"

	"github.com/quarksgroup/payment-client/airtel"
)

//Balance responsible for returning the airtel account balance
func (c *Client) Balance(ctx context.Context) (*airtel.Balance, *airtel.Response, error) {

	if err := c.renewToken(ctx); err != nil {
		return nil, nil, err
	}

	endpoint := "standard/v1/users/balance"
	out := new(balanceResponse)
	res, err := c.do(ctx, "GET", endpoint, nil, out, nil)
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

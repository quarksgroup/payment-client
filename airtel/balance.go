package airtel

import (
	"context"
	"net/http"
	"strconv"

	"github.com/quarksgroup/payment-client/client"
)

// Balance ...
type Balance struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
	Status   string `json:"status"`
}

//Balance responsible for returning the airtel account balance
func (c *Client) Balance(ctx context.Context) (*Balance, *client.Response, error) {
	endpoint := "standard/v1/users/balance"

	headers := http.Header{
		"X-country":  []string{c.Country},
		"X-Currency": []string{c.Currency},
	}

	out := new(balanceResponse)

	res, err := c.do(ctx, "GET", endpoint, nil, out, headers, true)

	// this is supposed to be taken care of from inside client.do
	if !out.Status.Success {

		code, _ := strconv.ParseInt(out.Status.Code, 10, 64)

		return nil, nil, &Error{Code: int(code), Message: out.ErrMsg}
	}

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
	ErrMsg string `json:"status_message"`
}

func convertBalance(res *balanceResponse) *Balance {
	return &Balance{
		Amount:   res.Data.Balance,
		Currency: res.Data.Currency,
		Status:   res.Data.AccountStatus,
	}
}

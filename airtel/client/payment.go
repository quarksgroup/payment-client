// payment.go this responsible for the implementation of refund and credit from your airtel portal or to your portal
// from or to specific phone account and vise-versa
package client

import (
	"context"
	"net/http"

	"github.com/quarksgroup/payment-client/airtel"
)

//Push this is responsible for the implementation of cash collection to your airtel portal from phone wallet
func (c *Client) Push(ctx context.Context, req *airtel.Payment) (*airtel.Status, *airtel.Response, error) {

	if err := c.renewToken(ctx); err != nil {
		return nil, nil, err
	}

	endpoint := "merchant/v1/payments/"

	sub := &subscriber{
		Country:  c.Client.Country,
		Currency: c.Client.Currency,
		Msisdn:   req.Phone,
	}

	tx := &tx{
		Id:     req.ID,
		Amount: req.Amount,
	}

	in := &pushRequest{
		Reference:   req.Ref,
		Transaction: tx,
		Subscriber:  sub,
	}

	header := http.Header{
		"X-country":  []string{c.Client.Country},
		"X-Currency": []string{c.Client.Currency},
	}

	out := new(pushResponse)

	res, err := c.do(ctx, "POST", endpoint, in, out, header)

	if err != nil {
		return nil, nil, err
	}
	// if !out.Status.Success {
	// 	return nil, nil, &airtel.Error{Code: http.StatusBadRequest, Message: out.Status.Message}
	// }

	return convertPush(out), res, err
}

// Pull this is responsible for cash distribution or disbursements from your phone wallet to your airtel portal
func (c *Client) Pull(ctx context.Context, req *airtel.Payment) (*airtel.Status, *airtel.Response, error) {

	if err := c.renewToken(ctx); err != nil {
		return nil, nil, err
	}

	endpoint := "standard/v1/disbursements/"

	tx := &tx{
		Id:     req.ID,
		Amount: req.Amount,
	}
	in := &pullRequest{
		Reference:   req.Ref,
		Pin:         c.Client.EncryptedPin,
		Transaction: tx,
	}
	in.Payee.Msisdn = req.Phone

	header := http.Header{
		"X-country":  []string{c.Client.Country},
		"X-Currency": []string{c.Client.Currency},
	}

	out := new(pullResponse)

	res, err := c.do(ctx, "POST", endpoint, in, out, header)

	if err != nil {
		return nil, nil, err
	}

	return convertPull(out), res, err
}

//subscriber...
type subscriber struct {
	Country  string `json:"country"`
	Currency string `json:"currency"`
	Msisdn   string `json:"msisdn"`
}

//tx..
type tx struct {
	Id     string `json:"id"`
	Amount int64  `json:"amount"`
}

//pushRequest represent push collection request payload
type pushRequest struct {
	Reference   string      `json:"reference"`
	Subscriber  *subscriber `json:"subscriber"`
	Transaction *tx         `json:"transaction"`
}

//Status...
type status struct {
	ResponseCode string `json:"response_code"`
	Code         string `json:"code"`
	Success      bool   `json:"success"`
	Result_code  string `json:"result_code"`
	Message      string `json:"message"`
}

//pushResponse represents push collection response body
type pushResponse struct {
	Data struct {
		Transaction struct {
			Id     string `json:"id"`
			Status string `json:"status"`
		} `json:"transaction"`
	} `json:"data"`
	Status *status `json:"status"`
}

func convertPush(res *pushResponse) *airtel.Status {
	return &airtel.Status{
		Ref:          res.Data.Transaction.Id,
		Status:       res.Status.Success,
		ResponseCode: res.Status.Result_code,
		Message:      res.Status.Message,
	}
}

//pullRequest represent push collection request payload
type pullRequest struct {
	Reference string `json:"reference"`
	Pin       string `json:"pin"`
	Payee     struct {
		Msisdn string `json:"msisdn"`
	} `json:"payee"`
	Transaction *tx `json:"transaction"`
}

//pullResponse represents push collection response body
type pullResponse struct {
	Data struct {
		Transaction struct {
			Id     string `json:"id"`
			Status string `json:"status"`
		} `json:"transaction"`
	} `json:"data"`
	Status *status `json:"status"`
}

func convertPull(res *pullResponse) *airtel.Status {
	return &airtel.Status{
		Ref:          res.Data.Transaction.Id,
		Status:       res.Status.Success,
		ResponseCode: res.Status.ResponseCode,
		Message:      res.Status.Message,
	}
}

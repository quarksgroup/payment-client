package airtel

import (
	"context"
	"net/http"

	"github.com/quarksgroup/payment-client/client"
)

// Payment represents as single transaction
type Payment struct {
	ID     string
	Ref    string
	Amount int64
	Phone  string
}

// Status reports the status of a requested transaction
type Status struct {
	Ref          string `json:"ref"`
	Status       bool   `json:"status"`
	Message      string `json:"message"`
	ResponseCode string `json:"response_code"`
}

//Pull this is responsible for the implementation of cash collection to your airtel portal from phone wallet
func (c *Client) Pull(ctx context.Context, req *Payment) (*Status, *client.Response, error) {
	endpoint := "merchant/v1/payments/"

	sub := &subscriber{
		Country:  c.Country,
		Currency: c.Currency,
		Msisdn:   req.Phone,
	}

	tx := &tx{
		Id:     req.ID,
		Amount: req.Amount,
	}

	in := &pullRequest{
		Reference:   req.Ref,
		Transaction: tx,
		Subscriber:  sub,
	}

	header := http.Header{
		"X-country":  []string{c.Country},
		"X-Currency": []string{c.Currency},
	}

	out := new(pullResponse)

	res, err := c.do(ctx, "POST", endpoint, in, out, header, true)

	if err != nil {
		return nil, nil, err
	}

	if !out.Status.Success {
		return nil, nil, &Error{Code: http.StatusBadRequest, Message: out.Status.Message}
	}

	return convertPull(out), res, err
}

// Push this is responsible for cash distribution or disbursements from your phone wallet to your airtel portal
func (c *Client) Push(ctx context.Context, req *Payment) (*Status, *client.Response, error) {
	endpoint := "standard/v1/disbursements/"

	tx := &tx{
		Id:     req.ID,
		Amount: req.Amount,
	}
	in := &pushRequest{
		Reference:   req.Ref,
		Pin:         c.EncryptedPin,
		Transaction: tx,
	}
	in.Payee.Msisdn = req.Phone

	header := http.Header{
		"X-country":  []string{c.Country},
		"X-Currency": []string{c.Currency},
	}

	out := new(pushResponse)

	res, err := c.do(ctx, "POST", endpoint, in, out, header, true)

	if err != nil {
		return nil, nil, err
	}

	if !out.Status.Success {
		return nil, nil, &Error{Code: http.StatusBadRequest, Message: out.Status.Message}
	}

	return convertPush(out), res, err
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

//pullRequest represent push collection request payload
type pullRequest struct {
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

func convertPull(res *pullResponse) *Status {
	return &Status{
		Ref:          res.Data.Transaction.Id,
		Status:       res.Status.Success,
		ResponseCode: res.Status.ResponseCode,
		Message:      res.Status.Message,
	}
}

//pushRequest represent push collection request payload
type pushRequest struct {
	Reference string `json:"reference"`
	Pin       string `json:"pin"`
	Payee     struct {
		Msisdn string `json:"msisdn"`
	} `json:"payee"`
	Transaction *tx `json:"transaction"`
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

func convertPush(res *pushResponse) *Status {
	return &Status{
		Ref:          res.Data.Transaction.Id,
		Status:       res.Status.Success,
		ResponseCode: res.Status.ResponseCode,
		Message:      res.Status.Message,
	}
}

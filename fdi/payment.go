package fdi

import (
	"context"

	"github.com/quarksgroup/payment-client/client"
)

// Provider is the supported payment channel
type Provider string

// Payment represents as single transaction
type Payment struct {
	ID       string
	Amount   float64
	Wallet   string
	Provider Provider
}

// Status reports the status of a requested transaction
type Status struct {
	Ref   string
	GRef  string
	State string
}

// FDI Supported providers
const (
	MTN    Provider = "momo-mtn-rw"
	Airtel Provider = "momo-airtel-rw"
)

func (c *Client) Pull(ctx context.Context, py *Payment) (*Status, *client.Response, error) {

	endpoint := "momo/pull"
	in := &paymentRequest{
		Ref:      py.ID,
		MSISDN:   py.Wallet,
		Amount:   py.Amount,
		Channel:  string(py.Provider),
		Callback: c.ReportURL.String(),
	}
	out := new(paymentResponse)
	res, err := c.do(ctx, "POST", endpoint, in, out)
	return convertResponse(out), res, err
}

func (c *Client) Push(ctx context.Context, py *Payment) (*Status, *client.Response, error) {

	endpoint := "momo/push"
	in := &paymentRequest{
		Ref:      py.ID,
		MSISDN:   py.Wallet,
		Amount:   py.Amount,
		Channel:  string(py.Provider),
		Callback: c.ReportURL.String(),
	}
	out := new(paymentResponse)
	res, err := c.do(ctx, "POST", endpoint, in, out)
	return convertResponse(out), res, err
}

type paymentRequest struct {
	Ref      string  `json:"trxRef"`
	Channel  string  `json:"channelId"`
	Account  string  `json:"accountId"`
	MSISDN   string  `json:"msisdn"`
	Amount   float64 `json:"amount"`
	Callback string  `json:"callback"`
}

type paymentResponse struct {
	Status string `json:"status"`
	Data   data   `json:"data"`
}

// Data ...
type data struct {
	Ref     string `json:"trxRef,omitempty"`
	Token   string `json:"token,omitempty"`
	Gateway string `json:"gwRef,omitempty"`
	State   string `json:"state,omitempty"`
}

func convertResponse(res *paymentResponse) *Status {
	return &Status{
		Ref:   res.Data.Ref,
		GRef:  res.Data.Gateway,
		State: res.Data.State,
	}
}

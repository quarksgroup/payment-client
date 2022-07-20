package client

import (
	"context"

	"github.com/quarksgroup/payment-client/fdi"
)

// FDI Supported providers
const (
	MTN    fdi.Provider = "momo-mtn-rw"
	Airtel fdi.Provider = "momo-airtel-rw"
)

func (c *Client) Pull(ctx context.Context, py *fdi.Payment) (*fdi.Status, *fdi.Response, error) {
	endpoint := "momo/pull"
	in := &paymentRequest{
		Ref:      py.ID,
		MSISDN:   py.Wallet,
		Amount:   py.Amount,
		Channel:  string(py.Provider),
		Callback: c.Client.ReportURL.String(),
	}
	out := new(paymentResponse)
	res, err := c.do(ctx, "POST", endpoint, in, out)
	return convertResponse(out), res, err
}

func (c *Client) Push(ctx context.Context, py *fdi.Payment) (*fdi.Status, *fdi.Response, error) {
	endpoint := "momo/push"
	in := &paymentRequest{
		Ref:      py.ID,
		MSISDN:   py.Wallet,
		Amount:   py.Amount,
		Channel:  string(py.Provider),
		Callback: c.Client.ReportURL.String(),
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

func convertResponse(res *paymentResponse) *fdi.Status {
	return &fdi.Status{
		Ref:   res.Data.Ref,
		GRef:  res.Data.Gateway,
		State: res.Data.State,
	}
}

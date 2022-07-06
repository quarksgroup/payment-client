package fdi

import (
	"context"

	"github.com/iradukunda1/payment-staging/payment/mtn"
)

// FDI Supported providers
const (
	MTN    mtn.Provider = "momo-mtn-rw"
	Airtel mtn.Provider = "momo-airtel-rw"
)

type paymentsService struct {
	client *wrapper
}

func (s *paymentsService) Pull(ctx context.Context, py *mtn.Payment) (*mtn.Status, *mtn.Response, error) {
	endpoint := "momo/pull"
	in := &paymentRequest{
		Ref:      py.ID,
		MSISDN:   py.Wallet,
		Amount:   py.Amount,
		Channel:  string(py.Provider),
		Callback: s.client.ReportURL.String(),
	}
	out := new(paymentResponse)
	res, err := s.client.do(ctx, "POST", endpoint, in, out)
	return convertResponse(out), res, err
}

func (s *paymentsService) Push(ctx context.Context, py *mtn.Payment) (*mtn.Status, *mtn.Response, error) {
	endpoint := "momo/push"
	in := &paymentRequest{
		Ref:      py.ID,
		MSISDN:   py.Wallet,
		Amount:   py.Amount,
		Channel:  string(py.Provider),
		Callback: s.client.ReportURL.String(),
	}
	out := new(paymentResponse)
	res, err := s.client.do(ctx, "POST", endpoint, in, out)
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

func convertResponse(res *paymentResponse) *mtn.Status {
	return &mtn.Status{
		Ref:   res.Data.Ref,
		GRef:  res.Data.Gateway,
		State: res.Data.State,
	}
}

var _ (mtn.PaymentsService) = (*paymentsService)(nil)

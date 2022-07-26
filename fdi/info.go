package fdi

import (
	"context"
	"fmt"
	"time"

	"github.com/quarksgroup/payment-client/client"
)

// Info represents information about a transaction
type Info struct {
	ID        string
	Amount    float64
	Cost      float64
	Status    string
	Type      string
	CreatedAt time.Time
}

func (c *Client) TransactionInfo(ctx context.Context, ref string) (*Info, *client.Response, error) {

	endpoint := fmt.Sprintf("momo/trx/%s/info", ref)
	out := new(infoResponse)
	res, err := c.do(ctx, "GET", endpoint, nil, out)
	return convertInfo(out), res, err
}

type infoResponse struct {
	Status string `json:"status"`
	Data   struct {
		ID         string    `json:"id"`
		Ref        string    `json:"trxRef"`
		Type       string    `json:"trxType"`
		ChannelID  string    `json:"channelId"`
		ChannelRef string    `json:"channelRef"`
		MSISDN     string    `json:"msisdn"`
		Amount     float64   `json:"amount"`
		Fees       float64   `json:"trxFees"`
		Currency   string    `json:"currency"`
		TrxStatus  string    `json:"trxStatus"`
		CreatedAt  time.Time `json:"created_at"`
		Callback   string    `json:"callback"`
	} `json:"data"`
}

func convertInfo(info *infoResponse) *Info {
	return &Info{
		ID:        info.Data.Ref,
		Amount:    info.Data.Amount,
		Cost:      info.Data.Fees,
		Status:    info.Data.TrxStatus,
		Type:      info.Data.Type,
		CreatedAt: info.Data.CreatedAt,
	}
}

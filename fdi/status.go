package fdi

import (
	"context"

	"github.com/quarksgroup/payment-client/client"
)

type APIStatus struct {
	State string
}

func (c *Client) Status(ctx context.Context) (*APIStatus, *client.Response, error) {
	endpoint := "status"
	out := new(statusResponse)
	res, err := c.do(ctx, "GET", endpoint, nil, out, true)
	return convertStatus(out), res, err
}

type statusResponse struct {
	Status string `json:"status"`
	Data   struct {
		State string `json:"state"`
	} `json:"data"`
}

func convertStatus(res *statusResponse) *APIStatus {
	return &APIStatus{
		State: res.Data.State,
	}
}

// Package fdi implements the payment.Client for the fdi(https://fdipaymentsapi.docs.apiary.io/)
package fdi

import (
	"bytes"
	"context"
	"encoding/json"
	"net/url"
	"strings"

	"github.com/quarksgroup/payment-client/payment"
)

// New creates a new payment.Client instance backed by the payment.DriverFDI
func New(uri, callback string) (*payment.Client, error) {
	base, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	if !strings.HasSuffix(base.Path, "/") {
		base.Path = base.Path + "/"
	}
	report, err := url.Parse(callback)
	if err != nil {
		return nil, err
	}

	client := &wrapper{new(payment.Client)}
	client.BaseURL = base
	client.ReportURL = report

	client.Driver = payment.DriverFDI

	// initialize services
	client.Payments = &paymentsService{client}
	client.Info = &infoService{client}
	client.Auth = &authService{client}

	return client.Client, nil
}

type wrapper struct {
	*payment.Client
}

// NewDefault returns a new FDI API client using the`
// default "https://payments-api.fdibiz.com/v2" address.
func NewDefault(callback string) *payment.Client {
	client, _ := New("https://payments-api.fdibiz.com/v2", callback)
	return client
}

// do wraps the Client.Do function by creating the Request and
// unmarshalling the response.
func (c *wrapper) do(ctx context.Context, method, path string, in, out interface{}) (*payment.Response, error) {
	req := &payment.Request{
		Method: method,
		Path:   path,
	}

	// if we are posting or putting data, we need to
	// write it to the body of the request.
	if in != nil {
		buf := new(bytes.Buffer)
		_ = json.NewEncoder(buf).Encode(in)
		req.Header = map[string][]string{
			"Content-Type": {"application/json"},
		}
		req.Body = buf
	}

	// execute the http request
	res, err := c.Client.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// if an error is encountered, unmarshal and return the
	// error response.
	if res.Status > 299 {
		err := new(Err)
		_ = json.NewDecoder(res.Body).Decode(err)
		return res, &payment.Error{Code: res.Status, Message: err.Data.Message}
	}

	if out == nil {
		return res, nil
	}

	// if a json response is expected, parse and return
	// the json response.
	return res, json.NewDecoder(res.Body).Decode(out)
}

// Error represents a Github error.
type Err struct {
	Status string `json:"status"`
	Data   struct {
		Message string `json:"message"`
	}
}

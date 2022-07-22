// Package airtel implements the payment.Client for the airtel(https://developers.airtel.africa/documentation)
package airtel

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/iradukunda1/payment-staging/payment/airtel"
)

const (
	baseUrl  = "https://openapi.airtel.africa"
	currency = "RWF"
	country  = "RW"
)

// New creates a new payment.Client instance backed by the payment.DriverAirtel
func New(uri, pin, currency, country string) (*airtel.Client, error) {
	base, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	if !strings.HasSuffix(base.Path, "/") {
		base.Path = base.Path + "/"
	}

	client := &wrapper{new(airtel.Client)}
	client.BaseURL = base
	client.Country = country
	client.Currency = currency
	client.EncryptedPin = pin
	client.Auth = &authService{client}
	client.Driver = airtel.DriverAirtel
	client.Account = &accountService{client}
	client.CheckNumber = &checkNumberService{client}
	client.Collections = &collectionService{client}
	client.Disbursement = &disbursementService{client}

	// initialize services

	return client.Client, nil
}

type wrapper struct {
	*airtel.Client
}

// NewDefault returns a new AIRTEL API client using the`
// default "https://openapi.airtel.africa" address, country RW(Rwanda) and RWF(Rwandan franc).
func NewDefault(pin string) *airtel.Client {
	client, _ := New(baseUrl, pin, currency, country)
	return client
}

// do wraps the Client.Do function by creating the Request and
// unmarshalling the response.
func (c *wrapper) do(ctx context.Context, method, path string, in, out interface{}, headers http.Header) (*airtel.Response, error) {
	req := &airtel.Request{
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

	for k, v := range headers {
		req.Header[k] = v
	}

	// execute the http request
	res, err := c.Client.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// if an error is encountered, unmarshal and return the
	// error response.
	if res.Status > 299 && res.Status < 499 {
		err := new(Err)
		_ = json.NewDecoder(res.Body).Decode(err)
		return res, &airtel.Error{Code: res.Status, Message: err.ErrorDescprition}
	}

	if res.Status > 499 {
		return res, &airtel.Error{Code: res.Status, Message: "Something went wrong"}
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
	Error            string `json:"error"`
	ErrorDescprition string `json:"error_description"`
}

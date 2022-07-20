// Package airtel implements the payment.Client for the airtel(https://developers.airtel.africa/documentation)
package airtel

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/quarksgroup/payment-client/airtel"
	"github.com/quarksgroup/payment-client/airtel/driver"
)

const (
	baseUrl   = "https://openapi.airtel.africa"
	currency  = "RWF"
	country   = "RW"
	userAgent = "paypack"
	retry     = 3
)

//This wrapper all client implentation of airtel
type wrapper struct {
	*airtel.Client
}

// New creates a new payment.Client instance backed by the payment.DriverAirtel
func New(uri, pin, id, sceret, grant, currency, country string, retry int) (*airtel.Client, error) {
	base, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	if !strings.HasSuffix(base.Path, "/") {
		base.Path = base.Path + "/"
	}

	transport := &airtel.RetryTransport{
		Next:       http.DefaultTransport,
		MaxRetries: retry,
		Logger:     os.Stdout,
		Delay:      time.Duration(1 * time.Second),
		Source:     ContextTokenSource(),
		ClientId:   id,
		Sceret:     sceret,
		Grant:      grant,
	}

	httpClient := &http.Client{
		Transport: transport,
	}

	client := &wrapper{new(airtel.Client)}
	client.Client.Client = httpClient
	client.BaseURL = base
	client.Country = country
	client.UserAgent = userAgent
	client.Currency = currency
	client.EncryptedPin = pin
	client.Auth = &authService{client}
	transport.Auth = &authService{client}
	client.Driver = driver.DriverAirtel
	client.Account = &accountService{client}
	client.CheckNumber = &checkNumberService{client}
	client.Payments = &paymentsService{client}

	return client.Client, nil
}

// NewDefault returns a new AIRTEL API client using the
//But it take payment credential parameter
// default "https://openapi.airtel.africa" address, country RW(Rwanda) and RWF(Rwandan franc).
func NewDefault(pin, clientId, secret, grant string) *airtel.Client {
	client, _ := New(baseUrl, pin, clientId, secret, grant, currency, country, retry)
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

	// set the request headers
	if headers != nil {
		for k, v := range headers {
			req.Header[k] = v
		}
	}

	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
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

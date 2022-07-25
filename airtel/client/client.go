// package client implements the payment.Client for the airtel(https://developers.airtel.africa/documentation)
package client

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
)

//This is const value for default configuration of airtel client
const (
	baseUrl   = "https://openapi.airtel.africa"
	currency  = "RWF"
	country   = "RW"
	userAgent = "paypack / 2"
	retry     = 3
)

//This Client all client implentation of airtel.CLient
type Client struct {
	*airtel.Client
}

// New creates a new airtel.Client instance backed by the http.Client
func New(uri, pin, clientId, clientSceret, grant, currency, country string, retry int) (*Client, error) {
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
	}

	httpClient := &http.Client{
		Transport: transport,
	}

	client := &Client{new(airtel.Client)}
	client.Client.Client = httpClient
	client.BaseURL = base
	client.Country = country
	client.UserAgent = userAgent
	client.Currency = currency
	client.EncryptedPin = pin
	client.ClientId = &clientId
	client.ClientSceret = &clientSceret
	client.GrantType = &grant

	token, _, err := client.login(context.Background(), clientId, clientSceret, grant)

	if err != nil {
		return nil, err
	}
	client.Client.Token = token

	return client, nil
}

// NewDefault returns a new AIRTEL API client using the http.Client
// But it take payment credential parameter
// default "https://openapi.airtel.africa" address, country RW(Rwanda) and RWF(Rwandan franc).
func NewDefault(pin, clientId, secret, grant string) (*Client, error) {
	return New(baseUrl, pin, clientId, secret, grant, currency, country, retry)
}

// do wraps the Client.Do function by creating the Request and
// unmarshalling the response according to user expected output.
func (c *Client) do(ctx context.Context, method, path string, in, out interface{}, headers http.Header) (*airtel.Response, error) {
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
		req.Header = headers
	}

	if c.UserAgent != "" {
		req.Header.Add("User-Agent", c.UserAgent)
	}

	// execute the http request using airtel.Client.Do()
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

// Error represents airtel error.
type Err struct {
	Error            string `json:"error"`
	ErrorDescprition string `json:"error_description"`
}

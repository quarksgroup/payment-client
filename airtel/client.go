package airtel

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/rehttp"
	"github.com/quarksgroup/payment-client/client"
	"github.com/quarksgroup/payment-client/token"
)

//This is const value for default configuration of airtel client
const (
	baseUrl        = "https://openapi.airtel.africa"
	currency       = "RWF"
	country        = "RW"
	userAgent      = "paypack"
	defaultRetries = 3
)

type Config struct {
	ClientId string
	Secret   string
	Grant    string
	Pin      string
	Currency string
	Country  string
}

//This Client all client implentation of airtel.CLient
type Client struct {
	inner        *client.Client
	TokenSource  token.TokenSource
	BaseURL      *url.URL
	Country      string
	Currency     string
	EncryptedPin string
	ClientId     string
	ClientSceret string
	GrantType    string
}

// New creates a new airtel.Client instance backed by the http.Client
func New(cfg *Config, source token.TokenSource, uri string, debug bool, retries int) (*Client, error) {
	base, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	if !strings.HasSuffix(base.Path, "/") {
		base.Path = base.Path + "/"
	}

	retryTransport := rehttp.NewTransport(
		http.DefaultTransport,
		rehttp.RetryAll(
			rehttp.RetryMaxRetries(retries),
			rehttp.RetryAny(
				rehttp.RetryTemporaryErr(),
				rehttp.RetryStatuses(502, 503),
			),
		),
		rehttp.ExpJitterDelay(100*time.Millisecond, 1*time.Second),
	)

	httpClient := &http.Client{
		Transport: retryTransport,
	}

	inner := &client.Client{
		Client:    httpClient,
		BaseURL:   base,
		UserAgent: userAgent,
	}

	client := new(Client)
	client.inner = inner

	client.BaseURL = base
	client.Country = cfg.Country
	client.Currency = cfg.Currency
	client.EncryptedPin = cfg.Pin
	client.ClientId = cfg.ClientId
	client.ClientSceret = cfg.Secret
	client.GrantType = cfg.Grant
	client.TokenSource = source

	if client.TokenSource == nil {
		client.TokenSource, err = newTokenSource(client, cfg)
		if err != nil {
			return nil, err
		}

	}

	return client, nil
}

// NewDefault returns a new AIRTEL API client using the http.Client
// But it take payment credential parameter
// default "https://openapi.airtel.africa" address, country RW(Rwanda) and RWF(Rwandan franc).
func NewDefault(pin, clientId, secret, grant string) (*Client, error) {
	config := &Config{
		Pin:      pin,
		ClientId: clientId,
		Secret:   secret,
		Grant:    grant,
		Currency: currency,
		Country:  country,
	}
	return New(config, nil, baseUrl, false, defaultRetries)
}

// do wraps the Client.Do function by creating the Request and
// unmarshalling the response according to user expected output.
func (c *Client) do(ctx context.Context, method, path string, in, out interface{}, headers http.Header, addToken bool) (*client.Response, error) {
	req := &client.Request{
		Method: method,
		Path:   path,
		Header: make(http.Header),
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
			req.Header.Set(k, v[0])
		}
	}

	if c.inner.UserAgent != "" {
		req.Header.Set("User-Agent", c.inner.UserAgent)
	}

	// set auth token from TokenSource
	if c.TokenSource != nil && addToken {
		token, err := c.TokenSource.Token(ctx)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+token.Token)
	}

	// execute the http request using airtel.Client.Do()
	res, err := c.inner.Do(ctx, req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	// if an error is encountered, unmarshal and return the
	// error response.
	switch res.Status {
	case http.StatusUnauthorized:
		_, err = c.TokenSource.Token(ctx)
		if err != nil {
			return nil, err
		}
		return c.do(ctx, method, path, in, out, headers, addToken)
	default:
		if res.Status > 299 && res.Status < 499 {
			err := new(Err)
			_ = json.NewDecoder(res.Body).Decode(err)
			return res, &Error{Code: res.Status, Message: err.ErrorDescprition}
		}

		if res.Status > 499 {
			return res, &Error{Code: res.Status, Message: "Something went wrong"}
		}
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

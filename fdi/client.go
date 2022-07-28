package fdi

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

const (
	baseUrl   = "https://payments-api.fdibiz.com/v2"
	retry     = 3 // this is the defualt retry for fdi.Transport MaxRetries of RoundTripp
	userAgent = "paypack"
)

type Config struct {
	ClientId string
	Secret   string
	CallBack string
}

// Client manages communication with a payment gateways API.
type Client struct {
	inner *client.Client

	// Base URL for API requests.
	BaseURL *url.URL

	// ReportURL is the url to callback for payment reports
	ReportURL *url.URL

	// Account crendential
	ClientId, ClientSceret string

	TokenSource token.TokenSource
}

// New creates a new fdi.Client instance backed by the  http.Client instance
func New(uri string, cfg *Config, source token.TokenSource, retry int) (*Client, error) {
	base, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	if !strings.HasSuffix(base.Path, "/") {
		base.Path = base.Path + "/"
	}
	report, err := url.Parse(cfg.CallBack)
	if err != nil {
		return nil, err
	}
	retryTransport := rehttp.NewTransport(
		http.DefaultTransport,
		rehttp.RetryAll(
			rehttp.RetryMaxRetries(retry),
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
	client.ReportURL = report
	client.ClientId = cfg.ClientId
	client.ClientSceret = cfg.Secret
	client.TokenSource = source

	if client.TokenSource == nil {
		client.TokenSource, err = newTokenSource(client, cfg)
		if err != nil {
			return nil, err
		}

	}

	return client, nil
}

// NewDefault returns a new FDI API client using the`
// default "https://payments-api.fdibiz.com/v2" address.
func NewDefault(callback, client_id, sceret string) (*Client, error) {
	config := &Config{
		ClientId: client_id,
		Secret:   sceret,
		CallBack: callback,
	}
	return New(baseUrl, config, nil, retry)
}

// do wraps the Client.Do function by creating the Request and
// unmarshalling the response according to user wish expected output.
func (c *Client) do(ctx context.Context, method, path string, in, out interface{}, addToken bool) (*client.Response, error) {
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
	// execute the http request
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
		return c.do(ctx, method, path, in, out, addToken)
	default:
		if res.Status > 299 && res.Status < 499 {
			err := new(Err)
			_ = json.NewDecoder(res.Body).Decode(err)
			return res, &Error{Code: res.Status, Message: err.Data.Message}
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

// Error represents a fdi error.
type Err struct {
	Status string `json:"status"`
	Data   struct {
		Message string `json:"message"`
	}
}

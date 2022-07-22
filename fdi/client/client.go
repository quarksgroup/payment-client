// package client implements the fdi.Client for the fdi(https://fdipaymentsapi.docs.apiary.io/)
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

	"github.com/quarksgroup/payment-client/fdi"
)

const (
	baseUrl = "https://payments-api.fdibiz.com/v2"
	retry   = 3 // this is the defualt retry for fdi.Transport MaxRetries of RoundTripp
)

type Client struct {
	*fdi.Client
}

// New creates a new fdi.Client instance backed by the  http.Client instance
func New(uri, client_id, sceret, callback string, retry int) (*Client, error) {
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

	transport := &fdi.RetryTransport{
		Next:       http.DefaultTransport,
		MaxRetries: retry,
		Logger:     os.Stdout,
		Delay:      time.Duration(1 * time.Second),
		Source:     ContextTokenSource(),
	}

	httpClient := &http.Client{
		Transport: transport,
	}

	client := &Client{new(fdi.Client)}
	client.Client.Client = httpClient
	client.Client.BaseURL = base
	client.Client.ReportURL = report
	client.Client.Client_id = &client_id
	client.Client.Client_Sceret = &sceret

	tk, _, err := client.login(context.Background(), client_id, sceret)

	if err != nil {
		return nil, err
	}

	client.Client.Token = tk

	return client, nil
}

// NewDefault returns a new FDI API client using the`
// default "https://payments-api.fdibiz.com/v2" address.
func NewDefault(callback, client_id, sceret string) (*Client, error) {
	return New(baseUrl, client_id, sceret, callback, retry)

}

// do wraps the Client.Do function by creating the Request and
// unmarshalling the response according to user wish expected output.
func (c *Client) do(ctx context.Context, method, path string, in, out interface{}) (*fdi.Response, error) {
	req := &fdi.Request{
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
	if res.Status > 299 && res.Status < 499 {
		err := new(Err)
		_ = json.NewDecoder(res.Body).Decode(err)
		return res, &fdi.Error{Code: res.Status, Message: err.Data.Message}
	}

	if res.Status > 499 {
		return res, &fdi.Error{Code: res.Status, Message: "Something went wrong"}
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

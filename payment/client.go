package payment

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
)

var (
	// ErrNotFound indicates a resource is not found.
	ErrNotFound = errors.New("Not Found")

	// ErrNotSupported indicates a resource endpoint is not
	// supported or implemented.
	ErrNotSupported = errors.New("Not Supported")

	// ErrNotAuthorized indicates the request is not
	// authorized or the user does not have access to the
	// resource.
	ErrNotAuthorized = errors.New("Not Authorized")
)

// Request represents an HTTP request.
type Request struct {
	Method string
	Path   string
	Header http.Header
	Body   io.Reader
}

// Response represents an HTTP response.
type Response struct {
	ID     string
	Status int
	Header http.Header
	Body   io.ReadCloser
}

// Client manages communication with a payment gateways API.
type Client struct {
	Client *http.Client

	// Base URL for API requests.
	BaseURL *url.URL

	// ReportURL is the url to callback for payment reports
	ReportURL *url.URL

	// Driver identifies the payment provider to use
	Driver Driver
	// Payments pulls and pushes funds from/to the underlying payment Provider/Driver
	Payments PaymentsService

	// Balances services returns information about account balance balance/now
	Balances BalanceService

	// Info services returns information about transactions/payments
	Info InfoService
	// Auth authenticates our http client against the payment provider.
	Auth AuthService

	// DumpResponse optionally specifies a function to
	// dump the the response body for debugging purposes.
	// This can be set to httputil.DumpResponse.
	DumpResponse func(*http.Response, bool) ([]byte, error)
}

// Do sends an API request and returns the API response.
func (c *Client) Do(ctx context.Context, in *Request) (*Response, error) {
	uri, err := c.BaseURL.Parse(in.Path)
	if err != nil {
		return nil, err
	}

	// creates a new http request with context.
	req, err := http.NewRequest(in.Method, uri.String(), in.Body)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	if in.Header != nil {
		req.Header = in.Header
	}

	client := c.Client
	if client == nil {
		client = http.DefaultClient
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// dumps the response for debugging purposes.
	if c.DumpResponse != nil {
		_, _ = c.DumpResponse(res, true)
	}

	return newResponse(res), nil
}

// newResponse creates a new Response for the provided
// http.Response. r must not be nil.
func newResponse(r *http.Response) *Response {
	res := &Response{
		Status: r.StatusCode,
		Header: r.Header,
		Body:   r.Body,
	}
	return res
}

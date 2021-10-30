package payment

import (
	"context"
	"errors"
	"io"
	"net/http"
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

type ClientOption = func(http.RoundTripper) http.RoundTripper

type funcTripper struct {
	roundTrip func(*http.Request) (*http.Response, error)
}

func (tr funcTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return tr.roundTrip(req)
}

// CheckToken is a client option tha renew the token if it is expired
func CheckToken(client *Client, token *Token) ClientOption {
	return func(tr http.RoundTripper) http.RoundTripper {
		return &tokenRenewer{
			client: client,
			token:  token,
			rt:     tr,
		}
	}
}

// AddHeader turns a RoundTripper into one that adds a request header
func AddHeader(name, value string) ClientOption {
	return func(tr http.RoundTripper) http.RoundTripper {
		return &funcTripper{roundTrip: func(req *http.Request) (*http.Response, error) {
			if req.Header.Get(name) == "" {
				req.Header.Add(name, value)
			}
			return tr.RoundTrip(req)
		}}
	}
}

type tokenRenewer struct {
	token  *Token
	client *Client
	rt     http.RoundTripper
}

func (tr tokenRenewer) RoundTrip(req *http.Request) (*http.Response, error) {

	// check if token is expired
	// if expired, renew it
	// if not, use the token
	if tr.token == nil {
		token, err := tr.client.RefreshToken(req.Context(), tr.token)
		if err != nil {
			return nil, err
		}
		tr.token = token
	}
	if tr.token.IsExpired() {
		token, err := tr.client.RefreshToken(req.Context(), tr.token)
		if err != nil {
			return nil, err
		}
		tr.token = token
	}
	req.Header.Set("Authorization", "Bearer "+tr.token.Token)

	return tr.rt.RoundTrip(req)
}

// Refresh token
func (c *Client) RefreshToken(ctx context.Context, in *Token) (*Token, error) {
	return &Token{}, nil
}

func NewClient(opts ...ClientOption) *Client {
	client := &Client{http: NewHTTPClient(opts...)}
	return client
}

// NewHTTPClient initializes an http.Client
func NewHTTPClient(opts ...ClientOption) *http.Client {

	tr := http.DefaultTransport
	for _, opt := range opts {
		tr = opt(tr)
	}
	return &http.Client{Transport: tr}
}

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
	http *http.Client
}

// Do sends an API request and returns the API response.
func (c *Client) Do(ctx context.Context, in *Request) (*Response, error) {
	req, err := http.NewRequest(in.Method, "TODO", in.Body)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	if in.Header != nil {
		req.Header = in.Header
	}

	res, err := c.http.Do(req)
	if err != nil {
		return nil, err
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

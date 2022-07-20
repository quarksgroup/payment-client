//roundTripper this override the default RoundTripper of net/http package
package fdi

import (
	"io"
	"net/http"
	"time"
)

// Supported authentication schemes.
const (
	SchemeBearer = "Bearer"
	SchemeToken  = "token"
)

//ReteryTransport  is an http.RoundTrip that refreshes oauth
// tokens, wrapping a base ReteryTransport and refreshing the
// token if expired with max retry.
//request that require authorization header by appending header on it roundTripper
type ReteryTransport struct {
	Next       http.RoundTripper
	MaxRetries int
	Logger     io.Writer
	Delay      time.Duration // delay between each retry
	Source     TokenSource
	Scheme     string
}

func (t ReteryTransport) RoundTrip(req *http.Request) (*http.Response, error) {

	ctx := req.Context()

	token, err := t.Source.Token(ctx)

	if err != nil {
		return nil, err
	}

	if token == nil {
		return t.base().RoundTrip(req)
	}

	var attempt int

	for {

		req.Header.Set("Authorization", t.scheme()+" "+token.Token)

		res, err := t.Next.RoundTrip(req)

		attempt++

		// max retries exceeded
		if attempt == t.MaxRetries {
			return res, err
		}

		if err == nil && res.StatusCode == http.StatusOK {
			return res, err
		}

		//Check if request response is not authorized.
		if err == nil && res.StatusCode == http.StatusUnauthorized {
			//Here this is where we referesh our token to renew it
			res, err = t.Next.RoundTrip(req)
		}

		if err == nil && res.StatusCode == http.StatusInternalServerError {
			return res, err
		}

		if res.StatusCode >= 299 && res.StatusCode <= 400 {
			return res, err
		}

		// delay and retry
		select {
		case <-ctx.Done():
			return res, ctx.Err()
		case <-time.After(t.Delay):
		}
	}
}

func (t *ReteryTransport) base() http.RoundTripper {
	if t.Next != nil {
		return t.Next
	}
	return http.DefaultTransport
}

// scheme returns the token scheme. If no scheme is
// configured, the bearer scheme is used.
func (t *ReteryTransport) scheme() string {
	if t.Scheme == "" {
		return SchemeBearer
	}
	return t.Scheme
}

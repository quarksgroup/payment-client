//roundTripper this override the default RoundTripper of net/http package
package airtel

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"time"
)

// Supported authentication schemes.
const (
	SchemeBearer = "Bearer"
	SchemeToken  = "token"
)

//Transport  is an http.RoundTrip that refreshes oauth
// tokens, wrapping a base Transport and refreshing the
// token if expired with max retry.
//request that require authorization header by appending header on it roundTripper
type Transport struct {
	Next       http.RoundTripper
	MaxRetries int
	Logger     io.Writer
	Delay      time.Duration // delay between each retry
	Source     TokenSource
	Scheme     string
}

func (t Transport) RoundTrip(req *http.Request) (*http.Response, error) {

	fmt.Fprintf(t.Logger, "[%s] Takes %s time %s %s\n", time.Now().Format(time.ANSIC), duration(time.Since(time.Now()), 2), req.Method, req.URL.String())

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
			//Here this is where we referesh our token to renew it tthe implementation will go/called here
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

func (t *Transport) base() http.RoundTripper {
	if t.Next != nil {
		return t.Next
	}
	return http.DefaultTransport
}

// scheme returns the token scheme. If no scheme is
// configured, the bearer scheme is used.
func (t *Transport) scheme() string {
	if t.Scheme == "" {
		return SchemeBearer
	}
	return t.Scheme
}

//Duration calculate how long request takes...
func duration(d time.Duration, dicimal int) time.Duration {
	shift := int(math.Pow10(dicimal))

	units := []time.Duration{time.Second, time.Millisecond, time.Microsecond, time.Nanosecond}
	for _, u := range units {
		if d > u {
			div := u / time.Duration(shift)
			if div == 0 {
				break
			}
			d = d / div * div
			break
		}
	}

	return d
}

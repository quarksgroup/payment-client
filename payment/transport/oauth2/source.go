package oauth2

import (
	"context"

	"github.com/quarksgroup/payment-client/payment/airtel"
	"github.com/quarksgroup/payment-client/payment/fdi"
)

// StaticTokenSource returns a TokenSource that always
// returns the same token. Because the provided token t
// is never refreshed, StaticTokenSource is only useful
// for tokens that never expire.
func StaticTokenSource(t *fdi.Token) fdi.TokenSource {
	return staticTokenSource{t}
}

type staticTokenSource struct {
	token *fdi.Token
}

func (s staticTokenSource) Token(context.Context) (*fdi.Token, error) {
	return s.token, nil
}

// ContextTokenSource returns a TokenSource that returns
// a token from the http.Request context.
func ContextTokenSource() fdi.TokenSource {
	return contextTokenSource{}
}

type contextTokenSource struct {
}

func (s contextTokenSource) Token(ctx context.Context) (*fdi.Token, error) {
	token, _ := ctx.Value(fdi.TokenKey{}).(*fdi.Token)
	return token, nil
}

// StaticTokenSource...
func AirtelStaticTokenSource(t *airtel.Token) airtel.TokenSource {
	return airtelStaticTokenSource{t}
}

// airtelStaticTokenSource...
type airtelStaticTokenSource struct {
	token *airtel.Token
}

func (s airtelStaticTokenSource) Token(context.Context) (*airtel.Token, error) {
	return s.token, nil
}

func AirtelContextTokenSource() airtel.TokenSource {
	return airtelContextTokenSource{}
}

// airtelContextTokenSource...
type airtelContextTokenSource struct {
}

func (s airtelContextTokenSource) Token(ctx context.Context) (*airtel.Token, error) {
	token, _ := ctx.Value(fdi.TokenKey{}).(*airtel.Token)
	return token, nil
}

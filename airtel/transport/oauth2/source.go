package oauth2

import (
	"context"

	"github.com/quarksgroup/payment-client/airtel"
)

// StaticTokenSource returns a TokenSource that always
// returns the same token. Because the provided token t
// is never refreshed, StaticTokenSource is only useful
// for tokens that never expire.
func StaticTokenSource(t *airtel.Token) airtel.TokenSource {
	return staticTokenSource{t}
}

type staticTokenSource struct {
	token *airtel.Token
}

func (s staticTokenSource) Token(context.Context) (*airtel.Token, error) {
	return s.token, nil
}

// ContextTokenSource returns a TokenSource that returns
// a token from the http.Request context.
func ContextTokenSource() airtel.TokenSource {
	return contextTokenSource{}
}

type contextTokenSource struct {
}

func (s contextTokenSource) Token(ctx context.Context) (*airtel.Token, error) {
	token, _ := ctx.Value(airtel.TokenKey{}).(*airtel.Token)
	return token, nil
}

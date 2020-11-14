package oauth2

import (
	"context"

	"github.com/quarksgroup/payment-client/payment"
)

// StaticTokenSource returns a TokenSource that always
// returns the same token. Because the provided token t
// is never refreshed, StaticTokenSource is only useful
// for tokens that never expire.
func StaticTokenSource(t *payment.Token) payment.TokenSource {
	return staticTokenSource{t}
}

type staticTokenSource struct {
	token *payment.Token
}

func (s staticTokenSource) Token(context.Context) (*payment.Token, error) {
	return s.token, nil
}

// ContextTokenSource returns a TokenSource that returns
// a token from the http.Request context.
func ContextTokenSource() payment.TokenSource {
	return contextTokenSource{}
}

type contextTokenSource struct {
}

func (s contextTokenSource) Token(ctx context.Context) (*payment.Token, error) {
	token, _ := ctx.Value(payment.TokenKey{}).(*payment.Token)
	return token, nil
}

package fdi

import (
	"context"

	"github.com/quarksgroup/payment-client/fdi"
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

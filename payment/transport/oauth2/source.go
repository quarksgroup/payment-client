package oauth2

import (
	"context"

	"github.com/iradukunda1/payment-staging/payment/mtn"
)

// StaticTokenSource returns a TokenSource that always
// returns the same token. Because the provided token t
// is never refreshed, StaticTokenSource is only useful
// for tokens that never expire.
func StaticTokenSource(t *mtn.Token) mtn.TokenSource {
	return staticTokenSource{t}
}

type staticTokenSource struct {
	token *mtn.Token
}

func (s staticTokenSource) Token(context.Context) (*mtn.Token, error) {
	return s.token, nil
}

// ContextTokenSource returns a TokenSource that returns
// a token from the http.Request context.
func ContextTokenSource() mtn.TokenSource {
	return contextTokenSource{}
}

type contextTokenSource struct {
}

func (s contextTokenSource) Token(ctx context.Context) (*mtn.Token, error) {
	token, _ := ctx.Value(mtn.TokenKey{}).(*mtn.Token)
	return token, nil
}

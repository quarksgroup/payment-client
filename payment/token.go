package payment

import (
	"context"
	"time"
)

// Token represents the credentials used to authorize
// the requests to access protected resources.
type Token struct {
	Token   string
	Expires time.Time
}

// IsExpired returns true if the token is expired.
func (t *Token) IsExpired() bool {
	return t.Expires.Before(time.Now())
}

// TokenKey is the key to use with the context.WithValue
// function to associate an Token value with a context.
type TokenKey struct{}

// TokenSource returns a token.
type TokenSource interface {
	Token(context.Context) (*Token, error)
}

// WithContext returns a copy of parent in which the token value is set
func WithContext(parent context.Context, token *Token) context.Context {
	return context.WithValue(parent, TokenKey{}, token)
}

// TokenFrom returns the login token rom the context.
func TokenFrom(ctx context.Context) *Token {
	token, _ := ctx.Value(TokenKey{}).(*Token)
	return token
}

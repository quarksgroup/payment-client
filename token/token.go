package token

import (
	"context"
	"time"
)

// Token represents the credentials used to authorize
// the requests to access protected resources.
type Token struct {
	Token   string    `json:"token"`
	Type    string    `json:"token_type"`
	Expires time.Time `json:"expires"`
}

// TokenSource is an interface that returns a token or an error.
// Use TokenSourceContext to adapt TokenSource to use a context.
type TokenSource interface {
	Token(ctx context.Context) (*Token, error)
}

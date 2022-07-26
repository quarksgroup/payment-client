package airtel

import (
	"context"
)

// Token represents the credentials used to authorize
// the requests to access protected resources.
type Token struct {
	Token   string `json:"token"`
	Type    string `json:"token_type"`
	Expires int64  `json:"expires"`
}

// TokenKey is the key to use with the context.WithValue
// function to associate an Token value with a context.
type TokenKey struct{}

// TokenSource returns a token.
type TokenSource interface {
	Token(context.Context) (*Token, error)
}

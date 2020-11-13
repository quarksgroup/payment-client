package payment

import (
	"context"
	"time"
)

// Token represents the credentials used to authorize
// the requests to access protected resources.
type Token struct {
	Token   string
	Refresh string
	Expires time.Time
}

// AuthService handles authentication to the underlying API
type AuthService interface {
	// Login with id and secret to the underlying API and get an JWT token
	Login(context.Context, string, string) (*Token, *Response, error)

	// Refresh the oauth2 token
	Refresh(ctx context.Context, token *Token, force bool) (*Token, *Response, error)
}

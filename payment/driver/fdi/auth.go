package fdi

import (
	"context"
	"errors"
	"time"

	"github.com/iradukunda1/payment-staging/payment/mtn"
)

type authService struct {
	client *wrapper
}

func (s *authService) Login(ctx context.Context, id, secret string) (*mtn.Token, *mtn.Response, error) {
	endpoint := "auth"
	in := tokenRequest{
		App:    id,
		Secret: secret,
	}
	out := new(tokenResponse)
	res, err := s.client.do(ctx, "POST", endpoint, in, out)
	return convertToken(out), res, err
}

func (s *authService) Refresh(ctx context.Context, token *mtn.Token, force bool) (*mtn.Token, *mtn.Response, error) {
	return nil, nil, errors.New("refresh not implemented")
}

type tokenRequest struct {
	App    string `json:"appId"`
	Secret string `json:"secret"`
}

type tokenResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Token     string `json:"token"`
		ExpiresAt string `json:"expires_at"`
	}
}

func convertToken(tk *tokenResponse) *mtn.Token {

	expires, _ := time.Parse("2006-01-02T15:04:05.99999", tk.Data.ExpiresAt)

	return &mtn.Token{
		Token:   tk.Data.Token,
		Expires: expires,
	}
}

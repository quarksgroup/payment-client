package fdi

import (
	"context"
	"errors"

	"github.com/quarksgroup/payment-client/payment"
)

type authService struct {
	client *wrapper
}

func (s *authService) Login(ctx context.Context, id, secret string) (*payment.Token, *payment.Response, error) {
	endpoint := "auth"
	in := tokenRequest{
		App:    id,
		Secret: secret,
	}
	out := new(tokenResponse)
	res, err := s.client.do(ctx, "POST", endpoint, in, out)
	return convertToken(out), res, err
}

func (s *authService) Refresh(ctx context.Context, token *payment.Token, force bool) (*payment.Token, *payment.Response, error) {
	return nil, nil, errors.New("refresh not implemented")
}

type tokenRequest struct {
	App    string `json:"appId"`
	Secret string `json:"secret"`
}

type tokenResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Token string `json:"token"`
	}
}

func convertToken(tk *tokenResponse) *payment.Token {
	return &payment.Token{
		Token: tk.Data.Token,
	}
}

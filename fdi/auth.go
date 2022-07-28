package fdi

import (
	"context"
	"time"

	"github.com/quarksgroup/payment-client/token"
)

type tokenSource struct {
	token  *token.Token
	client *Client
	id     string
	secret string
}

func newTokenSource(client *Client, cfg *Config) (token.TokenSource, error) {
	tks := &tokenSource{
		client: client,
		id:     cfg.ClientId,
		secret: cfg.Secret,
	}
	token, err := tks.Login(context.Background(), cfg.ClientId, cfg.Secret)
	if err != nil {
		return nil, err
	}
	tks.token = token
	return tks, nil
}

func (tk *tokenSource) Token(ctx context.Context) (*token.Token, error) {
	if tk.token != nil {
		if tk.token.Expires.Before(time.Now().Local()) {
			token, err := tk.Login(ctx, tk.id, tk.secret)
			if err != nil {
				return nil, err
			}
			tk.token = token
		}
		return tk.token, nil

	}
	return tk.Login(ctx, tk.id, tk.secret)

}

func (tk *tokenSource) Login(ctx context.Context, id, secret string) (*token.Token, error) {
	endpoint := "auth"
	in := tokenRequest{
		App:    id,
		Secret: secret,
	}
	out := new(tokenResponse)
	_, err := tk.client.do(ctx, "POST", endpoint, in, out, false)
	return convertToken(out), err
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

func convertToken(tk *tokenResponse) *token.Token {

	expires, _ := time.Parse("2006-01-02T15:04:05.99999", tk.Data.ExpiresAt)

	return &token.Token{
		Token:   tk.Data.Token,
		Expires: expires,
	}
}

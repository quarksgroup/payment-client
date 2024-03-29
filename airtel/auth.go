package airtel

import (
	"context"
	"time"

	"github.com/quarksgroup/payment-client/token"
)

type tokenSource struct {
	token     *token.Token
	client    *Client
	id        string
	secret    string
	grantType string
}

func newTokenSource(client *Client, cfg *Config) (token.TokenSource, error) {
	tks := &tokenSource{
		client:    client,
		id:        cfg.ClientId,
		secret:    cfg.Secret,
		grantType: cfg.Grant,
	}

	token, err := tks.Login(context.Background(), cfg.ClientId, cfg.Secret, cfg.Grant)

	if err != nil {
		return nil, err
	}

	tks.token = token

	return tks, nil
}

func (tk *tokenSource) Token(ctx context.Context) (*token.Token, error) {
	if tk.token != nil {
		if tk.token.Expires.Before(time.Now().Local()) {

			token, err := tk.Login(ctx, tk.id, tk.secret, tk.grantType)
			if err != nil {
				return nil, err
			}
			tk.token = token
		}
		return tk.token, nil

	}
	return tk.Login(ctx, tk.id, tk.secret, tk.grantType)

}

//login responsible client api authentication to your airtel account portal but doesn't exposed outside
func (tk tokenSource) Login(ctx context.Context, id, secret, grantType string) (*token.Token, error) {
	endpoint := "auth/oauth2/token"
	in := tokenRequest{
		ClientId:     id,
		ClientSecret: secret,
		GrantType:    grantType,
	}
	out := new(tokenResponse)
	_, err := tk.client.do(ctx, "POST", endpoint, in, out, nil, false)
	return convertToken(out), err
}

type tokenRequest struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Expires     int64  `json:"expires_in"`
}

func convertToken(tk *tokenResponse) *token.Token {

	expires := time.Now().Local().Add(time.Duration(tk.Expires) * time.Second)

	return &token.Token{
		Token:   tk.AccessToken,
		Type:    tk.TokenType,
		Expires: expires,
	}
}

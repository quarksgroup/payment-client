package client

import (
	"context"

	"github.com/quarksgroup/payment-client/airtel"
)

//login responsible client api authentication to your airtel account portal but doesn't exposed outside
func (c *Client) login(ctx context.Context, id, secret, grantType string) (*airtel.Token, *airtel.Response, error) {
	endpoint := "auth/oauth2/token"
	in := tokenRequest{
		ClientId:     id,
		ClientSecret: secret,
		GrantType:    grantType,
	}
	out := new(tokenResponse)
	res, err := c.do(ctx, "POST", endpoint, in, out, nil)
	return convertToken(out), res, err
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

func convertToken(tk *tokenResponse) *airtel.Token {

	// expires := time.Now().Local().Add(180 * time.Second)

	return &airtel.Token{
		Token:   tk.AccessToken,
		Type:    tk.TokenType,
		Expires: tk.Expires,
	}
}

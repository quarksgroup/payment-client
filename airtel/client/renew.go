package client

import (
	"context"
	"fmt"
	"time"
)

//renewToken responsible to check client token and renew it if is expired
func (c *Client) renewToken(ctx context.Context) error {

	fmt.Println(c.Token.Token)

	//Check if the token is still valid and refresh it if expired
	if !time.Unix(c.Token.Expires, 0).After(time.Now().UTC().Add(-time.Second * 30)) {

		_, _, err := c.login(ctx, *c.ClientId, *c.ClientSceret, *c.GrantType)

		if err != nil {
			return err
		}
	}

	return nil
}

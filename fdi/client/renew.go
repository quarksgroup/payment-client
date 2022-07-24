package client

import (
	"context"
	"time"
)

const timeDuration = time.Hour // 1hour

//renewToken responsible to check client token and renew it if is expired
func (c *Client) renewToken(ctx context.Context) error {
	//Check if the token is still valid and refresh it if expired
	if !time.Unix(c.Token.Expires.Unix(), 0).After(time.Now().UTC().Add(-timeDuration)) {

		_, _, err := c.login(ctx, *c.Client_id, *c.Client_Sceret)

		if err != nil {
			return err
		}
	}

	return nil
}

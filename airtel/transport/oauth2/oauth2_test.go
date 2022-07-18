package oauth2

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/quarksgroup/payment-client/airtel"
	"gopkg.in/h2non/gock.v1"
)

func TestTransport(t *testing.T) {
	defer gock.Off()

	gock.New("https://openapi.airtel.africa").
		Get("/standard/v1/users/balance").
		MatchHeader("Authorization", "Bearer mF_9.B5f-4.1JqM").
		Reply(200)

	client := &http.Client{
		Transport: &Transport{
			Source: StaticTokenSource(
				&airtel.Token{
					Token: "mF_9.B5f-4.1JqM",
				},
			),
		},
	}

	res, err := client.Get("https://openapi.airtel.africa/standard/v1/users/balance")
	if err != nil {
		t.Error(err)
		return
	}
	defer res.Body.Close()
}

func TestTransport_TokenError(t *testing.T) {
	want := errors.New("Cannot retrieve token")
	client := &http.Client{
		Transport: &Transport{
			Source: mockErrorSource{want},
		},
	}

	_, err := client.Get("https://openapi.airtel.africa/standard/v1/users/balance")
	if err == nil {
		t.Errorf("Expect token source error, got nil")
	}
}

type mockErrorSource struct {
	err error
}

func (s mockErrorSource) Token(ctx context.Context) (*airtel.Token, error) {
	return nil, s.err
}

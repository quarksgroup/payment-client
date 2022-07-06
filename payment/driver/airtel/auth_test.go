package airtel

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/iradukunda1/payment-staging/payment/airtel"
	"github.com/stretchr/testify/require"
	"gopkg.in/h2non/gock.v1"
)

func TestLogin(t *testing.T) {
	defer gock.Off()

	gock.New(baseUrl).
		Post("/auth/oauth2/token").
		Reply(200).
		Type("application/json").
		File("testdata/auth.json")
	client := NewDefault("encrypted-pin")

	got, _, err := client.Auth.Login(context.Background(), "id", "secret", "grant_type")

	require.Nil(t, err, fmt.Sprintf("unexpected error %v", err))

	want := new(airtel.Token)
	raw, _ := ioutil.ReadFile("testdata/auth.json.golden")
	_ = json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

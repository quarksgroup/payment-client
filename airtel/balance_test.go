package airtel

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/quarksgroup/payment-client/mock"
	"github.com/stretchr/testify/require"
	"gopkg.in/h2non/gock.v1"
)

func TestBalance(t *testing.T) {
	defer gock.Off()
	gock.Observe(gock.DumpRequest)

	gock.New(baseUrl).
		Get("/standard/v1/users/balance").
		Reply(200).
		Type("application/json").
		File("testdata/account.json")

	cfg := &Config{
		ClientId: "client_id",
		Secret:   "client_secret",
		Grant:    "client_credentials",
		Pin:      "pin",
		Currency: "RWF",
		Country:  "RW",
	}
	tokenSource := mock.NewMockTokenSource()

	client, err := New(cfg, tokenSource, baseUrl, true, defaultRetries, nil)

	require.Nil(t, err, fmt.Sprintf("unexpected error %v", err))

	got, _, err := client.Balance(context.Background())

	require.Nil(t, err, fmt.Sprintf("unexpected error %v", err))

	want := new(Balance)
	raw, _ := ioutil.ReadFile("testdata/account.json.golden")
	_ = json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

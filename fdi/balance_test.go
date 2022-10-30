package fdi

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

	//For debugging
	gock.Observe(gock.DumpRequest)

	gock.New("https://payments-api.fdibiz.com/v2").
		Get("/balance/now").
		Reply(200).
		Type("application/json").
		File("testdata/balance.json")

	cfg := &Config{
		ClientId: "client_id",
		Secret:   "client_secret",
		CallBack: "https://test-callback.io",
	}
	tokenSource := mock.NewMockTokenSource()

	client, err := New(baseUrl, cfg, tokenSource, retry, nil)

	require.Nil(t, err, fmt.Sprintf("client initialization error %v", err))

	got, _, err := client.Balance(context.Background())

	require.Nil(t, err, fmt.Sprintf("unexpected error %v", err))

	want := new(Balance)
	raw, _ := ioutil.ReadFile("testdata/balance.json.golden")
	_ = json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

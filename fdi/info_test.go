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

var ref = "xxxxx"

func TestInfo(t *testing.T) {
	defer gock.Off()

	gock.New("https://payments-api.fdibiz.com/v2").
		Get(fmt.Sprintf("/momo/trx/%s/info", ref)).
		Reply(200).
		Type("application/json").
		File("testdata/info.json")

	cfg := &Config{
		ClientId: "client_id",
		Secret:   "client_secret",
		CallBack: "https://test-callback.io",
	}
	tokenSource := mock.NewMockTokenSource()

	client, err := New(baseUrl, cfg, tokenSource, retry, nil)

	require.Nil(t, err, fmt.Sprintf("client initialization error %v", err))

	got, _, err := client.TransactionInfo(context.Background(), ref)

	require.Nil(t, err, fmt.Sprintf("unexpected error %v", err))

	want := new(Info)
	raw, _ := ioutil.ReadFile("testdata/info.json.golden")
	_ = json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

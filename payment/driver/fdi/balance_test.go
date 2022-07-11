package fdi

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/quarksgroup/payment-client/payment/fdi"
	"github.com/stretchr/testify/require"
	"gopkg.in/h2non/gock.v1"
)

func TestBalance(t *testing.T) {
	defer gock.Off()

	gock.New("https://payments-api.fdibiz.com/v2").
		Get("/balance/now").
		Reply(200).
		Type("application/json").
		File("testdata/balance.json")
	client := NewDefault("https://test-callback.io")

	got, _, err := client.Balances.Balance(context.Background())

	require.Nil(t, err, fmt.Sprintf("unexpected error %v", err))

	want := new(fdi.Balance)
	raw, _ := ioutil.ReadFile("testdata/balance.json.golden")
	_ = json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/quarksgroup/payment-client/airtel"
	"github.com/stretchr/testify/require"
	"gopkg.in/h2non/gock.v1"
)

func TestBalance(t *testing.T) {
	defer gock.Off()

	gock.New(baseUrl).
		Get("/standard/v1/users/balance").
		Reply(200).
		Type("application/json").
		File("testdata/account.json")
	client, err := NewDefault("encrypted-pin", "client_id", "sceret", "grant_type")

	require.Nil(t, err, fmt.Sprintf("unexpected error %v", err))

	got, _, err := client.Balance(context.Background())

	require.Nil(t, err, fmt.Sprintf("unexpected error %v", err))

	want := new(airtel.Balance)
	raw, _ := ioutil.ReadFile("testdata/account.json.golden")
	_ = json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

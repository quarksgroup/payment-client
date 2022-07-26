package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/quarksgroup/payment-client/airtel"
	"github.com/quarksgroup/payment-client/mock"
	"github.com/stretchr/testify/require"
	"gopkg.in/h2non/gock.v1"
)

var num = "72xxxxx"

func TestCheck(t *testing.T) {
	defer gock.Off()
	gock.Observe(gock.DumpRequest)

	gock.New(baseUrl).
		Get("/standard/v1/users/").
		Reply(200).
		Type("application/json").
		File("testdata/check_number.json")

	cfg := &Config{
		ClientId: "client_id",
		Secret:   "client_secret",
		Grant:    "client_credentials",
		Pin:      "pin",
		Currency: "RWF",
		Country:  "RW",
	}
	tokenSource := mock.NewMockTokenSource()

	client, err := New(cfg, tokenSource, baseUrl, true, defaultRetries)

	require.Nil(t, err, fmt.Sprintf("unexpected error %v", err))

	got, _, err := client.NumberInfo(context.Background(), num)

	require.Nil(t, err, fmt.Sprintf("unexpected error %v", err))

	want := new(airtel.Number)

	raw, _ := ioutil.ReadFile("testdata/check_number.json.golden")

	_ = json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

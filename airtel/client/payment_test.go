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

func TestPush(t *testing.T) {
	defer gock.Off()
	gock.Observe(gock.DumpRequest)

	in := &airtel.Payment{
		ID:     "xxxx",
		Amount: 100,
		Ref:    "xxxx",
		Phone:  num,
	}

	authClientMock()

	gock.New(baseUrl).
		Post("/standard/v1/disbursements/").
		Reply(200).
		Type("application/json").
		File("testdata/push.json")

	client, err := NewDefault("encrypted-pin", "client_id", "sceret", "grant_type")

	require.Nil(t, err, fmt.Sprintf("unexpected error %v", err))

	got, _, err := client.Push(context.Background(), in)

	require.Nil(t, err, fmt.Sprintf("unexpected error %v", err))

	want := new(airtel.Status)
	raw, _ := ioutil.ReadFile("testdata/push.json.golden")
	_ = json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestPull(t *testing.T) {
	defer gock.Off()
	gock.Observe(gock.DumpRequest)

	in := &airtel.Payment{
		ID:     "xxxx",
		Amount: 100,
		Ref:    "xxxx",
		Phone:  num,
	}

	authClientMock()

	gock.New(baseUrl).
		Post("/merchant/v1/payments/").
		Reply(200).
		Type("application/json").
		File("testdata/pull.json")

	client, err := NewDefault("encrypted-pin", "client_id", "sceret", "grant_type")

	require.Nil(t, err, fmt.Sprintf("unexpected error %v", err))

	got, _, err := client.Pull(context.Background(), in)

	require.Nil(t, err, fmt.Sprintf("unexpected error %v", err))

	want := new(airtel.Status)

	raw, _ := ioutil.ReadFile("testdata/pull.json.golden")

	_ = json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

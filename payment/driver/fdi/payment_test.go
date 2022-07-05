package fdi

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/quarksgroup/payment-client/payment/mtn"
	"github.com/stretchr/testify/require"
	"gopkg.in/h2non/gock.v1"
)

func TestPull(t *testing.T) {
	defer gock.Off()

	py := &mtn.Payment{
		ID:       "xxxx",
		Amount:   1000,
		Wallet:   "xxxx",
		Provider: MTN,
	}

	gock.New("https://payments-api.fdibiz.com/v2").
		Post("/momo/pull").
		Reply(200).
		Type("application/json").
		File("testdata/status.json")
	client := NewDefault("https://test-callback.io")

	got, _, err := client.Payments.Pull(context.Background(), py)

	require.Nil(t, err, fmt.Sprintf("unexpected error %v", err))

	want := new(mtn.Status)
	raw, _ := ioutil.ReadFile("testdata/status.json.golden")
	_ = json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func TestPush(t *testing.T) {
	defer gock.Off()

	py := &mtn.Payment{
		ID:       "xxxx",
		Amount:   1000,
		Wallet:   "xxxx",
		Provider: MTN,
	}

	gock.New("https://payments-api.fdibiz.com/v2").
		Post("/momo/push").
		Reply(202).
		Type("application/json").
		File("testdata/status.json")
	client := NewDefault("https://test-callback.io")

	got, _, err := client.Payments.Push(context.Background(), py)

	require.Nil(t, err, fmt.Sprintf("unexpected error %v", err))

	want := new(mtn.Status)
	raw, _ := ioutil.ReadFile("testdata/status.json.golden")
	_ = json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

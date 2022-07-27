package fdi

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/nbio/st"
	"github.com/quarksgroup/payment-client/mock"
	"github.com/stretchr/testify/require"
	"gopkg.in/h2non/gock.v1"
)

func TestPull(t *testing.T) {
	defer gock.Off()

	//For debugging proposer mode
	gock.Observe(gock.DumpRequest)

	py := &Payment{
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

	cfg := &Config{
		ClientId: "client_id",
		Secret:   "client_secret",
		CallBack: "https://test-callback.io",
	}
	tokenSource := mock.NewMockTokenSource()

	client, err := New(baseUrl, cfg, tokenSource, retry)

	require.Nil(t, err, fmt.Sprintf("client initialization error %v", err))

	got, _, err := client.Pull(context.Background(), py)

	require.Nil(t, err, fmt.Sprintf("unexpected error %v", err))

	want := new(Status)
	raw, _ := ioutil.ReadFile("testdata/status.json.golden")
	_ = json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	// Verify that we don't have pending mocks
	st.Expect(t, gock.IsDone(), true)
}

func TestPush(t *testing.T) {
	defer gock.Off()

	//For debugging proposer mode
	gock.Observe(gock.DumpRequest)

	py := &Payment{
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

	cfg := &Config{
		ClientId: "client_id",
		Secret:   "client_secret",
		CallBack: "https://test-callback.io",
	}
	tokenSource := mock.NewMockTokenSource()

	client, err := New(baseUrl, cfg, tokenSource, retry)

	require.Nil(t, err, fmt.Sprintf("client initialization error %v", err))

	got, _, err := client.Push(context.Background(), py)

	require.Nil(t, err, fmt.Sprintf("unexpected error %v", err))

	want := new(Status)
	raw, _ := ioutil.ReadFile("testdata/status.json.golden")
	_ = json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

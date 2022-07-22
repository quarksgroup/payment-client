package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/quarksgroup/payment-client/fdi"
	"github.com/stretchr/testify/require"
	"gopkg.in/h2non/gock.v1"
)

func TestLogin(t *testing.T) {
	defer gock.Off()
	gock.Observe(gock.DumpRequest)

	gock.New("https://payments-api.fdibiz.com/v2").
		Post("/auth").
		Reply(200).
		Type("application/json").
		File("testdata/token.json")

	gock.New("https://payments-api.fdibiz.com/v2").
		Post("/auth").
		Reply(200).
		Type("application/json").
		File("testdata/token.json")

	client, err := NewDefault("https://test-callback.io", "client_id", "screte")

	require.Nil(t, err, fmt.Sprintf("client initialization error %v", err))

	got, _, err := client.login(context.Background(), "id", "secret")

	require.Nil(t, err, fmt.Sprintf("unexpected error %v", err))

	want := new(fdi.Token)
	raw, _ := ioutil.ReadFile("testdata/token.json.golden")
	_ = json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

func AuthClientMock() {
	gock.New("https://payments-api.fdibiz.com/v2").
		Post("/auth").
		Reply(200).
		Type("application/json").
		File("testdata/token.json")

	gock.New("https://payments-api.fdibiz.com/v2").
		Post("/auth").
		Reply(200).
		Type("application/json").
		File("testdata/token.json")
}

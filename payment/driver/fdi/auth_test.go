package fdi

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/iradukunda1/payment-staging/payment/mtn"
	"github.com/stretchr/testify/require"
	"gopkg.in/h2non/gock.v1"
)

func TestLogin(t *testing.T) {
	defer gock.Off()

	gock.New("https://payments-api.fdibiz.com/v2").
		Post("/auth").
		Reply(200).
		Type("application/json").
		File("testdata/token.json")
	client := NewDefault("https://test-callback.io")

	got, _, err := client.Auth.Login(context.Background(), "id", "secret")

	require.Nil(t, err, fmt.Sprintf("unexpected error %v", err))

	want := new(mtn.Token)
	raw, _ := ioutil.ReadFile("testdata/token.json.golden")
	_ = json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

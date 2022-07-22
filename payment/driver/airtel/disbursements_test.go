package airtel

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/iradukunda1/payment-staging/payment/airtel"
	"github.com/stretchr/testify/require"
	"gopkg.in/h2non/gock.v1"
)

func TestPull(t *testing.T) {
	defer gock.Off()

	in := &airtel.PaymentReq{
		ID:     "xxxx",
		Amount: 1000,
		Ref:    "xxxx",
		Phone:  num,
	}

	gock.New(baseUrl).
		Post("/standard/v1/disbursements/").
		Reply(200).
		Type("application/json").
		File("testdata/pull.json")
	client := NewDefault("encrypted-pin")

	got, _, err := client.Disbursement.Pull(context.Background(), in)

	require.Nil(t, err, fmt.Sprintf("unexpected error %v", err))

	want := new(airtel.PaymentResp)
	raw, _ := ioutil.ReadFile("testdata/pull.json.golden")
	_ = json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

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

var num = "72xxxxx"

func TestCheck(t *testing.T) {
	defer gock.Off()

	gock.New(baseUrl).
		Get("/standard/v1/users/").
		Reply(200).
		Type("application/json").
		File("testdata/check_number.json")
	client := NewDefault("encrypted-pin")

	got, _, err := client.CheckNumber.Check(context.Background(), num)

	require.Nil(t, err, fmt.Sprintf("unexpected error %v", err))

	want := new(airtel.Number)

	raw, _ := ioutil.ReadFile("testdata/check_number.json.golden")

	_ = json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}
}

package fdi

import (
	"testing"
)

func TestBalance(t *testing.T) {
	// defer gock.Off()

	// gock.New("https://payments-api.fdibiz.com/v2").
	// 	Get("/balance/now").
	// 	Reply(200).
	// 	Type("application/json").
	// 	File("testdata/balance.json")
	// client := NewDefault("https://test-callback.io")

	// got, _, err := client.Balances.Balance(context.Background())

	// require.Nil(t, err, fmt.Sprintf("unexpected error %v", err))

	// want := new(payment.Balance)
	// raw, _ := ioutil.ReadFile("testdata/balance.json.golden")
	// _ = json.Unmarshal(raw, want)

	// if diff := cmp.Diff(got, want); diff != "" {
	// 	t.Errorf("Unexpected Results")
	// 	t.Log(diff)
	// }
}

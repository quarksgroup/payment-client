package fdi

import (
	"testing"
)

var ref = "xxxxx"

func TestInfo(t *testing.T) {
	// defer gock.Off()

	// gock.New("https://payments-api.fdibiz.com/v2").
	// 	Get(fmt.Sprintf("/momo/trx/%s/info", ref)).
	// 	Reply(200).
	// 	Type("application/json").
	// 	File("testdata/info.json")
	// client := NewDefault("https://test-callback.io")

	// got, _, err := client.Info.Info(context.Background(), ref)

	// require.Nil(t, err, fmt.Sprintf("unexpected error %v", err))

	// want := new(payment.Info)
	// raw, _ := ioutil.ReadFile("testdata/info.json.golden")
	// _ = json.Unmarshal(raw, want)

	// if diff := cmp.Diff(got, want); diff != "" {
	// 	t.Errorf("Unexpected Results")
	// 	t.Log(diff)
	// }
}

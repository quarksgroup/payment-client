package fdi

import (
	"testing"
)

func TestLogin(t *testing.T) {
	// defer gock.Off()

	// gock.New("https://payments-api.fdibiz.com/v2").
	// 	Post("/auth").
	// 	Reply(200).
	// 	Type("application/json").
	// 	File("testdata/token.json")
	// client := NewDefault("https://test-callback.io")

	// got, _, err := client.Auth.Login(context.Background(), "id", "secret")

	// require.Nil(t, err, fmt.Sprintf("unexpected error %v", err))

	// want := new(payment.Token)
	// raw, _ := ioutil.ReadFile("testdata/token.json.golden")
	// _ = json.Unmarshal(raw, want)

	// if diff := cmp.Diff(got, want); diff != "" {
	// 	t.Errorf("Unexpected Results")
	// 	t.Log(diff)
	// }
}

package fdi

import (
	"testing"

	"gopkg.in/h2non/gock.v1"
)

func TestLogin(t *testing.T) {
	defer gock.Off()
	gock.Observe(gock.DumpRequest)

	t.Skip()
}

package main

import (
	"testing"
)

func TestTokenRCContents(t *testing.T) {
	conf := Config{
		Registry: "https://npm.someorg.com/",
		Token:    "token",
	}
	actual := npmrcContentsToken(conf)
	expected := "//npm.someorg.com/:_authToken=token"
	if actual != expected {
		t.Errorf("Unexpected token config (Got: %s, Expected: %s)", actual, expected)
	}

	conf.Registry = "https://npm.someorg.com/with/path/"
	actual = npmrcContentsToken(conf)
	expected = "//npm.someorg.com/with/path/:_authToken=token"
	if actual != expected {
		t.Errorf("Unexpected token config (Got: %s, Expected: %s)", actual, expected)
	}
}

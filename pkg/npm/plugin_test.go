// Copyright (c) 2019, the Drone Plugins project authors.
// Please see the AUTHORS file for details. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package npm

import "testing"

func TestTokenRCContents(t *testing.T) {
	settings := Settings{
		Registry: "https://npm.someorg.com/",
		Token:    "token",
	}
	actual := npmrcContentsToken(settings)
	expected := "//npm.someorg.com/:_authToken=token"
	if actual != expected {
		t.Errorf("Unexpected token settings (Got: %s, Expected: %s)", actual, expected)
	}

	settings.Registry = "https://npm.someorg.com/with/path/"
	actual = npmrcContentsToken(settings)
	expected = "//npm.someorg.com/with/path/:_authToken=token"
	if actual != expected {
		t.Errorf("Unexpected token settings (Got: %s, Expected: %s)", actual, expected)
	}

	settings.Registry = globalRegistry
	actual = npmrcContentsToken(settings)
	expected = "//registry.npmjs.org/:_authToken=token"
	if actual != expected {
		t.Errorf("Unexpected token settings (Got: %s, Expected: %s)", actual, expected)
	}

	settings.Registry = "https://npm.someorg.com"
	actual = npmrcContentsToken(settings)
	expected = "//npm.someorg.com/:_authToken=token"
	if actual != expected {
		t.Errorf("Unexpected token settings (Got: %s, Expected: %s)", actual, expected)
	}

	settings.Registry = "https://npm.someorg.com/with/path"
	actual = npmrcContentsToken(settings)
	expected = "//npm.someorg.com/with/path/:_authToken=token"
	if actual != expected {
		t.Errorf("Unexpected token settings (Got: %s, Expected: %s)", actual, expected)
	}
}

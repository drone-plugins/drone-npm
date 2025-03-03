// Copyright (c) 2020, the Drone Plugins project authors.
// Please see the AUTHORS file for details. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package plugin

import (
	"context"
	"net/url"
	"testing"

	"github.com/drone-plugins/drone-plugin-lib/drone"
	"github.com/stretchr/testify/assert"
)

func initFakeSettings() Settings {
	nc := npmConfig{
		// Note: this registry is the one that would come from publishConfig in package.json
		Registry: "SetByFn_readPackageFile",
	}
	np := npmPackage{
		Name:    "Test Package",
		Version: "1.33.7",
		Config:  nc,
	}
	return Settings{
		Username:   "fakeUser",
		Password:   "fakePass",
		Token:      "",
		SkipWhoami: false,
		Email:      "fake@user.tst",
		// Note: this registry is the one that would come from drone yaml
		Registry:               "https://fakenpm.reg.org/good/path",
		Folder:                 "__test__",
		FailOnVersionConflict:  true,
		Tag:                    "",
		Access:                 "",
		SkipRegistryValidation: false,
		npm:                    &np,
	}
}

func initFakeNetwork() drone.Network {
	return drone.Network{
		SkipVerify: true,
		Client:     nil,
		Context:    context.TODO(),
	}
}

func initFakePipeline() drone.Pipeline {
	return drone.Pipeline{
		Build:  drone.Build{},
		Repo:   drone.Repo{},
		Commit: drone.Commit{},
		Stage:  drone.Stage{},
		Step:   drone.Step{},
		SemVer: drone.SemVer{},
		CalVer: drone.CalVer{},
		System: drone.System{},
	}
}

func initPlugin() *Plugin {
	return &Plugin{
		settings: initFakeSettings(),
		pipeline: initFakePipeline(),
		network:  initFakeNetwork(),
	}
}

func getParsedURI(s string) *url.URL {
	rslt, _ := url.Parse(s)
	return rslt
}

func TestIsDefaultOrNilPort(t *testing.T) {
	p := initPlugin()

	resultWithoutPort := isNilPortOrStandardSchemePort(getParsedURI(p.settings.Registry))
	assert.Equal(t, true, resultWithoutPort)

	p.settings.Registry = "https://fakenpm.reg.org:443"
	resultWithPort := isNilPortOrStandardSchemePort(getParsedURI(p.settings.Registry))
	assert.Equal(t, true, resultWithPort)

	p.settings.Registry = "http://fakenpm.reg.org:80"
	resultWithPortHTTP := isNilPortOrStandardSchemePort(getParsedURI(p.settings.Registry))
	assert.Equal(t, true, resultWithPortHTTP)

	p.settings.Registry = "fakenpm.reg.org"
	resultWithoutSchemeOrPort := isNilPortOrStandardSchemePort(getParsedURI(p.settings.Registry))
	// npm requires scheme to be part of the url; so this function will return false for any missing a scheme
	assert.Equal(t, false, resultWithoutSchemeOrPort)

	p.settings.Registry = "fakenpm.reg.org:80"
	resultWithoutScheme := isNilPortOrStandardSchemePort(getParsedURI(p.settings.Registry))
	assert.Equal(t, false, resultWithoutScheme)

	p.settings.Registry = "https://fakenpm.reg.org:8443"
	resultWithNonStandardPort := isNilPortOrStandardSchemePort(getParsedURI(p.settings.Registry))
	assert.Equal(t, false, resultWithNonStandardPort)

	p.settings.Registry = "https://fakenpm.reg.org:8080"
	resultWithNonStandardPortHTTP := isNilPortOrStandardSchemePort(getParsedURI(p.settings.Registry))
	assert.Equal(t, false, resultWithNonStandardPortHTTP)
}

func TestCompareRegistries(t *testing.T) {
	p := initPlugin()
	goodReg := "https://fakenpm.reg.org/good/path"
	goodRegWithPort := "https://fakenpm.reg.org:443/good/path"
	goodRegWithNonStandardPort := "https://fakenpm.reg.org:8443/good/path"

	p.settings.Registry = goodReg
	p.settings.npm.Config.Registry = goodReg
	validNoPorts, _ := p.CompareRegistries(p.settings.npm.Config)
	assert.Equal(t, true, validNoPorts)

	p.settings.Registry = goodRegWithPort
	sameURLOneWithPort, _ := p.CompareRegistries(p.settings.npm.Config)
	assert.Equal(t, true, sameURLOneWithPort)

	p.settings.Registry = goodRegWithPort
	p.settings.npm.Config.Registry = goodRegWithPort
	sameURLBothWithPort, _ := p.CompareRegistries(p.settings.npm.Config)
	assert.Equal(t, true, sameURLBothWithPort)

	p.settings.Registry = "invalidUri"
	invalidURITest, _ := p.CompareRegistries(p.settings.npm.Config)
	assert.Equal(t, false, invalidURITest)

	p.settings.Registry = goodRegWithNonStandardPort
	nonStandardPortTest, _ := p.CompareRegistries(p.settings.npm.Config)
	assert.Equal(t, false, nonStandardPortTest)
}

func TestValidateWithInvalidFields(t *testing.T) {
	p := initPlugin()
	// Validation tests with fields missing
	p.settings.Email = ""
	noEmailErr := p.Validate()
	if assert.NotNil(t, noEmailErr) {
		assert.Contains(t, noEmailErr.Error(), "email")
	}

	p.settings.Email = "fakeemail"
	p.settings.Username = ""
	noUserErr := p.Validate()
	if assert.NotNil(t, noUserErr) {
		assert.Contains(t, noUserErr.Error(), "username")
	}

	p.settings.Username = "fakeuser"
	p.settings.Password = ""
	noPassErr := p.Validate()
	if assert.NotNil(t, noPassErr) {
		assert.Contains(t, noPassErr.Error(), "password")
	}

	p.settings.Token = "fakeToken"
	p.settings.Password = ""
	p.settings.Username = ""
	p.settings.Email = ""
	tokenErr := p.Validate()
	assert.Nil(t, tokenErr)
}

func TestValidateWithRegistryVariations(t *testing.T) {
	p := initPlugin()

	// Validation tests with an invalid registry
	p.settings.Registry = "fakenpm.reg.org/good/path"
	missingSchemeErr := p.Validate()
	if assert.NotNil(t, missingSchemeErr) {
		assert.Contains(t, missingSchemeErr.Error(), "fakenpm.reg.org")
	}

	p.settings.Registry = "https://fakenpm.reg.org:7894/good/path"
	weirdPortErr := p.Validate()
	if assert.NotNil(t, weirdPortErr) {
		assert.Contains(t, weirdPortErr.Error(), "7894")
	}

	// Validation tests with default/no ports defined
	p.settings.Registry = "https://fakenpm.reg.org:443/good/path"
	defaultPortErr := p.Validate()
	assert.Nil(t, defaultPortErr)

	// Validation tests with failure conditions on registry
	p.settings.Registry = "https://registry.npmjs.org/good/path"
	diffRegistry := p.Validate()
	if assert.NotNil(t, diffRegistry) {
		assert.Contains(t, diffRegistry.Error(), "npmjs.org")
	}

	p.settings.Registry = "https://registry.npmjs.org:443/good/path"
	diffRegistryWithPort := p.Validate()
	if assert.NotNil(t, diffRegistryWithPort) {
		assert.Contains(t, diffRegistryWithPort.Error(), "npmjs.org:443")
	}

	// Validation failures with standard ports but different paths
	p.settings.Registry = "https://registry.npmjs.org:443/bad/path"
	p.settings.Folder = "__testwithport__"
	diffRegistryWithPort = p.Validate()
	if assert.NotNil(t, diffRegistryWithPort) {
		assert.Contains(t, diffRegistryWithPort.Error(), "npmjs.org:443/bad")
	}

	// Same path different ports
	p.settings.Registry = "https://registry.npmjs.org:8443/good/path"
	p.settings.Folder = "__testwithport__"
	diffRegistryWithPort = p.Validate()
	if assert.NotNil(t, diffRegistryWithPort) {
		assert.Contains(t, diffRegistryWithPort.Error(), "npmjs.org:8443/good")
	}

	// Same path different ports and schemes
	p.settings.Registry = "http://registry.npmjs.org:80/good/path"
	p.settings.Folder = "__testwithport__"
	diffRegistryWithPort = p.Validate()
	if assert.NotNil(t, diffRegistryWithPort) {
		assert.Contains(t, diffRegistryWithPort.Error(), "npmjs.org:80/good")
	}

	// Validation tests with SkipRegistryValidation
	p.settings.SkipRegistryValidation = true
	p.settings.Registry = "fakenpm.reg.org/good/path"
	skipMissingSchemeErr := p.Validate()
	assert.Nil(t, skipMissingSchemeErr)

	p.settings.SkipRegistryValidation = true
	p.settings.Registry = "https://fakenpm.reg.org:7894/good/path"
	skipWeirdPortErr := p.Validate()
	assert.Nil(t, skipWeirdPortErr)
}

func TestExecute(t *testing.T) {
	t.Skip()
}

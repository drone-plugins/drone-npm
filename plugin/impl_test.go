// Copyright (c) 2020, the Drone Plugins project authors.
// Please see the AUTHORS file for details. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package plugin

import (
	"context"
	"net/url"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/drone-plugins/drone-plugin-lib/drone"
)

func initFakeSettings() Settings {
	nc := npmConfig{
		// Note: this registry is the one that would come from publishConfig in package.json
		Registry: "https://fakenpm.reg.org",
	}
	np := npmPackage{
		Name:    "Test Package",
		Version: "1.33.7",
		Config: nc,
	}
	return Settings{
		Username:   "fakeUser",
		Password:   "fakePass",
		Token:      "",
		SkipWhoami: false,
		Email:      "fake@user.tst",
		// Note: this registry is the one that would come from drone yaml
		Registry:                  "https://fakenpm.reg.org",
		Folder:                    "folderpath",
		FailOnVersionConflict:     true,
		Tag:                       "",
		Access:                    "",
		SkipRegistryUriValidation: false,
		npm: &np,
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

func getParsedUri(s string) *url.URL{
	rslt, _ := url.Parse(s)
	return rslt
}

func TestisDefaultOrNilPort(t *testing.T) {
	p := initPlugin()

	resultWithoutPort := isDefaultOrNilPort(getParsedUri(p.settings.Registry))
	assert.Equal(t, true, resultWithoutPort)

	p.settings.Registry = "https://fakenpm.reg.org:443"
	resultWithPort := isDefaultOrNilPort(getParsedUri(p.settings.Registry))
	assert.Equal(t, true, resultWithPort)

	p.settings.Registry = "http://fakenpm.reg.org:80"
	resultWithPortHTTP := isDefaultOrNilPort(getParsedUri(p.settings.Registry))
	assert.Equal(t, true, resultWithPortHTTP)

	p.settings.Registry = "fakenpm.reg.org"
	resultWithoutSchemeOrPort := isDefaultOrNilPort(getParsedUri(p.settings.Registry))
	assert.Equal(t, true, resultWithoutSchemeOrPort)

	p.settings.Registry = "fakenpm.reg.org:80"
	resultWithoutScheme := isDefaultOrNilPort(getParsedUri(p.settings.Registry))
	assert.Equal(t, false, resultWithoutScheme)

	p.settings.Registry = "https://fakenpm.reg.org:8443"
	resultWithNonStandardPort := isDefaultOrNilPort(getParsedUri(p.settings.Registry))
	assert.Equal(t, false, resultWithNonStandardPort)

	p.settings.Registry = "https://fakenpm.reg.org:8080"
	resultWithNonStandardPortHTTP := isDefaultOrNilPort(getParsedUri(p.settings.Registry))
	assert.Equal(t, false, resultWithNonStandardPortHTTP)
}

func TestCheckMatchingUrlWithDefaultPorts(t *testing.T) {
	t.Skip()
}

func TestValidateMissingUsername(t *testing.T) {
	t.Skip()
}
func TestValidateMissingPassword(t *testing.T) {
	t.Skip()
}
func TestValidateMissingEmail(t *testing.T) {
	t.Skip()
}
func TestValidateUsingTokenWithMissingFields(t *testing.T) {
	t.Skip()
}
func TestValidateMissingRegistry(t *testing.T) {
	t.Skip()
}
func TestValidateInvalidReg(t *testing.T) {
	t.Skip()
}
func TestValidateRegDefaultPorts(t *testing.T) {
	t.Skip()
}
func TestValidateInvalidRegWithSkipReg(t *testing.T) {
	t.Skip()
}
func TestValidate(t *testing.T) {
	t.Skip()
}

func TestExecute(t *testing.T) {
	t.Skip()
}

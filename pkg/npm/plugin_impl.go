// Copyright (c) 2019, the Drone Plugins project authors.
// Please see the AUTHORS file for details. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package npm

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
)

// Settings for the Plugin.
type Settings struct {
	Username              string
	Password              string
	Token                 string
	Email                 string
	Registry              string
	Folder                string
	FailOnVersionConflict bool
	Tag                   string
	Access                string
}

// globalRegistry defines the default NPM registry.
const globalRegistry = "https://registry.npmjs.org"

func (p *pluginImpl) Validate() error {
	// Validate the Config and return an error if there are issues.
	return nil
}

func (p *pluginImpl) Exec() error {
	// Implementation of the plugin.
	return nil
}

// npmrcContentsUsernamePassword creates the contents from a username and
// password
func npmrcContentsUsernamePassword(config Settings) string {
	// get the base64 encoded string
	authString := fmt.Sprintf("%s:%s", config.Username, config.Password)
	encoded := base64.StdEncoding.EncodeToString([]byte(authString))

	// create the file contents
	return fmt.Sprintf("_auth = %s\nemail = %s", encoded, config.Email)
}

/// Writes npmrc contents when using a token
func npmrcContentsToken(config Settings) string {
	registry, _ := url.Parse(config.Registry)
	registry.Scheme = "" // Reset the scheme to empty. This makes it so we will get a protocol relative URL.
	registryString := registry.String()

	if !strings.HasSuffix(registryString, "/") {
		registryString = registryString + "/"
	}
	return fmt.Sprintf("%s:_authToken=%s", registryString, config.Token)
}

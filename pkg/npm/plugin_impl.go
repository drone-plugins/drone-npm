// Copyright (c) 2019, the Drone Plugins project authors.
// Please see the AUTHORS file for details. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package npm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/sirupsen/logrus"
)

// Settings for the Plugin.
type (
	Settings struct {
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

	npmPackage struct {
		Name    string    `json:"name"`
		Version string    `json:"version"`
		Config  npmConfig `json:"publishConfig"`
	}

	npmConfig struct {
		Registry string `json:"registry"`
	}
)

// globalRegistry defines the default NPM registry.
const globalRegistry = "https://registry.npmjs.org"

func (p *pluginImpl) Validate() error {
	// Check authentication options
	if len(p.settings.Token) == 0 {
		if len(p.settings.Username) == 0 {
			return fmt.Errorf("No username provided")
		}
		if len(p.settings.Email) == 0 {
			return fmt.Errorf("No email address provided")
		}
		if len(p.settings.Password) == 0 {
			return fmt.Errorf("No password provided")
		}

		logrus.WithFields(logrus.Fields{
			"username": p.settings.Username,
			"email":    p.settings.Email,
		}).Info("Specified credentials")
	} else {
		logrus.Info("Token credentials being used")
	}

	// Verify package.json file
	npm, err := readPackageFile(p.settings.Folder)
	if err != nil {
		return fmt.Errorf("Invalid package.json %w", err)
	}

	// Verify the same registry is being used
	if len(p.settings.Registry) == 0 {
		p.settings.Registry = globalRegistry
	}

	if strings.Compare(p.settings.Registry, npm.Config.Registry) != 0 {
		return fmt.Errorf("Registry values do not match .drone.yml: %s package.json: %s", p.settings.Registry, npm.Config.Registry)
	}

	return nil
}

func (p *pluginImpl) Exec() error {
	// Implementation of the plugin.
	return nil
}

/// readPackageFile reads the package file at the given path.
func readPackageFile(folder string) (*npmPackage, error) {
	// Verify package.json file exists
	packagePath := path.Join(folder, "package.json")
	info, err := os.Stat(packagePath)

	if os.IsNotExist(err) {
		return nil, fmt.Errorf("No package.json at %s %w", packagePath, err)
	}
	if info.IsDir() {
		return nil, fmt.Errorf("The package.json at %s is a directory", packagePath)
	}

	// Read the file
	file, err := ioutil.ReadFile(packagePath)

	if err != nil {
		return nil, fmt.Errorf("Could not read package.json at %s %w", packagePath, err)
	}

	// Unmarshal the json data
	npm := npmPackage{}
	err = json.Unmarshal(file, &npm)

	if err != nil {
		return nil, err
	}

	// Make sure values are present
	if len(npm.Name) == 0 {
		return nil, fmt.Errorf("No package name present")
	}
	if len(npm.Version) == 0 {
		return nil, fmt.Errorf("No package version present")
	}

	// Set the default registry
	if len(npm.Config.Registry) == 0 {
		npm.Config.Registry = globalRegistry
	}

	logrus.WithFields(logrus.Fields{
		"name":    npm.Name,
		"version": npm.Version,
		"path":    packagePath,
	}).Info("Found package.json")

	return &npm, nil
}

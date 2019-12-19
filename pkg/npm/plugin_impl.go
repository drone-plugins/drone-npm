// Copyright (c) 2019, the Drone Plugins project authors.
// Please see the AUTHORS file for details. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package npm

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

func (p *pluginImpl) Validate() error {
	// Validate the Config and return an error if there are issues.
	return nil
}

func (p *pluginImpl) Exec() error {
	// Implementation of the plugin.
	return nil
}

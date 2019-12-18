// Copyright (c) 2019, the Drone Plugins project authors.
// Please see the AUTHORS file for details. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package npm

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Settings for the Plugin.
type Settings struct {
	Username              string
	Password              string
	Token                 string
	Email                 string
	Registry              string
	Folder                string
	SkipVerify            bool
	FailOnVersionConflict bool
	Tag                   string
	Access                string
}

func (p *pluginImpl) Validate() error {
	// Check authentication options
	if len(p.settings.Token) == 0 {
        if len(p.settings.Username) == 0 {
			return errors.New("No username provided")
		}
		if len(p.settings.Email) == 0 {
			return errors.New("No email address provided")
		}
		if len(p.settings.Password) == 0 {
			return errors.New("No password provided")
		}

		logrus.WithFields(logrus.Fields{
			"username": p.settings.Username,
			"email":    p.settings.Email,
		}).Info("Specified credentials")
	} else {
		logrus.Info("Token credentials being used")
	}

	return nil
}

func (p *pluginImpl) Exec() error {
	// Implementation of the plugin.
	return nil
}

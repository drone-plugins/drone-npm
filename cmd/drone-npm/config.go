// Copyright (c) 2020, the Drone Plugins project authors.
// Please see the AUTHORS file for details. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package main

import (
	"github.com/drone-plugins/drone-npm/plugin"
	"github.com/urfave/cli/v2"
)

// settingsFlags has the cli.Flags for the plugin.Settings.
func settingsFlags(settings *plugin.Settings) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "username",
			Usage:       "NPM username",
			EnvVars:     []string{"PLUGIN_USERNAME", "NPM_USERNAME"},
			Destination: &settings.Username,
		},
		&cli.StringFlag{
			Name:        "password",
			Usage:       "NPM password",
			EnvVars:     []string{"PLUGIN_PASSWORD", "NPM_PASSWORD"},
			Destination: &settings.Password,
		},
		&cli.StringFlag{
			Name:        "email",
			Usage:       "NPM email",
			EnvVars:     []string{"PLUGIN_EMAIL", "NPM_EMAIL"},
			Destination: &settings.Email,
		},
		&cli.StringFlag{
			Name:        "token",
			Usage:       "NPM deploy token",
			EnvVars:     []string{"PLUGIN_TOKEN", "NPM_TOKEN"},
			Destination: &settings.Token,
		},
		&cli.BoolFlag{
			Name:        "skip-whoami",
			Usage:       "Skip credentials verification by running npm whoami command",
			EnvVars:     []string{"PLUGIN_SKIP_WHOAMI", "NPM_SKIP_WHOAMI"},
			Destination: &settings.SkipWhoami,
		},
		&cli.StringFlag{
			Name:        "registry",
			Usage:       "NPM registry",
			EnvVars:     []string{"PLUGIN_REGISTRY", "NPM_REGISTRY"},
			Destination: &settings.Registry,
		},
		&cli.StringFlag{
			Name:        "folder",
			Usage:       "folder containing package.json",
			EnvVars:     []string{"PLUGIN_FOLDER"},
			Destination: &settings.Folder,
		},
		&cli.BoolFlag{
			Name:        "fail-on-version-conflict",
			Usage:       "fail NPM publish if version already exists in NPM registry",
			EnvVars:     []string{"PLUGIN_FAIL_ON_VERSION_CONFLICT"},
			Destination: &settings.FailOnVersionConflict,
		},
		&cli.StringFlag{
			Name:        "tag",
			Usage:       "NPM publish tag",
			EnvVars:     []string{"PLUGIN_TAG"},
			Destination: &settings.Tag,
		},
		&cli.StringFlag{
			Name:        "access",
			Usage:       "NPM scoped package access",
			EnvVars:     []string{"PLUGIN_ACCESS"},
			Destination: &settings.Access,
		},
	}
}

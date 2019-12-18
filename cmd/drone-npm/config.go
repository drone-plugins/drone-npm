// Copyright (c) 2019, the Drone Plugins project authors.
// Please see the AUTHORS file for details. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package main

import (
	"github.com/urfave/cli/v2"

	"github.com/drone-plugins/drone-npm/pkg/npm"
)

const (
	// Add all the flag names here as const strings.
)

// settingsFlags has the cli.Flags for the plugin.Settings.
func settingsFlags() []cli.Flag {
	// Replace below with all the flags required for the plugin's specific
	// settings.
	return []cli.Flag{
		&cli.StringFlag{
			Name:   "username",
			Usage:  "NPM username",
			EnvVars: []string{"PLUGIN_USERNAME","NPM_USERNAME"},
		},
		&cli.StringFlag{
			Name:   "password",
			Usage:  "NPM password",
			EnvVars: []string{"PLUGIN_PASSWORD","NPM_PASSWORD"},
		},
		&cli.StringFlag{
			Name:   "email",
			Usage:  "NPM email",
			EnvVars: []string{"PLUGIN_EMAIL","NPM_EMAIL"},
		},
		&cli.StringFlag{
			Name:   "token",
			Usage:  "NPM deploy token",
			EnvVars: []string{"PLUGIN_TOKEN","NPM_TOKEN"},
		},
		&cli.StringFlag{
			Name:   "registry",
			Usage:  "NPM registry",
			EnvVars: []string{"PLUGIN_REGISTRY","NPM_REGISTRY"},
		},
		&cli.StringFlag{
			Name:   "folder",
			Usage:  "folder containing package.json",
			EnvVars: []string{"PLUGIN_FOLDER"},
		},
		&cli.BoolFlag{
			Name:   "skip_verify",
			Usage:  "skip SSL verification",
			EnvVars: []string{"PLUGIN_SKIP_VERIFY"},
		},
		&cli.BoolFlag{
			Name:   "fail_on_version_conflict",
			Usage:  "fail NPM publish if version already exists in NPM registry",
			EnvVars: []string{"PLUGIN_FAIL_ON_VERSION_CONFLICT"},
		},
		&cli.StringFlag{
			Name:   "tag",
			Usage:  "NPM publish tag",
			EnvVars: []string{"PLUGIN_TAG"},
		},
		&cli.StringFlag{
			Name:   "access",
			Usage:  "NPM scoped package access",
			EnvVars: []string{"PLUGIN_ACCESS"},
		},
	}
}

// settingsFromContext creates a plugin.Settings from the cli.Context.
func settingsFromContext(ctx *cli.Context) npm.Settings {
	// Replace below with the parsing of the
	return npm.Settings{
		Username:              ctx.String("username"),
		Password:              ctx.String("password"),
		Token:                 ctx.String("token"),
		Email:                 ctx.String("email"),
		Registry:              ctx.String("registry"),
		Folder:                ctx.String("folder"),
		SkipVerify:            ctx.Bool("skip_verify"),
		FailOnVersionConflict: ctx.Bool("fail_on_version_conflict"),
		Tag:                   ctx.String("tag"),
		Access:                ctx.String("access"),
	}
}

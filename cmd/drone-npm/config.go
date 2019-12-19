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
	usernameFlag              = "username"
	passwordFlag              = "password"
	emailFlag                 = "email"
	tokenFlag                 = "token"
	registryFlag              = "registry"
	folderFlag                = "folder"
	failOnVersionConflictFlag = "fail-on-version-conflict"
	tagFlag                   = "tag"
	accessFlag                = "access"
)

// settingsFlags has the cli.Flags for the plugin.Settings.
func settingsFlags() []cli.Flag {
	// Replace below with all the flags required for the plugin's specific
	// settings.
	return []cli.Flag{
		&cli.StringFlag{
			Name:    usernameFlag,
			Usage:   "NPM username",
			EnvVars: []string{"PLUGIN_USERNAME", "NPM_USERNAME"},
		},
		&cli.StringFlag{
			Name:    passwordFlag,
			Usage:   "NPM password",
			EnvVars: []string{"PLUGIN_PASSWORD", "NPM_PASSWORD"},
		},
		&cli.StringFlag{
			Name:    emailFlag,
			Usage:   "NPM email",
			EnvVars: []string{"PLUGIN_EMAIL", "NPM_EMAIL"},
		},
		&cli.StringFlag{
			Name:    tokenFlag,
			Usage:   "NPM deploy token",
			EnvVars: []string{"PLUGIN_TOKEN", "NPM_TOKEN"},
		},
		&cli.StringFlag{
			Name:    registryFlag,
			Usage:   "NPM registry",
			EnvVars: []string{"PLUGIN_REGISTRY", "NPM_REGISTRY"},
		},
		&cli.StringFlag{
			Name:    folderFlag,
			Usage:   "folder containing package.json",
			EnvVars: []string{"PLUGIN_FOLDER"},
		},
		&cli.BoolFlag{
			Name:    failOnVersionConflictFlag,
			Usage:   "fail NPM publish if version already exists in NPM registry",
			EnvVars: []string{"PLUGIN_FAIL_ON_VERSION_CONFLICT"},
		},
		&cli.StringFlag{
			Name:    tagFlag,
			Usage:   "NPM publish tag",
			EnvVars: []string{"PLUGIN_TAG"},
		},
		&cli.StringFlag{
			Name:    accessFlag,
			Usage:   "NPM scoped package access",
			EnvVars: []string{"PLUGIN_ACCESS"},
		},
	}
}

// settingsFromContext creates a plugin.Settings from the cli.Context.
func settingsFromContext(ctx *cli.Context) npm.Settings {
	// Replace below with the parsing of the
	return npm.Settings{
		Username:              ctx.String(usernameFlag),
		Password:              ctx.String(passwordFlag),
		Token:                 ctx.String(tokenFlag),
		Email:                 ctx.String(emailFlag),
		Registry:              ctx.String(registryFlag),
		Folder:                ctx.String(folderFlag),
		FailOnVersionConflict: ctx.Bool(failOnVersionConflictFlag),
		Tag:                   ctx.String(tagFlag),
		Access:                ctx.String(accessFlag),
	}
}

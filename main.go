// Copyright (c) 2020, the Drone Plugins project authors.
// Please see the AUTHORS file for details. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

// DO NOT MODIFY THIS FILE DIRECTLY

package main

import (
	"os"

	"github.com/drone-plugins/drone-npm/plugin"
	"github.com/drone-plugins/drone-plugin-lib/errors"
	"github.com/drone-plugins/drone-plugin-lib/urfave"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

var version = "unknown"

func main() {
	settings := &plugin.Settings{}

	if _, err := os.Stat("/run/drone/env"); err == nil {
		godotenv.Overload("/run/drone/env") //nolint:errcheck
	}

	app := &cli.App{
		Name:    "drone-npm",
		Usage:   "push a package to a npm repository",
		Version: version,
		Flags:   append(settingsFlags(settings), urfave.Flags()...),
		Action:  run(settings),
	}

	if err := app.Run(os.Args); err != nil {
		errors.HandleExit(err)
	}
}

func run(settings *plugin.Settings) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		urfave.LoggingFromContext(ctx)

		p := plugin.New(
			*settings,
			urfave.PipelineFromContext(ctx),
			urfave.NetworkFromContext(ctx),
		)

		if err := p.Validate(); err != nil {
			if e, ok := err.(errors.ExitCoder); ok {
				return e
			}

			return errors.ExitMessagef("validation failed: %w", err)
		}

		if err := p.Execute(); err != nil {
			if e, ok := err.(errors.ExitCoder); ok {
				return e
			}

			return errors.ExitMessagef("execution failed: %w", err)
		}

		return nil
	}
}

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
		&cli.BoolFlag{
			Name:        "skip-registry-validation",
			Usage:       "skips validation for uri in package.json and the currently configured registry",
			EnvVars:     []string{"PLUGIN_SKIP_REGISTRY_VALIDATION"},
			Destination: &settings.SkipRegistryValidation,
		},
	}
}

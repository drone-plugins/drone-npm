// Copyright (c) 2019, the Drone Plugins project authors.
// Please see the AUTHORS file for details. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

// DO NOT MODIFY THIS FILE DIRECTLY

package main

import (
	"fmt"
	"os"

	"github.com/drone-plugins/drone-plugin-lib/pkg/urfave"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/drone-plugins/drone-npm/pkg/npm"
)

var (
	version = "unknown"
)

func main() {
	app := cli.NewApp()
	app.Name = "npm plugin"
	app.Usage = "pushes a package to a npm repository"
	app.Action = run
	app.Flags = append(settingsFlags(), urfave.Flags()...)

	// Run the application
	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func run(ctx *cli.Context) error {
	urfave.LoggingFromContext(ctx)

	plugin := npm.New(
		settingsFromContext(ctx),
		urfave.PipelineFromContext(ctx),
		urfave.NetworkFromContext(ctx),
	)

	// Validate the settings
	if err := plugin.Validate(); err != nil {
		return fmt.Errorf("Validation failed %w", err)
	}

	// Run the plugin
	if err := plugin.Exec(); err != nil {
		return fmt.Errorf("Execution failed %w", err)
	}

	return nil
}

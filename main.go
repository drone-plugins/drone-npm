package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	version = "unknown"
)

func main() {
	app := cli.NewApp()
	app.Name = "npm plugin"
	app.Usage = "npm plugin"
	app.Action = run
	app.Version = version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "username",
			Usage:  "NPM username",
			EnvVar: "PLUGIN_USERNAME,NPM_USERNAME",
		},
		cli.StringFlag{
			Name:   "password",
			Usage:  "NPM password",
			EnvVar: "PLUGIN_PASSWORD,NPM_PASSWORD",
		},
		cli.StringFlag{
			Name:   "email",
			Usage:  "NPM email",
			EnvVar: "PLUGIN_EMAIL,NPM_EMAIL",
		},
		cli.StringFlag{
			Name:   "token",
			Usage:  "NPM deploy token",
			EnvVar: "PLUGIN_TOKEN,NPM_TOKEN",
		},
		cli.StringFlag{
			Name:   "registry",
			Usage:  "NPM registry",
			Value:  GlobalRegistry,
			EnvVar: "PLUGIN_REGISTRY,NPM_REGISTRY",
		},
		cli.StringFlag{
			Name:   "folder",
			Usage:  "folder containing package.json",
			EnvVar: "PLUGIN_FOLDER",
		},
		cli.BoolFlag{
			Name:   "skip_verify",
			Usage:  "skip SSL verification",
			EnvVar: "PLUGIN_SKIP_VERIFY",
		},
		cli.BoolFlag{
			Name:   "fail_on_version_conflict",
			Usage:  "fail NPM publish if version already exists in NPM registry",
			EnvVar: "PLUGIN_FAIL_ON_VERSION_CONFLICT",
		},
		cli.StringFlag{
			Name:  "env-file",
			Usage: "source env file",
		},
		cli.StringFlag{
			Name:   "tag",
			Usage:  "NPM publish tag",
			EnvVar: "PLUGIN_TAG",
		},
		cli.StringFlag{
			Name:   "access",
			Usage:  "NPM scoped package access",
			EnvVar: "PLUGIN_ACCESS",
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func run(c *cli.Context) error {
	if c.String("env-file") != "" {
		_ = godotenv.Load(c.String("env-file"))
	}

	plugin := Plugin{
		Config: Config{
			Username:              c.String("username"),
			Password:              c.String("password"),
			Token:                 c.String("token"),
			Email:                 c.String("email"),
			Registry:              c.String("registry"),
			Folder:                c.String("folder"),
			SkipVerify:            c.Bool("skip_verify"),
			FailOnVersionConflict: c.Bool("fail_on_version_conflict"),
			Tag:                   c.String("tag"),
			Access:                c.String("access"),
		},
	}

	return plugin.Exec()
}

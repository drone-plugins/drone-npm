package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
	_ "github.com/joho/godotenv/autoload"
)

var version string // build number set at compile-time

func main() {
	app := cli.NewApp()
	app.Name = "npm"
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
	}

	app.Run(os.Args)
}

func run(c *cli.Context) error {
	plugin := Plugin{
		Config: Config{
			Username:   c.String("username"),
			Password:   c.String("password"),
			Token:      c.String("token"),
			Email:      c.String("email"),
			Registry:   c.String("registry"),
			Folder:     c.String("folder"),
			SkipVerify: c.Bool("skip_verify"),
		},
	}

	if err := plugin.Exec(); err != nil {
		log.Error(err)
		os.Exit(1)
	}

	return nil
}

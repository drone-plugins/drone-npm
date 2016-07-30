package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"os/user"
	"path"
	"strings"

	log "github.com/Sirupsen/logrus"
)

type (
	// Config for the plugin.
	Config struct {
		Username   string
		Password   string
		Token      string
		Email      string
		Registry   string
		Folder     string
		SkipVerify bool
	}

	npmPackage struct {
		Name    string    `json:"name"`
		Version string    `json:"version"`
		Config  npmConfig `json:"publishConfig"`
	}

	npmConfig struct {
		Registry string `json:"registry"`
	}

	// Plugin values
	Plugin struct {
		Config Config
	}
)

// GlobalRegistry defines the default NPM registry.
const GlobalRegistry = "https://registry.npmjs.org"

// Exec executes the plugin.
func (p Plugin) Exec() error {
	// write npmrc for authentication
	err := writeNpmrc(p.Config)

	if err != nil {
		return err
	}

	// attempt to authenticate
	err = authenticate(p.Config)

	if err != nil {
		return err
	}

	// read the package
	npmPackage, err := readPackageFile(p.Config)

	if err != nil {
		return err
	}

	// determine whether to publish
	publish, err := shouldPublishPackage(p.Config, npmPackage)

	if err != nil {
		return err
	}

	if publish {
		log.Info("Publishing package")

		// run the publish command
		return runCommand(publishCommand(), p.Config.Folder)
	}

	return nil
}

/// writeNpmrc creates a .npmrc in the folder for authentication
func writeNpmrc(config Config) error {
	var npmrcContents string

	// check for an auth token
	if len(config.Token) == 0 {
		// check for a username
		if len(config.Username) == 0 {
			return errors.New("No username provided")
		}

		// check for an email
		if len(config.Email) == 0 {
			return errors.New("No email address provided")
		}

		// check for a password
		if len(config.Password) == 0 {
			log.Warning("No password provided")
		}

		log.WithFields(log.Fields{
			"username": config.Username,
			"email":    config.Email,
		}).Info("Specified credentials")

		npmrcContents = npmrcContentsUsernamePassword(config)
	} else {
		log.Info("Token credentials being used")

		npmrcContents = npmrcContentsToken(config)
	}

	// write npmrc file
	home := "/root"
	user, err := user.Current()
	if err == nil {
		home = user.HomeDir
	}
	npmrcPath := path.Join(home, ".npmrc")

	log.WithFields(log.Fields{
		"path": npmrcPath,
	}).Info("Writing npmrc")

	return ioutil.WriteFile(npmrcPath, []byte(npmrcContents), 0644)
}

/// authenticate atempts to authenticate with the NPM registry.
func authenticate(config Config) error {
	var cmds []*exec.Cmd

	// write the version command
	cmds = append(cmds, versionCommand())

	// write registry command
	if config.Registry != GlobalRegistry {
		cmds = append(cmds, registryCommand(config.Registry))
	}

	// write auth command
	cmds = append(cmds, alwaysAuthCommand())

	// write skip verify command
	if config.SkipVerify {
		cmds = append(cmds, skipVerifyCommand())
	}

	// write whoami command to verify credentials
	cmds = append(cmds, whoamiCommand())

	// run commands
	err := runCommands(cmds, config.Folder)

	if err != nil {
		return errors.New("Could not authenticate")
	}

	return nil
}

/// readPackageFile reads the package file at the given path.
func readPackageFile(config Config) (*npmPackage, error) {
	// read the file
	packagePath := path.Join(config.Folder, "package.json")
	file, err := ioutil.ReadFile(packagePath)

	if err != nil {
		return nil, err
	}

	// unmarshal the json data
	npm := npmPackage{}
	err = json.Unmarshal(file, &npm)

	if len(npm.Config.Registry) == 0 {
		log.Info("No registry specified in the package.json")
		npm.Config.Registry = GlobalRegistry
	}

	if err != nil {
		return nil, err
	}

	// make sure values are present
	if len(npm.Name) == 0 {
		return nil, errors.New("No package name present")
	}

	if len(npm.Version) == 0 {
		return nil, errors.New("No package version present")
	}

	// check package registry
	if strings.Compare(config.Registry, npm.Config.Registry) != 0 {
		return nil, fmt.Errorf("Registry values do not match .drone.yml: %s package.json: %s", config.Registry, npm.Config.Registry)
	}

	log.WithFields(log.Fields{
		"name":    npm.Name,
		"version": npm.Version,
	}).Info("Found package.json")

	return &npm, nil
}

/// shouldPublishPackage determines if the package should be published
func shouldPublishPackage(config Config, npm *npmPackage) (bool, error) {
	cmd := packageVersionsCommand(npm.Name)
	cmd.Dir = config.Folder

	trace(cmd)
	out, err := cmd.CombinedOutput()

	// see if there was an error
	// if there is an error its likely due to the package never being published
	if err == nil {
		// parse the json output
		var versions []string
		err = json.Unmarshal(out, &versions)

		if err != nil {
			log.Debug("Could not parse into array of string. Likely single value")

			var version string
			err := json.Unmarshal(out, &version)

			if err != nil {
				return false, err
			}

			versions = append(versions, version)
		}

		for _, value := range versions {
			log.WithFields(log.Fields{
				"version": value,
			}).Debug("Found version of package")

			if strings.Compare(npm.Version, value) == 0 {
				return false, nil
			}
		}

		log.Info("Version not found in the registry")
	} else {
		log.Info("Name was not found in the registry")
	}

	return true, nil
}

// npmrcContentsUsernamePassword creates the contents from a username and
// password
func npmrcContentsUsernamePassword(config Config) string {
	// get the base64 encoded string
	authString := fmt.Sprintf("%s:%s", config.Username, config.Password)
	encoded := base64.StdEncoding.EncodeToString([]byte(authString))

	// create the file contents
	return fmt.Sprintf("_auth = %s\nemail = %s", encoded, config.Email)
}

/// Writes npmrc contents when using a token
func npmrcContentsToken(config Config) string {
	registry, _ := url.Parse(config.Registry)
	return fmt.Sprintf("//%s/:_authToken=%s", registry.Host, config.Token)
}

// versionCommand gets the npm version
func versionCommand() *exec.Cmd {
	return exec.Command("npm", "--version")
}

// registryCommand sets the NPM registry.
func registryCommand(registry string) *exec.Cmd {
	return exec.Command("npm", "config", "set", "registry", registry)
}

// alwaysAuthCommand forces authentication.
func alwaysAuthCommand() *exec.Cmd {
	return exec.Command("npm", "config", "set", "always-auth", "true")
}

// skipVerifyCommand disables ssl verification.
func skipVerifyCommand() *exec.Cmd {
	return exec.Command("npm", "config", "set", "strict-ssl", "false")
}

// whoamiCommand creates a command that gets the currently logged in user.
func whoamiCommand() *exec.Cmd {
	return exec.Command("npm", "whoami")
}

// packageVersionsCommand gets the versions of the npm package.
func packageVersionsCommand(name string) *exec.Cmd {
	return exec.Command("npm", "view", name, "versions", "--json")
}

// publishCommand runs the publish command
func publishCommand() *exec.Cmd {
	return exec.Command("npm", "publish")
}

// trace writes each command to standard error (preceded by a ‘$ ’) before it
// is executed. Used for debugging your build.
func trace(cmd *exec.Cmd) {
	fmt.Fprintf(os.Stdout, "+ %s\n", strings.Join(cmd.Args, " "))
}

// runCommands executes the list of cmds in the given directory.
func runCommands(cmds []*exec.Cmd, dir string) error {
	for _, cmd := range cmds {
		err := runCommand(cmd, dir)

		if err != nil {
			return err
		}
	}

	return nil
}

func runCommand(cmd *exec.Cmd, dir string) error {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = dir
	trace(cmd)

	return cmd.Run()
}

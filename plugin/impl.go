// Copyright (c) 2020, the Drone Plugins project authors.
// Please see the AUTHORS file for details. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package plugin

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"os"
	"os/exec"
	"os/user"
	"path"
	"strings"

	"github.com/sirupsen/logrus"
)

type (
	// Settings for the Plugin.
	Settings struct {
		Username              string
		Password              string
		Token                 string
		SkipWhoami            bool
		Email                 string
		Registry              string
		Folder                string
		FailOnVersionConflict bool
		Tag                   string
		Access                string
		SkipRegistryUriValidation bool

		npm *npmPackage
	}

	npmPackage struct {
		Name    string    `json:"name"`
		Version string    `json:"version"`
		Config  npmConfig `json:"publishConfig"`
	}

	npmConfig struct {
		Registry string `json:"registry"`
	}
)

// globalRegistry defines the default NPM registry.
const globalRegistry = "https://registry.npmjs.org/"
const defaultPortMap = map[string]string{
	"http":80,
	"https":443,
}


func (p *Plugin) CheckMatchingUrlWithDefaultPorts() bool, error{
	parsedConifgReg, err:=url.Parse(npm.Config.Registry)
	if err != nil{
		return false,fmt.Errorf("package.json registry: %s failed to parse.")
	}
	parsedSettingsReg:=url.Parse(p.settings.Registry)
	if err != nil{
		return false, fmt.Errorf("Drone yaml npm Registry: %s failed to parse.")
	}
	compareWithoutDefaultPorts := strings.Compare(parsedConifgReg.Hostname(),parsedSettingsReg.Hostname()) &&
								  strings.Compare(parsedConifgReg.Scheme, parsedSettingsReg.Scheme) &&
								  parsedConifgReg.isDefaultOrNilPort() &&
								  parsedConifgReg.isDefaultOrNilPort()
	return compareWithoutDefaultPorts, nil
								  
}

func (u *URL) isDefaultOrNilPort() bool{
	if u.Port() != nil{
		if port, ok:=defaultPortMap[u.Scheme]; ok{
			return port == u.Port()
		}
		return false // this only happens if the scheme isn't in the above map. In this case the standard validation logic would apply
	}
	return true
}

// Validate handles the settings validation of the plugin.
func (p *Plugin) Validate() error {
	// Check authentication options
	if p.settings.Token == "" {
		if p.settings.Username == "" {
			return fmt.Errorf("no username provided")
		}
		if p.settings.Email == "" {
			return fmt.Errorf("no email address provided")
		}
		if p.settings.Password == "" {
			return fmt.Errorf("no password provided")
		}

		logrus.WithFields(logrus.Fields{
			"username": p.settings.Username,
			"email":    p.settings.Email,
		}).Info("Specified credentials")
	} else {
		logrus.Info("Token credentials being used")
	}

	// Verify package.json file
	npm, err := readPackageFile(p.settings.Folder)
	if err != nil {
		return fmt.Errorf("invalid package.json: %w", err)
	}

	// Verify the same registry is being used
	if p.settings.Registry == "" {
		p.settings.Registry = globalRegistry
	}	

	if p.settings.SkipRegistryUriValidation {
		p.settings.npm = npm
		return nil
	} else if strings.Compare(p.settings.Registry, npm.Config.Registry) != 0 || !p.CheckMatchingUrlWithDefaultPorts() {
		return fmt.Errorf("registry values do not match .drone.yml: %s package.json: %s", p.settings.Registry, npm.Config.Registry)
	}

	p.settings.npm = npm

	return nil
}

// Execute provides the implementation of the plugin.
func (p *Plugin) Execute() error {
	// Write the npmrc file
	if err := p.writeNpmrc(); err != nil {
		return fmt.Errorf("could not create npmrc: %w", err)
	}

	// Attempt authentication
	if err := p.authenticate(); err != nil {
		return fmt.Errorf("could not authenticate: %w", err)
	}

	// Determine whether to publish
	publish, err := p.shouldPublishPackage()

	if err != nil {
		return fmt.Errorf("could not determine if package should be published: %w", err)
	}

	if publish {
		logrus.Info("Publishing package")
		if err = runCommand(publishCommand(&p.settings), p.settings.Folder); err != nil {
			return fmt.Errorf("could not publish package: %w", err)
		}
	} else {
		logrus.Info("Not publishing package")
	}

	return nil
}

// / writeNpmrc creates a .npmrc in the folder for authentication
func (p *Plugin) writeNpmrc() error {
	var f func(settings *Settings) string
	if p.settings.Token == "" {
		logrus.WithFields(logrus.Fields{
			"username": p.settings.Username,
			"email":    p.settings.Email,
		}).Info("Specified credentials")
		f = npmrcContentsUsernamePassword
	} else {
		logrus.Info("Token credentials being used")
		f = npmrcContentsToken
	}

	// write npmrc file
	home := "/root"
	currentUser, err := user.Current()
	if err == nil {
		home = currentUser.HomeDir
	}
	npmrcPath := path.Join(home, ".npmrc")

	logrus.WithField("path", npmrcPath).Info("Writing npmrc")

	return os.WriteFile(npmrcPath, []byte(f(&p.settings)), 0644) //nolint:gomnd
}

// / shouldPublishPackage determines if the package should be published
func (p *Plugin) shouldPublishPackage() (bool, error) {
	cmd := packageVersionsCommand(p.settings.npm.Name)
	cmd.Dir = p.settings.Folder

	trace(cmd)
	out, err := cmd.CombinedOutput()

	// see if there was an error
	// if there is an error its likely due to the package never being published
	if err == nil {
		// parse the json output
		var versions []string
		err = json.Unmarshal(out, &versions)

		if err != nil {
			logrus.Debug("Could not parse into array of string. Likely single value")

			var version string
			err := json.Unmarshal(out, &version)

			if err != nil {
				return false, err
			}

			versions = append(versions, version)
		}

		for _, value := range versions {
			logrus.WithField("version", value).Debug("Found version of package")

			if p.settings.npm.Version == value {
				logrus.Info("Version found in the registry")
				if p.settings.FailOnVersionConflict {
					return false, fmt.Errorf("cannot publish package due to version conflict")
				}
				return false, nil
			}
		}

		logrus.Info("Version not found in the registry")
	} else {
		logrus.Info("Name was not found in the registry")
	}

	return true, nil
}

// / authenticate atempts to authenticate with the NPM registry.
func (p *Plugin) authenticate() error {
	var cmds []*exec.Cmd

	// Write the version command
	cmds = append(cmds, versionCommand())

	// write registry command
	if p.settings.Registry != globalRegistry {
		cmds = append(cmds, registryCommand(p.settings.Registry))
	}

	// Write skip verify command
	if p.network.SkipVerify {
		cmds = append(cmds, skipVerifyCommand())
	}

	// Write whoami command to verify credentials
	if !p.settings.SkipWhoami {
		cmds = append(cmds, whoamiCommand())
	}

	// Run commands
	err := runCommands(cmds, p.settings.Folder)

	if err != nil {
		return err
	}

	return nil
}

// / readPackageFile reads the package file at the given path.
func readPackageFile(folder string) (*npmPackage, error) {
	// Verify package.json file exists
	packagePath := path.Join(folder, "package.json")
	info, err := os.Stat(packagePath)

	if os.IsNotExist(err) {
		return nil, fmt.Errorf("no package.json at %s: %w", packagePath, err)
	}
	if info.IsDir() {
		return nil, fmt.Errorf("the package.json at %s is a directory", packagePath)
	}

	// Read the file
	file, err := os.ReadFile(packagePath)
	if err != nil {
		return nil, fmt.Errorf("could not read package.json at %s: %w", packagePath, err)
	}

	// Unmarshal the json data
	npm := npmPackage{}
	err = json.Unmarshal(file, &npm)
	if err != nil {
		return nil, err
	}

	// Make sure values are present
	if npm.Name == "" {
		return nil, fmt.Errorf("no package name present")
	}
	if npm.Version == "" {
		return nil, fmt.Errorf("no package version present")
	}

	// Set the default registry
	if npm.Config.Registry == "" {
		npm.Config.Registry = globalRegistry
	}

	logrus.WithFields(logrus.Fields{
		"name":    npm.Name,
		"version": npm.Version,
		"path":    packagePath,
	}).Info("Found package.json")

	return &npm, nil
}

// npmrcContentsUsernamePassword creates the contents from a username and
// password
func npmrcContentsUsernamePassword(config *Settings) string {
	// get the base64 encoded string
	authString := fmt.Sprintf("%s:%s", config.Username, config.Password)
	encoded := base64.StdEncoding.EncodeToString([]byte(authString))

	// create the file contents
	return fmt.Sprintf("_auth = %s\nemail = %s", encoded, config.Email)
}

// / Writes npmrc contents when using a token
func npmrcContentsToken(config *Settings) string {
	registry, _ := url.Parse(config.Registry)
	registry.Scheme = "" // Reset the scheme to empty. This makes it so we will get a protocol relative URL.
	host, port, _ := net.SplitHostPort(registry.Host)
	if port == "80" || port == "443" {
		registry.Host = host // Remove standard ports as they're not supported in authToken since NPM 7.
	}
	registryString := registry.String()

	if !strings.HasSuffix(registryString, "/") {
		registryString += "/"
	}
	return fmt.Sprintf("%s:_authToken=%s", registryString, config.Token)
}

// versionCommand gets the npm version
func versionCommand() *exec.Cmd {
	return exec.Command("npm", "--version")
}

// registryCommand sets the NPM registry.
func registryCommand(registry string) *exec.Cmd {
	return exec.Command("npm", "config", "set", "registry", registry)
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
func publishCommand(settings *Settings) *exec.Cmd {
	commandArgs := []string{"publish"}

	if settings.Tag != "" {
		commandArgs = append(commandArgs, "--tag", settings.Tag)
	}

	if settings.Access != "" {
		commandArgs = append(commandArgs, "--access", settings.Access)
	}

	return exec.Command("npm", commandArgs...)
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

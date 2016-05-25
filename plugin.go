package main

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strings"

	log "github.com/Sirupsen/logrus"
)

type (
	Config struct {
		Username   string
		Password   string
		Email      string
		Registry   string
		Folder     string
		AlwaysAuth bool
    SkipVerify bool
	}

	NpmPackage struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}

	Plugin struct {
		Config Config
	}
)

const GlobalRegistry = "https://registry.npmjs.org"

func (p Plugin) Exec() error {
	// check for a username
	if len(p.Config.Username) == 0 {
		log.Error("No username provided")
		return errors.New("No username provided")
	}

	// check for an email
	if len(p.Config.Email) == 0 {
		log.Error("No email address provided")
		return errors.New("No email address provided")
	}

	// check for a password
	if len(p.Config.Password) == 0 {
		log.Warning("No password provided")
	}

	log.WithFields(log.Fields{
		"username": p.Config.Username,
		"email": p.Config.Email,
	}).Info("Specified credentials")

	// read the package
	packagePath := path.Join(p.Config.Folder, "package.json")

	npmPackage, err := readPackageFile(packagePath)

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Could not read package.json")
		return err
	}

	log.WithFields(log.Fields{
		"name":    npmPackage.Name,
		"version": npmPackage.Version,
	}).Info("Found package")

	// see if the package should be published
	publish, err := shouldPublishPackage(p.Config, npmPackage)

	if publish {
		log.Info("Attempting to publish package")

		// write the npmrc file
		err := writeNpmrcFile(p.Config)

		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Problem creating npmrc file")

			return err
		}

		var cmds []*exec.Cmd

    // write the version command
    cmds = append(cmds, versionCommand())

		// write registry command
		if p.Config.Registry != GlobalRegistry {
			cmds = append(cmds, registryCommand(p.Config.Registry))
		}

		// write auth command
		if p.Config.AlwaysAuth {
			cmds = append(cmds, alwaysAuthCommand())
		}

    // write skip verify command
    if p.Config.SkipVerify {
      cmds = append(cmds, skipVerifyCommand())
    }

		// write the publish command
		cmds = append(cmds, publishCommand())

		// run the commands
		err = runCommands(cmds)

		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Package was not published")
		}
	} else {
		log.WithFields(log.Fields{
			"reason": err,
		}).Info("Package will not be published")
	}

	return nil
}

/// Reads the package file at the given path
func readPackageFile(path string) (*NpmPackage, error) {
	// read the file
	file, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	// unmarshal the json data
	npmPackage := NpmPackage{}
	err = json.Unmarshal(file, &npmPackage)

	if err != nil {
		return nil, err
	}

	// make sure values are present
	if len(npmPackage.Name) == 0 {
		return nil, errors.New("No package name present")
	}

	if len(npmPackage.Version) == 0 {
		return nil, errors.New("No package version present")
	}

	return &npmPackage, nil
}

/// Determines if the package should be published
func shouldPublishPackage(config Config, npmPackage *NpmPackage) (bool, error) {
	// get the url for the package
	//
	// encoding the portion of the url for the package name in case of scopes
	packageUrl := fmt.Sprintf("%s/%s/%s", config.Registry, url.QueryEscape(npmPackage.Name), npmPackage.Version)

	log.WithFields(log.Fields{
		"url": packageUrl,
	}).Info("Requesting package information")

	// create a request for the package
	req, err := http.NewRequest("GET", packageUrl, nil)

	if err != nil {
		return false, err
	}

	// set authentication if needed
	// NPM uses basic http auth
	if config.AlwaysAuth {
		req.SetBasicAuth(config.Username, config.Password)
	}

  // skip verify if necessary
  if config.SkipVerify {
    http.DefaultTransport = &http.Transport{
      TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }

    log.Warning("Skipping SSL verification")
  }

	// get the response
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return false, err
	}

	defer resp.Body.Close()

	statusCode := resp.StatusCode

	// look for a 404 to see if the package should be published
	if statusCode == http.StatusNotFound {
		return true, nil
	} else if statusCode == http.StatusOK {
		return false, errors.New("Package already published")
	} else {
		contents, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return false, err
		}

		return false, errors.New(fmt.Sprintf("Error Occurred. Status %d\nBody:%s", statusCode, contents))
	}
}

/// Writes the npmrc file
func writeNpmrcFile(config Config) error {
	// get the base64 encoded string
	authString := fmt.Sprintf("%s:%s", config.Username, config.Password)
	encoded := base64.StdEncoding.EncodeToString([]byte(authString))

	// create the file contents
	contents := fmt.Sprintf("_auth = %s\nemail = %s", encoded, config.Email)

	// write the file
	return ioutil.WriteFile("/root/.npmrc", []byte(contents), 0644)
}

// Gets the npm version
func versionCommand() *exec.Cmd {
  return exec.Command("npm", "--version")
}

// Sets the npm registry
func registryCommand(registry string) *exec.Cmd {
	return exec.Command("npm", "config", "set", "registry", registry)
}

// Sets the always off flag
func alwaysAuthCommand() *exec.Cmd {
	return exec.Command("npm", "config", "set", "always-auth", "true")
}

// Skip ssl verification
func skipVerifyCommand() *exec.Cmd {
  return exec.Command("npm", "config", "set", "ca=\"\"")
}

// Publishes the package
func publishCommand() *exec.Cmd {
	return exec.Command("npm", "publish")
}

// Trace writes each command to standard error (preceded by a ‘$ ’) before it
// is executed. Used for debugging your build.
func trace(cmd *exec.Cmd) {
	fmt.Fprintf(os.Stdout, "+ %s\n", strings.Join(cmd.Args, " "))
}

func runCommands(cmds []*exec.Cmd) error {
	for _, cmd := range cmds {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		trace(cmd)

		err := cmd.Run()
		if err != nil {
			return err
		}
	}

	return nil
}

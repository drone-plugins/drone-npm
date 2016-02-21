package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin"
)

type Npm struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	Email      string `json:"email"`
	Registry   string `json:"registry"`
	Folder     string `json:"folder"`
	AlwaysAuth bool   `json:"always_auth"`
}

type NpmPackage struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

var (
	buildCommit string
)

func main() {
	fmt.Printf("Drone NPM Plugin built from %s\n", buildCommit)

	repo := drone.Repo{}
	build := drone.Build{}
	workspace := drone.Workspace{}
	vargs := Npm{}

	plugin.Param("build", &build)
	plugin.Param("repo", &repo)
	plugin.Param("workspace", &workspace)
	plugin.Param("vargs", &vargs)

	// parse the parameters
	if err := plugin.Parse(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// check for required parameters
	if len(vargs.Username) == 0 {
		fmt.Println("Username not provided")
		os.Exit(1)
	}

	if len(vargs.Password) == 0 {
		fmt.Println("Password not provided")
		os.Exit(1)
	}

	if len(vargs.Email) == 0 {
		fmt.Println("Email not provided")
		os.Exit(1)
	}

	// set defaults
	var globalRegistry bool

	if len(vargs.Registry) == 0 {
		vargs.Registry = "https://registry.npmjs.org"
		globalRegistry = true
	} else {
		globalRegistry = false
	}

	// get the package info
	var packagePath string

	if len(vargs.Folder) == 0 {
		packagePath = path.Join(workspace.Path)
	} else {
		packagePath = path.Join(workspace.Path, vargs.Folder)
	}

	packageFile := path.Join(packagePath, "package.json")

	npmPackage, err := readPackageFile(packageFile)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// see if the package should be published
	publish, err := shouldPublishPackage(vargs, npmPackage)

	if publish {
		fmt.Println("Attempting to publish package")

		// write the npmrc file
		err := writeNpmrcFile(vargs)

		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		var cmds []*exec.Cmd

		// write registry command
		if !globalRegistry {
			cmds = append(cmds, registryCommand(vargs))
		}

		// write auth command
		if vargs.AlwaysAuth {
			cmds = append(cmds, alwaysAuthCommand())
		}

		// write the publish command
		cmds = append(cmds, publishCommand())

		// run the commands
		for _, cmd := range cmds {
			cmd.Dir = packagePath
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			trace(cmd)
			err := cmd.Run()
			if err != nil {
				os.Exit(1)
			}
		}
	} else {
		fmt.Println("Package already published")
	}
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
func shouldPublishPackage(vargs Npm, npmPackage *NpmPackage) (bool, error) {
	// get the url for the package
	packageUrl := fmt.Sprintf("%s/%s/%s", vargs.Registry, npmPackage.Name, npmPackage.Version)
	fmt.Printf("Requesting %s\n", packageUrl)

	// create a request for the package
	req, err := http.NewRequest("GET", packageUrl, nil)

	if err != nil {
		return false, err
	}

	// set authentication if needed
	// NPM uses basic http auth
	if vargs.AlwaysAuth {
		req.SetBasicAuth(vargs.Username, vargs.Password)
	}

	// get the response
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return false, err
	}

	defer resp.Body.Close()

	// look for a 404 to see if the package should be published
	if resp.StatusCode != http.StatusNotFound {
		return false, nil
	} else {
		return true, nil
	}
}

/// Writes the npmrc file
func writeNpmrcFile(vargs Npm) error {

	// get the base64 encoded string
	authString := fmt.Sprintf("%s:%s", vargs.Username, vargs.Password)
	encoded := base64.StdEncoding.EncodeToString([]byte(authString))

	// create the file contents
	contents := fmt.Sprintf("_auth = %s\nemail = %s", encoded, vargs.Email)

	// write the file
	return ioutil.WriteFile("/root/.npmrc", []byte(contents), 0644)
}

// Sets the npm registry
func registryCommand(vargs Npm) *exec.Cmd {
	return exec.Command("npm", "config", "set", "registry", vargs.Registry)
}

// Sets the always off flag
func alwaysAuthCommand() *exec.Cmd {
	return exec.Command("npm", "config", "set", "always-auth", "true")
}

// Publishes the package
func publishCommand() *exec.Cmd {
	return exec.Command("npm", "publish")
}

// Trace writes each command to standard error (preceded by a ‘$ ’) before it
// is executed. Used for debugging your build.
func trace(cmd *exec.Cmd) {
	fmt.Println("$", strings.Join(cmd.Args, " "))
}

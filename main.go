package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/drone/drone-plugin-go/plugin"
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

func main() {
	repo := plugin.Repo{}
	build := plugin.Build{}
	vargs := Npm{}

	plugin.Param("build", &build)
	plugin.Param("repo", &repo)
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
	if len(vargs.Registry) == 0 {
		vargs.Registry = "https://registry.npmjs.org"
	}

	// get the package info
	npmPackage, err := readPackageFile("package.json")

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// see if the package should be published
	publish, err := shouldPublishPackage(vargs, npmPackage)

	if publish {
		fmt.Println("Attempting to publish package")
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

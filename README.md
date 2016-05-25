# drone-npm

[![Build Status](http://beta.drone.io/api/badges/drone-plugins/drone-npm/status.svg)](http://beta.drone.io/drone-plugins/drone-npm)
[![Coverage Status](https://aircover.co/badges/drone-plugins/drone-npm/coverage.svg)](https://aircover.co/drone-plugins/drone-npm)
[![](https://badge.imagelayers.io/plugins/drone-npm:latest.svg)](https://imagelayers.io/?images=plugins/drone-npm:latest 'Get your own badge on imagelayers.io')

Drone plugin to publish files and artifacts to a NPM registry. For the usage information and a listing of the available options please take a look at [the docs](DOCS.md).

## Build

Build the binary with the following commands:

```
export GO15VENDOREXPERIMENT=1
go build
go test
```

## Docker

Build the docker image with the following commands:

```
export GO15VENDOREXPERIMENT=1
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -tags netgo
```

Please note incorrectly building the image for the correct x64 linux and with GCO disabled will result in an error when running the Docker image:

```
docker: Error response from daemon: Container command
'/bin/drone-git' not found or does not exist..
```

## Usage

Push to public NPM registry

```sh
docker run --rm \
  -e NPM_USERNAME=drone \
  -e NPM_PASSWORD=password \
  -e NPM_EMAIL=drone@drone.io \
  plugins/npm
``

Push to private NPM registry

```sh
docker run --rm \
  -e NPM_USERNAME=drone \
  -e NPM_PASSWORD=password \
  -e NPM_EMAIL=drone@drone.io \
  -e NPM_REGISTRY=http://myregistry.com \
  -e NPM_ALWAYS_AUTH=true
  plugins/npm
``

# drone-npm

[![Build Status](http://beta.drone.io/api/badges/drone-plugins/drone-npm/status.svg)](http://beta.drone.io/drone-plugins/drone-npm)
[![Go Doc](https://godoc.org/github.com/drone-plugins/drone-npm?status.svg)](http://godoc.org/github.com/drone-plugins/drone-npm)
[![Go Report](https://goreportcard.com/badge/github.com/drone-plugins/drone-npm)](https://goreportcard.com/report/github.com/drone-plugins/drone-npm)
[![Join the chat at https://gitter.im/drone/drone](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/drone/drone)

Drone plugin to publish files and artifacts to a private or public NPM
registry. For the usage information and a listing of the available options
please take a look at [the docs](DOCS.md).

## Build

Build the binary with the following commands:

```
go build
go test
```

## Docker

Build the Docker image with the following commands:

```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo
docker build --rm=true -t plugins/npm .
```

Please note incorrectly building the image for the correct x64 linux and with
CGO disabled will result in an error when running the Docker image:

```
docker: Error response from daemon: Container command
'/bin/drone-npm' not found or does not exist..
```

## Usage

Push to public NPM registry:

```sh
docker run --rm \
  -e NPM_USERNAME=drone \
  -e NPM_PASSWORD=password \
  -e NPM_EMAIL=drone@drone.io \
  -v $(pwd):$(pwd) \
  -w $(pwd) \  
  plugins/npm
```

Push to private NPM registry:

```sh
docker run --rm \
  -e NPM_USERNAME=drone \
  -e NPM_PASSWORD=password \
  -e NPM_EMAIL=drone@drone.io \
  -e NPM_REGISTRY=http://myregistry.com \
  -e NPM_ALWAYS_AUTH=true \
  -v $(pwd):$(pwd) \
  -w $(pwd) \  
  plugins/npm
```

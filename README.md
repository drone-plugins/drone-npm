# drone-npm

[![Build Status](http://cloud.drone.io/api/badges/drone-plugins/drone-npm/status.svg)](http://cloud.drone.io/drone-plugins/drone-npm)
[![Gitter chat](https://badges.gitter.im/drone/drone.png)](https://gitter.im/drone/drone)
[![Join the discussion at https://discourse.drone.io](https://img.shields.io/badge/discourse-forum-orange.svg)](https://discourse.drone.io)
[![Drone questions at https://stackoverflow.com](https://img.shields.io/badge/drone-stackoverflow-orange.svg)](https://stackoverflow.com/questions/tagged/drone.io)
[![](https://images.microbadger.com/badges/image/plugins/npm.svg)](https://microbadger.com/images/plugins/npm "Get your own image badge on microbadger.com")
[![Go Doc](https://godoc.org/github.com/drone-plugins/drone-npm?status.svg)](http://godoc.org/github.com/drone-plugins/drone-npm)
[![Go Report](https://goreportcard.com/badge/github.com/drone-plugins/drone-npm)](https://goreportcard.com/report/github.com/drone-plugins/drone-npm)

Drone plugin to publish files and artifacts to a private or public NPM registry. For the usage information and a listing of the available options please take a look at [the docs](http://plugins.drone.io/drone-plugins/drone-npm/).

## Build

Build the binary with the following commands:

```
go build
```

## Docker

Build the Docker image with the following commands:

```
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -tags netgo -o release/linux/amd64/drone-npm
docker build --rm -t plugins/npm .
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

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

Build the binary with the following command:

```console
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0
export GO111MODULE=on

go build -v -a -tags netgo -o release/linux/amd64/drone-npm
```

## Docker

Build the Docker image with the following command:

```console
docker build \
  --label org.label-schema.build-date=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
  --label org.label-schema.vcs-ref=$(git rev-parse --short HEAD) \
  --file docker/Dockerfile.linux.amd64 --tag plugins/npm .
```

## Usage
### Standard Local Usage
```console
docker run --rm \
  -e NPM_USERNAME=drone \
  -e NPM_PASSWORD=password \
  -e NPM_EMAIL=drone@drone.io \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  plugins/npm
```
#### With a specified registry for validation
This will allow the setting of the defautl publishing registry. This will also raise a validation error if the publish configuration of the npm package is not pointing to the specified registry.
```console
docker run --rm \
  -e NPM_USERNAME=drone \
  -e NPM_PASSWORD=password \
  -e NPM_EMAIL=drone@drone.io \
  -e NPM_REGISTRY="https://fakenpm.reg.org/good/path" \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  plugins/npm
```

#### Ignore registry validation
This will all the setting of a default publishing registry but will skip the verification of it being the same as the one in the npmrc. In this instance no validation error is raised and the registry in the npm rc is used
```console
docker run --rm \
  -e NPM_USERNAME=drone \
  -e NPM_PASSWORD=password \
  -e NPM_EMAIL=drone@drone.io \
  -e NPM_REGISTRY="https://fakenpm.reg.org/good/path" \
  -e PLUGIN_SKIP_REGISTRY_VALIDATION=true \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  plugins/npm
```
# Docker image for the Drone NPM plugin
#
#     cd $GOPATH/src/github.com/drone-plugins/drone-npm
#     CGO_ENABLED=0 go build -a -tags netgo
#     docker build --rm=true -t plugins/drone-npm .

FROM alpine:3.2
RUN apk add -U ca-certificates git nodejs && rm -rf /var/cache/apk/*
ADD drone-npm /bin/
ENTRYPOINT ["/bin/drone-npm"]

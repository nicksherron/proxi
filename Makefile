.PHONY: all docker-build build clean test help default test-verify run build-darwin build-linux

BIN_NAME=proxi
# dev build for latest release see releases at
VERSION := 99.99.99
GIT_COMMIT=$(shell git rev-parse HEAD)
GIT_DIRTY=$(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)
BUILD_DATE=$(shell date '+%Y-%m-%d-%H:%M:%S')

PWD := $(shell pwd)
GOPATH := $(shell go env GOPATH)
GOARCH := $(shell go env GOARCH)
GOOS := $(shell go env GOOS)

default: help

all: build run

help:
	@echo 'Management commands for proxi:'
	@echo
	@echo 'Usage:'
	@echo '    make docker-build    Build docker image.'
	@echo '    make build           Compile the project.'
	@echo '    make test            Run tests on a compiled project.'
	@echo '    make test-providers  Run tests and verify providers return results instead of just checking format.'
	@echo '    make clean           Clean the directory tree.'
	@echo


docker-build:
	docker build --no-cache=true --build-arg VERSION=${VERSION} --build-arg BUILD_DATE=${BUILD_DATE} --build-arg GIT_COMMIT=${GIT_COMMIT} -t proxi .

build:
	@echo "building ${BIN_NAME} ${VERSION}"
	@echo "GOPATH=${GOPATH}"
	go build -ldflags "-X github.com/nicksherron/proxi/internal.Version=${VERSION} -X github.com/nicksherron/proxi/cmd.Build=${GIT_COMMIT}${GIT_DIRTY} -X github.com/nicksherron/proxi/cmd.BuildDate=${BUILD_DATE}" -o bin/${BIN_NAME}


run:
	bin/${BIN_NAME} server -p -d=1

clean:
	@test ! -e bin/${BIN_NAME} || rm bin/${BIN_NAME}

test:
	go test ./...

race:
	go run --race *.go server -d=10

test-providers:
	go test -v  github.com/nicksherron/proxi/internal   -args -verify

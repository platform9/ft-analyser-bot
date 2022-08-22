# Usage:
# make build                # builds the artifact
# make clean           # removes the artifact and the vendored packages

SHELL := /usr/bin/env bash
GITHASH := $(shell git rev-parse --short HEAD)
BIN_DIR := $(shell pwd)/bin
CMD_DIR := $(shell pwd)/cmd
BIN := ft-analyser-bot
REPO := ft-analyser-bot
LDFLAGS := ""
NAME ?= ft-analyser-bot# default lable for docker build
TAG ?= latest
XDG_CACHE_HOME := /tmp
CONT_USER := $(shell id -u)
CONT_GRP := $(shell id -g)
GOFLAGS ?= ""
PACKAGE_GOPATH := /go/src/github.com/platform9/$(REPO)

.PHONY: clean format test build container-build docker-build

default: clean format test build container-build docker-build


container-build:
	docker run --rm --env XDG_CACHE_HOME=$(XDG_CACHE_HOME)  --env GOPATH=/tmp --env GOFLAGS=$(GOFLAGS) --user $(CONT_USER):$(CONT_GRP) --volume $(PWD):$(PACKAGE_GOPATH) --workdir $(PACKAGE_GOPATH) golang:1.18 make

docker-build:
	docker build -t $(NAME):$(TAG) .

format:
	gofmt -w -s */*.go

clean:
	rm -rf $(BIN_DIR)

build:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -a -o $(BIN_DIR)/$(BIN) $(CMD_DIR)/main.go

test:
	go test -v ./...

.DEFAULT_GOAL := build
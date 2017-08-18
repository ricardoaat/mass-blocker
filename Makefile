.PHONY: \
	build

GO_PATH = $(shell echo $(GOPATH) | awk -F':' '{print $$1}')
GO_SRC = $(shell pwd | xargs dirname | xargs dirname | xargs dirname)
DEPLOY_PATH := ~/ric/dev/go/mass-blocker/compiled/
BIN_NAME :=massblocker
VERSION=1.0
BUILD=$(shell date +%FT%T%z)
LDFLAGS=-ldflags "-X main.version=${VERSION} -X main.buildDate=${BUILD}"
build:
	go build -i -v $(LDFLAGS) -o $(DEPLOY_PATH)$(BIN_NAME) main.go

install:
	go install -i -v $(LDFLAGS) -o $(DEPLOY_PATH)$(BIN_NAME) main.go

.PHONY: test $(PKGS) clean vendor
VERSION := $(shell cat VERSION)
SHELL := /bin/bash
PKG := github.com/Clever/whackanop
PKGS := $(shell go list ./... | grep -v /vendor)
EXECUTABLE := whackanop
BUILDS := \
	build/$(EXECUTABLE)-v$(VERSION)-darwin-amd64 \
	build/$(EXECUTABLE)-v$(VERSION)-linux-amd64
COMPRESSED_BUILDS := $(BUILDS:%=%.tar.gz)
RELEASE_ARTIFACTS := $(COMPRESSED_BUILDS:build/%=release/%)
GODEP := $(GOPATH)/bin/godep

GOVERSION := $(shell go version | grep 1.5)
ifeq "$(GOVERSION)" ""
  $(error must be running Go version 1.5)
endif

export GO15VENDOREXPERIMENT = 1

test: $(PKGS)

$(GOPATH)/bin/golint:
	go get github.com/golang/lint/golint

$(PKGS): version.go $(GOPATH)/bin/golint
	$(GOPATH)/bin/golint $(GOPATH)/src/$@*/**.go
	go get -d -t $@
	go test -cover -coverprofile=$(GOPATH)/src/$@/c.out $@ -test.v
ifeq ($(HTMLCOV),1)
	go tool cover -html=$(GOPATH)/src/$@/c.out
endif

build/*: version.go
version.go: VERSION
	echo 'package main' > version.go
	echo '' >> version.go # Write a go file that lints :)
	echo 'const Version = "$(VERSION)"' >> version.go

build/$(EXECUTABLE)-v$(VERSION)-darwin-amd64:
	GOARCH=amd64 GOOS=darwin go build -o "$@/$(EXECUTABLE)"
build/$(EXECUTABLE)-v$(VERSION)-linux-amd64:
	GOARCH=amd64 GOOS=linux go build -o "$@/$(EXECUTABLE)"

%.tar.gz: %
	tar -C `dirname $<` -zcvf "$<.tar.gz" `basename $<`

$(RELEASE_ARTIFACTS): release/% : build/%
	mkdir -p release
	cp $< $@

release: $(RELEASE_ARTIFACTS)

clean:
	rm -rf build release


$(GODEP):
	go get -u github.com/tools/godep

vendor: $(GODEP)
	$(GODEP) save $(PKGS)
	find vendor/ -path '*/vendor' -type d | xargs -IX rm -r X # remove any nested vendor directories

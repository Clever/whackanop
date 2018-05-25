include golang.mk
.DEFAULT_GOAL := test # override default goal set in library makefile

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

$(eval $(call golang-version-check,1.10))

test: $(PKGS)

$(PKGS): version.go golang-test-all-deps
	$(call golang-test-all,$@)

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

install_deps: golang-dep-vendor-deps
	$(call golang-dep-vendor)

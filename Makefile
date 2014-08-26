SHELL := /bin/bash
PKG = github.com/Clever/whackanop
SUBPKGS =
PKGS = $(PKG) $(SUBPKGS)

.PHONY: test $(PKGS)

test: $(PKG)

$(PKG):
	go get github.com/golang/lint/golint
	$(GOPATH)/bin/golint $(GOPATH)/src/$@*/**.go
	go get -d -t $@
	go test -cover -coverprofile=$(GOPATH)/src/$@/c.out $@ -test.v
ifeq ($(HTMLCOV),1)
	go tool cover -html=$(GOPATH)/src/$@/c.out
endif

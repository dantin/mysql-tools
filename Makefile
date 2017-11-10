TOOLS_PKG := github.com/dantin/mysql-tools
VENDOR := $(shell pwd)/_vendor

LDFLAGS += -X "$(TOOLS_PKG)/pkg/utils.Version=0.0.1+git.$(shell git rev-parse --short HEAD)"
LDFLAGS += -X "$(TOOLS_PKG)/pkg/utils.BuildTS=$(shell date -u '+%Y-%m-%d %I:%M:%S')"
LDFLAGS += -X "$(TOOLS_PKG)/pkg/utils.GitHash=$(shell git rev-parse HEAD)"
LDFLAGS += -X "$(TOOLS_PKG)/pkg/utils.GitBranch=$(shell git rev-parse --abbrev-ref HEAD)"

.PHONY: update clean

default: build

build:
	GOPATH=$(VENDOR) go build -ldflags '$(LDFLAGS)' -o bin/drc cmd/drc/main.go

update:
	which glide >/dev/null || curl https://glide.sh/get | sh
	which glide-vc || go get -v -u github.com/sgotti/glide-vc
	rm -rf vendor && mv _vendor/src vendor || true
	rm -rf _vendor
	glide update --strip-vendor --skip-test
	@echo "removing test files"
	glide vc --only-code --no-tests
	mkdir -p _vendor
	mv vendor _vendor/src
	mkdir -p _vendor/src/$(shell dirname $(TOOLS_PKG))
	ln -s $(shell pwd) _vendor/src/$(TOOLS_PKG)

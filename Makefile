TOOLS_PKG := github.com/dantin/mysql-tools
VENDOR    := $(shell pwd)/_vendor

LDFLAGS += -X "$(TOOLS_PKG)/pkg/utils.Version=0.0.1+git.$(shell git rev-parse --short HEAD)"
LDFLAGS += -X "$(TOOLS_PKG)/pkg/utils.BuildTS=$(shell date -u '+%Y-%m-%d %I:%M:%S')"
LDFLAGS += -X "$(TOOLS_PKG)/pkg/utils.GitHash=$(shell git rev-parse HEAD)"
LDFLAGS += -X "$(TOOLS_PKG)/pkg/utils.GitBranch=$(shell git rev-parse --abbrev-ref HEAD)"

#TEST_DIRS := $(shell find . -iname "*_test.go" -exec dirname {} \; | uniq)

GOFILTER  := grep -vE 'vendor'
GOCHECKER := $(GOFILTER) | awk '{ print } END { if (NR > 0) { exit 1 } }'

GO      := GO15VENDOREXPERIMENT="1" go
GOBUILD := $(GO) build
GOTEST  := $(GO) test

PACKAGES := $$(go list ./...| grep -vE 'vendor')

.PHONY: update clean check build test setup-check

default: build

build:
	$(GOBUILD) -ldflags '$(LDFLAGS)' -o bin/drc cmd/drc/main.go

test:
	@echo "test"
	@$(GOTEST) --race --cover $(PACKAGES)

check:
	@echo "gofmt"
	@ gofmt -s -l . 2>&1 | $(GOCHECKER)

update:
	@ which glide >/dev/null || curl https://glide.sh/get | sh
	@ which glide-vc || go get -v -u github.com/sgotti/glide-vc
	glide update --strip-vendor --skip-test
	@echo "removing test files"
	glide vc --only-code --no-tests

setup-check:
	@ go get -u gopkg.in/check.v1

clean:
	rm -rf bin/*

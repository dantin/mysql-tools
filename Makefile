TOOLS_PKG := github.com/dantin/mysql-tools
VENDOR := $(shell pwd)/_vendor

LDFLAGS += -X "$(TOOLS_PKG)/pkg/utils.Version=0.0.1+git.$(shell git rev-parse --short HEAD)"
LDFLAGS += -X "$(TOOLS_PKG)/pkg/utils.BuildTS=$(shell date -u '+%Y-%m-%d %I:%M:%S')"
LDFLAGS += -X "$(TOOLS_PKG)/pkg/utils.GitHash=$(shell git rev-parse HEAD)"
LDFLAGS += -X "$(TOOLS_PKG)/pkg/utils.GitBranch=$(shell git rev-parse --abbrev-ref HEAD)"

GOFILTER := grep -vE 'vendor'
GOCHECKER := $(GOFILTER) | awk '{ print } END { if (NR > 0) { exit 1 } }'

GO := GO15VENDOREXPERIMENT="1" go

.PHONY: update clean check

default: build

build:
	$(GO) build -ldflags '$(LDFLAGS)' -o bin/drc cmd/drc/main.go

check:
	@echo "gofmt"
	@ gofmt -s -l . 2>&1 | $(GOCHECKER)

update:
	which glide >/dev/null || curl https://glide.sh/get | sh
	which glide-vc || go get -v -u github.com/sgotti/glide-vc
	glide update --strip-vendor --skip-test
	@echo "removing test files"
	glide vc --only-code --no-tests

clean:
	rm -rf bin/*

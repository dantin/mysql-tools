TOOLS_PKG := github.com/dantin/mysql-tools

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

.PHONY: update clean check build test init up down

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
	glide update --strip-vendor --skip-test
	@echo "removing test files"
	glide vc --only-code --no-tests

up:
	@docker-compose -f docker-compose.yml up -d --force-recreate
	@until mysql -h 127.0.0.1 -u root -e 'select 1' >/dev/null 2>&1; do sleep 1; echo "Waiting for DB to come up..."; done
	bin/drc -config resources/conf/syncer.toml

down:
	@docker-compose -f docker-compose.yml down --remove-orphans

init:
	@ which glide >/dev/null || curl https://glide.sh/get | sh
	@ which glide-vc >/dev/null || go get -v -u github.com/sgotti/glide-vc
	@echo "update testing framework"
	@ go get -u gopkg.in/check.v1

clean:
	rm -rf bin/*

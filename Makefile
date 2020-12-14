PKG := github.com/ExchangeUnion/xud-launcher

GO_BIN := ${GOPATH}/bin

GOBUILD := go build -v

VERSION := v1.0.0
COMMIT := $(shell git rev-parse HEAD)
ifeq ($(OS),Windows_NT)
	TIMESTAMP := $(shell powershell.exe scripts\get_timestamp.ps1)
else
	TIMESTAMP := $(shell date +%s)
endif
LDFLAGS := -ldflags "-w -s \
-X $(PKG)/build.Version=$(VERSION) \
-X $(PKG)/build.GitCommit=$(COMMIT) \
-X $(PKG)/build.Timestamp=$(TIMESTAMP)"

default: build


#
# Building
#

build:
	$(GOBUILD) -o xud-launcher $(LDFLAGS)

.PHONY: build

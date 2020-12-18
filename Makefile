PKG := github.com/ExchangeUnion/xud-launcher

GO_BIN := ${GOPATH}/bin

GOBUILD := go build -v

VERSION := local
COMMIT := $(shell git rev-parse HEAD)
ifeq ($(OS),Windows_NT)
	TIMESTAMP := $(shell powershell.exe scripts\get_timestamp.ps1)
else
	TIMESTAMP := $(shell date +%s)
endif

ifeq ($(GOOS), windows)
	OUTPUT := xud-launcher.exe
else
	OUTPUT := xud-launcher
endif


LDFLAGS := -ldflags "-w -s \
-X $(PKG)/build.Version=$(VERSION) \
-X $(PKG)/build.GitCommit=$(COMMIT) \
-X $(PKG)/build.Timestamp=$(TIMESTAMP)"

default: build

build:
	$(GOBUILD) $(LDFLAGS)

zip:
	zip --junk-paths xud-launcher.zip $(OUTPUT)

clean:
	rm -f xud-launcher
	rm -f xud-launcher.zip

.PHONY: build

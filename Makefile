.PHONY: build api 
PKGS := $(shell go list ./... | grep -v /vendor/ | grep -v numberserver/api | grep -v /migrations | grep -v /static)
VERSION := $(shell git describe --always)
GOOS ?= linux
GOARCH ?= $(arch) 

build:
	@echo "Compiling source for $(GOOS) $(GOARCH)"
	@mkdir -p build
	@GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags "-X main.version=$(VERSION)" -o build/numberserver_$(arch)$(BINEXT) cmd/numberserver/main.go

clean:
	@echo "Cleaning up workspace"
	@rm -rf build

api:
	@echo "Generating API code from .proto files"
	@go generate api/ns/ns.go

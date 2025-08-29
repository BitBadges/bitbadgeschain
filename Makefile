# Version must be provided as a CLI argument
VERSION ?= $(error VERSION is required. Usage: make build-linux/amd64 VERSION=v13)

# Common ldflags for version information
LDFLAGS := -X github.com/bitbadges/bitbadgeschain/app.Version=$(VERSION) \
	-X github.com/bitbadges/bitbadgeschain/app.Commit=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown") \
	-X github.com/bitbadges/bitbadgeschain/app.BuildTags=$(shell go list -f '{{.BuildTags}}' ./... 2>/dev/null | head -1 | tr ' ' ',' | sed 's/,$$//' || echo "") \
	-X github.com/bitbadges/bitbadgeschain/app.GoVersion=$(shell go version 2>/dev/null | sed 's/^go version //' | cut -d' ' -f1 || echo "unknown") \
	-X github.com/bitbadges/bitbadgeschain/app.CosmosSDKVersion=v0.50.13

build-linux/amd64:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o ./build/bitbadgeschain-linux-amd64 ./cmd/bitbadgeschaind/main.go

build-linux/arm64:
	CGO_ENABLED=1 CC=aarch64-linux-gnu-gcc GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o ./build/bitbadgeschain-linux-arm64 ./cmd/bitbadgeschaind/main.go

build-darwin:
	CGO_ENABLED=1 CC="o64-clang" GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o ./build/bitbadgeschain-darwin-amd64 ./cmd/bitbadgeschaind/main.go

build-all: 
	make build-linux/amd64
	make build-linux/arm64

do-checksum:
	cd build && sha256sum bitbadgeschain-linux-amd64 bitbadgeschain-linux-arm64 > bitbadgeschain_checksum

build-with-checksum: build-all do-checksum
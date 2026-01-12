# Version must be provided as a CLI argument
VERSION := v22

# Common ldflags for version information
LDFLAGS := -X github.com/cosmos/cosmos-sdk/version.Name=bitbadgeschain \
	-X github.com/cosmos/cosmos-sdk/version.AppName=bitbadgeschaind \
	-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown") \
	-X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(shell go list -f '{{.BuildTags}}' ./... 2>/dev/null | head -1 | tr ' ' ',' | sed 's/,$$//' || echo "")"

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

lint:
	@echo "Running golangci-lint (excluding test files, ai_test, pb.go, and simulation files)..."
	@golangci-lint run --skip-dirs='ai_test|simulation' --skip-files='.*_test\.go$$|.*\.pb\.go$$|.*\.pb\.gw\.go$$|.*\.pulsar\.go$$|.*_gen\.go$$' ./x/badges/... ./x/custom-hooks/... || ( \
		echo ""; \
		echo "If you see a 'Go language version' error, rebuild golangci-lint with:"; \
		echo "  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		echo "This project uses Go 1.24.5, so golangci-lint must be built with Go 1.24.5 or later."; \
		exit 1 \
	)

lint-security:
	@echo "Running security & bug-focused linting..."
	@golangci-lint run --config .golangci-security.yml ./x/badges/... ./x/custom-hooks/... || ( \
		echo ""; \
		echo "If you see a 'Go language version' error, rebuild golangci-lint with:"; \
		echo "  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		exit 1 \
	)

lint-fix:
	@echo "Running golangci-lint with --fix (excluding test files, ai_test, pb.go, and simulation files)..."
	@golangci-lint run --fix ./x/badges/... ./x/custom-hooks/... || ( \
		echo ""; \
		echo "If you see a 'Go language version' error, rebuild golangci-lint with:"; \
		echo "  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		exit 1 \
	)

lint-ci:
	@echo "Running golangci-lint for CI (excluding test files, ai_test, pb.go, and simulation files)..."
	@golangci-lint run --out-format=github-actions ./x/badges/... ./x/custom-hooks/... || ( \
		echo ""; \
		echo "If you see a 'Go language version' error, rebuild golangci-lint with:"; \
		echo "  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		exit 1 \
	)
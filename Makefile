# Version must be provided as a CLI argument
VERSION := v24

# Common ldflags for version information
LDFLAGS := -X github.com/cosmos/cosmos-sdk/version.Name=bitbadgeschain \
	-X github.com/cosmos/cosmos-sdk/version.AppName=bitbadgeschaind \
	-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown") \
	-X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(shell go list -f '{{.BuildTags}}' ./... 2>/dev/null | head -1 | tr ' ' ',' | sed 's/,$$//' || echo "")"

# EVM Chain ID ldflags
# Mainnet: 50024
LDFLAGS_MAINNET := $(LDFLAGS) \
	-X github.com/bitbadges/bitbadgeschain/app/params.BuildTimeEVMChainID=50024

# Testnet: 50025
LDFLAGS_TESTNET := $(LDFLAGS) \
	-X github.com/bitbadges/bitbadgeschain/app/params.BuildTimeEVMChainID=50025

# Generic build (no chain ID set - defaults to 90123 for local development)
build-linux/amd64:
	@echo "Building binary (EVM Chain ID: 90123 - local dev) for linux/amd64..."
	@mkdir -p ./build
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o ./build/bitbadgeschain-linux-amd64 ./cmd/bitbadgeschaind/main.go
	@echo "✓ Built: ./build/bitbadgeschain-linux-amd64"

build-linux/arm64:
	@echo "Building binary (EVM Chain ID: 90123 - local dev) for linux/arm64..."
	@mkdir -p ./build
	CGO_ENABLED=1 CC=aarch64-linux-gnu-gcc GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o ./build/bitbadgeschain-linux-arm64 ./cmd/bitbadgeschaind/main.go
	@echo "✓ Built: ./build/bitbadgeschain-linux-arm64"

build-darwin:
	@echo "Building binary (EVM Chain ID: 90123 - local dev) for darwin/amd64..."
	@mkdir -p ./build
	CGO_ENABLED=1 CC="o64-clang" GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o ./build/bitbadgeschain-darwin-amd64 ./cmd/bitbadgeschaind/main.go
	@echo "✓ Built: ./build/bitbadgeschain-darwin-amd64"

# Mainnet builds (with EVM Chain ID 50024 compiled in)
build-mainnet-linux/amd64:
	@echo "Building mainnet binary (EVM Chain ID: 50024) for linux/amd64..."
	@mkdir -p ./build
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS_MAINNET)" -o ./build/bitbadgeschain-linux-amd64 ./cmd/bitbadgeschaind/main.go
	@echo "✓ Built: ./build/bitbadgeschain-linux-amd64"

build-mainnet-linux/arm64:
	@echo "Building mainnet binary (EVM Chain ID: 50024) for linux/arm64..."
	@mkdir -p ./build
	CGO_ENABLED=1 CC=aarch64-linux-gnu-gcc GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS_MAINNET)" -o ./build/bitbadgeschain-linux-arm64 ./cmd/bitbadgeschaind/main.go
	@echo "✓ Built: ./build/bitbadgeschain-linux-arm64"

build-mainnet-darwin:
	@echo "Building mainnet binary (EVM Chain ID: 50024) for darwin/amd64..."
	@mkdir -p ./build
	CGO_ENABLED=1 CC="o64-clang" GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS_MAINNET)" -o ./build/bitbadgeschain-darwin-amd64 ./cmd/bitbadgeschaind/main.go
	@echo "✓ Built: ./build/bitbadgeschain-darwin-amd64"

# Testnet builds (with EVM Chain ID 50025 compiled in)
build-testnet-linux/amd64:
	@echo "Building testnet binary (EVM Chain ID: 50025) for linux/amd64..."
	@mkdir -p ./build
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS_TESTNET)" -o ./build/bitbadgeschain-testnet-linux-amd64 ./cmd/bitbadgeschaind/main.go
	@echo "✓ Built: ./build/bitbadgeschain-testnet-linux-amd64"

build-testnet-linux/arm64:
	@echo "Building testnet binary (EVM Chain ID: 50025) for linux/arm64..."
	@mkdir -p ./build
	CGO_ENABLED=1 CC=aarch64-linux-gnu-gcc GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS_TESTNET)" -o ./build/bitbadgeschain-testnet-linux-arm64 ./cmd/bitbadgeschaind/main.go
	@echo "✓ Built: ./build/bitbadgeschain-testnet-linux-arm64"

build-testnet-darwin:
	@echo "Building testnet binary (EVM Chain ID: 50025) for darwin/amd64..."
	@mkdir -p ./build
	CGO_ENABLED=1 CC="o64-clang" GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS_TESTNET)" -o ./build/bitbadgeschain-testnet-darwin-amd64 ./cmd/bitbadgeschaind/main.go
	@echo "✓ Built: ./build/bitbadgeschain-testnet-darwin-amd64"

# build-all builds all 4 production binaries:
# - 2 mainnet binaries (bitbadgeschain-linux-amd64, bitbadgeschain-linux-arm64) with chain ID 50024
# - 2 testnet binaries (bitbadgeschain-testnet-linux-amd64, bitbadgeschain-testnet-linux-arm64) with chain ID 50025
build-all: 
	make build-mainnet-linux/amd64
	make build-mainnet-linux/arm64
	make build-testnet-linux/amd64
	make build-testnet-linux/arm64

build-all-mainnet:
	make build-mainnet-linux/amd64
	make build-mainnet-linux/arm64

build-all-testnet:
	make build-testnet-linux/amd64
	make build-testnet-linux/arm64

do-checksum-testnet:
	cd build && sha256sum bitbadgeschain-testnet-linux-amd64 bitbadgeschain-testnet-linux-arm64 > bitbadgeschain-testnet_checksum

do-checksum-mainnet:
	cd build && sha256sum bitbadgeschain-linux-amd64 bitbadgeschain-linux-arm64 > bitbadgeschain-mainnet_checksum

# Checksum for all 4 binaries
do-checksum:
	cd build && sha256sum bitbadgeschain-linux-amd64 bitbadgeschain-linux-arm64 bitbadgeschain-testnet-linux-amd64 bitbadgeschain-testnet-linux-arm64 > bitbadgeschain_checksum

build-with-checksum: build-all do-checksum

lint:
	@echo "Running golangci-lint (excluding test files, ai_test, pb.go, and simulation files)..."
	@golangci-lint run --skip-dirs='ai_test|simulation' --skip-files='.*_test\.go$$|.*\.pb\.go$$|.*\.pb\.gw\.go$$|.*\.pulsar\.go$$|.*_gen\.go$$' ./x/tokenization/... ./x/custom-hooks/... || ( \
		echo ""; \
		echo "If you see a 'Go language version' error, rebuild golangci-lint with:"; \
		echo "  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		echo "This project uses Go 1.24.5, so golangci-lint must be built with Go 1.24.5 or later."; \
		exit 1 \
	)

lint-security:
	@echo "Running security & bug-focused linting..."
	@golangci-lint run --config .golangci-security.yml ./x/tokenization/... ./x/custom-hooks/... || ( \
		echo ""; \
		echo "If you see a 'Go language version' error, rebuild golangci-lint with:"; \
		echo "  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		exit 1 \
	)

lint-fix:
	@echo "Running golangci-lint with --fix (excluding test files, ai_test, pb.go, and simulation files)..."
	@golangci-lint run --fix ./x/tokenization/... ./x/custom-hooks/... || ( \
		echo ""; \
		echo "If you see a 'Go language version' error, rebuild golangci-lint with:"; \
		echo "  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		exit 1 \
	)

lint-ci:
	@echo "Running golangci-lint for CI (excluding test files, ai_test, pb.go, and simulation files)..."
	@golangci-lint run --out-format=github-actions ./x/tokenization/... ./x/custom-hooks/... || ( \
		echo ""; \
		echo "If you see a 'Go language version' error, rebuild golangci-lint with:"; \
		echo "  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		exit 1 \
	)

.PHONY: compile-contracts
compile-contracts:
	@echo "Compiling Solidity contracts..."
	@if command -v solcjs >/dev/null 2>&1; then \
		echo "Using solcjs..."; \
		cd contracts && \
		solcjs --bin --abi --base-path . --include-path . --output-dir test \
			types/TokenizationTypes.sol \
			interfaces/ITokenizationPrecompile.sol \
			test/MinimalTestContract.sol && \
		echo "Contracts compiled successfully"; \
	elif command -v solc >/dev/null 2>&1; then \
		echo "Using solc..."; \
		cd contracts && \
		solc --combined-json bin,abi --allow-paths . \
			types/TokenizationTypes.sol \
			interfaces/ITokenizationPrecompile.sol \
			test/MinimalTestContract.sol > test/MinimalTestContract.json && \
		echo "Contracts compiled successfully"; \
	else \
		echo "Error: Neither solcjs nor solc found. Please install one:"; \
		echo "  npm install -g solc@0.8.24"; \
		echo "  or"; \
		echo "  Install solc via your package manager"; \
		exit 1; \
	fi
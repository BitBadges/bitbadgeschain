# Version must be provided as a CLI argument
VERSION := v32

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

# Cross-compile CC detection:
# On macOS: uses native cc (supports both amd64 and arm64).
# On Linux: uses osxcross (o64-clang / oa64-clang). Requires macOS SDK.
# Override with: make build-mainnet-darwin/amd64 CC_DARWIN_AMD64=o64-clang
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Darwin)
  CC_DARWIN_AMD64 ?= cc
  CC_DARWIN_ARM64 ?= cc
else
  CC_DARWIN_AMD64 ?= o64-clang
  CC_DARWIN_ARM64 ?= oa64-clang
endif
CC_WINDOWS_AMD64 ?= x86_64-w64-mingw32-gcc

# Generic build (no chain ID set - defaults to 90123 for local development)
# Alias to mainnet builds for production use
build-linux/amd64: build-mainnet-linux/amd64

build-linux/arm64: build-mainnet-linux/arm64

# Local development builds (with EVM Chain ID 90123 for local dev)
build-local-linux/amd64:
	@echo "Building binary (EVM Chain ID: 90123 - local dev) for linux/amd64..."
	@mkdir -p ./build
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "$(LDFLAGS)" -o ./build/bitbadgeschain-linux-amd64 ./cmd/bitbadgeschaind/main.go
	@echo "✓ Built: ./build/bitbadgeschain-linux-amd64"

build-local-linux/arm64:
	@echo "Building binary (EVM Chain ID: 90123 - local dev) for linux/arm64..."
	@mkdir -p ./build
	CGO_ENABLED=1 CC=aarch64-linux-gnu-gcc GOOS=linux GOARCH=arm64 go build -trimpath -ldflags "$(LDFLAGS)" -o ./build/bitbadgeschain-linux-arm64 ./cmd/bitbadgeschaind/main.go
	@echo "✓ Built: ./build/bitbadgeschain-linux-arm64"

build-local-darwin/amd64:
	@echo "Building binary (EVM Chain ID: 90123 - local dev) for darwin/amd64..."
	@mkdir -p ./build
	CGO_ENABLED=1 CC="$(CC_DARWIN_AMD64)" GOOS=darwin GOARCH=amd64 go build -trimpath -ldflags "$(LDFLAGS)" -o ./build/bitbadgeschain-darwin-amd64 ./cmd/bitbadgeschaind/main.go
	@echo "✓ Built: ./build/bitbadgeschain-darwin-amd64"

build-local-darwin/arm64:
	@echo "Building binary (EVM Chain ID: 90123 - local dev) for darwin/arm64..."
	@mkdir -p ./build
	CGO_ENABLED=1 CC="$(CC_DARWIN_ARM64)" GOOS=darwin GOARCH=arm64 go build -trimpath -ldflags "$(LDFLAGS)" -o ./build/bitbadgeschain-darwin-arm64 ./cmd/bitbadgeschaind/main.go
	@echo "✓ Built: ./build/bitbadgeschain-darwin-arm64"

build-local-windows/amd64:
	@echo "Building binary (EVM Chain ID: 90123 - local dev) for windows/amd64..."
	@mkdir -p ./build
	CGO_ENABLED=1 CC="$(CC_WINDOWS_AMD64)" GOOS=windows GOARCH=amd64 go build -trimpath -ldflags "$(LDFLAGS)" -o ./build/bitbadgeschain-windows-amd64.exe ./cmd/bitbadgeschaind/main.go
	@echo "✓ Built: ./build/bitbadgeschain-windows-amd64.exe"

# Legacy aliases
build-local-darwin: build-local-darwin/amd64
build-darwin: build-local-darwin/amd64

# Mainnet builds (with EVM Chain ID 50024 compiled in)
build-mainnet-linux/amd64:
	@echo "Building mainnet binary (EVM Chain ID: 50024) for linux/amd64..."
	@mkdir -p ./build
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "$(LDFLAGS_MAINNET)" -o ./build/bitbadgeschain-linux-amd64 ./cmd/bitbadgeschaind/main.go
	@echo "✓ Built: ./build/bitbadgeschain-linux-amd64"

build-mainnet-linux/arm64:
	@echo "Building mainnet binary (EVM Chain ID: 50024) for linux/arm64..."
	@mkdir -p ./build
	CGO_ENABLED=1 CC=aarch64-linux-gnu-gcc GOOS=linux GOARCH=arm64 go build -trimpath -ldflags "$(LDFLAGS_MAINNET)" -o ./build/bitbadgeschain-linux-arm64 ./cmd/bitbadgeschaind/main.go
	@echo "✓ Built: ./build/bitbadgeschain-linux-arm64"

build-mainnet-darwin/amd64:
	@echo "Building mainnet binary (EVM Chain ID: 50024) for darwin/amd64..."
	@mkdir -p ./build
	CGO_ENABLED=1 CC="$(CC_DARWIN_AMD64)" GOOS=darwin GOARCH=amd64 go build -trimpath -ldflags "$(LDFLAGS_MAINNET)" -o ./build/bitbadgeschain-darwin-amd64 ./cmd/bitbadgeschaind/main.go
	@echo "✓ Built: ./build/bitbadgeschain-darwin-amd64"

build-mainnet-darwin/arm64:
	@echo "Building mainnet binary (EVM Chain ID: 50024) for darwin/arm64..."
	@mkdir -p ./build
	CGO_ENABLED=1 CC="$(CC_DARWIN_ARM64)" GOOS=darwin GOARCH=arm64 go build -trimpath -ldflags "$(LDFLAGS_MAINNET)" -o ./build/bitbadgeschain-darwin-arm64 ./cmd/bitbadgeschaind/main.go
	@echo "✓ Built: ./build/bitbadgeschain-darwin-arm64"

build-mainnet-windows/amd64:
	@echo "Building mainnet binary (EVM Chain ID: 50024) for windows/amd64..."
	@mkdir -p ./build
	CGO_ENABLED=1 CC="$(CC_WINDOWS_AMD64)" GOOS=windows GOARCH=amd64 go build -trimpath -ldflags "$(LDFLAGS_MAINNET)" -o ./build/bitbadgeschain-windows-amd64.exe ./cmd/bitbadgeschaind/main.go
	@echo "✓ Built: ./build/bitbadgeschain-windows-amd64.exe"

# Legacy alias
build-mainnet-darwin: build-mainnet-darwin/amd64

# Testnet builds (with EVM Chain ID 50025 compiled in)
build-testnet-linux/amd64:
	@echo "Building testnet binary (EVM Chain ID: 50025) for linux/amd64..."
	@mkdir -p ./build
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "$(LDFLAGS_TESTNET)" -o ./build/bitbadgeschain-testnet-linux-amd64 ./cmd/bitbadgeschaind/main.go
	@echo "✓ Built: ./build/bitbadgeschain-testnet-linux-amd64"

build-testnet-linux/arm64:
	@echo "Building testnet binary (EVM Chain ID: 50025) for linux/arm64..."
	@mkdir -p ./build
	CGO_ENABLED=1 CC=aarch64-linux-gnu-gcc GOOS=linux GOARCH=arm64 go build -trimpath -ldflags "$(LDFLAGS_TESTNET)" -o ./build/bitbadgeschain-testnet-linux-arm64 ./cmd/bitbadgeschaind/main.go
	@echo "✓ Built: ./build/bitbadgeschain-testnet-linux-arm64"

build-testnet-darwin/amd64:
	@echo "Building testnet binary (EVM Chain ID: 50025) for darwin/amd64..."
	@mkdir -p ./build
	CGO_ENABLED=1 CC="$(CC_DARWIN_AMD64)" GOOS=darwin GOARCH=amd64 go build -trimpath -ldflags "$(LDFLAGS_TESTNET)" -o ./build/bitbadgeschain-testnet-darwin-amd64 ./cmd/bitbadgeschaind/main.go
	@echo "✓ Built: ./build/bitbadgeschain-testnet-darwin-amd64"

build-testnet-darwin/arm64:
	@echo "Building testnet binary (EVM Chain ID: 50025) for darwin/arm64..."
	@mkdir -p ./build
	CGO_ENABLED=1 CC="$(CC_DARWIN_ARM64)" GOOS=darwin GOARCH=arm64 go build -trimpath -ldflags "$(LDFLAGS_TESTNET)" -o ./build/bitbadgeschain-testnet-darwin-arm64 ./cmd/bitbadgeschaind/main.go
	@echo "✓ Built: ./build/bitbadgeschain-testnet-darwin-arm64"

build-testnet-windows/amd64:
	@echo "Building testnet binary (EVM Chain ID: 50025) for windows/amd64..."
	@mkdir -p ./build
	CGO_ENABLED=1 CC="$(CC_WINDOWS_AMD64)" GOOS=windows GOARCH=amd64 go build -trimpath -ldflags "$(LDFLAGS_TESTNET)" -o ./build/bitbadgeschain-testnet-windows-amd64.exe ./cmd/bitbadgeschaind/main.go
	@echo "✓ Built: ./build/bitbadgeschain-testnet-windows-amd64.exe"

# Legacy alias
build-testnet-darwin: build-testnet-darwin/amd64

# build-all builds all production binaries for all platforms:
# - Mainnet (chain ID 50024): linux/amd64, linux/arm64, darwin/amd64, darwin/arm64, windows/amd64
# - Testnet (chain ID 50025): linux/amd64, linux/arm64, darwin/amd64, darwin/arm64, windows/amd64
#
# Note: Darwin targets require macOS SDK headers (CGO deps: IOKit, CoreFoundation, secp256k1).
#   - On macOS: works natively (CC=cc)
#   - On Linux: requires osxcross. zig cc alone is NOT sufficient (missing Apple frameworks).
#   - Recommended: build darwin binaries on macOS, everything else on Linux.
#
# Quick reference:
#   make build-all              # All platforms (needs macOS SDK for darwin)
#   make build-all-cross        # Linux + Windows only (no macOS SDK needed)
#   make build-all-linux        # Linux only
build-all:
	make build-all-mainnet
	make build-all-testnet

# Linux + Windows (buildable from Linux without macOS SDK)
build-all-cross:
	make build-all-cross-mainnet
	make build-all-cross-testnet

build-all-mainnet:
	make build-mainnet-linux/amd64
	make build-mainnet-linux/arm64
	make build-mainnet-darwin/amd64
	make build-mainnet-darwin/arm64
	make build-mainnet-windows/amd64

build-all-cross-mainnet:
	make build-mainnet-linux/amd64
	make build-mainnet-linux/arm64
	make build-mainnet-windows/amd64

build-all-testnet:
	make build-testnet-linux/amd64
	make build-testnet-linux/arm64
	make build-testnet-darwin/amd64
	make build-testnet-darwin/arm64
	make build-testnet-windows/amd64

build-all-cross-testnet:
	make build-testnet-linux/amd64
	make build-testnet-linux/arm64
	make build-testnet-windows/amd64

build-all-local:
	make build-local-linux/amd64
	make build-local-linux/arm64
	make build-local-darwin/amd64
	make build-local-darwin/arm64
	make build-local-windows/amd64

# Linux-only build targets (for CI environments without cross-compile toolchains)
build-all-linux:
	make build-mainnet-linux/amd64
	make build-mainnet-linux/arm64
	make build-testnet-linux/amd64
	make build-testnet-linux/arm64

do-checksum-testnet:
	cd build && sha256sum bitbadgeschain-testnet-linux-amd64 bitbadgeschain-testnet-linux-arm64 bitbadgeschain-testnet-darwin-amd64 bitbadgeschain-testnet-darwin-arm64 bitbadgeschain-testnet-windows-amd64.exe > bitbadgeschain-testnet_checksum

do-checksum-mainnet:
	cd build && sha256sum bitbadgeschain-linux-amd64 bitbadgeschain-linux-arm64 bitbadgeschain-darwin-amd64 bitbadgeschain-darwin-arm64 bitbadgeschain-windows-amd64.exe > bitbadgeschain-mainnet_checksum

# Checksum for all binaries
do-checksum:
	cd build && sha256sum \
		bitbadgeschain-linux-amd64 bitbadgeschain-linux-arm64 \
		bitbadgeschain-darwin-amd64 bitbadgeschain-darwin-arm64 \
		bitbadgeschain-windows-amd64.exe \
		bitbadgeschain-testnet-linux-amd64 bitbadgeschain-testnet-linux-arm64 \
		bitbadgeschain-testnet-darwin-amd64 bitbadgeschain-testnet-darwin-arm64 \
		bitbadgeschain-testnet-windows-amd64.exe \
		> bitbadgeschain_checksum

# Run all tests. Requires -tags=test for EVM/cosmos-evm test helpers (ResetTestConfig).
test:
	@echo "Running all tests (with -tags=test for EVM test config)..."
	@go test ./... -count=1 -tags=test

# Run tokenization module tests only
test-tokenization:
	@go test ./x/tokenization/... -count=1 -tags=test

# Run keeper tests only
test-keeper:
	@go test ./x/tokenization/keeper/... -count=1 -tags=test

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
		echo "Using solcjs with optimizer..."; \
		cd contracts && \
		echo "Compiling PrecompileTransferTestContract..." && \
		solcjs --bin --abi --optimize --base-path . --include-path . --output-dir test \
			types/TokenizationTypes.sol \
			interfaces/ITokenizationPrecompile.sol \
			libraries/TokenizationJSONHelpers.sol \
			test/PrecompileTransferTestContract.sol && \
		echo "Compiling PrecompileCollectionTestContract..." && \
		solcjs --bin --abi --optimize --base-path . --include-path . --output-dir test \
			types/TokenizationTypes.sol \
			interfaces/ITokenizationPrecompile.sol \
			libraries/TokenizationJSONHelpers.sol \
			test/PrecompileCollectionTestContract.sol && \
		echo "Compiling PrecompileDynamicStoreTestContract..." && \
		solcjs --bin --abi --optimize --base-path . --include-path . --output-dir test \
			types/TokenizationTypes.sol \
			interfaces/ITokenizationPrecompile.sol \
			libraries/TokenizationJSONHelpers.sol \
			test/PrecompileDynamicStoreTestContract.sol && \
		echo "Compiling MinBankBalanceChecker..." && \
		solcjs --bin --abi --optimize --base-path . --include-path . --output-dir test \
			test/MinBankBalanceChecker.sol && \
		echo "Compiling MaxUniqueHoldersChecker..." && \
		solcjs --bin --abi --optimize --base-path . --include-path . --output-dir test \
			interfaces/ITokenizationPrecompile.sol \
			test/MaxUniqueHoldersChecker.sol && \
		echo "Compiling GammHelperLibrariesTestContract..." && \
		solcjs --bin --abi --optimize --base-path . --include-path . --output-dir test \
			types/GammTypes.sol \
			interfaces/IGammPrecompile.sol \
			libraries/GammWrappers.sol \
			libraries/GammBuilders.sol \
			libraries/GammHelpers.sol \
			libraries/GammJSONHelpers.sol \
			libraries/GammErrors.sol \
			test/GammHelperLibrariesTestContract.sol && \
		echo "All contracts compiled successfully"; \
		echo "Note: MinimalTestContract and HelperLibrariesTestContract use pre-compiled artifacts."; \
		echo "      If you need to recompile them, use solc with --via-ir flag."; \
	elif command -v solc >/dev/null 2>&1; then \
		echo "Using solc with --via-ir for complex contracts..."; \
		cd contracts && \
		echo "Compiling test contracts..." && \
		solc --via-ir --optimize --combined-json bin,abi --allow-paths . \
			types/TokenizationTypes.sol \
			interfaces/ITokenizationPrecompile.sol \
			libraries/TokenizationJSONHelpers.sol \
			test/PrecompileTransferTestContract.sol \
			test/PrecompileCollectionTestContract.sol \
			test/PrecompileDynamicStoreTestContract.sol > test/precompile_split_contracts.json && \
		solc --via-ir --optimize --combined-json bin,abi --allow-paths . \
			types/GammTypes.sol \
			interfaces/IGammPrecompile.sol \
			libraries/GammWrappers.sol \
			libraries/GammBuilders.sol \
			libraries/GammHelpers.sol \
			libraries/GammJSONHelpers.sol \
			libraries/GammErrors.sol \
			test/GammHelperLibrariesTestContract.sol > test/GammHelperLibrariesTestContract.json && \
		echo "All contracts compiled successfully"; \
	else \
		echo "Error: Neither solcjs nor solc found. Please install one:"; \
		echo "  npm install -g solc@0.8.24"; \
		echo "  or"; \
		echo "  Install solc via your package manager"; \
		exit 1; \
	fi

# Quint formal specification targets (approval system)
.PHONY: quint-check quint-run quint-verify quint-run-all

# Type-check all Quint specs
quint-check:
	@echo "Type-checking Quint specifications..."
	@quint typecheck specs/quint/tokenization/approval_hierarchy.qnt
	@quint typecheck specs/quint/tokenization/amount_limits.qnt
	@quint typecheck specs/quint/tokenization/replay_protection.qnt
	@echo "All specs type-check successfully"

# Run simulation on approval hierarchy (quick check)
quint-run:
	@echo "Running Quint simulation (approval hierarchy)..."
	@quint run specs/quint/tokenization/approval_hierarchy.qnt \
		--invariant=inv_all \
		--max-steps=50

# Run all approval system invariant simulations
quint-run-all:
	@echo "Running all Quint simulations..."
	@echo "=== Approval Hierarchy ===" && \
	quint run specs/quint/tokenization/approval_hierarchy.qnt \
		--invariant=inv_all --max-steps=30
	@echo "=== Amount Limits ===" && \
	quint run specs/quint/tokenization/amount_limits.qnt \
		--invariant=inv_all --max-steps=30
	@echo "=== Replay Protection ===" && \
	quint run specs/quint/tokenization/replay_protection.qnt \
		--invariant=inv_all --max-steps=30
	@echo "All simulations passed"

# Full verification (requires JDK 17+)
quint-verify:
	@echo "Verifying Quint specifications..."
	@quint verify specs/quint/tokenization/approval_hierarchy.qnt \
		--invariant=inv_all

# IBC E2E Testing targets
# Run all IBC E2E tests
test-ibc-e2e:
	@echo "Running all IBC E2E tests..."
	@go test ./testing/ibc/e2e/... -count=1 -tags=test -v

# Run IBC transfer tests only
test-ibc-transfer:
	@echo "Running IBC transfer tests..."
	@go test ./testing/ibc/e2e/... -run "Transfer" -count=1 -tags=test -v

# Run IBC hooks tests only
test-ibc-hooks:
	@echo "Running IBC hooks tests..."
	@go test ./testing/ibc/e2e/... -run "Hooks" -count=1 -tags=test -v

# Run IBC rate limit tests only
test-ibc-rate-limit:
	@echo "Running IBC rate limit tests..."
	@go test ./testing/ibc/e2e/... -run "RateLimit" -count=1 -tags=test -v

# Verify IBC TestingApp interface compliance (build check)
verify-ibc-testing-app:
	@echo "Verifying TestingApp interface compliance..."
	@go build -tags=test ./app/...
	@echo "TestingApp interface compliance verified"

.PHONY: test-ibc-e2e test-ibc-transfer test-ibc-hooks test-ibc-rate-limit verify-ibc-testing-app

# CLI Message Testing targets
# Run CLI message marshal/unmarshal tests
test-cli:
	@echo "Running CLI message tests..."
	@go test ./testing/cli/... -count=1 -tags=test -v

# Cross-Module E2E Testing targets
# Run all cross-module E2E tests
test-cross-module:
	@echo "Running all cross-module E2E tests..."
	@go test ./testing/cross_module/e2e/... -count=1 -tags=test -v

# Run pool tests only
test-cross-module-pool:
	@echo "Running cross-module pool tests..."
	@go test ./testing/cross_module/e2e/... -run "Pool" -count=1 -tags=test -v

# Run collection tests only
test-cross-module-collection:
	@echo "Running cross-module collection tests..."
	@go test ./testing/cross_module/e2e/... -run "Collection" -count=1 -tags=test -v

# Run combined cross-module tests only
test-cross-module-combined:
	@echo "Running combined cross-module tests..."
	@go test ./testing/cross_module/e2e/... -run "Combined" -count=1 -tags=test -v

# Run all E2E tests (IBC + cross-module)
test-e2e-all:
	@echo "Running all E2E tests..."
	@go test ./testing/... -count=1 -tags=test -v

.PHONY: test-cli test-cross-module test-cross-module-pool test-cross-module-collection test-cross-module-combined test-e2e-all
build-linux/amd64:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o ./build/bitbadgeschain-linux-amd64 ./cmd/bitbadgeschaind/main.go

build-linux/arm64:
	CGO_ENABLED=1 CC=aarch64-linux-gnu-gcc GOOS=linux GOARCH=arm64 go build -o ./build/bitbadgeschain-linux-arm64 ./cmd/bitbadgeschaind/main.go

build-darwin:
	CGO_ENABLED=1 CC="o64-clang" GOOS=darwin GOARCH=amd64 go build -o ./build/bitbadgeschain-darwin-amd64 ./cmd/bitbadgeschaind/main.go

build-all: 
	make build-linux/amd64
	make build-linux/arm64

do-checksum:
	cd build && sha256sum bitbadgeschain-linux-amd64 bitbadgeschain-linux-arm64 > bitbadgeschain_checksum

build-with-checksum: build-all do-checksum
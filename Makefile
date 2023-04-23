build-all:
	GOOS=linux GOARCH=amd64 go build -o ./build/bitbadgeschain-linux-amd64 ./cmd/bitbadgeschaind/main.go
	GOOS=linux GOARCH=arm64 go build -o ./build/bitbadgeschain-linux-arm64 ./cmd/bitbadgeschaind/main.go
	GOOS=darwin GOARCH=amd64 go build -o ./build/bitbadgeschain-darwin-amd64 ./cmd/bitbadgeschaind/main.go

do-checksum:
	cd build && sha256sum bitbadgeschain-linux-amd64 bitbadgeschain-linux-arm64 bitbadgeschain-darwin-amd64 > bitbadgeschain_checksum

build-with-checksum: build-all do-checksum
# Ubuntu image for rehearsing a cosmovisor upgrade between two git refs.
#
# The chain binaries are built inside the container: go-ethereum's secp256k1
# needs CGO (no pure-Go fallback for RecoverPubkey/VerifySignature), so
# cross-compiling from macOS is not an option. Building in-container also
# matches what validators actually run.
#
# GO_VERSIONS is the space-separated list of toolchains the rehearsal needs,
# computed by rehearse.sh from each ref's go.mod. Each is installed under
# /usr/local/go-<version>; container/build.sh puts the right one on PATH.
FROM ubuntu:24.04

ARG GO_VERSIONS="1.26.6"
ARG TARGETARCH

RUN apt-get update && apt-get install -y --no-install-recommends \
        build-essential ca-certificates curl git jq procps python3 \
    && rm -rf /var/lib/apt/lists/*

RUN set -eu; for v in $GO_VERSIONS; do \
      curl -fsSL "https://go.dev/dl/go${v}.linux-${TARGETARCH}.tar.gz" -o /tmp/go.tgz; \
      mkdir -p "/usr/local/go-${v}"; \
      tar -C "/usr/local/go-${v}" --strip-components=1 -xzf /tmp/go.tgz; \
      rm /tmp/go.tgz; \
    done

ENV GOMODCACHE=/gomodcache
ENV GOCACHE=/gobuild
# Modules resolve from the host's download cache, mounted read-only as a file
# proxy. Nothing is fetched from the network, which is also how private modules
# pinned in go.mod are available without shipping a token into the image.
ENV GOPROXY=file:///hostmod/cache/download
ENV GOFLAGS=-mod=mod
ENV GOSUMDB=off
ENV GONOSUMCHECK=1
ENV GOPRIVATE=*
ENV GONOSUMDB=*
# Never let Go try to fetch a toolchain: build.sh selects one explicitly.
ENV GOTOOLCHAIN=local

WORKDIR /work

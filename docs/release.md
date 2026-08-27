# Cutting a release

Normally `.github/workflows/release.yml` does everything: push a `v*` tag and it
builds the six binaries, generates `checksums.txt`, and creates the GitHub
release.

This document covers the case where it cannot — an **embargoed release**.

## When CI cannot build

If `go.mod` carries a `replace` pointing `github.com/cosmos/evm` at a private
hotfix repository, the Actions runner has no access to that module and every
build job fails on `go mod download`. v33 and v34 were both in this position.

`release.yml`'s `preflight` job detects that replace and skips the build and
release jobs, so the tag produces a clean skip instead of a wall of red — and,
more importantly, the release job cannot run with partial artifacts and
overwrite hand-uploaded assets.

Binaries for those releases are built locally and uploaded by hand.

## Building release binaries locally

Match the CI environment rather than building on your workstation: the workflow
pins `ubuntu-22.04` (glibc 2.35) so the binaries run on Ubuntu 22.04+, Debian
12+ and RHEL 9+. Building on a newer glibc silently produces binaries that fail
on older hosts, and nothing catches it until an operator reports a crash.

`scripts-dev/upgrade-rehearsal/` carries a Dockerfile for the test harness; the
release build needs its own, with the cross toolchains the workflow installs:

```dockerfile
FROM ubuntu:22.04
RUN apt-get update && apt-get install -y \
        build-essential ca-certificates curl git make \
        gcc-aarch64-linux-gnu libc6-dev-arm64-cross gcc-mingw-w64-x86-64
# + Go matching go.mod
```

Two things that will cost you an hour if you miss them:

- **Do not pass `--no-install-recommends`** to that apt line. The workflow does
  not, and `gcc-aarch64-linux-gnu` pulls `libc6-dev-arm64-cross` as a
  *recommend*. Without it, cgo fails on `bits/wordsize.h: No such file`.
- **Run the container as `linux/amd64`.** `make build-mainnet-linux/amd64` sets
  no `CC`, so it assumes a native amd64 host, exactly like the runner. On an
  arm64 machine use `docker run --platform linux/amd64` (emulated).

Then, with the embargoed module already in your local module cache:

```sh
make build-all-cross VERSION=v34
cd build && sha256sum \
  bitbadgeschain-linux-amd64 bitbadgeschain-linux-arm64 \
  bitbadgeschain-windows-amd64.exe \
  bitbadgeschain-testnet-linux-amd64 bitbadgeschain-testnet-linux-arm64 \
  bitbadgeschain-testnet-windows-amd64.exe > checksums.txt
```

`build-all-cross` deliberately omits macOS — it needs an SDK the Linux image
does not have. Upgrade proposals reference only the Linux binaries.

To avoid putting a token in the image, mount the host module cache read-only and
point Go at it as a file proxy:

```
-v ~/go/pkg/mod:/hostmod:ro  -e GOPROXY=file:///hostmod/cache/download
```

Leave `GOPRIVATE` and `GONOPROXY` **empty** when you do — they imply "bypass the
proxy" and send Go to the network, where the private module then fails to
authenticate.

## Verify before uploading

Do not upload a binary you have only compiled. At minimum:

```sh
./bitbadgeschain-linux-amd64 version          # prints the tag
go version -m ./bitbadgeschain-linux-amd64 | grep cosmos/evm   # shows the replace
```

The second one is the check that matters for an embargoed release: it proves the
hotfix module is actually linked into the artifact operators will run, rather
than the public module the tag *claims* to replace.

Then start a single-validator chain from the real artifact and confirm it
produces blocks. `scripts-dev/upgrade-rehearsal/` has the harness.

## Publishing

Create the release as a **draft** so the tag is not created until you publish:

```sh
gh release create v34 --draft --title v34 --notes-file notes.md build/*
```

Check the upgrade height and time in the notes against the governance proposal
before publishing — those are the two fields operators act on, and they are the
easiest to leave stale from the previous release.

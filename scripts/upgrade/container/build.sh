#!/usr/bin/env bash
# Builds the FROM and TO chain binaries plus cosmovisor inside the rehearsal
# image. Runs as: build.sh <label> <src dir> <go version> <out binary> [commit]
#
# Sources are mounted read-only, so the tree is copied to a writable scratch
# dir first: `-mod=mod` may touch go.sum.
set -euo pipefail

LABEL=$1 SRC=$2 GOV=$3 OUT=$4 COMMIT=${5:-unknown}

export PATH="/usr/local/go-${GOV}/bin:$PATH"
[ -x "/usr/local/go-${GOV}/bin/go" ] || { echo "go ${GOV} is not installed in the image" >&2; exit 1; }

WORK=/work/$LABEL
rm -rf "$WORK"; mkdir -p "$WORK"
cp -a "$SRC/." "$WORK/"
cd "$WORK"

LDFLAGS="-X github.com/cosmos/cosmos-sdk/version.Name=bitbadgeschain \
 -X github.com/cosmos/cosmos-sdk/version.AppName=bitbadgeschaind \
 -X github.com/cosmos/cosmos-sdk/version.Version=$LABEL \
 -X github.com/cosmos/cosmos-sdk/version.Commit=$COMMIT"

echo "building $LABEL with $(go version) -> $OUT"
if [ -f cmd/bitbadgeschaind/main.go ]; then
  CGO_ENABLED=1 go build -trimpath -ldflags "$LDFLAGS" -o "$OUT" ./cmd/bitbadgeschaind
else
  # cosmovisor wrapper module: build the pinned tool.
  go build -trimpath -o "$OUT" cosmossdk.io/tools/cosmovisor/cmd/cosmovisor
fi
chmod +x "$OUT"
"$OUT" version 2>/dev/null | head -1 || true

# The evmtx helper is built inside the TO tree so go-ethereum resolves from
# the chain's own go.mod. EVMTX_SRC points at scripts/upgrade/evmtx.
if [ -n "${EVMTX_SRC:-}" ] && [ -f cmd/bitbadgeschaind/main.go ]; then
  rm -rf scripts/upgrade/evmtx; mkdir -p scripts/upgrade/evmtx
  cp "$EVMTX_SRC"/*.go scripts/upgrade/evmtx/
  echo "building evmtx -> /out/evmtx"
  CGO_ENABLED=1 go build -trimpath -o /out/evmtx ./scripts/upgrade/evmtx
fi

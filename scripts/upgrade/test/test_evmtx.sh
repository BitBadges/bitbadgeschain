#!/usr/bin/env bash
# evmtx: vets, builds, is gofmt-clean, and refuses to run without its flags.
set -uo pipefail
HERE=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
REPO=$(cd "$HERE/../../.." && pwd)
FAILED=0
ok()   { echo "  PASS  $1"; }
fail() { echo "  FAIL  $1"; FAILED=$((FAILED+1)); }
assert() { local l=$1; shift; if "$@"; then ok "$l"; else fail "$l"; fi; }
refute() { local l=$1; shift; if "$@" >/dev/null 2>&1; then fail "$l"; else ok "$l"; fi; }

if ! command -v go >/dev/null; then echo "  SKIP  go not installed"; exit 0; fi
T=$(mktemp -d); trap 'rm -rf "$T"' EXIT
cd "$REPO" || exit 1
assert "gofmt-clean"  [ -z "$(gofmt -l scripts/upgrade/evmtx)" ]
assert "go vet"       go vet ./scripts/upgrade/evmtx
assert "go build"     go build -o "$T/evmtx" ./scripts/upgrade/evmtx
if [ -x "$T/evmtx" ]; then
  refute "refuses to run without --key/--to/--value-wei" "$T/evmtx" --rpc http://127.0.0.1:1
  refute "rejects a non-numeric --value-wei" "$T/evmtx" --key 00 --to 0x0 --value-wei abc
  refute "fails cleanly when the RPC is unreachable" "$T/evmtx" --rpc http://127.0.0.1:1 \
    --key 4c0883a69102937d6231471b5dbb6204fe5129617082792ae468d01a3f362318 \
    --to 0x0000000000000000000000000000000000000001 --value-wei 1
fi

echo
if [ "$FAILED" = 0 ]; then echo "evmtx: all passed"; else echo "evmtx: $FAILED failed"; exit 1; fi

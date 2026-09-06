#!/usr/bin/env bash
set -euo pipefail

verify() {
  local binary=$1 expected=$2 actual
  actual=$("$binary" version --long --output json | jq -er '.commit | select(type == "string" and length > 0)') || {
    echo "cannot read commit identity from $binary; rebuild without --skip-build" >&2
    return 1
  }
  if [ "$actual" != "$expected" ]; then
    echo "$binary reports commit $actual, expected $expected; rebuild without --skip-build" >&2
    return 1
  fi
}

verify "$1" "$2"
verify "$3" "$4"

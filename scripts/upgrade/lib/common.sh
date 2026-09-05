#!/usr/bin/env bash
# Shared helpers for the upgrade rehearsal scripts. Source, do not execute.

red()  { printf '\033[0;31m%s\033[0m\n' "$*"; }
grn()  { printf '\033[0;32m%s\033[0m\n' "$*"; }
ylw()  { printf '\033[1;33m%s\033[0m\n' "$*"; }
step() { printf '\n\033[0;34m=== %s ===\033[0m\n' "$*"; }
die()  { red "ERROR: $*" >&2; exit 1; }

# check <label> <command...>
# Runs the command; prints PASS/FAIL and counts failures in $FAILED.
FAILED=0
check() {
  local label=$1; shift
  if "$@"; then
    grn "  PASS  $label"
  else
    red "  FAIL  $label"
    FAILED=$((FAILED + 1))
  fi
}

# toml_get <file> <section> <key> -> prints the raw value (quotes included)
toml_get() {
  sed -n "/^\[$2\]/,/^\[/p" "$1" | grep -E "^\s*$3\s*=" | head -1 | sed 's/^[^=]*= *//'
}

# genesis_json_set <genesis> <jq filter>
genesis_json_set() {
  jq "$2" "$1" > "$1.tmp" && mv "$1.tmp" "$1"
}

# wait_for_height <bin> <min height> [timeout seconds] [extra args...]
# Polls `<bin> status` until latest_block_height >= min. Sets $H.
wait_for_height() {
  local bin=$1 min=$2 timeout=${3:-120}; shift 3 || shift $#
  local deadline=$((SECONDS + timeout))
  H=""
  while [ $SECONDS -lt $deadline ]; do
    H=$("$bin" status "$@" 2>/dev/null | jq -r '.sync_info.latest_block_height // empty' 2>/dev/null || true)
    if [ -n "$H" ] && [ "$H" -ge "$min" ] 2>/dev/null; then return 0; fi
    sleep 2
  done
  return 1
}

# tx_hash_of <cli output> -> first 64-hex token
tx_hash_of() { grep -oE '[A-F0-9]{64}' <<<"$1" | head -1; }

# go_version_of_gomod <path or - for stdin> -> "1.26.6"
go_version_of_gomod() {
  local content tc go
  content=$(cat "$1")
  tc=$(grep -E '^toolchain go' <<<"$content" | head -1 | sed 's/^toolchain go//')
  go=$(grep -E '^go [0-9]' <<<"$content" | head -1 | sed 's/^go //')
  echo "${tc:-$go}"
}

summary() {
  echo
  if [ "$FAILED" -eq 0 ]; then grn "$1"; else red "$FAILED check(s) failed"; fi
}

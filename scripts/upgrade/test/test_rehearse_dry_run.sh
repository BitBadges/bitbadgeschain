#!/usr/bin/env bash
# rehearse.sh --dry-run resolves refs, toolchains and the upgrade name without Docker.
set -uo pipefail
HERE=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
RH="$HERE/../rehearse.sh"
FAILED=0
ok()   { echo "  PASS  $1"; }
fail() { echo "  FAIL  $1"; FAILED=$((FAILED+1)); }
# assert <label> <cmd...>: PASS when the command succeeds; refute: when it fails.
assert() { local l=$1; shift; if "$@"; then ok "$l"; else fail "$l"; fi; }
refute() { local l=$1; shift; if "$@" >/dev/null 2>&1; then fail "$l"; else ok "$l"; fi; }
has()  { if grep -qE -- "$2" <<<"$3"; then ok "$1"; else fail "$1"; fi; }

set +e; OUT=$(PATH=/usr/bin:/bin "$RH" --from v34 --to HEAD --evmcheck --rollback --dry-run 2>&1); RC=$?; set -e
assert "exits 0 without docker on PATH" [ "$RC" = 0 ]
has "derives upgrade name from app/upgrades on --to"  'upgrade +: v35' "$OUT"
has "reads the FROM toolchain from its go.mod"       'from +: v34 \([0-9a-f]{40}\) go1\.[0-9]+' "$OUT"
has "reads the TO toolchain from its go.mod"         'to +: HEAD \([0-9a-f]{40}\) go1\.[0-9]+' "$OUT"
has "detects the host arch"                          'arch +: (arm64|amd64)' "$OUT"
has "lists the requested stages"                     'stages +: build upgrade rollback evmcheck' "$OUT"
has "finds the v35 checks hook"                      'checks hook : checks/v35\.sh' "$OUT"
has "reports the dry run"                            'dry run: nothing executed' "$OUT"

set +e; OUT2=$("$RH" --from v34 --to HEAD --name v99 --dry-run 2>&1); RC2=$?; set -e
assert "--name without app/upgrades/<name> is rejected" [ "$RC2" != 0 ]
assert "...with a clear message" grep -q 'no app/upgrades/v99' <<<"$OUT2"
refute "unknown git ref is rejected" "$RH" --from no-such-ref --to HEAD --dry-run

echo
if [ "$FAILED" = 0 ]; then echo "rehearse dry-run: all passed"; else echo "rehearse dry-run: $FAILED failed"; exit 1; fi

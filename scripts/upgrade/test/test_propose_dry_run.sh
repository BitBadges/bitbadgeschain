#!/usr/bin/env bash
# propose.sh --dry-run renders the proposal without a node; deposit defaults to ubadge.
set -uo pipefail
HERE=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
PR="$HERE/../propose.sh"
FAILED=0
ok()   { echo "  PASS  $1"; }
fail() { echo "  FAIL  $1"; FAILED=$((FAILED+1)); }
# assert <label> <cmd...>: PASS when the command succeeds; refute: when it fails.
assert() { local l=$1; shift; if "$@"; then ok "$l"; else fail "$l"; fi; }
refute() { local l=$1; shift; if "$@" >/dev/null 2>&1; then fail "$l"; else ok "$l"; fi; }
has()  { if grep -qF -- "$2" <<<"$3"; then ok "$1"; else fail "$1"; fi; }
hasnt(){ if grep -qF -- "$2" <<<"$3"; then fail "$1"; else ok "$1"; fi; }

OUT=$("$PR" --name v35 --home /tmp/x --from faucet --dry-run 2>&1); RC=$?
assert "dry run exits 0" [ "$RC" = 0 ]
has   "plan name"                    '"name": "v35"' "$OUT"
has   "default deposit is ubadge"    '"deposit": "10000000ubadge"' "$OUT"
hasnt "never ustake"                 'ustake' "$OUT"
has   "relative height is computed from the node" '"height": "<current height> + 30"' "$OUT"
has   "not expedited by default"     '"expedited": false' "$OUT"
has   "submits with --gas auto"      'tx gov submit-proposal <proposal> --from faucet --keyring-backend os --gas auto' "$OUT"
has   "votes yes by default"         'tx gov vote <id> yes' "$OUT"
has   "reports the height"           'UPGRADE_HEIGHT=<current height> + 30' "$OUT"
has   "default gas price meets v35 floor" '--gas-prices 10ubadge' "$OUT"

OUT_FEES=$("$PR" --name v35 --home /tmp/x --from k --fees 1000000ubadge --dry-run 2>&1)
has   "explicit fees remain supported" '--fees 1000000ubadge' "$OUT_FEES"
hasnt "explicit fees replace default gas price" '--gas-prices' "$OUT_FEES"

OUT2=$("$PR" --name v35 --home /tmp/x --from k --height 12345 --expedited --deposit 20000000000ubadge --no-vote --gas-prices 10ubadge --chain-id bitbadges-1 --dry-run 2>&1)
has   "absolute height"              '"height": "12345"' "$OUT2"
has   "expedited flag"               '"expedited": true' "$OUT2"
has   "custom deposit"               '"deposit": "20000000000ubadge"' "$OUT2"
has   "gas prices instead of fees"   '--gas-prices 10ubadge' "$OUT2"
has   "explicit chain id"            '--chain-id bitbadges-1' "$OUT2"
hasnt "--no-vote skips the vote"     'tx gov vote' "$OUT2"

refute "ustake deposit is rejected" "$PR" --name v35 --home /tmp/x --from k --deposit 1000ustake --dry-run
refute "missing --home is rejected" "$PR" --name v35 --from k --dry-run

echo
if [ "$FAILED" = 0 ]; then echo "propose dry-run: all passed"; else echo "propose dry-run: $FAILED failed"; exit 1; fi

#!/usr/bin/env bash
# Single-validator cosmovisor upgrade rehearsal: FROM binary -> TO binary.
#
# The question this answers: does a validator have to do anything by hand?
# After cosmovisor starts, nothing here touches the config, restarts a node,
# or intervenes. The only inputs are the governance proposal and the vote,
# which is what a real upgrade looks like from the operator's side.
#
# Environment (set by rehearse.sh):
#   FROM_BIN TO_BIN CV   binaries
#   UPGRADE_NAME         plan name, e.g. v35
#   CHECKS_DIR           directory holding optional checks/<name>.sh
#   LOG_DIR              where logs go
set -euo pipefail
# shellcheck source=../lib/common.sh
. "$(dirname "$0")/../lib/common.sh"

DH=${DAEMON_HOME_DIR:-/root/.bitbadgeschain}
CHAIN_ID=bitbadges-rehearsal-1
KR=(--keyring-backend test)
LOG=$LOG_DIR/cosmovisor.log
PROPOSE="$(dirname "$0")/../propose.sh"

step "1. Initialise a chain with the FROM binary"
rm -rf "$DH"
"$FROM_BIN" init testnode --chain-id "$CHAIN_ID" --home "$DH" >/dev/null 2>&1
"$FROM_BIN" keys add validator "${KR[@]}" --home "$DH" >/dev/null 2>&1
"$FROM_BIN" keys add alice "${KR[@]}" --home "$DH" >/dev/null 2>&1

G="$DH/config/genesis.json"
BOND_DENOM=$(jq -r '.app_state.staking.params.bond_denom' "$G")
DEPOSIT_DENOM=$(jq -r '.app_state.gov.params.min_deposit[0].denom' "$G")
FUNDS="100000000000000${BOND_DENOM},100000000000000ubadge"
[ "$DEPOSIT_DENOM" = "$BOND_DENOM" ] || [ "$DEPOSIT_DENOM" = ubadge ] || FUNDS="$FUNDS,100000000000000${DEPOSIT_DENOM}"
"$FROM_BIN" genesis add-genesis-account validator "$FUNDS" "${KR[@]}" --home "$DH" >/dev/null 2>&1
"$FROM_BIN" genesis add-genesis-account alice "$FUNDS" "${KR[@]}" --home "$DH" >/dev/null 2>&1

# Short governance periods are the only genesis deviation from a stock chain.
genesis_json_set "$G" '.app_state.gov.params.voting_period="15s"
  | .app_state.gov.params.expedited_voting_period="10s"
  | .app_state.gov.params.max_deposit_period="30s"
  | .app_state.gov.params.min_deposit[0].amount="1"
  | .app_state.gov.params.expedited_min_deposit[0].amount="2"'

"$FROM_BIN" genesis gentx validator "50000000000${BOND_DENOM}" --chain-id "$CHAIN_ID" "${KR[@]}" --home "$DH" >/dev/null 2>&1
"$FROM_BIN" genesis collect-gentxs --home "$DH" >/dev/null 2>&1
"$FROM_BIN" genesis validate-genesis --home "$DH" >/dev/null 2>&1 || die "genesis did not validate"
grn "chain initialised (bond denom $BOND_DENOM, deposit denom $DEPOSIT_DENOM)"

# The REST API server is enabled so checks/<name>.sh can read module params
# that have no CLI query (cosmos/evm's feemarket, for one). This is the only
# app.toml edit, and it happens before the snapshot below.
python3 - "$DH/config/app.toml" <<'PYEOF'
import re,sys
p=sys.argv[1]; s=open(p).read()
m=re.search(r'(?ms)^\[api\].*?(?=^\[|\Z)', s)
blk=re.sub(r'(?m)^enable\s*=.*$', 'enable = true', m.group(0))
open(p,'w').write(s[:m.start()]+blk+s[m.end():])
PYEOF
API_URL=http://127.0.0.1:1317

# Snapshot the config files. Nothing in this script writes to them after this
# point, so any change at the end was made by the TO binary's pre-upgrade hook.
mkdir -p "$LOG_DIR/config-before"
cp "$DH/config/config.toml" "$DH/config/app.toml" "$LOG_DIR/config-before/"

step "2. Lay out cosmovisor"
mkdir -p "$DH/cosmovisor/genesis/bin" "$DH/cosmovisor/upgrades/$UPGRADE_NAME/bin"
cp "$FROM_BIN" "$DH/cosmovisor/genesis/bin/bitbadgeschaind"
cp "$TO_BIN" "$DH/cosmovisor/upgrades/$UPGRADE_NAME/bin/bitbadgeschaind"
export DAEMON_NAME=bitbadgeschaind DAEMON_HOME=$DH DAEMON_ALLOW_DOWNLOAD_BINARIES=false \
       DAEMON_RESTART_AFTER_UPGRADE=true UNSAFE_SKIP_BACKUP=true
grn "cosmovisor laid out: genesis=FROM, upgrades/$UPGRADE_NAME=TO"

step "3. Start cosmovisor (FROM) - nothing after this point is hands-on"
"$CV" run start --home "$DH" --minimum-gas-prices "0${BOND_DENOM}" > "$LOG" 2>&1 &
CV_PID=$!
wait_for_height "$FROM_BIN" 2 120 --home "$DH" || { tail -30 "$LOG"; die "node never produced a block"; }
grn "FROM node is producing blocks (height $H)"

step "4. Submit and pass the $UPGRADE_NAME upgrade proposal"
PROPOSE_OUT=$("$PROPOSE" --name "$UPGRADE_NAME" --home "$DH" --from validator \
  --bin "$FROM_BIN" --chain-id "$CHAIN_ID" --keyring-backend test \
  --deposit "10${DEPOSIT_DENOM}" --fees "0${BOND_DENOM}" --height +25)
echo "$PROPOSE_OUT"
UPGRADE_HEIGHT=$(sed -n 's/^UPGRADE_HEIGHT=//p' <<<"$PROPOSE_OUT")
PID=$(sed -n 's/^PROPOSAL_ID=//p' <<<"$PROPOSE_OUT")
[ -n "$UPGRADE_HEIGHT" ] && [ -n "$PID" ] || die "propose.sh did not report a proposal"

sleep 20
STATUS=$("$FROM_BIN" query gov proposal "$PID" --home "$DH" --output json 2>/dev/null | jq -r '.proposal.status')
[ "$STATUS" = "PROPOSAL_STATUS_PASSED" ] || die "proposal $PID did not pass (status $STATUS)"
grn "proposal $PID passed, upgrade height $UPGRADE_HEIGHT"

step "5. Wait for the upgrade to apply - no intervention"
DEADLINE=$((SECONDS + 420))
UPGRADED=0
while [ $SECONDS -lt $DEADLINE ]; do
  CH=$("$TO_BIN" status --home "$DH" 2>/dev/null | jq -r '.sync_info.latest_block_height // empty' 2>/dev/null || true)
  if [ -n "${CH:-}" ] && [ "$CH" -gt $((UPGRADE_HEIGHT + 5)) ] 2>/dev/null; then UPGRADED=1; break; fi
  kill -0 "$CV_PID" 2>/dev/null || { red "cosmovisor exited"; break; }
  sleep 5
done
FINAL_HEIGHT=$("$TO_BIN" status --home "$DH" 2>/dev/null | jq -r '.sync_info.latest_block_height // "0"' 2>/dev/null || echo 0)

step "6. Results"
echo "upgrade height : $UPGRADE_HEIGHT"
echo "final height   : $FINAL_HEIGHT"
echo
check "chain advanced past the upgrade height"          [ "$UPGRADED" = 1 ]
check "FROM binary halted with UPGRADE \"$UPGRADE_NAME\" NEEDED" grep -q "UPGRADE \"$UPGRADE_NAME\" NEEDED" "$LOG"
check "cosmovisor ran the pre-upgrade hook"             grep -q 'pre-upgrade' "$LOG"
check "TO binary applied the $UPGRADE_NAME handler"     grep -q "applying upgrade \"$UPGRADE_NAME\"" "$LOG"
CURRENT_BIN=$(readlink -f "$DH/cosmovisor/current/bin/bitbadgeschaind" 2>/dev/null || echo "")
check "cosmovisor 'current' points at upgrades/$UPGRADE_NAME" grep -q "/upgrades/$UPGRADE_NAME/" <<<"$CURRENT_BIN"
check "node is still producing blocks after the upgrade" [ "$FINAL_HEIGHT" -gt "$UPGRADE_HEIGHT" ]

# The pre-upgrade hook may rewrite config files (v34 did, for mempool.type).
# Report whatever changed; the hook must leave a restorable backup when it does.
for f in config.toml app.toml; do
  if diff -q "$LOG_DIR/config-before/$f" "$DH/config/$f" >/dev/null; then
    ylw "  $f: unchanged by the upgrade"
  else
    ylw "  $f: rewritten by the pre-upgrade hook:"
    diff "$LOG_DIR/config-before/$f" "$DH/config/$f" | sed 's/^/      /' || true
    check "$f.bak restores the pre-upgrade $f" diff -q "$LOG_DIR/config-before/$f" "$DH/config/$f.bak"
  fi
done

step "7. Post-upgrade assertions for $UPGRADE_NAME"
HOOK="$CHECKS_DIR/$UPGRADE_NAME.sh"
if [ -f "$HOOK" ] && [ "$UPGRADED" = 1 ]; then
  set +e
  BIN="$TO_BIN" HOME_DIR="$DH" UPGRADE_HEIGHT="$UPGRADE_HEIGHT" LOG="$LOG" API_URL="$API_URL" bash "$HOOK"
  HOOK_FAILED=$?
  set -e
  FAILED=$((FAILED + HOOK_FAILED))
elif [ -f "$HOOK" ]; then
  red "  FAIL  skipped $HOOK: the upgrade did not apply"; FAILED=$((FAILED + 1))
else
  ylw "  no checks/$UPGRADE_NAME.sh - only the generic checks ran"
fi

kill "$CV_PID" 2>/dev/null || true
wait "$CV_PID" 2>/dev/null || true
summary "ALL CHECKS PASSED - the $UPGRADE_NAME upgrade was hands-off"
exit "$FAILED"

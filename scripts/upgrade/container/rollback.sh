#!/usr/bin/env bash
# Rollback check: can an operator go back to the FROM binary after the TO
# binary's pre-upgrade hook has run?
#
# Asserts that the hook exits 0 and is idempotent; that when it rewrites a
# config file it leaves a <file>.bak matching the pristine config; and that
# the FROM binary produces blocks again once the documented rollback order
# (restore config, then binary) is followed. Whether the FROM binary also
# starts against the hook-rewritten config is reported, not asserted: it is
# what tells the release notes whether the order matters for this upgrade.
set -euo pipefail
# shellcheck source=../lib/common.sh
. "$(dirname "$0")/../lib/common.sh"

DH=/root/.bitbadgeschain-rb
CHAIN_ID=bitbadges-rollback-1
KR=(--keyring-backend test)

produces_blocks() { grep -qiE 'finalizing commit|executed block|indexed block' "$1"; }
start_from_for() { # <seconds> <log>
  set +e; timeout "$1" "$FROM_BIN" start --home "$DH" > "$2" 2>&1; set -e
}

step "1. Stock FROM config (real single-validator genesis)"
rm -rf "$DH"
"$FROM_BIN" init rbnode --chain-id "$CHAIN_ID" --home "$DH" >/dev/null 2>&1
"$FROM_BIN" keys add v "${KR[@]}" --home "$DH" >/dev/null 2>&1
BOND_DENOM=$(jq -r '.app_state.staking.params.bond_denom' "$DH/config/genesis.json")
"$FROM_BIN" genesis add-genesis-account v "100000000000000${BOND_DENOM}" "${KR[@]}" --home "$DH" >/dev/null 2>&1
"$FROM_BIN" genesis gentx v "50000000000${BOND_DENOM}" --chain-id "$CHAIN_ID" "${KR[@]}" --home "$DH" >/dev/null 2>&1
"$FROM_BIN" genesis collect-gentxs --home "$DH" >/dev/null 2>&1
mkdir -p /tmp/pristine
cp "$DH/config/config.toml" "$DH/config/app.toml" /tmp/pristine/

# Control: this node DOES start before anything is changed, so a later
# failure can be attributed to the hook rather than to the fixture.
start_from_for 25 "$LOG_DIR/rb-control.log"
check "control: stock FROM node produces blocks before any edit" produces_blocks "$LOG_DIR/rb-control.log"

step "2. Run the TO pre-upgrade hook exactly as cosmovisor does"
# No --home argument; home comes only from DAEMON_HOME.
set +e; DAEMON_HOME=$DH "$TO_BIN" pre-upgrade > "$LOG_DIR/rb-hook1.log" 2>&1; RC1=$?; set -e
cat "$LOG_DIR/rb-hook1.log"
check "hook exits 0" [ "$RC1" = 0 ]
cp "$DH/config/config.toml" /tmp/after1.toml
set +e; DAEMON_HOME=$DH "$TO_BIN" pre-upgrade > "$LOG_DIR/rb-hook2.log" 2>&1; RC2=$?; set -e
check "hook is idempotent (second run exits 0)" [ "$RC2" = 0 ]
check "second run leaves config.toml unchanged" diff -q /tmp/after1.toml "$DH/config/config.toml"

CHANGED=()
for f in config.toml app.toml; do
  diff -q "/tmp/pristine/$f" "$DH/config/$f" >/dev/null || CHANGED+=("$f")
done

if [ ${#CHANGED[@]} -eq 0 ]; then
  step "3. Hook left both config files untouched"
  ylw "  rollback is binary-only for this upgrade: no config step in the release notes"
  start_from_for 25 "$LOG_DIR/rb-after-hook.log"
  check "FROM binary still produces blocks after the hook" produces_blocks "$LOG_DIR/rb-after-hook.log"
else
  step "3. Hook rewrote: ${CHANGED[*]}"
  for f in "${CHANGED[@]}"; do
    diff "/tmp/pristine/$f" "$DH/config/$f" | sed 's/^/      /' || true
    check "$f.bak exists" [ -f "$DH/config/$f.bak" ]
    check "$f.bak matches the pristine config" diff -q "/tmp/pristine/$f" "$DH/config/$f.bak"
  done

  step "4. Binary-first rollback (config still rewritten) - informational"
  start_from_for 25 "$LOG_DIR/rb-wrong-order.log"
  if produces_blocks "$LOG_DIR/rb-wrong-order.log"; then
    ylw "  FROM binary accepts the rewritten config: rollback order does not matter"
  else
    ylw "  FROM binary does NOT start against the rewritten config - release notes must say: restore config first, then binary"
    grep -iE 'panic|error' "$LOG_DIR/rb-wrong-order.log" | head -3 | sed 's/^/      /' || true
  fi

  step "5. Documented order: restore config, then binary"
  for f in "${CHANGED[@]}"; do cp "$DH/config/$f.bak" "$DH/config/$f"; done
  start_from_for 25 "$LOG_DIR/rb-right-order.log"
  check "FROM binary produces blocks against the restored config" produces_blocks "$LOG_DIR/rb-right-order.log"
fi

summary "ROLLBACK PROCEDURE VERIFIED"
exit "$FAILED"

#!/usr/bin/env bash
# Fresh-chain EVM balance check on the TO binary.
#
# Starts a chain from a plain `init` with JSON-RPC enabled and compares
# eth_getBalance with the bank balance (ubadge x 10^9 on a 9-decimal chain).
# The upgrade rehearsal cannot see this: it starts from a config the FROM
# binary wrote and the hook repaired, so it says nothing about the config
# `init` produces or about the EVM bridge on a running node.
set -euo pipefail
# shellcheck source=../lib/common.sh
. "$(dirname "$0")/../lib/common.sh"

BIN=$TO_BIN
DH=/root/.bitbadgeschain-evm
CHAIN_ID=bitbadges-evm-1
KR=(--keyring-backend test)
RPC=http://127.0.0.1:8545

step "1. Single-validator chain with JSON-RPC enabled (TO binary)"
rm -rf "$DH"
"$BIN" init evmnode --chain-id "$CHAIN_ID" --home "$DH" >/dev/null 2>&1
"$BIN" keys add v "${KR[@]}" --home "$DH" >/dev/null 2>&1
"$BIN" keys add user "${KR[@]}" --home "$DH" --key-type eth_secp256k1 >/dev/null 2>&1 || \
  "$BIN" keys add user "${KR[@]}" --home "$DH" >/dev/null 2>&1
BOND_DENOM=$(jq -r '.app_state.staking.params.bond_denom' "$DH/config/genesis.json")
FUND=100000000000000
"$BIN" genesis add-genesis-account v "${FUND}${BOND_DENOM},${FUND}ubadge" "${KR[@]}" --home "$DH" >/dev/null 2>&1
"$BIN" genesis add-genesis-account user "${FUND}${BOND_DENOM},${FUND}ubadge" "${KR[@]}" --home "$DH" >/dev/null 2>&1
"$BIN" genesis gentx v "50000000000${BOND_DENOM}" --chain-id "$CHAIN_ID" "${KR[@]}" --home "$DH" >/dev/null 2>&1
"$BIN" genesis collect-gentxs --home "$DH" >/dev/null 2>&1

enable_rpc_and_api "$DH/config/app.toml"

step "2. Start the node"
"$BIN" start --home "$DH" --minimum-gas-prices "0${BOND_DENOM}" > "$LOG_DIR/evmnode.log" 2>&1 &
PID=$!
wait_for_height "$BIN" 2 120 --home "$DH" || { tail -20 "$LOG_DIR/evmnode.log"; die "node never started from a fresh init"; }
grn "node at height $H"

step "3. Compare bank balance with eth_getBalance"
ADDR=$("$BIN" keys show user -a "${KR[@]}" --home "$DH")
BANK=$("$BIN" query bank balances "$ADDR" --home "$DH" --output json | jq -r '.balances[] | select(.denom=="ubadge") | .amount')
ylw "bank ubadge : $BANK"
HEX=$("$BIN" debug addr "$ADDR" 2>&1 | grep -i 'hex' | head -1 | awk '{print $NF}')
[ -n "$HEX" ] || { kill "$PID"; die "could not derive hex address"; }
case "$HEX" in 0x*) ;; *) HEX=0x$HEX ;; esac
ylw "evm address : $HEX"
sleep 3
RESP=$(curl -s --max-time 15 -X POST -H 'Content-Type: application/json' \
  --data "{\"jsonrpc\":\"2.0\",\"method\":\"eth_getBalance\",\"params\":[\"$HEX\",\"latest\"],\"id\":1}" "$RPC" || true)
ylw "eth_getBalance response: $RESP"
WEI_HEX=$(jq -r '.result // empty' <<<"$RESP")
[ -n "$WEI_HEX" ] || { tail -20 "$LOG_DIR/evmnode.log"; kill "$PID" 2>/dev/null || true; die "JSON-RPC returned no result"; }
WEI=$(python3 -c "print(int('$WEI_HEX',16))")
EXPECTED=$(python3 -c "print($BANK * 10**9)")
ylw "eth_getBalance : $WEI"
ylw "expected       : $EXPECTED  (bank ubadge x 10^9)"
check "eth_getBalance is non-zero (precisebank is wired)" [ "$WEI" != 0 ]
check "eth_getBalance equals bank ubadge x 10^9"          [ "$WEI" = "$EXPECTED" ]

step "4. A transaction still works"
RECV=$("$BIN" keys show v -a "${KR[@]}" --home "$DH")
SEND_OUT=$("$BIN" tx bank send "$ADDR" "$RECV" 1000ubadge --from user "${KR[@]}" --home "$DH" \
  --chain-id "$CHAIN_ID" --gas auto --gas-adjustment 1.5 --fees "0${BOND_DENOM}" -y 2>&1 || true)
sleep 6
TXH=$(tx_hash_of "$SEND_OUT")
if [ -n "$TXH" ]; then
  CODE=$("$BIN" query tx "$TXH" --home "$DH" --output json 2>/dev/null | jq -r '.code // "?"')
  check "bank send succeeded (code 0)" [ "$CODE" = 0 ]
else
  red "  FAIL  could not parse tx hash"; head -5 <<<"$SEND_OUT"; FAILED=$((FAILED+1))
fi

step "5. EVM value transfers over JSON-RPC"
# shellcheck source=../lib/evm_transfer.sh
. "$(dirname "$0")/../lib/evm_transfer.sh"
HOME_DIR=$DH RPC_URL=$RPC API_URL=http://127.0.0.1:1317 EVMTX=${EVMTX:-/out/evmtx} FUNDER=v evm_transfer_checks

kill "$PID" 2>/dev/null || true
wait "$PID" 2>/dev/null || true
summary "EVM BRIDGE VERIFIED ON A FRESH $UPGRADE_NAME NODE"
exit "$FAILED"

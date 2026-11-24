#!/usr/bin/env bash
# Post-upgrade assertions for v37 (9 -> 18 decimal migration: ubadge -> abadge).
# Run by container/upgrade.sh after the upgrade applied.
set -uo pipefail
# shellcheck source=../lib/common.sh
. "$(dirname "$0")/../lib/common.sh"

KRT=(--keyring-backend test)

# 1. Verify staking params migrated (bond_denom & key_rotation_fee -> abadge).
STAKING_JSON=$("$BIN" query staking params --home "$HOME_DIR" --output json 2>/dev/null || true)
BOND_DENOM=$(jq -r '.params.bond_denom // empty' <<<"$STAKING_JSON")
ROT_DENOM=$(jq -r '.params.key_rotation_fee.denom // empty' <<<"$STAKING_JSON")
ylw "  staking bond_denom      : ${BOND_DENOM:-<none>}"
ylw "  key_rotation_fee.denom  : ${ROT_DENOM:-<none>}"
check "staking bond_denom migrated to abadge" [ "$BOND_DENOM" = "abadge" ]
check "staking key_rotation_fee.denom migrated to abadge" [ "$ROT_DENOM" = "abadge" ]

# 2. Verify bank supply: ubadge == 0, abadge > 0, and denom metadata exponent == 18.
UBADGE_SUP=$(curl -s --max-time 10 "$API_URL/cosmos/bank/v1beta1/supply/by_denom?denom=ubadge" | jq -r '.amount.amount // "0"')
ABADGE_SUP=$(curl -s --max-time 10 "$API_URL/cosmos/bank/v1beta1/supply/by_denom?denom=abadge" | jq -r '.amount.amount // "0"')
META_EXP=$(curl -s --max-time 10 "$API_URL/cosmos/bank/v1beta1/denoms_metadata/abadge" | jq -r '[.metadata.denom_units[] | select(.denom=="badge") | .exponent][0] // 0')
OLD_META=$(curl -s --max-time 10 "$API_URL/cosmos/bank/v1beta1/denoms_metadata/ubadge" | jq -r '.metadata.base // empty')
ylw "  bank ubadge supply      : $UBADGE_SUP"
ylw "  bank abadge supply      : $ABADGE_SUP"
ylw "  badge display exponent  : $META_EXP"
check "ubadge supply is 0 after v37" [ "$UBADGE_SUP" = "0" ]
check "abadge supply is positive after v37" [ "$ABADGE_SUP" != "0" ] && [ -n "$ABADGE_SUP" ]
check "abadge metadata display exponent is 18" [ "$META_EXP" = "18" ]
check "legacy ubadge metadata removed" [ -z "$OLD_META" ]

# 3. Verify poolmanager params migrated to abadge.
PM_JSON=$("$BIN" query poolmanager params --home "$HOME_DIR" --output json 2>/dev/null || true)
PM_SWAP_DENOM=$(jq -r '.params.taker_fee_params.community_pool_denom_to_swap_non_whitelisted_assets_to // empty' <<<"$PM_JSON")
ylw "  poolmanager swap denom  : ${PM_SWAP_DENOM:-<none>}"
check "poolmanager community_pool swap denom is abadge" [ "$PM_SWAP_DENOM" = "abadge" ]

# 4. Verify Cosmos bank send with 18-decimal abadge fee + 1:1 EVM eth_getBalance & transfer.
hex_addr() { "$BIN" debug addr "$1" 2>&1 | grep -i 'hex' | head -1 | awk '{print $NF}' | sed 's/^/0x/'; }
bank_abadge() { "$BIN" query bank balances "$1" --home "$HOME_DIR" --output json 2>/dev/null | jq -r '[.balances[] | select(.denom=="abadge") | .amount][0] // "0"'; }
eth_balance() {
  curl -s --max-time 15 -X POST -H 'Content-Type: application/json' \
    --data "{\"jsonrpc\":\"2.0\",\"method\":\"eth_getBalance\",\"params\":[\"$1\",\"latest\"],\"id\":1}" "$RPC_URL" \
    | jq -r '.result // "0x0"' | python3 -c 'import sys; print(int(sys.stdin.read().strip(), 16))'
}
bigcalc() { python3 -c 'import sys; print(eval(sys.argv[1]))' "$1"; }

"$BIN" keys add evm-v37-sender "${KRT[@]}" --home "$HOME_DIR" --key-type eth_secp256k1 >/dev/null 2>&1 || true
"$BIN" keys add evm-v37-recv   "${KRT[@]}" --home "$HOME_DIR" --key-type eth_secp256k1 >/dev/null 2>&1 || true
A1=$("$BIN" keys show evm-v37-sender -a "${KRT[@]}" --home "$HOME_DIR")
A2=$("$BIN" keys show evm-v37-recv   -a "${KRT[@]}" --home "$HOME_DIR")
H1=$(hex_addr "$A1"); H2=$(hex_addr "$A2")
KEY=$("$BIN" keys unsafe-export-eth-key evm-v37-sender "${KRT[@]}" --home "$HOME_DIR" 2>/dev/null)

FUND="5000000000000000000000" # 5000 BADGE in abadge
OUT=$("$BIN" tx bank send "$FUNDER" "$A1" "${FUND}abadge" --from "$FUNDER" "${KRT[@]}" --home "$HOME_DIR" \
  --chain-id "$CHAIN_ID" --gas auto --gas-adjustment 1.5 --gas-prices 10000000000abadge -y 2>&1 || true)
TXH=$(tx_hash_of "$OUT"); sleep 6
CODE=$("$BIN" query tx "$TXH" --home "$HOME_DIR" --output json 2>/dev/null | jq -r '.code // "?"')
check "post-v37 Cosmos bank send with abadge fee succeeded" [ "$CODE" = "0" ]
check "post-v37 eth_getBalance equals bank abadge 1:1" [ "$(eth_balance "$H1")" = "$FUND" ]

# Transfer an arbitrary 18-decimal wei amount (not divisible by 10^9).
V1="1234567890123456789"
EVM_OUT=$("$EVMTX" --rpc "$RPC_URL" --key "$KEY" --to "$H2" --value-wei "$V1" 2>&1 || true)
TX_STATUS=$(sed -n 's/^status=//p' <<<"$EVM_OUT")
check "post-v37 EVM transfer of 18-decimal wei succeeded" [ "$TX_STATUS" = "0x1" ]
check "post-v37 EVM recipient eth_getBalance == value" [ "$(eth_balance "$H2")" = "$V1" ]
check "post-v37 EVM recipient bank abadge == value (1:1)" [ "$(bank_abadge "$A2")" = "$V1" ]

exit "$FAILED"

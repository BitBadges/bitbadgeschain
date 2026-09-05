#!/usr/bin/env bash
# EVM value-transfer assertions against a running node. Source, then call
# evm_transfer_checks with these in the environment:
#   BIN HOME_DIR CHAIN_ID FUNDER  chain CLI, node home, chain id, funded key (keyring test)
#   RPC_URL API_URL EVMTX         JSON-RPC, REST, and the evmtx helper binary
# Every assertion goes through check(), so $FAILED accumulates in the caller.
#
# Two transfers: one whose value is an exact multiple of 1e9 wei (so the bank
# arithmetic is exact), then one of 1e9+5 wei to walk the fractional path
# through precisebank. abadge, the 18-decimal extended denom, must never show
# up as a bank balance or in supply: precisebank keeps the fraction in its own
# store and x/bank only ever sees whole ubadge.

KRT=(--keyring-backend test)

hex_addr()  { "$BIN" debug addr "$1" 2>&1 | grep -i 'hex' | head -1 | awk '{print $NF}' | sed 's/^/0x/'; }
bank_ubadge() { "$BIN" query bank balances "$1" --home "$HOME_DIR" --output json 2>/dev/null | jq -r '[.balances[] | select(.denom=="ubadge") | .amount][0] // "0"'; }
eth_balance() {
  curl -s --max-time 15 -X POST -H 'Content-Type: application/json' \
    --data "{\"jsonrpc\":\"2.0\",\"method\":\"eth_getBalance\",\"params\":[\"$1\",\"latest\"],\"id\":1}" "$RPC_URL" \
    | jq -r '.result // "0x0"' | python3 -c 'import sys; print(int(sys.stdin.read().strip(), 16))'
}
abadge_owners() { curl -s --max-time 10 "$API_URL/cosmos/bank/v1beta1/denom_owners/abadge" | jq -r '.denom_owners | length'; }
abadge_supply() { curl -s --max-time 10 "$API_URL/cosmos/bank/v1beta1/supply/by_denom?denom=abadge" | jq -r '.amount.amount // "0"'; }
bigcalc() { python3 -c 'import sys; print(eval(sys.argv[1]))' "$1"; }
eth_gas_price() {
  curl -s --max-time 15 -X POST -H 'Content-Type: application/json' \
    --data '{"jsonrpc":"2.0","method":"eth_gasPrice","params":[],"id":1}' "$RPC_URL" \
    | jq -r '.result // "0x0"' | python3 -c 'import sys; print(int(sys.stdin.read().strip(), 16))'
}

# evm_send <privkey hex> <to 0x> <value wei>: sets TX_STATUS TX_GAS_USED TX_GAS_PRICE
evm_send() {
  local out
  TX_STATUS="" TX_GAS_USED=0 TX_GAS_PRICE=0
  out=$("$EVMTX" --rpc "$RPC_URL" --key "$1" --to "$2" --value-wei "$3" 2>&1) || true
  while IFS= read -r line; do echo "     $line"; done <<<"$out"
  TX_STATUS=$(sed -n 's/^status=//p' <<<"$out")
  TX_GAS_USED=$(sed -n 's/^gasUsed=//p' <<<"$out"); TX_GAS_USED=${TX_GAS_USED:-0}
  TX_GAS_PRICE=$(sed -n 's/^gasPrice=//p' <<<"$out"); TX_GAS_PRICE=${TX_GAS_PRICE:-0}
}

evm_transfer_checks() {
  local a1 a2 h1 h2 key v1 v2 fee1 s_before r_before s_after r_after
  "$BIN" keys add evm-sender "${KRT[@]}" --home "$HOME_DIR" --key-type eth_secp256k1 >/dev/null 2>&1 || true
  "$BIN" keys add evm-recv   "${KRT[@]}" --home "$HOME_DIR" --key-type eth_secp256k1 >/dev/null 2>&1 || true
  a1=$("$BIN" keys show evm-sender -a "${KRT[@]}" --home "$HOME_DIR")
  a2=$("$BIN" keys show evm-recv   -a "${KRT[@]}" --home "$HOME_DIR")
  h1=$(hex_addr "$a1"); h2=$(hex_addr "$a2")
  key=$("$BIN" keys unsafe-export-eth-key evm-sender "${KRT[@]}" --home "$HOME_DIR" 2>/dev/null)
  ylw "  evm sender $h1 ($a1)"
  ylw "  evm recv   $h2 ($a2)"

  # Fund the sender: 2 BADGE plus fees for two 21000-gas transfers at the
  # current eth_gasPrice. The EIP-1559 base fee starts at genesis
  # (1e9 ubadge/gas, i.e. ~1e18 wei/gas) and decays 12.5% per empty block,
  # so on a young chain a single transfer can cost thousands of BADGE; the
  # base fee only falls over time, so a fee budget read now stays sufficient.
  local out txh code gp fund
  gp=$(eth_gas_price)
  fund=$(bigcalc "2*10**9 + (2*21000*(($gp+10**9-1)//10**9*10**9))//10**9 + 1")
  ylw "  eth_gasPrice $gp wei/gas; funding sender with ${fund}ubadge"
  out=$("$BIN" tx bank send "$FUNDER" "$a1" "${fund}ubadge" --from "$FUNDER" "${KRT[@]}" --home "$HOME_DIR" \
    --chain-id "$CHAIN_ID" --gas auto --gas-adjustment 1.5 --gas-prices 10ubadge -y 2>&1 || true)
  txh=$(tx_hash_of "$out"); sleep 6
  code=$("$BIN" query tx "$txh" --home "$HOME_DIR" --output json 2>/dev/null | jq -r '.code // "?"')
  check "funding bank send to the eth key succeeded (code 0)" [ "$code" = 0 ]
  check "funded eth key is visible over eth_getBalance (fund x 1e9)" [ "$(eth_balance "$h1")" = "$(bigcalc "$fund*10**9")" ]

  # --- transfer 1: exact multiple of 1e9 wei -----------------------------
  v1=500000000000000000   # 0.5 BADGE = 500000000 ubadge exactly
  s_before=$(bank_ubadge "$a1"); r_before=$(bank_ubadge "$a2")
  evm_send "$key" "$h2" "$v1"
  check "transfer 1: receipt status 0x1" [ "$TX_STATUS" = 0x1 ]
  fee1=$(bigcalc "$TX_GAS_USED*$TX_GAS_PRICE")
  s_after=$(bank_ubadge "$a1"); r_after=$(bank_ubadge "$a2")
  ylw "  sender bank $s_before -> $s_after, recipient bank $r_before -> $r_after, fee $fee1 wei"
  check "transfer 1: recipient eth_getBalance == value"          [ "$(eth_balance "$h2")" = "$v1" ]
  check "transfer 1: recipient bank ubadge increased by value/1e9" [ "$r_after" = "$(bigcalc "$r_before+$v1//10**9")" ]
  check "transfer 1: sender bank ubadge decreased by (value+fee)/1e9" [ "$s_after" = "$(bigcalc "$s_before-($v1+$fee1)//10**9")" ]
  check "transfer 1: fee is a whole number of ubadge"             [ "$(bigcalc "$fee1%10**9")" = 0 ]
  check "no account holds abadge in x/bank"                       [ "$(abadge_owners)" = 0 ]
  check "abadge bank supply is 0"                                 [ "$(abadge_supply)" = 0 ]

  # --- transfer 2: fractional ubadge ---------------------------------------
  v2=1000000005
  s_before=$s_after
  evm_send "$key" "$h2" "$v2"
  check "transfer 2 (1e9+5 wei): receipt status 0x1" [ "$TX_STATUS" = 0x1 ]
  s_after=$(bank_ubadge "$a1"); r_after=$(bank_ubadge "$a2")
  check "transfer 2: recipient eth_getBalance == v1+v2 exactly"   [ "$(eth_balance "$h2")" = "$(bigcalc "$v1+$v2")" ]
  check "transfer 2: recipient bank ubadge holds the whole part"  [ "$r_after" = "$(bigcalc "$r_before+($v1+$v2)//10**9")" ]
  check "transfer 2: sender bank ubadge decreased by >= 1"        [ "$s_after" -le "$(bigcalc "$s_before-1")" ]
  check "transfer 2: still no abadge bank owners"                 [ "$(abadge_owners)" = 0 ]
  check "transfer 2: abadge bank supply still 0"                  [ "$(abadge_supply)" = 0 ]
}

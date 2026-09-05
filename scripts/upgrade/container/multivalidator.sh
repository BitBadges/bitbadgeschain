#!/usr/bin/env bash
# 4-validator coordinated upgrade: FROM -> TO.
#
# The single-validator rehearsal proves the handler and the operator
# experience. It cannot prove consensus: one validator is consensus, so a
# nondeterministic result is invisible. Four independent processes, each with
# its own state and Go runtime, must agree on the app hash at every checkpoint
# (before the upgrade, in a block that provably contains a transaction, at the
# upgrade height, after it, and in a post-upgrade transaction block).
set -euo pipefail
# shellcheck source=../lib/common.sh
. "$(dirname "$0")/../lib/common.sh"

CHAIN_ID=bitbadges-mv-1
N=4
KR=(--keyring-backend test)
PROPOSE="$(dirname "$0")/../propose.sh"

home(){ echo "/root/.n$1"; }
rpc(){  echo $((26657 + $1 * 10)); }
p2p(){  echo $((26656 + $1 * 10)); }

step "1. Initialise a $N-validator genesis with the FROM binary"
rm -rf /root/.n*
for i in $(seq 0 $((N-1))); do
  "$FROM_BIN" init "node$i" --home "$(home "$i")" --chain-id "$CHAIN_ID" >/dev/null 2>&1
done
# One shared keyring holds the account keys; each node signs its gentx with
# its own priv_validator_key, which is what makes them distinct validators.
for i in $(seq 0 $((N-1))); do
  "$FROM_BIN" keys add "val$i" "${KR[@]}" --home "$(home 0)" >/dev/null 2>&1
done
G0=$(home 0)/config/genesis.json
BOND_DENOM=$(jq -r '.app_state.staking.params.bond_denom' "$G0")
DEPOSIT_DENOM=$(jq -r '.app_state.gov.params.min_deposit[0].denom' "$G0")
FUNDS="100000000000000${BOND_DENOM},100000000000000ubadge"
[ "$DEPOSIT_DENOM" = "$BOND_DENOM" ] || [ "$DEPOSIT_DENOM" = ubadge ] || FUNDS="$FUNDS,100000000000000${DEPOSIT_DENOM}"
for i in $(seq 0 $((N-1))); do
  "$FROM_BIN" genesis add-genesis-account "val$i" "$FUNDS" --home "$(home 0)" "${KR[@]}" >/dev/null 2>&1
done
genesis_json_set "$G0" '.app_state.gov.params.voting_period="20s"
  | .app_state.gov.params.expedited_voting_period="15s"
  | .app_state.gov.params.max_deposit_period="30s"
  | .app_state.gov.params.min_deposit[0].amount="1"
  | .app_state.gov.params.expedited_min_deposit[0].amount="2"'

mkdir -p "$(home 0)/config/gentx"
for i in $(seq 0 $((N-1))); do
  if [ "$i" -ne 0 ]; then
    cp "$G0" "$(home "$i")/config/genesis.json"
    cp -r "$(home 0)/keyring-test" "$(home "$i")/"
  fi
  "$FROM_BIN" genesis gentx "val$i" "50000000000${BOND_DENOM}" --home "$(home "$i")" \
    --chain-id "$CHAIN_ID" "${KR[@]}" >/dev/null 2>&1
  if [ "$i" -ne 0 ]; then cp "$(home "$i")"/config/gentx/*.json "$(home 0)/config/gentx/"; fi
done
"$FROM_BIN" genesis collect-gentxs --home "$(home 0)" >/dev/null 2>&1
"$FROM_BIN" genesis validate-genesis --home "$(home 0)" >/dev/null 2>&1 || die "genesis invalid"
VALCOUNT=$(jq '.app_state.genutil.gen_txs | length' "$G0")
check "genesis carries $N gentxs" [ "$VALCOUNT" = "$N" ]
for i in $(seq 1 $((N-1))); do cp "$G0" "$(home "$i")/config/genesis.json"; done

step "2. Wire ports and peers"
PEERS=""
for i in $(seq 0 $((N-1))); do
  ID=$("$FROM_BIN" comet show-node-id --home "$(home "$i")" 2>/dev/null)
  [ -n "$PEERS" ] && PEERS="$PEERS,"
  PEERS="$PEERS$ID@127.0.0.1:$(p2p "$i")"
done
for i in $(seq 0 $((N-1))); do
  C=$(home "$i")/config/config.toml
  A=$(home "$i")/config/app.toml
  sed -i "s|^laddr = \"tcp://127.0.0.1:26657\"|laddr = \"tcp://0.0.0.0:$(rpc "$i")\"|" "$C"
  sed -i "s|^laddr = \"tcp://0.0.0.0:26656\"|laddr = \"tcp://0.0.0.0:$(p2p "$i")\"|" "$C"
  sed -i "s|^persistent_peers = \".*\"|persistent_peers = \"$PEERS\"|" "$C"
  # All four share 127.0.0.1, which CometBFT rejects by default.
  sed -i 's|^allow_duplicate_ip = false|allow_duplicate_ip = true|' "$C"
  sed -i 's|^addr_book_strict = true|addr_book_strict = false|' "$C"
  sed -i 's|^pprof_laddr = ".*"|pprof_laddr = ""|' "$C"
  python3 - "$A" "$i" <<'PY'
import re,sys
p,i=sys.argv[1],int(sys.argv[2])
s=open(p).read()
def sec(name,key,val,s):
    m=re.search(r'(?ms)^\['+re.escape(name)+r'\].*?(?=^\[|\Z)',s)
    if not m: return s
    b=re.sub(r'(?m)^'+re.escape(key)+r'\s*=.*$',f'{key} = {val}',m.group(0))
    return s[:m.start()]+b+s[m.end():]
s=sec('grpc','address',f'"127.0.0.1:{9090+i}"',s)
s=sec('grpc-web','address',f'"127.0.0.1:{9091+i*10}"',s)
s=sec('api','address',f'"tcp://127.0.0.1:{1317+i}"',s)
s=sec('api','enable','true',s)
s=sec('json-rpc','address',f'"127.0.0.1:{8545+i*10}"',s)
s=sec('json-rpc','ws-address',f'"127.0.0.1:{8546+i*10}"',s)
open(p,'w').write(s)
PY
  mkdir -p "$(home "$i")/cosmovisor/genesis/bin" "$(home "$i")/cosmovisor/upgrades/$UPGRADE_NAME/bin"
  cp "$FROM_BIN" "$(home "$i")/cosmovisor/genesis/bin/bitbadgeschaind"
  cp "$TO_BIN" "$(home "$i")/cosmovisor/upgrades/$UPGRADE_NAME/bin/bitbadgeschaind"
done
grn "peers: $PEERS"

step "3. Start all $N validators under cosmovisor"
export DAEMON_NAME=bitbadgeschaind DAEMON_ALLOW_DOWNLOAD_BINARIES=false \
       DAEMON_RESTART_AFTER_UPGRADE=true UNSAFE_SKIP_BACKUP=true
PIDS=()
for i in $(seq 0 $((N-1))); do
  DAEMON_HOME=$(home "$i") "$CV" run start --home "$(home "$i")" --minimum-gas-prices "0${BOND_DENOM}" \
    > "$LOG_DIR/mv-node$i.log" 2>&1 &
  PIDS+=($!)
done
height_of(){ curl -s --max-time 5 "http://127.0.0.1:$(rpc "$1")/status" 2>/dev/null | jq -r '.result.sync_info.latest_block_height // empty' 2>/dev/null || true; }
apphash_at(){ curl -s --max-time 5 "http://127.0.0.1:$(rpc "$1")/block?height=$2" 2>/dev/null | jq -r '.result.block.header.app_hash // empty' 2>/dev/null || true; }
NODE0="tcp://127.0.0.1:$(rpc 0)"
API0="http://127.0.0.1:1317"

# fee_flags: the gas price the chain demands right now, read from node0. An
# upgrade may raise the fee-market floor (v35 does), so a fee that cleared
# before the upgrade is not evidence it clears after it. For Cosmos txs the
# ante enforces feemarket min_gas_price in ubadge per gas; the EIP-1559
# base_fee (18-decimal EVM units) gates EVM txs only, so it is not used here.
fee_flags(){
  local mgp
  mgp=$(curl -s --max-time 5 "$API0/cosmos/evm/feemarket/v1/params" | jq -r '.params.min_gas_price // "0"' 2>/dev/null || echo 0)
  echo "--gas-prices $(python3 -c 'import sys; from decimal import Decimal; print(format(Decimal(sys.argv[1] or 0), "f"))' "$mgp")ubadge"
}

for _ in $(seq 1 90); do
  H=$(height_of 0); [ -n "${H:-}" ] && [ "$H" -ge 3 ] 2>/dev/null && break; sleep 2
done
[ -n "${H:-}" ] || { tail -20 "$LOG_DIR/mv-node0.log"; die "network never started"; }
UP=0; for i in $(seq 0 $((N-1))); do [ -n "$(height_of "$i")" ] && UP=$((UP+1)); done
check "all $N validators are online and syncing" [ "$UP" = "$N" ]
grn "network producing blocks at height $H"

compare_apphashes(){
  local label=$1 h=$2 base="" ok=1
  for i in $(seq 0 $((N-1))); do
    local a; a=$(apphash_at "$i" "$h")
    if [ -z "$a" ]; then ylw "     node$i had no block $h"; ok=0; continue; fi
    if [ -z "$base" ]; then base=$a; elif [ "$a" != "$base" ]; then
      red "     DIVERGENCE at height $h: node$i=$a vs node0=$base"; ok=0
    fi
  done
  if [ "$ok" = 1 ] && [ -n "$base" ]; then
    grn "  PASS  $label: all $N agree at height $h (${base:0:16}...)"
  else
    red "  FAIL  $label: app hashes differ or missing at height $h"; FAILED=$((FAILED+1))
  fi
}

# send_and_confirm <label> <bin> <amount> -- asserts the transfer executed
# with code 0 and sets TX_HEIGHT to its block. Uses a global on purpose: a
# command substitution would swallow the PASS/FAIL output and the FAILED count.
TX_HEIGHT=""
send_and_confirm(){
  local label=$1 bin=$2 amount=$3 out code h txh fees
  TX_HEIGHT=""
  read -r -a fees <<<"$(fee_flags)"
  ylw "     paying ${fees[*]}"
  out=$("$bin" tx bank send "$A0" "$A1" "${amount}ubadge" --from val0 "${KR[@]}" \
    --home "$(home 0)" --node "$NODE0" --chain-id "$CHAIN_ID" \
    --gas auto --gas-adjustment 1.5 "${fees[@]}" -y 2>&1) || true
  txh=$(tx_hash_of "$out" || true)
  if [ -z "$txh" ]; then
    red "  FAIL  $label: transfer was never submitted"; echo "$out" | tail -3
    FAILED=$((FAILED+1)); return
  fi
  sleep 6
  local q; q=$("$bin" query tx "$txh" --home "$(home 0)" --node "$NODE0" --output json 2>/dev/null)
  code=$(jq -r '.code // "?"' <<<"$q"); h=$(jq -r '.height // empty' <<<"$q")
  if [ "$code" != 0 ] || [ -z "$h" ]; then
    red "  FAIL  $label: transfer did not execute (code=$code height=$h)"
    FAILED=$((FAILED+1)); return
  fi
  grn "  PASS  $label: executed in block $h (code 0)"
  TX_HEIGHT=$h
}
A0=$("$FROM_BIN" keys show val0 -a "${KR[@]}" --home "$(home 0)")
A1=$("$FROM_BIN" keys show val1 -a "${KR[@]}" --home "$(home 0)")

step "4. Pre-upgrade determinism"
compare_apphashes "pre-upgrade" 3

step "5. Real transaction traffic, pre-upgrade"
send_and_confirm "pre-upgrade transfer" "$FROM_BIN" 1000
[ -n "$TX_HEIGHT" ] && compare_apphashes "block containing the pre-upgrade tx" "$TX_HEIGHT"

step "6. Propose and pass the $UPGRADE_NAME upgrade"
PROPOSE_OUT=$("$PROPOSE" --name "$UPGRADE_NAME" --home "$(home 0)" --from val0 --voters val1,val2,val3 \
  --bin "$FROM_BIN" --node "$NODE0" --chain-id "$CHAIN_ID" --keyring-backend test \
  --deposit "10${DEPOSIT_DENOM}" --fees "0${BOND_DENOM}" --height +30)
echo "$PROPOSE_OUT"
UPGRADE_HEIGHT=$(sed -n 's/^UPGRADE_HEIGHT=//p' <<<"$PROPOSE_OUT")
PID_=$(sed -n 's/^PROPOSAL_ID=//p' <<<"$PROPOSE_OUT")
sleep 22
STATUS=$("$FROM_BIN" query gov proposal "$PID_" --home "$(home 0)" --node "$NODE0" --output json 2>/dev/null | jq -r '.proposal.status')
check "proposal $PID_ passed with all $N voting" [ "$STATUS" = "PROPOSAL_STATUS_PASSED" ]
grn "upgrade height $UPGRADE_HEIGHT (current $(height_of 0))"

step "7. Coordinated upgrade - no intervention"
DEADLINE=$((SECONDS + 600))
while [ $SECONDS -lt $DEADLINE ]; do
  CH=$(height_of 0)
  [ -n "${CH:-}" ] && [ "$CH" -gt $((UPGRADE_HEIGHT + 8)) ] 2>/dev/null && break
  sleep 5
done

step "8. Results"
for i in $(seq 0 $((N-1))); do ylw "  node$i height: $(height_of "$i")"; done
ALIVE=0; for i in $(seq 0 $((N-1))); do
  h=$(height_of "$i"); [ -n "$h" ] && [ "$h" -gt "$UPGRADE_HEIGHT" ] 2>/dev/null && ALIVE=$((ALIVE+1))
done
check "all $N validators resumed past the upgrade height" [ "$ALIVE" = "$N" ]
compare_apphashes "at the upgrade height" "$UPGRADE_HEIGHT"
compare_apphashes "post-upgrade"          "$((UPGRADE_HEIGHT + 5))"
APPLIED=0; HOOKED=0
for i in $(seq 0 $((N-1))); do
  grep -q "applying upgrade \"$UPGRADE_NAME\"" "$LOG_DIR/mv-node$i.log" && APPLIED=$((APPLIED+1))
  grep -q 'pre-upgrade' "$LOG_DIR/mv-node$i.log" && HOOKED=$((HOOKED+1))
done
check "all $N applied the $UPGRADE_NAME handler" [ "$APPLIED" = "$N" ]
check "all $N ran the pre-upgrade hook" [ "$HOOKED" = "$N" ]

step "9. Post-upgrade transaction"
send_and_confirm "post-upgrade transfer" "$TO_BIN" 500
[ -n "$TX_HEIGHT" ] && compare_apphashes "block containing the post-upgrade tx" "$TX_HEIGHT"

for p in "${PIDS[@]}"; do kill "$p" 2>/dev/null || true; done
wait 2>/dev/null || true
summary "ALL CHECKS PASSED - $N validators agreed on every app hash through the upgrade"
exit "$FAILED"

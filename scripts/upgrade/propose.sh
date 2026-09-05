#!/usr/bin/env bash
# Submit a MsgSoftwareUpgrade proposal and vote yes on it. Non-interactive.
#
# Works inside the rehearsal container and against a real node. Reads nothing
# from other repositories; every input is a flag.
#
#   propose.sh --name v35 --home ~/.bitbadgeschain --from mykey \
#     [--deposit 10000000ubadge] [--height +N | --height 12345] \
#     [--chain-id X] [--node tcp://...] [--bin bitbadgeschaind] \
#     [--keyring-backend os] [--fees 0ubadge | --gas-prices 10ubadge] \
#     [--expedited] [--info '<binaries json>'] [--voters key1,key2] \
#     [--no-vote] [--wait 6] [--dry-run]
#
# Prints PROPOSAL_ID=<id> and UPGRADE_HEIGHT=<height> on stdout when done.
set -euo pipefail

NAME="" HOME_DIR="" FROM="" DEPOSIT="10000000ubadge" HEIGHT="+30" CHAIN_ID="" NODE=""
BIN="bitbadgeschaind" KR="os" FEES="" GAS_PRICES="" EXPEDITED=false INFO="" VOTERS=""
VOTE=true WAIT=6 DRY_RUN=false AUTHORITY=""

usage() { sed -n '2,15p' "$0"; exit "${1:-0}"; }
while [ $# -gt 0 ]; do
  case $1 in
    --name) NAME=$2; shift 2 ;;
    --home) HOME_DIR=$2; shift 2 ;;
    --from) FROM=$2; shift 2 ;;
    --deposit) DEPOSIT=$2; shift 2 ;;
    --height) HEIGHT=$2; shift 2 ;;
    --chain-id) CHAIN_ID=$2; shift 2 ;;
    --node) NODE=$2; shift 2 ;;
    --bin) BIN=$2; shift 2 ;;
    --keyring-backend) KR=$2; shift 2 ;;
    --fees) FEES=$2; shift 2 ;;
    --gas-prices) GAS_PRICES=$2; shift 2 ;;
    --expedited) EXPEDITED=true; shift ;;
    --info) INFO=$2; shift 2 ;;
    --voters) VOTERS=$2; shift 2 ;;
    --no-vote) VOTE=false; shift ;;
    --wait) WAIT=$2; shift 2 ;;
    --authority) AUTHORITY=$2; shift 2 ;;
    --dry-run) DRY_RUN=true; shift ;;
    -h|--help) usage ;;
    *) echo "unknown flag: $1" >&2; usage 1 ;;
  esac
done
[ -n "$NAME" ] && [ -n "$HOME_DIR" ] && [ -n "$FROM" ] || { echo "--name, --home and --from are required" >&2; usage 1; }
[[ $DEPOSIT =~ ^[0-9]+[a-z][a-z0-9/]*$ ]] || { echo "--deposit must look like 10000000ubadge, got '$DEPOSIT'" >&2; exit 1; }
[[ $DEPOSIT == *ustake ]] && { echo "--deposit uses 'ustake'; this chain's denom is 'ubadge'" >&2; exit 1; }

COMMON=(--home "$HOME_DIR")
[ -n "$NODE" ] && COMMON+=(--node "$NODE")
TXFLAGS=(--from "$FROM" --keyring-backend "$KR" --gas auto --gas-adjustment 1.5 -y --output json)
if [ -n "$GAS_PRICES" ]; then TXFLAGS+=(--gas-prices "$GAS_PRICES"); else TXFLAGS+=(--fees "${FEES:-0ubadge}"); fi

q() { "$BIN" "$@" "${COMMON[@]}" --output json 2>/dev/null; }
# The CLI prints "gas estimate: N" ahead of the JSON when --gas auto is used.
json_of() { sed -n '/^{/,$p' <<<"$1"; }

if $DRY_RUN; then
  CHAIN_ID=${CHAIN_ID:-"<from node status>"}
  CURRENT="<current height>"
  AUTHORITY=${AUTHORITY:-"<gov module account>"}
else
  CHAIN_ID=${CHAIN_ID:-$("$BIN" status "${COMMON[@]}" 2>/dev/null | jq -r '.node_info.network')}
  CURRENT=$("$BIN" status "${COMMON[@]}" 2>/dev/null | jq -r '.sync_info.latest_block_height')
  [ -n "$CHAIN_ID" ] && [ -n "$CURRENT" ] && [ "$CURRENT" != null ] || { echo "cannot reach the node (--node/--home wrong?)" >&2; exit 1; }
  AUTHORITY=${AUTHORITY:-$(q query auth module-account gov | jq -r '.account.value.address // .account.base_account.address // empty')}
  [ -n "$AUTHORITY" ] || { echo "could not resolve the gov module account; pass --authority" >&2; exit 1; }
fi
TXFLAGS+=(--chain-id "$CHAIN_ID")

case $HEIGHT in
  +*) if $DRY_RUN; then UPGRADE_HEIGHT="$CURRENT + ${HEIGHT#+}"; else UPGRADE_HEIGHT=$((CURRENT + ${HEIGHT#+})); fi ;;
  *)  UPGRADE_HEIGHT=$HEIGHT ;;
esac

PROPOSAL="$HOME_DIR/proposal-$NAME.json"
$DRY_RUN && PROPOSAL=/dev/stdout
INFO_JSON=$(jq -cn --arg i "$INFO" '$i')
cat > "$PROPOSAL" <<JSON
{
  "messages": [
    {
      "@type": "/cosmos.upgrade.v1beta1.MsgSoftwareUpgrade",
      "authority": "$AUTHORITY",
      "plan": { "name": "$NAME", "height": "$UPGRADE_HEIGHT", "info": $INFO_JSON, "upgraded_client_state": null }
    }
  ],
  "expedited": $EXPEDITED,
  "deposit": "$DEPOSIT",
  "title": "Upgrade to $NAME",
  "summary": "This proposal upgrades the chain to version $NAME at height $UPGRADE_HEIGHT."
}
JSON

if $DRY_RUN; then
  echo "would run: $BIN tx gov submit-proposal <proposal> ${TXFLAGS[*]} ${COMMON[*]}"
  $VOTE && echo "would run: $BIN tx gov vote <id> yes ${TXFLAGS[*]} ${COMMON[*]}"
  echo "PROPOSAL_ID=<dry-run>"
  echo "UPGRADE_HEIGHT=$UPGRADE_HEIGHT"
  exit 0
fi

OUT=$("$BIN" tx gov submit-proposal "$PROPOSAL" "${TXFLAGS[@]}" "${COMMON[@]}" 2>&1) || { echo "$OUT" >&2; exit 1; }
TXHASH=$(json_of "$OUT" | jq -r '.txhash // empty' 2>/dev/null || true)
[ -n "$TXHASH" ] || { echo "submit-proposal returned no txhash: $OUT" >&2; exit 1; }
sleep "$WAIT"

TX=$(q query tx "$TXHASH" || true)
CODE=$(jq -r '.code // 1' <<<"$TX")
[ "$CODE" = 0 ] || { echo "submit-proposal tx $TXHASH failed: $(jq -r '.raw_log // "not found yet"' <<<"$TX")" >&2; exit 1; }
PROPOSAL_ID=$(jq -r '[.events[] | select(.type=="submit_proposal") | .attributes[] | select(.key=="proposal_id") | .value][0] // empty' <<<"$TX")
[ -n "$PROPOSAL_ID" ] || PROPOSAL_ID=$(q query gov proposals | jq -r '.proposals | sort_by(.id|tonumber) | last.id')
[ -n "$PROPOSAL_ID" ] && [ "$PROPOSAL_ID" != null ] || { echo "no proposal id found for tx $TXHASH" >&2; exit 1; }

if $VOTE; then
  IFS=, read -r -a EXTRA <<<"$VOTERS"
  for voter in "$FROM" "${EXTRA[@]}"; do
    [ -n "$voter" ] || continue
    VFLAGS=("${TXFLAGS[@]}")
    VFLAGS[1]=$voter   # --from <voter>
    VOUT=$("$BIN" tx gov vote "$PROPOSAL_ID" yes "${VFLAGS[@]}" "${COMMON[@]}" 2>&1) || { echo "$VOUT" >&2; exit 1; }
    [ "$(json_of "$VOUT" | jq -r '.code // 1')" = 0 ] || { echo "vote from $voter rejected: $VOUT" >&2; exit 1; }
    sleep 2
  done
fi

echo "PROPOSAL_ID=$PROPOSAL_ID"
echo "UPGRADE_HEIGHT=$UPGRADE_HEIGHT"

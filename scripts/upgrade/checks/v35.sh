#!/usr/bin/env bash
# Post-upgrade assertions for v35. Run by container/upgrade.sh after the
# upgrade applied, with BIN (the v35 binary), HOME_DIR, UPGRADE_HEIGHT, LOG
# and API_URL (REST) in the environment. Exit code = number of failed checks.
set -uo pipefail
# shellcheck source=../lib/common.sh
. "$(dirname "$0")/../lib/common.sh"

MAX_GAS=$("$BIN" query consensus params --home "$HOME_DIR" --output json 2>/dev/null | jq -r '.params.block.max_gas // empty')
ylw "  consensus block.max_gas : ${MAX_GAS:-<none>}"
check "block max_gas is 100000000 after v35" [ "${MAX_GAS:-}" = 100000000 ]

# The chain CLI registers no feemarket query command, so read it over REST.
MIN_GAS_PRICE=$(curl -s --max-time 10 "$API_URL/cosmos/evm/feemarket/v1/params" | jq -r '.params.min_gas_price // empty')
ylw "  feemarket min_gas_price : ${MIN_GAS_PRICE:-<none>}"
check "feemarket min_gas_price is 10 after v35" \
  python3 -c 'import sys; from decimal import Decimal; sys.exit(0 if Decimal(sys.argv[1]) == 10 else 1)' "${MIN_GAS_PRICE:-0}"

exit "$FAILED"

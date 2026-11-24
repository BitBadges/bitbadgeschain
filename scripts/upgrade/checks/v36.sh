#!/usr/bin/env bash
# Post-upgrade assertions for v36. Run by container/upgrade.sh after the
# upgrade applied, with BIN (the v36 binary), HOME_DIR, UPGRADE_HEIGHT, LOG,
# API_URL (REST), RPC_URL (JSON-RPC), EVMTX, CHAIN_ID and FUNDER in the
# environment. Exit code = number of failed checks.
set -uo pipefail
# shellcheck source=../lib/common.sh
. "$(dirname "$0")/../lib/common.sh"

STAKING_JSON=$("$BIN" query staking params --home "$HOME_DIR" --output json 2>/dev/null || true)
BOND_DENOM=$(jq -r '.params.bond_denom // empty' <<<"$STAKING_JSON")
ROT_DENOM=$(jq -r '.params.key_rotation_fee.denom // empty' <<<"$STAKING_JSON")
ylw "  staking bond_denom      : ${BOND_DENOM:-<none>}"
ylw "  key_rotation_fee.denom  : ${ROT_DENOM:-<none>}"
check "staking key_rotation_fee.denom matches bond_denom after v36" \
  [ -n "${BOND_DENOM:-}" ] && [ "${ROT_DENOM:-}" = "${BOND_DENOM:-}" ]

PM_TAKER=$("$BIN" query poolmanager params --home "$HOME_DIR" --output json 2>/dev/null | jq -r '.params.taker_fee_params.default_taker_fee // empty')
ylw "  poolmanager taker_fee   : ${PM_TAKER:-<none>}"
check "poolmanager params migrated and queryable after v36" [ -n "${PM_TAKER:-}" ]

# EVM value transfers on the upgraded node (see lib/evm_transfer.sh).
# shellcheck source=../lib/evm_transfer.sh
. "$(dirname "$0")/../lib/evm_transfer.sh"
evm_transfer_checks

exit "$FAILED"

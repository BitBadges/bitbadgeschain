#!/usr/bin/env bash
set -euo pipefail
HERE=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
TEST_HOME=$(mktemp -d)
trap 'rm -rf "$TEST_HOME"' EXIT
FAKE="$HERE/testdata/proposal_node.sh"
if OUT=$(PROPOSAL_TEST_NETWORK=bitbadges-1 PROPOSAL_TEST_EVENTS=one "$HERE/../propose.sh" --name v35 --home "$TEST_HOME" --from test --chain-id bitbadges-1 --authority gov --height 123 --bin "$FAKE" --wait 0 --no-vote 2>&1); then
 echo "FAIL: mainnet proposal accepted without an explicit mainnet flag"
 exit 1
fi
[[ $OUT == *"--allow-mainnet"* ]]
echo "mainnet proposal guard: passed"

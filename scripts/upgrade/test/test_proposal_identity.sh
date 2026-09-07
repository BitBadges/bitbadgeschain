#!/usr/bin/env bash
set -euo pipefail
HERE=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
TEST_HOME=$(mktemp -d)
trap 'rm -rf "$TEST_HOME"' EXIT
FAKE="$HERE/testdata/proposal_node.sh"
chmod +x "$FAKE"

for events in missing ambiguous; do
  if OUT=$(PROPOSAL_TEST_EVENTS=$events "$HERE/../propose.sh" --name v35 --home "$TEST_HOME" --from test --chain-id local --authority gov --height 123 --bin "$FAKE" --wait 0 --no-vote 2>&1); then
    echo "FAIL: $events proposal identity was accepted"
    exit 1
  fi
  [[ $OUT == *"proposal id"* ]] || { echo "$OUT"; exit 1; }
done
OUT=$(PROPOSAL_TEST_EVENTS=one "$HERE/../propose.sh" --name v35 --home "$TEST_HOME" --from test --chain-id local --authority gov --height 123 --bin "$FAKE" --wait 0 --no-vote)
[[ $OUT == *"PROPOSAL_ID=7"* ]]
OUT=$(PROPOSAL_TEST_EVENTS=one "$HERE/../propose.sh" --name v35 --home "$TEST_HOME" --from test --chain-id local --authority gov --height 123 --bin "$FAKE" --wait 0)
[[ $OUT == *"PROPOSAL_ID=7"* ]]
echo "proposal identity: all passed"

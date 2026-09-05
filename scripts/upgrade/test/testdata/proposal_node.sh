#!/usr/bin/env bash
set -euo pipefail
case "$1 $2" in
  'status --home') echo '{"sync_info":{"latest_block_height":"100"}}' ;;
  'query tx')
    case ${PROPOSAL_TEST_EVENTS:-missing} in
      one) echo '{"code":0,"events":[{"type":"submit_proposal","attributes":[{"key":"proposal_id","value":"7"}]}]}' ;;
      ambiguous) echo '{"code":0,"events":[{"type":"submit_proposal","attributes":[{"key":"proposal_id","value":"7"},{"key":"proposal_id","value":"8"}]}]}' ;;
      *) echo '{"code":0,"events":[]}' ;;
    esac ;;
  'query gov') echo '{"proposals":[{"id":"999"}]}' ;;
  'tx gov') echo '{"code":0,"txhash":"ABC"}' ;;
  *) exit 1 ;;
esac

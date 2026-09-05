#!/usr/bin/env bash
# Runs every scripts/upgrade self-test plus shellcheck (local binary, else the
# koalaman/shellcheck Docker image). Exit 0 only when everything passes.
set -euo pipefail
HERE=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
UP=$(cd "$HERE/.." && pwd)
RC=0
for t in "$HERE"/test_*.sh; do
  echo "### $(basename "$t")"
  bash "$t" || RC=1
  echo
done
echo "### shellcheck"
cd "$UP"
FILES=(rehearse.sh new-version.sh propose.sh lib/common.sh container/*.sh checks/*.sh test/*.sh)
if command -v shellcheck >/dev/null; then
  if shellcheck -x -P SCRIPTDIR -S style "${FILES[@]}"; then echo "shellcheck: clean"; else RC=1; fi
elif command -v docker >/dev/null; then
  if docker run --rm -v "$UP:/mnt:ro" -w /mnt koalaman/shellcheck:stable -x -P SCRIPTDIR -S style "${FILES[@]}"; then echo "shellcheck: clean"; else RC=1; fi
else
  echo "shellcheck: skipped (no shellcheck or docker)"
fi
exit $RC

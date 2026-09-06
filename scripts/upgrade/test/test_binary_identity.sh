#!/usr/bin/env bash
set -euo pipefail
HERE=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
CHECK="$HERE/../container/verify-binaries.sh"
T=$(mktemp -d)
trap 'rm -rf "$T"' EXIT
cat > "$T/chain" <<'SH'
#!/usr/bin/env bash
[ "$*" = 'version --long --output json' ] || exit 2
printf '%s\n' "${VERSION_JSON:-}"
exit "${VERSION_RC:-0}"
SH
chmod +x "$T/chain"

VERSION_JSON='{"commit":"expected"}' bash "$CHECK" "$T/chain" expected "$T/chain" expected
for value in '{"commit":"stale"}' '{}' '{"commit":null}' 'invalid'; do
  if VERSION_JSON="$value" bash "$CHECK" "$T/chain" expected "$T/chain" expected; then
    echo "FAIL: accepted incorrect or missing binary identity" >&2
    exit 1
  fi
done
if VERSION_JSON='{"commit":"expected"}' VERSION_RC=3 bash "$CHECK" "$T/chain" expected "$T/chain" expected; then
  echo "FAIL: accepted failed version command" >&2
  exit 1
fi
if VERSION_JSON='{"commit":"expected"}' bash "$CHECK" "$T/chain" expected "$T/chain" different; then
  echo "FAIL: did not validate the destination binary" >&2
  exit 1
fi
echo "binary identity: all passed"

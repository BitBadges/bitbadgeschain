#!/usr/bin/env bash
# Start a new chain version: scripts/upgrade/new-version.sh vNN [options]
#
#   --repo DIR          repository root (default: this checkout)
#   --dry-run           print what would change, write nothing
#   --snapshot-proto    snapshot the PREVIOUS version's tokenization protos
#                       into proto/tokenization/v<N-1> and
#                       x/tokenization/types/v<N-1> (for state migrations
#                       that must decode the old types). Only needed when the
#                       tokenization proto types changed since the last
#                       version. Reads the protos from git ref v<N-1>
#                       (override with --proto-ref).
#   --proto-ref REF     git ref to snapshot protos from (default v<N-1>)
#   --no-proto-gen      skip `ignite generate proto-go` after the snapshot
#
# Unlike the legacy handle_upgrade_logic.sh this never renames or deletes
# app/upgrades/v<N-1>: every version's handler stays in history. Everything
# it does is idempotent; running it twice is a no-op.
set -euo pipefail
SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
# shellcheck source=lib/common.sh
. "$SCRIPT_DIR/lib/common.sh"

VERSION="" REPO="" DRY_RUN=false SNAPSHOT=false PROTO_REF="" PROTO_GEN=true
usage() { sed -n '2,19p' "$0"; exit "${1:-0}"; }
while [ $# -gt 0 ]; do
  case $1 in
    v[0-9]*) VERSION=$1; shift ;;
    --repo) REPO=$2; shift 2 ;;
    --dry-run) DRY_RUN=true; shift ;;
    --snapshot-proto) SNAPSHOT=true; shift ;;
    --proto-ref) PROTO_REF=$2; shift 2 ;;
    --no-proto-gen) PROTO_GEN=false; shift ;;
    -h|--help) usage ;;
    *) echo "unknown argument: $1" >&2; usage 1 ;;
  esac
done
[[ $VERSION =~ ^v[0-9]+$ ]] || { echo "first argument must be vNN (e.g. v35)" >&2; usage 1; }
REPO=${REPO:-$(git -C "$SCRIPT_DIR" rev-parse --show-toplevel)}
[ -f "$REPO/Makefile" ] && [ -f "$REPO/app/upgrades.go" ] || die "$REPO does not look like the chain repo"
N=${VERSION#v}
PREV="v$((N - 1))"

step "Wire $VERSION into $REPO$($DRY_RUN && echo ' (dry run)')"
PYFLAGS=(); $DRY_RUN && PYFLAGS+=(--dry-run)
python3 "$SCRIPT_DIR/lib/wire_upgrade.py" --repo "$REPO" --version "$VERSION" ${PYFLAGS[@]+"${PYFLAGS[@]}"}

if $SNAPSHOT; then
  REF=${PROTO_REF:-$PREV}
  DEST="$REPO/proto/tokenization/$PREV"
  TYPES="$REPO/x/tokenization/types/$PREV"
  step "Snapshot tokenization protos from $REF into proto/tokenization/$PREV"
  if ls "$DEST"/*.proto >/dev/null 2>&1 && [ -f "$TYPES/math.go" ]; then
    ylw "proto/tokenization/$PREV and x/tokenization/types/$PREV already exist; skipping"
  elif $DRY_RUN; then
    ylw "would snapshot proto/tokenization/*.proto from $REF into $DEST and write $TYPES/math.go"
  else
    git -C "$REPO" rev-parse --verify --quiet "$REF^{commit}" >/dev/null || die "git ref $REF not found (pass --proto-ref)"
    TMP=$(mktemp -d); trap 'rm -rf "$TMP"' EXIT
    git -C "$REPO" archive "$REF" proto/tokenization | tar -x -C "$TMP"
    rm -rf "$DEST"; mkdir -p "$DEST"
    find "$TMP/proto/tokenization" -maxdepth 1 -name '*.proto' -exec cp {} "$DEST/" \;
    rm -f "$DEST/legacytx.proto" "$DEST/query.proto"
    python3 "$SCRIPT_DIR/lib/rewrite_proto_snapshot.py" "$DEST" "$PREV"
    rm -rf "$TYPES"; mkdir -p "$TYPES"
    sed "s/__VERSION__/$PREV/" "$SCRIPT_DIR/templates/math.go.tmpl" > "$TYPES/math.go"
    MIG="$REPO/x/tokenization/keeper/migrations.go"
    if [ -f "$MIG" ]; then
      python3 - "$MIG" "$PREV" <<'PY'
import re, sys
p, v = sys.argv[1], sys.argv[2]
s = open(p, encoding="utf-8").read()
n = re.sub(r'oldtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types/v\d+"',
           f'oldtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types/{v}"', s)
if n != s:
    open(p, "w", encoding="utf-8").write(n)
    print(f"changed:   x/tokenization/keeper/migrations.go (oldtypes -> {v})")
else:
    print("unchanged: x/tokenization/keeper/migrations.go")
PY
    fi
    grn "snapshot written: $(find "$DEST" -name '*.proto' | wc -l | tr -d ' ') protos, $TYPES/math.go"
    if $PROTO_GEN; then
      if command -v ignite >/dev/null; then
        (cd "$REPO" && ignite generate proto-go -y)
      else
        ylw "ignite not found: run 'ignite generate proto-go -y' in $REPO to generate x/tokenization/types/$PREV/*.pb.go"
      fi
    fi
    ylw "reminder: bump ConsensusVersion in x/tokenization/module/module.go and register the migration if $VERSION changes tokenization state"
  fi
fi

grn "done. Next: implement app/upgrades/$VERSION/upgrades.go, then 'make upgrade-rehearsal FROM=$PREV TO=HEAD'"

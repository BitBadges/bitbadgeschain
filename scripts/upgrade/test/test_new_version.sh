#!/usr/bin/env bash
# new-version.sh: idempotent scaffolding + wiring, and a no-op (Makefile only)
# when the version is already wired by hand.
set -euo pipefail
HERE=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
NV="$HERE/../new-version.sh"
T=$(mktemp -d); trap 'rm -rf "$T"' EXIT
FAILED=0
ok()   { echo "  PASS  $1"; }
fail() { echo "  FAIL  $1"; FAILED=$((FAILED+1)); }
# assert <label> <cmd...>: PASS when the command succeeds; refute: when it fails.
assert() { local l=$1; shift; if "$@"; then ok "$l"; else fail "$l"; fi; }
refute() { local l=$1; shift; if "$@" >/dev/null 2>&1; then fail "$l"; else ok "$l"; fi; }
assert_grep()  { if grep -qF -- "$2" "$3"; then ok "$1"; else fail "$1"; fi; }
assert_ngrep() { if grep -qF -- "$2" "$3"; then fail "$1"; else ok "$1"; fi; }

# --- fixture: v35 already wired by hand, Makefile still says v34 ------------
mk_fixture() { rm -rf "${T:?}/$1"; cp -R "$HERE/testdata/$1" "$T/$1"; }
mk_fixture wired
echo "== v35 on a tree where app/upgrades/v35 is already wired"
BEFORE=$(shasum "$T/wired/app/upgrades.go" "$T/wired/app/upgrades/v35/upgrades.go")
"$NV" v35 --repo "$T/wired" >/dev/null
AFTER=$(shasum "$T/wired/app/upgrades.go" "$T/wired/app/upgrades/v35/upgrades.go")
assert_grep "Makefile VERSION bumped to v35" 'VERSION := v35' "$T/wired/Makefile"
assert "Makefile keeps its line count (no swallowed blank line)" [ "$(wc -l < "$T/wired/Makefile")" = "$(wc -l < "$HERE/testdata/wired/Makefile")" ]
if [ "$BEFORE" = "$AFTER" ]; then ok "app/upgrades.go and app/upgrades/v35 untouched"; else fail "app/upgrades.go or v35 handler were modified"; fi
assert "previous app/upgrades/v34 kept" [ -d "$T/wired/app/upgrades/v34" ]

# --- fixture: nothing for v36 exists yet -----------------------------------
mk_fixture fresh
echo "== v36 on a tree with only v34/v35"
"$NV" v36 --repo "$T/fresh" --dry-run >/dev/null
refute "--dry-run writes nothing" [ -e "$T/fresh/app/upgrades/v36" ]
assert_grep "--dry-run leaves Makefile alone" 'VERSION := v35' "$T/fresh/Makefile"

"$NV" v36 --repo "$T/fresh" >/dev/null
assert_grep "Makefile VERSION bumped to v36" 'VERSION := v36' "$T/fresh/Makefile"
assert "scaffolded app/upgrades/v36/upgrades.go" [ -f "$T/fresh/app/upgrades/v36/upgrades.go" ]
assert_grep "scaffold declares UpgradeName" 'UpgradeName = "v36"' "$T/fresh/app/upgrades/v36/upgrades.go"
assert_grep "scaffold package is v36" 'package v36' "$T/fresh/app/upgrades/v36/upgrades.go"
assert_grep "import added"      'v36 "github.com/bitbadges/bitbadgeschain/app/upgrades/v36"' "$T/fresh/app/upgrades.go"
assert_grep "handler registered" 'v36.CreateUpgradeHandler(' "$T/fresh/app/upgrades.go"
assert_grep "store-upgrade case added" 'case v36.UpgradeName:' "$T/fresh/app/upgrades.go"
assert_grep "v35 registration preserved" 'v35.CreateUpgradeHandler(' "$T/fresh/app/upgrades.go"
assert "previous app/upgrades/v35 kept" [ -d "$T/fresh/app/upgrades/v35" ]
if command -v gofmt >/dev/null; then
  BAD=$(gofmt -l "$T/fresh/app/upgrades/v36/upgrades.go")
  assert "scaffold is gofmt-clean (gofmt -l: '$BAD')" [ -z "$BAD" ]
  assert "edited app/upgrades.go parses" gofmt -e "$T/fresh/app/upgrades.go" >/dev/null 2>&1
  # The fixture mirrors the real file, which already has import-order drift;
  # the script must not add any: gofmt's diff before and after may differ only
  # by the v36 lines it inserted.
  DRIFT_BEFORE=$(gofmt -d "$HERE/testdata/fresh/app/upgrades.go" | grep -cE '^[-+][^-+]' || true)
  DRIFT_AFTER=$(gofmt -d "$T/fresh/app/upgrades.go" | grep -E '^[-+][^-+]' | grep -vc v36 || true)
  assert "no new gofmt drift in app/upgrades.go ($DRIFT_BEFORE -> $DRIFT_AFTER)" [ "$DRIFT_BEFORE" = "$DRIFT_AFTER" ]
fi

SNAP1=$(find "$T/fresh" -type f -exec shasum {} + | sort)
"$NV" v36 --repo "$T/fresh" >/dev/null
SNAP2=$(find "$T/fresh" -type f -exec shasum {} + | sort)
assert "second run is a no-op" [ "$SNAP1" = "$SNAP2" ]

# --- proto snapshot from a git ref ----------------------------------------
echo "== --snapshot-proto copies the previous tag's tokenization protos"
R="$T/protorepo"; mkdir -p "$R/proto/tokenization" "$R/x/tokenization/keeper" "$R/app/upgrades/v35"
cp -R "$HERE/testdata/fresh/." "$R/"
cat > "$R/proto/tokenization/balances.proto" <<'P'
syntax = "proto3";
package tokenization;
import "tokenization/params.proto";
option go_package = "github.com/bitbadges/bitbadgeschain/x/tokenization/types";
message Balance { string amount = 1; }
P
printf 'syntax = "proto3";\npackage tokenization;\noption go_package = "github.com/bitbadges/bitbadgeschain/x/tokenization/types";\n' > "$R/proto/tokenization/params.proto"
printf 'syntax = "proto3";\npackage tokenization;\n' > "$R/proto/tokenization/query.proto"
printf 'syntax = "proto3";\npackage tokenization;\n' > "$R/proto/tokenization/legacytx.proto"
printf 'package keeper\n\nimport (\n\toldtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types/v34"\n)\n\nvar _ = oldtypes.Balance{}\n' > "$R/x/tokenization/keeper/migrations.go"
git -C "$R" init -q && git -C "$R" add -A && git -C "$R" -c user.name=t -c user.email=t@t commit -qm init && git -C "$R" tag v35
"$NV" v36 --repo "$R" --snapshot-proto --no-proto-gen >/dev/null
assert_grep "package renamed"      'package tokenization.v35;' "$R/proto/tokenization/v35/balances.proto"
assert_grep "imports renamed"      'import "tokenization/v35/params.proto";' "$R/proto/tokenization/v35/balances.proto"
assert_grep "go_package renamed"   'x/tokenization/types/v35";' "$R/proto/tokenization/v35/balances.proto"
refute "query.proto dropped" [ -e "$R/proto/tokenization/v35/query.proto" ]
refute "legacytx.proto dropped" [ -e "$R/proto/tokenization/v35/legacytx.proto" ]
assert_grep "math.go written"      'package v35' "$R/x/tokenization/types/v35/math.go"
assert_grep "migrations.go oldtypes now v35" 'x/tokenization/types/v35"' "$R/x/tokenization/keeper/migrations.go"
assert_grep "current protos untouched" 'package tokenization;' "$R/proto/tokenization/balances.proto"

echo
if [ "$FAILED" = 0 ]; then echo "new-version: all passed"; else echo "new-version: $FAILED failed"; exit 1; fi

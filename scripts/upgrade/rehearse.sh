#!/usr/bin/env bash
# Cosmovisor upgrade rehearsal between two git refs, inside Docker.
#
#   scripts/upgrade/rehearse.sh --from v34 --to HEAD [--name v35]
#       [--multivalidator] [--evmcheck] [--rollback] [--all]
#       [--workdir DIR] [--skip-build] [--dry-run]
#
# Builds both chain binaries and cosmovisor from the given refs (never from a
# working tree: sources come from `git archive`), then runs the single-
# validator cosmovisor rehearsal plus any optional stages. Exit code 0 only if
# every check in every stage passed. --dry-run resolves refs, toolchains and
# the upgrade name, prints the plan and exits without touching Docker.
set -euo pipefail

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
# shellcheck source=lib/common.sh
. "$SCRIPT_DIR/lib/common.sh"
REPO=$(git -C "$SCRIPT_DIR" rev-parse --show-toplevel)

FROM="" TO="" NAME="" MULTI=false EVM=false ROLLBACK=false DRY_RUN=false SKIP_BUILD=false
WORKDIR=""
usage() { sed -n '2,12p' "$0"; exit "${1:-0}"; }
while [ $# -gt 0 ]; do
  case $1 in
    --from) FROM=$2; shift 2 ;;
    --to) TO=$2; shift 2 ;;
    --name) NAME=$2; shift 2 ;;
    --multivalidator) MULTI=true; shift ;;
    --evmcheck) EVM=true; shift ;;
    --rollback) ROLLBACK=true; shift ;;
    --all) MULTI=true; EVM=true; ROLLBACK=true; shift ;;
    --workdir) WORKDIR=$2; shift 2 ;;
    --skip-build) SKIP_BUILD=true; shift ;;
    --dry-run) DRY_RUN=true; shift ;;
    -h|--help) usage ;;
    *) echo "unknown flag: $1" >&2; usage 1 ;;
  esac
done
[ -n "$FROM" ] && [ -n "$TO" ] || { echo "--from and --to are required" >&2; usage 1; }

# --- resolve inputs --------------------------------------------------------
resolve() { git -C "$REPO" rev-parse --verify --quiet "$1^{commit}" || die "cannot resolve git ref '$1'"; }
label() { echo "$1" | tr '/' '-' | tr -c 'A-Za-z0-9._-\n' '_'; }
FROM_SHA=$(resolve "$FROM"); TO_SHA=$(resolve "$TO")
FROM_LABEL=$(label "$FROM"); TO_LABEL=$(label "$TO")
[ "$FROM_LABEL" != "$TO_LABEL" ] || TO_LABEL="${TO_LABEL}-to"

if [ -z "$NAME" ]; then
  NAME=$(git -C "$REPO" ls-tree --name-only "$TO_SHA" app/upgrades/ 2>/dev/null \
    | sed 's|.*/||' | grep -E '^v[0-9]+$' | sort -V | tail -1 || true)
  [ -n "$NAME" ] || die "no app/upgrades/vNN directory on $TO; pass --name"
fi
git -C "$REPO" cat-file -e "$TO_SHA:app/upgrades/$NAME/upgrades.go" 2>/dev/null \
  || die "$TO has no app/upgrades/$NAME/upgrades.go"

gomod_go() { git -C "$REPO" show "$1:go.mod" | go_version_of_gomod /dev/stdin; }
FROM_GO=$(gomod_go "$FROM_SHA"); TO_GO=$(gomod_go "$TO_SHA")
CV_GO=$(go_version_of_gomod "$SCRIPT_DIR/cosmovisor/go.mod")
GO_VERSIONS=$(printf '%s\n' "$FROM_GO" "$TO_GO" "$CV_GO" | sort -uV | tr '\n' ' ' | sed 's/ $//')

case $(uname -m) in
  arm64|aarch64) ARCH=arm64 ;;
  x86_64|amd64) ARCH=amd64 ;;
  *) die "unsupported host arch $(uname -m)" ;;
esac
HOST_MODCACHE=$(go env GOMODCACHE 2>/dev/null || echo "$HOME/go/pkg/mod")
TMP=${TMPDIR:-/tmp}
WORKDIR=${WORKDIR:-${TMP%/}/bb-upgrade-rehearsal/${FROM_LABEL}-to-${TO_LABEL}}
IMAGE="bb-upgrade-rehearsal:go-$(echo "$GO_VERSIONS" | tr ' ' '_')"
FROM_BIN=/out/bitbadgeschaind-$FROM_LABEL
TO_BIN=/out/bitbadgeschaind-$TO_LABEL

step "Plan"
cat <<PLAN
  from        : $FROM ($FROM_SHA) go$FROM_GO
  to          : $TO ($TO_SHA) go$TO_GO
  upgrade     : $NAME
  cosmovisor  : v1.7.1 via scripts/upgrade/cosmovisor (go$CV_GO)
  arch        : $ARCH
  image       : $IMAGE  (GO_VERSIONS="$GO_VERSIONS")
  module cache: $HOST_MODCACHE (mounted read-only as a file proxy)
  workdir     : $WORKDIR
  stages      : build upgrade$($ROLLBACK && echo ' rollback')$($EVM && echo ' evmcheck')$($MULTI && echo ' multivalidator')
  checks hook : $([ -f "$SCRIPT_DIR/checks/$NAME.sh" ] && echo "checks/$NAME.sh" || echo "none (generic checks only)")
PLAN
$DRY_RUN && { grn "dry run: nothing executed"; exit 0; }

command -v docker >/dev/null || die "docker is required"
[ -d "$HOST_MODCACHE/cache/download" ] || die "no module download cache at $HOST_MODCACHE; run 'go mod download' on the host first"
mkdir -p "$WORKDIR/bin" "$WORKDIR/logs" "$WORKDIR/src"

RESULTS=()
record() { RESULTS+=("$1=$2"); }

# --- build -----------------------------------------------------------------
step "Build image"
docker build --quiet -t "$IMAGE" --build-arg "GO_VERSIONS=$GO_VERSIONS" --build-arg "TARGETARCH=$ARCH" \
  "$SCRIPT_DIR" >/dev/null || die "docker build failed"
docker volume create bb-gomod >/dev/null; docker volume create bb-gobuild >/dev/null

export_src() { # <sha> <dir>
  rm -rf "$2"; mkdir -p "$2"
  git -C "$REPO" archive "$1" | tar -x -C "$2"
}
build_in_docker() { # <label> <src> <go> <out> <commit> [evmtx src]
  docker run --rm -e "EVMTX_SRC=${6:-}" -v "$SCRIPT_DIR/evmtx:/scripts/evmtx:ro" \
    -v "$HOST_MODCACHE:/hostmod:ro" -v bb-gomod:/gomodcache -v bb-gobuild:/gobuild \
    -v "$2:/src/$1:ro" -v "$WORKDIR/bin:/out" -v "$SCRIPT_DIR/container:/scripts/container:ro" \
    "$IMAGE" bash /scripts/container/build.sh "$1" "/src/$1" "$3" "$4" "$5"
}
if $SKIP_BUILD && [ -x "$WORKDIR/bin/bitbadgeschaind-$FROM_LABEL" ] && [ -x "$WORKDIR/bin/bitbadgeschaind-$TO_LABEL" ] && [ -x "$WORKDIR/bin/cosmovisor" ] && [ -x "$WORKDIR/bin/evmtx" ]; then
  ylw "--skip-build: reusing binaries in $WORKDIR/bin"
else
  step "Build $FROM_LABEL (go$FROM_GO)"
  export_src "$FROM_SHA" "$WORKDIR/src/$FROM_LABEL"
  build_in_docker "$FROM_LABEL" "$WORKDIR/src/$FROM_LABEL" "$FROM_GO" "$FROM_BIN" "$FROM_SHA" || die "build of $FROM failed"
  step "Build $TO_LABEL (go$TO_GO)"
  export_src "$TO_SHA" "$WORKDIR/src/$TO_LABEL"
  build_in_docker "$TO_LABEL" "$WORKDIR/src/$TO_LABEL" "$TO_GO" "$TO_BIN" "$TO_SHA" /scripts/evmtx || die "build of $TO failed"
  step "Build cosmovisor (go$CV_GO)"
  build_in_docker cosmovisor "$SCRIPT_DIR/cosmovisor" "$CV_GO" /out/cosmovisor "$(git -C "$REPO" rev-parse HEAD)" || die "build of cosmovisor failed"
fi
docker run --rm -v "$WORKDIR/bin:/out:ro" -v "$SCRIPT_DIR/container:/scripts/container:ro" \
  "$IMAGE" bash /scripts/container/verify-binaries.sh "$FROM_BIN" "$FROM_SHA" "$TO_BIN" "$TO_SHA" \
  || die "chain binary identity check failed"
record build pass

# --- stages ----------------------------------------------------------------
run_stage() { # <name> <script>
  step "Stage: $1"
  local rc
  set +e
  docker run --rm \
    -e FROM_BIN="$FROM_BIN" -e TO_BIN="$TO_BIN" -e CV=/out/cosmovisor \
    -e UPGRADE_NAME="$NAME" -e CHECKS_DIR=/scripts/checks -e LOG_DIR=/logs -e EVMTX=/out/evmtx \
    -v "$WORKDIR/bin:/out:ro" -v "$SCRIPT_DIR:/scripts:ro" -v "$WORKDIR/logs:/logs" \
    "$IMAGE" bash "/scripts/container/$2" 2>&1 | tee "$WORKDIR/logs/stage-$1.txt"
  rc=${PIPESTATUS[0]}
  set -e
  if [ "$rc" = 0 ]; then record "$1" pass; else record "$1" "fail($rc)"; fi
}
run_stage upgrade upgrade.sh
$ROLLBACK && run_stage rollback rollback.sh
$EVM && run_stage evmcheck evmcheck.sh
$MULTI && run_stage multivalidator multivalidator.sh

step "Summary ($FROM -> $TO, upgrade $NAME)"
EXIT=0
for r in "${RESULTS[@]}"; do
  case $r in
    *=pass) grn "  PASS  ${r%=*}" ;;
    *) red "  FAIL  ${r%=*} (${r#*=})"; EXIT=1 ;;
  esac
done
echo "logs: $WORKDIR/logs"
exit $EXIT

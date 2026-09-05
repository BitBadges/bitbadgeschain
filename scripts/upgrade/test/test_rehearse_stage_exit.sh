#!/usr/bin/env bash
# A failing stage must show as FAIL in the summary and make rehearse.sh exit 1,
# and a passing one as PASS with exit 0. Uses a fake `docker` on PATH whose
# `run` exits with $FAKE_DOCKER_RC, so no image or container is involved.
set -uo pipefail
HERE=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
RH="$HERE/../rehearse.sh"
FAILED=0
ok()   { echo "  PASS  $1"; }
fail() { echo "  FAIL  $1"; FAILED=$((FAILED+1)); }
assert() { local l=$1; shift; if "$@"; then ok "$l"; else fail "$l"; fi; }

T=$(mktemp -d); trap 'rm -rf "$T"' EXIT
mkdir -p "$T/bin" "$T/modcache/cache/download" "$T/work/bin"
cat > "$T/bin/docker" <<'SH'
#!/usr/bin/env bash
case $1 in
  run) echo "fake stage output"; exit "${FAKE_DOCKER_RC:-0}" ;;
  *) exit 0 ;;
esac
SH
chmod +x "$T/bin/docker"
# Fake `go` too, so the host toolchain is never invoked and the module cache
# rehearse.sh mounts is the empty one under $T.
printf '#!/usr/bin/env bash\necho "%s"\n' "$T/modcache" > "$T/bin/go"; chmod +x "$T/bin/go"
for b in bitbadgeschaind-v34 bitbadgeschaind-HEAD cosmovisor; do printf '#!/bin/sh\n' > "$T/work/bin/$b"; chmod +x "$T/work/bin/$b"; done

run() { PATH="$T/bin:$PATH" FAKE_DOCKER_RC=$1 \
  "$RH" --from v34 --to HEAD --workdir "$T/work" --skip-build --rollback 2>&1; }

OUT=$(run 3); RC=$?
assert "exit 1 when a stage fails"          [ "$RC" = 1 ]
assert "summary marks upgrade as FAIL(3)"   grep -q 'FAIL  upgrade (fail(3))' <<<"$OUT"
assert "summary marks rollback as FAIL(3)"  grep -q 'FAIL  rollback (fail(3))' <<<"$OUT"
assert "build reused (--skip-build)"        grep -q 'reusing binaries' <<<"$OUT"

OUT=$(run 0); RC=$?
assert "exit 0 when every stage passes"     [ "$RC" = 0 ]
assert "summary marks upgrade as PASS"      grep -q 'PASS  upgrade' <<<"$OUT"
assert "stage output is tee'd to logs"      grep -q 'fake stage output' "$T/work/logs/stage-upgrade.txt"

echo
if [ "$FAILED" = 0 ]; then echo "rehearse stage exit: all passed"; else echo "rehearse stage exit: $FAILED failed"; exit 1; fi

#!/usr/bin/env bash
# Test runner for the veg port. Iterates without resetting the tree.
#
# Usage:
#   ./port/tests.sh            - runnable Go unit tests
#   ./port/tests.sh <pkg...>   - specific packages, e.g. ./port/tests.sh ./pkg/structural/...
#   ./port/tests.sh flow       - self-hosting: CUE flow tests from test.cue via `veg flow`
set -uo pipefail

cd ~/verdverm/veg

# packages whose tests cannot compile in the port (known, pre-existing):
#   pkg/gotils/lockedfile/...      - imports Go stdlib internal/testenv
#   pkg/repos/utils                - test references undefined parseModURL
#   pkg/yagu/configdir             - test references undefined reset
# script-harness suites (flow/dm/fmt/cuecmd) are excluded pending testscript migration.
RUNNABLE=(
  ./pkg/structural/...
  ./pkg/diff3/...
  ./pkg/gotils/par/...
  ./pkg/gotils/txtar/...
  ./pkg/create/...
  ./pkg/repos/git/...
  ./pkg/repos/cache/...
  ./pkg/env/incept/...
  ./pkg/datamodel/...
  ./pkg/cuetils/...
  ./pkg/templates/...
  ./pkg/watch/...
  ./pkg/yagu/...
  ./pkg/singletons/...
)

# known environment-dependent (need live registry probing):
#   ./pkg/repos/remote/...   ./pkg/repos/oci/...

if [ $# -gt 0 ]; then
  if [ "$1" = "flow" ]; then
    shift
    TARGETS=${@:-test/st}
    for t in $TARGETS; do
      echo "== veg flow $t"
      veg flow "$t" || echo "FAIL: veg flow $t"
    done
    exit $?
  fi
  RUNNABLE=("$@")
fi

FAILS=0

# known issue: flags.arch (@tag(arch,var=arch)) has no default and is not
# auto-injected; dist targets need it concrete. Derive from the machine,
# without assuming go is installed.
case "$(uname -m)" in
  arm64|aarch64) ARCH=arm64 ;;
  x86_64)        ARCH=amd64 ;;
  *)             ARCH="" ;;
esac
SYNC_ARGS=()
if [ -n "$ARCH" ]; then
  SYNC_ARGS=(-t "arch=$ARCH")
fi

echo "== veg env sync dist-"
if ! veg env sync dist- ${SYNC_ARGS[@]+"${SYNC_ARGS[@]}"}; then
  FAILS=$((FAILS+1))
fi

# regression: builder container include list must use the veg layout
# (cmd, pkg), else the dagger inception reports missing go packages.
echo "== veg env sync veg-cli-darwin"
if ! veg env sync veg-cli-darwin ${SYNC_ARGS[@]+"${SYNC_ARGS[@]}"}; then
  FAILS=$((FAILS+1))
fi

for p in "${RUNNABLE[@]}"; do
  echo "== go test $p"
  if ! go test "$p"; then
    FAILS=$((FAILS+1))
  fi
done

echo
if [ $FAILS -gt 0 ]; then
  echo "$FAILS check(s) FAILED"
  exit 1
fi
echo "ALL TESTS PASSED"

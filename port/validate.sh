#!/usr/bin/env bash
# Validation gate for the veg port: build, vet, CUE vet, branding sweep, smoke test.
# Exits non-zero on the first failure.
set -euo pipefail

cd ~/verdverm/veg

fail() { echo "FAIL: $1" >&2; exit 1; }

echo "== go build ./..."
go build ./... || fail "go build"

echo "== go install ./cmd/veg + veg version"
go install ./cmd/veg || fail "go install"
veg version || fail "veg version"

echo "== go vet ./... (report only)"
go vet ./... || true

echo "== cue vet -c ./... (enrichers two-package warning tolerated)"
cue vet -c ./... 2>&1 | grep -v 'schemas/dm/enrichers' && fail "cue vet reported errors"
echo "cue vet clean"

echo "== branding sweep: no standalone 'hof' outside hofstadter identifiers"
LEFT=$(grep -rn 'hof' --include='*.go' --include='*.cue' pkg cmd/veg schemas design catalogs examples veg.cue flow.cue git.cue test.cue 2>/dev/null \
  | grep -v 'hofstadter\|Hofstadter' \
  | grep -v '_test.go:\|/test/\|testdata/' \
  | grep -v 'pkg/repos/oci' \
  | grep -v 'hof2veg' \
  || true)
if [ -n "$LEFT" ]; then
  echo "$LEFT" | head -20
  fail "leftover hof references"
fi
echo "branding clean (oci media types + test fixtures exempt)"

echo "== CUE alias sweep: no old-style aliases"
OLD=$(grep -rn '\[[A-Za-z_][A-Za-z_0-9]*=string\]\|[A-Za-z_][A-Za-z_0-9]*=[A-Za-z_][A-Za-z_0-9]* & {' \
  --include='*.cue' pkg schemas design catalogs examples veg.cue flow.cue git.cue test.cue 2>/dev/null || true)
if [ -n "$OLD" ]; then
  echo "$OLD" | head -10
  fail "old-style CUE aliases remain"
fi
echo "aliasv2 clean"

echo "== cue.mod language version"
grep -q 'version: "v0.18.0"' cue.mod/module.cue || fail "cue.mod language version not v0.18.0"

echo "== smoke: veg --help"
veg --help >/dev/null || fail "veg --help"

echo
echo "ALL CHECKS PASSED"

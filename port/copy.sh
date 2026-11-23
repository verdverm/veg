#!/usr/bin/env bash
# Mechanical copy-rename: ~/hof/hof -> ~/verdverm/veg
# Idempotent: wipes output dirs it owns, then re-copies. Hand-written files
# outside these outputs are never touched. Patches applied separately.
set -euo pipefail

SRC=~/hof/hof
DST=~/verdverm/veg

cd "$DST"

# ---------- outputs this script owns (wiped on each run) ----------
rm -rf pkg cmd/veg schemas design catalogs examples ci
rm -rf cmd/veg/verinfo  # replaced by overlay pkg/verinfo (stdlib)
rm -f veg.cue flow.cue git.cue test.cue envs.cue

# ---------- excluded from source ----------
# dirs: website/docs (not copied at all), tcell TUI, connector, script harness, ga
EXCL=(
  --exclude '.git'
  --exclude 'node_modules'
  --exclude '.hof'
  --exclude '.veg'
  --exclude 'docs'
  --exclude 'docs-old'
  --exclude 'script'
  --exclude 'hack'
  --exclude 'images'
  --exclude 'notes'
  --exclude 'repro'
  --exclude '/tui'
  --exclude '/connector'
  --exclude '/test'
  --exclude '/packer'
  --exclude '/extern'
  --exclude '/datautils/diff'
  --exclude '/datautils/visit'
  --exclude '/flow_test.go'
  --exclude '/fmt/test/00_test.go'
  --exclude '/cuecmd/cli_test.go'
  --exclude '/datamodel/test/testscript_test.go'
  --exclude '/cmd/hof/ga'
  --exclude 'goreleaser*'
  --exclude 'orig-gorel.yml'
  --exclude 'package.json'
  --exclude 'pnpm-*'
)

# ---------- path mapping ----------
# lib/...       -> pkg/...
# lib/hof/      -> pkg/veg/
# flow/         -> pkg/flow/
# cmd/hof/      -> cmd/veg/
rsync -a "${EXCL[@]}" \
  --exclude 'cmd/hof' \
  "$SRC/lib/" "$DST/pkg/"

rsync -a --exclude '.git' --exclude 'node_modules' \
  --exclude 'tui' --exclude 'ga' \
  --exclude '/cmd/run.go' --exclude '/cmd/tui.go' --exclude '/flags/run.go' \
  "$SRC/cmd/hof/" "$DST/cmd/veg/"

rsync -a "${EXCL[@]}" --exclude '/cmd' \
  "$SRC/flow/" "$DST/pkg/flow/"
rsync -a "${EXCL[@]}" "$SRC/flow/cmd/" "$DST/pkg/flow/cmd/"

rsync -a "${EXCL[@]}" "$SRC/schemas/" "$DST/schemas/"
rsync -a "${EXCL[@]}" "$SRC/design/" "$DST/design/"
rsync -a "${EXCL[@]}" "$SRC/catalogs/" "$DST/catalogs/"
rsync -a "${EXCL[@]}" "$SRC/examples/" "$DST/examples/"
# ci/ excluded: repo workflow/packer tooling, not part of the CLI

rsync -a "$SRC/hof.cue" "$DST/veg.cue"

# root pnpm workspace files (drive extensions/vscode build via `pnpm vscode:build`)
rsync -a "$SRC/package.json" "$SRC/pnpm-workspace.yaml" "$SRC/pnpm-lock.yaml" "$DST/"

rm -rf .veg-data.tmp
if [ -d .veg/data ]; then mv .veg/data .veg-data.tmp; fi
rm -rf .veg
rsync -a "$SRC/.veg/" .veg/
[ -d .veg-data.tmp ] && mv .veg-data.tmp .veg/data
rm -f port/patches/.rejects

# lib/hof -> pkg/veg (dir rename; hof.go -> veg.go)
rm -rf "$DST/pkg/veg"
mv "$DST/pkg/hof" "$DST/pkg/veg"
mv "$DST/pkg/veg/hof.go" "$DST/pkg/veg/veg.go"

# flow/tasks/hof (package hof) -> tasks/veg (package veg)
rm -rf "$DST/pkg/flow/tasks/veg"
mv "$DST/pkg/flow/tasks/hof" "$DST/pkg/flow/tasks/veg"
mv "$DST/pkg/flow/tasks/veg/hof.go" "$DST/pkg/flow/tasks/veg/veg.go" 2>/dev/null || true
rsync -a "$SRC/flow.cue" "$DST/flow.cue"
rsync -a "$SRC/git.cue" "$DST/git.cue"
rsync -a "$SRC/test.cue" "$DST/test.cue"
rsync -a "$SRC/veg.cue" "$DST/envs.cue"

rsync -a --delete \
  --exclude 'node_modules' --exclude 'out' --exclude '.vscode-test' \
  "$SRC/extensions/" "$DST/extensions/"
# out/ is build output, rebuilt via pnpm; node_modules is installed locally

# ---------- overlay: hand-written files owned by the port ----------
rsync -a "$DST/port/overlay/" "$DST/"

# ---------- patches: hand edits, applied in order ----------
shopt -s nullglob
for p in "$DST"/port/patches/*.patch; do
  # exit 1 = hunks already applied (benign on re-runs, --forward skips), 2 = real failure
  patch -s -p1 --forward -r "$DST/port/patches/.rejects" -d "$DST" < "$p" || [ $? -eq 1 ]
done
find "$DST" -name '*.orig' -not -path '*/port/patches/orig/*' -delete

echo "copy done"

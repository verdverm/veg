#!/usr/bin/env bash
# Rewrite phase: module path, CUE imports, #hof->#veg, Hof->Veg, .hof paths, GA removal.
# Idempotent: all substitutions are rename-shaped and stable on re-run.
set -euo pipefail

DST=~/verdverm/veg
cd "$DST"

TARGETS=$(find pkg cmd/veg schemas design catalogs examples .veg veg.cue envs.cue flow.cue git.cue test.cue \
  -type f \( -name '*.go' -o -name '*.cue' -o -name '*.md' -o -name '*.yml' -o -name '*.yaml' -o -name '*.tmpl' -o -name '*.txt' -o -name '*.mod' \) 2>/dev/null)

# ---------- Go / CUE import paths (ordered: specific -> generic) ----------
for f in $TARGETS; do
  sed -i '' \
    -e 's|github.com/hofstadter-io/hof/lib/hof|github.com/verdverm/veg/pkg/veg|g' \
    -e 's|github.com/hofstadter-io/hof/lib/|github.com/verdverm/veg/pkg/|g' \
    -e 's|github.com/hofstadter-io/hof/flow/|github.com/verdverm/veg/pkg/flow/|g' \
    -e 's|github.com/hofstadter-io/hof/cmd/hof|github.com/verdverm/veg/cmd/veg|g' \
    -e 's|github.com/hofstadter-io/hof/schemas|github.com/verdverm/veg/schemas|g' \
    -e 's|github.com/hofstadter-io/hof|github.com/verdverm/veg|g' \
    -e 's|github.com/hofstadter-io/schemas|github.com/verdverm/veg/schemas|g' \
    -e 's|github.com/hofstadter-io/hofmod-cli|github.com/verdverm/veg/catalogs/gen/cli|g' \
    "$f"
done

# ---------- Go identifiers (ordered: longer first) ----------
GOFILES=$(find pkg cmd/veg -type f -name '*.go')
for f in $GOFILES; do
  sed -i '' \
    -e 's|#hof|#veg|g' \
    -e 's|FindHofsOrig|FindVegsOrig|g' \
    -e 's|FindHofs|FindVegs|g' \
    -e 's|ParseHof|ParseVeg|g' \
    -e 's|HofsOrig|VegsOrig|g' \
    -e 's|Hofs|Vegs|g' \
    -e 's|[[:<:]]Hof[[:>:]]|Veg|g' \
    "$f"
done

# undo sed \b artifact on Hofstadter (word boundary mis-match)
sed -i '' 's|Vegstadter|Hofstadter|g' pkg/repos/oci/oci.go

# package qualifier: pkg/veg package name hof -> veg
for f in $GOFILES; do
  sed -i '' \
    -e 's|^[[:space:]]*package hof[[:space:]]*$|package veg|' \
    -e 's|[[:<:]]hof[.]Veg[[:>:]]|veg.Veg|g' \
    -e 's|[[:<:]]hof\(\.[A-Z]\)|veg\1|g' \
    "$f"
done

# BSD sed [[:<:]] did not apply cleanly earlier runs are idempotent; re-fix leftovers
for f in $GOFILES; do
  sed -i '' \
    -e 's|Vegstadter|Hofstadter|g' \
    -e 's|\.[H]of\.|.Veg.|g' \
    -e 's|u\.T\.Hof|u.T.Veg|g' \
    "$f"
done

# ---------- CUE: module paths + identifiers ----------
CUEFILES=$(find pkg cmd/veg schemas design catalogs examples .veg veg.cue envs.cue flow.cue git.cue test.cue -type f -name '*.cue' 2>/dev/null)
for f in $CUEFILES; do
  sed -i '' \
    -e 's|github.com/verdverm/vegmod-cli|github.com/verdverm/veg/catalogs/gen/cli|g' \
    -e 's|schema\.Hof|schemas.#veg|g' \
    -e 's|schema\.DHof|schemas.#veg|g' \
    -e 's|ghcr.io/hofstadter-io|ghcr.io/verdverm|g' \
    -e 's|docs.hofstadter.io|github.com/verdverm/veg/wiki|g' \
    -e 's|https://hofstadter.io|https://github.com/verdverm/veg|g' \
    -e 's|ghcr.io/hofstadter-io/hof|ghcr.io/verdverm/veg|g' \
    -e 's|Hofstadter, Inc|verdverm|g' \
    -e 's|open-source@hofstadter.io|verdverm@users.noreply.github.com|g' \
    -e 's|Owner: "hofstadter-io"|Owner: "verdverm"|g' \
    -e 's|Repo:       "hofstadter"|Repo:       "veg"|g' \
    -e 's|github.com/hofstadter-io/demos|github.com/verdverm/demos|g' \
    -e 's|[[:<:]]hof[[:>:]]|veg|g' \
    "$f"
done

# creators/templates keeps the published module name (consumed remotely by generated modules)
for f in $(find catalogs/gen/cli/creators/templates -type f -name '*.cue'); do
  sed -i '' 's|github.com/verdverm/veg/catalogs/gen/cli|github.com/verdverm/vegmod-cli|g' "$f"
done

# ---------- string path references: lib/ -> pkg/ (dir paths, globs, tags) ----------
for f in $CUEFILES; do
  sed -i '' 's|"lib/|"pkg/|g' "$f"
done

# ---------- CUE: aliasv2 pattern aliases [K=T]: -> [T]~(K, _): ----------
for f in $CUEFILES; do
  sed -i '' \
    -e 's|^\([A-Z][A-Z_0-9]*\)=\([a-z][A-Za-z_0-9]*\): |\2~\1: |' \
    -e 's|\[\([A-Za-z_][A-Za-z_0-9]*\)=\(string\)\]|[\2]~(\1, _)|g' \
    -e 's|\[\([A-Za-z_][A-Za-z_0-9]*\)= \(!="[^"*]*"\)\]|[\2]~(\1, _)|g' \
    "$f"
done

# ---------- CUE metadata: #hof -> #veg ----------
for f in $CUEFILES; do
  sed -i '' -e 's|#hof|#veg|g' "$f"
done

# strip ga imports from generated cmd files (package is dropped)
for f in $(find cmd/veg -type f -name '*.go'); do
  sed -i '' \
    -e '/"github.com\/verdverm\/veg\/cmd\/veg\/ga"/d' \
    "$f"
done

# drop run/tui commands (decided keep-list): remove their wiring from generated root.go
sed -i '' -e '/RootCmd.AddCommand(RunCmd)/d' -e '/RootCmd.AddCommand(TuiCmd)/d' cmd/veg/cmd/root.go

# verinfo: replaced by pkg/verinfo (stdlib runtime/debug)
for f in $(find pkg cmd/veg -type f -name '*.go'); do
  sed -i '' 's|"github.com/verdverm/veg/cmd/veg/verinfo"|"github.com/verdverm/veg/pkg/verinfo"|g' "$f"
done

# verinfo: HofVersion/HofSchemaVersion -> Veg names
for f in $(find pkg cmd/veg -type f -name '*.go'); do
  sed -i '' -e 's|verinfo.HofVersion|verinfo.VegVersion|g' -e 's|verinfo.HofSchemaVersion|verinfo.VegSchemaVersion|g' "$f"
done

# ---------- user-facing rebrand: binary name in cmd strings ----------
for f in $(find cmd/veg -type f -name '*.go'); do
  sed -i '' \
    -e 's|[[:<:]]hofLong[[:>:]]|vegLong|g' \
    -e 's|[[:<:]]hof[[:>:]]|veg|g' \
    "$f"
done

# ---------- user-facing rebrand: author/homepage/update URLs ----------
for f in $(find cmd/veg -type f -name '*.go'); do
  sed -i '' \
    -e 's|Author:      Hofstadter, Inc|Author:      verdverm|' \
    -e 's|https://hofstadter.io|https://github.com/verdverm/veg|' \
    -e 's|api.github.com/repos/hofstadter-io/veg/|api.github.com/repos/verdverm/veg/|' \
    -e 's|docs.hofstadter.io|github.com/verdverm/veg/wiki|' \
    "$f"
done

# ---------- rebrand leftovers: identifiers, strings, comments ----------
# catalogs/gen/cuelm Go tooling files too
GOFILES_ALL="$GOFILES $(find catalogs/gen/cuelm -type f -name '*.go' 2>/dev/null)"
for f in $GOFILES_ALL; do
  sed -i '' \
    -e 's|hofcontext|vegcontext|g' \
    -e 's|hofruntime|vegruntime|g' \
    -e 's|HofRuntime|VegRuntime|g' \
    -e 's|hofDir|vegDir|g' \
    -e 's|\.hofmod-ignore|.vegmod-ignore|g' \
    -e 's|\$hof|$veg|g' \
    -e 's|"hof", "mod", "tidy"|"veg", "mod", "tidy"|g' \
    -e 's|"hof mod tidy"|"veg mod tidy"|g' \
    -e 's|"hof mod vendor"|"veg mod vendor"|g' \
    -e 's|hof-create-|veg-create-|g' \
    -e 's|hof-stdout-|veg-stdout-|g' \
    -e 's|"hof env |"veg env |g' \
    -e 's|changes to hof|changes to veg|g' \
    -e 's|HOF_VERSION|VEG_VERSION|g' \
    -e 's|hof_v|veg_v|g' \
    -e 's|[[:<:]]hofp[[:>:]]|vegp|g' \
    -e 's|[[:<:]]hof[[:>:]]|veg|g' \
    -e 's|[[:<:]]hof[[:>:]]|veg|g' \
    "$f"
done

for f in $CUEFILES; do
  sed -i '' -e 's|\$hof|$veg|g' -e 's|[[:<:]]hofmods[[:>:]]|vegmods|g' "$f"
done

# .veg files: package veg (half-done rebrand in source) collides with root package veg -> import cycle
for f in $(find .veg -type f -name '*.cue'); do
  sed -i '' 's|^package veg$|package config|' "$f"
done
# envs.cue (hof's root veg.cue): re-point the link at the renamed package
sed -i '' \
  -e 's|github.com/verdverm/veg/.veg:veg|github.com/verdverm/veg/.veg:config|' \
  -e 's|^veg$|veg: config|' \
  envs.cue

# ---------- CUE fix: value aliases (F=_flow) -> let clauses via cue fix ----------
# cue fix parses at the module's language version, so drop to v0.15.1 for the run
sed -i '' 's|language: version: "v0.18.0"|language: version: "v0.15.1"|' cue.mod/module.cue
cue fix . --exp=aliasv2 -f >/dev/null 2>&1 || true
sed -i '' 's|language: version: "v0.15.1"|language: version: "v0.18.0"|' cue.mod/module.cue

# ---------- targeted string fixes (idempotent) ----------
sed -i '' 's|Mkdir(".hof")|Mkdir(".veg")|; s|TempFile(".hof"|TempFile(".veg"|' pkg/yagu/exec.go
sed -i '' 's|ContainerPrefix = "hof-fmt-"|ContainerPrefix = "veg-fmt-"|' pkg/fmt/fmtrs.go
sed -i '' 's|module := "hof.io"|module := "veg.io"|; s|"hof mod tidy"|"veg mod tidy"|' pkg/gen/cmd/init.go
sed -i '' 's|filepath.Join(dir, ".hof", "dm")|filepath.Join(dir, ".veg", "dm")|' pkg/datamodel/history.go
sed -i '' 's|cacheBaseDir = filepath.Join(d, "hof")|cacheBaseDir = filepath.Join(d, "veg")|' pkg/repos/cache/load.go
sed -i '' 's|hofDir             = "hof"|hofDir             = "veg"|' pkg/repos/remote/mirrors.go
sed -i '' 's|case "hof":|case "veg":|' pkg/veg/find.go

# lib/ root-level files (lib/cmd.go, lib/feedback.go): package lib -> package pkg
sed -i '' 's|^package lib$|package pkg|' pkg/cmd.go pkg/feedback.go
sed -i '' 's|"github.com/verdverm/veg/lib"|"github.com/verdverm/veg/pkg"|; s|lib\.SendFeedback|pkg.SendFeedback|g' cmd/veg/cmd/feedback.go

# ---------- .hof paths -> .veg ----------
for f in $TARGETS; do
  sed -i '' \
    -e 's|\.hof/shadow|.veg/shadow|g' \
    -e 's|\.hof/dm|.veg/dm|g' \
    "$f"
done

# ---------- Google Analytics: strip SendCommandPath lines ----------
GAFILES=$(find cmd/veg -type f -name '*.go')
for f in $GAFILES; do
  sed -i '' -e '/ga\.SendCommandPath/d' "$f"
done

echo "rewrite done"

# AGENTS.md

This is the `veg` repo. (~/verdverm/veg)

The current major task is renaming and cleaning up `hof` into this new `veg` project and tool.

- `hof`: ~/hof/hof
- `veg`: ~/verdverm/veg

The core is a Go code base. The experience is centered around CUE. It is a big project built up over nearly a decade.

- `#veg` (formerly `#hof`) special metadata fields that align to `veg` features
- Go's `text/template` for code generation
- CUE's `cue/flow` for workflows
- Dagger cue-sdk for containers
- ADK for LLM and agentic
- TUI coding agent

Long-form subsystem reference: NOTES.md (source map, internals, gotchas, decisions, open questions).

## Status

- exploration and planning done (see NOTES.md)
- port not started yet

## Vision

- no concerns about backwards incompatibility, consider this greenfield
- two distinct phases, do NOT conflate them:
    1. mechanical rename/port (bash scripts using cp and go tooling) - NO dependency changes (version or replace), with TWO exceptions: CUE bumps to latest `v0.18.0-rc.3` (get out in front; newly caught silent errors surface during the port) and Go updates to 1.27.4; Dagger and ADK stay pinned
    2. after the rename is done: remaining dependency updates (Dagger, ADK), config revamp, cleanup
- `go fix` should be core to the port (mechanical Go API migration); `cue fix` is the CUE equivalent (old syntax/API rewrites, especially after the CUE bump)

## Important guidance

- scripts and tooling (bash, cp, go fix, cue fix, go tooling) do the vast majority of the port - hand-editing is the exception, not the method
- all files associated with the porting effort (scripts, manifests, notes) live under `./port/...` - this directory is temporary and will eventually be deleted
- the copy-rename phase must be resettable: able to wipe its output and rerun after fixing the porting scripts - keep the mechanical pass idempotent and its outputs clearly delimited so a reset never touches hand-written work
- the full port cannot be done mechanically alone: the agent also edits files directly (API changes, refactors, fixups the scripts cannot express)
- every hand edit is captured as a patch under `port/patches/...`; patches are a permanent, first-class part of the pipeline, applied in order as the last step of every run: reset = wipe outputs -> run scripts -> apply patches
- expect MOST patches to survive to the end of the port - that is the normal case, not a backlog to burn down; when the port completes, the final run's output (scripts + patches) is the tree and `port/` is deleted
- occasionally a fixup turns out to be a trivial deterministic transformation (e.g. a missed string rename) - it may be moved into the scripts and its patch dropped

## Path mapping

| Source (~/hof/hof) | Dest (~/verdverm/veg) |
|---|---|
| `lib/...` | `pkg/...` (no "go" path segment) |
| `lib/hof/` | `pkg/veg/` (`#hof` -> `#veg`, Go `Hof` -> `Veg`) |
| `flow/` | `pkg/flow/` |
| `cmd/hof/` | `cmd/veg/` |
| `{catalogs,design,examples,schemas}` | same names |
| `catalogs/gen/cli/` (in-tree hofmod-cli copy, external import paths) | `catalogs/gen/cli/` - rewrite imports to veg module path |

Module path: `github.com/verdverm/veg` (go.mod and cue.mod already initialized).
Per-package keep/drop details: NOTES.md section 0 and 8.

## Decisions

- `.veg` is the single config dir; codegen shadow lives at `.veg/shadow`; `.hof` disappears
- `lib/prompt` stays - it is survey-based questionnaires (NOT LLM), needed by `create` and the flow prompt task; legacy LLM code is already gone from hof (all live LLM is ADK under `lib/agent`)
- CLI codegen: hand-port the generated cmd files for now; the hofmod-cli copy already exists in hof at `catalogs/gen/cli` - port it and rewrite its imports so self-codegen works again
- CLI keep-list: keep feedback and update; drop `run` (script-based) and `tui` (tcell) commands; charm TUI remains reachable via agent commands until lifted
- schemas: repo-root `schemas/` is the ONLY source of truth; the external `hofstadter-io/schemas` CUE module goes away (imports rewrite to root schemas); eventually published as a CUE module from a naked repo `verdverm/schema`
- ADK: keep the local `../adk` replace for the first pass; investigate its patches against upstream v1+ later (patches likely still needed)
- ADK auth: the tool must NOT manage Google/Gemini auth; use the OpenAI-compatible model path pointed at a user-run gateway that holds credentials
- Dagger: port as-is for the first pass (`incept` + distconsts replace included); reassess against modern Dagger later
- `lib/container` + `lib/fmt` formatters: keep, core subsystems
- flow tasks: drop IRC; keep watermill (important for pkg/flow workflows) and the rest (`db.*`, etc.)
- datamodel subsystem: keep; `.hof/dm` -> `.veg/dm`
- Google Analytics (`cmd/hof/ga`, `ga.SendCommandPath` in every generated command, `catalogs/gen/cli/templates/ga.go`): remove entirely
- `.veg` internal layout: keep flat - CUE and users can nest as they like; reserved dirs: `.veg/{shadow,dm,data}` (`data` = agent sqlite DBs + artifacts; global DB default is `UserConfigDir/veg/data/veg.db`, overridable via `VEG_DATABASE_CONN`/`VEG_USER_DATA_DIR`)
- `lib/connector`: drop (verified: only used by the excluded tcell TUI)
- `lib/cuetils`: keep the useful helpers, remove the load-logic duplication with the core runtime
- `lib/repos`: keep for `veg create` (also used by `lib/fmt`); `mod` is already a `cue mod` proxy and does not use repos
- `lib/database`: readme-only; aspirational unified data layer across datastore/datamodel/flow - carry the readme as a design note
- config system revamp (phase 2): merging as close to CUE semantics as possible, with system-level config that can override / not be updated by lower levels
- TUI: charm only (`lib/agent/cmd/tui`); `veg tui` currently wires the excluded tcell TUI and must be rewired; eventually lift the charm TUI up from its agent focus
- tests: migrate the hof/script-based tests to `github.com/rogpeppe/go-internal` testscript rather than dropping them outright

## Exclude

- most anything unrelated to the core hof/veg binary (website, docs, github)
- tcell TUI (`lib/tui`, vendored tview fork)
- `~/hof/hof/script` (test harness itself; its tests migrate to go-internal testscript)
- IRC flow task
- `run` and `tui` CLI commands (tui returns when charm is lifted)
- Google Analytics tracking

## Gotchas

- go.mod local replaces to preserve during the rename: `google.golang.org/adk => ../adk` and `github.com/dagger/dagger/engine/distconsts => ./lib/env/incept/distconsts`; the patched ADK checkout exists at `~/hof/adk` but NOT at `~/verdverm/adk` - it must be placed beside veg for the build to work
- all CUE imports use the Go module path - the module rename forces a sweep of every CUE import (schemas, catalogs, design, examples, .veg)
- half-done rebrand already in hof: `VEG_*` env vars, `// [[ hof2veg ]]` markers, `.veg` CUE files with `package veg`
- two distinct "schemas": repo-local `schemas/` dir vs external CUE module `github.com/hofstadter-io/schemas`
- pinned versions during rename: Dagger SDK v0.20.0, ADK v0.3.0; known drift: dagger engine image const v0.19.10 vs SDK v0.20.x
- CUE version: hof pins v0.15.4 (`language.version: v0.15.1`); target is latest `v0.18.0-rc.3` (language v0.18.0, matches veg cue.mod) - local binary verified at rc.3; run `cue fix` after the bump; expect previously silent errors to surface
- Go toolchain: local go is 1.27.4 (verified); veg go.mod still says `go 1.23.0`, hof pinned 1.24.4 - bump the go directive during the port; modern `go fix` is available
- validate loop: `go install ./cmd/veg` + `veg version` only, unless instructed otherwise



---


The problem, plainly:

1. What broke: The port decision bumped CUE from v0.15 to v0.18. In CUE v0.18.0 the old alias syntax (`F=_flow & {...}, [N=string]: {...}`) became a hard parse error — the new "aliasv2" syntax uses postfix `~` instead. About 26 sites across `flow.cue`, `git.cue`, `schemas/`, `catalogs/gen/cuelm/`, `pkg/datamodel/test/` use the old form.
2. Why `cue fix` can't save us: The tool refuses to rewrite these — it can't parse files containing the old syntax, so it skips them and just reports the errors. Chicken-and-egg.
3. What I was doing in /tmp: I don't know the exact `aliasv2` postfix spelling, and my guesses (`expr~F, [string~N], {[K]~(V)}`...) all fail to parse. So I was probing the parser with one-liners to learn the correct form before touching the repo — sloppy-looking, agreed, but it was to avoid hand-editing 26 sites with a wrong syntax.
4. What I actually know works: `{(_flow)~F}` inside braces parses (verified). The bare top-level form and the `[K=V]` pattern-alias form I have not found the correct spelling for yet.

The correct move is to get the aliasv2 grammar from the CUE spec (cuelang.org spec section on aliases, or the v0.18 release notes) instead of reverse-engineering the parser.

Final state of the repo, for the record:

- Builds and runs: `go install ./cmd/veg` succeeds; `veg version`, `--help`, and subcommands work with veg branding, Go 1.27.4, CUE v0.18.0-rc.3.
- Pipeline: `port/copy.sh`, `port/rewrite.sh`, `port/overlay/pkg/verinfo/` — the whole port is re-runnable from those.
- Unfinished, left as-is: ~15 old-style pattern aliases (`[N=string]:`) across `schemas/`, `catalogs/gen/cuelm/`, `pkg/veg/testdata/` — `cue fix --exp=aliasv2` handled the value aliases but not these. Note: `cue.mod/module.cue` is currently set to `language.version: "v0.15.1"` from the fix run and will need to go back to `v0.18.0` as part of finishing that work.
- No commits were made; all changes are working-tree only.

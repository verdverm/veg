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

## Status

nothing started yet

## Vision

- no concerns about backwards incompatibility, consider this greenfield
- initial port should be mostly mechanical (bash scripts using cp and go tooling)
- cleanup and update for core deps and porting exclusions

## Elements to port

Primary:

- The CLI, API, TUI, VS Code, and other elements for the core subsystems that come over
- The core runtime and `#hof` -> `#veg`, abstractions into the core subsystems
- Update core dependencies (CUE, Dagger, ADK), these will require changes
    - we expose several subcommands from `cue` and `dagger`, Dagger is notably a bit hacky by copying/vendoring files from their repo and messing with imports, this was pre-v1 Dagger and may have changed significantly
    - ADK has changed significantly as well with a v1+ (maybe even v2?)
- the cli scaffolding is self code-gen'd, let's mechanically port for now until we have code gen in veg working
- There are a number of packages in `~/hof/hof/lib/...` (renamed to `~/verdverm/veg/pkg/go/...`)
- `~/hof/hof/{catalogs,design,examples,schemas}`
- `~/hof/hof/flow` -> `~/verdverm/veg/pkg/go/flow`

Exclude:

- most anything unrelated the core hof/veg binary (website, docs, github)
- tcell tui, we only want the charm based one
- `~/hof/hof/script` and any tests that use it

Conflicts/Confusing:

- `.hof` vs `.veg` - there was a forward looking addition when adding ADK and the hof agent, these should be merged, not sure if any organization by subdir would help / make sense
- there are much older versions of LLM related things (prompt vs agentic), we only want the newer ADK based agentic systems, exclude the older stuff
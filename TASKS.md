We are working on a large port, read AGENTS.md and NOTES.md for context.

Major porting tasks:

- update ADK for given code
- port docs, change tech stack
- port tests (that make sense), get running in CI
- CUE SDKs for Dagger and ADK (more thorough / accurate)

Major port-porting tasks:

- refactor `veg` agentic subsystem
- rewrite the TUI
- embed/port `gdb` (separate project)
- integrate SpiceDB for permissions, especially agents
- regenerate AGENTS.md for project and core dependencies

## Update ADK

patches: ../adk
latest: ../adk2

- Investigate the current patches in local ADK fork (in particular there are modifications to the API/SDK to support the sandbox/dagger integration into ADK, and a custom session implementation so we can associate chat and env session history together in the database)
- Investigate the latest ADK features, interfaces, subsystems
- Make plan for updating existing veg code for updated ADK

## Port Docs

tbd... base on tanstack, react, tailwind, cloudflare

- ../verdverm.com as a guideline, but let's use filesystem based routing instead of the dynamic ones there
- similar cloudflare setup, domain should be https://veg.verdverm.com (some pnpm script names should be changed for clarity)

## Port Tests

tbd... but should be running on github actions for now

## CUE SDKs for Dagger and ADK

Develop a `./gen` folder/scripts/package for Dagger and ADK which generates the `veg` internals (CUE & Go) to create more complete mappings from CUE -> veg engine(s) -> respective SDKs. Dagger is already pretty close (iirc), unsure about ADK (might be the start of something) but also ADK has changed. Neither has been done mechanically (aiui).

## Refactor agentic and tui subsystems

- Investigate opencode database schema and session management
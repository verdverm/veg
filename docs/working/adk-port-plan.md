# ADK Port - Plan

Goal: move veg from the patched `google.golang.org/adk v0.3.0` (via the `../adk` fork at `v0.4.0-30`) to `google.golang.org/adk/v2 v2.8.0` (`~/verdverm/adk2`), re-expressing the fork's patches as either upstream features or a minimized patch set. Context and evidence: see `docs/working/adk-port-report.md`.

Guiding decisions (from AGENTS.md):
- ADK auth stays out of the tool: OpenAI-compatible model path pointed at a user-run gateway.
- The Dagger sandbox lives entirely in veg (`pkg/container`, `pkg/env/incept`). What exists in the ADK fork is interface changes and modifications to the built-in options and examples; none of that code is used by or required for the sandbox.
- First pass keeps a local `replace`; we are NOT doing a full fork of adk v2 unless forced.
- `.veg/data` sqlite DBs are where session history lives; global default `UserConfigDir/veg/data/veg.db`, overridable via `VEG_DATABASE_CONN` / `VEG_USER_DATA_DIR`.

## Phase 0 - Decide the fork strategy (do first)

The current approach is "fork adk, patch it, replace it". v2 makes most of the fork redundant:

| Fork patch | v2 fate |
|---|---|
| A. shared stateDelta + delete-from-state | re-check behaviorally against v2; if still buggy, upstream issue + minimal patch; else drop |
| B. `Service` Clone/Splice + partial accumulation | interface change to a public interface; v2's `session.Service` does not have Clone/Splice. Options: (1) keep a fork of adk/v2 carrying ONLY this + A (recommend); (2) re-express the operations in veg without changing the interface. Decide in Phase 0 |
| C. GORM session service | upstream `session/database` covers the basics incl. `NewSessionServiceFromDB`; the fork's `Clone`/`Splice` implementations ride along with whatever is decided for B |
| D. runner user-yield + finish injection | re-evaluate against v2 runner (it now yields/persists differently); likely droppable or a small wrapper |
| E. plugin subsystem | replace with upstream `plugin/`; re-register `functioncallmodifier` as the upstream built-in |
| F. filesystem memory/artifact services | re-apply as-is; small, additive (consider upstreaming) |
| G. httprr | drop if present upstream |
| H. housekeeping | drop |

Target: **minimal-fork** - veg depends on `adk/v2` directly for everything except the fork's interface changes. Carry a small, documented patch series on `../adk2` containing: (A) shared stateDelta if still needed, (B) `Clone`/`Splice` on `session.Service` + partial accumulation, (C) their implementation in `session/database`. Everything else (plugins, runner, filesystem services) uses upstream as-is. File upstream issues proposing Clone/Splice as a first-class feature; if accepted upstream, the fork shrinks to zero.

## Phase 1 - Stand up v2 alongside (module mechanics)

1. `go.mod`: add `google.golang.org/adk/v2 v2.8.0`; point the `replace` at `../adk2`. Keep the old `google.golang.org/adk` dependency temporarily (both modules can coexist) so the tree compiles during the port.
2. Port veg code package-by-package (below), switching imports `google.golang.org/adk` -> `google.golang.org/adk/v2` per package as it is migrated (a package cannot mix both).
3. Remove the old adk dependency and the `../adk` replace only at the end.

Mechanical passes (scripts under `port/`, per repo convention):
- import-path rewrite `google.golang.org/adk` -> `google.golang.org/adk/v2`
- `session.NewEvent(x)` -> `session.NewEvent(ctx, x)` (manual review, signature change)
- `tool.Context` -> `agent.Context` in functiontool handlers and toolset methods
- `agent.CallbackContext` -> `agent.Context` in callback registrations

## Phase 2 - Port order (leaf-to-root, compile after each step)

1. **Tools** (`pkg/agent/tools/*`): handler context `tool.Context` -> `agent.Context`; add confirmation fields only if wanted. `mcptoolset` API check.
2. **Models** (`pkg/agent/models/`):
   - Add `Capabilities() model.Capabilities` to the OpenAI and Vertex models.
   - Decision point: keep the hand-rolled openai.go (already on openai-go/v3, works with gateways) vs adopt `model/openaimodel` with `APIChatCompletions`. Recommendation: adopt `openaimodel` for the gateway path (less code to own), keep veg's implementation only if gateway quirks demand it.
   - Vertex/gemini: `model/gemini` still exists in v2; port as-is.
   - Optionally register models in `model/registry` so CUE agent configs resolve models by name/pattern.
3. **Agent construction** (`pkg/agent/agents/cue.go`): update `llmagent.Config` field/callback types; verify `agenttool` signature; no `Agent` interface implementations to fix (veg uses constructors).
4. **Session service** (`pkg/agent/services/session/`):
   - Base persistence: upstream `session/database` (`NewSessionServiceFromDB` on veg's sqlite `*gorm.DB` + `AutoMigrate`).
   - Port the fork's `Clone`/`Splice` interface additions and their in-memory + database implementations (carried by the minimal fork from Phase 0), plus partial-accumulation in `AppendEvent`.
   - Port veg's `localSession`/`storage_session` expectations to the new event shape (`NewEvent(ctx, id)`, `session.ErrNotFound`, new Event fields).
   - Validate against the upstream table layout; run the upstream `session/sessiontestsuite`.
5. **Runtime/handlers** (`pkg/agent/runtime/...`): `runner.Config` new fields; `Run` options; event JSON changes; `session.ErrNotFound` handling; user-event yielding from Patch D - confirm v2 runner behavior, remove veg workarounds.
6. **Plugins**: swap fork `plugin/` usage for upstream `plugin.New(plugin.Config{...})`; wire via `runner.Config.PluginConfig`; re-implement/adopt built-in `functioncallmodifier`.
7. **TUI** (`pkg/agent/cmd/tui`): session type imports only; should compile with little change.

## Phase 3 - Validate

- `go build ./...` and `go install ./cmd/veg` + `veg version` (per repo validate loop).
- Unit tests: `pkg/agent/services/session` tests ported to the new service; upstream `sessiontestsuite` for any custom service.
- Smoke: run an agent over the OpenAI-compatible gateway (Ollama-style base URL) with tools + MCP, check streaming events, state deltas, and delete-from-state semantics; confirm the v2 `MaxLLMCalls` default (500) does not break long TUI sessions.
- Database migration check: existing `.veg/data` sqlite files must still open; if upstream `session/database` table names/columns differ from the fork's gorm models, write a small migration or accept a reset (decide explicitly).

## Phase 4 - Cleanup

- Drop the `../adk` replace and old module; decide whether `../adk2` replace stays (it should, matching the ADK-first-pass posture) and whether the `~/hof/adk` fork can be retired (it can, once Phases 0-3 land).
- File upstream issues/PRs for anything veg had to patch in v2 (Patch A/B semantics, filesystem services from Patch F are good upstream candidates).
- Update `AGENTS.md` (ADK version, stale `~/hof/adk` gotcha) and `NOTES.md`.

## Open questions (need answers before/during Phase 0)

1. Is the shared-stateDelta bug (Patch A) still present in v2? Needs a targeted test.
2. Clone/Splice as upstream features: propose to adk-go; until then the minimal fork carries them.
3. Do we still need the fork's user-event yield + finish-injection (Patch D), or does v2 runner already stream user turns correctly?
4. `session/database` table-layout drift vs existing veg sqlite files: migrate or reset?
5. Does the v2 Event shape (IsolationScope, Routes, RequestedInput, NodeInfo) overlap with anything the fork's Clone/Splice machinery covers?

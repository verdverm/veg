# ADK Port - Investigation Report

Status: investigation for the "Update ADK" task in `TASKS.md`.
Checkouts examined:

- `~/verdverm/adk` (identical to `~/hof/adk`, same commit): the patched fork, branch `veg`, `v0.4.0-30-gbd45719` (`bd45719`)
- `~/verdverm/adk2`: latest upstream, module `google.golang.org/adk/v2`, tag `v2.8.0` (HEAD slightly past, `eb2f3543`)
- `veg` current pin: `google.golang.org/adk v0.3.0` via `replace => ../adk` (resolves to the patched fork above)

## 1. What the fork actually patches (vs upstream base `03705a53`, ~v0.4.0 + 11)

19 local commits on the `veg` branch. Grouped into patches:

### A. State correctness (multi-function-call deltas, delete-from-state)

<!-- user: state correctness is hopefully fixed, not sure about delete support -->

- `internal/context/callback_context.go`: adds `NewCallbackContextWithDelta(ctx, stateDelta)` - one shared stateDelta map per turn.
- `internal/llminternal/base_flow.go`: threads the shared `stateDelta` through `callLLM`, after/on-model callbacks, `finalizeModelResponseEvent` (6 sites).
- `internal/toolinternal/context.go`: uses the new constructor.
- Delete-from-state semantics: nil / empty-string value in `State.Set` deletes the key - in `session/inmemory.go` and `session/database/session.go`, plus `trimTempDeltaState` / `updateSessionState` helpers.
- New `internal/sessionutils/` package (`ExtractStateDeltas`, `MergeStates`) factored out of duplicated logic.

### B. `session.Service` interface changes - Clone + Splice + partial accumulation

<!-- this is the main change we are after -->

Interface changes to ADK: the public `session.Service` interface gains `Clone(ctx, Session) (Session, error)` and `Splice(ctx, Session, start, count, fill) (Session, error)`. These are modifications to the SDK surface and to the built-in implementations and examples; none of this code is used by the Dagger sandbox, which lives in veg.

- `AppendEvent` gains partial-message accumulation (merge `Parts` + `StateDelta` into the last event when `Partial && same author`) in both in-memory and database implementations.
- Tests in `session/inmemory_test.go`, `session/database/service_test.go`, `session/delete_test.go`.

### C. GORM database session service (`session/database/`)

<!-- is this more than an in-repo example to make sure things compile? The real implementation we use is in our code -->

Modifications and additions to the built-in GORM session implementation:

- `databaseService` with `NewSessionService(dialector, opts...)`, `AutoMigrate`.
- Fork addition `NewSessionServiceGorm(db *gorm.DB)` - reuse an already-open `*gorm.DB` (also on the `new-session-db` branch).
- gorm models for sessions/events/app-user-state, dialector-aware JSON serializers (`gorm_datatypes.go`).
- NOTE: this landed upstream around v0.3/v0.4; upstream v2 has since matured it (Postgres/Spanner/SQLite, `NewSessionServiceFromDB`, row locking, compaction persistence).

### D. Runner: user-event yielding + finish injection

<!-- user events / interrupts should be handled by ADK now -->

- `runner/runner.go`: after `appendMessageToSession`, yields the user event to the client stream (with a fix for incorrect indexing).
- Final-response events get `FinishReason = "STOP"`, `TurnComplete = true` in the runner; the equivalent block in `base_flow.go` `Flow.Run` is **commented out** (HEAD commit disabled it - it double-set finish against the debug REST session).

### E. Plugin subsystem (fork addition)

<!-- the question for me here is, what are we doing with the plugin system and functionacllmodifier? -->

- `plugin/` + `plugins/functioncallmodifier/` (used by veg) + `internal/plugininternal/` (PluginManager with ~13 hooks, plugin context).
- Wired into `runner`, `agent`, `llminternal/base_flow`, `testutil`.
- Biggest reconcile-or-replace area: upstream v2 now has its own first-class `plugin/` package with the same hook shape, including a built-in `functioncallmodifier` plugin.

### F. One-shot filesystem services

<!-- these changes are likely just for the Session interface change -->

- `memory/filesystem.go` and `artifact/filesystem.go`: `FilesystemService(baseDir)` implementations. Self-contained.

### G. `internal/httprr/` HTTP record/replay

<!-- we don't need to be recording http here, I'm not sure I had any patches tied to this, would need more details -->

- Likely exists in newer upstream already; verify before carrying forward.

### H. Housekeeping

<!-- we can ignore any AGENTS.md files, we can remake them later -->

- `.gitignore` additions, gorm deps in go.mod, generated `AGENTS.md` files (doc-only).

Patch A is an additional interface/behavior change: `NewCallbackContextWithDelta` threads one shared stateDelta map through multi-function-call turns.

## 3. Latest ADK (v2.8.0): what changed since v0.3.0 (veg's pin)

### Module / deps
- Module path is now `google.golang.org/adk/v2` - every import and the `replace` must change.
- go directive 1.26.10 (veg is already on 1.27.4, fine); genai v1.40 -> v1.77; new heavy deps: openai-go/v3, anthropic-sdk-go, gorm v1.31.2, a2a-go/v2, mcp go-sdk v1.8.

### Breaking API changes (checklist)
1. `model.LLM` gains a required `Capabilities() Capabilities` method.
2. `tool.Context` and `agent.CallbackContext` are GONE, unified into `agent.Context`; toolset `Tools(ctx agent.ReadonlyContext)`; all callback types take `agent.Context`.
3. `session.NewEvent(invocationID)` -> `session.NewEvent(ctx, invocationID)`.
4. `agent.Agent` now has an unexported `internal()` method - external code cannot implement it; must use constructors (`agent.New`, `llmagent.New`, `workflowagent.New`, `workflowagents/*`).
5. Custom `session.Service` implementers must handle: `session.ErrNotFound` sentinel, JSON-tagged `Event` with new fields (`IsolationScope`, `Routes`, `RequestedInput`, `Output`, `NodeInfo`, `EventActions.Compaction`, `RequestedAuthConfigs`), and should pass `session/sessiontestsuite`.
6. `llmagent.Config` gains `Mode`, `OutputArtifact`, `OnModelErrorCallbacks`, `OnToolErrorCallbacks`; callback signatures changed context type.
7. `runner.Config` gains `PluginConfig`, `AutoCreateSession`, `Compaction`, `ContextCache`; `Run` is variadic `RunOption`; new `runner.NewInMemory`.
8. `functiontool` handler signature: `func(tool.Context, ...)` -> `func(agent.Context, ...)`; config gains confirmation fields; new streaming function tool.

### New subsystems relevant to veg

<!-- I assume we will still need a custom session, and likely the patch(es) to support the dagger/timetravelling - in short to have a custom DB schema to support our features -->
- **`model/openaimodel`** - built-in OpenAI-compatible model: `openaimodel.NewModel(ctx, name, &openaimodel.ClientConfig{APIKey, BaseURL, API: openaimodel.APIChatCompletions, ...})`. This is exactly veg's "user-run gateway holds the credentials" decision; it can likely replace `pkg/agent/models/openai.go`.
- **`plugin/`** upstream - same lifecycle-hook shape as the fork's Patch E, with built-ins `loggingplugin`, `retryandreflect`, `functioncallmodifier` (the fork plugin veg actually uses).
- **`session/database`** upstream - GORM-backed `session.Service` with SQLite/Postgres/Spanner, `NewSessionServiceFromDB(db *gorm.DB)` (the fork's `NewSessionServiceGorm`, upstreamed in some form), `AutoMigrate`. Plus `session/sessiontestsuite` conformance tests and `session/compaction`.
- **`workflow/` + `agent/workflowagent`** - node/graph engine with routing, typed schemas, HITL pause/resume. Potential future home for flow<->agent bridging; not required for the upgrade.
- **`model/registry`** - name-pattern based LLM factories (`model.Register`, `model.NewLLM`).
- **`platform/`** - injectable time/UUID seams; **`auth/`**, **`agentregistry/`**, **`client/agentengine`** - not needed now.

## 4. veg's current ADK usage surface (what has to change)

Import census across `pkg/agent`:
- `session` (17 files), `model` (9), `tool` (8), `functiontool` (4), `mcptoolset` (3), `artifact` (3), `agent` (2), `agenttool`, `runner`, `model/gemini`, `agent/llmagent`.

Touch points:
- `pkg/agent/models/openai.go` - hand-rolled `model.LLM` via openai-go/v3 (OpenAI-compatible base URL). Add `Capabilities()` or replace with `model/openaimodel`. <!-- replace -->
- `pkg/agent/models/vertex.go` - `model/gemini`; add `Capabilities()`. <!-- drop gemini specifics, only generic openai -->
- `pkg/agent/agents/cue.go` - builds agents from CUE config via `llmagent.Config`, `agenttool`; uses `InstructionProvider`, tools, MCP. Uses constructors, so no `internal()` problem. Callback/tool context types change.
- `pkg/agent/tools/*` - `functiontool` handlers with `tool.Context` -> `agent.Context`; `mcptoolset` users.
- `pkg/agent/services/session/` - veg's custom session service (localSession, storage_session, service) -> decide: replace with upstream `session/database` SQLite, or port forward as a fork-patched tree (see plan). <!-- absolutely update, the entire point is that their DB is trash and does not support our needs -->
- `pkg/agent/runtime/handlers/{common,api,ws}` - `runner.Config`, `session.Event` construction (`NewEvent(ctx, id)`), event iteration. <!-- ideally this can be improved -->
- `pkg/agent/cmd/tui` - reads `session` types; mostly unaffected.

## 5. Risks / open questions

- The fork adds `Clone`/`Splice` to the `session.Service` interface; upstream v2 does not have them. Porting forward means either a fork of adk/v2 carrying this interface change, or re-expressing the operations without the interface change. <!-- almost certainly need to carry this part of the patch forward -->
- Patch A (shared stateDelta) may be partially fixed upstream by now - needs a behavioral check against v2 before re-porting.
- Whether upstream's `session/database` matches veg's DB schema needs (chat + env session history in one DB) - veg's tables wrap gorm models directly; likely fine, but the veg `storage_session` mapping and `.veg/data` sqlite paths must be reconciled. <!-- I doubt very much we can use their schema, we want full customization anyway -->
- Plugin subsystem: replacing the fork plugin with upstream `plugin/` changes where `functioncallmodifier` hooks are registered (runner `PluginConfig` vs fork's internal manager) - low risk, mechanical.

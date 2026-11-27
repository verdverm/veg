# ADK Port - Plan

Goal: replace the patched `google.golang.org/adk v0.3.0` dependency (`../adk`) with `google.golang.org/adk/v2 v2.8.0` (`../adk2`), retaining only the SDK changes required by veg. Investigation and API details are in `docs/working/adk-port-report.md`.

## Decisions

- Use upstream `model/openaimodel` with `APIChatCompletions` and a user-run OpenAI-compatible gateway. Remove the hand-written OpenAI adapter and all Vertex/Gemini-specific model setup.
- Keep session persistence in `pkg/agent/services/session`. Veg owns its GORM schema and shares the database with environment history; upstream `session/database` is a reference, not a replacement.
- Preserve Clone/Splice and partial-event accumulation. Carry a minimal `../adk2` patch for the public `session.Service` interface and its built-in implementations, while the primary behavior remains in veg.
- Use upstream state-delta handling, runner finalization, and `runner.WithYieldUserMessage()`.
- Do not port the fork's plugin manager. Veg has no current plugin registrations; any future function-call modifier uses upstream `plugin/functioncallmodifier`.
- Drop ADK artifact and memory integration. Remove the runtime artifact plumbing and omit `ArtifactService` from runner construction. Filesystem access belongs to the veg environment/sandbox, not a parallel artifact store.
- Session data remains under `.veg/data`, with the global default at `UserConfigDir/veg/data/veg.db` and overrides through `VEG_DATABASE_CONN` and `VEG_USER_DATA_DIR`.

## Phase 0 - Establish the v2 baseline and minimal patch

1. Create a clean, reproducible v2 patch series in `../adk2`, based on the selected v2.8.0-compatible checkout.
2. Add Clone/Splice to `session.Service` and update every built-in implementation needed for ADK to compile and its tests to run. Keep this patch limited to the session operation surface; do not port the old runner, plugin, state-delta, filesystem, or HTTPRR patches.
3. Add focused tests for the new operations. Test the full transactional behavior in veg's service rather than duplicating its schema in ADK.
4. Confirm with targeted upstream tests that shared state deltas survive multiple function calls. Treat state deletion as veg behavior: nil deletes a key, while an empty string remains a value.

## Phase 1 - Change the module dependency

1. Add `google.golang.org/adk/v2 v2.8.0` and `replace google.golang.org/adk/v2 => ../adk2`.
2. Port packages leaf-to-root, switching imports from `google.golang.org/adk` to `google.golang.org/adk/v2`. Keep the old module temporarily so unmigrated packages continue to compile.
3. Remove the old requirement and `../adk` replacement after all imports are migrated; run `go mod tidy`.
4. Capture mechanical import and signature rewrites under `port/` where practical. Ensure repeated script runs do not append a second `/v2`, and capture manual fixups as post-rewrite patches per the repository porting rules.

Mechanical changes include:

- `session.NewEvent(id)` to `session.NewEvent(ctx, id)`.
- `tool.Context` and `agent.CallbackContext` to `agent.Context`.
- Toolset and callback signature updates.
- Import-path changes, followed by package-specific API fixes.

## Phase 2 - Port veg packages

Compile after each step. Keep automated tests focused on concrete service behavior, and validate the integrated runtime in an actual agent session. Do not build a mock-heavy ADK test harness for this port.

### 1. Tools and agent construction

- Update `pkg/agent/tools` function handlers and MCP toolsets.
- Update `pkg/agent/agents/cue.go` for v2 `llmagent.Config`, callbacks, `agenttool`, and toolset types.
- Continue using ADK constructors; do not implement `agent.Agent` externally.

### 2. Models

- Replace `pkg/agent/models/openai.go` with a small constructor around `openaimodel.NewModel` using `APIChatCompletions`, the configured model ID, and base URL.
- Delete `pkg/agent/models/vertex.go` and the `vertex` branch in `Runtime.initModels`.
- Keep model configuration generic and OpenAI-compatible. Do not add Google/Gemini authentication or provider-specific setup.
- Defer optional `model/registry` integration until the CUE agent configuration needs named model factories.

### 3. Custom session service

- Port `pkg/agent/services/session` to the v2 session interfaces and complete event shape.
- Preserve the veg-owned tables and shared `*gorm.DB`; do not run upstream `session/database.AutoMigrate` against veg databases.
- Implement `session.ErrNotFound`, deterministic list ordering, event ID assignment, temporary-state filtering, and all fields required by `session/sessiontestsuite`.
- Port transactional Clone/Splice and partial-event accumulation, including in-memory session consistency and timestamp ordering.
- Implement deletion consistently for app, user, and session state: nil deletes, empty strings persist.
- Extend storage models for the complete v2 event and action fields. Legacy schema migration and backward compatibility are not requirements.

### 4. CUE configuration

- Treat `.veg/models.cue`, `.veg/agents.cue`, `.veg/tools.cue`, `schemas/agentic`, and `pkg/agent/config/types.go` as one configuration contract and update them together with the Go builders.
- Make the model shape OpenAI-compatible only: remove the provider switch, Gemini model entries, and Gemini defaults; retain the model ID and gateway base URL needed by `openaimodel`.
- Align agent fields and defaults with the v2 features veg actually constructs: model, description, instructions, tools, toolsets, MCP, and subagents. Keep veg-specific environment configuration explicit rather than mirroring every ADK option.
- Load the updated CUE and construct the default agents to catch schema, decoding, default-model, and model-reference errors.

### 5. Runtime, handlers, and service cleanup

- Update runner configuration and `Run` calls in `pkg/agent/runtime/handlers/common`.
- Remove `artifact.Service` fields, getters, initialization, and interface requirements from the runtime and handlers. Remove the empty `pkg/agent/services/artifact` scaffold and the obsolete filesystem artifact fork dependency.
- Pass `runner.WithYieldUserMessage()` and remove reliance on fork runner modifications.
- Review channel production/consumption in the common, REST, WebSocket, and TUI paths so cancellation, iterator errors, terminal events, and interruptions propagate without duplicate user events or blocked sends.
- Update session and event JSON handling for v2 while preserving veg's Clone/Splice API behavior.
- Do not register plugins during this port unless a concrete veg hook is added; use upstream `runner.Config.PluginConfig` when that occurs.

## Phase 3 - Validate

- Run `go install ./cmd/veg` and `veg version` as the repository validation loop.
- Run the existing session service tests and `session/sessiontestsuite` against the custom service using real SQLite databases. Add focused cases only for uncovered v2 behavior, Clone/Splice, deletion, and partial accumulation.
- Load the updated CUE configuration and construct the default agents without provider or model-reference errors.
- Validate an actual agent session through an OpenAI-compatible gateway with streaming, function tools, and MCP. Check state deltas, deletion, user-event ordering, final/interrupted events, and the 500-call default.
- Exercise Clone/Splice through both REST and WebSocket handlers and verify the associated environment-history workflow.
- Verify a fresh database serves both session and environment records. Legacy database migration is not a validation requirement.

## Phase 4 - Clean up

- Remove the old ADK module, `../adk` replacement, unused OpenAI conversion code, Vertex configuration, artifact plumbing, and fork-only workarounds.
- Keep the `../adk2` replacement only for the documented Clone/Splice patch series. Upstream those operations if a suitable public API is accepted; otherwise document patch maintenance and upgrade steps.
- Update `AGENTS.md`, `NOTES.md`, and subsystem documentation with the final module version and service ownership.

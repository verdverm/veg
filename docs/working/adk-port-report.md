# ADK Port - Investigation Report

Status: investigation complete; implementation pending. The implementation sequence is in `docs/working/adk-port-plan.md`.

## Baseline

- Current veg dependency: `google.golang.org/adk v0.3.0`, replaced by `../adk`.
- Patched fork: `~/verdverm/adk`, branch `veg`, commit `bd45719` (`v0.4.0-30-gbd45719`). The `~/hof/adk` checkout is identical.
- Target: `google.golang.org/adk/v2 v2.8.0`, developed through `~/verdverm/adk2`. The local checkout is at `eb2f3543`, 78 commits after the `v2.8.0` tag.
- The fork comparison base is upstream commit `03705a53`; the `veg` branch contains 19 local commits.

## Fork patch inventory and disposition

### A. Shared state deltas and state deletion

The fork threads one state-delta map through model and tool callbacks with `NewCallbackContextWithDelta`, preventing function calls in the same turn from overwriting each other's deltas. ADK v2 now has this mechanism in `internal/context` and `internal/llminternal`, including artifact deltas. Do not carry this part of the fork forward.

The fork also treats nil and empty-string state values as deletion. ADK v2's in-memory and database `State.Set` implementations still assign values directly and do not implement deletion. Veg needs deletion for `SessionStateDel`, so its custom session service must own this behavior. Use nil as the deletion marker and preserve empty strings as valid values; cover app, user, and session state in tests.

The fork's `internal/sessionutils` extraction and merge helpers have upstream v2 equivalents.

### B. `session.Service` Clone/Splice and partial accumulation

The fork adds these methods to the public `session.Service` interface:

```go
Clone(ctx context.Context, sess session.Session) (session.Session, error)
Splice(ctx context.Context, sess session.Session, start, count int, fill session.Events) (session.Session, error)
```

These operations are required by veg's REST and WebSocket session handlers. They support branching and time-travel workflows alongside the Dagger environment history; the sandbox implementation itself remains in veg.

ADK v2 does not expose Clone/Splice. Carry a minimal `../adk2` patch that extends the interface and keeps its built-in implementations compiling. Implement the real persistence semantics in veg's custom session service, including transactional event copying/replacement and correct ordering.

The fork also accumulates partial messages in `AppendEvent`, merging content parts and state deltas into the previous event when the partial event and author match. Upstream v2 does not provide that behavior. Port it to the veg service and test streaming, retries, state updates, and final-event handling.

### C. GORM session service

The fork's `session/database` package contains a GORM service, JSON serializers, and `NewSessionServiceGorm` for reusing an existing database connection. Veg's real implementation already lives in `pkg/agent/services/session` and shares a `*gorm.DB` with environment persistence.

Upstream v2 has a more mature `session/database` package with `NewSessionServiceFromDB`, migrations, locking, and compaction support. It is a useful API and testing reference, but it is not the persistence layer for veg. Retain a veg-owned schema so chat events, session state, environment history, cloning, and splicing can evolve together.

### D. Runner user events and finish injection

The fork yields the persisted user event and injects final-response finish metadata. ADK v2 provides `runner.WithYieldUserMessage()` for the current LLM-agent node path, with tests for user-event state deltas. Use that option and remove the fork's runner modifications. Rely on v2 event finalization rather than reapplying the finish-injection patch; validate interruption and terminal-event behavior through the veg handlers.

### E. Plugins and function-call modification

The fork adds a plugin manager and a `plugins/functioncallmodifier` package. ADK v2 now has a first-class `plugin` package and `plugin/functioncallmodifier`, configured through `runner.Config.PluginConfig`.

Current veg code does not import or register the fork's plugin system, so there are no veg call sites to migrate. Do not carry the fork implementation. If function-call modification is needed, construct the upstream plugin and register it through the runner.

### F. Filesystem memory and artifact services

The fork adds self-contained filesystem-backed memory and artifact services. They are unrelated to the session interface changes.

Veg does not need either service. The artifact service is only wired into runtime scaffolding and the runner; there are no artifact API workflows to preserve. Remove that plumbing and omit `runner.Config.ArtifactService`. Filesystem operations belong to the veg environment/sandbox, not a separate ADK artifact store. Drop the unused memory service and the empty `pkg/agent/services/artifact` scaffold.

### G. HTTP record/replay

`internal/httprr` exists upstream in v2 and supports ADK's own tests. It is internal test infrastructure, not a veg runtime dependency. No port is needed.

### H. Housekeeping

Ignore generated `AGENTS.md` files, ignore-rule changes, and other fork housekeeping. Dependency changes come from the v2 module and normal Go module tidy operations.

## ADK v2 migration surface

### Module and model APIs

- Change imports and the local replacement from `google.golang.org/adk` to `google.golang.org/adk/v2`.
- Replace `pkg/agent/models/openai.go` with upstream `model/openaimodel`, explicitly selecting `openaimodel.APIChatCompletions` for gateway compatibility.
- Remove `pkg/agent/models/vertex.go`, the `model/gemini` import, and the Vertex provider branch. Veg's model configuration uses the generic OpenAI-compatible path, with credentials held by the user-run gateway.
- A custom `model.LLM` would need `Capabilities() model.Capabilities`; using `openaimodel` avoids maintaining that implementation in veg.

### Agents and tools

- `tool.Context` and `agent.CallbackContext` are replaced by `agent.Context`; toolsets use `agent.ReadonlyContext`.
- Update function-tool handlers, callbacks, `llmagent.Config`, `agenttool`, and MCP toolset construction to their v2 signatures.
- External packages cannot implement `agent.Agent` directly because v2 adds an unexported interface method. Veg uses constructors, so no agent-interface reimplementation is required.

### Sessions and events

- `session.NewEvent(invocationID)` becomes `session.NewEvent(ctx, invocationID)`.
- Custom services must wrap `session.ErrNotFound` for missing sessions and satisfy `session/sessiontestsuite`.
- Extend veg's event storage and JSON mapping for the complete v2 event shape, including `IsolationScope`, `Routes`, `RequestedInput`, `Output`, `NodeInfo`, and all `EventActions` fields such as compaction, tool confirmations, and requested auth configs.
- Preserve Clone/Splice and partial accumulation in addition to the upstream service contract.

### Runner

- `runner.Config` adds plugin, session creation, compaction, and context-cache configuration.
- `Run` accepts variadic `RunOption` values; veg should pass `runner.WithYieldUserMessage()`.
- The default model-call budget is 500 per invocation. Validate long-running sessions and set an explicit veg policy if that default is unsuitable.

### CUE configuration

The CUE schemas, concrete `.veg` configuration, decoded Go types, and agent/model builders must be updated as one contract. Remove the model provider switch, Gemini entries, and Gemini defaults; the remaining model shape supplies an ID and OpenAI-compatible gateway base URL. Keep agent configuration limited to the v2 capabilities veg uses, with environment settings remaining veg-specific. Validate by loading the CUE and constructing the default agents.

## Current veg touch points

The import census under `pkg/agent` is 17 files using `session`, 9 using `model`, 8 using `tool`, 4 using `functiontool`, 3 using `mcptoolset`, 3 using `artifact`, 2 using `agent`, and one each using `agenttool`, `runner`, `model/gemini`, and `agent/llmagent`. The artifact imports are removal targets rather than migration targets.

The main implementation areas are:

- `pkg/agent/models` and `pkg/agent/runtime/runtime.go`: model construction and removal of provider selection and artifact plumbing.
- `pkg/agent/agents/cue.go` and `pkg/agent/tools`: agent, callback, function-tool, and MCP APIs.
- `.veg`, `schemas/agentic`, and `pkg/agent/config`: CUE shapes, defaults, decoding, and model references.
- `pkg/agent/services/session`: the custom GORM service, event mapping, Clone/Splice, partial accumulation, and tests.
- `pkg/agent/runtime/handlers/common`: runner construction, user-event streaming, Clone/Splice calls, and state deletion.
- `pkg/agent/runtime/handlers/{api,ws}` and `pkg/agent/cmd/tui`: updated event and session types plus streaming/error behavior.

## Validation requirements

- Run upstream session conformance tests against the veg service and retain dedicated Clone/Splice, deletion, and partial-stream tests.
- Load the updated CUE and construct the default agents with valid model references.
- Exercise a real agent session through the OpenAI-compatible gateway with streaming and tools, including MCP; do not create a mock-heavy harness for integration validation.
- Confirm user events are yielded once, terminal and interrupted events retain correct metadata, and no fork finish injection is necessary.
- Verify a fresh veg SQLite database supports sessions and environment history. Legacy migration is not required.

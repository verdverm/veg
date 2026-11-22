# NOTES.md - hof -> veg port reference

Source repo: `~/hof/hof` (module `github.com/hofstadter-io/hof`, go 1.24.4).
Dest repo: `~/verdverm/veg`. Task: port core binary elements; exclude website/docs, tcell TUI, script tests.

## 0. Quick orientation map

| Area | Source path | Notes |
|---|---|---|
| CLI (generated) | `cmd/hof/{main.go,cmd/,flags/}` | cobra; generated from `design/` via external CUE module `hofmod-cli` |
| Core runtime | `lib/runtime/` | `Runtime` struct: CUE load, dagger, deps; `Load()`, `EnsureInfra()` |
| `#hof` metadata | `lib/hof/` | `Hof` struct, `ParseHof`, `FindHofs`, attr->struct upgrade, `Node[T]`, `Upgrade` |
| Codegen | `lib/gen/` + `lib/gen/cmd/` | `Generator`/`File`; shadow dir + diff3; text/template only |
| Templates | `lib/templates/` | Go text/template + big FuncMap (`AddGolangHelpers`) |
| Flow engine | `flow/` (top-level) | cue/flow based; `flow/context`, `flow/tasker`, `flow/tasks`, `flow/middleware` |
| Agent (ADK) | `lib/agent/` | current LLM system; ADK `google.golang.org/adk v0.3.0` via local replace `../adk` |
| Agent TUI (charm) | `lib/agent/cmd/tui/` | bubbletea/lipgloss/glamour - KEEP |
| tcell TUI | `lib/tui/` (vendored tview fork!) | DROP; `cmd/hof/cmd/tui.go` must be rewired/dropped |
| Dagger CUE layer | `lib/env/dag/` | clean, public `dagger.io/dagger` SDK only - safe port |
| Dagger hack | `lib/env/incept/` + `lib/env/incept/distconsts/` | vendored Dagger engine internals (go.mod replace); reassess vs modern Dagger |
| Container CLI | `lib/container/` | docker/nerdctl/podman wrapper (not Dagger); used by `lib/fmt` formatters |
| CUE cmd wrappers | `lib/cuecmd/` + `lib/cuetils/` | eval/def/export/vet; `lib/cuetils` has duplicated load logic |
| Mod (cue proxy) | `cmd/hof/cmd/mod.go` | imports `cuelang.org/go/cmd/cue/cmd` directly |
| create | `lib/create/` + `lib/repos/` | creator scaffolding; repos = git/OCI module fetch + cache |
| Schemas (CUE) | `schemas/` | `#hof` metadata schema tree (env, dm, gen, agentic, create, prompt, common) |
| CLI design | `design/` | drives CLI codegen; `design/main.cue` has Module/Outdir; `design/cmds/*.cue` |
| Env catalogs | `catalogs/env/{bases,packs,utils}` | container packs for `veg env` |
| Config dirs | `.hof/` (only shadow + .gitignore) vs `.veg/` (agent/env config, `package veg`) | `.hof/shadow` = old codegen artifacts; `.veg` is the live agent config |

## 1. CLI structure & self-codegen

- Entry: `cmd/hof/main.go` -> `cmd/hof/cmd.RunExit()` (root.go has `RootCmd`, `RootInit()`, `RunExit/Int/Err`; pprof via `HOF_CPU_PROFILE`).
- Every `cmd/hof/cmd/*.go` and `cmd/hof/flags/*.go` is GENERATED. Source of truth: `design/` (CUE). `design/main.cue` defines `CLI: schema.Cli & {...}` with Name/Module/Outdir; per-command specs in `design/cmds/*.cue` each have `Body:"err = libX.Y(...)"` strings that become generated cobra Run funcs.
- Driven by `hof.cue` (repo root): `@gen(cli,hof)`, `WatchGlobs:["./design/**/*"]`, `WatchExec: go install ./cmd/hof`. Generator schema+templates come from EXTERNAL CUE module `github.com/hofstadter-io/hofmod-cli` (declared in `cue.mod/module.cue`; NOT vendored in-tree - lives in CUE module cache).
- Implication for port: to keep mechanical port, either (a) fork/port hofmod-cli too, or (b) hand-port the generated cmd files for now and revisit. Renaming Module in `design/main.cue` regenerates everything.
- `.hof/shadow/cli/**` = shadow copies of generated files (bookkeeping, not source; excluded per port rules).
- Flags pattern: `XxxFlagpole` struct + `SetupXxxFlags` + package var `XxxFlags` (see `cmd/hof/flags/root.go` for `RootPflagpole`, includes `-p/--package`).
- Command -> backing package: create->`lib/create`; datamodel->`lib/datamodel/cmd`; gen->`lib/gen/cmd`; flow->`flow/cmd`; fmt->`lib/fmt`; def/eval/export/vet->`lib/cuecmd`; mod->cue proxy; extension->`lib/extension`; agent->`lib/agent/cmd` (`libagentcmd`); env->`lib/env/cmd` (`libenvcmd`); tui->`lib/tui/cmd` (DROP, rewire to charm or drop); run->adhoc runtime; version->`lib/container`.

## 2. Core runtime

- `lib/runtime/runtime.go`: `Runtime` central struct - `Ctx`, `DB *gorm.DB`, `API *echo.Echo`, `DAG *dagger.Client`, `Flags flags.RootPflagpole`, dir bookkeeping (`CueModuleRoot`, `CueExtractDir`, `RootToCwd`, ...), CUE fields (`CueContext *cue.Context`, `CueConfig *load.Config`, `BuildInstances`, `Value cue.Value`), result collections `Nodes []*hof.Node[any]`, `Envs`, `Agentics`, `Datamodels`, `Generators`, `Workflows`.
- `New(entrypoints, rflags)`: builds `load.Config` (AcceptLegacyModules, Overlay, TagVars), injects git tags (`lib/yagu`) + CI tags (hofstadter-io/cinful), `--inject-env` for os.Environ. `CueExtractDir = UserCacheDir/cue/mod/extract`.
- `lib/runtime/load.go`: `Load()` = prepPlacedDatafiles (`file@path.to.field` placement) -> load -> prepPlacedUserfiles -> Validate -> `hof.FindHofsOrig(R.Value)`. Shared `cue.Context` via `lib/singletons.CueContext()`. Orphaned-file placement is "hacky" (`DontPlaceOrphanedFiles` opt-out).
- GOTCHA: `lib/runtime/load.go` and `lib/cuetils` (`CueRuntime`) have parallel/duplicated load logic.
- `lib/runtime/infra.go`: `EnsureInfra()` = local OCI registry + dagger engine containers. Env var already renamed: `VEG_EXPERIMENTAL_DAGGER_RUNNER_HOST = container://veg-dagger-engine` (half-done rebrand, `// [[ hof2veg ]]` markers exist e.g. `lib/hof/hof.go`).
- Enrich dispatch pattern (`lib/runtime/{generator,flow,env,datamodel,agentic}.go`): filter `R.Nodes` by `node.Hof.<X>.Root` + name regex/keepFilter (`common.go`), `hof.Upgrade` node -> typed (e.g. `gen.Generator`), collect into `R.<X>s`, run enricher func. `GeneratorEnricher`, `FlowEnricher`, `EnvEnricher`, `DatamodelEnricher`, `AgenticEnricher` types.

## 3. `#hof` metadata (`lib/hof/`)

- `type Hof struct` (`hof.go`): `ID`, `Label`, `Path`, `APIVersion`, `Metadata`, `Memo`, per-feature sub-structs `Datamodel`, `Gen`, `Flow`(+`Pool`,`Print`), `Env`, `Agentic`, `Extra map[string]any`.
- `ParseHof[T](cue.Value)`: `val.LookupPath(cue.ParsePath("#hof"))` -> decode -> `upgradeAttrs(node, label)` converts sugar attrs (`@gen()`, `@flow()`, `@task()`, `@env()`, `@agent()`, `@pool()`, `@dm()`) into `Hof` fields. `FindHofs`/`FindHofsOrig` walk tree into `R.Nodes`.
- `Node[T]` generic wrapper (Hof, Value, parent/children); `Upgrade[any,T]` converts to typed node. This is the heart of `#hof` -> `#veg` rename (also rename Go field `Hof` -> `Veg`, struct `lib/hof` -> `lib/veg` package decision needed).

## 4. Codegen (`lib/gen`)

- `Generator` (`lib/gen/generator.go`): embeds `*hof.Node[Generator]`; `Outdir`, `In`, `Val`, `Templates/Partials []*TemplateGlobs`, `Statics []*StaticGlobs`, `Embedded*` maps, `Generators map[string]*Generator` (subgen composition), `ModuleName` (vendor lookup, `CUE_VENDOR_DIR = cue.mod/pkg`), `UseDiff3`, `Out []*File`.
- `File` (`lib/gen/file.go`): `Filepath` + mutually exclusive selectors: `TemplateContent` (inline, precedence), `TemplatePath`, `DatafileFormat` (CUE export), `StaticFile` (copy). Plus CUE export knobs (`Raw,Final,Concrete,Definitions,...`), delims, shadow/diff fields.
- Template engine: ONLY Go `text/template` (`lib/templates`). "Static" = copy-through, not an engine. Delims configurable per glob (`Delims{LHS,RHS}`) - that is the "alternate template style". FuncMap in `lib/templates/helpers.go` (`AddGolangHelpers`): yaml/toml/json/xml, case helpers (kace), `dict`, `file`, `indent`, `render`, `lookup`, `builtin`.
- Run flow (`lib/gen/cmd/`): own `Runtime` wrapping `lib/runtime.Runtime`; `genOnce()` = `RunGenerators()` + `WriteOutput()` + `CleanupRemainingShadow()` + `PrintMergeConflicts()` + optional PostFlow (gated `--exec`); watch mode via `lib/watch`. Diff3 forced on for real generators. Shadow dir const `SHADOW_DIR = ".hof/shadow/"` (`runtime.go:83`) - rename site.
- `lib/gen/cmd/init.go`: `InitModule` scaffolds gen module dirs (`templates,partials,statics,examples,creators,gen,schema`); default module hardcoded `"hof.io"` (lines ~21,122 - rebrand site). `adhoc.go` backs `veg gen --template`.
- Generator module resolution: `calcBasePath()` uses `R.DepMapping` + `CueModuleRoot` to find external generator modules in CUE cache/vendor.
- Load-order inconsistency (TODO comment): statics disk->embed, templates/partials embed->disk.

## 5. Flow engine (`flow/`, package `github.com/hofstadter-io/hof/flow/...`)

- Built on `cuelang.org/go/tools/flow`. `flow/flow/flow.go`: `type Flow struct{ *hof.Node[Flow]; Root,Orig,Final cue.Value; FlowCtx *flowctx.Context; Ctrl *cueflow.Controller }`. `Start()` -> cueflow.Config{IgnoreConcrete:true, FindHiddenTasks:true, UpdateFunc}, `tasker.NewTasker`, run.
- `flow/context/context.go`: `Context` (per-flow-exec state: RootValue, FlowStack, Middlewares, TaskRegistry *sync.Map, Pools, CUELock, Mailbox). `Middleware{ Apply(*Context, RunnerFunc) RunnerFunc }`; `Register(key, RunnerFunc)` / `Lookup`.
- `flow/tasker/tasker.go`: for each value `hof.ParseHof`, check `node.Hof.Flow.Task != ""`, lookup runner, apply middlewares in REVERSE order, run, `Fill(value)` results back into CUE.
- Task ops registered in `flow/tasks/reg.go` `RegisterDefaults`: noop, nest, api.Call/Serve, csp.Chan, cue.Format, db.* (sqlite/postgres), gen.* (uuid/slug/cuid...), hof.Template, kv.Mem, msg.* (irc/watermill), os.* (Exec/ReadFile/WriteFile/Glob/Getenv/Watch/FileLock...), prompt.Prompt, st.* (Diff/Insert/Mask/Patch/Pick/Replace/Upsert via `lib/structural`). Schemas embedded (`embed.go`) + in `schemas/`.
- Middlewares (`flow/middleware/`): `defaults.go` currently enables info.NewProgress + sync.NewPool only (print/bookkeeping/dummy commented out). `sync/pool.go` implements `@pool()`.
- Command: `flow/cmd/run.go` - workerpool over `R.Workflows`, `-B/--bulk src@dst`, `--stats`. Selection via `@flow` attr / `-F` flag.
- Chat: NOT in flow. `lib/runtime/chat.go` is fully commented (dead). Live chat = `veg agent chat` -> `lib/agent/cmd.Chat` -> bubbletea TUI.

## 6. Agent system (ADK) - `lib/agent/`

- ADK: `google.golang.org/adk v0.3.0` with `replace google.golang.org/adk => ../adk` (LOCAL SIBLING CHECKOUT - find out what patches it carries before moving to upstream v1+/v2).
- `lib/agent/veg-node.go`: `Agentic` struct (the `#hof agentic` node type). `lib/agent/config/types.go`: `Config{Presets, Models, Agents, Tools, Toolsets, Environs, Embeds (EmbedDir ".veg/embed"), AgentsMD, Templates}`; `Agent{Name,Model,Instruction,Tools,Toolsets,Mcp,SubAgents,AutoLoadWorkdir,Environ,AgentsMD}`.
- `lib/agent/agents/cue.go`: CUE->ADK bridge. `OldAgenticCUE(agentDir, models)` loads CUE dir (package `veg`, `.veg/` -> `./` prefix); `BuildAgent` -> `llmagent.New(llmagent.Config{...})`; `buildTools` (tools/cache, tools/filesys, tools/exec); `buildMcp` (github/tavily/quickbooks); `addCallbacks` (ADK Before/After hooks); `RenderInstructions` via `lib/templates` with agentsMd/state injection.
- Models: `lib/agent/models/{vertex.go, openai.go}` - vertex via `adk/model/gemini`, openai is hand-written ADK `model.LLM` wrapping `openai-go/v3`. Both ADK-native (keep).
- Runtime: `lib/agent/runtime/` - singleton with echo server (`:2257`), gorm DB, Models map, websocket client; `ReadEnvConfig()` (CUE dir source), `BackfillAgentic()` (merges `#veg agentic` nodes from core runtime scan into config). GOTCHA: two config sources (CUE dir `.veg` + `#veg` metadata) to reconcile.
- Handlers: `lib/agent/runtime/handlers/{ws,api,common}/` - chat/sessions/filesys/runtime; uses `adk/runner`, `adk/session`.
- Services: `session/` (adk session + gorm), `environ/` (Dagger-backed virtual envs; singleton env store, gorm + OCI registry versioning; sets `_EXPERIMENTAL_DAGGER_RUNNER_HOST=container://veg-dagger-engine`), `artifact/` (blob store), `container.go` (stores interface).
- Tools: `lib/agent/tools/{cache,exec,filesys,mathy.go,mcp/*}`.
- Cmd: `lib/agent/cmd/{main.go,runtime.go,serve.go,chat.go,bulk.go,list.go,eval.go}` + `cmd/hof/cmd/agent/*` cobra proxies.
- `.veg/` dir at repo root = live agent/env config (`agents.cue`, `models.cue`, `tools.cue`, `envs.cue`, `spaces.cue`, `presets.cue`, `containers.cue`, `services.cue`, `commands.cue`, `flags.cue`, `embed/`, `tasks/`, `plans/`), all `package veg` but still importing `github.com/hofstadter-io/hof/...` (mid-rebrand). `.hof/` only contains `.gitignore` + `shadow/` (old codegen artifacts). Merge plan: `.veg` becomes the single dir; shadow regeneration moves under it or is dropped.
- `.hof` hardcoded paths scattered: `lib/datamodel/history.go:17` (`.hof/dm`), `lib/yagu/exec.go:32,38`, `lib/gen/cmd/runtime.go:83` (SHADOW_DIR), `lib/tui/modules/eval/persist.go:43` (dropped anyway), flow middleware defaults (commented).
- LEGACY to EXCLUDE: `lib/prompt/` is NOT LLM - it is `AlecAivazis/survey/v2` interactive questionnaires used by `lib/create` and `flow/tasks/prompt`. The AGENTS exclusion means legacy prompt-based LLM paths (none found as separate SDK; all live LLM is ADK). Decide with user whether survey prompts stay (needed by `create`). `schemas/prompt/` + `flow/tasks/prompt/` depend on it. `hack/genai.go` is scratch.
- `schemas/agentic/agent.cue` still uses `#hof agentic kind:"agent"` with stale comment - rename to `#veg`.

## 7. Dagger / containers

- Two distinct things: `lib/container/` = docker/nerdctl/podman CLI wrapper (used by `lib/fmt` formatters); Dagger = `lib/env/` + `lib/agent/services/environ/`.
- CLEAN Dagger layer: `lib/env/dag/` (`container.go, exec.go, filesys.go, host.go, git.go, volume.go, service.go, export.go, diff.go, sbom.go, ...`) - imports only public `dagger.io/dagger`. Safe port.
- HACK: `lib/env/incept/` copies Dagger engine client internals (imports `github.com/dagger/dagger/engine/client`, `slog`, `telemetry`, `dagql/dagui/idtui`) so hof can run dagger "in itself" for pretty TUI logs. Supported by `go.mod` replace: `github.com/dagger/dagger/engine/distconsts => ./lib/env/incept/distconsts` (copied subtree with its own go.mod; hand-copied `consts.go` from Dagger repo). Commented-out `// replace github.com/dagger/dagger => ../dagger`.
- Versions: SDK `dagger.io/dagger v0.20.0`, engine internals `v0.20.1-0.20260227...`, engine image const `VEG_DAGGER_ENGINE_DEFAULT_IMAGE = registry.dagger.io/engine:v0.19.10` (`lib/consts/veg.go`) - drift to reconcile. Dagger is stable now; evaluate dropping `incept` + distconsts replace for upstream.
- Dagger surfaces in CLI via `veg env` only (no direct dagger subcommand proxy).

## 8. Other lib packages (port candidates)

- `lib/structural/` - CUE value surgery (require/insert/upsert/replace/pick/walk/mask/patch/diff). CUE-only. Clean, core.
- `lib/diff3/` - self-contained 3-way merge (from node-diff3/Synchrotron). Clean.
- `lib/datamodel/` - CUE datamodel snapshots/history/diff/tree/log. Couples to `.hof/dm` path + `lib/hof` nodes.
- `lib/watch/` - fsnotify + debouncer, takes a Runner func. Small.
- `lib/yagu/` - grab bag: exec, dir, glob, copy, git, netrc, ssh, http, port, uuid, hash, browser, stripansi, table, billy. Generic; fix hardcoded `.hof` paths.
- `lib/fmt/` - formatters: built-in (cue/json/yaml/toml/xml) + container-based external via `lib/container`. Couples to runtime Env.
- `lib/connector/` - reflection-based method-call Connector graph, standalone.
- `lib/gotils/` - par, queue, txtar, lockedfile, intern, diff, auth/netrc. Utility.
- `lib/repos/` - module/repo fetching: `remote/` (git vs OCI Kind, mirrors), `git/` (go-git), `oci/` (go-containerregistry), `cache/`. Used by `lib/create` + `hof mod`.
- `lib/database/` - only a readme (aspirational). Skip or minimal.
- `lib/cuetils/` - CUE helpers: CueRuntime, printing, path/attr helpers, `argsplit.go`. Note duplication with lib/runtime.

## 9. Assets & consumption

- All CUE imports use full Go module paths (`github.com/hofstadter-io/hof/schemas/env`), resolved via `cue.mod/module.cue` `source.kind: "self"`. Renaming module => must update every CUE import in schemas/catalogs/design/examples too. No `go:embed` of these dirs found.
- `schemas/` defines `#hof` metadata schema tree (apiVersion v1beta1). Note two distinct "schemas": repo-local dir vs external CUE package `github.com/hofstadter-io/schemas` (used by design/cmds/datamodel.cue, lib/gen/cmd/init.go). External hofmod-cli + hofstadter-io/schemas are separate repos - cannot rename in-tree.
- `catalogs/env/` = container packs (alpine/wolfi/rocky/debian; node/go/python/cue; postgres; dagger/docker/registry...). Needed by `veg env`.
- `examples/` = fixtures (optional). `design/` = CLI codegen (core). `docs/`, `docs-old/`, website = EXCLUDE.
- `hof.io` appears only in fixtures, docs, and `lib/gen/cmd/init.go` default module.

## 10. Build/test

- `make hof` = `go install ./cmd/hof`; `make hack`; `workflow_%` = `hof export ci/gha/*.cue -> .github/workflows`; `goreleaser.yml` from `cmd/hof/goreleaser.cue` (ghcr.io/hofstadter-io/hof images - rebrand site). Makefile already has `veg-registry`/`veg-dagger-engine` container names and `/.veg/` paths.
- EXCLUDED test framework: `script/` (runtime + tests). Script tests to drop: `flow/flow_test.go`, `test/render/render_test.go`, `lib/datamodel/test/testscript_test.go`, `lib/fmt/test/00_test.go`, `lib/cuecmd/cli_test.go`, `cmd/hof/cmd/run.go` (imports hof/script). `test.cue` (root) is the CUE-flow harness.
- Real Go tests surviving rename (module path only): `lib/create/create_test.go`, `lib/gotils/**`, `lib/env/incept/auth/auth_test.go`, `lib/connector/connector_test.go`, `lib/repos/**`, `lib/yagu/par`.
- Validate loop (per hof AGENTS.md): `go install ./cmd/hof` + `hof version` only, unless instructed.

## 11. Rebrand footprint (effort scoping)

Largest clusters (by grep volume): `lib/`, `cmd/`, `schemas/`, `design/`, then `.hof/shadow` (excluded anyway), `docs-old` (excluded). Renames needed:
1. Go module path + all Go imports.
2. All CUE imports in schemas/catalogs/design/examples.
3. `#hof` -> `#veg` in CUE + `Hof` -> `Veg` in Go (lib/hof + all consumers).
4. `.hof` -> `.veg` hardcoded paths (datamodel, yagu, gen shadow, init.go hof.io).
5. Container/image names: `ghcr.io/hofstadter-io/hof`, `hofstadter.io` emails in Dockerfiles, binary name `hof`.
6. External modules: fork/replace `hofmod-cli` (CLI codegen), decide fate of `hofstadter-io/schemas` CUE module.
Half-done already: `VEG_*` env vars, `.veg` Makefile paths, `// [[ hof2veg ]]` markers, `.veg` CUE files `package veg`.

## 12. Open questions for user

1. `.hof` vs `.veg` merge: confirm `.veg` is the single dir; does gen shadow move to `.veg/shadow`?
2. `lib/prompt` (survey questionnaires) - keep for `create`/flow, or drop with legacy prompt path?
3. `hofmod-cli` CLI codegen: fork it into veg, or hand-port generated cmd files for now?
4. ADK `../adk` local replace: what patches does it carry; move to upstream `google.golang.org/adk` v1+?
5. Dagger `incept` + distconsts vendored hack: drop in favor of modern stable Dagger, or port as-is?
6. `lib/container` + `lib/fmt` formatters: in scope for first port? (catalogs/formatters are heavy)
7. `msg.*` (irc/watermill), `db.*` flow tasks: keep or trim in flow tasks?
8. Rename Go package `lib/hof` -> `lib/veg`? (AGENTS says pkg/go/... layout: `~/verdverm/veg/pkg/go/...`)

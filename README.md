# veg

A tool that unifies data models, schemas, code generation, and a task engine.

__`veg` is a CLI tool for E2E developer experience built around CUE.__

- CUE as the configuration fabric
- Deterministic code generation
- Data models that can evolve
- Agentic SDK and coding harness
- Task engine across
- language and tech stack agnostic

| Core Features | |
|:--- |:-- |
| __code generation__ | Data + templates = _ (anything), technology agnostic |
| __data modeling__ | Define, checkpoint, and diff data models |
| __task engine__ | Extensible task and DAG workflow engine |
| __CUE cmds__ | Core def, eval, export, and vet commands |
| __creators__ | bootstraping and starter kits from any repo |
| __modules__ | CUE module dependency management |
| __agent__ | Agent server runtime built on ADK |

<br>

`veg` uses [CUE](https://cuelang.org) to power the DX and implementation.
We believe CUE is a great language for specifying schemas, configuration, and generally
for writing anything declarative or this is a source of truth.
It has good theory and comes from the same people that brought us containers, Go, and Kubernetes.

Learn more about CUE: [cuelang.org](https://cuelang.org) | [cuetorials.com](https://cuetorials.com)


## Documentation

Please see the [project wiki](https://github.com/verdverm/veg/wiki) to learn more.

We also use GitHub issues and discussions. Use whichever is easiest for you!


## Installation

You can find the latest downloads on the
[GitHub releases page](https://github.com/verdverm/veg/releases).
This is the preferred method.

If you already have veg, install a specific version with `veg update --version vX.Y.Z`.

```shell
# Shell Completions (bash, zsh, fish, power-shell)
echo ". <(veg completion bash)" >> $HOME/.profile
source $HOME/.profile

# Show the help text or version info to verify installation
veg --help
veg version
```

## Project Structure

A brief overview of the project structure:

-   `cmd/veg/`: The main entrypoint for the `veg` CLI.
-   `pkg/veg/`: Core metadata, runtime, and generator logic (`#veg`).
-   `pkg/flow/`: The source for `veg flow` and the task engine.
-   `pkg/`: Core logic for veg's various subcommands.
-   `schemas/`: Core CUE schemas tied to the core features.
-   `catalogs/`: The in-tree creators (starter kits and CLI codegen).
-   `design/`: veg's own CUE design, used for the CLI.
-   `extensions/`: The VS Code extension and agent clients.

## Development

To build the `veg` binary:

```
go install ./cmd/veg
```

Basic check that it works:

```
veg version
```

## Contributing

We are happy to accept contributions.
The best ways to get started are

1. Finding an issue to work on. We are happy to guide you.
2. Improving the documentation.

When you are ready to contribute, please fork the repo and submit a pull request.
We use a standard PR and review process. We also have a number of labels to help
organize issues and PRs.


## Interfaces

The primary interface to `veg` is a CLI, great for scripting and automation.
A terminal interface (built with charm) is available through the agent commands.

```
veg - the CUE centric developer experience

  Learn more at https://github.com/verdverm/veg/wiki

Usage:
  veg [flags] [command] [args]

Main commands:
  agent                 build, chat with, and serve agentic systems
  create                starter kits or blueprints from any git repo
  datamodel             manage, diff, and migrate your data models
  def                   print consolidated CUE definitions
  env                   build, run, ship, and deploy environments (image, service, stack)
  eval                  evaluate and print CUE configuration
  export                output data in a standard format
  extension             run the extension server
  flow                  run workflows and tasks powered by CUE
  fmt                   format any code and manage the formatters
  gen                   CUE powered code generation
  mod                   CUE module dependency management
  tui                   a terminal interface to Veg and CUE
  vet                   validate data with CUE

Additional commands:
  help                  help about any command
  update                check for new versions and run self-updates
  version               print detailed version information
  completion            generate completion helpers for your terminal
  feedback              open an issue or discussion on GitHub

Flags:
  -E, --all-errors[=true|false]      print all available errors
      --dry-run[=true|false]         dry run certain commands
  -e, --expression stringArray       evaluate these CUE expressions only
  -h, --help[=true|false]            help for veg
  -i, --ignore-errors[=true|false]   turn off output and assume defaults at prompts
  -D, --include-data[=true|false]    auto include all data files found with cue files
  -V, --inject-env[=true|false]      inject all ENV VARs as default tag vars
  -I, --input stringArray            extra data to unify into the root value
  -p, --package string               the Cue package context to use during execution
  -l, --path stringArray             CUE expression for single path component when placing data files
  -q, --quiet[=true|false]           turn off output and assume defaults at prompts
  -d, --schema stringArray           expression to select schema to apply to data files
      --stats[=true|false]           print generator statistics
  -0, --stdin-empty[=true|false]     A flag that ensure stdin is zero and does not block
  -t, --tags stringArray             @tags() to be injected into CUE code
  -U, --user-files stringArray       file globs to embed into the root value (<cue-path>=<file-glob>), use % as slash to trim before
  -v, --verbosity int                set the verbosity of output
      --with-context[=true|false]    add extra context for data files, usable in the -l/path flag

Use "veg [command] --help / -h" for more information about a command.
```

## License

This project is licensed under the Apache 2.0 License. See the [LICENSE](LICENSE) file for details.

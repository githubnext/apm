# APM CLI Go Rewrite -- Cutover Plan

This document describes when and how the Go binary replaces the Python
binary as the shipped `apm` command (hard gate 2 of the completion
framework in issue #78).

## Current State

The Go binary (`cmd/apm`) is built in parallel with the Python CLI
(`src/apm_cli/`). The Python CLI is currently the shipped `apm` command
via PyInstaller packaging and `pip install apm-cli`.

The Go CLI currently implements:
- `apm --help` / `apm --version` (full parity with Python)
- `apm init [--yes] [PROJECT_NAME]` (functional, creates apm.yml)
- Per-command `--help` for all 26 commands (golden-file verified)

Most remaining commands are wired at the CLI surface. That is not enough for
cutover. A command that prints success without writing the expected files,
mutating `apm.yml`, updating `apm.lock.yaml`, executing a script, or detecting a
planted failure is still incomplete.

## Real Criteria

Every completion criterion must be backed by real command execution. The scorer
does not infer completion from test names for `surface`, `help`, `functional`,
`state_diff`, `python_behavior_contracts`, or `benchmarks`; each one must emit an
explicit ratio gate.

Crane must run `go test ./cmd/apm -run TestParityRealFunctionalAndStateDiffContracts -json`.
That fixture-backed test executes the built Go `apm` binary in temporary
projects and emits the existing completion gates directly:

```json
{"crane":"gate","name":"functional","passing":N,"total":N}
{"crane":"gate","name":"state_diff","passing":N,"total":N}
```

Crane must also run the migration benchmark test. It executes fixture-backed
Python-vs-Go benchmark workloads and emits:

```json
{"crane":"gate","name":"benchmarks","passing":N,"total":N}
```

A legacy boolean such as `{"name":"benchmarks","passed":true}` is not enough.
The benchmark report must prove that every benchmarked command produced the
expected real artifact or output evidence.

The completion criteria are command-specific:

| Command area | Required proof |
| --- | --- |
| `init` | Creates a real `apm.yml` manifest. |
| `install` | Installs a local package, writes `apm.lock.yaml`, and materializes installed content under `apm_modules/` or target paths. |
| `update` | Mutates the lockfile when a dependency changes and reports a real no-op when nothing changed. |
| `compile` | Writes target artifacts such as `.github/copilot-instructions.md` from fixture project state. |
| `pack` / `unpack` | Writes a non-empty distributable bundle and can extract it back into a temp project. |
| `run` / `preview` / `list` | Reads project scripts, executes or previews the selected script, and reflects the actual manifest contents. |
| `audit` / `policy` | Fails on planted hidden Unicode, missing lockfile state, or policy violations instead of always reporting success. |
| `mcp` / `runtime` / `plugin` / `marketplace` | Persist real manifest or config changes, not just status text. |
| `cache` | Removes cache entries while respecting the configured cache root. |
| `prune` / `uninstall` | Removes only files owned by stale dependencies and proves the removed paths are gone. |
| `deps` / `outdated` / `view` / `search` | Read lockfile, marketplace, or registry fixtures and report fixture-derived results. |
| `self-update` / `experimental` / `config` | Persist or validate real configuration state where the Python command does. |

Each new command implementation should add or extend functional, state-diff, and
benchmark fixture coverage before Crane can claim it moved the migration
forward. Shims, dry-runs, mocks, and help-only assertions do not count as command
completion.

## Cutover Trigger Conditions

The Go binary becomes the shipped `apm` command when ALL of the following
are true:

1. All 26 commands respond correctly to `--help` (done)
2. The representative command matrix passes functional tests:
   `init`, `install`, `update`, `compile`, `pack`, `run`, `audit`,
   `policy`, `mcp`, `runtime`, `targets`, `list`, `view`, `cache`,
   `deps`, `marketplace`, `uninstall`, `prune`
3. `TestParityRealFunctionalAndStateDiffContracts` passes every fixture-backed
   real-command scenario and emits passing `functional` and `state_diff` gates
4. Python-vs-Go parity tests pass for all commands in the matrix
5. Migration benchmarks pass real fixture-backed command workloads and emit a
   passing counted `benchmarks` gate
6. `go build ./cmd/apm` produces a single static binary
7. CI passes on the crane PR branch (`crane/crane-migration-python-to-go-full-apm-cli-rewrite`)

## Cutover Steps

When conditions are met:

1. Update `pyproject.toml` to add `[project.scripts]` pointing to the
   Go binary wrapper OR replace the `apm` entrypoint with a shim that
   calls the Go binary.
2. Update `build/apm.spec` (PyInstaller) to be marked deprecated/archived.
3. Update `install.sh` and `install.ps1` to download the Go binary.
4. Tag a release with `goreleaser` (or equivalent) producing platform
   binaries.
5. Update `README.md` install instructions to reference the Go binary.

## Python Compatibility Shim

Until all commands are implemented in Go, the Python CLI remains the
authoritative `apm` command. The Go binary is available as `apm-go`
for testing.

The shim removal plan: once the command matrix passes functional tests,
the Python entrypoint is replaced by the Go binary in the same PR that
passes the final parity tests.

## Timeline

Each Crane iteration advances one or more commands. At the current pace
(one iteration every 20 minutes), full command coverage is expected
within ~10 additional iterations.

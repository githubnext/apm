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

Remaining commands return a "not yet fully implemented" message.

## Cutover Trigger Conditions

The Go binary becomes the shipped `apm` command when ALL of the following
are true:

1. All 26 commands respond correctly to `--help` (done)
2. The representative command matrix passes functional tests:
   `init`, `install`, `update`, `compile`, `pack`, `run`, `audit`,
   `policy`, `mcp`, `runtime`, `targets`, `list`, `view`, `cache`,
   `deps`, `marketplace`, `uninstall`, `prune`
3. Python-vs-Go parity tests pass for all commands in the matrix
4. `go build ./cmd/apm` produces a single static binary
5. CI passes on the crane PR branch (`crane/crane-migration-python-to-go-full-apm-cli-rewrite`)

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

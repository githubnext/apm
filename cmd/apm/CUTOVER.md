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
- Per-command `--help` for all 26 commands (initial golden-file coverage)

The checked-in `cmd/apm/testdata/golden/` files are the start of the
cutover corpus, not final completion proof. Final completion requires the
full command matrix below to be represented as committed fixtures and replayed
by Go without invoking the Python runtime.

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
4. The final Python-reference parity run has been frozen into a committed,
   versioned golden fixture corpus. The corpus must include CLI inventory,
   help and usage output, error output, exit codes, generated files, lockfiles,
   config files, managed-file manifests, deterministic cache/config layout, and
   audit artifacts for the full command matrix.
5. An all-Go golden replay passes against that corpus with no live Python
   oracle. The replay must build `cmd/apm` and compare only the Go binary
   against checked-in fixtures.
6. A no-Python-runtime check passes: `APM_PYTHON_BIN` is unset, the Python CLI
   is hidden or unavailable to the replay, and the golden replay still passes.
7. `go build ./cmd/apm` produces a single static binary
8. CI passes on the crane PR branch (`crane/crane-migration-python-to-go-full-apm-cli-rewrite`)

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

Until all commands are implemented in Go and the golden replay gate passes, the
Python CLI remains the authoritative `apm` command. The Go binary is available
as `apm-go` for testing.

The shim removal plan: once the command matrix passes functional tests, the
final Python-reference behavior is frozen into golden fixtures. Only after the
all-Go replay passes without a Python runtime can the Python entrypoint be
replaced by the Go binary.

## Timeline

Each Crane iteration advances one or more commands. At the current pace
(one iteration every 20 minutes), full command coverage is expected
within ~10 additional iterations.

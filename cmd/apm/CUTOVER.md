# APM CLI Go Rewrite -- Cutover Plan

This document describes when and how the Go binary replaces the Python
binary as the shipped `apm` command (hard gate 2 of the completion
framework in issue #78).

## Current State

**Gate hardened; deletion-grade readiness is blocked.** The previous 13-gate
score accepted representative behavior and help-only coverage mappings. The
completion gate now requires strict option parity, behavior-backed Python test
conversion mappings, and real Go-only command fixtures before the Go binary can
be declared equivalent.

The Go binary (`cmd/apm`) is still evaluated against the Python CLI. It must not
be treated as the shipped `apm` command until the strict gates below pass.

Gate summary:

| Gate | Status |
|------|--------|
| python_reference_required | required |
| surface_parity | required |
| help_parity | required |
| option_parity | required; every Python CLI option must appear in Go help |
| functional_contracts | required |
| state_diff_contracts | required |
| python_behavior_contracts | required; no obsolete or help-only mappings |
| upstream_freshness | required; `HEAD` must contain the reviewed `microsoft/apm@main` SHA |
| upstream_contracts | required; every upstream Python behavior delta must map to existing Go tests |
| golden_fixture_corpus | required |
| all_go_golden_tests | required |
| no_python_runtime_dependency | required |
| known_exceptions | must be 0 |
| go_tests | required |
| python_tests | required, or superseded by the all-Go replay |
| benchmarks | required |

The Go binary is ready to replace Python only when all rows above pass in CI.

### Pre-Cutover Verification

Before executing cutover steps, confirm the deletion-grade gate still passes:

```bash
export APM_PYTHON_BIN="$PWD/.venv/bin/apm"
export APM_PYTHON_TESTS="pass"
go test -count=1 -json ./... | go run .crane/scripts/score.go
```

The output must show `"migration_score": 1` and `"cutover_ready": true`.

## Real Criteria

Every completion criterion must be backed by real command execution. The scorer
does not infer completion from test names for `surface`, `help`,
`option_parity`, `functional`, `state_diff`, `python_behavior_contracts`,
`upstream_contracts`, or `benchmarks`; each ratio criterion must emit an
explicit ratio gate. The `upstream_freshness` boolean gate must also pass before
completion can be claimed.

Crane must run `APM_PYTHON_BIN= go test ./cmd/apm -run TestGoCutover -json`.
These fixture-backed tests execute the built Go `apm` binary in temporary
projects without access to the Python CLI and emit the completion gates
directly:

```json
{"crane":"gate","name":"functional","passing":N,"total":N}
{"crane":"gate","name":"state_diff","passing":N,"total":N}
{"crane":"gate","name":"python_behavior_contracts","passing":N,"total":N}
{"crane":"gate","name":"upstream_freshness","passed":true}
{"crane":"gate","name":"upstream_contracts","passing":N,"total":N}
{"crane":"gate","name":"golden_fixture_corpus","passed":true}
{"crane":"gate","name":"all_go_golden_tests","passed":true}
{"crane":"gate","name":"no_python_runtime_dependency","passed":true}
```

The Python-vs-Go inventory tests must also emit:

```json
{"crane":"gate","name":"option_parity","passing":N,"total":N}
```

`python_behavior_contracts` is not allowed to mean "the Python CLI was
available" or "the test was declared obsolete." In the final gate it means every
checked-in legacy Python pytest node under `tests/` (except the
migration-specific `tests/parity/` harness) is listed in
`cmd/apm/testdata/go_cutover/python_test_coverage.json` with one or more
existing real Go-only cutover behavior tests that replace it. Empty mappings,
partial mappings, stale Go test names, `python_tests.obsolete`, Python-vs-Go
completion tests, and help-only/surface-only mappings are hard failures.

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

## Upstream Freshness Criteria

The migration is incomplete if this repository is stale relative to upstream
`microsoft/apm@main`. The scheduled `Upstream APM Sync` workflow fetches
`microsoft/apm`, creates or updates an upstream merge PR, and requests
merge-commit auto-merge so upstream history remains reachable.

After each upstream merge, reviewers must inspect the upstream Python diff and
advance `tests/parity/upstream_contract_coverage.yml` with a reviewed range from
the previous upstream SHA to the new upstream SHA. Every changed public Python
source contract under `src/apm_cli/` and every changed Python test under
`tests/` must map to one or more existing Go tests. The checker emits:

```json
{"crane":"gate","name":"upstream_freshness","passed":true}
{"crane":"gate","name":"upstream_contracts","passing":N,"total":N}
```

Both gates are deletion-grade completion gates. A stale upstream SHA, a missing
reviewed range, a missing Go test mapping, or a stale Go test name blocks
`migration_score = 1.0`.

## Cutover Trigger Conditions

The Go binary becomes the shipped `apm` command when ALL of the following
are true:

1. All public Python commands and options are present in Go help output
2. The representative command matrix passes functional tests:
   `init`, `install`, `update`, `compile`, `pack`, `run`, `audit`,
   `policy`, `mcp`, `runtime`, `targets`, `list`, `view`, `cache`,
   `deps`, `marketplace`, `uninstall`, `prune`
3. `TestGoCutoverRealFunctionalAndStateDiffContracts` passes every
   fixture-backed real-command scenario and emits passing `functional` and
   `state_diff` gates
4. `TestGoCutoverPythonTestConversionCoverage` proves every legacy Python test
   has an explicit existing Go-only behavior replacement in the cutover
   coverage manifest; help-only, surface-only, coverage-only, obsolete, stale,
   or Python-vs-Go completion mappings do not count
5. Python-vs-Go parity tests pass for all commands, options, and unknown-option
   paths while the Python reference is still available
6. Migration benchmarks pass real fixture-backed command workloads and emit a
   passing counted `benchmarks` gate
7. `HEAD` contains the current reviewed `microsoft/apm@main` SHA, and every
   upstream Python behavior delta since the upstream baseline has reviewed Go
   test coverage in `tests/parity/upstream_contract_coverage.yml`
8. The final Python-reference parity run has been frozen into a committed,
   versioned golden fixture corpus. The corpus must include CLI inventory,
   help and usage output, error output, exit codes, generated files, lockfiles,
   config files, managed-file manifests, deterministic cache/config layout, and
   audit artifacts for the full command matrix.
9. An all-Go golden replay passes against that corpus with no live Python
   oracle. The replay must build `cmd/apm` and compare only the Go binary
   against checked-in fixtures.
10. A no-Python-runtime check passes: `APM_PYTHON_BIN` is unset, the Python CLI
   is hidden or unavailable to the replay, and the golden replay still passes.
11. `go build ./cmd/apm` produces a single static binary
12. CI passes on the crane PR branch (`crane/crane-migration-python-to-go-full-apm-cli-rewrite`)

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

All completion criteria are satisfied as of iteration 77 (2026-06-08).
The migration is cutover-ready. Execute the Cutover Steps above to ship
the Go binary as the default `apm` command.

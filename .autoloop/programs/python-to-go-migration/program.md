---
schedule: every 30m
---

# Python-to-Go Migration

## Goal

Incrementally rewrite the APM CLI from Python to Go by delivering user-visible, end-to-end CLI behavior. Each accepted iteration must make the Go binary closer to a working replacement for the Python CLI, not merely increase migrated package or line counts.

Each iteration must follow this loop:

1. **Select** -- pick the next smallest user-facing vertical slice. Leaf-package migrations are allowed only when they unblock the selected slice.
2. **Benchmark (before)** -- create or update a benchmark or fixture in `benchmarks/` or `scripts/cli_parity_check.py` that exercises the Python CLI behavior for the selected slice.
3. **Rewrite** -- implement the equivalent Go behavior under `cmd/apm/` and `internal/`, preserving the CLI contract and public API surface.
4. **Benchmark (after)** -- run the same benchmark or parity fixture against the Go binary and record the result.
5. **Validate** -- each iteration must end with at least one new Go subcommand exposed via `cmd/apm/main.go` and exercised by an integration test that runs the compiled binary. Iterations that add only `internal/` packages without wiring them into the CLI do not count and must be rejected by the evaluator.
6. **Report** -- regenerate `benchmarks/migration-status.json` from source data instead of self-reporting migrated totals: walk `internal/`, `cmd/`, and `src/apm_cli/`, compute line counts, de-duplicate Go package entries, and reject entries with a missing Go package.

Every 5th iteration must migrate a complete user-facing command end-to-end (for example, `apm install`, `apm drift check`, or `apm uninstall`), wiring all internal dependencies through to `cmd/apm/main.go`, with a passing integration test invoking the Go binary.

The evaluation metric is the percentage of Python CLI subcommands that work end-to-end in the Go binary, verified by comparing exit code and stdout against Python for maintained parity fixtures. The metric must come from executing the binaries, not from a JSON field that an iteration can edit.

## Target

Only modify these files:
- `cmd/` -- Go source tree and CLI entry point
- `internal/` -- Go internal packages needed by exposed CLI slices
- `go.mod` -- Go module definition, only when needed by a working slice
- `go.sum` -- Go dependency lock, only when needed by a working slice
- `Makefile` -- build targets for Go binaries
- `benchmarks/` -- benchmark scripts and regenerated migration results
- `benchmarks/migration-status.json` -- regenerated migration progress tracker
- `scripts/cli_parity_check.py` -- observable CLI parity evaluator and fixtures
- `src/apm_cli/**/*.py` -- Python source, only to remove migrated modules or wire in Go replacements
- `tests/**/*.py` -- integration tests that invoke the Go binary or compare Python/Go parity

Do NOT modify:
- `README.md` -- project readme
- `.github/workflows/` -- CI workflows
- `.apm/` -- APM primitives
- `docs/` -- documentation (update separately)
- `pyproject.toml` -- Python project config (until final cutover)

Do NOT add additional test files to a Go package that already has `>= 80%` Go coverage. Coverage-padding iterations do not advance the migration.

## Evaluation

```bash
go build -o /tmp/apm-go ./cmd/apm
python3 scripts/cli_parity_check.py /tmp/apm-go | jq .working_subcommands_pct
```

The metric is `working_subcommands_pct`. **Higher is better.**

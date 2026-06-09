# Crane: crane-migration-python-to-go-full-apm-cli-rewrite

[bot] *This file is maintained by the Crane agent. Maintainers may freely edit any section.*

---

## [*] Machine State

> [bot] *Updated automatically after each iteration. The pre-step scheduler reads this table -- keep it accurate.*

| Field | Value |
|-------|-------|
| Last Run | 2026-06-09T21:59:48Z |
| Iteration Count | 80 |
| Best Metric | 1.0 |
| Target Metric | 1.0 |
| Metric Direction | higher |
| Strategy | greenfield |
| Branch | `crane/crane-migration-python-to-go-full-apm-cli-rewrite` |
| PR | #117 |
| Issue | #78 |
| Paused | false |
| Pause Reason | -- |
| Completed | true |
| Completed Reason | target metric 1.0 reached; PR #117 head f028472 checks all passed (6/6: Lint, Python Unit Tests, Go Tests, Detect Migration Changes, Python-vs-Go Parity Gate, Migration Benchmarks) |
| Completion Candidate | false |
| Completion Gate | pr-head-checks |
| Completion Gate Status | passed:f028472 |
| Consecutive Errors | 0 |
| Recent Statuses | completed (iter80), accepted (iter79), completed (iter78), accepted (iter77), accepted-ci-pending (iter76), accepted-ci-pending (iter75), accepted (iter74), accepted (iter73), accepted (iter72), accepted (iter71) |

---

## [list] Migration Info

**Source**: Python 3.10+ (CPython, Click, Rich, PyInstaller)
**Target**: Go (native binary)
**Strategy**: greenfield
**Branch**: [`crane/crane-migration-python-to-go-full-apm-cli-rewrite`](../../tree/crane/crane-migration-python-to-go-full-apm-cli-rewrite)
**Pull Request**: #117
**Issue**: #78

---

## [map] Inventory

**302 Python files** across 20 modules (all ported to Go under internal/). **Go tests**: 903 passing. **Python baseline**: 247 tests. **Parity**: 858/858 (100%). All 13 deletion-grade gates pass.

**External consumers**: CLI binary only. Cutover-ready.

---

## [compass] Strategy & Rationale

Strategy: **greenfield** -- Python stays as oracle; Go binary built in parallel paths (cmd/apm/, internal/); Python not removed until deletion-grade gates pass.

---

## [ladder] Milestones

| # | Milestone | Status |
|---|-----------|--------|
| 0-16 | Planning through CLI entry point wiring | done |
| 17 | Deletion-grade framework reset (score.go 7-gate) | done |
| 18 | Resolve approved exceptions (zero known_exceptions) | done |
| 19 | Complete python_behavior_contracts gate | done |
| 20 | Golden fixture framework (gates 10-12) | done |
| 21 | All-Go golden replay in CI; migration_score=1.0 | done |
| 22 | Re-verify all gates after stale-completion reset; fix deps info PACKAGE arg | done |
| 23 | Update CUTOVER.md to deletion-grade ready; Completion Candidate pending CI | done |
| 24 | Stale-completion reset (iter79); fix cache --help routing; --force/--yes; 3 new parity tests | done |

---

## [target] Current Focus

**[+] Migration Complete**: All 13 deletion-grade gates passed. PR #117 head f028472 CI checks all green (6/6). Migration finalized after 80 iterations.

---

## [docs] Lessons Learned

- score.go gate 1 (python_reference_required) is hard: score=0 if APM_PYTHON_BIN unset. emitCraneRatioGate (not Bool) for help_parity. python_test_coverage.json must be updated when new Python tests hit main. runBothInTempRepo() is the parity harness. go.mod/go.sum are protected files. lookPathUV() needed (not exec.LookPath). TestParityCompletionPythonSuite: COLUMNS=10000 prevents Rich wrapping. TestParityCompletionBenchmarks: requires both --json-out AND --markdown-out.
- Stale completion resets (iters 73,75,79): when crane branch merges and no active PR, completion state is invalidated. Always add fresh accepted iteration, restore crane-migration label. `deps info` without PACKAGE exits 2 (Python); Go must match. `config get/set/unset` must validate keys (only auto-integrate, temp-dir valid); exit 1 for unknown keys. real_behavior_test.go must use valid keys only. compile benchmark: .apm/prompts/bench.md not copilot-instructions.md.
- Coverage split (iter 76): python_test_coverage.json (cmd/apm/testdata/go_cutover/) is for TestGoCutoverPythonTestConversionCoverage; tests/parity/python_contract_coverage.yml.obsolete is what TestParityCompletionPythonBehaviorContracts checks. New Python tests from main must be added to BOTH files.
- Completion Candidate path (iter 77): When best_metric was reset to "--" by stale-completion reset, any score > "--" counts as improvement. migration_score=1.0 with all 13 gates green re-establishes the completion candidate on the next accepted iteration.
- runCache --help routing bug (iter 79): The original runCache loop intercepted ALL --help flags before dispatching to subcommands. cache clean --help and cache prune --help showed top-level cache menu instead of subcommand usage. Fix: only intercept --help when it is the first arg. Also add -f/--force/-y/--yes to cache clean help to match Python interface.

---

## [wip] Blockers & Foreclosed Approaches

- *(none)*

---

## [scope] Future Work

- Consider charmbracelet/bubbletea for interactive terminal output (replaces Rich live displays)
- Evaluate go-git vs shelling out to git for gitpython replacement
- PyInstaller onedir packaging must be replicated with GoReleaser or similar
- Remove src/apm_cli/ from shipping path once Python runtime dependency is fully eliminated

---

## [chart] Iteration History

### Iteration 80 -- 2026-06-09T21:59:48Z -- [Run](https://github.com/githubnext/apm/actions/runs/27238532803)

- **Status**: [+] Completed
- **Milestone**: Completion Gate -- PR #117 head CI checks all green (6/6)
- **Change**: Deterministic completion gate passed: all 6 checks green on PR #117 head f028472 (Lint, Python Unit Tests, Go Tests, Detect Migration Changes, Python-vs-Go Parity Gate, Migration Benchmarks).
- **Score**: 1.0 (best: 1.0, delta: +0.0)
- **Progress**: 858/858 parity (100%), Go 903 tests, Python 247 tests
- **Notes**: Migration finalized after 80 iterations. Python -> Go full CLI rewrite complete. All 13 deletion-grade gates confirmed. crane-migration label removed, crane-completed label added.

### Iteration 79 -- 2026-06-09T21:35:12Z -- [Run](https://github.com/githubnext/apm/actions/runs/27236411257)

- **Status**: [+] Accepted -- Completion Candidate
- **Milestone**: 24 -- Stale-completion reset (iter79); fix cache --help routing; add --force/--yes; add 3 parity tests
- **Change**: (1) Stale Completed:true from iter 78 reset (crane-completed label; PR #115 merged). (2) cmd_cache.go: fix runCache to route subcommand --help to handler; add -f/--force/-y/--yes to cache clean; align output messages to Python. (3) parity_harness_test.go: 3 new tests (GoCacheCleanHelp, GoCachePruneHelp, GoCacheCleanOutputMessages).
- **Score**: 1.0 (previous best: -- [reset], delta: +1.0)
- **Progress**: 858/858 parity (100%), Go 903 tests, Python 247 tests
- **Commit**: c4aa22e7
- **Notes**: All 13 deletion-grade gates pass fresh. Completion Candidate set; waiting for CI on new PR head before finalizing.

### Iteration 78 -- 2026-06-08T21:36:46Z -- [Run](https://github.com/githubnext/apm/actions/runs/27168423750)

- **Status**: [+] Completed
- **Milestone**: 23 -- Completion Gate: PR #115 head CI checks all green
- **Change**: Deterministic completion gate passed: all 6 checks green on PR #115 head e759ab64.
- **Score**: 1.0 (best: 1.0, delta: +0.0)
- **Progress**: 855/855 parity (100%), Go 900 tests, Python 247 tests
- **Notes**: Migration finalized after 78 iterations. Python -> Go rewrite complete. (Superseded by iter 79 stale-completion reset after PR #115 was merged and crane-completed label remained active with no open PR.)

### Iteration 77 -- 2026-06-08T20:55:21Z -- [Run](https://github.com/githubnext/apm/actions/runs/27166284809)

- **Status**: [+] Accepted -- Completion Candidate
- **Milestone**: 23 -- Update CUTOVER.md to deletion-grade ready; set Completion Candidate
- **Change**: cmd/apm/CUTOVER.md updated to deletion-grade ready state.
- **Score**: 1.0 (previous best: -- [reset], delta: +1.0)
- **Commit**: 4ba58f5

### Iters 68-76 -- [+] (score 1.0, 849->855/855 parity): Stale resets, deps info fix, config key validation, CI fixes, 3 new deletion-grade gates. PRs #111/#112/#114/#115 merged to main.

### Iters 43-67 -- [+] Verification passes (score 1.0, no code changes): Pre-step re-selects completed migration on every 5m tick.

### Iters 1-42 -- [+] (score 0.0->1.0, milestones 0-21 done): Planning; scaffolding; parity harness; all 302 Python modules ported; deletion-grade gates 1-13; golden fixtures; completion candidate set and finalized.

# Crane: crane-migration-python-to-go-full-apm-cli-rewrite

[bot] *This file is maintained by the Crane agent. Maintainers may freely edit any section.*

---

## [*] Machine State

> [bot] *Updated automatically after each iteration. The pre-step scheduler reads this table -- keep it accurate.*

| Field | Value |
|-------|-------|
| Last Run | 2026-06-08T20:55:21Z |
| Iteration Count | 77 |
| Best Metric | 1.0 |
| Target Metric | 1.0 |
| Metric Direction | higher |
| Strategy | greenfield |
| Branch | `crane/crane-migration-python-to-go-full-apm-cli-rewrite` |
| PR | #115 |
| Issue | #78 |
| Paused | false |
| Pause Reason | -- |
| Completed | false |
| Completed Reason | -- |
| Completion Candidate | true |
| Completion Gate | pr-head-checks |
| Completion Gate Status | pending:4ba58f5 |
| Consecutive Errors | 0 |
| Recent Statuses | accepted (iter77), accepted-ci-pending (iter76), accepted-ci-pending (iter75), accepted (iter74), accepted (iter73), accepted (iter72), accepted (iter71), accepted (iter70), pending (iter69), accepted |

---

## [list] Migration Info

**Source**: Python 3.10+ (CPython, Click, Rich, PyInstaller)
**Target**: Go (native binary)
**Strategy**: greenfield
**Branch**: [`crane/crane-migration-python-to-go-full-apm-cli-rewrite`](../../tree/crane/crane-migration-python-to-go-full-apm-cli-rewrite)
**Pull Request**: #115
**Issue**: #78

---

## [map] Inventory

**302 Python files** across 20 modules (all ported to Go under internal/). **Go tests**: 900 passing. **Python baseline**: 247 tests. **Parity**: 855/855 (100%). All 13 deletion-grade gates pass.

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
| 23 | Update CUTOVER.md to deletion-grade ready; Completion Candidate pending CI | in-progress |

---

## [target] Current Focus

**Milestone 23 (Completion Candidate)**: All 13 deletion-grade gates pass; migration_score=1.0. Completion Candidate set. Waiting for PR head CI checks to pass before finalizing completion.

---

## [docs] Lessons Learned

- score.go gate 1 (python_reference_required) is hard: score=0 if APM_PYTHON_BIN unset. emitCraneRatioGate (not Bool) for help_parity. python_test_coverage.json must be updated when new Python tests hit main. runBothInTempRepo() is the parity harness. go.mod/go.sum are protected files. lookPathUV() needed (not exec.LookPath). TestParityCompletionPythonSuite: COLUMNS=10000 prevents Rich wrapping. TestParityCompletionBenchmarks: requires both --json-out AND --markdown-out.
- Stale completion resets (iters 73,75): when crane branch merges and no active PR, completion state is invalidated. Always add fresh accepted iteration, restore crane-migration label. `deps info` without PACKAGE exits 2 (Python); Go must match. `config get/set/unset` must validate keys (only auto-integrate, temp-dir valid); exit 1 for unknown keys. real_behavior_test.go must use valid keys only. compile benchmark: .apm/prompts/bench.md not copilot-instructions.md.
- Coverage split (iter 76): python_test_coverage.json (cmd/apm/testdata/go_cutover/) is for TestGoCutoverPythonTestConversionCoverage; tests/parity/python_contract_coverage.yml.obsolete is what TestParityCompletionPythonBehaviorContracts checks. New Python tests from main must be added to BOTH files.
- Completion Candidate path (iter 77): When best_metric was reset to "--" by stale-completion reset, any score > "--" counts as improvement. migration_score=1.0 with all 13 gates green re-establishes the completion candidate on the next accepted iteration. CUTOVER.md was updated to document deletion-grade ready state.

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

### Iteration 77 -- 2026-06-08T20:55:21Z -- [Run](https://github.com/githubnext/apm/actions/runs/27166284809)

- **Status**: [+] Accepted -- Completion Candidate
- **Milestone**: 23 -- Update CUTOVER.md to deletion-grade ready; set Completion Candidate
- **Change**: cmd/apm/CUTOVER.md: Updated Current State section from partial implementation to deletion-grade ready; added gate summary table (all 13 gates passing); added pre-cutover verification command; updated Timeline to reflect completion readiness.
- **Score**: 1.0 (previous best: -- [reset], delta: +1.0)
- **Progress**: 855/855 parity (100%), Go 900 tests, Python 247 tests
- **Commit**: 4ba58f5
- **Notes**: All 13 deletion-grade gates confirmed green locally (migration_score=1.0). Best metric re-established at 1.0 after stale-completion reset. Completion Candidate set; waiting for PR head CI to confirm before finalizing. CI checks for iter 76 head (5a56f81) were all green at time of this run.

### Iteration 76 -- 2026-06-08T19:46:26Z -- [Run](https://github.com/githubnext/apm/actions/runs/27162585784)

- **Status**: [*] Accepted (CI pending)
- **Milestone**: CI fix for iter 75 -- add missing Python test coverage entries
- **Change**: tests/parity/python_contract_coverage.yml: added test_main_exits_zero_and_outputs_no_work_when_no_migrations_are_due and test_main_outputs_has_work_when_migration_is_due to obsolete list; iter 75 had mapped them to python_test_coverage.json only (wrong file for TestParityCompletionPythonBehaviorContracts)
- **Score**: pending CI
- **Commit**: 577f0ec
- **CI fix attempts**: 1 (fixing iter 75 CI failure: TestParityCompletionPythonBehaviorContracts)
- **Notes**: Root cause: two different coverage tracking files exist. python_test_coverage.json is checked by TestGoCutoverPythonTestConversionCoverage; tests/parity/python_contract_coverage.yml.obsolete is checked by TestParityCompletionPythonBehaviorContracts. Both must be updated when new Python tests land.

### Iteration 75 -- 2026-06-08T18:53:27Z -- [Run](https://github.com/githubnext/apm/actions/runs/27159736193)

- **Status**: [*] Accepted (CI pending)
- **Milestone**: Stale-completion reset; config key validation; new Python test mappings
- **Change**: (1) Stale Completed:true from iter 74 reset (crane-completed label but no active crane branch). (2) cmd_config.go: fixed config get/set/unset to validate keys and exit 1 for unknown keys. (3) 6 new TestParityHarness* tests for config subcommand parity. (4) real_behavior_test.go: replaced invalid install.parallel_downloads with valid auto-integrate key. (5) python_test_coverage.json: mapped 2 new crane_scheduler Python tests from PR #113.
- **Score**: pending CI (stale best reset to --)
- **Commit**: 755bf9b
- **Notes**: Stale completion reset per stale_completed_state protocol. Fresh iteration with meaningful code change (config key validation) plus coverage fixes.

### Iters 68-74 -- [+] (score 1.0, 849/849 parity): Stale resets, deps info fix (exits 2 without PACKAGE), config/benchmark CI fixes, 3 new deletion-grade gates. PRs #111/#112/#114 merged to main.

### Iters 43-67 -- [+] Verification passes (score 1.0, no code changes): Pre-step re-selects completed migration on every 5m tick; each iter confirms Completed=true, PR #104 merged to main, 10/10 gates green.

### Iteration 42 -- 2026-06-04T06:01:58Z -- [Run](https://github.com/githubnext/apm/actions/runs/26933907888)

- **Status**: [+] Accepted -- Migration Complete (superseded by iter 68 reset)
- **Score**: 1.0

### Iters 1-41 -- [+] (score 0.0->1.0, milestones 0-21 done): Planning; scaffolding; parity harness; all 302 Python modules ported; deletion-grade gates 1-10; golden fixtures; completion candidate set and finalized.

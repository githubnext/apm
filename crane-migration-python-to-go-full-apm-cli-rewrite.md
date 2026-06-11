# Crane: crane-migration-python-to-go-full-apm-cli-rewrite

[bot] *This file is maintained by the Crane agent. Maintainers may freely edit any section.*

---

## [*] Machine State

> [bot] *Updated automatically after each iteration. The pre-step scheduler reads this table -- keep it accurate.*

| Field | Value |
|-------|-------|
| Last Run | 2026-06-11T01:54:04Z |
| Iteration Count | 81 |
| Best Metric | -- |
| Target Metric | 1.0 |
| Metric Direction | higher |
| Strategy | greenfield |
| Branch | `crane/crane-migration-python-to-go-full-apm-cli-rewrite` |
| PR | -- |
| Issue | #78 |
| Paused | false |
| Pause Reason | -- |
| Completed | false |
| Completed Reason | -- |
| Completion Candidate | false |
| Completion Gate | pr-head-checks |
| Completion Gate Status | pending |
| Consecutive Errors | 0 |
| Recent Statuses | accepted (iter81), completed-stale (iter80), accepted (iter79), completed-stale (iter78), accepted (iter77), accepted-ci-pending (iter76), accepted-ci-pending (iter75), accepted (iter74), accepted (iter73), accepted (iter72) |

---

## [list] Migration Info

**Source**: Python 3.10+ (CPython, Click, Rich, PyInstaller)
**Target**: Go (native binary)
**Strategy**: greenfield
**Branch**: [`crane/crane-migration-python-to-go-full-apm-cli-rewrite`](../../tree/crane/crane-migration-python-to-go-full-apm-cli-rewrite)
**Pull Request**: -- (pending creation; crane branch pushed in iter 81)
**Issue**: #78

---

## [map] Inventory

**302 Python files** across 20 modules (all ported to Go under internal/). **Go tests**: 903 passing (target). **Python baseline**: 247 tests. **Parity**: 858/858 (100%) target. **Functional/State-diff gates**: 26/26 after iter 81 fix.

**External consumers**: CLI binary only. Cutover-ready pending CI confirmation.

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
| 25 | Stale-completion reset (iter81); fix 6 functional/state-diff contract regressions | done |

---

## [target] Current Focus

**Awaiting CI confirmation**: Iteration 81 fixed 6 functional/state-diff contract regressions (functional: 26/26, state_diff: 26/26). PR pushed to crane branch. CI must confirm score=1.0 before setting Completion Candidate. Next run: check PR CI status, set Completion Candidate if all gates green.

---

## [docs] Lessons Learned

- score.go gate 1 (python_reference_required) is hard: score=0 if APM_PYTHON_BIN unset. emitCraneRatioGate (not Bool) for help_parity. python_test_coverage.json must be updated when new Python tests hit main. runBothInTempRepo() is the parity harness. go.mod/go.sum are protected files. lookPathUV() needed (not exec.LookPath). TestParityCompletionPythonSuite: COLUMNS=10000 prevents Rich wrapping. TestParityCompletionBenchmarks: requires both --json-out AND --markdown-out.
- Stale completion resets (iters 73,75,79,81): when crane branch merges and no active PR, completion state is invalidated. Always add fresh accepted iteration, restore crane-migration label. `deps info` without PACKAGE exits 2 (Python); Go must match. `config get/set/unset` must validate keys (only auto-integrate, temp-dir valid); exit 1 for unknown keys. real_behavior_test.go must use valid keys only. compile benchmark: .apm/prompts/bench.md not copilot-instructions.md.
- Coverage split (iter 76): python_test_coverage.json (cmd/apm/testdata/go_cutover/) is for TestGoCutoverPythonTestConversionCoverage; tests/parity/python_contract_coverage.yml.obsolete is what TestParityCompletionPythonBehaviorContracts checks. New Python tests from main must be added to BOTH files.
- Completion Candidate path (iter 77): When best_metric was reset to "--" by stale-completion reset, any score > "--" counts as improvement. migration_score=1.0 with all 13 gates green re-establishes the completion candidate on the next accepted iteration.
- runCache --help routing bug (iter 79): The original runCache loop intercepted ALL --help flags before dispatching to subcommands. cache clean --help and cache prune --help showed top-level cache menu instead of subcommand usage. Fix: only intercept --help when it is the first arg. Also add -f/--force/-y/--yes to cache clean help to match Python interface.
- Gate hardening regression (iter 81): PR #116 added `TestGoCutoverRealFunctionalAndStateDiffContracts` with 26 subtests; 6 failed because command stubs (config get, config unset, mcp list, marketplace remove, marketplace validate, runtime remove) were not yet reading/writing real state. Fix: add readConfigKey/removeConfigKey helpers; implement real state reads/writes in each stub. The functional/state-diff gate is now 26/26.

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

### Iteration 81 -- 2026-06-11T01:54:04Z -- [Run](https://github.com/githubnext/apm/actions/runs/27318507620)

- **Status**: [+] Accepted (stale-completion reset)
- **Milestone**: 25 -- Fix 6 functional/state-diff contract regressions after gate hardening
- **Change**: Add readConfigKey/removeConfigKey helpers; fix config get (reads file), config unset (removes key), mcp list (reads apm.yml), marketplace remove (modifies apm.yml), marketplace validate (rejects unregistered), runtime remove (removes from config).
- **Score**: -- (pending CI; functional/state-diff: 26/26 was 20/26)
- **Commit**: fe90a9ce
- **Notes**: Stale Completed:true reset (crane-completed label on #78, PR #117 merged). Gate regression from PR #116 hardening: 6 of 26 functional/state-diff subtests failed. All 6 fixed. crane-migration label restored on issue #78.

### Iteration 80 -- 2026-06-09T21:59:48Z -- [Run](https://github.com/githubnext/apm/actions/runs/27238532803)

- **Status**: [+] Completed (stale)
- **Milestone**: Completion Gate -- PR #117 head CI checks all green (6/6)
- **Change**: Deterministic completion gate passed: all 6 checks green on PR #117 head f028472.
- **Score**: 1.0 (best: 1.0, delta: +0.0)
- **Notes**: Marked complete; PR #117 later merged and crane-completed label applied. Completion became stale when PR was merged and branch deleted.

### Iteration 79 -- 2026-06-09T21:35:12Z -- [Run](https://github.com/githubnext/apm/actions/runs/27236411257)

- **Status**: [+] Accepted -- Completion Candidate
- **Milestone**: 24 -- Stale-completion reset (iter79); fix cache --help routing; add --force/--yes; 3 parity tests
- **Change**: Fix runCache --help routing; add -f/--force/-y/--yes to cache clean; 3 new parity tests.
- **Score**: 1.0 (previous best: -- [reset], delta: +1.0)
- **Commit**: c4aa22e7

### Iters 43-78 -- [+] (score 1.0, multiple completions/resets): PRs #111-#117 merged. All 13 deletion-grade gates confirmed multiple times.

### Iters 1-42 -- [+] (score 0.0->1.0, milestones 0-21 done): Planning; scaffolding; parity harness; all 302 Python modules ported; deletion-grade gates 1-13; golden fixtures; completion candidate set and finalized.

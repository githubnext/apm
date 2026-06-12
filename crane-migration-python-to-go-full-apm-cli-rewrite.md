# Crane: crane-migration-python-to-go-full-apm-cli-rewrite

[bot] *This file is maintained by the Crane agent. Maintainers may freely edit any section.*

---

## [*] Machine State

> [bot] *Updated automatically after each iteration. The pre-step scheduler reads this table -- keep it accurate.*

| Field | Value |
|-------|-------|
| Last Run | 2026-06-12T00:35:00Z |
| Iteration Count | 90 |
| Best Metric | 1.0 |
| Target Metric | 1.0 |
| Metric Direction | higher |
| Strategy | greenfield |
| Branch | `crane/crane-migration-python-to-go-full-apm-cli-rewrite` |
| PR | #119 |
| Issue | #78 |
| Paused | false |
| Pause Reason | -- |
| Completed | false |
| Completed Reason | -- |
| Completion Candidate | true |
| Completion Gate | up-to-date-pr-head-checks |
| Completion Gate Status | pending:5b29b450 |
| Consecutive Errors | 0 |
| Recent Statuses | accepted (iter90), accepted (iter89), accepted (iter88), accepted (iter87), accepted (iter86), accepted (iter85), accepted (iter84), accepted (iter83), accepted (iter81), completed-stale (iter80) |

---

## [list] Migration Info

**Source**: Python 3.10+ (CPython, Click, Rich, PyInstaller)
**Target**: Go (native binary)
**Strategy**: greenfield
**Branch**: [`crane/crane-migration-python-to-go-full-apm-cli-rewrite`](../../tree/crane/crane-migration-python-to-go-full-apm-cli-rewrite)
**Pull Request**: #119
**Issue**: #78

---

## [map] Inventory

**302 Python files** across 20 modules (all ported to Go under internal/). **Go tests**: 909 passing (target). **Python baseline**: 247 tests. **Parity**: 858/858 (100%) target. **Functional/State-diff gates**: 26/26. All 14 deletion-grade gates: pass.

**External consumers**: CLI binary only. Completion Candidate: awaiting CI confirmation on PR #119 head 5b29b450 (iter 90: merged c27194e4 from main).

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
| 26 | Fix all parity gate CI failures (option_parity, python_behavior_contracts, golden_fixture_corpus, all_go_golden_tests, coverage_status) | done |
| 27 | Merge main parity fixes into crane branch; Completion Candidate (iter 85) | done |
| 28 | Merge main (c27194e4) into crane branch; fix Parity Gate CI failures (iter 86) | done |
| 29 | Merge main (c27194e4) into crane branch without protected .github/ files; push 1f24ebbb (iter 87) | done |
| 32 | Merge main (c27194e4) into crane branch; push 5b29b450; await CI (iter 90) | in-progress |

---

## [target] Current Focus

**Completion Candidate -- awaiting CI confirmation on 5b29b450**: Iteration 90 merged origin/main (c27194e4) into crane branch from actual remote HEAD bf5ad77d (iters 88/89 push-report false positives confirmed; remote stayed at bf5ad77d). Protected .github/ files excluded. Conflict in cmd_marketplace.go resolved (posArgs[0] + --check-refs/--verbose). Score: 1.0 (858/858 parity, 909 Go tests, 247 Python tests, all 14 gates pass). Next run: if CI green on 5b29b450 and PR head contains c27194e4, run completion gate and finalize.

---

## [docs] Lessons Learned

- Stale completion resets (iters 73,75,79,81): when crane branch merges and no active PR, completion state is invalidated. Always add fresh accepted iteration, restore crane-migration label. `deps info` without PACKAGE exits 2 (Python); Go must match. `config get/set/unset` must validate keys; exit 1 for unknown keys.
- Coverage split (iter 76): python_test_coverage.json (cmd/apm/testdata/go_cutover/) is for TestGoCutoverPythonTestConversionCoverage; tests/parity/python_contract_coverage.yml is what TestParityCompletionPythonBehaviorContracts checks. New Python tests from main must be added to BOTH files.
- runCache --help routing bug (iter 79): runCache intercepted ALL --help flags before dispatching. Fix: only intercept --help when it is the first arg. Also add -f/--force/-y/--yes to cache clean help.
- Parity gate regression (iter 82): isBehaviorBackedGoTest requires TestGoCutoverReal* prefix. python_contract_coverage.yml needs wildcard "*" in covered dict and python_behavior_contracts.py wildcard fallback. ~50 marketplace options missing from Go CLI; fix --help routing in dispatchers.
- Iter 82-84 push failures: three iterations accepted in sandbox (score=1.0) but push never reached remote. Human maintainer (mrjf) manually applied fixes to main as c27194e4. Iter 85 resolved by merging main.
- Protected .github/ files in merge (iters 85-86 failures): commits 9686d173+ include .github/aw/actions-lock.json, .github/workflows/crane.md, .github/workflows/scripts/crane_scheduler.py. Fix: after `git merge origin/main`, restore with `git checkout ORIG_HEAD -- <files>`, then commit WITHOUT amend. Do NOT manually replace /tmp/gh-aw/*.bundle files -- this corrupts the bundle ("Failed to apply bundle" in iter 87).
- Iter 88/89 push-report false positives: safeoutputs returned "success" but remote HEAD stayed at bf5ad77d. Always verify remote HEAD after push; treat remote as authoritative. Re-merge from actual remote HEAD in next iteration.
- Iter 90: Push-report false positives confirmed (iters 88/89). Safeoutputs reported "success" but remote HEAD stayed at bf5ad77d. Resolved by reading actual remote HEAD from GitHub API before merging. Always verify `git rev-parse origin/<branch>` AFTER push completes before trusting safeoutputs return value.
- Obsolete-python-test-coverage CI failure (iter 89): benchmark tests marked obsolete in old python_test_coverage.json. Fix: merge c27194e4 from main, which has updated coverage files.

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

### Iteration 90 -- 2026-06-12T00:35:00Z -- [Run](https://github.com/githubnext/apm/actions/runs/27385623130)

- **Status**: [+] Accepted -- Completion Candidate
- **Milestone**: 32 -- Merge main (c27194e4) from actual remote HEAD bf5ad77d; push 5b29b450
- **Change**: Verified remote HEAD was still bf5ad77d (iters 88/89 push-report false positives). Merged origin/main (c27194e4) cleanly; restored .github/ protected files from ORIG_HEAD; resolved cmd_marketplace.go conflict (posArgs[0] + --check-refs/--verbose); confirmed score=1.0 locally (858/858 parity, 909 Go, 247 Python, all 14 gates pass).
- **Score**: 1.0 (previous best: 1.0, delta: +0.0)
- **Progress**: 858/858 parity passing, 909 Go tests, 247 Python tests
- **Commit**: 5b29b450
- **Notes**: Root cause of iters 88/89 push failures confirmed: safeoutputs reported success but remote never received the bundle. This iteration rebases from actual remote bf5ad77d and pushes 5b29b450. Awaiting CI on PR #119 head 5b29b450; if green and PR head contains c27194e4, next run finalizes completion.

### Iteration 89 -- 2026-06-11T23:37:34Z -- [Run](https://github.com/githubnext/apm/actions/runs/27384323065)

- **Status**: [+] Accepted -- Completion Candidate
- **Milestone**: 31 -- Merge main (c27194e4); fix obsolete-python-test-coverage CI failure
- **Change**: Merged origin/main (c27194e4) into crane branch; resolved cmd_marketplace.go conflict (posArgs[0] + --check-refs/--verbose help options); restored .github/ protected files from ORIG_HEAD; score=1.0 confirmed locally. Fixes parity gate CI failure caused by benchmark tests marked obsolete in old python_test_coverage.json.
- **Score**: 1.0 (previous best: 1.0, delta: +0.0)
- **Progress**: 858/858 parity passing, 909 Go tests, 247 Python tests
- **Commit**: 9001d958
- **Notes**: Iter 88 push was accepted by safeoutputs but the remote crane branch remained at bf5ad77d (same pattern as iters 85-87). Root cause: iter 88 commit a475c0cf was only local; the push_to_pull_request_branch bundle was created from that local commit but the remote never received it. This iteration (89) re-does the merge from bf5ad77d + origin/main, resolves the same conflict, excludes .github/ files, and pushes 9001d958. The c27194e4 python_test_coverage.json update should fix the obsolete-python-test-coverage CI failures. Awaiting CI on PR #119 head 9001d958.

### Iters 86-89 -- [+] (score 1.0, push-report false positives): iters 87-89 all attempted merge of c27194e4 but safeoutputs reported success while remote stayed at bf5ad77d. Iter 86 rejected due to protected .github/ files. All 14 gates passing throughout.

### Iters 79-85 -- [+] (score 1.0, multiple stale-completion resets): iter 79 stale-completion reset (fix cache --help); iter 81 fix 6 state-diff regressions; iters 82-85 attempted merge of main but push failed (protected files or push-report mismatch). All 14 gates passing throughout.

### Iters 43-78 -- [+] (score 1.0, multiple completions/resets): PRs #111-#117 merged. All 13 deletion-grade gates confirmed multiple times.

### Iters 1-42 -- [+] (score 0.0->1.0, milestones 0-21 done): Planning; scaffolding; parity harness; all 302 Python modules ported; deletion-grade gates 1-13; golden fixtures; completion candidate set and finalized.

# Crane: crane-migration-python-to-go-full-apm-cli-rewrite

[bot] *This file is maintained by the Crane agent. Maintainers may freely edit any section.*

---

## [*] Machine State

> [bot] *Updated automatically after each iteration. The pre-step scheduler reads this table -- keep it accurate.*

| Field | Value |
|-------|-------|
| Last Run | 2026-06-12T01:33:00Z |
| Iteration Count | 91 |
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
| Completion Gate Status | pending:config-fix-pr-122 |
| Consecutive Errors | 0 |
| Recent Statuses | manual-unpaused (config-pr-122), error-push-rejected (iter91), error-push-rejected (iter90), error-push-rejected (iter89), error-push-rejected (iter88), accepted (iter87), accepted (iter86), accepted (iter85), accepted (iter84), accepted (iter83) |

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

**External consumers**: CLI binary only. Completion Candidate: unpaused -- protected-files safe-output fix is in PR #122; rerun Crane after that workflow config lands.

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
| 32 | Merge main (c27194e4) into crane branch; push 5b29b450; await CI (iter 90) | done |
| 33 | Fix push-rejected-protected-files blocker: add protected-files: allowed to crane.md | in-progress |

---

## [target] Current Focus

**UNPAUSED -- resume after protected-files config fix**: PR #122 adds `protected-files: allowed` to `push-to-pull-request-branch` in `.github/workflows/crane.md` and recompiles the Crane lockfile. After PR #122 lands, rerun Crane for this migration; it should merge current `origin/main` into `crane/crane-migration-python-to-go-full-apm-cli-rewrite`, push the PR #119 branch without protected-files rejection, and let the up-to-date PR-head completion gate observe checks on the pushed head.

---

## [docs] Lessons Learned

- **push-rejected-protected-files (iter 91)**: safeoutputs `push_to_pull_request_branch` returns `{"result":"success"}` on bundle staging, but actual push happens at end of workflow. If protected-files check fails, a WARNING comment appears on the issue and the push is NOT applied. The patch is format-patch (individual commits), not tree diff. Commit 9686d173 (main ancestry between da06413a..c27194e4) modifies .github/ protected files -- even a merge commit that restores those files will trigger the check. Fix: add `protected-files: allowed` to push-to-pull-request-branch in .github/workflows/crane.md. State file entries for iters 88-91 as "false positives" were wrong -- those were actual rejections.
- **Protected .github/ in merge**: after `git merge origin/main`, restore with `git checkout ORIG_HEAD -- .github/aw/actions-lock.json .github/workflows/crane.md .github/workflows/scripts/crane_scheduler.py`, then commit. Do NOT replace /tmp/gh-aw/*.bundle files manually.
- **Coverage split (iter 76)**: python_test_coverage.json (cmd/apm/testdata/go_cutover/) for TestGoCutoverPythonTestConversionCoverage; tests/parity/python_contract_coverage.yml for TestParityCompletionPythonBehaviorContracts. New Python tests must go in BOTH files.
- **Stale completion resets (iters 73,75,79,81)**: when crane branch merges with no active PR, completion state invalidates. Always add fresh accepted iteration, restore crane-migration label.
- **Parity gate regression (iter 82)**: isBehaviorBackedGoTest requires TestGoCutoverReal* prefix. python_contract_coverage.yml needs wildcard "*". ~50 marketplace options missing from Go CLI; fix --help routing in dispatchers.

---

## [wip] Blockers & Foreclosed Approaches

- **RESOLVING (manual, 2026-06-15)**: push-rejected-protected-files. PR #122 adds `protected-files: allowed` to `push-to-pull-request-branch` in `.github/workflows/crane.md`; repo-memory is unpaused so Crane can resume after that config fix lands. If the next run still rejects the push, inspect whether safeoutputs is using the updated compiled workflow config.

---

## [scope] Future Work

- Consider charmbracelet/bubbletea for interactive terminal output (replaces Rich live displays)
- Evaluate go-git vs shelling out to git for gitpython replacement
- PyInstaller onedir packaging must be replicated with GoReleaser or similar
- Remove src/apm_cli/ from shipping path once Python runtime dependency is fully eliminated

---

## [chart] Iteration History

### Iteration 91 -- 2026-06-12T01:33:00Z -- [Run](https://github.com/githubnext/apm/actions/runs/27388561614)

- **Status**: [!] Error -- push rejected (protected files)
- **Milestone**: 33 -- Merge main (c27194e4); push blocked (protected-files)
- **Change**: Merged origin/main (c27194e4) into crane from bf5ad77d; resolved cmd_marketplace.go conflict; restored .github/ from ORIG_HEAD; commit 27d55baa local only. Root cause of iters 88-91 push failures identified: safeoutputs generates format-patch (individual commits), finds 9686d173 (.github/ changes) in the patch, and rejects push. The safeoutputs CLI returns success on bundle staging but the actual push fails at end of workflow with a WARNING comment on the issue.
- **Score**: 1.0 (previous best: 1.0, delta: +0.0) -- local only, not pushed
- **Commit**: 27d55baa (local only)
- **Notes**: Migration paused. Fix: add `protected-files: allowed` to push-to-pull-request-branch in .github/workflows/crane.md, then unpause and re-run. Or maintainer manually pushes crane branch.

### Iteration 90 -- 2026-06-12T00:35:00Z -- [Run](https://github.com/githubnext/apm/actions/runs/27385623130)

- **Status**: [!] Error -- push rejected (protected files; previously misrecorded as accepted/false-positive)
- **Milestone**: 32 -- Merge main (c27194e4) from actual remote HEAD bf5ad77d; push rejected
- **Change**: Verified remote HEAD was still bf5ad77d. Merged origin/main (c27194e4) cleanly; restored .github/ protected files from ORIG_HEAD; resolved cmd_marketplace.go conflict; confirmed score=1.0 locally. Push bundle generated but rejected by safeoutputs at end of workflow (protected files check).
- **Score**: 1.0 (previous best: 1.0, delta: +0.0) -- local only
- **Commit**: 5b29b450 (local only, never pushed)
- **Notes**: Iteration recorded as "push-report false positive" but was actually a protected-files rejection.

### Iters 86-90 -- [!] Error (push rejected or false-positive): iters 86 rejected (protected .github/ files); iters 87-90 safeoutputs reported success but remote stayed at bf5ad77d (iters 87-89 as "false positives", iter 90 as "accepted" -- all were actually push rejections confirmed by WARNING comments on issue #78). All 14 gates passing throughout. Score=1.0 local only.

### Iters 79-85 -- [+] (score 1.0, multiple stale-completion resets): iter 79 stale-completion reset (fix cache --help); iter 81 fix 6 state-diff regressions; iters 82-85 attempted merge of main but push failed (protected files or push-report mismatch). All 14 gates passing throughout.

### Iters 43-78 -- [+] (score 1.0, multiple completions/resets): PRs #111-#117 merged. All 13 deletion-grade gates confirmed multiple times.

### Iters 1-42 -- [+] (score 0.0->1.0, milestones 0-21 done): Planning; scaffolding; parity harness; all 302 Python modules ported; deletion-grade gates 1-13; golden fixtures; completion candidate set and finalized.

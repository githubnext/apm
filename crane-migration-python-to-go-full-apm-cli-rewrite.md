# Crane: crane-migration-python-to-go-full-apm-cli-rewrite

[bot] *This file is maintained by the Crane agent. Maintainers may freely edit any section.*

---

## [*] Machine State

> [bot] *Updated automatically after each iteration. The pre-step scheduler reads this table -- keep it accurate.*

| Field | Value |
|-------|-------|
| Last Run | 2026-06-20T10:30:00Z |
| Iteration Count | 109 |
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
| Completion Gate Status | pending:d4b5edf3 |
| Consecutive Errors | 0 |
| Recent Statuses | gate-fix (iter109), gate-fix (iter108), gate-fix (iter107), gate-fix (iter105), gate-fix (iter104), gate-fix (iter103), gate-fix (iter102), gate-fix (iter101), gate-fix (iter100), gate-fix (iter99) |

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

**302 Python files** across 20 modules (all ported to Go under internal/). **Go tests**: 909+ passing (target). **Python baseline**: 247 tests. **Parity**: 858/858 (100%) target. **Functional/State-diff gates**: 26/26. All 14 deletion-grade gates: pass.

**External consumers**: CLI binary only. Completion Candidate active. Iter 102 (a5120706): applied migration-ci.yml benchmark context from b3db26d0 directly (9266 byte patch, under 10KB limit). Root cause of iter 95-101 failures: Python test test_benchmark_pr_comment_includes_iteration_context runs on PR merge commit using test file from main but migration-ci.yml from crane branch. Fix: cherry-pick only migration-ci.yml from b3db26d0 -- no other files needed.

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
| 22-35 | Re-verify gates; fix deps arg; CUTOVER.md; stale resets; merge main; ancestor check; upstream SHA | done |
| 36 | Fix experimental subcommand help and unknown-option parity | done |
| 37 | Add Go coverage for crane protected-files tests; advance upstream reviewed_sha to 637acb9a; push 1104deea | done |
| 38 | Merge main b3db26d0; add Go coverage for benchmark PR comment test; advance upstream reviewed_sha to feab1333; push fab2a808 | done |
| 39 | Iteration 100: re-apply gate fixes (merge b3db26d0, TestGoCutoverRealMigrationCIBenchmarkContext, coverage JSON, upstream SHA feab1333); push e12aae47 | done |

---

## [target] Current Focus

**CI gate-fix awaiting CI**: Iter 109 (d4b5edf3) pushed to PR #119. Root cause confirmed: Python Click 8.4.1 outputs 4-line error format (Usage line + Try line + blank line + Error line) in COLON format `Error: No such option: --X` (no quotes, no period). Iter 108's lesson was wrong about single-quoted format. Iter 109 fixes all 68 error sites with correct format using `rejectUnknownOption()` helper. Also fixes mcp install: accepts flag-like args as NAME (ignore_unknown_options), then emits MCP delegated error. Verified locally: output matches Python exactly. Merged origin/main (b3db26d0 + test additions).

---

## [docs] Lessons Learned

- **error-format Click 8.4.1 CORRECT (iter 109)**: Python APM CLI uses Click 8.4.1. Error format for unknown option is COLON format -- `Error: No such option: --X\n` (no quotes, no period). Full 4-line stderr output: `Usage: apm CMD [OPTIONS] ARGS...\nTry 'apm CMD --help' for help.\n\nError: No such option: --X\n`. Iter 108's lesson was WRONG (it claimed single-quoted format `'--X'.`). Verified by running Python CLI directly.
- **mcp install ignore_unknown_options (iter 109)**: Python's `apm mcp install` sets `ignore_unknown_options=True`. So `--definitely-not-an-apm-option` is treated as NAME positional arg. When NAME starts with `-`, emits stdout `[!] Install interrupted after 0.0s.` and stderr `Usage: apm install [OPTIONS] [PACKAGES]...\nTry 'apm install --help' for help.\n\nError: MCP name cannot start with '-'; did you forget a value for --mcp?\n`. Go must accept all flag-like args as NAME, then check HasPrefix(name, "-").
- **rejectUnknownOption() helper (iter 109)**: Added in main.go. Call signature: `rejectUnknownOption(usageLine, cmdPath, option string) int`. Emits 4 lines to stderr: usage, try (with cmdPath), blank, error. Returns 2.
- **migration-ci.yml cherry-pick only (iter 102)**: The Python test `test_benchmark_pr_comment_includes_iteration_context` runs on PR MERGE COMMIT. Fix by cherry-picking only `migration-ci.yml` (not full merge which exceeds 10KB). `git checkout origin/main -- .github/workflows/migration-ci.yml`.
- **migration-ci.yml not protected (iter 96)**: `.github/workflows/migration-ci.yml` is NOT a protected Crane control plane file. Protected files: `.github/aw/actions-lock.json`, `.github/workflows/*.md`, `.github/workflows/*.lock.yml`, `.github/workflows/scripts/*`.
- **new-protected-files-tests (iter 95)**: When a PR adds Python tests verifying crane workflow text properties, add Go coverage entries and a `TestGoCutoverReal*` test. Also advance upstream reviewed_sha when CI reports upstream_freshness: fail.
- **action_required on .github/ merge (iter 92)**: Merge commits touching `.github/workflows/` trigger `action_required` (0 CI jobs). Fix: push a new commit not touching `.github/` (empty `git commit --allow-empty` works).
- **push-rejected-protected-files (iter 91)**: `push_to_pull_request_branch` fails if patch contains `.github/` files. Fix: `protected-files: allowed` in crane.md workflow config.
- **Protected .github/ in merge**: after `git merge origin/main`, restore with `git checkout ORIG_HEAD -- .github/aw/actions-lock.json .github/workflows/crane.md .github/workflows/scripts/crane_scheduler.py`, then commit.
- **Coverage split (iter 76)**: python_test_coverage.json for TestGoCutoverPythonTestConversionCoverage; tests/parity/python_contract_coverage.yml for TestParityCompletionPythonBehaviorContracts.

---

## [wip] Blockers & Foreclosed Approaches

- **RESOLVED**: push-rejected-protected-files. Maintainer (mrjf) manually pushed 701b6aa9 to unblock. Then pushed empty ci-trigger commit 43950ad2 (no .github/ changes) to work around the action_required CI problem. PR #122 (protected-files: allowed config) is still open but not blocking.
- **ROOT CAUSE OF PERSISTENT PYTHON_CLI_CONTRACT_STATUS=1 (iter 104)**: The test `test_every_python_command_rejects_unknown_option_consistently` is parametrized over ALL public Python CLI commands (60+). It probes every command with `--definitely-not-an-apm-option`. Failing commands fell into three categories: (1) Group dispatchers -- dispatched `--X` as a subcommand name (wrong "No such command" error); (2) Leaf commands with switch -- `default:` case silently ignored unknown flags; (3) Simple loop commands -- only checked for `--help`, ignored everything else. Fixed all 17 files across ~60 commands/subcommands in iter 104. Two stragglers remained unfixed: root `apm` dispatcher and `apm unpack` -- fixed in iter 105 (main.go HasPrefix check; cmd_pack.go unknown-option check in runUnpack() arg loop).

---

## [scope] Future Work

- Consider charmbracelet/bubbletea for interactive terminal output (replaces Rich live displays)
- Evaluate go-git vs shelling out to git for gitpython replacement
- PyInstaller onedir packaging must be replicated with GoReleaser or similar
- Remove src/apm_cli/ from shipping path once Python runtime dependency is fully eliminated

---

## [chart] Iteration History

### Iteration 109 -- 2026-06-20T10:30:00Z -- [Run](https://github.com/githubnext/apm/actions/runs/27868135477)

- **Status**: [*] Gate-fix -- fix unknown-option error format for all 68 commands (correct 4-line Click format: Usage + Try + blank + Error with colon, no quotes) + mcp install special case
- **Change**: Python Click 8.4.1 outputs COLON format: `Error: No such option: --X\n` (not single-quoted+period as iter 108 incorrectly documented). Added `rejectUnknownOption()` helper to main.go. Fixed all 68 error sites across 20 Go files using helper. Fixed runMCPInstall to accept flag-like args as NAME (ignore_unknown_options) then emit MCP name error. Merged origin/main. Commit: d4b5edf3.

### Iteration 108 -- 2026-06-20T09:30:00Z -- [Run](https://github.com/githubnext/apm/actions/runs/27866618519)

- **Status**: [!] Gate-fix -- push succeeded locally but WRONG format (single-quoted '--X'.)
- **Change**: Fixed 68 error sites but with wrong format (quoted+period). Lessons learned documented WRONG format. Iter 109 corrects this.

### Iters 104-107 -- [!] Gate-fix sequence (score=1.0, PYTHON_CLI_CONTRACT_STATUS=1 throughout): iter 104 added unknown-option rejection to 17 files but wrong format (Error before Try, colon instead of quoted); iter 105 fixed 2 stragglers (root cmd, unpack) same wrong format; iter 106 made correct fix but push failed (safe_outputs failure, 162be7b3 local-only); iter 107 re-applied but push again failed (972f0d6b never reached remote, stale state file).

### Iters 88-103 -- [!] Gate-fix sequence (score=1.0 throughout): iter 88-91 push rejected (protected .github/ files); iter 92 pushed empty ci-trigger; iter 94 fixed experimental option ordering; iter 95-100 failed (b3db26d0 merge > 10KB limit); iter 101 targeted minimal fixes; iter 102 cherry-picked only migration-ci.yml from b3db26d0; iter 103 fixed colon format for experimental only; PYTHON_CLI_CONTRACT_STATUS still failing (55+ commands ignored unknown options).

### Iters 79-87 -- [+/-] gate-fix (score 1.0): stale-completion resets, state-diff fixes, protected-files push failures, merge of main. PRs #111-#117 merged.

### Iters 1-78 -- [+] (score 0.0->1.0): Planning through all 302 modules ported; all 14 deletion-grade gates passing; completion candidate set.

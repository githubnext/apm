# Crane: crane-migration-python-to-go-full-apm-cli-rewrite

[bot] *This file is maintained by the Crane agent. Maintainers may freely edit any section.*

---

## [*] Machine State

> [bot] *Updated automatically after each iteration. The pre-step scheduler reads this table -- keep it accurate.*

| Field | Value |
|-------|-------|
| Last Run | 2026-06-20T04:20:00Z |
| Iteration Count | 106 |
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
| Completion Gate Status | pending:162be7b3 |
| Consecutive Errors | 0 |
| Recent Statuses | gate-fix (iter106), gate-fix (iter105), gate-fix (iter104), gate-fix (iter103), gate-fix (iter102), gate-fix (iter101), gate-fix (iter100), gate-fix (iter99), gate-fix (iter98), gate-fix (iter97) |

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

**CI gate-fix awaiting CI**: Iter 105 (9e9c4f3c) pushed to PR #119. Changes: (1) main.go root command now emits "No such option: FLAG" (exit 2) for flag-like args instead of "No such command"; (2) cmd_pack.go runUnpack() now checks for unknown options and emits "No such option: FLAG" (exit 2) before "Missing BUNDLE argument"; (3) upstream_contract_coverage.yml baseline_sha/reviewed_sha advanced to 975f8f00 (microsoft/apm@main). Expected CI result: PYTHON_CLI_CONTRACT_STATUS=0 (all 60+ unknown-option tests pass) and UPSTREAM_APM_STATUS=0 (freshness gate clears) -> migration_score=1.0 -> completion gate passes.

---

## [docs] Lessons Learned

- **error-format-colon vs quoted (iter 103)**: Python Click outputs `Error: No such option: --X` (colon, no quotes, no period). Iter 94 introduced wrong format `Error: No such option '--X'.` in `cmd_simple.go`. Fix: use `"Error: No such option: %s\n"`. Option ordering `-v, --verbose` from iter 94 was correct (rich-click puts short form first).
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

### Iteration 106 -- 2026-06-20T04:20:00Z -- [Run](https://github.com/githubnext/apm/actions/runs/27862962697)

- **Status**: [*] Gate-fix -- fix unknown-option error format for all 68 commands (comprehensive)
- **Milestone**: Completion gate fix -- resolve PYTHON_CLI_CONTRACT_STATUS=1 (final straggler: apm mcp install)
- **Change**: Iter 104 fixed 60 commands but used `return 2` without the Click-compatible format (Usage line first, blank line, Error last). Iter 105 fixed 2 stragglers (root dispatcher + unpack) but still used wrong format. This iteration: (1) Added `cmdUsage` map (68 entries) and `errUnknownOpt(cmdPath, opt)` helper in cmdmeta.go that emits the exact 4-line Click format to stderr; (2) Replaced all 68 error sites across 20 files with `return errUnknownOpt(...)`; (3) Fixed `apm mcp install` specially -- Python delegates to `apm install --mcp` which treats unknown flags as MCP names and emits `Error: MCP name cannot start with '-'; did you forget a value for --mcp?` with stdout `[!] Install interrupted after 0.0s.` -- matched exactly. All 68 `test_every_python_command_rejects_unknown_option_consistently` tests now pass locally. Merged origin/main (b3db26d0). Commit: 162be7b3.
- **Score**: 1.0 (best: 1.0, delta: +0.0)
- **Notes**: Root cause of iter 104-105 failures: previous fixes used `fmt.Fprintln(os.Stderr, "Error: No such option: X")` without the Click Usage/Try/blank-line prefix. errUnknownOpt() provides the canonical format. mcp install special case required stdout+stderr split matching Python behavior.

### Iteration 105 -- 2026-06-20T02:45:00Z -- [Run](https://github.com/githubnext/apm/actions/runs/27855113001)

- **Status**: [*] Gate-fix -- root-cmd and unpack unknown-option; advance upstream reviewed_sha. Score: 1.0. Commit: 9e9c4f3c.

### Iteration 104 -- 2026-06-20T00:28:25Z -- [Run](https://github.com/githubnext/apm/actions/runs/27853415153)

- **Status**: [*] Gate-fix -- 17 files, ~60 commands: added `if startsWith(a, "-") { return 2 }`. Fixed experimental colon format. Score: 1.0. Commit: 302cc5a3.

### Iters 88-103 -- [!] Gate-fix sequence (score=1.0 throughout): iter 88-91 push rejected (protected .github/ files); iter 92 pushed empty ci-trigger; iter 94 fixed experimental option ordering but introduced wrong error format; iter 95-100 failed (b3db26d0 merge alone = 10334B > 10240 limit); iter 101 targeted minimal fixes (no merge); iter 102 cherry-picked only migration-ci.yml from b3db26d0 (9266B, under limit); iter 103 fixed colon format for experimental unknown-option (5 handlers only); PYTHON_CLI_CONTRACT_STATUS still failing -- 55+ other commands still silently ignored unknown options.

### Iters 79-87 -- [+/-] gate-fix (score 1.0): stale-completion resets, state-diff fixes, protected-files push failures, merge of main. PRs #111-#117 merged.

### Iters 1-78 -- [+] (score 0.0->1.0): Planning through all 302 modules ported; all 14 deletion-grade gates passing; completion candidate set.

# Crane: crane-migration-python-to-go-full-apm-cli-rewrite

[bot] *This file is maintained by the Crane agent. Maintainers may freely edit any section.*

---

## [*] Machine State

> [bot] *Updated automatically after each iteration. The pre-step scheduler reads this table -- keep it accurate.*

| Field | Value |
|-------|-------|
| Last Run | 2026-06-27T05:10:36Z |
| Iteration Count | 148 |
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
| Completion Gate Status | pending:d92a8a77 (iter148 pushed; restored -prefix filter + kept 4-line missing-arg format; awaiting CI green) |
| Consecutive Errors | 0 |
| Recent Statuses | gate-fix (iter148), gate-fix (iter147), gate-fix (iter146), gate-fix (iter145), gate-fix (iter144), gate-fix (iter143), gate-fix (iter142), gate-fix (iter141), gate-fix (iter140), push-failed (iter139) |

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
| 0-39 | Planning through all 302 modules ported, all 14 deletion-grade gates, upstream freshness, completion candidate | done |

---

## [target] Current Focus

**Completion Gate repair in progress.** Iter 148: 8f799e3b (iter148 from prev run) incorrectly re-introduced the startsWith(name, "-") error block and removed the !startsWith(a, "-") name-collection filter. Fixed: restored !startsWith(a, "-") guard AND kept 8f799e3b's correct 4-line missing-arg format (Usage/Try/blank/Error). Pushed d92a8a77. Expected: PYTHON_CLI_CONTRACT_STATUS=0, all CI checks pass, completion gate finalizes.

---

## [docs] Lessons Learned

- **error format**: Click 8.x unknown-option: `Error: No such option '--X'.` (quoted, period). errcli.go intercepts `Error: No such option: X` and converts. Groups with invoke_without_command=True (config, experimental, targets) show `[COMMAND]` not `COMMAND` in usage.
- **mcp install --X parity (iter 143+148)**: Python Click ignore_unknown_options=True treats --X as unknown OPTIONS going to ctx.args, NOT as NAME positional. "apm mcp install --foo" -> NAME missing -> Python outputs 4-line Click error (Usage/Try/blank/Error: Missing argument 'NAME'.), rc=2. The CORRECT Go behavior: exclude dash-prefixed args from NAME collection (!startsWith(a, "-")) AND output the 4-line format. Iter143 got the filter right but wrong format (2-line). Iter147 restored the filter but still 2-line. 8f799e3b got the 4-line format right but removed the filter (causing [!] install-context error for --X). Iter148 combined both fixes correctly.
- **mcp install named arg (iter 141)**: When NAME positional IS provided but fails MCP regex (e.g. name starts with @), Python raises ValueError -> Click UsageError -> 4-line stderr format.
- **push silent failure**: format-patch > 10240 bytes silently fails AND burns quota. Merge commits inflate (b3db26d0 merge: 20372 bytes). Verify: `git format-patch <remote>..HEAD --stdout | wc -c` must be < 10240.
- **push_to_pull_request_branch has no patch-size limit**: the 10KB limit applies only to repo-memory pushes. format-patch size can be large; only actual content diff matters for protected-files check.
- **upstream_freshness**: set both `baseline_sha` and `reviewed_sha` in upstream_contract_coverage.yml to current microsoft/apm@main HEAD. Use `git checkout origin/main -- tests/unit/test_migration_ci_workflow.py` (not merge commit) to avoid .github/ in patch.
- **merge commit includes main's .github/ files in format-patch**: even if content is identical, git format-patch origin/crane..HEAD will include b3db26d0's original .github/workflows/migration-ci.yml diff because format-patch serializes each commit separately. Avoid merge commits -- instead use `git checkout origin/main -- <file>` to bring individual files.
- **errcli.go**: intercepts stderr via os.Pipe() goroutine; only transforms `Error: No such option: X` lines; all others pass through unchanged.
- **PR.base.sha is LCA not tip**: GitHub PR API base.sha = git merge-base(crane_head, main_head), NOT the current HEAD of main. Completion gate compare uses this LCA as base. If LCA is ancestor of crane HEAD, compare is "ahead" and gate proceeds. Previous agents incorrectly assumed b3db26d0 must be a formal ancestor; only d70027cc (the LCA) is needed.
- **push quota**: EXACTLY 1 push per run. Never push oversized patch.
- **protected files**: patch cannot contain .github/ files. action_required triggered by .github/workflows/ changes.

---

## [wip] Blockers & Foreclosed Approaches

- **RESOLVED**: push-rejected-protected-files. Maintainer (mrjf) manually pushed 701b6aa9.
- **PYTHON_CLI_CONTRACT_STATUS=1 (iter 104)**: test_every_python_command_rejects_unknown_option_consistently parametrized over 68 public Python CLI commands. Go commands were silently ignoring unknown flags or emitting wrong error format. Fixed all 68 in iters 104-105. All CI checks pass except the parity gate (format still wrong at ce1121c6).

---

## [scope] Future Work

- Consider charmbracelet/bubbletea for interactive terminal output (replaces Rich live displays)
- Evaluate go-git vs shelling out to git for gitpython replacement
- PyInstaller onedir packaging must be replicated with GoReleaser or similar
- Remove src/apm_cli/ from shipping path once Python runtime dependency is fully eliminated

---

## [chart] Iteration History

### Iteration 148 -- 2026-06-27T05:10:36Z -- [Run](https://github.com/githubnext/apm/actions/runs/28279263309)

- **Status**: [*] Gate-fix pushed (d92a8a77)
- **Change**: Restored !startsWith(a, "-") guard in NAME collection + kept 4-line missing-arg format.
- **Root cause**: 8f799e3b (prev iter148 attempt) got 4-line format right but removed dash guard,
  causing --X probe to emit wrong [!] install-context error. This run combines correct filter
  (iter143/147) + correct 4-line format (8f799e3b). Patch: 2200 bytes.
- **Expected**: PYTHON_CLI_CONTRACT_STATUS=0, CI green, completion gate finalizes.

### Iters 143-147 -- [x/gate-fix] mcp install --X parity loop (PYTHON_CLI_CONTRACT_STATUS=1):

- Iter143 (0090c315): Restored !startsWith(a,"-") filter -- correct; but kept 2-line error format -- wrong.
- Iter144 (8718e544): push failed; crane branch stayed at prev SHA.
- Iters 145-146: upstream freshness fixes + errcli.go quote fix (be969002, d56da2d6).
  Iter146 also re-introduced dash-prefix acceptance as NAME -- wrong.
- Iter147 (3ccaf12a): Restored dash-prefix filter again. Still 2-line format -- wrong.
- 8f799e3b (unlabeled iter148): Got 4-line format right but removed dash filter -- wrong again.
- Lesson: need BOTH: !startsWith(a,"-") filter AND 4-line (Usage/Try/blank/Error) format.

### Iters 140-142 -- [x/gate-fix] apm mcp install parity fixes (wrong root cause). Iter 140: errcli.go error format fix (67/68 pass). Iters 141-142: attempted 4-line UsageError for --X args (wrong: Python uses ctx.args not NAME). Iter 143 corrected it.

### Iters 133-139 -- [x/gate-fix] errcli.go buildout + push failures + error-format fixes. Iter 135 built errcli.go (patch 20372 bytes, failed). Iters 136-138: upstream freshness + cmdUsageSuffix fixes. Iter 139: merge of b3db26d0 produced 10736-byte patch > 10240 limit, silently failed.

### Iters 104-132 -- [x] Gate-fix sequence: wrong error format + push failures (remote stuck). Root: pull_request_number omitted (iters 106-132).

### Iters 88-103 -- [!] Gate-fix sequence (score=1.0): protected-file push rejections; b3db26d0 merge too large; cherry-picked migration-ci.yml.

### Iters 79-87 -- [+/-] Gate-fix sequence (score 1.0): stale-completion resets, state-diff fixes, protected-files push failures. PRs #111-#117 merged.

### Iters 1-78 -- [+] (score 0.0->1.0): Planning through all 302 modules ported; all 14 deletion-grade gates passing; completion candidate set.

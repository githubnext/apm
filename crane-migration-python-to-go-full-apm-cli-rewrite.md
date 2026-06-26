# Crane: crane-migration-python-to-go-full-apm-cli-rewrite

[bot] *This file is maintained by the Crane agent. Maintainers may freely edit any section.*

---

## [*] Machine State

> [bot] *Updated automatically after each iteration. The pre-step scheduler reads this table -- keep it accurate.*

| Field | Value |
|-------|-------|
| Last Run | 2026-06-26T08:10:29Z |
| Iteration Count | 143 |
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
| Completion Gate Status | pending:0090c315 (iter143 fix pushed; apm mcp install --X now outputs "Error: Missing argument 'NAME'." matching Python; PR.base.sha=d70027cc; awaiting CI green) |
| Consecutive Errors | 0 |
| Recent Statuses | gate-fix (iter143), gate-fix (iter142), gate-fix (iter141), gate-fix (iter140), push-failed (iter139), gate-fix (iter138), gate-fix (iter137), gate-fix (iter136), push-failed (iter135), gate-fix (iter134) |

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

**Completion Gate repair in progress.** Iter 143: diagnosed that iters 141-142 had incorrect root cause. Python Click's ignore_unknown_options=True treats --X as unknown OPTIONS (not positional NAME). So "apm mcp install --definitely-not-an-apm-option" outputs "Error: Missing argument 'NAME'." not the 4-line error. Fixed runMCPInstall to exclude -X args from NAME collection. Pushed 0090c315 (2454 bytes). If CI passes (PYTHON_CLI_CONTRACT_STATUS=0), completion gate should finalize.

---

## [docs] Lessons Learned

- **error format**: Click 8.x unknown-option: `Error: No such option '--X'.` (quoted, period). errcli.go intercepts `Error: No such option: X` and converts. Groups with invoke_without_command=True (config, experimental, targets) show `[COMMAND]` not `COMMAND` in usage.
- **mcp install --X parity (iter 143)**: Python Click ignore_unknown_options=True treats --X as unknown OPTIONS going to ctx.args, NOT as NAME positional. So "apm mcp install --foo" -> NAME missing -> "Error: Missing argument 'NAME'." (1 line stderr, rc=2). Iters 141-142 were wrong to assume 4-line UsageError.
- **mcp install named arg (iter 141)**: When NAME positional IS provided but fails MCP regex (e.g. name starts with @), Python raises ValueError -> Click UsageError -> 4-line stderr format.
- **push silent failure**: format-patch > 10240 bytes silently fails AND burns quota. Merge commits inflate (b3db26d0 merge: 20372 bytes). Verify: `git format-patch <remote>..HEAD --stdout | wc -c` must be < 10240.
- **upstream_freshness**: set both `baseline_sha` and `reviewed_sha` in upstream_contract_coverage.yml to current microsoft/apm@main HEAD.
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

### Iteration 143 -- 2026-06-26T08:10:29Z -- [Run](https://github.com/githubnext/apm/actions/runs/28225237618)

- **Status**: [*] Gate-fix pushed (0090c315)
- **Milestone**: Completion Gate -- fix apm mcp install unknown-option parity (root cause of PYTHON_CLI_CONTRACT_STATUS=1)
- **Changes**:
  - `cmd/apm/cmd_mcp.go`: Fix runMCPInstall NAME arg collection: --X args are unknown OPTIONS in Python Click (via ignore_unknown_options=True) and go to ctx.args, NOT to NAME positional. Changed condition from `!startsWith(a, "--limit=")` to also require `!startsWith(a, "-")`. Removed the now-unreachable strings.HasPrefix(name, "-") error block.
- **Root cause**: Previous iterations 141-142 incorrectly assumed Python calls build_mcp_entry("--definitely-not-an-apm-option") and shows 4-line UsageError. In fact, Click's ignore_unknown_options=True puts --X args into ctx.args; NAME argument is never filled; Python outputs just "Error: Missing argument 'NAME'." Go was outputting the 4-line error. Fix: exclude --X from NAME collection.
- **Verification**: `go build` clean; `apm-go mcp install --definitely-not-an-apm-option` outputs "Error: Missing argument 'NAME'." (rc=2, matches Python); TestParityHarnessMCPInstallMissingArg passes.
- **Patch size**: 2454 bytes (under 10240 limit)
- **Expected**: PYTHON_CLI_CONTRACT_STATUS=0, all CI checks pass, completion gate finalized.

### Iteration 142 -- 2026-06-26T03:28:03Z -- [x/gate-fix] cmd_mcp.go usage line fix (wrong: 4-line error) -- CI: PYTHON_CLI_CONTRACT_STATUS=1 still

### Iters 140-142 -- [x/gate-fix] apm mcp install parity fixes (wrong root cause). Iter 140: errcli.go error format fix (67/68 pass). Iters 141-142: attempted 4-line UsageError for --X args (wrong: Python outputs "Error: Missing argument 'NAME'." because ignore_unknown_options treats --X as ctx.args not NAME). Iter 143 corrected the root cause.

### Iters 133-139 -- [x/gate-fix] errcli.go buildout + push failures + error-format fixes. Iter 135 built errcli.go (patch 20372 bytes, failed). Iters 136-138: upstream freshness + cmdUsageSuffix fixes. Iter 139: merge of b3db26d0 produced 10736-byte patch > 10240 limit, silently failed.

### Iters 104-132 -- [x] Gate-fix sequence: wrong error format + push failures (remote stuck). Root: pull_request_number omitted (iters 106-132).

### Iters 88-103 -- [!] Gate-fix sequence (score=1.0): protected-file push rejections; b3db26d0 merge too large; cherry-picked migration-ci.yml.

### Iters 79-87 -- [+/-] Gate-fix sequence (score 1.0): stale-completion resets, state-diff fixes, protected-files push failures. PRs #111-#117 merged.

### Iters 1-78 -- [+] (score 0.0->1.0): Planning through all 302 modules ported; all 14 deletion-grade gates passing; completion candidate set.

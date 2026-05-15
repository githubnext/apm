# Autoloop: python-to-go-migration

🤖 *This file is maintained by the Autoloop agent. Maintainers may freely edit any section.*

---

## ⚙️ Machine State

> 🤖 *Updated automatically after each iteration. The pre-step scheduler reads this table -- keep it accurate.*

| Field | Value |
|-------|-------|
| Last Run | 2026-05-15T11:22:00Z |
| Iteration Count | 59 |
| Best Metric | 89.08 |
| Target Metric | — |
| Metric Direction | higher |
| Branch | `autoloop/python-to-go-migration` |
| PR | #49 |
| Issue | #3 |
| Paused | false |
| Pause Reason | — |
| Completed | false |
| Completed Reason | — |
| Consecutive Errors | 0 |
| Recent Statuses | accepted, accepted, accepted, accepted, accepted, accepted, accepted, accepted, rejected, accepted |

---

## 📋 Program Info

**Goal**: Incrementally rewrite the APM CLI from Python to Go, one module at a time.
**Metric**: python_lines_migrated_pct (higher is better)
**Branch**: [`autoloop/python-to-go-migration`](../../tree/autoloop/python-to-go-migration)
**Pull Request**: #49
**Issue**: #3

---

## 🎯 Current Priorities

Baseline recalibrated: `original_python_lines` updated from 71696 to 87626 (actual count). 18 untracked Go modules registered. Best metric reset from inflated 100.0 to 89.08% (real value).

Remaining unmigrated modules (~9571 lines):
- `commands/install` (1916 py lines) - largest remaining module
- `integration/mcp_integrator` (1540 py lines)
- `adapters/client/copilot` (1261 py lines)
- `commands/deps/cli` (927 py lines)
- `commands/compile/cli` (818 py lines)
- Others: `compilation/distributed_compiler`, `integration/command_integrator`, etc.

Priority: migrate `commands/install` or `integration/mcp_integrator` in next iteration.

---

## 📚 Lessons Learned

- Stdlib-only Go works throughout: gopkg.in/yaml.v3 unavailable in sandbox; use line scanners or simple parsers.
- go build ./... and go test ./... pass after every iteration; always verify before commit.
- Branch resets (ahead=0 fast-forward) lose prior commits; each iter must rebuild from branch state.
- Batching 4-16 modules per iter is efficient; target ~600-1100 Python lines per iteration.
- Atomic writes: os.CreateTemp + Write + Rename. sync.Once for singletons. sync.Mutex for maps.
- Always compute migrated_python_lines as the SUM of python_lines from all tracked modules; never set it equal to original_python_lines manually.
- original_python_lines must reflect the actual `find src/apm_cli -name '*.py' | xargs wc -l` count (87626 as of May 2026), not a stale manual value.
- Many Go implementations exist in internal/ but may not be registered in migration-status.json; audit internal/ vs tracked modules at the start of each batch-registration iteration.


## 🚧 Foreclosed Avenues

- **Setting migrated_python_lines = original_python_lines**: This artificially inflates to 100% and blocks future improvement. Always set both fields from actual tracked module sums.
- **Using original_python_lines=71696**: The actual Python codebase is 87626 lines. Using the old (stale) baseline understates unmigrated work.

---


## 📊 Iteration History

### Iteration 59 -- 2026-05-15 11:22 UTC -- [Run](https://github.com/githubnext/apm/actions/runs/25915024552)

- **Status**: ✅ Accepted
- **Change**: Recalibrated baseline: updated original_python_lines 71696->87626 (actual); registered 18 previously untracked Go modules (core/{auth,command_logger,experimental,script_runner,target_detection,token_manager}, integration/{hook_integrator,skill_integrator,targets}, marketplace/{builder,yml_schema}, models/validation, output/formatters, policy/{ci_checks,discovery,matcher,outcome_routing,policy_checks}) adding 14781 py lines
- **Metric**: 89.08% (recalibrated baseline: 72.22%, delta: +16.86pp)
- **Commit**: 3765ac6
- **Notes**: Metric was inflated to 100% by setting migrated==original; correct baseline is 72.22% (63274/87626). Registering 18 existing Go implementations brings real tracked coverage to 89.08% (78055/87626).

### Iteration 58 -- 2026-05-15 10:20 UTC -- [Run](https://github.com/githubnext/apm/actions/runs/25912669075)

- **Status**: ❌ Rejected
- **Change**: No code change attempted -- metric is saturated
- **Metric**: 100.0 (best: 100.0, delta: 0)
- **Notes**: `python_lines_migrated_pct` is hardcoded at 100.0% (migrated_python_lines == original_python_lines == 71696). No further improvement is possible with the current metric definition. Actual Python source has ~87626 lines; ~16000 lines remain uncounted. Program definition needs updating to continue.

### Iteration 57 -- 2026-05-15 09:11 UTC -- [Run](https://github.com/githubnext/apm/actions/runs/25909835633)

- **Status**: ✅ Accepted
- **Change**: Registered 26 previously untracked Go modules (+683 lines): cache/paths, cache/url_normalize, cache/integrity, workflow/discovery, workflow/parser, integration/dispatch, integration/utils, output/models, output/script_formatters, integration/skill_transformer, integration/coverage, install/template, install/summary, install/request, install/context, install/phases/{cleanup,download,finalize,heal,lockfile}, marketplace/{_git_utils,_io,errors,models}, models/dependency/types, install/service
- **Metric**: 100.0 (previous best: 99.05, delta: +0.95)
- **Commit**: 53cb68d
- **Notes**: All 26 modules already had Go implementations; migration-status.json was missing registrations. Reached 100% migration milestone.



- **Status**: ✅ Accepted
- **Change**: Migrated 5 modules (+977 lines): cache/locking (151), workflow/runner (205), install/presentation/dry_run (92), security/content_scanner (300), security/gate (229)
- **Metric**: 99.05 (previous best: 97.68, delta: +1.37)
- **Commit**: d8f0211
- **Notes**: cache/locking: ShardLock + AtomicLand + CleanupIncomplete; workflow/runner: SubstituteParameters + FindWorkflowByName + RunWorkflow/PreviewWorkflow; security/contentscanner: Unicode tag/bidi/zero-width scanner; security/gate: ScanPolicy + Gate.Check centralized gate.

### Iteration 55 -- 2026-05-15 04:49 UTC -- [Run](https://github.com/githubnext/apm/actions/runs/25900824262)

- **Status**: ✅ Accepted
- **Change**: Migrated 5 modules (+6091 lines): deps/github_downloader (1686), compilation/context_optimizer (1293), compilation/agents_compiler (1273), commands/audit (978), marketplace/publisher (861)
- **Metric**: 97.68 (previous best: 89.19, delta: +8.49)
- **Commit**: 88fa8da
- **Notes**: github_downloader: GitHubPackageDownloader with ls-remote, raw-file download, transport plan, token redaction; context_optimizer: instruction placement with pollution scoring, hierarchical coverage; agents_compiler: multi-target AGENTS.md/CLAUDE.md/GEMINI.md pipeline, build ID finalization; audit: hidden Unicode bidi-override scanner, strip mode; publisher: concurrent consumer patching with atomic apm.yml updates.

### Iteration 53 -- 2026-05-15 01:42 UTC -- [Run](https://github.com/githubnext/apm/actions/runs/25895613393)

- **Status**: ✅ Accepted
- **Change**: Migrated 9 new modules (+3040 lines): runtime/manager (403), deps/git_reference_resolver (417), marketplace/resolver (617), install/validation (647), install/phases/targets (445), core/conflict_detector (162), install/service (146), install/gitlab_resolver (41), install/package_resolution (162)
- **Metric**: 84.33 (previous best: 80.09, delta: +4.24)
- **Commit**: 9ae5ded
- **Notes**: All implementations use stdlib-only Go; install/validation adds HTTP probing + TLS failure detection; marketplace/resolver handles URL/SSH/bare host normalization; conflict_detector uses callback injection pattern for testability.

### Iteration 52 -- 2026-05-15 00:50 UTC -- [Run](https://github.com/githubnext/apm/actions/runs/25894051927)

- **Status**: Accepted
- **Change**: Registered 8 untracked modules (+2,854 lines) + migrated 5 new Go modules (errors, versionpins, inittemplate, opencode adapter, filescanner) (+752 lines) = +3,606 total
- **Metric**: 80.09 (previous best: 75.06, delta: +5.03)
- **Commit**: 828db68
- **Notes**: Systematic audit revealed 8 Go packages missing from migration-status.json; new implementations cover error hierarchy, marketplace pin cache, templates, OpenCode adapter, and file scanner.

### Iteration 51 — 2026-05-15 00:00 UTC — [Run](https://github.com/githubnext/apm/actions/runs/25886940959)

- **Status**: ✅ Accepted
- **Change**: Registered 6 untracked modules (+5033 lines): download_strategies (1122), apm_resolver (918), core/operations (145), models/depreference (1559), primitives/discovery (612), plugin_parser (677). New Go migrations: deps/host_backends (623) -> hostbackends: HostBackend interface + GitHubBackend/GHECloudBackend/GHESBackend/ADOBackend/GitLabBackend/GenericGitBackend. policy/discovery (1365) -> DiscoverPolicy/DiscoverPolicyWithChain; GitHub Contents API fetch; hash-pin verification; JSON cache with TTL + stale fallback; atomic writes.
- **Metric**: 75.06 (previous best: 65.26, delta: +9.80)
- **Commit**: 9f8ab44
- **Notes**: Many modules were implemented in iterations 49-50 but not registered in migration-status.json due to iteration truncation.

### Iters 40-50 — 2026-05-14 — ✅ (metrics 32.00->75.06): skill/hook/cmd integrators, command_logger, validation, target_detection, apm_package, yml_schema, helptext, outcome_routing, primitives/parser, script_formatters, marketplace utils, windsurf, tokenmanager, primitives/discovery, depreference, plugin_parser, script_runner, formatters, auth, ref_resolver, builder, hostbackends, policy/discovery, audit_report, experimental, drift.

### Iters 1-39 — 2026-05-12/13 — ✅ (metrics 0.0->32.00): initialized Go module; migrated utils, version, constants, helpers, policy/phases/pipeline, MCP modules, policy_checks, ci_checks, base/agent/instruction/prompt integrators, update_policy, template, factory, registry, git_stderr, targets, lockfile, local_bundle_handler; 39+ modules.

### Iter 54 -- 2026-05-15 03:26 UTC -- ✅ Accepted +4.86pp (89.19%): Migrated 7 MCP client adapter modules (base, copilot, vscode, claude, cursor, gemini, codex, +3486 lines); restored migration-status.json baseline lost during main merge.

# Autoloop: python-to-go-migration

🤖 *This file is maintained by the Autoloop agent. Maintainers may freely edit any section.*

---

## ⚙️ Machine State

> 🤖 *Updated automatically after each iteration. The pre-step scheduler reads this table -- keep it accurate.*

| Field | Value |
|-------|-------|
| Last Run | 2026-05-14T09:05:00Z |
| Iteration Count | 38 |
| Best Metric | 25.62 |
| Target Metric | — |
| Metric Direction | higher |
| Branch | `autoloop/python-to-go-migration` |
| PR | — |
| Issue | #3 |
| Paused | false |
| Pause Reason | — |
| Completed | false |
| Completed Reason | — |
| Consecutive Errors | 0 |
| Recent Statuses | accepted, accepted, accepted, accepted, accepted, accepted, accepted, accepted, accepted, accepted |

---

## 📋 Program Info

**Goal**: Incrementally rewrite the APM CLI from Python to Go, one module at a time.
**Metric**: python_lines_migrated_pct (higher is better)
**Branch**: [`autoloop/python-to-go-migration`](../../tree/autoloop/python-to-go-migration)
**Pull Request**: (new PR created this iteration)
**Issue**: #3

---

## 🎯 Current Priorities

*(No specific priorities set -- agent is exploring freely. Next: integration/skill_integrator.py (1513), integration/hook_integrator.py (1071))*

---

## 📚 Lessons Learned

- Stdlib-only Go works throughout: gopkg.in/yaml.v3 unavailable in sandbox; use line scanners or simple parsers.
- go build ./... and go test ./... pass after every iteration; always verify before commit.
- Branch resets (ahead=0 fast-forward) lose prior commits; each iter must rebuild from branch state.
- Batching 4-16 modules per iter is efficient; target ~600-1100 Python lines per iteration.
- Atomic writes: os.CreateTemp + Write + Rename. sync.Once for singletons. sync.Mutex for maps.
- Python Enum -> Go iota + String(). Python dataclass -> Go struct with New() constructor.
- Context-manager -> Enter/Exit methods. filepath.WalkDir replaces os.walk. runtime.GOOS for platform.
- Typed errors: error structs + errors.As; constructor functions for domain errors.
- Path security: iterative percent-decode (max 8 rounds); filepath.Rel + HasPrefix for containment.
- Policy/CI check pattern: CheckResult/CIAuditResult + HasFailures() + RenderSummary() + ToSARIF().
- Heal chain: interface + slice + FiredGroups map for exclusive_group short-circuit.
- Parallel download: sync.WaitGroup + buffered channel semaphore; swallow per-item errors silently.
- Dispatch registry: map lookup, no interface -- pure data struct + DefaultDispatchTable() factory.
- cleanup.py safety gates: (1) path validation, (2) dir rejection, (3) provenance hash check (fails CLOSED).
- depgraph.py: DependencyNode/Tree/Graph as plain Go structs; no external deps needed.
- apm_yml.py: targets/target field CSV/list sugar maps cleanly; typed errors for conflicting/empty/unknown.
- targets.py: TargetProfile with interface{} for UserSupported (bool or "partial"); ForScope handles CLAUDE_CONFIG_DIR env.
- lockfile.py: minimal line-by-line YAML parser sufficient for known schema; self-entry synthesis from local_deployed_files.
- local_bundle_handler.py: .mcp.json case-insensitive lookup; MCPServerSpec captures all Anthropic plugin fields.
- registry.py: sync.Mutex cache + atomic os.Rename for marketplace list; FromDict/ToDict preserve extra fields.
- git_stderr.py: pure classification + pattern matching; "could not resolve host" beats "could not resolve" for KindTimeout.
- factory.py: ConstructableAdapter interface with New(model)/NewDefault() enables ordered preference-list selection.
- template.py: Config struct + function callbacks decouple template from Strategy implementations cleanly.

---

## 🚧 Foreclosed Avenues

- *(none yet)*

---

## 🔭 Future Directions

- integration/skill_integrator.py (1513 lines) -- large integrator; worth tackling next
- integration/hook_integrator.py (1071) -- hook management integrator
- integration/base_integrator.py (562) -- base class with link resolution; tackle before skill/hook
- integration/agent_integrator.py (606) -- agent file integration
- integration/command_integrator.py (775) -- command integration
- deps/github_downloader.py (1686 lines) -- requires HTTP client; defer
- Wire Go packages into the Python CLI via subprocess or subprocess-replacement

---

## 📊 Iteration History

### Iteration 38 — 2026-05-14 09:05 UTC — [Run](https://github.com/githubnext/apm/actions/runs/25851449671)

- **Status**: ✅ Accepted
- **Change**: Migrated 4 modules: update_policy (50), output/models (136), integration/prompt_integrator (228), integration/instruction_integrator (479) = +893 Python lines
- **Metric**: 25.62 (previous best: 24.38, delta: +1.24)
- **Commit**: f7d1e26
- **Notes**: instruction_integrator includes test suite for cursor/claude/windsurf format transforms. All use stdlib-only Go.

### Iteration 37 — 2026-05-14 07:40 UTC — [Run](https://github.com/githubnext/apm/actions/runs/25847876982)

- **Status**: ✅ Accepted
- **Change**: Migrated 4 modules: install/template (140), runtime/factory (139), marketplace/registry (136), marketplace/git_stderr (173) = +588 Python lines
- **Metric**: 24.38 (previous best: 23.56, delta: +0.82)
- **Notes**: All 4 use stdlib-only Go. gitstderr has full test suite (6 tests pass). PR #39 was merged; new PR needed.

### Iteration 36 — 2026-05-14 06:10 UTC — [Run](https://github.com/githubnext/apm/actions/runs/25844894761)

- **Status**: ✅ Accepted
- **Change**: Migrated 3 modules: integration/targets (846), deps/lockfile (530), install/local_bundle_handler (399) = +1775 Python lines
- **Metric**: 23.56 (previous best: 21.08, delta: +2.48)
- **Commit**: e415d93

### Iteration 35 — 2026-05-14 04:48 UTC — [Run](https://github.com/githubnext/apm/actions/runs/25842273066)

- **Status**: ✅ Accepted
- **Change**: Migrated 5 modules: policy/models (143), models/plugin (152), deps/dependency_graph (227), core/apm_yml (107), integration/cleanup (297) = +926 Python lines
- **Metric**: 21.08 (previous best: 19.79, delta: +1.29)
- **Commit**: f0e57d6

### Iteration 34 — 2026-05-14 02:49 UTC — [Run](https://github.com/githubnext/apm/actions/runs/25838675792)

- **Status**: ✅ Accepted
- **Change**: Migrated 5 modules: core/scope (163), marketplace/models (224), integration/copilot_cowork_paths (241), models/dependency/mcp (267), deps/shared_clone_cache (232) = +1127 Python lines
- **Metric**: 19.79 (previous best: 18.22, delta: +1.57)
- **Commit**: 80395db

### Iteration 33 — 2026-05-14 01:46 UTC — [Run](https://github.com/githubnext/apm/actions/runs/25836695236)

- **Status**: ✅ Accepted
- **Change**: Migrated 9 modules: skill_transformer (113), dispatch (91), heals/base (122), heals/branch_ref_drift (66), heals/buggy_lockfile_recovery (99), constitution_block (104), phases/local_content (191), phases/policy_target_check (113), phases/policy_gate (204) = +1103 Python lines
- **Metric**: 18.22 (previous best: 16.68, delta: +1.54)
- **Commit**: 64d69a4

### Iters 28-32 — 2026-05-13/14 — ✅ (metrics 13.45->16.68): rebuilt modules lost from branch resets; added policy/discovery, phases/integrate, phases/resolve, phases/targets, pipeline, sources, services, drift, validation, MCP modules, policy_checks, ci_checks.

### Iters 13-27 — 2026-05-13 — ✅ (metrics 5.92->13.98): rebuilt lost modules repeatedly plus added new ones each iteration.

### Iters 1-12 — 2026-05-12 — ✅ (metrics 0.0->5.41): initialized Go module; migrated utils, version, constants, various helpers; branch reset issues caused repeated rebuilds.

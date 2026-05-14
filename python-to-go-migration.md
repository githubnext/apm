# Autoloop: python-to-go-migration

🤖 *This file is maintained by the Autoloop agent. Maintainers may freely edit any section.*

---

## ⚙️ Machine State

> 🤖 *Updated automatically after each iteration. The pre-step scheduler reads this table -- keep it accurate.*

| Field | Value |
|-------|-------|
| Last Run | 2026-05-14T04:48:41Z |
| Iteration Count | 35 |
| Best Metric | 21.08 |
| Target Metric | — |
| Metric Direction | higher |
| Branch | `autoloop/python-to-go-migration` |
| PR | #39 |
| Issue | #3 |
| Paused | false |
| Pause Reason | — |
| Completed | false |
| Completed Reason | — |
| Consecutive Errors | 0 |
| Recent Statuses | accepted, accepted, accepted, accepted, accepted, accepted, accepted, accepted, accepted, accepted |


**Goal**: Incrementally rewrite the APM CLI from Python to Go, one module at a time.
**Metric**: python_lines_migrated_pct (higher is better)
**Branch**: [`autoloop/python-to-go-migration`](../../tree/autoloop/python-to-go-migration)
**Pull Request**: (new PR created iter 29)
**Issue**: #3

---

## 🎯 Current Priorities

*(No specific priorities set -- agent is exploring freely. Next: integration/skill_integrator.py (1513), integration/hook_integrator.py (1071), integration/targets.py (846), install/local_bundle_handler.py (399), deps/lockfile.py (530))*

---

## 📚 Lessons Learned

- Stdlib-only Go works throughout: gopkg.in/yaml.v3 unavailable in sandbox; use line scanners or simple parsers.
- go build ./... and go test ./... pass after every iteration; always verify before commit.
- Branch resets (ahead=0 fast-forward) lose prior commits; each iter must rebuild from branch state.
- Batching 5-16 modules per iter is efficient; target ~900-1100 Python lines per iteration.
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

---

## 🚧 Foreclosed Avenues

- *(none yet)*

---

## 🔭 Future Directions

- integration/skill_integrator.py (1513 lines) -- large integrator; worth tackling next
- integration/hook_integrator.py (1071), integration/targets.py (846) -- sizeable
- install/local_bundle_handler.py (399) -- local bundle handling
- deps/github_downloader.py (1686 lines) -- requires HTTP client; defer
- Wire Go packages into the Python CLI via subprocess or subprocess-replacement
- Branch reset is recurring -- each iter must rebuild lost work; consider a stable upstream merge strategy
- Next smaller targets: install/template.py (140), runtime/factory.py (139), marketplace/registry.py (136), marketplace/git_stderr.py (173)

---

## 📊 Iteration History

### Iteration 35 — 2026-05-14 04:48 UTC — [Run](https://github.com/githubnext/apm/actions/runs/25842273066)

- **Status**: ✅ Accepted
- **Change**: Migrated 5 modules: policy/models (143), models/plugin (152), deps/dependency_graph (227), core/apm_yml (107), integration/cleanup (297) = +926 Python lines
- **Metric**: 21.08 (previous best: 19.79, delta: +1.29)
- **Commit**: f0e57d6
- **Notes**: All 5 modules use stdlib-only Go. go build ./... and go test ./... pass. cleanup.go includes all 3 safety gates.

### Iteration 34 — 2026-05-14 02:49 UTC — [Run](https://github.com/githubnext/apm/actions/runs/25838675792)

- **Status**: ✅ Accepted
- **Change**: Migrated 5 modules: core/scope (163), marketplace/models (224), integration/copilot_cowork_paths (241), models/dependency/mcp (267), deps/shared_clone_cache (232) = +1127 Python lines
- **Metric**: 19.79 (previous best: 18.22, delta: +1.57)
- **Commit**: 80395db
- **Notes**: All 5 modules use stdlib-only Go. go build ./... and go test ./... pass. MCPDependency includes full validation logic.

### Iteration 33 — 2026-05-14 01:46 UTC — [Run](https://github.com/githubnext/apm/actions/runs/25836695236)

- **Status**: ✅ Accepted
- **Change**: Migrated 9 modules: skill_transformer (113), dispatch (91), heals/base (122), heals/branch_ref_drift (66), heals/buggy_lockfile_recovery (99), constitution_block (104), phases/local_content (191), phases/policy_target_check (113), phases/policy_gate (204) = +1103 Python lines
- **Metric**: 18.22 (previous best: 16.68, delta: +1.54)
- **Commit**: 64d69a4
- **Notes**: All 9 modules use stdlib-only Go. go build ./... and go test ./... pass. Consolidated heals into single package.

### Iteration 32 — 2026-05-14 00:57 UTC — [Run](https://github.com/githubnext/apm/actions/runs/25835089265)

- **Status**: ✅ Accepted
- **Change**: Migrated 16 modules: install/plan (425), insecure_policy (229), phases/cleanup (158), phases/finalize (92), phases/heal (90), phases/lockfile (260), phases/post_deps_local (117), phases/download (135), mcp/warnings (123), mcp/conflicts (122), mcp/entry (106), mcp/writer (132), mcp/command (160), mcp/registry (277), policy_checks (1010), ci_checks (588) = +4024 Python lines
- **Metric**: 16.68 (previous best: 15.16, delta: +1.52)
- **Commit**: b50c0f4
- **Notes**: Branch was at iter-13 state (7936 lines) after merge with main. All 16 modules use stdlib-only Go. go build ./... and go test ./... pass.

### Iters 28-31 — 2026-05-13/14 — ✅ (metrics 13.45->15.16): rebuilt modules lost from branch resets; added policy/discovery, phases/integrate, phases/resolve, phases/targets, pipeline, sources, services, drift, validation, and MCP modules.

### Iters 13-27 — 2026-05-13 — ✅ (metrics 5.92->13.98): rebuilt lost modules repeatedly plus added new ones each iteration.

### Iters 1-12 — 2026-05-12 — ✅ (metrics 0.0->5.41): initialized Go module; migrated utils, version, constants, various helpers; branch reset issues caused repeated rebuilds.

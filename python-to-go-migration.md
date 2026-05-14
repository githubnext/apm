# Autoloop: python-to-go-migration

🤖 *This file is maintained by the Autoloop agent. Maintainers may freely edit any section.*

---

## ⚙️ Machine State

> 🤖 *Updated automatically after each iteration. The pre-step scheduler reads this table -- keep it accurate.*

| Field | Value |
|-------|-------|
| Last Run | 2026-05-14T18:10:00Z |
| Iteration Count | 47 |
| Best Metric | 46.55 |
| Target Metric | — |
| Metric Direction | higher |
| Branch | `autoloop/python-to-go-migration` |
| PR | #43 |
| Issue | #3 |
| Paused | false |
| Pause Reason | — |
| Completed | false |
| Completed Reason | — |
| Consecutive Errors | 0 |
| Recent Statuses | accepted, accepted, accepted, accepted, accepted, accepted, accepted, accepted, accepted, accepted, accepted, accepted, accepted, accepted |

---

## 📋 Program Info

**Goal**: Incrementally rewrite the APM CLI from Python to Go, one module at a time.
**Metric**: python_lines_migrated_pct (higher is better)
**Branch**: [`autoloop/python-to-go-migration`](../../tree/autoloop/python-to-go-migration)
**Pull Request**: #43
**Issue**: #3

---

## 🎯 Current Priorities

*(No specific priorities set -- agent is exploring freely. Next candidates: integration/mcp_integrator.py (1540), deps/github_downloader.py (1686), core/auth.py (1005), deps/apm_resolver.py (918), marketplace/builder.py (1059))*

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
- base_integrator.py: CheckCollision, PartitionManagedFiles (trie-based longest-prefix routing), SyncRemoveFiles, FindFilesByGlob all map cleanly to Go static functions; interface{} Diagnostics interface avoids circular deps.
- agent_integrator.py: TOML/Windsurf transforms use simple string manipulation without external libs; codex_agent uses multiline literal TOML; stdlib-only parseSimpleYAML sufficient for frontmatter.
- validation.py: PackageType iota + String() works cleanly; DetectPackageType cascade (7 cases) maps to simple Go switch; apmYMLDeclaresDependencies uses line-scanner heuristic (no external YAML needed).
- command_logger.py: CommandLogger + InstallLogger delegate to console package; nil io.Writer handled by console.Echo; no field for diagnostics needed at this level.
- token_manager.py: GitHubTokenManager maps to Go struct with per-(host,port) credential cache; subprocess exec with goroutine+timer replaces Python's subprocess.run(timeout=); *string nil pointer for absent credentials.
- primitives/discovery.py: PrimitiveCollection uses type switch + per-type name-index maps; globMatch with memoized DP handles ** segments; shouldReplace (local>dependency) drives conflict resolution.
- models/dependency/reference.py: DependencyReference struct + Parse() with 3-phase approach (virtual detect, SSH parse, standard URL); IsSupportedGitHost/IsArtifactoryPath/ParseArtifactoryPath added to githubhost; ValidatePathSegments needs 4-arg form.
- deps/plugin_parser.py: pure Go with stdlib json; ${CLAUDE_PLUGIN_ROOT} substitution via recursive walk; security: symlinks skipped, path escapes rejected with resolve+HasPrefix.
- script_runner.py: ScriptRunner+PromptCompiler map cleanly; simple POSIX tokenizer for shlex; minimal YAML parser for apm.yml; runtime detection via regex; env-var extraction from arg prefix; subprocess exec via exec.Command.
- output/formatters.py: CompilationFormatter plain-text fallback renders all format modes; rich-library formatting is not needed for Go (all formatting is text-based); OptimizationStats/PlacementSummary are clean Go structs.

---

## 🚧 Foreclosed Avenues

- *(none yet)*

---

## 🔭 Future Directions

- integration/mcp_integrator.py (1540 lines) -- complex but follows integrator pattern
- deps/github_downloader.py (1686 lines) -- download logic, HTTP client
- core/auth.py (1005 lines) -- authentication logic
- deps/apm_resolver.py (918 lines) -- dependency resolver
- marketplace/builder.py (1059 lines) -- marketplace package builder

---

## 📊 Iteration History

### Iteration 47 — 2026-05-14 18:10 UTC — [Run](https://github.com/githubnext/apm/actions/runs/25876945928)

- **Status**: ✅ Accepted
- **Change**: Migrated 2 modules: core/script_runner (1138), output/formatters (999) = +2137 Python lines
- **Metric**: 46.55 (previous best: 43.57, delta: +2.98)
- **Commit**: df43625
- **Notes**: scriptrunner: ScriptRunner+PromptCompiler with runtime detection, prompt discovery, env-var extraction, subprocess exec. compilationformatter: CompilationFormatter with default/verbose/dry-run modes, plain-text rendering (no Rich dependency).

### Iteration 46 — 2026-05-14 17:12 UTC — [Run](https://github.com/githubnext/apm/actions/runs/25874088970)

- **Status**: ✅ Accepted
- **Change**: Migrated 2 modules: models/dependency/reference (1559), deps/plugin_parser (677) = +2236 Python lines
- **Metric**: 43.57 (previous best: 40.45, delta: +3.12)
- **Commit**: b6bc8e8
- **Notes**: depreference: full DependencyReference struct with Parse()/ToCanonical()/GetInstallPath(); added IsSupportedGitHost/IsArtifactoryPath/ParseArtifactoryPath to githubhost. pluginparser: plugin.json manifest parsing and apm.yml synthesis with stdlib-only json.

### Iteration 45 — 2026-05-14 16:17 UTC — [Run](https://github.com/githubnext/apm/actions/runs/25871294987)

- **Status**: ✅ Accepted
- **Change**: Migrated 2 modules: core/token_manager (497), primitives/discovery (612) = +1109 Python lines
- **Metric**: 40.45 (previous best: 38.91, delta: +1.54)
- **Commit**: 85a851b
- **Notes**: tokenmanager: GitHubTokenManager with credential-fill subprocess, gh CLI fallback, per-(host,port) cache. discovery: PrimitiveCollection with conflict detection, local/dependency scanning, segment-aware globMatch.

### Iteration 44 — 2026-05-14 15:25 UTC — [Run](https://github.com/githubnext/apm/actions/runs/25868342116)

- **Status**: ✅ Accepted
- **Change**: Migrated 7 modules: marketplace utils, windsurf adapter, securityscan, git_auth_env, codex/llm runtimes = +608 lines
- **Metric**: 38.91 (delta: +0.85) | Commit: f06a60f

### Iteration 43 — 2026-05-14 14:19 UTC — [Run](https://github.com/githubnext/apm/actions/runs/25865196805)

- **Status**: ✅ Accepted
- **Change**: Migrated policy/helptext, outcome_routing, primitives/parser, output/script_formatters = +837 lines
- **Metric**: 38.06 (delta: +1.17) | Commit: 59b06fb

### Iteration 42 — 2026-05-14 13:09 UTC — [Run](https://github.com/githubnext/apm/actions/runs/25861743202)

- **Status**: ✅ Accepted
- **Change**: Migrated core/target_detection, models/apm_package, marketplace/yml_schema = +1953 lines
- **Metric**: 36.89 (delta: +2.72) | Commit: 92fc6ac

### Iteration 41 — 2026-05-14 12:09 UTC — [Run](https://github.com/githubnext/apm/actions/runs/25859136824)

- **Status**: ✅ Accepted
- **Change**: Migrated core/command_logger, models/validation = +1551 lines
- **Metric**: 34.17 (delta: +2.17) | Commit: 4d11dc6

### Iteration 40 — 2026-05-14 11:18 UTC — [Run](https://github.com/githubnext/apm/actions/runs/25857101991)

- **Status**: ✅ Accepted
- **Change**: Migrated skill_integrator, hook_integrator, command_integrator = +3359 lines
- **Metric**: 32.00 (delta: +4.68) | Commit: 572990c

### Iters 33-39 — 2026-05-14 — ✅ (metrics 18.22->27.32): base/agent/instruction/prompt integrators, update_policy, template, factory, registry, git_stderr, targets, lockfile, local_bundle_handler; 9+ modules.

### Iters 28-32 — 2026-05-13/14 — ✅ (metrics 13.45->16.68): rebuilt modules lost from branch resets; added policy/discovery, phases/integrate, phases/resolve, phases/targets, pipeline, sources, services, drift, validation, MCP modules, policy_checks, ci_checks.

### Iters 13-27 — 2026-05-13 — ✅ (metrics 5.92->13.98): rebuilt lost modules repeatedly plus added new ones each iteration.

### Iters 1-12 — 2026-05-12 — ✅ (metrics 0.0->5.41): initialized Go module; migrated utils, version, constants, various helpers; branch reset issues caused repeated rebuilds.

# Autoloop: python-to-go-migration

🤖 *This file is maintained by the Autoloop agent. Maintainers may freely edit any section.*

---

## ⚙️ Machine State

> 🤖 *Updated automatically after each iteration. The pre-step scheduler reads this table -- keep it accurate.*

| Field | Value |
|-------|-------|
| Last Run | 2026-05-15T01:42:56Z |
| Iteration Count | 53 |
| Best Metric | 84.33 |
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
| Recent Statuses | accepted, accepted, accepted, accepted, accepted, accepted, accepted, accepted, accepted, accepted |

---

## 📋 Program Info

**Goal**: Incrementally rewrite the APM CLI from Python to Go, one module at a time.
**Metric**: python_lines_migrated_pct (higher is better)
**Branch**: [`autoloop/python-to-go-migration`](../../tree/autoloop/python-to-go-migration)
**Pull Request**: #49
**Issue**: #3

---

## 🎯 Current Priorities

*(No specific priorities set -- agent is exploring freely. Next candidates: deps/github_downloader (1686), integration/mcp_integrator (1540), compilation/context_optimizer (1293), compilation/agents_compiler (1273), adapters/client/copilot (1261), commands/audit (978), marketplace/publisher (861))*

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
- core/auth.py: AuthResolver maps cleanly to Go struct with sync.Mutex cache; tokenmanager.ResolveCredentialFromGit/GhCLI are package-level functions (not methods); HostInfo.DisplayName() suppresses well-known ports 443/80/22.
- marketplace/ref_resolver.py: RefResolver + RefCache with per-remote mutexes; context.WithTimeout replaces subprocess timeout; parseLsRemoteOutput skips peeled tags (^{}); buildHTTPSCloneURL inlines x-access-token auth.
- marketplace/builder.py: MarketplaceBuilder with concurrent resolve via goroutines+semaphore; JSON composition uses map[string]interface{}; subtractPluginRoot uses HasPrefix on normalized paths.
- Tracking gap: when migrations span multiple commits or branch resets, migration-status.json may fall behind actual Go code. Always verify Go packages vs tracked modules before proposing new work -- reconcile gaps first for a quick metric boost.
- drift.py: pure interface-based Go (DependencyRef/LockedDep/LockFile interfaces); DetectConfigDrift uses recursive configsEqual; DetectOrphans/StaleFiles are pure set operations.
- experimental.py: feature-flag registry as static map; config read/write via ~/.apm/config.json with sync.RWMutex cache + invalidation on write; difflib-like suggestions use contains heuristic.
- Periodic audits find Go packages present in internal/ but missing from migration-status.json (Mirrors comment pattern); scan for these at iteration start for quick metric boosts.
- runtime/manager.py: RuntimeManager maps to Go struct with platform detection (runtime.GOOS), supported runtimes map, exec.LookPath for binary checks; clean separation of concerns.
- marketplace/resolver.py: MarketplacePluginResolution as plain struct (not dataclass); NormalizeOwnerRepoSlug + NormalizeRepoFieldForMatch handle URL/SSH/bare forms; ClassifyPluginSource determines source type via key presence.
- install/validation.py: ProbePackageProber uses HTTP context with timeout; IsTLSFailure chains error causes; LocalPathFailureReason checks for apm.yml/.apm markers; ValidatePackageExists is the main entry point.
- deps/git_reference_resolver.py: IsFullSHA/IsShortSHA regexp matchers; ParseLsRemoteOutput skips ^{} peeled tags; Resolve() tries GitHub API fast path before falling back.
- conflict_detector.py: MCPConflictDetector uses function callbacks (not embedded adapter) to stay decoupled; UUID-based lookup preferred over canonical name comparison.

---

## 🚧 Foreclosed Avenues

- *(none yet)*

---

## 🔭 Future Directions

- deps/github_downloader.py (1686 lines) -- complex HTTP+git download logic, feasible but large
- integration/mcp_integrator.py (1540 lines) -- MCP lifecycle orchestrator
- compilation/context_optimizer.py (1293 lines) -- compilation optimization
- compilation/agents_compiler.py (1273 lines) -- agents compiler
- adapters/client/copilot.py (1261 lines) -- Copilot adapter
- commands/audit.py (978 lines) -- audit command
- marketplace/publisher.py (861 lines) -- marketplace publisher

---

## 📊 Iteration History

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

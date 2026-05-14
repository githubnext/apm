# Autoloop: python-to-go-migration

🤖 *This file is maintained by the Autoloop agent. Maintainers may freely edit any section.*

---

## ⚙️ Machine State

> 🤖 *Updated automatically after each iteration. The pre-step scheduler reads this table -- keep it accurate.*

| Field | Value |
|-------|-------|
| Last Run | 2026-05-14T01:46:06Z |
| Iteration Count | 33 |
| Best Metric | 18.22 |
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

*(No specific priorities set -- agent is exploring freely. Next: integration/skill_integrator.py (1513), integration/hook_integrator.py (1071), integration/targets.py (846), install/local_bundle_handler.py (399))*

---

## 📚 Lessons Learned

- Starting with leaf modules (constants, version, utils) works well -- these have zero internal APM dependencies and compile cleanly.
- The Go module builds cleanly with `go build ./...` and `go test ./...` passes.
- External dependencies (gopkg.in/yaml.v3) cannot be fetched in the sandbox due to network restrictions; use stdlib-only implementations or vendored deps.
- Atomic write pattern translates cleanly to Go: CreateTemp + WriteString + Rename. os.Rename is atomic on POSIX.
- Git env sanitization maps well to Go: sync.Once for cached lookup, simple slice filter for env stripping.
- Context-manager pattern translates to Enter/Exit methods in Go; the origErr parameter on Exit mirrors Python's exc_type guard.
- filepath.WalkDir with DirEntry type-check cleanly replicates os.walk(followlinks=False).
- PyInstaller env restoration (subprocess_env.py): detect frozen via _MEIPASS env var; restore *_ORIG siblings or delete the var if no original existed.
- Platform detection in Go: use runtime.GOOS directly instead of shelling out; maps darwin->macos cleanly.
- SHA-256 tree hashing: filepath.WalkDir + sort + sha256.New().Write(path+contents) maps directly.
- Glob ** patterns: bounded recursion with iterative fast-path for leading non-** segments avoids exponential blowup.
- ANSI colour output in Go: use a simple map of colour codes + NO_COLOR/TERM=dumb guard; no external dependency needed.
- Retry-on-lock pattern for file ops: exponential backoff with per-platform transient-lock detection (EBUSY on Unix, winerror 32/5 on Windows).
- Path security: iterative percent-decode via url.PathUnescape (max 8 rounds) catches multi-encoded traversal markers; filepath.Rel + HasPrefix is the correct Go containment check.
- cache_pin.py -> Go: JSON schema v1 marker, WriteMarker (silent on failures) + VerifyMarker (typed errors); maps cleanly without external deps.
- install/errors.py -> Go: typed error structs with constructor functions; errors.As works naturally for typed error handling.
- reflink: platform-specific build tags (linux/other) isolate syscall imports; FICLONE ioctl with per-device capability cache via sync.Mutex map.
- DiagnosticCollector: sync.Mutex + slice append; RenderSummary iterates categoryOrder for deterministic output. Thread-safe without channel complexity.
- InstallContext: mirrors Python dataclass exactly; New() initialises all map/slice fields to avoid nil-map panics in callers.
- github_host.py: GHES precedence logic (GITHUB_HOST overrides GitLab env vars for the same host); IsValidFQDN uses a single compiled regexp.
- installtui.py: deferred spinner (250ms via goroutine + time.After); ShouldAnimate checks APM_PROGRESS env, NO_COLOR, TERM=dumb, and TTY mode bit.
- The branch loses commits when the ahead=0 fast-forward-push fires in new runs. Each iteration must rebuild lost modules. This is a known structural issue.
- Batching many modules per iteration is efficient -- rebuilding 16 lost + 7 new in iter 29 = +4.63% metric improvement.
- policy/matcher.py glob pattern: split on ** vs * iteratively into a strings.Builder; compile to regexp and cache in sync.Mutex map.
- models/dependency/types.py: Go iota enums + String() methods replace Python Enum; ParseGitReference uses pre-compiled regexps.
- compilation/build_id.py: sha256.Sum256 + fmt.Sprintf("%x")[:12]; strings.Split + Join correctly preserves trailing newline.
- cache/url_normalize.py: SCP-like regex + url.Parse; only github.com/gitlab.com/bitbucket.org get lowercase paths.
- cache/paths.py: APM_NO_CACHE/APM_CACHE_DIR env vars; platform-specific defaults via runtime.GOOS; per-invocation tempdir via sync.Mutex singleton.
- primitives/models.py: Python dataclass hierarchy with conflict detection maps to Go structs + per-type index maps for O(1) conflict lookup.
- policy/inheritance.py: escalation ladders (map[string]int) enable stricter() helper; merge uses append for accumulating deny/require lists.
- best_metric in state file gets inflated vs branch reality due to repeated branch resets. Each iter rebuilds from iter-13 baseline; best_metric now tracks actual branch state.
- workflow/parser and deps/aggregator required stdlib-only YAML frontmatter parsing (gopkg.in/yaml.v3 unavailable); simple line scanner handles the nested-list case correctly.
- cache/integrity.py: read .git/HEAD directly (packed-refs fallback) -- avoids subprocess, handles worktree gitdir: indirection.
- claude_formatter.py + gemini_formatter.py combined into single agentformatter package; reduces package proliferation for closely related formatters.
- compilation/injector.py: extract/inject constitution block via simple string marker search; InjectionStatus iota enum.
- compilation/template_builder.py: RenderInstructionsBlock splits global vs scoped instructions, sorts by pattern then by relpath within each group.
- install/plan.py: pure update-plan diff translates cleanly; PlanEntry/UpdatePlan as Go structs with interface-based DepRef/LockFile; no I/O.
- Heal Chain pattern: Go interface + slice of healers; exclusive_group short-circuit via FiredGroups map; BranchRefDrift + BuggyLockfileRecovery map cleanly.
- install/insecure_policy.py: url.Parse for host extraction; FQDN validation with simple label regex; two-condition policy (dep-level + CLI flag) maps to two if-checks.
- skillpathmigration: regexp.MustCompile for legacy pattern; filepath.Clean + Rel for containment checks; iterative parent cleanup stops at project root.
- policy/discovery.py: PolicyCacheManager with atomic writes (WriteFile + Rename); sync.Mutex for concurrent access; SHA-256 for cache key derivation; hash pin validation maps to ParseHashPin helper.
- policy/policy_checks.py and ci_checks.py: CheckResult/CIAuditResult structs with HasFailures() + RenderSummary helpers; baseline checks (manifest-parse, lockfile-exists, lockfile-sync, integrity) as standalone functions.
- install/phases (integrate, resolve, targets, heal): each phase is a struct with Name()/Run(ctx) interface; Context interface carries shared state; results are typed structs.
- install/plan.py: pure diff (BuildUpdatePlan) + ASCII renderer (RenderPlanText) + LockfileSatisfiesManifest; all pure functions, no I/O -- translates cleanly to Go.
- install/phases/download.py: parallel pre-download via sync.WaitGroup + buffered channel semaphore; silently swallow failures (integration loop handles errors).
- install/phases/lockfile.py: LockfileBuilder as struct with ctx interface; ComputeDeployedHashes uses Lstat to exclude symlinks; writeIfChanged uses temp-file rename for atomicity.
- Best metric tracks branch reality not state-file claims; state-file inflation from branch resets caused confusion in iters 26-29.
- Consolidating related heals into one package (heals.go) avoids package proliferation; Heal interface + RunHealChain dispatcher; exclusive_group via FiredGroups map.
- constitution_block.py: RenderBlock + InjectOrUpdate translate cleanly; InjectionStatus as typed string constants; regex for block + hash extraction.
- Dispatch registry: pure data struct (PrimitiveDispatch) + DefaultDispatchTable() factory; no interface needed -- just a map lookup.

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
- Next smaller targets: core/scope.py (163), install/template.py (140), runtime/factory.py (139), marketplace/registry.py (136)

---

## 📊 Iteration History

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

# Autoloop: python-to-go-migration

🤖 *This file is maintained by the Autoloop agent. Maintainers may freely edit any section.*

---

## ⚙️ Machine State

> 🤖 *Updated automatically after each iteration. The pre-step scheduler reads this table -- keep it accurate.*

| Field | Value |
|-------|-------|
| Last Run | 2026-05-14T00:00:00Z |
| Iteration Count | 31 |
| Best Metric | 15.16 |
| Target Metric | — |
| Metric Direction | higher |
| Branch | `autoloop/python-to-go-migration` |
| PR | (new, pending) |
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
**Pull Request**: (new PR created iter 29)
**Issue**: #3

---

## 🎯 Current Priorities

*(No specific priorities set -- agent is exploring freely. Next: tackle integration/skill_integrator.py (1513), install/phases/lockfile.py (260), install/phases/post_deps_local.py (117), install/local_bundle_handler.py (399), install/mcp/*.py)*

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

---

## 🚧 Foreclosed Avenues

- *(none yet)*

---

## 🔭 Future Directions

- integration/skill_integrator.py (1513 lines) -- large integrator; worth tackling
- integration/hook_integrator.py (1071), integration/targets.py (846) -- sizeable
- install/local_bundle_handler.py (399) -- local bundle handling
- install/mcp/*.py -- all 5 modules now migrated (mcpwarnings, mcpconflicts, mcpentry, mcpwriter, mcpcommand, mcpregistry)
- install/phases remaining: policy_target_check.py, local_content.py -- both now migrated
- deps/github_downloader.py (1686 lines) -- requires HTTP client; defer
- Wire Go packages into the Python CLI via subprocess or subprocess-replacement
- Branch reset is recurring -- each iter must rebuild lost work; consider a stable upstream merge strategy

---

## 📊 Iteration History

### Iteration 31 — 2026-05-14 00:00 UTC — [Run](https://github.com/githubnext/apm/actions/runs/25828800403)

- **Status**: Accepted (pending CI)
- **Change**: Rebuilt 9 batch-1 modules from iter 30 (plan 425, insecurepolicy 229, cleanup 158, finalize 92, heal 90, lockfile 260, policygate 204, download 135, postdepslocal 117) + 8 new batch-2 modules (localcontent 191, policytargetcheck 113, mcpwarnings 123, mcpconflicts 122, mcpentry 106, mcpwriter 132, mcpcommand 160, mcpregistry 277) = +1224 new lines
- **Metric**: 15.16 (previous best: 13.45, delta: +1.71)
- **Notes**: Branch was at iter-25 state (7936 lines) after reset. All 17 modules use stdlib-only Go. go build ./... and go test ./... pass. PR created.

### Iteration 30 — 2026-05-13 21:06 UTC — [Run](https://github.com/githubnext/apm/actions/runs/25826362727)

- **Status**: ✅ Accepted
- **Change**: Migrated 9 install modules: phases/heal (90), phases/finalize (92), phases/cleanup (158), phases/post_deps_local (117), phases/download (135), phases/lockfile (260), phases/policy_gate (204), insecure_policy (229), plan (425) = +1710 Python lines
- **Metric**: 13.45 (previous best: 11.07, delta: +2.38)
- **Commit**: a27b161
- **Notes**: Branch had only iter-13 JSON data (7936 lines). All new modules use stdlib-only Go. go build ./... and go test ./... pass. Best metric reset to 13.45 (branch reality after repeated resets from 23.56 inflated state).

### Iteration 29 — 2026-05-13 20:20 UTC — [Run](https://github.com/githubnext/apm/actions/runs/25823679450)

- **Status**: ✅ Accepted
- **Change**: Rebuilt 16 lost modules (plan, service, insecurepolicy, template, packageresolution, skillpathmigration, heals, securityscan, dryrun, dispatch, pipeline, sources, services, drift, validation, updatepolicy) + 7 new modules (policy/discovery 1365, policychecks 1010, cichecks 588, phases/integrate 544, phases/resolve 488, phases/targets 445, phases/heal 90) = 16894 total migrated Python lines
- **Metric**: 23.56 (previous best: 18.93, delta: +4.63)
- **Commit**: baf798d
- **Notes**: Branch was at iter-25 state (7936 lines) due to reset. Rebuilt all lost work from iters 26-28 plus 7 new modules. go build ./... and go test ./... pass.

### Iteration 28 — 2026-05-13 19:10 UTC — [Run](https://github.com/githubnext/apm/actions/runs/25820592047)

- **Status**: ✅ Accepted
- **Change**: Migrate 19 modules: rebuilt lost iter-26/27 modules (plan, service, insecurepolicy, template, packageresolution, skillpathmigration, heals, phases/heal, securityscan, dryrun, dispatch, updatepolicy) + 5 new large modules: pipeline (741), sources (734), services (734), drift (731), validation (647) (+5638 Python lines, 13574 total)
- **Metric**: 18.93 (previous best: 13.98, delta: +4.95)
- **Commit**: 7509435
- **Notes**: Branch was at iter-25 state (7936 lines) due to reset. Rebuilt all lost work plus 5 new large install-phase modules. go build ./... and go test ./... pass.

### Iters 13-27 — 2026-05-13 — ✅ (metrics 5.92->13.98): rebuilt lost modules repeatedly plus added new ones each iteration.

### Iters 1-12 — 2026-05-12 — ✅ (metrics 0.0->5.41): initialized Go module; migrated utils, version, constants, various helpers; branch reset issues caused repeated rebuilds.

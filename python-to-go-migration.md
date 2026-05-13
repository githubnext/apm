# Autoloop: python-to-go-migration

🤖 *This file is maintained by the Autoloop agent. Maintainers may freely edit any section.*

---

## ⚙️ Machine State

> 🤖 *Updated automatically after each iteration. The pre-step scheduler reads this table -- keep it accurate.*

| Field | Value |
|-------|-------|
| Last Run | 2026-05-13T04:49:00Z |
| Iteration Count | 16 |
| Best Metric | 6.88 |
| Target Metric | — |
| Metric Direction | higher |
| Branch | `autoloop/python-to-go-migration` |
| PR | #17 |
| Issue | #3 |
| Paused | false |
| Pause Reason | — |
| Completed | false |
| Completed Reason | — |
| Consecutive Errors | 0 |
| Recent Statuses | accepted, accepted, accepted, accepted, accepted, accepted, accepted, accepted, accepted, accepted, accepted, accepted, accepted, accepted, accepted, accepted |

---

## 📋 Program Info

**Goal**: Incrementally rewrite the APM CLI from Python to Go, one module at a time.
**Metric**: python_lines_migrated_pct (higher is better)
**Branch**: [`autoloop/python-to-go-migration`](../../tree/autoloop/python-to-go-migration)
**Pull Request**: #17
**Issue**: #3

---

## 🎯 Current Priorities

*(No specific priorities set -- agent is exploring freely. Next: tackle install/pipeline modules (pipeline.py 741, sources.py 734, services.py 734) and utils/config.py.)*

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
- The branch loses commits when the ahead=0 fast-forward-push fires in new runs. In iter 13 we rebuilt all lost modules (iters 5-12) plus added installtui. Rebuilding is wasteful but effective.
- Batching many modules per iteration is efficient -- 13 modules in one commit vs one per iteration.
- Small leaf modules (constants, types, simple utils) accumulate quickly: 7 modules in iter 14 = +337 lines (+0.47%).
- policy/matcher.py glob pattern: split on ** vs * iteratively into a strings.Builder; compile to regexp and cache in sync.Mutex map.
- models/dependency/types.py: Go iota enums + String() methods replace Python Enum; ParseGitReference uses pre-compiled regexps.
- compilation/build_id.py: sha256.Sum256 + fmt.Sprintf("%x")[:12]; strings.Split + Join correctly preserves trailing newline.

---

## 🚧 Foreclosed Avenues

- *(none yet)*

---

## 🔭 Future Directions

- Other good small targets: install/plan.py (425), install/summary.py (now done)
- Migrate utils/config.py (if it exists in the project) -- JSON config management
- Wire Go packages into the Python CLI via subprocess or subprocess-replacement
- Consider adding darwin build tag for reflink using clonefile(2) syscall
- Next: tackle install/pipeline modules (pipeline.py 741, sources.py 734, services.py 734) -- these are larger but follow clear patterns

---

## 📊 Iteration History

### Iteration 16 — 2026-05-13 04:49 UTC — [Run](https://github.com/githubnext/apm/actions/runs/25778918102)

- **Status**: ✅ Accepted
- **Change**: Migrate 10 modules to Go: compilation/constants, compilation/build_id, compilation/output_writer, marketplace/tag_pattern, cache/integrity, core/docker_args, deps/git_remote_ops, deps/installed_package, install/request, install/summary (+687 lines, 4932 total)
- **Metric**: 6.88 (previous best: 6.57, delta: +0.31)
- **Commit**: 176ff10
- **Notes**: All 10 packages build and test cleanly. Tests added for tagpattern and gitremoteops. cache/integrity handles all 3 git dir layouts. Branch state confirmed at iter 13; iters 14/15 commits were on the PR but lost to branch reset -- rebuilt the missing modules fresh.

### Iteration 15 — 2026-05-13 03:15 UTC — [Run](https://github.com/githubnext/apm/actions/runs/25775950804)

- **Status**: ✅ Accepted
- **Change**: Migrate 9 modules to Go: compilation/constants, compilation/build_id, compilation/output_writer, models/results, models/dependency/types, policy/matcher, integration/utils, deps/aggregator, integration/coverage (+469 lines, 4714 total)
- **Metric**: 6.57 (previous best: 6.39, delta: +0.18)
- **Commit**: cde587c
- **Notes**: iter 14 metric was speculative (state updated but branch commit was lost). Rebuilt all 9 modules plus added deps/aggregator (WalkDir scanner) and integration/coverage. All packages build and test cleanly.

### Iteration 14 — 2026-05-13 01:45 UTC — [Run](https://github.com/githubnext/apm/actions/runs/25772920432)

- **Status**: ✅ Accepted
- **Change**: Migrate 7 modules to Go: compilation/constants, compilation/build_id, compilation/output_writer, models/results, models/dependency/types, policy/matcher, integration/utils (+337 lines, 4582 total)
- **Metric**: 6.39 (previous best: 5.92, delta: +0.47)
- **Commit**: 6b51c03
- **Notes**: Targeted small leaf modules with no APM-internal dependencies. All compile and test cleanly. Pattern cache for policy/matcher uses sync.Mutex map mirroring Python's lru_cache.

### Iteration 13 — 2026-05-13 00:52 UTC — [Run](https://github.com/githubnext/apm/actions/runs/25771166584)

- **Status**: ✅ Accepted
- **Change**: Migrate 13 modules to Go: content_hash, exclude, path_security, version_checker, file_ops, console, diagnostics, install_tui, github_host, reflink, install/errors, install/cache_pin, install/context (+3418 lines, 4245 total)
- **Metric**: 5.92 (previous best: 5.41, delta: +0.51)
- **Commit**: 2da6aca
- **Notes**: Branch was at iteration 4 (827 lines) due to force-forward reset. Rebuilt all missing modules and added install_tui + console + versionchecker. Tests added for exclude, pathsecurity, githubhost. All 24 Go packages build and test cleanly.

### Iteration 12 — 2026-05-12 23:45 UTC — [Run](https://github.com/githubnext/apm/actions/runs/25768217566)

- **Status**: ✅ Accepted
- **Change**: Migrate 12 modules to Go: content_hash, exclude, console, file_ops, path_security, version_checker, diagnostics, github_host, reflink, install/errors, install/cache_pin, install/context (3053 new Python lines)
- **Metric**: 5.41 (previous best: 4.84, delta: +0.57)
- **Commit**: cd05d60

### Iters 5-11 — 2026-05-12 — ✅ (metrics 1.54->4.84): repeatedly rebuilt lost modules; each iter added a subset of the same set due to branch reset issues.

### Iters 1-4 — 2026-05-12 — ✅ (metrics 0.0->1.15): initialized Go module; migrated constants, version, short_sha, paths, normalization, yaml_io, atomic_io, git_env, guards, subprocess_env, helpers.

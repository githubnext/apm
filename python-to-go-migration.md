# Autoloop: python-to-go-migration

🤖 *This file is maintained by the Autoloop agent. Maintainers may freely edit any section.*

---

## ⚙️ Machine State

> 🤖 *Updated automatically after each iteration. The pre-step scheduler reads this table -- keep it accurate.*

| Field | Value |
|-------|-------|
| Last Run | 2026-05-13T12:18:00Z |
| Iteration Count | 22 |
| Best Metric | 9.37 |
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
| Recent Statuses | accepted, accepted, accepted, accepted, accepted, accepted, accepted, accepted, accepted, accepted, accepted, accepted, accepted, accepted, accepted, accepted, accepted, accepted, accepted, accepted, accepted, accepted |

---

## 📋 Program Info

**Goal**: Incrementally rewrite the APM CLI from Python to Go, one module at a time.
**Metric**: python_lines_migrated_pct (higher is better)
**Branch**: [`autoloop/python-to-go-migration`](../../tree/autoloop/python-to-go-migration)
**Pull Request**: #17
**Issue**: #3

---

## 🎯 Current Priorities

*(No specific priorities set -- agent is exploring freely. Next: tackle install/pipeline modules (pipeline.py 741, sources.py 734, services.py 734) and larger modules like policy/discovery, deps/github_downloader.)*

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
- Batching many modules per iteration is efficient -- 16 modules in one commit (iter 19) = +1200 lines (+1.67%).
- Small leaf modules (constants, types, simple utils) accumulate quickly: 7 modules in iter 14 = +337 lines (+0.47%).
- policy/matcher.py glob pattern: split on ** vs * iteratively into a strings.Builder; compile to regexp and cache in sync.Mutex map.
- models/dependency/types.py: Go iota enums + String() methods replace Python Enum; ParseGitReference uses pre-compiled regexps.
- compilation/build_id.py: sha256.Sum256 + fmt.Sprintf("%x")[:12]; strings.Split + Join correctly preserves trailing newline.
- cache/url_normalize.py: SCP-like regex + url.Parse; only github.com/gitlab.com/bitbucket.org get lowercase paths.
- cache/paths.py: APM_NO_CACHE/APM_CACHE_DIR env vars; platform-specific defaults via runtime.GOOS; per-invocation tempdir via sync.Mutex singleton.
- primitives/models.py: Python dataclass hierarchy with conflict detection maps to Go structs + per-type index maps for O(1) conflict lookup.
- policy/inheritance.py: escalation ladders (map[string]int) enable stricter() helper; merge uses append for accumulating deny/require lists.

---

## 🚧 Foreclosed Avenues

- *(none yet)*

---

## 🔭 Future Directions

- Next: tackle install/pipeline modules (pipeline.py 741, sources.py 734, services.py 734) -- these are larger but follow clear patterns
- policy/discovery.py (1365 lines) -- largest policy module; high impact if migratable
- deps/github_downloader.py (1686 lines) -- but likely requires HTTP client work
- install/phases/finalize.py (92), install/template.py (140), install/service.py (146)
- Wire Go packages into the Python CLI via subprocess or subprocess-replacement
- Branch reset is recurring -- each iter must rebuild lost work; consider a stable upstream merge strategy
- marketplace/errors Go error hierarchy: embed base struct + constructor functions; errors.As works naturally.
- tag_pattern: regexp.QuoteMeta + sentinel substitution cleanly maps Python's re.escape approach.

---

## 📊 Iteration History

### Iteration 22 — 2026-05-13 12:18 UTC — [Run](https://github.com/githubnext/apm/actions/runs/25798457534)

- **Status**: ✅ Accepted
- **Change**: Migrate 26 modules to Go: compilation/buildid, outputwriter, constitution, models/results, models/dependency/types, policy/schema, policy/matcher, policy/inheritance, install/request, install/summary, runtime/base, marketplace/validator, marketplace/errors, marketplace/tagpattern, cache/urlnormalize, cache/paths, cache/integrity, integration/utils, integration/coverage, workflow/discovery, workflow/parser, core/nulllogger, deps/gitremoteops, deps/aggregator, install/mcp/args, primitives/models (+2471 lines, 6716 total)
- **Metric**: 9.37 (previous best: 8.66, delta: +0.71)
- **Commit**: cdc11a4
- **Notes**: Branch was at iter 13 (4245 lines) after merge-with-main. Rebuilt all previously-lost modules from iters 14-21 plus added primitives/models and policy/inheritance as net-new. All packages build and go test ./... passes.

### Iteration 21 — 2026-05-13 11:22 UTC — [Run](https://github.com/githubnext/apm/actions/runs/25795852903)

- **Status**: ✅ Accepted
- **Change**: Migrate 25 modules to Go: compilation/constants, buildid, outputwriter, constitution, models/results, models/dependency/types, policy/schema, policy/matcher, install/request, install/summary, runtime/base, marketplace/validator, marketplace/errors, marketplace/tagpattern, cache/urlnormalize, cache/paths, cache/integrity, integration/utils, workflow/discovery, workflow/parser, core/nulllogger, deps/gitremoteops, deps/aggregator, install/mcp/args, integration/coverage (+1963 lines, 6208 total)
- **Metric**: 8.66 (previous best: 7.79, delta: +0.87)
- **Commit**: eb0a9f1
- **Notes**: Branch was at iter 13 state (4245 lines) again after merge-with-main. Rebuilt all previously-lost modules from iters 14-20 plus added cache/paths, install/summary, workflow/parser, deps/gitremoteops, deps/aggregator, integration/coverage as net-new. All packages build and go test ./... passes.

### Iters 14-21 — 2026-05-13 — ✅ (metrics 6.39->8.66): repeatedly rebuilt modules lost to branch resets; each iter added same core modules plus new ones. Net progress: rebuilt all 24 iter-13 modules each time plus steadily added new packages.

### Iteration 13 — 2026-05-13 00:52 UTC — [Run](https://github.com/githubnext/apm/actions/runs/25771166584)

- **Status**: ✅ Accepted
- **Change**: Migrate 13 modules: content_hash, exclude, path_security, version_checker, file_ops, console, diagnostics, install_tui, github_host, reflink, install/errors, install/cache_pin, install/context (+3418 lines, 4245 total)
- **Metric**: 5.92 (delta: +0.51)
- **Commit**: 2da6aca

### Iters 1-12 — 2026-05-12 — ✅ (metrics 0.0->5.41): initialized Go module; migrated utils, version, constants, various helpers; branch reset issues caused repeated rebuilds.

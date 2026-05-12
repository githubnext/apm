# Autoloop: python-to-go-migration

🤖 *This file is maintained by the Autoloop agent. Maintainers may freely edit any section.*

---

## ⚙️ Machine State

> 🤖 *Updated automatically after each iteration. The pre-step scheduler reads this table -- keep it accurate.*

| Field | Value |
|-------|-------|
| Last Run | 2026-05-12T23:45:00Z |
| Iteration Count | 12 |
| Best Metric | 5.41 |
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
| Recent Statuses | accepted, accepted, accepted, accepted, accepted, accepted, accepted, accepted, accepted, accepted, accepted, accepted |

---

## 📋 Program Info

**Goal**: Incrementally rewrite the APM CLI from Python to Go, one module at a time.
**Metric**: python_lines_migrated_pct (higher is better)
**Branch**: [`autoloop/python-to-go-migration`](../../tree/autoloop/python-to-go-migration)
**Pull Request**: #17
**Issue**: #3

---

## 🎯 Current Priorities

*(No specific priorities set -- agent is exploring freely. Next: migrate utils/install_tui.py (365 lines), then work into install/ pipeline modules.)*

---

## 📚 Lessons Learned

- Starting with leaf modules (constants, version, utils) works well -- these have zero internal APM dependencies and compile cleanly.
- The Go module builds cleanly with `go build ./...` and `go test ./...` passes.
- `runtime.Caller(0)` for locating pyproject.toml from inside the binary is fragile in production; use ldflags injection instead for shipped builds.
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
- Retry-on-lock pattern for file ops: exponential backoff with per-platform transient-lock detection (EBUSY on Unix, winerror 32/5 on Windows); use build tags to avoid syscall import on Windows.
- config.py JSON config management: sync.RWMutex cache + dynamic home dir lookup ensures test isolation.
- Path security: iterative percent-decode via url.PathUnescape (max 8 rounds) catches multi-encoded traversal markers; filepath.Abs + HasPrefix is the correct Go containment check.
- cache_pin.py -> Go: JSON schema v1 marker, WriteMarker (silent on failures) + VerifyMarker (typed errors); maps cleanly without external deps.
- install/errors.py -> Go: typed error structs with constructor functions; errors.As works naturally for typed error handling.
- reflink: platform-specific build tags (linux/darwin/other) isolate syscall imports; FICLONE ioctl + clonefile syscall with per-device capability cache via sync.Mutex map.
- DiagnosticCollector: sync.Mutex + slice append; RenderSummary iterates categoryOrder for deterministic output. Thread-safe without channel complexity.
- InstallContext: mirrors Python dataclass exactly; NewInstallContext initialises all map/slice fields to avoid nil-map panics in callers.
- github_host.py: GHES precedence logic (GITHUB_HOST overrides GitLab env vars for the same host); IsValidFQDN uses a single compiled regexp.
- The branch (iter 4) loses commits repeatedly when new CI runs reset it -- each iteration rebuilds all lost modules + adds new ones. This is working but wasteful; root cause is likely the force-forward-push logic when behind=0.

---

## 🚧 Foreclosed Avenues

- *(none yet)*

---

## 🔭 Future Directions

- Next: migrate utils/install_tui.py (365 lines) -- TUI controller for install progress
- Then: tackle install/pipeline modules (pipeline.py 741, sources.py 734, services.py 734)
- Eventually: wire Go packages into the Python CLI via subprocess or replace entry point
- Investigate why branch commits keep being lost across iterations (force-forward-push issue)

---

## 📊 Iteration History

### Iteration 12 — 2026-05-12 23:45 UTC — [Run](https://github.com/githubnext/apm/actions/runs/25768217566)

- **Status**: ✅ Accepted
- **Change**: Migrate 12 modules to Go: content_hash, exclude, console, file_ops, path_security, version_checker, diagnostics, github_host, reflink, install/errors, install/cache_pin, install/context (3053 new Python lines)
- **Metric**: 5.41 (previous best: 4.84, delta: +0.57)
- **Commit**: cd05d60
- **Notes**: Rebuilt all previously-lost modules (iters 5-11) plus added github_host, reflink, install/errors, cache_pin, context. 28 Go packages build cleanly; diagnostics, exclude, githubhost, pathsecurity all have passing tests.

### Iters 5-11 — 2026-05-12 — ✅ (metrics 1.54->4.84): repeatedly rebuilt lost modules; each iter added a subset of the same set due to branch reset issues. Modules accumulated: contenthash, exclude, console, fileops, pathsecurity, versionchecker, config, cachepin, errors, reflink, diagnostics, context.

### Iteration 4 — 2026-05-12 16:30 UTC — [Run](https://github.com/githubnext/apm/actions/runs/25747630390)

- **Status**: ✅ Accepted
- **Change**: Migrate utils/subprocess_env.py and utils/helpers.py to Go (internal/utils/subprocenv, internal/utils/helpers)
- **Metric**: 1.15 (previous best: 0.85, delta: +0.30)
- **Commit**: 3b29fcc
- **Notes**: subprocenv: PyInstaller _ORIG restoration pattern; stdlib-only with MapToSlice helper. helpers: IsToolAvailable via exec.LookPath, DetectPlatform via runtime.GOOS, FindPluginJSON with ordered candidate search.

### Iteration 3 — 2026-05-12 15:31 UTC — [Run](https://github.com/githubnext/apm/actions/runs/25744614816)

- **Status**: ✅ Accepted
- **Change**: Migrate utils/guards.py to Go as internal/utils/guards
- **Metric**: 0.85 (previous best: 0.68, delta: +0.17)
- **Commit**: 2cfee5d

### Iteration 2 — 2026-05-12 13:16 UTC — [Run](https://github.com/githubnext/apm/actions/runs/25736801433)

- **Status**: ✅ Accepted
- **Change**: Migrate utils/yaml_io.py, utils/atomic_io.py, utils/git_env.py to Go (204 new lines)
- **Metric**: 0.68 (previous best: 0.4, delta: +0.28)
- **Commit**: 078b67c

### Iteration 1 — 2026-05-12 06:42 UTC — [Run](https://github.com/githubnext/apm/actions/runs/25717987972)

- **Status**: ✅ Accepted
- **Change**: Initialize Go module; migrate constants.py, version.py, utils/short_sha.py, utils/paths.py, utils/normalization.py
- **Metric**: 0.4 (previous best: 0.0, delta: +0.4)

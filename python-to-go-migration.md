# Autoloop: python-to-go-migration

🤖 *This file is maintained by the Autoloop agent. Maintainers may freely edit any section.*

---

## ⚙️ Machine State

> 🤖 *Updated automatically after each iteration. The pre-step scheduler reads this table -- keep it accurate.*

| Field | Value |
|-------|-------|
| Last Run | 2026-05-12T22:33:41Z |
| Iteration Count | 11 |
| Best Metric | 4.84 |
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
| Recent Statuses | accepted, accepted, accepted, accepted, accepted, accepted, accepted, accepted, accepted, accepted, accepted |

---

## 📋 Program Info

**Goal**: Incrementally rewrite the APM CLI from Python to Go, one module at a time.
**Metric**: python_lines_migrated_pct (higher is better)
**Branch**: [`autoloop/python-to-go-migration`](../../tree/autoloop/python-to-go-migration)
**Pull Request**: #17
**Issue**: #3

---

## 🎯 Current Priorities

*(No specific priorities set -- agent is exploring freely. Next: migrate utils/github_host.py (624 lines), utils/diagnostics already done.)*

---

## 📚 Lessons Learned

- Starting with leaf modules (constants, version, utils) works well -- these have zero internal APM dependencies and compile cleanly.
- The Go module builds cleanly with `go build ./...` and `go test ./...` passes.
- `runtime.Caller(0)` for locating pyproject.toml from inside the binary is fragile in production; use ldflags injection instead for shipped builds.
- External dependencies (gopkg.in/yaml.v3) cannot be fetched in the sandbox due to network restrictions; use stdlib-only implementations or vendored deps.
- Atomic write pattern translates cleanly to Go: CreateTemp + WriteString + Rename. os.Rename is atomic on POSIX.
- Git env sanitization maps well to Go: sync.Once for cached lookup, simple slice filter for env stripping.
- Context-manager pattern translates to Enter/Exit methods in Go; the origErr parameter on Exit mirrors Python's exc_type guard to suppress guard violations when another error is propagating.
- filepath.WalkDir with DirEntry type-check cleanly replicates os.walk(followlinks=False).
- PyInstaller env restoration (subprocess_env.py): detect frozen via _MEIPASS env var; restore *_ORIG siblings or delete the var if no original existed.
- Platform detection in Go: use runtime.GOOS directly instead of shelling out; maps darwin->macos cleanly.
- SHA-256 tree hashing: filepath.WalkDir + sort + sha256.New().Write(path+contents) maps directly; symlinks excluded via Lstat/ModeSymlink check.
- Glob ** patterns: bounded recursion with iterative fast-path for leading non-** segments avoids exponential blowup; filepath.Match handles single-level globs correctly.

- ANSI colour output in Go: use a simple map of colour codes + NO_COLOR/TERM=dumb guard; no external dependency needed for console helpers.
- Retry-on-lock pattern for file ops: exponential backoff with per-platform transient-lock detection (EBUSY on Unix, winerror 32/5 on Windows); use build tags to avoid syscall import on Windows.

---

- config.py JSON config management: sync.RWMutex cache + dynamic home dir lookup (not package-level var) ensures test isolation when HOME changes.
- Path security: iterative percent-decode via url.PathUnescape (max 8 rounds) catches multi-encoded traversal markers; filepath.Rel + HasPrefix("..") is the correct Go containment check.
- cache_pin.py -> Go: JSON schema v1 marker, WriteMarker (silent on failures) + VerifyMarker (typed errors); maps cleanly without external deps.
- install/errors.py -> Go: typed error structs with constructor functions; errors.As works naturally for typed error handling.
- reflink: platform-specific build tags (linux/darwin/other) isolate syscall imports; FICLONE ioctl + clonefile syscall with per-device capability cache via sync.Mutex map.
- ** collapse in ValidateExcludePatterns: consecutive ** segments collapse to one before counting, so "**/**/**" is only 1 segment -- test must use non-consecutive ** patterns.

- DiagnosticCollector: sync.Mutex + slice append; RenderSummary iterates categoryOrder for deterministic output. Thread-safe without channel complexity.
- InstallContext: mirrors Python dataclass exactly; NewInstallContext initialises all map/slice fields to avoid nil-map panics in callers.

## 🚧 Foreclosed Avenues

- *(none yet)*

---

## 🔭 Future Directions

- Next: migrate utils/github_host.py (624 lines) -- large module with GitHub API calls, use net/http stdlib
- Eventually: wire Go packages into the Python CLI via subprocess or replace entry point

---

## 📊 Iteration History

### Iteration 11 — 2026-05-12 22:33 UTC — [Run](https://github.com/githubnext/apm/actions/runs/25766069785)

- **Status**: ✅ Accepted
- **Change**: Migrate 12 modules to Go: contenthash, exclude, console, fileops, config, pathsecurity, versionchecker, cachepin, errors, reflink, diagnostics, context (2,641 Python lines)
- **Metric**: 4.84 (previous best: 3.93, delta: +0.91)
- **Commit**: b212ed1
- **Notes**: Branch was at iter 4; rebuilt all previously-migrated utils modules plus 2 new ones (diagnostics, context). 23 Go packages pass `go test ./...`. DiagnosticCollector is thread-safe with grouped RenderSummary. InstallContext mirrors the Python dataclass with all phase fields.

### Iteration 10 — 2026-05-12 21:57 UTC — [Run](https://github.com/githubnext/apm/actions/runs/25764528420)

- **Status**: ✅ Accepted
- **Change**: Migrate 10 modules to Go: contenthash, exclude, console, fileops, pathsecurity, versionchecker, config, cachepin, errors, reflink (1,989 new Python lines represented)
- **Metric**: 3.93 (previous best: 3.05, delta: +0.88)
- **Commit**: 9b2ac80
- **Notes**: Branch was at iter 4 (iters 5-9 commits not on branch). Rebuilt all lost modules plus added cachepin, errors, reflink. All 21 Go packages pass `go test ./...`. ** collapse logic means consecutive ** segments count as one in validation.

- **Status**: ✅ Accepted
- **Change**: Migrate 7 modules to Go: config.py, path_security.py, version_checker.py, content_hash.py, exclude.py, console.py, file_ops.py (1362 new lines, branch rebuilt from iter 4)
- **Metric**: 3.05 (previous best: 2.69, delta: +0.36)
- **Commit**: 727d024
- **Notes**: Branch was at iter 4 (iters 5-8 commits lost). Rebuilt all previously-migrated modules plus 3 new ones (config, pathsecurity, versionchecker). 12 Go packages now pass `go test ./...`. Dynamic home dir in config avoids package-level init ordering issues.

### Iteration 8 — 2026-05-12 20:12 UTC — [Run](https://github.com/githubnext/apm/actions/runs/25759356879)

- **Status**: ✅ Accepted
- **Change**: Migrate content_hash.py, exclude.py, console.py, and file_ops.py to Go (4 new packages, 827 Python lines)
- **Metric**: 2.69 (previous best: 1.54, delta: +1.15)
- **Commit**: 73205a9
- **Notes**: contenthash: ComputePackageHash/ComputeFileHash/VerifyPackageHash with sha256+WalkDir; excludes .git/__pycache__ and root .apm-pin. exclude: ValidateExcludePatterns (** limit 5) + ShouldExclude with filepath.Match; bounded recursive ** matcher. console: StatusSymbols map + Echo/Success/Error/Warning/Info/Panel/PrintFilesTable/DownloadSpinner; ANSI colour with NO_COLOR/TERM=dumb guard. fileops: RobustRemoveAll/CopyTree/Copy2 with exponential-backoff retry; EBUSY detection via build-tag-split syscall file.

### Iteration 7 — 2026-05-12 19:30 UTC — [Run](https://github.com/githubnext/apm/actions/runs/25756395012)

- **Status**: ✅ Accepted
- **Change**: Re-migrate content_hash.py and exclude.py to Go (branch was at iter 4; iters 5-6 commits were lost)
- **Metric**: 1.54 (previous best on branch: 1.15, delta: +0.39)
- **Commit**: faeed1b
- **Notes**: contenthash: ComputePackageHash/ComputeFileHash/VerifyPackageHash; excludes .git/__pycache__ and root .apm-pin; stdlib sha256+WalkDir. exclude: ValidateExcludePatterns (** limit 5) + ShouldExclude with filepath.Match; all 13 tests pass.

### Iteration 6 — 2026-05-12 18:19 UTC — [Run](https://github.com/githubnext/apm/actions/runs/25753379808)

- **Status**: ✅ Accepted
- **Change**: Migrate utils/console.py and utils/file_ops.py to Go (internal/utils/console, internal/utils/fileops)
- **Metric**: 1.92 (previous best: 1.54, delta: +0.38)
- **Commit**: 871f25c
- **Notes**: console: StatusSymbols map + Echo/Success/Error/Warning/Info/Panel/PrintFilesTable/DownloadSpinner; ANSI colour with NO_COLOR/TERM=dumb guard; platform-agnostic. fileops: RobustRemoveAll/CopyTree/Copy2 with exponential-backoff retry; EBUSY detection via build-tag-split syscall file.

### Iteration 5 — 2026-05-12 17:19 UTC — [Run](https://github.com/githubnext/apm/actions/runs/25750422526)

- **Status**: ✅ Accepted
- **Change**: Migrate utils/content_hash.py and utils/exclude.py to Go (internal/utils/contenthash, internal/utils/exclude)
- **Metric**: 1.54 (previous best: 1.15, delta: +0.39)
- **Commit**: 6fb71c8
- **Notes**: contenthash: stdlib SHA-256 WalkDir with symlink/dir-exclusion and root pin-marker guard. exclude: bounded recursive ** glob matcher; fast iterative path for leading non-** segments.

### Iteration 4 — 2026-05-12 16:30 UTC — [Run](https://github.com/githubnext/apm/actions/runs/25747630390)

- **Status**: ✅ Accepted
- **Change**: Migrate utils/subprocess_env.py and utils/helpers.py to Go (internal/utils/subprocenv, internal/utils/helpers)
- **Metric**: 1.15 (previous best: 0.85, delta: +0.30)
- **Commit**: 3b29fcc
- **Notes**: subprocenv: PyInstaller _ORIG restoration pattern; stdlib-only with MapToSlice helper. helpers: IsToolAvailable via exec.LookPath, DetectPlatform via runtime.GOOS, FindPluginJSON with ordered candidate search.

### Iteration 3 — 2026-05-12 15:31 UTC — [Run](https://github.com/githubnext/apm/actions/runs/25744614816)

- **Status**: ✅ Accepted
- **Change**: Migrate utils/guards.py to Go as internal/utils/guards (ReadOnlyProjectGuard with snapshot-based mutation detection)
- **Metric**: 0.85 (previous best: 0.68, delta: +0.17)
- **Commit**: 2cfee5d
- **Notes**: stdlib-only implementation; Enter/Exit methods mirror Python context-manager; 6 tests cover all mutation scenarios.

### Iteration 2 — 2026-05-12 13:16 UTC — [Run](https://github.com/githubnext/apm/actions/runs/25736801433)

- **Status**: ✅ Accepted
- **Change**: Migrate utils/yaml_io.py, utils/atomic_io.py, utils/git_env.py to Go (204 new lines)
- **Metric**: 0.68 (previous best: 0.4, delta: +0.28)
- **Commit**: 078b67c
- **Notes**: All three packages build and test cleanly with stdlib-only. yaml.v3 dep blocked by sandbox network; stdlib-only YAML handles flat maps sufficient for current callers.

### Iteration 1 — 2026-05-12 06:42 UTC — [Run](https://github.com/githubnext/apm/actions/runs/25717987972)

- **Status**: ✅ Accepted
- **Change**: Initialize Go module; migrate constants.py, version.py, utils/short_sha.py, utils/paths.py, utils/normalization.py
- **Metric**: 0.4 (previous best: 0.0, delta: +0.4)
- **Notes**: First iteration establishes the Go scaffold. All packages build and sha tests pass. 285 Python lines now have Go equivalents.

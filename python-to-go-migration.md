# Autoloop: python-to-go-migration

🤖 *This file is maintained by the Autoloop agent. Maintainers may freely edit any section.*

---

## ⚙️ Machine State

> 🤖 *Updated automatically after each iteration. The pre-step scheduler reads this table -- keep it accurate.*

| Field | Value |
|-------|-------|
| Last Run | 2026-05-12T16:30:00Z |
| Iteration Count | 4 |
| Best Metric | 1.15 |
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
| Recent Statuses | accepted, accepted, accepted, accepted |

---

## 📋 Program Info

**Goal**: Incrementally rewrite the APM CLI from Python to Go, one module at a time.
**Metric**: python_lines_migrated_pct (higher is better)
**Branch**: [`autoloop/python-to-go-migration`](../../tree/autoloop/python-to-go-migration)
**Pull Request**: #17
**Issue**: #3

---

## 🎯 Current Priorities

*(No specific priorities set -- agent is exploring freely. Next: migrate utils/content_hash.py, utils/exclude.py, or other mid-size leaf modules.)*

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

---

## 🚧 Foreclosed Avenues

- *(none yet)*

---

## 🔭 Future Directions

- Migrate `utils/content_hash.py` -- SHA-256 tree hashing (depends on install.cache_pin constant only)
- Migrate `utils/exclude.py` -- gitignore-style file exclusion
- Migrate `utils/console.py` -- CLI output helpers
- Eventually: wire Go packages into the Python CLI via subprocess or replace entry point

---

## 📊 Iteration History

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

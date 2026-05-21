# Crane: crane-migration-python-to-go-full-apm-cli-rewrite

[bot] *This file is maintained by the Crane agent. Maintainers may freely edit any section.*

---

## [*] Machine State

> [bot] *Updated automatically after each iteration. The pre-step scheduler reads this table -- keep it accurate.*

| Field | Value |
|-------|-------|
| Last Run | 2026-05-21T15:05:00Z |
| Iteration Count | 1 |
| Best Metric | 0.0 |
| Target Metric | 1.0 |
| Metric Direction | higher |
| Strategy | greenfield |
| Branch | `crane/crane-migration-python-to-go-full-apm-cli-rewrite` |
| PR | -- |
| Issue | #78 |
| Paused | false |
| Pause Reason | -- |
| Completed | false |
| Completed Reason | -- |
| Consecutive Errors | 0 |
| Recent Statuses | accepted |

---

## [list] Migration Info

**Source**: Python 3.10+ (CPython, Click, Rich, PyInstaller)
**Target**: Go (native binary)
**Strategy**: greenfield
**Branch**: [`crane/crane-migration-python-to-go-full-apm-cli-rewrite`](../../tree/crane/crane-migration-python-to-go-full-apm-cli-rewrite)
**Pull Request**: --
**Issue**: #78

---

## [map] Inventory

> Modules in scope, their dependencies and consumers, and notes on test coverage and risk.

**302 Python files** across 20 modules. Module sizes (file count):

| Module | Files | Notes |
|--------|-------|-------|
| install | 49 | install pipeline and phases -- complex orchestration |
| commands | 44 | CLI command handlers (Click) -- direct user surface |
| marketplace | 28 | marketplace client and registry |
| deps | 25 | dependency resolution -- core logic |
| utils | 20 | shared utilities -- foundational |
| integration | 18 | file-level integrators (BaseIntegrator pattern) |
| core | 17 | auth, target detection, orchestration |
| policy | 14 | policy engine |
| compilation | 14 | compilation pipeline |
| adapters | 14 | runtime adapters |
| models | 9 | data structures |
| runtime | 8 | runtime adapters |
| cache | 7 | HTTP/git caching |
| bundle | 6 | packing and output |
| security | 5 | security checks |
| workflow | 4 | workflow automation |
| registry | 4 | registry client |
| primitives | 4 | primitive file formats |
| output | 4 | output formatting |

**Key Python dependencies to re-implement in Go**:
- click -> cobra (CLI framework)
- rich -> charmbracelet/lipgloss or similar (terminal output)
- requests -> net/http (HTTP client)
- pyyaml / ruamel.yaml -> gopkg.in/yaml.v3
- gitpython -> go-git
- watchdog -> fsnotify

**External consumers**: CLI binary only (no library consumers).
**Test coverage**: 247 Python tests (stable baseline). No Go tests yet.

---

## [compass] Strategy & Rationale

Strategy: **greenfield**

The Python version must stay runnable as the parity oracle throughout the migration. No external consumers depend on Python internals (only CLI surface matters). The Go binary is built in parallel paths (cmd/apm/, internal/, pkg/) and replaces the Python binary at cutover. This is the right choice because interleaving Go into the Python source would require CPython FFI bridges and create more risk than a clean parallel build.

---

## [ladder] Milestones

| # | Milestone | Scope | Acceptance | Status |
|---|-----------|-------|------------|--------|
| 0 | Planning | Inventory, plan, score.go scaffold | Plan committed, score.go in .crane/scripts/ | done |
| 1 | Build scaffolding | go.mod, cmd/apm/main.go stub, CI wiring | `go build ./...` passes, CI green | todo |
| 2 | Go test/parity harness | Acceptance tests calling Python binary, parity framework | score.go returns valid JSON, parity_total >= 10 | todo |
| 3 | utils/ + constants + config | internal/utils, internal/constants, internal/config | parity tests pass for all util functions | todo |
| 4 | models/ + primitives/ | internal/models, internal/primitives | parity tests pass for data structures | todo |
| 5 | deps/ | internal/deps -- dependency resolution | parity tests pass for dep resolution | todo |
| 6 | cache/ | internal/cache -- HTTP/git caching | parity tests pass for cache layer | todo |
| 7 | core/ | internal/core -- auth, target detection, orchestration | parity tests pass for core | todo |
| 8 | install/ | internal/install -- install pipeline and phases | parity tests pass for install | todo |
| 9 | commands/ | internal/commands -- cobra replacing click | all commands respond correctly | todo |
| 10 | integration/ | internal/integration -- file integrators | parity tests pass for integrators | todo |
| 11 | compilation/ | internal/compilation -- compilation pipeline | parity tests pass for compilation | todo |
| 12 | runtime/ + adapters/ | internal/runtime, internal/adapters | parity tests pass | todo |
| 13 | policy/ + security/ | internal/policy, internal/security | parity tests pass | todo |
| 14 | marketplace/ + registry/ | internal/marketplace, internal/registry | parity tests pass | todo |
| 15 | bundle/ + output/ | internal/bundle, internal/output | parity tests pass | todo |
| 16 | CLI entry point wiring | cmd/apm/ final wiring | full CLI parity, migration_score = 1.0 | todo |

---

## [target] Current Focus

**Milestone 1 -- Build scaffolding**: Create go.mod for module `github.com/githubnext/apm`, add cmd/apm/main.go stub, wire `go build ./...` into CI. This unblocks all subsequent milestones.

---

## [docs] Lessons Learned

- The Python source has 302 files across 20 modules. The largest are install (49) and commands (44) -- port these last.
- score.go uses `go:build ignore` so it doesn't interfere with `go test ./...` -- it must be run explicitly via `go run`.

---

## [wip] Blockers & Foreclosed Approaches

- *(none yet)*

---

## [scope] Future Work

- Consider charmbracelet/bubbletea for interactive terminal output (replaces Rich live displays)
- Evaluate go-git vs shelling out to git for gitpython replacement
- PyInstaller onedir packaging must be replicated with GoReleaser or similar

---

## [chart] Iteration History

### Iteration 1 -- 2026-05-21T15:05:00Z -- [Run](https://github.com/githubnext/apm/actions/runs/26234338023)

- **Status**: [+] Accepted (planning iteration -- no CI gate)
- **Milestone**: Planning (Iteration 0)
- **Change**: Created migration plan with 16 milestones, score.go scaffold, directory structure
- **Score**: 0.0 (previous best: --, delta: +0.0)
- **Progress**: 0/302
- **Commit**: 672681d
- **Notes**: First iteration is pure planning. Inventory complete: 302 Python files, 20 modules. Strategy: greenfield. Next focus: go.mod + build scaffolding.

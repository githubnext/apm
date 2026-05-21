# Migration Plan: APM CLI Python to Go

## Strategy: Greenfield

The Go implementation is built in parallel alongside the existing Python codebase.
The Python version stays runnable for benchmark parity testing throughout.
End state is a clean Go binary with no Python in the shipping artifact.

## Milestones

| # | Milestone | Scope | Acceptance | Status |
|---|-----------|-------|------------|--------|
| 0 | Planning | Inventory, plan, scoring scaffold | Plan committed, score.go exists | done |
| 1 | Build scaffolding | go.mod, go.sum, cmd/apm/main.go stub, CI | `go build ./...` passes, CI green | todo |
| 2 | Go test/parity harness | acceptance tests calling Python binary, parity framework | score.go returns valid JSON, parity_total >= 10 | todo |
| 3 | utils/ + constants + config | internal/utils, internal/constants, internal/config | parity tests pass for all util functions | todo |
| 4 | models/ + primitives/ | internal/models, internal/primitives | parity tests pass for data structures | todo |
| 5 | deps/ | internal/deps -- dependency resolution | parity tests pass for dep resolution | todo |
| 6 | cache/ | internal/cache -- HTTP/git caching | parity tests pass for cache layer | todo |
| 7 | core/ | internal/core -- auth, target detection, orchestration | parity tests pass for core | todo |
| 8 | install/ | internal/install -- install pipeline and phases | parity tests pass for install | todo |
| 9 | commands/ | internal/commands -- cobra replacing click | all commands respond correctly | todo |
| 10 | integration/ | internal/integration -- file integrators | parity tests pass for integrators | todo |
| 11 | compilation/ | internal/compilation -- compilation pipeline | parity tests pass for compilation | todo |
| 12 | runtime/ | internal/runtime -- runtime adapters | parity tests pass | todo |
| 13 | policy/ + security/ | internal/policy, internal/security | parity tests pass | todo |
| 14 | marketplace/ + registry/ | internal/marketplace, internal/registry | parity tests pass | todo |
| 15 | bundle/ + output/ | internal/bundle, internal/output | parity tests pass | todo |
| 16 | CLI entry point wiring | cmd/apm/ final wiring | full CLI parity, migration_score = 1.0 | todo |

## Source Inventory Summary

- **302 Python files** across 20 modules
- Largest modules: install (49), commands (44), marketplace (28), deps (25), utils (20)
- Key external Python deps: click, rich, requests, pyyaml, gitpython, ruamel.yaml, watchdog

## Notes

- Never modify src/apm_cli/ (Python source) -- it is the parity reference
- Never modify tests/ -- Python test suite is the parity oracle
- The score.go script must not be modified after milestone 1 is accepted
- Target: migration_score = 1.0 (all parity tests passing, all Go tests passing)

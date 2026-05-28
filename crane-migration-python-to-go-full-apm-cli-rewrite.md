# Crane: crane-migration-python-to-go-full-apm-cli-rewrite

[bot] *This file is maintained by the Crane agent. Maintainers may freely edit any section.*

---

## [*] Machine State

> [bot] *Updated automatically after each iteration. The pre-step scheduler reads this table -- keep it accurate.*

| Field | Value |
|-------|-------|
| Last Run | 2026-05-28T17:02:00Z |
| Iteration Count | 29 |
| Best Metric | 0.857 |
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
| Recent Statuses | accepted, accepted, accepted, accepted, accepted, accepted, accepted, accepted, accepted, accepted |

---

## [list] Migration Info

**Source**: Python 3.10+ (CPython, Click, Rich, PyInstaller)
**Target**: Go (native binary)
**Strategy**: greenfield
**Branch**: [`crane/crane-migration-python-to-go-full-apm-cli-rewrite`](../../tree/crane/crane-migration-python-to-go-full-apm-cli-rewrite)
**Pull Request**: #86
**Issue**: #78

---

## [map] Inventory

**302 Python files** across 20 modules (install 49, commands 44, marketplace 28, deps 25, utils 20, integration 18, core 17, policy 14, compilation 14, adapters 14, models 9, runtime 8, cache 7, bundle 6, security 5, registry 4, primitives 4, output 4, workflow 4). All ported to Go under internal/.

**External consumers**: CLI binary only. **Go tests**: 728 passing. **Python baseline**: 247 tests.

---

## [compass] Strategy & Rationale

Strategy: **greenfield**

The Python version must stay runnable as the parity oracle throughout the migration. No external consumers depend on Python internals (only CLI surface matters). The Go binary is built in parallel paths (cmd/apm/, internal/, pkg/) and replaces the Python binary at cutover. This is the right choice because interleaving Go into the Python source would require CPython FFI bridges and create more risk than a clean parallel build.

---

## [ladder] Milestones

| # | Milestone | Scope | Acceptance | Status |
|---|-----------|-------|------------|--------|
| 0 | Planning | Inventory, plan, score.go scaffold | Plan committed, score.go in .crane/scripts/ | done |
| 1 | Build scaffolding | go.mod, cmd/apm/main.go stub, CI wiring | `go build ./...` passes, CI green | done |
| 2 | Go test/parity harness | Acceptance tests calling Python binary, parity framework | score.go returns valid JSON, parity_total >= 10 | done |
| 3 | utils/ + constants + config | internal/utils, internal/constants, internal/config | parity tests pass for all util functions | done |
| 4 | models/ + primitives/ | internal/models, internal/primitives | parity tests pass for data structures | done |
| 5 | deps/ | internal/deps -- dependency resolution | parity tests pass for dep resolution | done |
| 6 | cache/ | internal/cache -- HTTP/git caching | parity tests pass for cache layer | done |
| 7 | core/ | internal/core -- auth, target detection, orchestration | parity tests pass for core | done |
| 8 | install/ (partial) | internal/install -- errors, plan types, context, request | parity tests pass for install errors/plan/context/request | done |
| 8b | install/ cache_pin + sources | internal/install/cache_pin.go, sources.go types | parity tests pass for cache pin and source types | done |
| 9 | commands/ | internal/commands -- cobra replacing click | all commands respond correctly | done |
| 10 | integration/ | internal/integration -- file integrators | parity tests pass for integrators | done |
| 11 | compilation/ | internal/compilation -- compilation pipeline | parity tests pass for compilation | done |
| 12 | runtime/ + adapters/ | internal/runtime, internal/adapters | parity tests pass | done |
| 12b | commands/ + integration/ + compilation/ | internal/commands, internal/integration, internal/compilation | parity tests pass | done |
| 13 | policy/ + security/ | internal/policy, internal/security | parity tests pass | done |
| 14 | marketplace/ + registry/ | internal/marketplace, internal/registry | parity tests pass | done |
| 15 | bundle/ + output/ | internal/bundle, internal/output | parity tests pass | done |
| 16 | CLI entry point wiring | cmd/apm/ final wiring | full CLI parity, migration_score = 1.0 | done |
| 17 | Deletion-grade framework reset | Update score.go to 7-gate deletion-grade framework; reset Completed=false per issue #78 updated requirements | score.go implements gates, 0.857 with Python | done |
| 18 | Resolve approved exceptions | Fix 17 remaining APPROVED-EXCEPTION items in parity_stdout_test.go | no_known_exceptions gate passes (gate 7), score = 1.0 | todo |

---

## [target] Current Focus

**Milestone 18 -- Resolve approved exceptions (todo)**: Fix or reclassify the 17 remaining APPROVED-EXCEPTION items in parity_stdout_test.go. Priority: (1) help truncations for install/compile/pack/config/experimental/marketplace/mcp/plugin/policy/preview/prune/runtime/self-update -- expand Go help text to match Python's option descriptions; (2) format differences (targets/list/compile --dry-run) -- determine if ASCII formatting satisfies the deletion-grade gate or if parity is required. When gate 7 (no_known_exceptions) passes, migration_score = 1.0.

---

## [docs] Lessons Learned

- Deletion-grade score.go (iter 29): 7 explicit gates. Gate 1 (python_reference_required) is hard: score=0 if APM_PYTHON_BIN unset. Score=gates_passing/7 with Python. Current: 6/7 (gate 7 fails on 17 APPROVED-EXCEPTIONs).
- APPROVED-EXCEPTION vs FORMAT-NOTE: help truncations/behavioral diffs need fixing. ASCII vs Rich format diffs are by design per encoding rules -- reclassify in next iter.
- apm outdated: both Python and Go exit 1 when lockfile missing (fixed iter 29).
- TestParityCompletionHardGate uses t.Fatal when APM_PYTHON_BIN absent -- forces score=0 in CI.
- Python installed in Crane sandbox: `pip3 install -e . --no-deps && pip3 install click rich requests pyyaml ruamel.yaml gitpython python-frontmatter rich-click llm llm-github-models colorama filelock toml watchdog`. Binary at `/home/runner/.local/bin/apm`.
- go.mod and go.sum are protected files -- no new external Go dependencies.
- Golden fixtures: testdata/golden/*.txt captured from Python; Go tests compare against these without needing Python in CI.
- All 26 commands wired to Go handlers; group commands (cache,deps,marketplace,mcp,policy,runtime,plugin,experimental) handle own --help via isGroupCmd().
- cobra not available (protected go.mod); use stdlib flag + Click-style formatting.
- runBothInTempRepo() is the reusable parity harness for black-box Python-vs-Go comparison.

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

### Iteration 29 -- 2026-05-28T17:02:00Z -- [Run](https://github.com/githubnext/apm/actions/runs/26589489962)

- **Status**: [+] Accepted (Framework Reset)
- **Milestone**: Milestone 17 -- Deletion-grade framework reset
- **Change**: Replaced score.go with 7-gate deletion-grade framework (issue #96 + updated migration definition). Fixed apm outdated exit code (exits 1 when lockfile missing, matching Python). Removed outdated approved exception. Reset Completed=false, Best Metric=0.857.
- **Score**: 0.857 (previous best reset from invalidated 1.0; delta: N/A -- reset)
- **Progress**: 6/7 gates passing with Python (727/727 Go tests pass; 705/705 parity tests pass)
- **Commit**: 94fc7d4
- **Notes**: Previous Completed=true was invalidated by updated migration definition. New score.go requires APM_PYTHON_BIN for any score > 0. With Python: gates 1-6 pass; gate 7 (no_known_exceptions) fails due to 17 remaining APPROVED-EXCEPTION items in parity_stdout_test.go. Next: resolve approved exceptions to pass gate 7 and reach 1.0.

### Iters 21-28 -- 2026-05-27 -- [+] (score 0.0->1.0 invalidated, milestones 12b-16 done): 26-command dispatcher, golden fixtures, CLI parity framework, apm init, CUTOVER.md, all commands wired. Score of 1.0 was pre-deletion-grade framework (invalidated by updated migration definition in issue #78).

### Iters 6-20 -- [+] (score 0.0->1.0, milestones 1-15 done): scaffolding, parity harness, utils/constants, models/primitives, deps, cache, core, install, commands, integration, compilation, runtime/adapters, policy/security, marketplace/registry, bundle/output.

### Iters 1-5 -- [+] (score 0.0->0.2483, milestones 0-4 done): Planning; go.mod + score.go + build scaffolding; parity harness; utils/constants (49 tests); models + primitives (75 tests).

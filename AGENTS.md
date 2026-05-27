# AGENTS.md

## Agentic Workflow Sources

When changing GitHub Agentic Workflow source files under
`.github/workflows/`, run `gh aw compile` before committing and
include the regenerated `.lock.yml` files in the same change.

This applies to workflow markdown files such as `.github/workflows/*.md`
and shared workflow markdown such as `.github/workflows/shared/*.md`.
The `gh-aw lock check` CI job runs the same compile step and fails if
the generated workflow files are stale.

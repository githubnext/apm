---
description: "Agentic workflow recompilation: always recompile after changing workflow files"
---

# Agentic Workflows

After modifying any `.md` workflow file under `.github/workflows/`, always
recompile both agentic workflows and APM integration files before committing:

```bash
gh aw compile
apm compile
```

Commit the regenerated `.lock.yml` and integration files together with your
changes. The CI `APM Self-Check` job will fail if generated files are stale.

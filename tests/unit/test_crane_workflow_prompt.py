from __future__ import annotations

from pathlib import Path


ROOT = Path(__file__).resolve().parents[2]
CRANE_WORKFLOW = ROOT / ".github" / "workflows" / "crane.md"


def _workflow_text() -> str:
    return CRANE_WORKFLOW.read_text(encoding="utf-8")


def test_crane_acceptance_requires_shared_iteration_summary_for_pr_updates() -> None:
    text = _workflow_text()

    assert "### Accepted Iteration Summary" in text
    assert "single shared source" in text
    assert "PR body, PR comment, migration issue comment, and repo-memory history" in text
    assert "add-comment" in text
    assert "push-to-pull-request-branch" in text
    assert "ci: trigger checks" in text
    assert "unless it is the only new commit" in text


def test_crane_commit_guidance_provides_structured_summary_fallback() -> None:
    text = _workflow_text()

    assert "Subject: `[Crane: {migration-name}] Iteration <N>: <short description" in text
    assert "Changes:" in text
    assert "Run: {run_url}" in text
    assert text.index("Changes:") < text.index("Run: {run_url}")


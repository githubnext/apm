from __future__ import annotations

from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
WORKFLOW = ROOT / ".github" / "workflows" / "upstream-apm-sync.yml"


def _workflow_text() -> str:
    return WORKFLOW.read_text(encoding="utf-8")


def test_upstream_sync_workflow_fetches_and_merges_microsoft_apm() -> None:
    text = _workflow_text()

    assert "https://github.com/microsoft/apm.git" in text
    assert "git fetch upstream" in text
    assert "git merge --no-ff --no-edit" in text
    assert "automation/upstream-microsoft-apm-main" in text


def test_upstream_sync_workflow_uses_pr_auto_merge_not_squash() -> None:
    text = _workflow_text()

    assert "gh pr create" in text
    assert "gh pr merge" in text
    assert "--auto --merge --delete-branch" in text
    assert "--squash" not in text


def test_upstream_sync_workflow_tells_reviewers_to_update_go_coverage() -> None:
    text = _workflow_text()

    assert "Review the upstream Python diff" in text
    assert "real Go behavior tests" in text
    assert "tests/parity/upstream_contract_coverage.yml" in text

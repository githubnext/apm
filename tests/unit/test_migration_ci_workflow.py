from __future__ import annotations

from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
WORKFLOW = ROOT / ".github" / "workflows" / "migration-ci.yml"


def _workflow_text() -> str:
    return WORKFLOW.read_text(encoding="utf-8")


def test_migration_ci_enforces_completion_for_crane_prs_and_manual_runs() -> None:
    text = _workflow_text()

    assert "MIGRATION_COMPLETION_ENFORCED=$enforce_completion" in text
    assert "APM_ENFORCE_COMPLETION_GATES=1" in text
    assert 'github.event_name }}" = "workflow_dispatch"' in text
    assert 'github.event.pull_request.head.ref }}" == crane/*' in text
    assert "completion gates are enforced only for crane/* PRs and manual runs" in text


def test_migration_ci_collects_incomplete_evidence_for_non_crane_prs() -> None:
    text = _workflow_text()

    assert "--allow-failures" in text
    assert "Non-enforcing migration evidence run" in text
    assert "Python behavior contract tests are incomplete in collection mode." in text
    assert "Go parity tests are incomplete in collection mode." in text

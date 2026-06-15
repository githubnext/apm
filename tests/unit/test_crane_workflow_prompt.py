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


def test_crane_prompt_blocks_stale_completed_state_from_finishing() -> None:
    text = _workflow_text()

    assert "stale_completed_state" in text
    assert "active label wins" in text
    assert "ignore any pre-existing `Completed: true`" in text
    assert "Restore the issue to active migration state" in text
    assert "`crane-completed` but the scheduler selected it in `stale_completed_state`" in text
    assert "deterministic completion gate passes on the pushed PR head" in text
    assert "same-run sandbox score is only evidence" in text


def test_crane_completion_is_two_phase_and_pr_head_gated() -> None:
    text = _workflow_text()

    assert "Reaching the target metric does **not** complete the migration in this run" in text
    assert "Completion Candidate: true" in text
    assert "Completion Gate: up-to-date-pr-head-checks" in text
    assert "leave the `crane-migration` label on the issue" in text
    assert "current PR head contains the current base branch SHA" in text
    assert "every check for the current up-to-date PR head is terminal success" in text
    assert "Completion Gate Status: passed:<sha>" in text


def test_crane_base_sync_strips_protected_workflow_files_from_push_patch() -> None:
    text = _workflow_text()

    assert "trusted base-branch sync noise" in text
    assert "git checkout ORIG_HEAD -- <path>" in text
    assert (
        "safe-output patch for an existing Crane PR must not include protected workflow/config files"
        in text
    )


def test_crane_push_to_pr_branch_allows_protected_files() -> None:
    text = _workflow_text()

    push_config = text.split("push-to-pull-request-branch:", 1)[1].split("create-issue:", 1)[0]
    assert "protected-files: allowed" in push_config


def test_crane_state_template_tracks_completion_candidate_gate() -> None:
    text = _workflow_text()

    assert "| Completion Candidate | false |" in text
    assert "| Completion Gate | up-to-date-pr-head-checks |" in text
    assert "| Completion Gate Status | -- |" in text
    assert "Whether the target metric has been reached and the migration is waiting" in text

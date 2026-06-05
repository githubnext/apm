from __future__ import annotations

import importlib.util
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
SCHEDULER_PATH = ROOT / ".github" / "workflows" / "scripts" / "crane_scheduler.py"

spec = importlib.util.spec_from_file_location("crane_scheduler", SCHEDULER_PATH)
assert spec is not None
crane_scheduler = importlib.util.module_from_spec(spec)
assert spec.loader is not None
spec.loader.exec_module(crane_scheduler)


def test_completed_state_skips_inactive_migration() -> None:
    should_skip, reason = crane_scheduler.check_skip_conditions({"completed": True})

    assert should_skip is True
    assert reason == "completed: target metric reached"


def test_active_issue_overrides_stale_completed_state() -> None:
    should_skip, reason = crane_scheduler.check_skip_conditions(
        {"completed": True},
        issue_active=True,
    )

    assert should_skip is False
    assert reason is None


def test_active_issue_does_not_override_pause() -> None:
    should_skip, reason = crane_scheduler.check_skip_conditions(
        {"completed": True, "paused": True, "pause_reason": "manual hold"},
        issue_active=True,
    )

    assert should_skip is True
    assert reason == "paused: manual hold"


def test_machine_state_completed_string_is_recognized() -> None:
    assert crane_scheduler.is_completed_state({"completed": "true"}) is True


def test_parse_machine_state_accepts_bracketed_status_heading() -> None:
    state = crane_scheduler.parse_machine_state(
        """# Crane: sample

## [*] Machine State

| Field | Value |
|-------|-------|
| Last Run | 2026-06-05T16:10:36Z |
| Iteration Count | 67 |
| PR | #104 |
| Completed | true |
| Recent Statuses | accepted, rejected |

---

## [list] Migration Info
"""
    )

    assert state["last_run"] == "2026-06-05T16:10:36Z"
    assert state["iteration_count"] == 67
    assert state["completed"] is True
    assert state["recent_statuses"] == ["accepted", "rejected"]
    assert "-------" not in state


def test_issue_label_detection_accepts_github_label_payloads() -> None:
    issue = {"labels": [{"name": "crane-completed"}, "automation"]}

    assert crane_scheduler._issue_has_label(issue, "crane-completed") is True
    assert crane_scheduler._issue_has_label(issue, "automation") is True
    assert crane_scheduler._issue_has_label(issue, "crane-migration") is False


def test_completed_label_without_open_pr_is_recovered_as_stale() -> None:
    stale, recovered, event = crane_scheduler.evaluate_completed_label_recovery(
        "crane-migration-python-to-go-full-apm-cli-rewrite",
        {"completed": True},
        issue_active=False,
        issue_completed_label=True,
        repo="githubnext/apm",
        github_token="token",
        find_pr=lambda *_args: None,
    )

    assert stale is True
    assert recovered is True
    assert event == ("stale_no_pr", None, "no-open-migration-pr")

    should_skip, reason = crane_scheduler.check_skip_conditions(
        {"completed": True},
        issue_active=recovered,
    )
    assert should_skip is False
    assert reason is None


def test_completed_label_with_unknown_pr_gate_is_recovered_as_stale() -> None:
    stale, recovered, event = crane_scheduler.evaluate_completed_label_recovery(
        "crane-migration-python-to-go-full-apm-cli-rewrite",
        {"completed": True},
        issue_active=False,
        issue_completed_label=True,
        repo="githubnext/apm",
        github_token="token",
        find_pr=lambda *_args: 104,
        check_gate=lambda *_args: (None, "checks-unavailable:2699b7d"),
    )

    assert stale is True
    assert recovered is True
    assert event == ("stale_gate", 104, "checks-unavailable:2699b7d")


def test_pr_head_gate_fails_when_any_check_is_not_success() -> None:
    def fake_http_get_json(url, _headers, timeout=30):
        del timeout
        if url.endswith("/pulls/102"):
            return {"head": {"sha": "abcdef1234567890"}}, None
        if "/check-runs" in url:
            return {
                "check_runs": [
                    {"name": "Go Tests", "status": "completed", "conclusion": "success"},
                    {"name": "Parity Gate", "status": "completed", "conclusion": "failure"},
                ]
            }, None
        raise AssertionError(f"unexpected URL: {url}")

    passed, reason = crane_scheduler.get_pr_head_check_gate(
        "githubnext/apm",
        102,
        "token",
        http_get_json=fake_http_get_json,
    )

    assert passed is False
    assert reason.startswith("failing:abcdef123456")
    assert "Parity Gate:completed:failure" in reason


def test_pr_head_gate_passes_only_when_all_checks_succeed() -> None:
    def fake_http_get_json(url, _headers, timeout=30):
        del timeout
        if url.endswith("/pulls/102"):
            return {"head": {"sha": "abcdef1234567890"}}, None
        if "/check-runs" in url:
            return {
                "check_runs": [
                    {"name": "Go Tests", "status": "completed", "conclusion": "success"},
                    {"name": "Parity Gate", "status": "completed", "conclusion": "success"},
                ]
            }, None
        raise AssertionError(f"unexpected URL: {url}")

    passed, reason = crane_scheduler.get_pr_head_check_gate(
        "githubnext/apm",
        102,
        "token",
        http_get_json=fake_http_get_json,
    )

    assert passed is True
    assert reason == "passed:abcdef123456"

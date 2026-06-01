from __future__ import annotations

import json
import os
import shutil
import subprocess
import tempfile
from pathlib import Path

import pytest

ROOT = Path(__file__).resolve().parents[2]


def _run_score(input_lines: list[str]) -> dict[str, object]:
    if shutil.which("go") is None:
        pytest.skip("Go toolchain is not installed")

    env = os.environ.copy()
    with tempfile.TemporaryDirectory(prefix="apm-go-cache-") as go_cache:
        env.setdefault("GOCACHE", go_cache)
        result = subprocess.run(
            ["go", "run", ".crane/scripts/score.go"],
            cwd=ROOT,
            input="\n".join(input_lines) + "\n",
            text=True,
            capture_output=True,
            check=True,
            env=env,
        )
    return json.loads(result.stdout)


def _event(action: str, test: str, *, output: str = "") -> str:
    return json.dumps(
        {
            "Action": action,
            "Package": "github.com/githubnext/apm/cmd/apm",
            "Test": test,
            "Output": output,
        }
    )


def _pass(test: str) -> list[str]:
    return [_event("run", test), _event("pass", test)]


def _gates(score: dict[str, object]) -> dict[str, dict[str, object]]:
    gates = score["gates"]
    assert isinstance(gates, list)
    return {gate["name"]: gate for gate in gates}


def _all_required_gate_events() -> list[str]:
    tests = [
        "TestParityCompletionHardGate",
        "TestParityCompletionSurfaceParity",
        "TestParityCompletionCommandMatrix",
        "TestParityCompletionHelpIdentical",
        "TestParityCompletionFunctionalContracts",
        "TestParityCompletionStateDiffContracts",
        "TestParityCompletionPythonSuite",
        "TestParityCompletionBenchmarks",
        "TestParityCompletionPythonBehaviorContracts",
    ]
    return [line for test in tests for line in _pass(test)]


def test_crane_score_blocks_help_only_completion() -> None:
    score = _run_score(
        [
            "not json",
            *_pass("TestParityCompletionHardGate"),
            *_pass("TestParityCompletionCommandMatrix"),
            *_pass("TestParityCompletionHelpIdentical"),
            *_pass("TestParityCompletionVersionEquivalent"),
            *_pass("TestParityCompletionInitParity"),
            *_pass("TestParityCompletionErrorParity"),
        ]
    )

    gates = _gates(score)

    assert score["migration_score"] < 1.0
    assert gates["python_reference_required"]["passing"] is True
    assert gates["go_tests_pass"]["passing"] is True
    assert gates["help_parity"]["passing"] is True
    assert gates["surface_parity"]["passing"] is False
    assert gates["functional_contracts"]["passing"] is False
    assert gates["state_diff_contracts"]["passing"] is False
    assert gates["python_tests_pass"]["passing"] is False
    assert gates["benchmarks_pass"]["passing"] is False
    assert gates["python_behavior_contracts"]["passing"] is False


def test_crane_score_reaches_one_only_when_all_deletion_grade_gates_pass() -> None:
    score = _run_score(_all_required_gate_events())

    assert score["migration_score"] == 1.0
    assert score["progress"] == 1.0
    assert all(gate["passing"] for gate in _gates(score).values())


def test_crane_score_forces_zero_without_python_reference() -> None:
    score = _run_score(
        [
            *_pass("TestParityCompletionSurfaceParity"),
            *_pass("TestParityCompletionCommandMatrix"),
            *_pass("TestParityCompletionHelpIdentical"),
            *_pass("TestParityCompletionFunctionalContracts"),
            *_pass("TestParityCompletionStateDiffContracts"),
            *_pass("TestParityCompletionPythonSuite"),
            *_pass("TestParityCompletionBenchmarks"),
            *_pass("TestParityCompletionPythonBehaviorContracts"),
        ]
    )

    gates = _gates(score)

    assert score["migration_score"] == 0
    assert gates["python_reference_required"]["passing"] is False
    assert "TestParityCompletionHardGate not found" in gates["python_reference_required"]["reason"]


def test_crane_score_blocks_known_exceptions() -> None:
    score = _run_score(
        [
            *_all_required_gate_events(),
            _event("output", "TestParityCompletionHelpIdentical", output="APPROVED-EXCEPTION: no"),
        ]
    )

    gates = _gates(score)

    assert score["migration_score"] < 1.0
    assert gates["no_known_exceptions"]["passing"] is False
    assert "approved exception" in gates["no_known_exceptions"]["reason"]

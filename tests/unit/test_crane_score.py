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

    with tempfile.TemporaryDirectory(prefix="apm-go-cache-") as go_cache:
        env = os.environ.copy()
        env.setdefault("GOCACHE", go_cache)
        if not env.get("HOME"):
            env["HOME"] = str(Path.home())
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


def _go_pass(test: str, package: str = "github.com/githubnext/apm/cmd/apm") -> list[str]:
    return [
        json.dumps({"Action": "run", "Package": package, "Test": test}),
        json.dumps({"Action": "pass", "Package": package, "Test": test}),
    ]


def _event(action: str, test: str, *, output: str = "") -> str:
    return json.dumps(
        {
            "Action": action,
            "Package": "github.com/githubnext/apm/cmd/apm",
            "Test": test,
            "Output": output,
        }
    )


def _package_pass(package: str = "github.com/githubnext/apm/cmd/apm") -> str:
    return json.dumps({"Action": "pass", "Package": package})


def _package_fail(package: str = "github.com/githubnext/apm/internal/config") -> str:
    return json.dumps({"Action": "fail", "Package": package})


def _parity_passes(count: int) -> list[str]:
    lines: list[str] = []
    for i in range(count):
        lines.extend(_go_pass(f"TestParity{i}"))
    return lines


def _deletion_gates() -> list[str]:
    return [
        '{"crane":"gate","name":"python_reference","passed":true}',
        '{"crane":"gate","name":"surface","passing":1,"total":1}',
        '{"crane":"gate","name":"help","passing":1,"total":1}',
        '{"crane":"gate","name":"functional","passing":1,"total":1}',
        '{"crane":"gate","name":"state_diff","passing":1,"total":1}',
        '{"crane":"gate","name":"python_behavior_contracts","passing":1,"total":1}',
        '{"crane":"gate","name":"known_exceptions","count":0}',
        '{"crane":"gate","name":"python_tests","passed":true}',
        '{"crane":"gate","name":"benchmarks","passed":true}',
    ]


def _completion_gate_events() -> list[str]:
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
    return [line for test in tests for line in _go_pass(test)]


def _behavior_contract_gate_output(passing: int, total: int) -> str:
    return _event(
        "output",
        "TestParityCompletionPythonBehaviorContracts",
        output=json.dumps(
            {
                "crane": "gate",
                "name": "python_behavior_contracts",
                "passing": passing,
                "total": total,
            }
        )
        + "\n",
    )


def _gates(score: dict[str, object]) -> dict[str, dict[str, object]]:
    gates = score["gates"]
    assert isinstance(gates, list)
    return {gate["name"]: gate for gate in gates}


def test_crane_score_counts_parity_events() -> None:
    score = _run_score(
        [
            "not json",
            *_go_pass("TestInstallParity"),
            *_go_pass("TestCompileParity"),
            _package_pass(),
        ]
    )

    assert score["migration_score"] == pytest.approx(2 / 302)
    assert score["progress"] == pytest.approx(2 / 302)
    assert score["parity_passing"] == 2
    assert score["parity_total"] == 302
    assert score["source_tests_passing"] == 247
    assert score["target_tests_passing"] == 2
    assert score["perf_ratio"] == 1.0
    assert score["deletion_grade_ready"] is False


def test_crane_score_applies_target_correctness_gate() -> None:
    score = _run_score(
        [
            *_go_pass("TestInstallParity"),
            '{"Action":"run","Package":"github.com/githubnext/apm/internal/config","Test":"TestConfig"}',
            '{"Action":"fail","Package":"github.com/githubnext/apm/internal/config","Test":"TestConfig"}',
            _package_fail(),
        ]
    )

    assert score["migration_score"] == 0
    assert score["progress"] == pytest.approx(1 / 302)
    assert score["target_tests_passing"] == 1
    assert score["go_tests_passing"] is False


def test_crane_score_can_reach_one_with_all_deletion_grade_gates() -> None:
    score = _run_score([*_parity_passes(302), _package_pass(), *_deletion_gates()])

    assert score["migration_score"] == 1.0
    assert score["deletion_grade_ready"] is True
    assert score["cutover_gates"] == {
        "python_reference_required": True,
        "surface_parity": 1.0,
        "help_parity": 1.0,
        "functional_contracts": 1.0,
        "state_diff_contracts": 1.0,
        "python_behavior_contracts": 1.0,
        "known_exceptions": 0,
        "go_tests": "pass",
        "python_tests": "pass",
        "benchmarks": "pass",
    }


@pytest.mark.parametrize(
    "bad_gate",
    [
        '{"crane":"gate","name":"python_reference","passed":false}',
        '{"crane":"gate","name":"surface","passing":0,"total":1}',
        '{"crane":"gate","name":"help","passing":0,"total":1}',
        '{"crane":"gate","name":"functional","passing":0,"total":1}',
        '{"crane":"gate","name":"state_diff","passing":0,"total":1}',
        '{"crane":"gate","name":"python_behavior_contracts","passing":0,"total":1}',
        '{"crane":"gate","name":"known_exceptions","count":1}',
        '{"crane":"gate","name":"python_tests","passed":false}',
        '{"crane":"gate","name":"benchmarks","passed":false}',
    ],
)
def test_crane_score_full_parity_but_bad_deletion_gate_cannot_reach_one(
    bad_gate: str,
) -> None:
    bad_gate_name = json.loads(bad_gate)["name"]
    gates = [line for line in _deletion_gates() if json.loads(line)["name"] != bad_gate_name]

    score = _run_score([*_parity_passes(302), _package_pass(), *gates, bad_gate])

    assert score["migration_score"] < 1.0
    assert score["deletion_grade_ready"] is False


def test_crane_score_full_parity_but_missing_deletion_gates_cannot_reach_one() -> None:
    score = _run_score([*_parity_passes(302), _package_pass()])

    assert score["migration_score"] < 1.0
    assert score["deletion_grade_ready"] is False


def test_crane_score_package_level_go_failure_blocks_one() -> None:
    score = _run_score([*_parity_passes(302), _package_fail(), *_deletion_gates()])

    assert score["migration_score"] == 0
    assert score["go_tests_passing"] is False
    assert score["deletion_grade_ready"] is False


def test_crane_score_rejects_empty_event_stream() -> None:
    if shutil.which("go") is None:
        pytest.skip("Go toolchain is not installed")

    with tempfile.TemporaryDirectory(prefix="apm-go-cache-") as go_cache:
        env = os.environ.copy()
        env.setdefault("GOCACHE", go_cache)
        if not env.get("HOME"):
            env["HOME"] = str(Path.home())
        result = subprocess.run(
            ["go", "run", ".crane/scripts/score.go"],
            cwd=ROOT,
            input="",
            text=True,
            capture_output=True,
            check=False,
            env=env,
        )

    assert result.returncode != 0
    assert "empty or incomplete" in result.stderr


def test_crane_score_reaches_one_with_completion_tests_and_explicit_behavior_gate() -> None:
    score = _run_score(
        [
            *_parity_passes(293),
            *_completion_gate_events(),
            _behavior_contract_gate_output(1, 1),
            _package_pass(),
        ]
    )

    assert score["migration_score"] == 1.0
    assert score["progress"] == 1.0
    assert score["deletion_grade_ready"] is True
    assert all(gate["passing"] for gate in _gates(score).values())


def test_crane_score_does_not_infer_behavior_contracts_from_test_name() -> None:
    score = _run_score([*_parity_passes(293), *_completion_gate_events(), _package_pass()])
    gates = _gates(score)

    assert score["progress"] == 1.0
    assert score["migration_score"] < 1.0
    assert score["deletion_grade_ready"] is False
    assert gates["python_behavior_contracts"]["passing"] is False


def test_crane_score_blocks_incomplete_behavior_contract_gate() -> None:
    score = _run_score(
        [
            *_parity_passes(293),
            *_completion_gate_events(),
            _behavior_contract_gate_output(0, 1),
            _package_pass(),
        ]
    )
    gates = _gates(score)

    assert score["progress"] == 1.0
    assert score["migration_score"] < 1.0
    assert score["deletion_grade_ready"] is False
    assert gates["python_behavior_contracts"]["passing"] is False


def test_crane_score_blocks_known_exceptions() -> None:
    score = _run_score(
        [
            *_parity_passes(293),
            *_completion_gate_events(),
            _event("output", "TestParityCompletionHelpIdentical", output="APPROVED-EXCEPTION: no"),
            _package_pass(),
        ]
    )

    gates = _gates(score)

    assert score["migration_score"] < 1.0
    assert score["deletion_grade_ready"] is False
    assert gates["no_known_exceptions"]["passing"] is False

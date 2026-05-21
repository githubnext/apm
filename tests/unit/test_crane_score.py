from __future__ import annotations

import json
import shutil
import subprocess
from pathlib import Path

import pytest

ROOT = Path(__file__).resolve().parents[2]


def _run_score(input_lines: list[str]) -> dict[str, object]:
    if shutil.which("go") is None:
        pytest.skip("Go toolchain is not installed")

    result = subprocess.run(
        ["go", "run", ".crane/scripts/score.go"],
        cwd=ROOT,
        input="\n".join(input_lines) + "\n",
        text=True,
        capture_output=True,
        check=True,
    )
    return json.loads(result.stdout)


def test_crane_score_counts_parity_events() -> None:
    score = _run_score(
        [
            "not json",
            '{"Action":"run","Package":"github.com/githubnext/apm/internal/parity","Test":"TestInstallParity"}',
            '{"Action":"pass","Package":"github.com/githubnext/apm/internal/parity","Test":"TestInstallParity"}',
            '{"Action":"run","Package":"github.com/githubnext/apm/internal/parity","Test":"TestCompileParity"}',
            '{"Action":"pass","Package":"github.com/githubnext/apm/internal/parity","Test":"TestCompileParity"}',
        ]
    )

    assert score["migration_score"] == pytest.approx(2 / 302)
    assert score["progress"] == pytest.approx(2 / 302)
    assert score["parity_passing"] == 2
    assert score["parity_total"] == 302
    assert score["source_tests_passing"] == 247
    assert score["target_tests_passing"] == 2
    assert score["perf_ratio"] == 1.0


def test_crane_score_applies_target_correctness_gate() -> None:
    score = _run_score(
        [
            '{"Action":"run","Package":"github.com/githubnext/apm/internal/parity","Test":"TestInstallParity"}',
            '{"Action":"pass","Package":"github.com/githubnext/apm/internal/parity","Test":"TestInstallParity"}',
            '{"Action":"run","Package":"github.com/githubnext/apm/internal/config","Test":"TestConfig"}',
        ]
    )

    assert score["migration_score"] == 0
    assert score["progress"] == pytest.approx(1 / 302)
    assert score["target_tests_passing"] == 1

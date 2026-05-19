from pathlib import Path

import yaml

REPO_ROOT = Path(__file__).resolve().parents[2]


def _load_ci_workflow():
    return yaml.safe_load((REPO_ROOT / ".github/workflows/ci.yml").read_text())


def test_go_build_test_matrix_includes_macos():
    job = _load_ci_workflow()["jobs"]["go-build-test"]

    assert job["runs-on"] == "${{ matrix.os }}"
    assert job["strategy"]["matrix"]["os"] == ["ubuntu-24.04", "macos-latest"]


def test_go_build_test_runs_build_and_tests():
    steps = _load_ci_workflow()["jobs"]["go-build-test"]["steps"]
    commands = [step.get("run") for step in steps]

    assert "go test ./..." in commands
    assert "go build ./..." in commands

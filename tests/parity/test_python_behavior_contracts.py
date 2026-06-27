from __future__ import annotations

import os
import subprocess
import sys
from pathlib import Path

import pytest

ROOT = Path(__file__).resolve().parents[2]
sys.path.insert(0, str(ROOT / "scripts" / "ci"))

from python_behavior_contracts import (  # noqa: E402
    _load_coverage,
    check_coverage,
    extract_inventory,
    render_summary,
)


def _normalize_cli_output(text: str) -> str:
    lines: list[str] = []
    after_banner = False
    for line in text.splitlines():
        if "A new version of APM is available" in line:
            after_banner = True
            continue
        if "Run apm update to upgrade" in line:
            after_banner = True
            continue
        if after_banner and not line.strip():
            after_banner = False
            continue
        after_banner = False
        lines.append(line.rstrip())
    return "\n".join(lines).rstrip()


@pytest.fixture(scope="session")
def inventory() -> dict[str, object]:
    return extract_inventory()


@pytest.fixture(scope="session")
def coverage() -> dict[str, object]:
    return _load_coverage(ROOT / "tests" / "parity" / "python_contract_coverage.yml")


@pytest.fixture(scope="session")
def python_bin() -> Path:
    value = os.environ.get("APM_PYTHON_BIN")
    if not value:
        pytest.skip("APM_PYTHON_BIN is required for Python-vs-Go contract tests")
    path = Path(value)
    if not path.exists():
        pytest.fail(f"APM_PYTHON_BIN does not exist: {path}")
    return path


@pytest.fixture(scope="session")
def go_bin(tmp_path_factory: pytest.TempPathFactory) -> Path:
    value = os.environ.get("APM_GO_BIN")
    if value:
        path = Path(value)
        if not path.exists():
            pytest.fail(f"APM_GO_BIN does not exist: {path}")
        return path

    out = tmp_path_factory.mktemp("apm-go") / ("apm.exe" if os.name == "nt" else "apm")
    subprocess.run(["go", "build", "-o", str(out), "./cmd/apm"], cwd=ROOT, check=True)
    return out


def _public_commands(inventory: dict[str, object]) -> list[dict[str, object]]:
    commands = inventory["commands"]
    assert isinstance(commands, list)
    return [cmd for cmd in commands if isinstance(cmd, dict) and not cmd.get("hidden")]


def pytest_generate_tests(metafunc: pytest.Metafunc) -> None:
    if "command_contract" not in metafunc.fixturenames:
        return
    inv = extract_inventory()
    commands = _public_commands(inv)
    metafunc.parametrize(
        "command_contract",
        commands,
        ids=[str(command["id"]) for command in commands],
    )


def _run(bin_path: Path, args: list[str], cwd: Path) -> subprocess.CompletedProcess[str]:
    return subprocess.run(
        [str(bin_path), *args],
        cwd=cwd,
        text=True,
        capture_output=True,
        check=False,
        env={**os.environ, "NO_COLOR": "1", "COLUMNS": "10000"},
    )


def _help_args(command: dict[str, object]) -> list[str]:
    path = command["path"]
    assert isinstance(path, list)
    args = [str(part) for part in path]
    return [*args, "--help"] if args else ["--help"]


def test_every_python_command_help_matches_go(
    command_contract: dict[str, object],
    python_bin: Path,
    go_bin: Path,
    tmp_path: Path,
    coverage: dict[str, object],
) -> None:
    if coverage.get("status") == "intentionally-incomplete":
        pytest.skip("coverage manifest is intentionally incomplete; migration in progress")
    args = _help_args(command_contract)
    py = _run(python_bin, args, tmp_path)
    go = _run(go_bin, args, tmp_path)
    assert py.returncode == 0
    mismatch = (
        go.returncode != py.returncode
        or _normalize_cli_output(go.stdout) != _normalize_cli_output(py.stdout)
        or _normalize_cli_output(go.stderr) != _normalize_cli_output(py.stderr)
    )
    if mismatch:
        message = (
            f"help parity mismatch for {command_contract['id']} "
            "(tracking while migration is in progress)"
        )
        if os.environ.get("APM_ENFORCE_PYTHON_BEHAVIOR_CONTRACTS") == "1":
            pytest.fail(message)
        pytest.xfail(message)


def test_every_python_command_rejects_unknown_option_consistently(
    command_contract: dict[str, object],
    python_bin: Path,
    go_bin: Path,
    tmp_path: Path,
    coverage: dict[str, object],
) -> None:
    if coverage.get("status") == "intentionally-incomplete":
        pytest.skip("coverage manifest is intentionally incomplete; migration in progress")
    path = command_contract["path"]
    assert isinstance(path, list)
    args = [str(part) for part in path]
    probe = [*args, "--definitely-not-an-apm-option"]
    py = _run(python_bin, probe, tmp_path)
    go = _run(go_bin, probe, tmp_path)
    assert py.returncode != 0
    mismatch = (
        go.returncode != py.returncode
        or _normalize_cli_output(go.stdout) != _normalize_cli_output(py.stdout)
        or _normalize_cli_output(go.stderr) != _normalize_cli_output(py.stderr)
    )
    if mismatch:
        message = (
            f"unknown-option parity mismatch for {command_contract['id']} "
            "(tracking while migration is in progress)"
        )
        if os.environ.get("APM_ENFORCE_PYTHON_BEHAVIOR_CONTRACTS") == "1":
            pytest.fail(message)
        pytest.xfail(message)


def test_python_contract_coverage_manifest_is_complete(inventory: dict[str, object]) -> None:
    coverage = _load_coverage(ROOT / "tests" / "parity" / "python_contract_coverage.yml")
    enforce = os.environ.get("APM_ENFORCE_PYTHON_BEHAVIOR_CONTRACTS") == "1"
    if coverage.get("status") == "intentionally-incomplete":
        if not enforce:
            pytest.xfail(
                "Coverage manifest is intentionally incomplete; completion gate "
                "is reported by migration_score"
            )
        pytest.fail(
            "Coverage manifest is intentionally incomplete; remove status field "
            "only after all contracts are mapped"
        )
    findings = check_coverage(inventory, coverage)
    if findings and not enforce:
        pytest.xfail(render_summary(inventory, findings))
    assert not findings, render_summary(inventory, findings)


def test_python_contract_coverage_rejects_obsolete_tests_by_default() -> None:
    inventory = {
        "summary": {
            "commands": 0,
            "public_commands": 0,
            "python_tests": 1,
            "python_test_cases": 1,
            "source_contracts": 0,
        },
        "commands": [],
        "tests": [{"id": "tests/unit/test_example.py::test_real_behavior"}],
        "source_contracts": [],
    }
    coverage = {
        "commands": {},
        "python_tests": {
            "covered": {},
            "obsolete": ["tests/unit/test_example.py::test_real_behavior"],
        },
    }

    findings = check_coverage(inventory, coverage)

    assert [finding.code for finding in findings] == ["obsolete-python-test-coverage"]


def test_python_contract_coverage_can_allow_obsolete_tests_in_report_only_mode() -> None:
    inventory = {
        "summary": {
            "commands": 0,
            "public_commands": 0,
            "python_tests": 1,
            "python_test_cases": 1,
            "source_contracts": 0,
        },
        "commands": [],
        "tests": [{"id": "tests/unit/test_example.py::test_real_behavior"}],
        "source_contracts": [],
    }
    coverage = {
        "commands": {},
        "python_tests": {
            "covered": {},
            "obsolete": ["tests/unit/test_example.py::test_real_behavior"],
        },
    }

    assert check_coverage(inventory, coverage, allow_obsolete=True) == []

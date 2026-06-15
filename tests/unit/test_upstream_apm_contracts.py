from __future__ import annotations

import subprocess
import sys
from pathlib import Path

import pytest
import yaml

ROOT = Path(__file__).resolve().parents[2]
sys.path.insert(0, str(ROOT / "scripts" / "ci"))

from upstream_apm_contracts import check_upstream_contracts  # noqa: E402


def _git(repo: Path, *args: str) -> str:
    return subprocess.run(
        ["git", *args],
        cwd=repo,
        text=True,
        capture_output=True,
        check=True,
    ).stdout.strip()


def _commit(repo: Path, message: str) -> str:
    _git(repo, "add", ".")
    _git(repo, "commit", "-m", message)
    return _git(repo, "rev-parse", "HEAD")


@pytest.fixture()
def repo(tmp_path: Path) -> Path:
    _git(tmp_path, "init")
    _git(tmp_path, "config", "user.email", "test@example.com")
    _git(tmp_path, "config", "user.name", "Test User")
    (tmp_path / "src" / "apm_cli").mkdir(parents=True)
    (tmp_path / "tests" / "unit").mkdir(parents=True)
    (tmp_path / "cmd" / "apm").mkdir(parents=True)
    (tmp_path / "src" / "apm_cli" / "__init__.py").write_text("", encoding="utf-8")
    (tmp_path / "cmd" / "apm" / "real_behavior_test.go").write_text(
        "package main\n\nfunc TestGoUpstreamBehavior(t *testing.T) {}\n",
        encoding="utf-8",
    )
    _commit(tmp_path, "base")
    return tmp_path


def _write_manifest(repo: Path, data: dict[str, object]) -> Path:
    path = repo / "coverage.yml"
    path.write_text(yaml.safe_dump(data, sort_keys=False), encoding="utf-8")
    return path


def test_upstream_contracts_pass_when_reviewed_sha_matches_head(repo: Path) -> None:
    head = _git(repo, "rev-parse", "HEAD")
    manifest = _write_manifest(
        repo,
        {
            "schema_version": 1,
            "upstream": {"baseline_sha": head, "reviewed_sha": head},
            "reviewed_ranges": [],
        },
    )

    result = check_upstream_contracts(
        root=repo,
        coverage_path=manifest,
        upstream_ref=head,
        head_ref=head,
    )

    assert result.freshness_ok is True
    assert result.contracts_passing == result.contracts_total == 1
    assert result.findings == []


def test_upstream_contracts_fail_when_upstream_adds_unreviewed_python_behavior(
    repo: Path,
) -> None:
    base = _git(repo, "rev-parse", "HEAD")
    (repo / "src" / "apm_cli" / "new_feature.py").write_text(
        "def kiro_target():\n    return 'kiro'\n",
        encoding="utf-8",
    )
    (repo / "tests" / "unit" / "test_new_feature.py").write_text(
        "def test_kiro_target():\n    assert True\n",
        encoding="utf-8",
    )
    upstream = _commit(repo, "upstream behavior")
    manifest = _write_manifest(
        repo,
        {
            "schema_version": 1,
            "upstream": {"baseline_sha": base, "reviewed_sha": base},
            "reviewed_ranges": [],
        },
    )

    result = check_upstream_contracts(
        root=repo,
        coverage_path=manifest,
        upstream_ref=upstream,
        head_ref=base,
    )

    assert result.freshness_ok is False
    assert result.contracts_passing == 0
    assert result.contracts_total == 2
    assert {finding.code for finding in result.findings} == {"missing-upstream-go-tests"}


def test_upstream_contracts_require_chained_reviewed_range_when_sha_advances(
    repo: Path,
) -> None:
    base = _git(repo, "rev-parse", "HEAD")
    (repo / "src" / "apm_cli" / "new_feature.py").write_text(
        "def source_base():\n    return 'marketplace'\n",
        encoding="utf-8",
    )
    upstream = _commit(repo, "upstream source")
    manifest = _write_manifest(
        repo,
        {
            "schema_version": 1,
            "upstream": {"baseline_sha": base, "reviewed_sha": upstream},
            "reviewed_ranges": [],
        },
    )

    result = check_upstream_contracts(
        root=repo,
        coverage_path=manifest,
        upstream_ref=upstream,
        head_ref=upstream,
    )

    assert result.freshness_ok is True
    assert [finding.code for finding in result.findings] == ["missing-reviewed-range"]
    assert result.contracts_passing == 0


def test_upstream_contracts_accept_reviewed_range_with_existing_go_tests(repo: Path) -> None:
    base = _git(repo, "rev-parse", "HEAD")
    (repo / "src" / "apm_cli" / "new_feature.py").write_text(
        "def optional_registry_inputs():\n    return True\n",
        encoding="utf-8",
    )
    (repo / "tests" / "unit" / "test_new_feature.py").write_text(
        "class TestOptionalRegistryInputs:\n"
        "    def test_preserves_optional_input(self):\n"
        "        assert True\n",
        encoding="utf-8",
    )
    upstream = _commit(repo, "upstream contracts")
    manifest = _write_manifest(
        repo,
        {
            "schema_version": 1,
            "upstream": {"baseline_sha": base, "reviewed_sha": upstream},
            "reviewed_ranges": [
                {
                    "from": base,
                    "to": upstream,
                    "source_contracts": {
                        "src/apm_cli/new_feature.py::optional_registry_inputs": {
                            "go_tests": ["TestGoUpstreamBehavior"]
                        }
                    },
                    "python_tests": {
                        (
                            "tests/unit/test_new_feature.py::"
                            "TestOptionalRegistryInputs::test_preserves_optional_input"
                        ): {"go_tests": ["TestGoUpstreamBehavior"]}
                    },
                }
            ],
        },
    )

    result = check_upstream_contracts(
        root=repo,
        coverage_path=manifest,
        upstream_ref=upstream,
        head_ref=upstream,
    )

    assert result.freshness_ok is True
    assert result.contracts_passing == result.contracts_total == 2
    assert result.findings == []

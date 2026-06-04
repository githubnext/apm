from __future__ import annotations

from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]


def test_readme_documents_go_cli_migration_usage() -> None:
    readme = (ROOT / "README.md").read_text(encoding="utf-8")

    required_snippets = [
        "The Python CLI remains the reference implementation and parity oracle.",
        "The Go CLI lives under `cmd/apm`",
        "githubnext/apm#91",
        "githubnext/apm#78",
        "githubnext/apm#93",
        "go build -o ./dist/apm-go ./cmd/apm",
        "./dist/apm-go --help",
        "./dist/apm-go init --yes",
        'export APM_PYTHON_BIN="$PWD/.venv/bin/apm"',
        "go test -json ./... | go run .crane/scripts/score.go",
        "gh workflow run migration-ci.yml --repo githubnext/apm --ref main",
        "`migration-benchmark-evidence`",
        "`Go/Python` ratio is the Go median duration divided by the Python median",
        "fixture-backed commands that read, write, execute, or fail against realistic APM",
        "`apm.yml`, `apm.lock.yaml`, installed `apm_modules`",
    ]

    for snippet in required_snippets:
        assert snippet in readme


def test_readme_distinguishes_parity_gate_from_full_historical_coverage() -> None:
    readme = (ROOT / "README.md").read_text(encoding="utf-8")

    assert "`migration_score=1`, `progress=1`, and `706/706` parity tests passing" in readme
    assert (
        "That gate is not the same thing as claiming all historical Python integration," in readme
    )
    assert "`APM_PYTHON_BIN` is required for the hard Python-vs-Go parity gate" in readme

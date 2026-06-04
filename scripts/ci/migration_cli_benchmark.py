#!/usr/bin/env python3
"""Compare Python and Go CLI latency for migration benchmark workloads."""

from __future__ import annotations

import argparse
import json
import os
import re
import statistics
import subprocess
import tempfile
import time
from dataclasses import dataclass
from pathlib import Path

FixtureName = str


@dataclass(frozen=True)
class BenchmarkCommand:
    name: str
    args: list[str]
    fixture: FixtureName
    workload: str


COMMANDS: list[BenchmarkCommand] = [
    BenchmarkCommand(
        name="startup help",
        args=["--help"],
        fixture="none",
        workload="Cold CLI startup and top-level help rendering.",
    ),
    BenchmarkCommand(
        name="startup version",
        args=["--version"],
        fixture="none",
        workload="Cold CLI startup and version rendering.",
    ),
    BenchmarkCommand(
        name="init scaffold",
        args=["init", "--yes"],
        fixture="empty-project",
        workload="Creates a new apm.yml in an otherwise empty project directory.",
    ),
    BenchmarkCommand(
        name="targets json",
        args=["targets", "--json"],
        fixture="installed-project",
        workload="Reads configured project targets from apm.yml and emits machine output.",
    ),
    BenchmarkCommand(
        name="script list",
        args=["list"],
        fixture="installed-project",
        workload="Reads apm.yml scripts and renders the runnable script inventory.",
    ),
    BenchmarkCommand(
        name="deps list",
        args=["deps", "list"],
        fixture="installed-project",
        workload="Scans apm_modules package directories and apm.lock.yaml metadata.",
    ),
    BenchmarkCommand(
        name="deps tree",
        args=["deps", "tree"],
        fixture="installed-project",
        workload="Builds a dependency tree from apm.lock.yaml and installed package metadata.",
    ),
    BenchmarkCommand(
        name="install dry-run",
        args=["install", "--dry-run", "--no-policy"],
        fixture="installed-project",
        workload="Builds an offline install preview from manifest dependencies.",
    ),
    BenchmarkCommand(
        name="compile dry-run",
        args=["compile", "--dry-run", "--all", "--local-only"],
        fixture="installed-project",
        workload="Discovers local primitives and plans compilation for all targets without writes.",
    ),
    BenchmarkCommand(
        name="pack dry-run",
        args=["pack", "--dry-run", "--offline", "--marketplace", "none"],
        fixture="installed-project",
        workload="Resolves local package contents and bundle metadata without writing artifacts.",
    ),
    BenchmarkCommand(
        name="audit file scan",
        args=["audit", "--file", ".apm/instructions/bench-00.instructions.md"],
        fixture="installed-project",
        workload="Scans a real prompt instruction file for hidden Unicode content.",
    ),
]


def _run_once(binary: str, args: list[str], cwd: Path, env: dict[str, str]) -> dict[str, object]:
    start = time.perf_counter()
    proc = subprocess.run(  # noqa: S603 -- benchmark intentionally executes supplied CLIs.
        [binary, *args],
        cwd=cwd,
        env=env,
        text=True,
        capture_output=True,
        timeout=30,
        check=False,
    )
    elapsed = time.perf_counter() - start
    return {
        "elapsed_seconds": elapsed,
        "returncode": proc.returncode,
        "stdout_bytes": len(proc.stdout.encode("utf-8")),
        "stderr_bytes": len(proc.stderr.encode("utf-8")),
    }


def _write(path: Path, content: str) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(content, encoding="utf-8")


def _safe_name(name: str) -> str:
    return re.sub(r"[^a-zA-Z0-9_.-]+", "-", name).strip("-")


def _write_empty_project(workdir: Path) -> None:
    _write(workdir / "README.md", "# Benchmark fixture\n")


def _write_installed_project(workdir: Path) -> None:
    _write_empty_project(workdir)
    for directory in [
        ".github",
        ".claude",
        ".cursor/rules",
        ".codex",
        "src/apm_bench",
        ".apm/instructions",
        ".apm/chatmodes",
        "apm_modules/microsoft/apm-package-alpha/.apm/instructions",
        "apm_modules/github/agent-toolkit/.apm/instructions",
    ]:
        (workdir / directory).mkdir(parents=True, exist_ok=True)

    _write(
        workdir / "apm.yml",
        """name: benchmark-project
version: 1.2.3
description: Realistic migration benchmark fixture
author: benchmark
targets:
  - copilot
  - claude
  - cursor
dependencies:
  apm:
    - microsoft/apm-package-alpha#v1.0.0
    - github/agent-toolkit#v2.3.4
  mcp: []
scripts:
  build: Build benchmark artifacts
  test: Run the test suite
  lint: Run lint checks
  release: Prepare release artifacts
includes: auto
""",
    )
    _write(
        workdir / "apm.lock.yaml",
        """lockfile_version: "1"
generated_at: "2026-01-01T00:00:00+00:00"
apm_version: benchmark
dependencies:
  - repo_url: microsoft/apm-package-alpha
    resolved_ref: v1.0.0
    resolved_commit: "1111111111111111111111111111111111111111"
    version: 1.0.0
    package_type: instructions
    deployed_files:
      - .github/copilot-instructions.md
  - repo_url: github/agent-toolkit
    resolved_ref: v2.3.4
    resolved_commit: "2222222222222222222222222222222222222222"
    version: 2.3.4
    depth: 2
    resolved_by: microsoft/apm-package-alpha
    package_type: instructions
    deployed_files:
      - CLAUDE.md
local_deployed_files:
  - .github/copilot-instructions.md
  - CLAUDE.md
  - .cursor/rules/AGENTS.md
""",
    )
    _write(
        workdir / ".github/copilot-instructions.md",
        "# Copilot Benchmark Instructions\n\nUse the local benchmark context.\n",
    )
    _write(
        workdir / "CLAUDE.md",
        "# Claude Benchmark Instructions\n\nUse the local benchmark context.\n",
    )
    _write(
        workdir / ".cursor/rules/AGENTS.md",
        "# Cursor Benchmark Instructions\n\nUse the local benchmark context.\n",
    )

    for index in range(16):
        _write(
            workdir / f".apm/instructions/bench-{index:02d}.instructions.md",
            f"""---
applyTo: "src/**/*.py"
description: Benchmark instruction {index}
---
# Benchmark Instruction {index}

Keep implementation clear and tested.

- Check input boundaries.
- Prefer small functions.
- Leave useful diagnostics for failures.
""",
        )
    for index in range(2):
        _write(
            workdir / f".apm/chatmodes/reviewer-{index}.chatmode.md",
            f"""---
description: Review benchmark fixture {index}
---
# Reviewer {index}

Review for correctness, maintainability, and test coverage.
""",
        )
    for index in range(24):
        _write(
            workdir / f"src/apm_bench/module_{index:02d}.py",
            f'"""Benchmark source module {index}."""\n\nVALUE_{index} = {index}\n',
        )

    packages = [
        ("microsoft", "apm-package-alpha", "1.0.0"),
        ("github", "agent-toolkit", "2.3.4"),
    ]
    for owner, repo, version in packages:
        package_dir = workdir / "apm_modules" / owner / repo
        _write(
            package_dir / "apm.yml",
            f"""name: {repo}
version: {version}
description: Fixture dependency package
author: benchmark
dependencies:
  apm: []
  mcp: []
""",
        )
        _write(
            package_dir / f".apm/instructions/{repo}.instructions.md",
            f"""---
applyTo: "**/*"
description: Installed package instruction for {repo}
---
# {repo}

Installed dependency instruction used by migration benchmarks.
""",
        )


def _workspace(base: Path, command: BenchmarkCommand, run_index: int) -> Path:
    if command.fixture == "none":
        return base

    workdir = base / _safe_name(command.name) / str(run_index)
    workdir.mkdir(parents=True, exist_ok=True)

    if command.fixture == "empty-project":
        _write_empty_project(workdir)
    elif command.fixture == "installed-project":
        _write_installed_project(workdir)
    else:
        raise ValueError(f"unknown benchmark fixture: {command.fixture}")

    return workdir


def _measure(
    *,
    binary: str,
    command: BenchmarkCommand,
    repeats: int,
    base: Path,
    env: dict[str, str],
) -> dict[str, object]:
    base.mkdir(parents=True, exist_ok=True)
    samples: list[dict[str, object]] = []
    for index in range(repeats):
        cwd = _workspace(base, command, index)
        samples.append(_run_once(binary, command.args, cwd, env))

    elapsed = [float(sample["elapsed_seconds"]) for sample in samples]
    return {
        "median_seconds": statistics.median(elapsed),
        "min_seconds": min(elapsed),
        "max_seconds": max(elapsed),
        "returncodes": sorted({int(sample["returncode"]) for sample in samples}),
        "samples": samples,
    }


def _speed_label(ratio: float) -> str:
    if ratio == 0:
        return "n/a"
    if ratio < 1:
        return f"{1 / ratio:.2f}x faster"
    if ratio > 1:
        return f"{ratio:.2f}x slower"
    return "same"


def _markdown(results: list[dict[str, object]], max_ratio: float) -> str:
    lines = [
        "## Migration CLI Benchmark",
        "",
        "Includes startup baselines plus fixture-backed real-world commands. "
        "The installed-project fixture contains apm.yml, apm.lock.yaml, "
        "apm_modules packages, local .apm primitives, target directories, "
        "deployed prompt files, and sample source files.",
        "The harness checks return-code parity for each command. Detailed stdout/stderr "
        "byte counts are kept in the JSON samples, but this is not an output-parity test.",
        "",
        f"Max allowed Go/Python median ratio: `{max_ratio:.2f}`",
        "",
        "| Benchmark | Command | Fixture | Python median | Go median | Go/Python | Result | Return codes |",
        "|---|---|---|---:|---:|---:|---|---|",
    ]
    for row in results:
        lines.append(
            "| {name} | `{command}` | {fixture} | {python:.4f}s | {go:.4f}s | {ratio:.2f}x | {result} | {codes} |".format(
                name=row["name"],
                command=row["command"],
                fixture=row["fixture"],
                python=row["python_median_seconds"],
                go=row["go_median_seconds"],
                ratio=row["ratio"],
                result=_speed_label(float(row["ratio"])),
                codes=row["returncodes"],
            )
        )
    lines.extend(["", "### Workloads", ""])
    for row in results:
        lines.append(f"- **{row['name']}**: {row['workload']}")
    lines.append("")
    return "\n".join(lines)


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--python-bin", required=True)
    parser.add_argument("--go-bin", required=True)
    parser.add_argument("--json-out", required=True)
    parser.add_argument("--markdown-out", required=True)
    parser.add_argument("--max-ratio", type=float, default=5.0)
    parser.add_argument("--repeats", type=int, default=5)
    args = parser.parse_args()

    env = os.environ.copy()
    env.update(
        {
            "NO_COLOR": "1",
            "TERM": "dumb",
            "PYTHONUNBUFFERED": "1",
        }
    )

    results: list[dict[str, object]] = []
    failures: list[str] = []
    with tempfile.TemporaryDirectory(prefix="apm-migration-bench-") as tmp:
        base = Path(tmp)
        for command in COMMANDS:
            python_result = _measure(
                binary=args.python_bin,
                command=command,
                repeats=args.repeats,
                base=base / "python" / _safe_name(command.name),
                env=env,
            )
            go_result = _measure(
                binary=args.go_bin,
                command=command,
                repeats=args.repeats,
                base=base / "go" / _safe_name(command.name),
                env=env,
            )

            python_median = float(python_result["median_seconds"])
            go_median = float(go_result["median_seconds"])
            ratio = go_median / max(python_median, 0.000001)
            returncodes = {
                "python": python_result["returncodes"],
                "go": go_result["returncodes"],
            }

            row = {
                "name": command.name,
                "command": " ".join(command.args),
                "fixture": command.fixture,
                "workload": command.workload,
                "python": python_result,
                "go": go_result,
                "python_median_seconds": python_median,
                "go_median_seconds": go_median,
                "ratio": ratio,
                "returncodes": returncodes,
            }
            results.append(row)

            if python_result["returncodes"] != go_result["returncodes"]:
                failures.append(f"{command.name}: return codes differ: {returncodes}")
            if ratio > args.max_ratio:
                failures.append(
                    f"{command.name}: Go median {ratio:.2f}x slower than Python "
                    f"(limit {args.max_ratio:.2f}x)"
                )

    json_path = Path(args.json_out)
    markdown_path = Path(args.markdown_out)
    json_path.write_text(
        json.dumps({"results": results, "failures": failures}, indent=2), encoding="utf-8"
    )
    markdown_path.write_text(_markdown(results, args.max_ratio), encoding="utf-8")

    print(markdown_path.read_text(encoding="utf-8"))
    if failures:
        for failure in failures:
            print(f"::error::{failure}")
        return 1
    return 0


if __name__ == "__main__":
    raise SystemExit(main())

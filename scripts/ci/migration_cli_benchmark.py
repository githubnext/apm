#!/usr/bin/env python3
"""Compare Python and Go CLI latency for migration smoke commands."""

from __future__ import annotations

import argparse
import json
import os
import statistics
import subprocess
import tempfile
import time
from pathlib import Path

COMMANDS: list[tuple[str, list[str], bool]] = [
    ("help", ["--help"], False),
    ("version", ["--version"], False),
    ("compile-help", ["compile", "--help"], False),
    ("install-help", ["install", "--help"], False),
    ("pack-help", ["pack", "--help"], False),
    ("audit-help", ["audit", "--help"], False),
    ("init-yes", ["init", "--yes"], True),
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


def _workspace(base: Path, name: str, run_index: int) -> Path:
    workdir = base / name / str(run_index)
    workdir.mkdir(parents=True, exist_ok=True)
    (workdir / "README.md").write_text("# Benchmark fixture\n", encoding="utf-8")
    return workdir


def _measure(
    *,
    binary: str,
    args: list[str],
    mutates_workspace: bool,
    repeats: int,
    base: Path,
    label: str,
    env: dict[str, str],
) -> dict[str, object]:
    base.mkdir(parents=True, exist_ok=True)
    samples: list[dict[str, object]] = []
    for index in range(repeats):
        cwd = _workspace(base, label, index) if mutates_workspace else base
        samples.append(_run_once(binary, args, cwd, env))

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
        f"Max allowed Go/Python median ratio: `{max_ratio:.2f}`",
        "",
        "| Command | Python median | Go median | Go/Python | Result | Return codes |",
        "|---|---:|---:|---:|---|---|",
    ]
    for row in results:
        lines.append(
            "| {command} | {python:.4f}s | {go:.4f}s | {ratio:.2f}x | {result} | {codes} |".format(
                command=row["command"],
                python=row["python_median_seconds"],
                go=row["go_median_seconds"],
                ratio=row["ratio"],
                result=_speed_label(float(row["ratio"])),
                codes=row["returncodes"],
            )
        )
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
        for command, command_args, mutates_workspace in COMMANDS:
            python_result = _measure(
                binary=args.python_bin,
                args=command_args,
                mutates_workspace=mutates_workspace,
                repeats=args.repeats,
                base=base / "python" / command,
                label=command,
                env=env,
            )
            go_result = _measure(
                binary=args.go_bin,
                args=command_args,
                mutates_workspace=mutates_workspace,
                repeats=args.repeats,
                base=base / "go" / command,
                label=command,
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
                "command": " ".join(command_args),
                "python": python_result,
                "go": go_result,
                "python_median_seconds": python_median,
                "go_median_seconds": go_median,
                "ratio": ratio,
                "returncodes": returncodes,
            }
            results.append(row)

            if python_result["returncodes"] != go_result["returncodes"]:
                failures.append(f"{command}: return codes differ: {returncodes}")
            if ratio > args.max_ratio:
                failures.append(
                    f"{command}: Go median {ratio:.2f}x slower than Python "
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

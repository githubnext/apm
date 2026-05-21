#!/usr/bin/env python3
"""Measure Go CLI parity against the Python CLI with observable executions."""

from __future__ import annotations

import argparse
import ast
import importlib
import json
import os
import shutil
import subprocess
import sys
from collections.abc import Callable, Mapping, Sequence
from pathlib import Path
from typing import Any

DEFAULT_TIMEOUT_SECONDS = 20


def _repo_root() -> Path:
    return Path(__file__).resolve().parents[1]


def _python_env() -> dict[str, str]:
    env = os.environ.copy()
    src = str(_repo_root() / "src")
    env["PYTHONPATH"] = src + os.pathsep + env["PYTHONPATH"] if env.get("PYTHONPATH") else src
    env["APM_E2E_TESTS"] = "1"
    return env


def enumerate_python_subcommands() -> list[str]:
    """Return visible top-level commands registered by the Python Click CLI."""
    sys.path.insert(0, str(_repo_root() / "src"))
    try:
        cli_module = importlib.import_module("apm_cli.cli")
        commands = getattr(cli_module.cli, "commands", {})
        return sorted(name for name, command in commands.items() if not getattr(command, "hidden", False))
    except ModuleNotFoundError:
        return enumerate_python_subcommands_static()


def enumerate_python_subcommands_static() -> list[str]:
    """Parse ``cli.py`` to enumerate commands without importing Click."""
    cli_path = _repo_root() / "src/apm_cli/cli.py"
    tree = ast.parse(cli_path.read_text())
    import_names: dict[str, str] = {}
    commands: set[str] = set()

    for node in tree.body:
        if isinstance(node, ast.ImportFrom):
            for alias in node.names:
                import_names[alias.asname or alias.name] = alias.name

    for node in ast.walk(tree):
        if not isinstance(node, ast.Call):
            continue
        if not (
            isinstance(node.func, ast.Attribute)
            and node.func.attr == "add_command"
            and isinstance(node.func.value, ast.Name)
            and node.func.value.id == "cli"
            and node.args
        ):
            continue
        explicit_name = next(
            (
                keyword.value.value
                for keyword in node.keywords
                if keyword.arg == "name" and isinstance(keyword.value, ast.Constant)
            ),
            None,
        )
        if explicit_name:
            commands.add(str(explicit_name))
        elif isinstance(node.args[0], ast.Name):
            commands.add(import_names.get(node.args[0].id, node.args[0].id))

    return sorted(commands)


def load_fixture_args(path: Path | None) -> dict[str, list[str]]:
    """Load optional command fixtures from JSON.

    The JSON format is ``{"command": ["arg", ...]}``. Commands without a
    fixture use ``--help`` so the metric is executable from a fresh checkout,
    while still allowing deeper command fixtures as slices are migrated.
    """
    if path is None:
        return {}
    raw = json.loads(path.read_text())
    if not isinstance(raw, dict):
        raise ValueError("fixture JSON must be an object mapping commands to argument lists")
    fixtures: dict[str, list[str]] = {}
    for command, args in raw.items():
        if (
            not isinstance(command, str)
            or not isinstance(args, list)
            or not all(isinstance(arg, str) for arg in args)
        ):
            raise ValueError("fixture JSON values must be lists of strings")
        fixtures[command] = args
    return fixtures


def fixture_args_for(command: str, fixtures: Mapping[str, Sequence[str]]) -> list[str]:
    return list(fixtures.get(command, ["--help"]))


def run_command(argv: Sequence[str], timeout: int = DEFAULT_TIMEOUT_SECONDS) -> dict[str, Any]:
    completed = subprocess.run(  # noqa: S603 -- argv is constructed from local CLI fixtures.
        list(argv),
        cwd=_repo_root(),
        env=_python_env(),
        text=True,
        capture_output=True,
        timeout=timeout,
        check=False,
    )
    return {
        "exit_code": completed.returncode,
        "stdout": completed.stdout,
    }


Runner = Callable[[Sequence[str]], dict[str, Any]]


def python_cli_argv(command: str, args: Sequence[str]) -> list[str]:
    """Build an argv for invoking the Python CLI with dependencies available."""
    module_args = ["-m", "apm_cli.cli", command, *args]
    try:
        importlib.import_module("click")
    except ModuleNotFoundError:
        uv = shutil.which("uv")
        if uv:
            return [uv, "run", "--extra", "dev", "python", *module_args]
    return [sys.executable, *module_args]


def compare_subcommands(
    go_binary: Path,
    commands: Sequence[str],
    fixtures: Mapping[str, Sequence[str]],
    python_runner: Runner | None = None,
    go_runner: Runner | None = None,
) -> dict[str, Any]:
    python_runner = python_runner or run_command
    go_runner = go_runner or run_command

    details = []
    working = 0
    for command in commands:
        args = fixture_args_for(command, fixtures)
        python_result = python_runner(python_cli_argv(command, args))
        go_result = go_runner([str(go_binary), command, *args])
        matches = (
            python_result["exit_code"] == go_result["exit_code"]
            and python_result["stdout"] == go_result["stdout"]
        )
        if matches:
            working += 1
        details.append(
            {
                "command": command,
                "args": args,
                "matches": matches,
                "python_exit_code": python_result["exit_code"],
                "go_exit_code": go_result["exit_code"],
            }
        )

    total = len(commands)
    pct = round((working / total) * 100, 2) if total else 0.0
    return {
        "working_subcommands_pct": pct,
        "working_subcommands": working,
        "total_subcommands": total,
        "details": details,
    }


def main(argv: Sequence[str] | None = None) -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("go_binary", type=Path, help="Path to the compiled Go apm binary")
    parser.add_argument(
        "--fixtures",
        type=Path,
        default=None,
        help="Optional JSON map of subcommand names to parity fixture arguments",
    )
    args = parser.parse_args(argv)

    if not args.go_binary.exists():
        parser.error(f"Go binary not found: {args.go_binary}")

    fixtures = load_fixture_args(args.fixtures)
    result = compare_subcommands(
        args.go_binary,
        enumerate_python_subcommands(),
        fixtures,
    )
    print(json.dumps(result, indent=2, sort_keys=True))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())

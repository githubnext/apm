#!/usr/bin/env python3
"""Extract and check Python behavior contracts for the Go migration.

This script is intentionally source-of-truth oriented: it reads the Python
Click command tree and the existing Python tests, then checks whether each
contract is explicitly covered by Go and CLI-agnostic parity tests.
"""

from __future__ import annotations

import argparse
import ast
import inspect
import json
import os
import sys
from dataclasses import dataclass
from pathlib import Path
from typing import Any

try:
    import click
except ImportError as exc:  # pragma: no cover - exercised in CI setup failures
    raise SystemExit(f"click is required to extract CLI contracts: {exc}") from exc

try:
    import yaml
except ImportError as exc:  # pragma: no cover - exercised in CI setup failures
    raise SystemExit(f"PyYAML is required to check coverage contracts: {exc}") from exc


ROOT = Path(__file__).resolve().parents[2]
SRC = ROOT / "src"
TESTS = ROOT / "tests"


@dataclass(frozen=True)
class Finding:
    code: str
    message: str
    contract: str


def _rel(path: Path) -> str:
    return path.resolve().relative_to(ROOT).as_posix()


def _source_location(obj: object | None) -> dict[str, Any] | None:
    if obj is None:
        return None
    obj = inspect.unwrap(obj)
    try:
        source = Path(inspect.getsourcefile(obj) or "")
        line = inspect.getsourcelines(obj)[1]
    except (OSError, TypeError):
        return None
    if not source:
        return None
    try:
        return {"file": _rel(source), "line": line}
    except ValueError:
        return {"file": str(source), "line": line}


def _param_contract(param: click.Parameter) -> dict[str, Any]:
    base: dict[str, Any] = {
        "name": param.name,
        "required": bool(getattr(param, "required", False)),
        "multiple": bool(getattr(param, "multiple", False)),
        "nargs": getattr(param, "nargs", None),
        "type": param.__class__.__name__,
    }
    if isinstance(param, click.Option):
        base.update(
            {
                "opts": list(param.opts),
                "secondary_opts": list(param.secondary_opts),
                "help": param.help or "",
                "default": repr(param.default),
                "is_flag": bool(param.is_flag),
                "flag_value": repr(param.flag_value),
            }
        )
    elif isinstance(param, click.Argument):
        base.update({"human_readable_name": param.human_readable_name})
    return base


def _command_id(path: tuple[str, ...]) -> str:
    return "apm" if not path else "apm " + " ".join(path)


def _iter_click_commands(
    command: click.Command, path: tuple[str, ...] = ()
) -> list[dict[str, Any]]:
    callback = getattr(command, "callback", None)
    contract = {
        "id": _command_id(path),
        "path": list(path),
        "name": command.name or "apm",
        "hidden": bool(getattr(command, "hidden", False)),
        "deprecated": bool(getattr(command, "deprecated", False)),
        "type": command.__class__.__name__,
        "help": command.help or "",
        "short_help": command.short_help or "",
        "params": [_param_contract(param) for param in command.params],
        "source": _source_location(callback),
        "subcommands": [],
    }

    contracts = [contract]
    if isinstance(command, click.Group):
        # Use the raw commands mapping so hidden aliases are included too.
        for name in sorted(command.commands):
            child = command.commands[name]
            contract["subcommands"].append(name)
            contracts.extend(_iter_click_commands(child, (*path, name)))
    return contracts


def extract_click_contracts() -> list[dict[str, Any]]:
    sys.path.insert(0, str(SRC))
    from apm_cli.cli import cli

    return _iter_click_commands(cli)


def _literal_string_sequence(node: ast.AST) -> list[str] | None:
    if not isinstance(node, (ast.List, ast.Tuple)):
        return None
    values: list[str] = []
    for elt in node.elts:
        if isinstance(elt, ast.Constant) and isinstance(elt.value, str):
            values.append(elt.value)
        else:
            return None
    return values


def _call_name(node: ast.AST) -> str:
    if isinstance(node, ast.Name):
        return node.id
    if isinstance(node, ast.Attribute):
        parent = _call_name(node.value)
        return f"{parent}.{node.attr}" if parent else node.attr
    return ""


def _extract_cli_invocations(node: ast.AST) -> list[dict[str, Any]]:
    invocations: list[dict[str, Any]] = []
    for child in ast.walk(node):
        if not isinstance(child, ast.Call):
            continue
        call = _call_name(child.func)
        if not call.endswith(".invoke") and call != "invoke":
            continue
        args: list[str] | None = None
        if len(child.args) >= 2:
            args = _literal_string_sequence(child.args[1])
        for keyword in child.keywords:
            if keyword.arg in {"args", "cli_args"}:
                args = _literal_string_sequence(keyword.value)
        invocations.append({"call": call, "args": args})
    return invocations


def _extract_import_roots(tree: ast.Module) -> list[str]:
    roots: set[str] = set()
    for node in tree.body:
        if isinstance(node, ast.Import):
            for alias in node.names:
                if alias.name.startswith("apm_cli"):
                    roots.add(alias.name)
        elif isinstance(node, ast.ImportFrom) and node.module:
            if node.module.startswith("apm_cli"):
                roots.add(node.module)
    return sorted(roots)


def _parametrize_count(node: ast.AST) -> int:
    count = 1
    decorators = getattr(node, "decorator_list", [])
    for decorator in decorators:
        if not isinstance(decorator, ast.Call):
            continue
        if not _call_name(decorator.func).endswith("parametrize"):
            continue
        if len(decorator.args) < 2:
            continue
        values = decorator.args[1]
        if isinstance(values, (ast.List, ast.Tuple)):
            count *= max(1, len(values.elts))
    return count


def _test_contract(
    file: Path,
    tree: ast.Module,
    node: ast.FunctionDef | ast.AsyncFunctionDef,
    class_name: str | None = None,
) -> dict[str, Any]:
    name = node.name if class_name is None else f"{class_name}::{node.name}"
    return {
        "id": f"{_rel(file)}::{name}",
        "file": _rel(file),
        "line": node.lineno,
        "name": node.name,
        "class": class_name,
        "doc": ast.get_docstring(node) or "",
        "import_roots": _extract_import_roots(tree),
        "cli_invocations": _extract_cli_invocations(node),
        "assertions": sum(isinstance(child, ast.Assert) for child in ast.walk(node)),
        "parametrize_cases": _parametrize_count(node),
    }


def extract_test_contracts() -> list[dict[str, Any]]:
    contracts: list[dict[str, Any]] = []
    for file in sorted(TESTS.rglob("test*.py")):
        if "parity" in file.relative_to(TESTS).parts:
            continue
        tree = ast.parse(file.read_text(encoding="utf-8"), filename=str(file))
        for node in tree.body:
            if isinstance(node, (ast.FunctionDef, ast.AsyncFunctionDef)) and node.name.startswith(
                "test_"
            ):
                contracts.append(_test_contract(file, tree, node))
            elif isinstance(node, ast.ClassDef) and node.name.startswith("Test"):
                for item in node.body:
                    if isinstance(
                        item, (ast.FunctionDef, ast.AsyncFunctionDef)
                    ) and item.name.startswith("test_"):
                        contracts.append(_test_contract(file, tree, item, node.name))
    return contracts


def extract_source_contracts() -> list[dict[str, Any]]:
    contracts: list[dict[str, Any]] = []
    for file in sorted(SRC.rglob("*.py")):
        tree = ast.parse(file.read_text(encoding="utf-8"), filename=str(file))
        for node in tree.body:
            if isinstance(node, (ast.FunctionDef, ast.AsyncFunctionDef, ast.ClassDef)):
                if node.name.startswith("_") and node.name not in {"__init__"}:
                    continue
                contracts.append(
                    {
                        "id": f"{_rel(file)}::{node.name}",
                        "file": _rel(file),
                        "line": node.lineno,
                        "name": node.name,
                        "type": node.__class__.__name__,
                        "doc": ast.get_docstring(node) or "",
                    }
                )
    return contracts


def extract_inventory() -> dict[str, Any]:
    commands = extract_click_contracts()
    tests = extract_test_contracts()
    source = extract_source_contracts()
    return {
        "schema_version": 1,
        "root": str(ROOT),
        "summary": {
            "commands": len(commands),
            "public_commands": sum(not c["hidden"] for c in commands),
            "python_tests": len(tests),
            "python_test_cases": sum(t["parametrize_cases"] for t in tests),
            "source_contracts": len(source),
        },
        "commands": commands,
        "tests": tests,
        "source_contracts": source,
    }


def _load_inventory(path: Path | None) -> dict[str, Any]:
    if path is None:
        return extract_inventory()
    return json.loads(path.read_text(encoding="utf-8"))


def _load_coverage(path: Path) -> dict[str, Any]:
    if not path.exists():
        raise SystemExit(f"coverage manifest not found: {path}")
    data = yaml.safe_load(path.read_text(encoding="utf-8")) or {}
    if not isinstance(data, dict):
        raise SystemExit(f"coverage manifest must be a mapping: {path}")
    return data


def _has_tests(entry: dict[str, Any], key: str) -> bool:
    value = entry.get(key)
    return (
        isinstance(value, list)
        and all(isinstance(item, str) and item for item in value)
        and bool(value)
    )


def check_coverage(inventory: dict[str, Any], coverage: dict[str, Any]) -> list[Finding]:
    findings: list[Finding] = []
    command_coverage = coverage.get("commands") or {}
    if not isinstance(command_coverage, dict):
        command_coverage = {}
    test_coverage = (coverage.get("python_tests") or {}).get("covered") or {}
    if not isinstance(test_coverage, dict):
        test_coverage = {}
    obsolete_tests = set((coverage.get("python_tests") or {}).get("obsolete") or [])

    for command in inventory["commands"]:
        command_id = command["id"]
        entry = command_coverage.get(command_id)
        if not isinstance(entry, dict):
            findings.append(
                Finding("missing-command-coverage", "missing command coverage entry", command_id)
            )
            continue
        if not _has_tests(entry, "go_tests"):
            findings.append(
                Finding("missing-command-go-tests", "command lacks mapped Go tests", command_id)
            )
        if not _has_tests(entry, "cli_agnostic_tests"):
            findings.append(
                Finding(
                    "missing-command-cli-tests",
                    "command lacks mapped CLI-agnostic tests",
                    command_id,
                )
            )

    for test in inventory["tests"]:
        test_id = test["id"]
        if test_id in obsolete_tests:
            continue
        entry = test_coverage.get(test_id)
        if not isinstance(entry, dict):
            findings.append(
                Finding("missing-python-test-coverage", "missing Python test mapping", test_id)
            )
            continue
        if not (_has_tests(entry, "go_tests") or _has_tests(entry, "cli_agnostic_tests")):
            findings.append(
                Finding(
                    "missing-python-test-tests",
                    "Python test mapping has neither Go nor CLI-agnostic tests",
                    test_id,
                )
            )

    return findings


def render_summary(inventory: dict[str, Any], findings: list[Finding], *, limit: int = 80) -> str:
    by_code: dict[str, int] = {}
    for finding in findings:
        by_code[finding.code] = by_code.get(finding.code, 0) + 1

    lines = [
        "# Python Behavior Contract Coverage",
        "",
        "## Inventory",
        "",
        f"- Commands: {inventory['summary']['commands']}",
        f"- Public commands: {inventory['summary']['public_commands']}",
        f"- Python tests: {inventory['summary']['python_tests']}",
        f"- Python parametrized test cases: {inventory['summary']['python_test_cases']}",
        f"- Source contracts: {inventory['summary']['source_contracts']}",
        "",
        "## Coverage Findings",
        "",
    ]
    if not findings:
        lines.append("No missing coverage findings.")
        return "\n".join(lines) + "\n"

    for code in sorted(by_code):
        lines.append(f"- {code}: {by_code[code]}")
    lines.extend(
        ["", f"Showing first {min(limit, len(findings))} of {len(findings)} findings:", ""]
    )
    for finding in findings[:limit]:
        lines.append(f"- `{finding.code}` `{finding.contract}`: {finding.message}")
    return "\n".join(lines) + "\n"


def cmd_extract(args: argparse.Namespace) -> int:
    inventory = extract_inventory()
    text = json.dumps(inventory, indent=2, sort_keys=True)
    if args.output:
        Path(args.output).write_text(text + "\n", encoding="utf-8")
    else:
        print(text)
    return 0


def cmd_check(args: argparse.Namespace) -> int:
    inventory = _load_inventory(Path(args.inventory) if args.inventory else None)
    coverage = _load_coverage(Path(args.coverage))
    findings = check_coverage(inventory, coverage)
    summary = render_summary(inventory, findings)
    if args.summary:
        Path(args.summary).write_text(summary, encoding="utf-8")
    print(summary)
    if coverage.get("status") == "intentionally-incomplete":
        # Manifest explicitly declared incomplete; report findings without failing.
        return 0
    return 1 if findings else 0


def main(argv: list[str] | None = None) -> int:
    parser = argparse.ArgumentParser()
    sub = parser.add_subparsers(dest="command", required=True)

    extract = sub.add_parser("extract", help="extract Python behavior/test contracts as JSON")
    extract.add_argument("--output", help="write inventory JSON to this path")
    extract.set_defaults(func=cmd_extract)

    check = sub.add_parser("check", help="check coverage manifest against extracted contracts")
    check.add_argument("--inventory", help="existing inventory JSON path; extracts live if omitted")
    check.add_argument(
        "--coverage",
        default=str(ROOT / "tests" / "parity" / "python_contract_coverage.yml"),
        help="coverage manifest path",
    )
    check.add_argument("--summary", help="write markdown coverage summary to this path")
    check.set_defaults(func=cmd_check)

    args = parser.parse_args(argv)
    os.chdir(ROOT)
    return args.func(args)


if __name__ == "__main__":
    raise SystemExit(main())

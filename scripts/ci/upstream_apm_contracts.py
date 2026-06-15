#!/usr/bin/env python3
"""Check upstream microsoft/apm freshness and Go migration coverage.

The Go migration is not complete just because it matches the Python code that
was present when the migration started. It must also be current with upstream
``microsoft/apm`` and every upstream Python behavior delta must be reviewed and
mapped to real Go tests before Crane can claim completion.
"""

from __future__ import annotations

import argparse
import ast
import json
import os
import re
import subprocess
from dataclasses import dataclass
from pathlib import Path
from typing import Any

import yaml

ROOT = Path(__file__).resolve().parents[2]
GO_TEST_RE = re.compile(r"^func\s+(Test[A-Za-z0-9_]*)\s*\(")


@dataclass(frozen=True)
class Contract:
    id: str
    kind: str


@dataclass(frozen=True)
class Finding:
    code: str
    contract: str
    message: str


@dataclass(frozen=True)
class CheckResult:
    upstream_sha: str
    reviewed_sha: str
    freshness_ok: bool
    contracts_passing: int
    contracts_total: int
    findings: list[Finding]
    freshness_findings: list[str]


def _run_git(
    root: Path, args: list[str], *, check: bool = True
) -> subprocess.CompletedProcess[str]:
    return subprocess.run(  # noqa: S603 - git args are fixed by callers in this CI checker.
        ["git", *args],  # noqa: S607 - git is expected on PATH in CI and local tests.
        cwd=root,
        text=True,
        capture_output=True,
        check=check,
    )


def _git_stdout(root: Path, args: list[str]) -> str:
    return _run_git(root, args).stdout.strip()


def _rev_parse(root: Path, ref: str) -> str:
    return _git_stdout(root, ["rev-parse", ref])


def _has_object(root: Path, ref: str) -> bool:
    return _run_git(root, ["cat-file", "-e", f"{ref}^{{commit}}"], check=False).returncode == 0


def _is_ancestor(root: Path, ancestor: str, descendant: str) -> bool:
    return (
        _run_git(
            root,
            ["merge-base", "--is-ancestor", ancestor, descendant],
            check=False,
        ).returncode
        == 0
    )


def _changed_python_files(root: Path, start: str, end: str) -> list[str]:
    out = _git_stdout(
        root,
        [
            "diff",
            "--name-only",
            "--diff-filter=ACMR",
            f"{start}..{end}",
            "--",
            "src/apm_cli",
            "tests",
        ],
    )
    return sorted(
        path
        for path in out.splitlines()
        if path.endswith(".py") and not path.startswith("tests/parity/")
    )


def _blob_text(root: Path, ref: str, path: str) -> str | None:
    proc = _run_git(root, ["show", f"{ref}:{path}"], check=False)
    if proc.returncode != 0:
        return None
    return proc.stdout


def _source_contracts(path: str, text: str) -> list[Contract]:
    try:
        tree = ast.parse(text, filename=path)
    except SyntaxError:
        return []
    contracts: list[Contract] = []
    for node in tree.body:
        if not isinstance(node, (ast.FunctionDef, ast.AsyncFunctionDef, ast.ClassDef)):
            continue
        if node.name.startswith("_") and node.name != "__init__":
            continue
        contracts.append(Contract(id=f"{path}::{node.name}", kind="source"))
    return contracts


def _test_contracts(path: str, text: str) -> list[Contract]:
    try:
        tree = ast.parse(text, filename=path)
    except SyntaxError:
        return []
    contracts: list[Contract] = []
    for node in tree.body:
        if isinstance(node, (ast.FunctionDef, ast.AsyncFunctionDef)) and node.name.startswith(
            "test_"
        ):
            contracts.append(Contract(id=f"{path}::{node.name}", kind="python_test"))
        elif isinstance(node, ast.ClassDef) and node.name.startswith("Test"):
            for item in node.body:
                if isinstance(
                    item, (ast.FunctionDef, ast.AsyncFunctionDef)
                ) and item.name.startswith("test_"):
                    contracts.append(
                        Contract(id=f"{path}::{node.name}::{item.name}", kind="python_test")
                    )
    return contracts


def changed_contracts(root: Path, start: str, end: str) -> list[Contract]:
    contracts: list[Contract] = []
    for path in _changed_python_files(root, start, end):
        text = _blob_text(root, end, path)
        if text is None:
            continue
        if path.startswith("src/apm_cli/"):
            contracts.extend(_source_contracts(path, text))
        elif path.startswith("tests/"):
            contracts.extend(_test_contracts(path, text))
    unique = {contract.id: contract for contract in contracts}
    return [unique[key] for key in sorted(unique)]


def discover_go_tests(root: Path) -> set[str]:
    tests: set[str] = set()
    for file in sorted((root / "cmd" / "apm").rglob("*_test.go")):
        with file.open(encoding="utf-8") as fh:
            for line in fh:
                match = GO_TEST_RE.match(line.strip())
                if match:
                    tests.add(match.group(1))
    return tests


def _load_yaml(path: Path) -> dict[str, Any]:
    data = yaml.safe_load(path.read_text(encoding="utf-8")) or {}
    if not isinstance(data, dict):
        raise ValueError(f"coverage manifest must be a mapping: {path}")
    return data


def _go_tests_for(entry: object) -> list[str]:
    if not isinstance(entry, dict):
        return []
    tests = entry.get("go_tests")
    if not isinstance(tests, list):
        return []
    return [test for test in tests if isinstance(test, str) and test]


def _coverage_for_contract(range_entry: dict[str, Any], contract: Contract) -> object:
    key = "source_contracts" if contract.kind == "source" else "python_tests"
    entries = range_entry.get(key)
    if not isinstance(entries, dict):
        return None
    return entries.get(contract.id)


def _validate_contracts(
    *,
    contracts: list[Contract],
    range_entry: dict[str, Any],
    go_tests: set[str],
    findings: list[Finding],
) -> int:
    passing = 0
    for contract in contracts:
        entry = _coverage_for_contract(range_entry, contract)
        mapped_tests = _go_tests_for(entry)
        if not mapped_tests:
            findings.append(
                Finding(
                    "missing-upstream-go-tests",
                    contract.id,
                    "upstream Python contract lacks mapped Go tests",
                )
            )
            continue
        unknown = [test for test in mapped_tests if test not in go_tests]
        if unknown:
            findings.append(
                Finding(
                    "unknown-upstream-go-test",
                    contract.id,
                    "mapped Go tests do not exist: " + ", ".join(sorted(unknown)),
                )
            )
            continue
        passing += 1
    return passing


def _range_chain(
    coverage: dict[str, Any],
    *,
    baseline_sha: str,
    reviewed_sha: str,
) -> tuple[list[dict[str, Any]], list[Finding]]:
    ranges = coverage.get("reviewed_ranges") or []
    if not isinstance(ranges, list):
        return [], [Finding("invalid-reviewed-ranges", "reviewed_ranges", "must be a list")]

    by_start: dict[str, dict[str, Any]] = {}
    for index, entry in enumerate(ranges):
        if not isinstance(entry, dict):
            return [], [
                Finding("invalid-reviewed-range", f"reviewed_ranges[{index}]", "must be a mapping")
            ]
        start = entry.get("from")
        end = entry.get("to")
        if not isinstance(start, str) or not isinstance(end, str):
            return [], [
                Finding(
                    "invalid-reviewed-range",
                    f"reviewed_ranges[{index}]",
                    "range must include string 'from' and 'to' SHAs",
                )
            ]
        if start in by_start:
            return [], [Finding("duplicate-reviewed-range", start, "multiple ranges start here")]
        by_start[start] = entry

    chain: list[dict[str, Any]] = []
    cursor = baseline_sha
    seen: set[str] = set()
    while cursor != reviewed_sha:
        if cursor in seen:
            return chain, [Finding("cycle-reviewed-range", cursor, "reviewed range chain loops")]
        seen.add(cursor)
        entry = by_start.get(cursor)
        if entry is None:
            return chain, [
                Finding(
                    "missing-reviewed-range",
                    f"{cursor}..{reviewed_sha}",
                    "reviewed_sha advanced without a chained reviewed_ranges entry",
                )
            ]
        chain.append(entry)
        cursor = str(entry["to"])
    return chain, []


def check_upstream_contracts(
    *,
    root: Path,
    coverage_path: Path,
    upstream_ref: str,
    head_ref: str,
) -> CheckResult:
    coverage = _load_yaml(coverage_path)
    upstream = coverage.get("upstream") or {}
    if not isinstance(upstream, dict):
        raise ValueError("coverage manifest must contain upstream mapping")

    baseline_sha = upstream.get("baseline_sha")
    reviewed_sha = upstream.get("reviewed_sha")
    if not isinstance(baseline_sha, str) or not isinstance(reviewed_sha, str):
        raise ValueError("upstream.baseline_sha and upstream.reviewed_sha are required")

    upstream_sha = _rev_parse(root, upstream_ref)
    freshness_findings: list[str] = []
    if reviewed_sha != upstream_sha:
        freshness_findings.append(
            f"reviewed upstream SHA {reviewed_sha} does not match {upstream_ref} at {upstream_sha}"
        )
    if not _has_object(root, reviewed_sha):
        freshness_findings.append(f"reviewed upstream SHA is not present locally: {reviewed_sha}")
    elif not _is_ancestor(root, reviewed_sha, upstream_sha):
        freshness_findings.append(
            f"reviewed SHA {reviewed_sha} is not reachable from upstream {upstream_sha}"
        )

    go_tests = discover_go_tests(root)
    findings: list[Finding] = []
    passing = 0
    total = 0

    chain, chain_findings = _range_chain(
        coverage,
        baseline_sha=baseline_sha,
        reviewed_sha=reviewed_sha,
    )
    findings.extend(chain_findings)

    for range_entry in chain:
        start = str(range_entry["from"])
        end = str(range_entry["to"])
        contracts = changed_contracts(root, start, end)
        total += len(contracts)
        passing += _validate_contracts(
            contracts=contracts,
            range_entry=range_entry,
            go_tests=go_tests,
            findings=findings,
        )

    if _has_object(root, reviewed_sha):
        pending_contracts = changed_contracts(root, reviewed_sha, upstream_sha)
        total += len(pending_contracts)
        passing += _validate_contracts(
            contracts=pending_contracts,
            range_entry={},
            go_tests=go_tests,
            findings=findings,
        )

    if total == 0:
        total = 1
        passing = 1 if not findings else 0

    return CheckResult(
        upstream_sha=upstream_sha,
        reviewed_sha=reviewed_sha,
        freshness_ok=not freshness_findings,
        contracts_passing=passing,
        contracts_total=total,
        findings=findings,
        freshness_findings=freshness_findings,
    )


def render_summary(result: CheckResult, *, limit: int = 80) -> str:
    lines = [
        "# Upstream APM Contract Coverage",
        "",
        f"- Current upstream SHA: `{result.upstream_sha}`",
        f"- Reviewed upstream SHA: `{result.reviewed_sha}`",
        f"- Freshness: {'pass' if result.freshness_ok else 'fail'}",
        f"- Contract coverage: {result.contracts_passing}/{result.contracts_total}",
        "",
        "## Freshness Findings",
        "",
    ]
    if result.freshness_findings:
        lines.extend(f"- {finding}" for finding in result.freshness_findings)
    else:
        lines.append("No freshness findings.")

    lines.extend(["", "## Contract Findings", ""])
    if result.findings:
        for finding in result.findings[:limit]:
            lines.append(f"- `{finding.code}` `{finding.contract}`: {finding.message}")
        if len(result.findings) > limit:
            lines.append(f"- ... {len(result.findings) - limit} more findings omitted")
    else:
        lines.append("No contract findings.")
    return "\n".join(lines) + "\n"


def _emit_gates(result: CheckResult) -> None:
    print(
        json.dumps(
            {
                "crane": "gate",
                "name": "upstream_freshness",
                "passed": result.freshness_ok,
            },
            sort_keys=True,
        )
    )
    print(
        json.dumps(
            {
                "crane": "gate",
                "name": "upstream_contracts",
                "passing": result.contracts_passing,
                "total": result.contracts_total,
            },
            sort_keys=True,
        )
    )


def cmd_check(args: argparse.Namespace) -> int:
    result = check_upstream_contracts(
        root=Path(args.root).resolve(),
        coverage_path=Path(args.coverage),
        upstream_ref=args.upstream_ref,
        head_ref=args.head_ref,
    )
    _emit_gates(result)
    summary = render_summary(result)
    print(summary)
    if args.summary:
        Path(args.summary).write_text(summary, encoding="utf-8")
    if args.enforce and (
        not result.freshness_ok or result.contracts_passing != result.contracts_total
    ):
        return 1
    return 0


def main(argv: list[str] | None = None) -> int:
    parser = argparse.ArgumentParser()
    sub = parser.add_subparsers(dest="command", required=True)
    check = sub.add_parser("check", help="check upstream APM freshness and coverage")
    check.add_argument("--root", default=str(ROOT), help="repository root")
    check.add_argument(
        "--coverage",
        default=str(ROOT / "tests" / "parity" / "upstream_contract_coverage.yml"),
        help="upstream contract coverage manifest",
    )
    check.add_argument("--upstream-ref", default="upstream/main")
    check.add_argument("--head-ref", default="HEAD")
    check.add_argument("--summary", help="write markdown summary to path")
    check.add_argument("--enforce", action="store_true", help="fail on stale or uncovered upstream")
    check.set_defaults(func=cmd_check)

    args = parser.parse_args(argv)
    os.chdir(Path(args.root).resolve())
    return args.func(args)


if __name__ == "__main__":
    raise SystemExit(main())

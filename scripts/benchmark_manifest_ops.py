#!/usr/bin/env python3
"""Benchmark: manifest-based collision detection and sync operations.

Simulates realistic scale to measure the impact of algorithmic optimizations:
- Optimization 1: Pre-normalized managed_files set (check_collision)
- Optimization 2: Pre-partitioned managed_files (sync_remove_files)
- Optimization 3: Batch empty-parent cleanup (cleanup_empty_parents)
- Optimization 4: Scoped uninstall file set (removed packages only)

Usage:
    python3 scripts/benchmark_manifest_ops.py
    python3 scripts/benchmark_manifest_ops.py --work-dir "$RUNNER_TEMP" --markdown --update-doc docs/src/content/docs/progress/autoloop-go-migration.mdx
"""

import argparse
import shutil
import tempfile
import time
from pathlib import Path

# ---------------------------------------------------------------------------
# Synthetic data — models a repo managing 50 packages × 5 files each = 250
# managed paths (the scale scenario from the review comment).
# ---------------------------------------------------------------------------

PREFIXES = [
    ".github/prompts/",
    ".github/agents/",
    ".claude/agents/",
    ".claude/commands/",
    ".github/skills/",
    ".github/hooks/",
]

PACKAGES = 50
FILES_PER_PACKAGE = 5
INTEGRATOR_TYPES = 6  # prompts, agents-gh, agents-cl, commands, skills, hooks
DOC_START = "{/* benchmark_manifest_ops:start */}"
DOC_END = "{/* benchmark_manifest_ops:end */}"
SCALES = [
    ("Current (10 pkgs x 5 files = 50 paths)", "Current: 10 pkgs, 50 paths", 10, 5),
    ("Growing (50 pkgs x 5 files = 250 paths)", "Growing: 50 pkgs, 250 paths", 50, 5),
    ("Large monorepo (100 pkgs x 20 files = 2000 paths)", "Large monorepo: 100 pkgs, 2,000 paths", 100, 20),
]


def build_managed_files(n_packages: int, files_per_pkg: int) -> set:
    """Generate a synthetic managed_files set."""
    paths = set()
    for i in range(n_packages):
        prefix = PREFIXES[i % len(PREFIXES)]
        for j in range(files_per_pkg):
            paths.add(f"{prefix}pkg-{i}-file-{j}.md")
    return paths


# ---------------------------------------------------------------------------
# OLD: check_collision rebuilds a normalized set on every call
# ---------------------------------------------------------------------------

def check_collision_OLD(rel_path: str, managed_files: set) -> bool:
    """Original O(M) per call — rebuilds normalized set."""
    return rel_path.replace("\\", "/") not in {p.replace("\\", "/") for p in managed_files}


# ---------------------------------------------------------------------------
# NEW: managed_files is pre-normalized — O(1) amortized lookup
# ---------------------------------------------------------------------------

def normalize_managed_files(managed_files: set) -> set:
    return {p.replace("\\", "/") for p in managed_files}


def check_collision_NEW(rel_path: str, managed_files_normalized: set) -> bool:
    """Optimized O(1) lookup against pre-normalized set."""
    return rel_path.replace("\\", "/") not in managed_files_normalized


# ---------------------------------------------------------------------------
# OLD: sync_remove_files scans all paths with startswith for each integrator
# ---------------------------------------------------------------------------

def sync_remove_old(managed_files: set, prefix: str) -> int:
    """Original: iterate full set, filter by prefix."""
    count = 0
    for rel_path in managed_files:
        normalized = rel_path.replace("\\", "/")
        if not normalized.startswith(prefix) or ".." in rel_path:
            continue
        count += 1  # simulate file operation
    return count


# ---------------------------------------------------------------------------
# NEW: pre-partitioned — each integrator receives only its own paths
# ---------------------------------------------------------------------------

def partition_managed_files(managed_files: set) -> dict:
    """Single-pass partition by prefix."""
    buckets = {p: set() for p in PREFIXES}
    for p in managed_files:
        for prefix in PREFIXES:
            if p.startswith(prefix):
                buckets[prefix].add(p)
                break
    return buckets


def sync_remove_new(bucket: set) -> int:
    """Optimized: iterate only matching paths (no prefix check needed)."""
    count = 0
    for rel_path in bucket:
        if ".." in rel_path:
            continue
        count += 1
    return count


# ---------------------------------------------------------------------------
# Benchmarks
# ---------------------------------------------------------------------------

def timeit(fn, *args, iterations: int = 1000) -> float:
    """Return total time in ms for *iterations* calls."""
    start = time.perf_counter()
    for _ in range(iterations):
        fn(*args)
    return (time.perf_counter() - start) * 1000


def _make_temp_dir(work_dir: Path | None) -> Path:
    """Create benchmark scratch space inside a caller-approved directory."""
    if work_dir is not None:
        work_dir.mkdir(parents=True, exist_ok=True)
    return Path(
        tempfile.mkdtemp(
            prefix="apm-manifest-ops-",
            dir=str(work_dir) if work_dir is not None else None,
        )
    )


def _format_speedup(speedup: float) -> str:
    if speedup == float("inf"):
        return "inf x"
    return f"{speedup:.1f}x"


def run_benchmarks(work_dir: Path | None = None, emit_text: bool = True) -> list[dict[str, str]]:
    results = []
    if emit_text:
        print("=" * 72)
        print("APM Manifest Operations Benchmark")
        print("=" * 72)

    for scale_label, markdown_label, n_pkgs, n_files in SCALES:
        managed = build_managed_files(n_pkgs, n_files)
        M = len(managed)
        if emit_text:
            print(f"\n{'─' * 72}")
            print(f"Scale: {scale_label}  (M={M})")
            print(f"{'─' * 72}")

        # -- Benchmark 1: check_collision ----------------------------------
        #
        # Simulate: P=n_pkgs packages × F=n_files files × I=6 integrators
        # Each call does one collision check.
        calls = n_pkgs * n_files * INTEGRATOR_TYPES
        test_path = ".github/prompts/pkg-0-file-0.md"

        old_time = timeit(check_collision_OLD, test_path, managed, iterations=calls)
        normalized = normalize_managed_files(managed)
        norm_time = timeit(normalize_managed_files, managed, iterations=1)
        new_time = norm_time + timeit(check_collision_NEW, test_path, normalized, iterations=calls)

        speedup = old_time / new_time if new_time > 0 else float("inf")
        if emit_text:
            print(f"\n  check_collision ({calls:,} calls):")
            print(f"    OLD (set rebuild per call):  {old_time:>8.2f} ms")
            print(f"    NEW (pre-normalized O(1)):   {new_time:>8.2f} ms")
            print(f"    Speedup:                     {speedup:>8.1f}x")

        # -- Benchmark 2: sync_remove_files --------------------------------
        #
        # Simulate uninstall: 6 integrators, each scanning full M paths.
        sync_prefixes = PREFIXES
        iters = 100

        # OLD: 6 calls, each scanning full set
        old_sync = 0.0
        for _ in range(iters):
            t0 = time.perf_counter()
            for prefix in sync_prefixes:
                sync_remove_old(managed, prefix)
            old_sync += (time.perf_counter() - t0) * 1000

        # NEW: 1 partition pass + 6 calls over subset
        new_sync = 0.0
        for _ in range(iters):
            t0 = time.perf_counter()
            buckets = partition_managed_files(managed)
            for prefix in sync_prefixes:
                sync_remove_new(buckets[prefix])
            new_sync += (time.perf_counter() - t0) * 1000

        speedup2 = old_sync / new_sync if new_sync > 0 else float("inf")
        if emit_text:
            print(f"\n  sync_remove_files ({iters} uninstall cycles x 6 integrators):")
            print(f"    OLD (6x full-set scan):      {old_sync:>8.2f} ms")
            print(f"    NEW (pre-partitioned):       {new_sync:>8.2f} ms")
            print(f"    Speedup:                     {speedup2:>8.1f}x")

        # -- Benchmark 3: empty-parent cleanup ----------------------------
        #
        # Create a real temp directory tree and compare per-file parent
        # walk-up vs. batch bottom-up cleanup.
        depth = 4  # nesting depth for hook-style paths
        n_deleted = n_pkgs * 2  # number of files to delete

        def _make_tree(base: Path, count: int, nest: int):
            """Create *count* files nested *nest* levels deep."""
            paths = []
            for i in range(count):
                parts = [f"d{i % 6}"] + [f"sub{j}" for j in range(nest - 1)]
                d = base.joinpath(*parts)
                d.mkdir(parents=True, exist_ok=True)
                f = d / f"file-{i}.md"
                f.write_text("")
                paths.append(f)
            return paths

        # OLD: per-file walk-up
        tmp_old = _make_temp_dir(work_dir)
        try:
            files_old = _make_tree(tmp_old, n_deleted, depth)
            for f in files_old:
                f.unlink()
            t0 = time.perf_counter()
            for f in files_old:
                parent = f.parent
                while parent != tmp_old and parent.exists():
                    try:
                        if not any(parent.iterdir()):
                            parent.rmdir()
                            parent = parent.parent
                        else:
                            break
                    except OSError:
                        break
            old_parent_ms = (time.perf_counter() - t0) * 1000
        finally:
            shutil.rmtree(tmp_old, ignore_errors=True)

        # NEW: batch bottom-up
        tmp_new = _make_temp_dir(work_dir)
        try:
            files_new = _make_tree(tmp_new, n_deleted, depth)
            for f in files_new:
                f.unlink()
            t0 = time.perf_counter()
            # Inline the algorithm (same as BaseIntegrator.cleanup_empty_parents)
            candidates = set()
            for p in files_new:
                parent = p.parent
                while parent != tmp_new:
                    candidates.add(parent)
                    parent = parent.parent
            for d in sorted(candidates, key=lambda p: len(p.parts), reverse=True):
                try:
                    if d.exists() and not any(d.iterdir()):
                        d.rmdir()
                except OSError:
                    pass
            new_parent_ms = (time.perf_counter() - t0) * 1000
        finally:
            shutil.rmtree(tmp_new, ignore_errors=True)

        speedup3 = old_parent_ms / new_parent_ms if new_parent_ms > 0 else float("inf")
        if emit_text:
            print(f"\n  cleanup_empty_parents ({n_deleted} deleted files, depth={depth}):")
            print(f"    OLD (per-file walk-up):      {old_parent_ms:>8.2f} ms")
            print(f"    NEW (batch bottom-up):       {new_parent_ms:>8.2f} ms")
            print(f"    Speedup:                     {speedup3:>8.1f}x")

        # -- Benchmark 4: scoped vs. union-all deployed files --------------
        #
        # Simulate uninstalling 5 out of n_pkgs packages. Compare iterating
        # all M paths vs. only the removed packages' paths.
        removed_count = min(5, n_pkgs)
        removed_pkgs = set(range(removed_count))

        # Build per-package deployed_files
        pkg_files: dict = {}
        for i in range(n_pkgs):
            prefix = PREFIXES[i % len(PREFIXES)]
            pkg_files[i] = {f"{prefix}pkg-{i}-file-{j}.md" for j in range(n_files)}

        iters4 = 1000

        # OLD: union ALL
        all_files = set()
        for v in pkg_files.values():
            all_files.update(v)
        t0 = time.perf_counter()
        for _ in range(iters4):
            for prefix in PREFIXES:
                _ = [p for p in all_files if p.startswith(prefix)]
        old_scope_ms = (time.perf_counter() - t0) * 1000

        # NEW: union only removed
        removed_files = set()
        for i in removed_pkgs:
            removed_files.update(pkg_files[i])
        t0 = time.perf_counter()
        for _ in range(iters4):
            for prefix in PREFIXES:
                _ = [p for p in removed_files if p.startswith(prefix)]
        new_scope_ms = (time.perf_counter() - t0) * 1000

        speedup4 = old_scope_ms / new_scope_ms if new_scope_ms > 0 else float("inf")
        if emit_text:
            print(f"\n  scoped uninstall set (removing {removed_count}/{n_pkgs} pkgs, {iters4} cycles):")
            print(f"    OLD (union ALL {len(all_files)} paths):     {old_scope_ms:>8.2f} ms")
            print(f"    NEW (union removed {len(removed_files)} paths): {new_scope_ms:>8.2f} ms")
            print(f"    Speedup:                     {speedup4:>8.1f}x")

        results.append(
            {
                "scale": markdown_label,
                "check_collision": _format_speedup(speedup),
                "sync_remove_files": _format_speedup(speedup2),
                "cleanup_empty_parents": _format_speedup(speedup3),
                "scoped_uninstall": _format_speedup(speedup4),
            }
        )

    if emit_text:
        print(f"\n{'=' * 72}")
        print("Done.")

    return results


def render_markdown(results: list[dict[str, str]]) -> str:
    lines = [
        DOC_START,
        "The table below is generated by `scripts/benchmark_manifest_ops.py` during the docs build.",
        "",
        "| Scale | `check_collision` speedup | `sync_remove_files` speedup | `cleanup_empty_parents` speedup | Scoped uninstall speedup |",
        "|---|---:|---:|---:|---:|",
    ]
    for result in results:
        lines.append(
            f"| {result['scale']} | {result['check_collision']} | {result['sync_remove_files']} | "
            f"{result['cleanup_empty_parents']} | {result['scoped_uninstall']} |"
        )
    lines.extend(
        [
            "",
            "`cleanup_empty_parents` may show a small regression at low deleted-file counts because the batch bottom-up algorithm has higher constant overhead than the legacy per-file walk-up. This is expected and acceptable given the gains on the other three operations.",
            DOC_END,
        ]
    )
    return "\n".join(lines)


def update_doc(path: Path, markdown: str) -> None:
    content = path.read_text()
    if DOC_START not in content:
        raise ValueError(f"Could not update {path}: missing {DOC_START} marker")
    start = content.index(DOC_START)
    if DOC_END not in content[start:]:
        raise ValueError(f"Could not update {path}: missing {DOC_END} marker")
    end = content.index(DOC_END, start) + len(DOC_END)
    path.write_text(f"{content[:start]}{markdown}{content[end:]}")


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument(
        "--work-dir",
        type=Path,
        default=None,
        help="Writable directory for temporary benchmark trees (defaults to Python's temp dir).",
    )
    parser.add_argument(
        "--markdown",
        action="store_true",
        help="Emit the docs markdown benchmark table instead of the console report.",
    )
    parser.add_argument(
        "--update-doc",
        type=Path,
        default=None,
        help="Replace the generated benchmark block in the given docs page.",
    )
    return parser.parse_args()


if __name__ == "__main__":
    args = parse_args()
    benchmark_results = run_benchmarks(work_dir=args.work_dir, emit_text=not args.markdown)
    if args.markdown or args.update_doc is not None:
        benchmark_markdown = render_markdown(benchmark_results)
        if args.update_doc is not None:
            update_doc(args.update_doc, benchmark_markdown)
            print(f"Updated {args.update_doc}")
        else:
            print(benchmark_markdown)

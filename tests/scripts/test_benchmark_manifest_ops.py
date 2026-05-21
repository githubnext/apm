import importlib.util
from pathlib import Path


def load_benchmark_module():
    script_path = Path(__file__).resolve().parents[2] / "scripts" / "benchmark_manifest_ops.py"
    spec = importlib.util.spec_from_file_location("benchmark_manifest_ops", script_path)
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return module


def test_run_benchmarks_uses_configured_work_dir(tmp_path, monkeypatch):
    benchmark = load_benchmark_module()
    created_paths = []
    original_mkdtemp = benchmark.tempfile.mkdtemp

    def recording_mkdtemp(*args, **kwargs):
        path = Path(original_mkdtemp(*args, **kwargs))
        created_paths.append(path)
        return str(path)

    monkeypatch.setattr(benchmark.tempfile, "mkdtemp", recording_mkdtemp)

    results = benchmark.run_benchmarks(work_dir=tmp_path, emit_text=False)

    assert [result["scale"] for result in results] == [
        "Current: 10 pkgs, 50 paths",
        "Growing: 50 pkgs, 250 paths",
        "Large monorepo: 100 pkgs, 2,000 paths",
    ]
    assert created_paths
    paths_outside_work_dir = [path for path in created_paths if not path.is_relative_to(tmp_path)]
    assert not paths_outside_work_dir, f"Paths outside work dir: {paths_outside_work_dir}"
    assert not any(tmp_path.iterdir())

    markdown = benchmark.render_markdown(results)
    assert benchmark.DOC_START in markdown
    assert benchmark.DOC_END in markdown
    assert "| Current: 10 pkgs, 50 paths |" in markdown


def test_update_doc_replaces_generated_block(tmp_path):
    benchmark = load_benchmark_module()
    doc_path = tmp_path / "page.mdx"
    doc_path.write_text(f"before\n{benchmark.DOC_START}\nold content\n{benchmark.DOC_END}\nafter\n")

    benchmark.update_doc(doc_path, f"{benchmark.DOC_START}\nnew content\n{benchmark.DOC_END}")

    assert doc_path.read_text() == (
        f"before\n{benchmark.DOC_START}\nnew content\n{benchmark.DOC_END}\nafter\n"
    )

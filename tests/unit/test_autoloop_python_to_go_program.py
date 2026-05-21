import importlib.util
import unittest
from pathlib import Path

REPO_ROOT = Path(__file__).resolve().parents[2]
PROGRAM_FILE = REPO_ROOT / ".autoloop/programs/python-to-go-migration/program.md"
SCHEDULER_FILE = REPO_ROOT / ".github/workflows/scripts/autoloop_scheduler.py"
PARITY_SCRIPT = REPO_ROOT / "scripts/cli_parity_check.py"


def _load_scheduler():
    spec = importlib.util.spec_from_file_location("autoloop_scheduler", SCHEDULER_FILE)
    module = importlib.util.module_from_spec(spec)
    assert spec.loader is not None
    spec.loader.exec_module(module)
    return module


def _load_parity_script():
    spec = importlib.util.spec_from_file_location("cli_parity_check", PARITY_SCRIPT)
    module = importlib.util.module_from_spec(spec)
    assert spec.loader is not None
    spec.loader.exec_module(module)
    return module


class PythonToGoAutoloopProgramTest(unittest.TestCase):
    def test_program_uses_observable_cli_parity_metric(self):
        body = PROGRAM_FILE.read_text()

        self.assertIn("scripts/cli_parity_check.py", body)
        self.assertIn("working_subcommands_pct", body)
        self.assertIn("go build -o /tmp/apm-go ./cmd/apm", body)
        self.assertNotIn("migrated_python_lines", body)
        self.assertNotIn("python_lines_migrated_pct", body)

    def test_program_requires_vertical_cli_wiring_and_rejects_padding(self):
        body = PROGRAM_FILE.read_text()

        self.assertIn("cmd/apm/main.go", body)
        self.assertIn("integration test that runs the compiled binary", body)
        self.assertIn("Every 5th iteration", body)
        self.assertIn("complete user-facing command", body)
        self.assertIn("Do NOT add additional test files", body)
        self.assertIn(">= 80%", body)
        self.assertIn("regenerate `benchmarks/migration-status.json` from source data", body)

    def test_scheduler_prefers_checked_in_program_over_issue_duplicate(self):
        scheduler = _load_scheduler()
        files = [
            ".autoloop/programs/python-to-go-migration/program.md",
            "/tmp/gh-aw/issue-programs/python-to-go-migration.md",
        ]

        self.assertEqual(scheduler.dedupe_program_files(files), files[:1])


class CliParityCheckTest(unittest.TestCase):
    def test_compare_subcommands_scores_exit_code_and_stdout_matches(self):
        parity = _load_parity_script()

        def python_runner(argv):
            command = next(token for token in ("ok", "mismatch") if token in argv)
            return {"exit_code": 0, "stdout": f"{command}\n"}

        def go_runner(argv):
            command = argv[1]
            stdout = f"{command}\n" if command == "ok" else "different\n"
            return {"exit_code": 0, "stdout": stdout}

        result = parity.compare_subcommands(
            Path("/tmp/apm-go"),
            ["ok", "mismatch"],
            {},
            python_runner=python_runner,
            go_runner=go_runner,
        )

        self.assertEqual(result["working_subcommands_pct"], 50.0)
        self.assertEqual(result["working_subcommands"], 1)
        self.assertEqual(result["total_subcommands"], 2)
        self.assertEqual([detail["matches"] for detail in result["details"]], [True, False])

    def test_fixture_args_default_to_help(self):
        parity = _load_parity_script()

        self.assertEqual(parity.fixture_args_for("install", {}), ["--help"])
        self.assertEqual(
            parity.fixture_args_for("install", {"install": ["--dry-run", "pkg"]}),
            ["--dry-run", "pkg"],
        )

    def test_static_enumerator_finds_python_cli_subcommands(self):
        parity = _load_parity_script()

        commands = parity.enumerate_python_subcommands_static()

        self.assertIn("audit", commands)
        self.assertIn("install", commands)
        self.assertIn("uninstall", commands)
        self.assertNotIn("info", commands)


if __name__ == "__main__":
    unittest.main()

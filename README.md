# APM – Agent Package Manager

Experimental go migration of APM.

## Go CLI migration status

The Python-to-Go CLI migration is now landed but still intentionally keeps both
implementations in the tree:

- The Python CLI remains the reference implementation and parity oracle.
- The Go CLI lives under `cmd/apm` and can be built as a local `apm-go` binary.
- The Crane migration PR was merged in
  [githubnext/apm#91](https://github.com/githubnext/apm/pull/91), and the
  migration issue is marked `crane-completed` in
  [githubnext/apm#78](https://github.com/githubnext/apm/issues/78).
- The restored migration workflow verifies Python-vs-Go parity with
  `APM_PYTHON_BIN` set. Latest verified parity evidence reports
  `migration_score=1`, `progress=1`, and `706/706` parity tests passing.
- The benchmark workflow follow-up in
  [githubnext/apm#93](https://github.com/githubnext/apm/pull/93) uploads a
  `migration-benchmark-evidence` artifact with Python-vs-Go CLI timings.

This means Python unit tests and Go parity tests pass for the migration gate.
That gate is not the same thing as claiming all historical Python integration,
live-service, and benchmark coverage is now green for every workflow.

### Build and run the Go CLI locally

From the repository root:

```bash
go build -o ./dist/apm-go ./cmd/apm
```

Then try the local binary:

```bash
./dist/apm-go --help
./dist/apm-go init --yes
```

### Run the definitive parity check locally

When the Python CLI is installed in the project virtual environment, run:

```bash
uv sync --extra dev
export APM_PYTHON_BIN="$PWD/.venv/bin/apm"
go test ./...
go test -json ./... | go run .crane/scripts/score.go
```

`APM_PYTHON_BIN` is required for the hard Python-vs-Go parity gate; without it,
Go-only tests are not completion evidence.

### Run the Actions parity and benchmark workflow

Maintainers can dispatch the migration workflow manually:

```bash
gh workflow run migration-ci.yml --repo githubnext/apm --ref main
```

That default manual run collects parity and benchmark evidence without treating
known migration gaps as a CI failure. To run the deterministic hard completion
gate, opt in explicitly:

```bash
gh workflow run migration-ci.yml --repo githubnext/apm --ref main -f enforce_completion=true
```

After it runs, open the **Migration Benchmarks** job summary for the timing
table. The same run uploads the `migration-benchmark-evidence` artifact with
JSON and Markdown copies of the benchmark data. In the benchmark table, the
`Go/Python` ratio is the Go median duration divided by the Python median
duration: values below `1.00x` mean Go is faster. The benchmark includes
fixture-backed commands that read, write, execute, or fail against realistic APM
project state: `apm.yml`, `apm.lock.yaml`, installed `apm_modules`, local
`.apm` primitives, target directories, deployed prompt files, and sample source
files.

# Python Behavior Contract Parity

The Go migration is not complete until every behavior contract from the
original Python CLI is covered by Go tests and by binary-level, CLI-agnostic
tests.

`scripts/ci/python_behavior_contracts.py` extracts three inventories:

- the Python Click command tree, including options, arguments, hidden aliases,
  help text, and callback source locations;
- every Python test function and parametrized test case;
- public Python source callables that describe implementation behavior.

`python_contract_coverage.yml` is the audited mapping from those extracted
contracts to parity evidence. The completion scorer must not reach
`migration_score = 1.0` while any extracted command or Python test lacks mapped
coverage.

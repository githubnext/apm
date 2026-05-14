---
description: |
  Evergreen -- keeps pull requests healthy by automatically fixing merge conflicts
  and failing CI checks. Runs on a short schedule, deterministically selects one
  PR per run, and gives up after 5 attempts that don't improve the same repo state.

on:
  schedule: every 5m
  workflow_dispatch:
    inputs:
      pr_number:
        description: "Fix a specific PR by number (bypasses scheduling)"
        required: false
        type: string

permissions: read-all

timeout-minutes: 360

network:
  allowed:
  - defaults
  - node
  - go
  - releaseassets.githubusercontent.com

safe-outputs:
  push-to-pull-request-branch:
    target: "*"
    max: 3
    protected-files: allowed
  add-comment:
    max: 3
    target: "*"

checkout:
  fetch: ["*"]
  fetch-depth: 0

tools:
  github:
    toolsets: [repos, pull_requests, issues, actions]
  bash: true
  repo-memory:
    branch-name: memory/evergreen
    file-glob: ["*.md"]

imports:
  - shared/reporting.md

steps:
  - name: Find a PR that needs attention
    id: find-pr
    env:
      GITHUB_TOKEN: ${{ github.token }}
      GH_AW_CI_TRIGGER_TOKEN: ${{ secrets.GH_AW_CI_TRIGGER_TOKEN }}
      GITHUB_REPOSITORY: ${{ github.repository }}
      FORCED_PR: ${{ github.event.inputs.pr_number }}
    run: |
      python3 - << 'PYEOF'
      import os, json, re, subprocess, sys
      import urllib.request, urllib.error

      def emit_selected_output(pr_number):
          """Expose `selected` as a step output for workflow gating.
          Empty string means no PR needs attention; otherwise the PR number."""
          gh_output = os.environ.get("GITHUB_OUTPUT")
          if gh_output:
              with open(gh_output, "a") as f:
                  f.write(f"selected={'' if pr_number is None else pr_number}\n")

      token = os.environ.get("GITHUB_TOKEN", "")
      repo = os.environ.get("GITHUB_REPOSITORY", "")
      forced_pr = os.environ.get("FORCED_PR", "").strip()

      repo_memory_dir = "/tmp/gh-aw/repo-memory/evergreen"
      output_file = "/tmp/gh-aw/evergreen.json"
      os.makedirs("/tmp/gh-aw", exist_ok=True)

      MAX_ATTEMPTS = 5

      def api_get(url):
          """Make an authenticated GET request to the GitHub API."""
          req = urllib.request.Request(url, headers={
              "Authorization": f"token {token}",
              "Accept": "application/vnd.github.v3+json",
          })
          with urllib.request.urlopen(req, timeout=30) as resp:
              return json.loads(resp.read().decode())

      def get_all_open_prs():
          """Fetch all open PRs, paginated."""
          prs = []
          page = 1
          while True:
              url = f"https://api.github.com/repos/{repo}/pulls?state=open&per_page=100&page={page}&sort=number&direction=asc"
              batch = api_get(url)
              if not batch:
                  break
              prs.extend(batch)
              if len(batch) < 100:
                  break
              page += 1
          return prs

      def get_check_status(pr):
          """Get combined CI check status for a PR's head commit."""
          head_sha = pr["head"]["sha"]
          url = f"https://api.github.com/repos/{repo}/commits/{head_sha}/status"
          try:
              status = api_get(url)
              return status.get("state", "unknown")
          except Exception as e:
              print(f"  Warning: could not fetch status for PR #{pr['number']}: {e}")
              return "unknown"

      def get_check_runs(pr):
          """Get check runs for a PR's head commit."""
          head_sha = pr["head"]["sha"]
          url = f"https://api.github.com/repos/{repo}/commits/{head_sha}/check-runs"
          try:
              data = api_get(url)
              return data.get("check_runs", [])
          except Exception as e:
              print(f"  Warning: could not fetch check runs for PR #{pr['number']}: {e}")
              return []

      def read_attempt_state(pr_number):
          """Read attempt tracking state from repo-memory."""
          state_file = os.path.join(repo_memory_dir, f"pr-{pr_number}.md")
          if not os.path.isfile(state_file):
              return {"attempts": 0, "head_sha": None}
          with open(state_file, encoding="utf-8") as f:
              content = f.read()
          state = {"attempts": 0, "head_sha": None}
          m = re.search(r'\|\s*head_sha\s*\|\s*(\S+)\s*\|', content)
          if m:
              state["head_sha"] = m.group(1)
          m = re.search(r'\|\s*attempts\s*\|\s*(\d+)\s*\|', content)
          if m:
              state["attempts"] = int(m.group(1))
          return state

      def get_commit_date(sha):
          """Return the committer date (ISO 8601) for a given commit SHA, or None."""
          url = f"https://api.github.com/repos/{repo}/commits/{sha}"
          try:
              data = api_get(url)
              return data.get("commit", {}).get("committer", {}).get("date")
          except Exception as e:
              print(f"  Warning: could not fetch commit {sha[:12]}: {e}")
              return None

      def is_autoloop_pr(pr):
          """Return True if the PR is from an autoloop branch.
          Branch name is the primary gate (labels can be added by anyone on
          public repos); the `autoloop` label is just an additional signal."""
          head_ref = pr.get("head", {}).get("ref", "") or ""
          return head_ref.startswith("autoloop/")

      def get_behind_by(pr):
          """Return how many commits the PR base branch is ahead of the PR head."""
          base = pr["base"]["ref"]
          head_sha = pr["head"]["sha"]
          url = f"https://api.github.com/repos/{repo}/compare/{base}...{head_sha}"
          try:
              data = api_get(url)
              return int(data.get("behind_by", 0) or 0)
          except Exception as e:
              print(f"  Warning: could not fetch compare for PR #{pr['number']}: {e}")
              return 0

      def trigger_ci_workflow(branch):
          """Trigger this repository's ci.yml for `branch` by pushing an empty
          commit with GH_AW_CI_TRIGGER_TOKEN. Because the push is attributed to
          a real user instead of GITHUB_TOKEN, GitHub emits the PR synchronize
          event that starts the pull_request CI workflow; this is separate from
          the manual workflow_dispatch trigger available in ci.yml."""
          ci_token = os.environ.get("GH_AW_CI_TRIGGER_TOKEN", "") or token
          try:
              # Get current HEAD SHA
              url = f"https://api.github.com/repos/{repo}/git/ref/heads/{branch}"
              req = urllib.request.Request(url, headers={
                  "Authorization": f"token {ci_token}",
                  "Accept": "application/vnd.github.v3+json",
              })
              with urllib.request.urlopen(req, timeout=30) as resp:
                  head_sha = json.loads(resp.read().decode())["object"]["sha"]

              # Create an empty commit on top of HEAD
              url = f"https://api.github.com/repos/{repo}/git/commits"
              payload = json.dumps({
                  "message": "chore: trigger CI [evergreen]",
                  "tree": json.loads(urllib.request.urlopen(
                      urllib.request.Request(
                          f"https://api.github.com/repos/{repo}/git/commits/{head_sha}",
                          headers={"Authorization": f"token {ci_token}", "Accept": "application/vnd.github.v3+json"},
                      ), timeout=30
                  ).read().decode())["tree"]["sha"],
                  "parents": [head_sha],
              }).encode()
              req = urllib.request.Request(url, data=payload, method="POST", headers={
                  "Authorization": f"token {ci_token}",
                  "Accept": "application/vnd.github.v3+json",
                  "Content-Type": "application/json",
              })
              with urllib.request.urlopen(req, timeout=30) as resp:
                  new_sha = json.loads(resp.read().decode())["sha"]

              # Update the branch ref to the new commit
              url = f"https://api.github.com/repos/{repo}/git/refs/heads/{branch}"
              payload = json.dumps({"sha": new_sha}).encode()
              req = urllib.request.Request(url, data=payload, method="PATCH", headers={
                  "Authorization": f"token {ci_token}",
                  "Accept": "application/vnd.github.v3+json",
                  "Content-Type": "application/json",
              })
              with urllib.request.urlopen(req, timeout=30) as resp:
                  return 200 <= resp.status < 300
          except urllib.error.HTTPError as e:
              print(f"  Warning: CI trigger (empty commit) failed for {branch}: HTTP {e.code} {e.reason}")
              return False
          except Exception as e:
              print(f"  Warning: CI trigger (empty commit) failed for {branch}: {e}")
              return False

      def post_pr_comment(pr_number, body):
          """Post a comment on a PR using the issues comments API."""
          url = f"https://api.github.com/repos/{repo}/issues/{pr_number}/comments"
          payload = json.dumps({"body": body}).encode()
          req = urllib.request.Request(url, data=payload, method="POST", headers={
              "Authorization": f"token {token}",
              "Accept": "application/vnd.github.v3+json",
              "Content-Type": "application/json",
          })
          try:
              with urllib.request.urlopen(req, timeout=30) as resp:
                  return 200 <= resp.status < 300
          except Exception as e:
              print(f"  Warning: could not post comment on PR #{pr_number}: {e}")
              return False

      def pr_needs_attention(pr):
          """Check if a PR has merge conflicts, is behind main, or has failing CI.
          Returns a list of issues."""
          issues = []

          # Check mergeable state
          # Need to fetch full PR details for mergeable info
          pr_url = f"https://api.github.com/repos/{repo}/pulls/{pr['number']}"
          try:
              full_pr = api_get(pr_url)
              mergeable = full_pr.get("mergeable")
              mergeable_state = full_pr.get("mergeable_state", "unknown")
              if mergeable is False:
                  issues.append("merge_conflict")
              elif mergeable_state == "dirty":
                  issues.append("merge_conflict")
          except Exception as e:
              print(f"  Warning: could not fetch mergeable state for PR #{pr['number']}: {e}")

          # Check if the PR branch is behind its base branch (e.g., main moved forward).
          # We always want to merge main first before fixing CI, so flag this explicitly.
          behind_by = get_behind_by(pr)
          if behind_by > 0 and "merge_conflict" not in issues:
              issues.append(f"behind_main: {behind_by} commit(s)")

          # Check CI status via check runs
          check_runs = get_check_runs(pr)
          failed_checks = []
          for cr in check_runs:
              conclusion = cr.get("conclusion")
              status = cr.get("status")
              name = cr.get("name", "unknown")
              if conclusion in ("failure", "timed_out", "action_required"):
                  failed_checks.append(name)
              elif status == "completed" and conclusion not in ("success", "neutral", "skipped"):
                  if conclusion is not None:
                      failed_checks.append(name)
          if failed_checks:
              issues.append(f"failing_checks: {', '.join(failed_checks)}")

          # Also check commit status API (some checks use the older status API)
          combined_status = get_check_status(pr)
          if combined_status == "failure":
              if not failed_checks:
                  issues.append("failing_status")

          # Detect missing/stale CI for autoloop PRs.
          # Pushes via GITHUB_TOKEN don't trigger workflows, so autoloop PRs
          # can sit indefinitely with no checks. Only autoloop branches are
          # eligible -- never trigger CI automatically on outside-contributor PRs.
          if is_autoloop_pr(pr):
              completed_runs = [cr for cr in check_runs if cr.get("status") == "completed"]
              # No check runs at all on this HEAD SHA
              if not check_runs:
                  issues.append("missing_checks: no check runs on HEAD")
              # All runs are still queued / in_progress and the HEAD has been
              # sitting around for a while -- likely a stuck/missing trigger.
              elif not completed_runs:
                  head_date = get_commit_date(pr["head"]["sha"])
                  if head_date:
                      # If HEAD is older than 15 minutes and nothing has completed,
                      # treat it as missing (a real CI run would have started by now).
                      try:
                          from datetime import datetime, timezone
                          ht = datetime.fromisoformat(head_date.replace("Z", "+00:00"))
                          age_s = (datetime.now(timezone.utc) - ht).total_seconds()
                          if age_s > 15 * 60:
                              issues.append("missing_checks: only queued/in-progress checks on HEAD")
                      except Exception:
                          pass

          return issues

      # --- Main logic ---

      print("=== Evergreen PR Health Check ===")
      print(f"Repository: {repo}")

      prs = get_all_open_prs()
      print(f"Found {len(prs)} open PR(s)")

      if not prs:
          print("No open PRs. Nothing to do.")
          with open(output_file, "w") as f:
              json.dump({"selected": None, "reason": "no_open_prs"}, f)
          emit_selected_output(None)
          sys.exit(0)

      # Evaluate each PR deterministically (sorted by PR number ascending)
      candidates = []
      skipped = []
      ci_triggered = []

      # If a specific PR is forced, only check that one
      if forced_pr:
          prs = [pr for pr in prs if str(pr["number"]) == forced_pr]
          if not prs:
              print(f"ERROR: PR #{forced_pr} not found among open PRs.")
              sys.exit(1)
          print(f"FORCED: checking only PR #{forced_pr}")

      for pr in sorted(prs, key=lambda p: p["number"]):
          pr_num = pr["number"]
          head_sha = pr["head"]["sha"]
          print(f"\nChecking PR #{pr_num}: {pr['title'][:60]}...")
          print(f"  Head SHA: {head_sha[:12]}")

          issues = pr_needs_attention(pr)
          if not issues:
              print(f"  Status: healthy (no issues)")
              continue

          print(f"  Issues: {issues}")

          # Handle `missing_checks` for autoloop PRs directly in the pre-flight,
          # without invoking the agent. The action is purely an API dispatch --
          # no code fix is needed -- and keeping the privileged CI trigger token
          # out of the agent context is a security win. We only do this for
          # autoloop branches (the detector enforces this), and only if
          # `missing_checks` is the *only* issue: any other issue (merge
          # conflict, behind main, failing checks) still needs the agent.
          if (
              len(issues) == 1
              and issues[0].startswith("missing_checks")
              and is_autoloop_pr(pr)
          ):
              branch = pr["head"]["ref"]
              # Cap retries on the same SHA so we don't spam-dispatch on a
              # truly-broken workflow.
              attempt_state = read_attempt_state(pr_num)
              prior_attempts = (
                  attempt_state["attempts"] if attempt_state["head_sha"] == head_sha else 0
              )
              if prior_attempts >= MAX_ATTEMPTS:
                  skipped.append({
                      "pr": pr_num,
                      "reason": (
                          f"missing_checks: max dispatch attempts ({MAX_ATTEMPTS}) "
                          f"reached on SHA {head_sha[:12]}"
                      ),
                  })
                  print(f"  SKIPPED: max missing_checks attempts reached")
                  continue
              print(f"  Triggering ci.yml on branch {branch} (attempt {prior_attempts + 1}/{MAX_ATTEMPTS})")
              ok = trigger_ci_workflow(branch)
              if ok:
                  print(f"  [+] Triggered ci.yml on {branch}")
                  workflow_url = (
                      f"https://github.com/{repo}/actions/workflows/ci.yml"
                      f"?query=branch%3A{branch}"
                  )
                  post_pr_comment(
                      pr_num,
                      (
                          "Evergreen: this PR's HEAD had no completed CI checks, "
                          "so I pushed an empty commit to trigger the `ci.yml` workflow on this branch. "
                          f"See [recent CI runs]({workflow_url}).\n\n"
                          "_(Triggered automatically because pushes via `GITHUB_TOKEN` "
                          "do not start workflows.)_"
                      ),
                  )
                  ci_triggered.append({
                      "pr": pr_num,
                      "branch": branch,
                      "head_sha": head_sha,
                  })
                  # Persist attempt count so we eventually give up if dispatching
                  # never produces check runs.
                  os.makedirs(repo_memory_dir, exist_ok=True)
                  state_path = os.path.join(repo_memory_dir, f"pr-{pr_num}.md")
                  from datetime import datetime, timezone
                  ts = datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
                  with open(state_path, "w", encoding="utf-8") as sf:
                      sf.write(
                          f"# Evergreen: PR #{pr_num}\n\n"
                          f"## State\n\n"
                          f"| Field | Value |\n"
                          f"|:---|:---|\n"
                          f"| head_sha | {head_sha} |\n"
                          f"| attempts | {prior_attempts + 1} |\n"
                          f"| last_run | {ts} |\n"
                          f"| last_result | ci_dispatched |\n"
                      )
              else:
                  print(f"  [x] Failed to dispatch ci.yml on {branch}")
                  skipped.append({
                      "pr": pr_num,
                      "reason": "missing_checks: CI trigger push failed",
                  })
              # Do NOT add to candidates -- the agent has nothing to fix here.
              continue

          # Check attempt tracking
          attempt_state = read_attempt_state(pr_num)
          if attempt_state["head_sha"] == head_sha:
              attempts = attempt_state["attempts"]
              print(f"  Attempts on this SHA: {attempts}/{MAX_ATTEMPTS}")
              if attempts >= MAX_ATTEMPTS:
                  skipped.append({
                      "pr": pr_num,
                      "reason": f"max attempts ({MAX_ATTEMPTS}) reached on SHA {head_sha[:12]}",
                  })
                  print(f"  SKIPPED: max attempts reached")
                  continue
          else:
              attempts = 0
              print(f"  New SHA detected -- resetting attempt counter")

          candidates.append({
              "pr_number": pr_num,
              "title": pr["title"],
              "head_sha": head_sha,
              "base_branch": pr["base"]["ref"],
              "head_branch": pr["head"]["ref"],
              "issues": issues,
              "attempts": attempts,
          })

      # Select the first candidate (lowest PR number -- deterministic)
      selected = candidates[0] if candidates else None

      result = {
          "selected": selected,
          "skipped": skipped,
          "ci_triggered": ci_triggered,
          "total_open_prs": len(prs),
          "candidates_found": len(candidates),
      }

      with open(output_file, "w") as f:
          json.dump(result, f, indent=2)

      if selected:
          branch = selected["head_branch"]
          print(f"Checking out PR branch before agent run: {branch}")
          subprocess.check_call(["git", "checkout", "-B", branch, f"origin/{branch}"])
          subprocess.check_call(["git", "branch", "--set-upstream-to", f"origin/{branch}", branch])
          print(f"\n>>> Selected PR #{selected['pr_number']}: {selected['title']}")
          print(f"    Issues: {selected['issues']}")
          print(f"    Attempt: {selected['attempts'] + 1}/{MAX_ATTEMPTS}")
          emit_selected_output(selected["pr_number"])
      else:
          print("\nNo PRs need attention. Nothing to do.")
          emit_selected_output(None)
          sys.exit(0)
      PYEOF

engine: copilot

features:
  copilot-requests: true
---

# Evergreen -- PR Health Keeper

You are the Evergreen agent. Your job is to fix pull requests that have merge conflicts or failing CI checks.

## Context

A pre-flight step has already identified a PR that needs attention. Read the selection data from `/tmp/gh-aw/evergreen.json` to understand which PR to fix and what issues it has.

## Workflow

1. **Read the selection file** at `/tmp/gh-aw/evergreen.json`. It contains:
   - `selected.pr_number` -- the PR to fix
   - `selected.issues` -- list of problems (e.g., `"merge_conflict"`, `"failing_checks: Test & Lint"`)
   - `selected.head_sha` -- current HEAD of the PR branch
   - `selected.head_branch` -- the PR's branch name
   - `selected.base_branch` -- the target branch (usually `main`)
   - `selected.attempts` -- how many times we've already tried on this SHA

   > If `selected` is `null`, no PRs need attention right now. Call the **noop** tool with a message like "All PRs are healthy -- nothing to fix." and stop.

2. The pre-flight step already checks out `selected.head_branch` as a named local tracking branch before you start. Keep working on that branch (do not switch back to `main` or use detached HEAD).

3. **Fix the issues** -- always follow this sequence, in order. Each push is a separate `push-to-pull-request-branch` call:

   ### Step 1 -- Merge `main` first if the PR is behind (or has conflicts)
   If `selected.issues` contains `"merge_conflict"` **or** any `"behind_main: ..."` entry, you must bring the branch up to date with `main` before doing anything else:

   - `git fetch origin main`
   - `git merge origin/main` (or `origin/<base_branch>` if the base isn't `main`)
   - Resolve any conflicts intelligently by understanding the intent of both sides. If the PR is from an autoloop branch, prefer the PR's changes for feature code and `main`'s changes for infrastructure/config.
   - Run tests/lint/typecheck locally to make sure the merge is clean.

   ### Step 2 -- Push the merge as its own commit
   - Push the merge commit using `push-to-pull-request-branch` **before doing anything else**.
   - This is the *first* push of the run. It contains *only* the merge with `main` (plus any conflict resolutions). Do **not** mix CI-fix changes into this patch.
   - Merging `main` often fixes CI on its own (the failure was just drift). After the push, re-check whether CI is still failing on the new HEAD.

   ### Step 3 -- Re-check CI after the merge
   - Look at the failing checks for the new HEAD SHA (the one you just pushed).
   - If everything is green or pending-but-likely-green, you're done -- skip to step 5.
   - If checks are still failing, continue to step 4.

   ### Step 4 -- Fix the failing checks (second push)
   - Read the failing check logs using GitHub tools.
   - Identify the root cause (test failures, lint errors, type errors, build failures).
   - Fix the code on the (now-merged) PR branch.
   - Run the relevant checks locally to verify the fix before pushing.
   - Push the fix using `push-to-pull-request-branch`. This is the *second* push of the run, and it contains *only* the CI fix -- no merge commits.

   ### Step 5 -- Update tracking and comment
   Continue to steps 4 and 5 below.

4. **Update attempt tracking** by writing to repo-memory. Write a file to the repo-memory directory at `/tmp/gh-aw/repo-memory/evergreen/pr-{number}.md` with this format:

   ```markdown
   # Evergreen: PR #{number}

   ## State

   | Field | Value |
   |:---|:---|
   | head_sha | {sha_after_push} |
   | attempts | {new_attempt_count} |
   | last_run | {ISO 8601 timestamp} |
   | last_result | {success or failure} |
   ```

   - If you **pushed a fix** (or a clean merge), set `head_sha` to the new SHA (post-push), reset `attempts` to `0`, and set `last_result` to `success`.
   - If you **could not fix** the issue, keep the original `head_sha`, increment `attempts` by 1, and set `last_result` to `failure`.

5. **Add a comment** on the PR summarizing what you did (or why you couldn't fix it). If you pushed both a merge and a CI fix, mention both.

## Rules

- **Be surgical**: make the minimum changes needed to fix the issue. Do not refactor, improve, or add features.
- **Don't break things**: always run tests/lint/typecheck locally before pushing.
- **Always merge `main` first, as its own push.** If the PR branch is behind `main` (or has merge conflicts), the *first* `push-to-pull-request-branch` call of the run must contain only the merge commit (and any conflict resolutions). Do **NOT** include a merge of `main` inside a CI-fix patch -- that's a separate, second push. Mixing the two causes patch conflicts when the remote PR branch hasn't been merged yet.
- **One concern per push**: the merge push contains only the merge; the fix push contains only the fix. Never combine them.
- **Give up gracefully**: if you cannot fix the issue after investigating, update the attempt counter and leave a comment explaining what went wrong. Do not force-push or make destructive changes.
- **One PR per run**: only fix the selected PR. Do not touch other PRs.
- **Respect the 5-attempt limit**: the pre-flight step will stop selecting this PR once attempts reach 5 on the same HEAD SHA. If the SHA changes (someone else pushes), the counter resets.

## Autoloop PRs

Evergreen treats Autoloop draft PRs (branch name `autoloop/*`, label `autoloop`) the same as human-authored PRs for CI-failure and merge-conflict fixing. These PRs are produced by the Autoloop agent (`.github/workflows/autoloop.md`), which has its own in-iteration fix-retry loop (up to 5 attempts per iteration). If Autoloop exhausts its budget or hits its per-iteration wall-clock cap, it sets `paused: true` in its state file (`{program-name}.md` on the `memory/autoloop` branch) with a `pause_reason` like `"ci-fix-exhausted: <signature>"` or `"stuck in CI fix loop: <signature>"`. The PR is left in a failing state -- deliberately, so Evergreen (or a human) can continue from there.

When Evergreen is selected for an Autoloop PR:

1. Identify it by the branch prefix `autoloop/` and/or the `autoloop` label.
2. Attempt the fix as usual -- read failing check logs, make the minimum change, run local checks, push via `push-to-pull-request-branch`.
3. If the push succeeds **and** you believe the fix is correct, also **un-pause the Autoloop program**:
   - Clone or checkout the `memory/autoloop` branch.
   - Find the state file `{program-name}.md` where `{program-name}` is the part of the branch name after `autoloop/`.
   - In the **Machine State** table, set `Paused` to `false` and `Pause Reason` to `--`.
   - Commit and push the state-file change to `memory/autoloop`.
   - Leave a comment on the Autoloop program issue (`[Autoloop: {program-name}]`, labeled `autoloop-program`) noting that Evergreen pushed a CI fix and un-paused the program, with links to the commit and the newly-green check run.
4. If you cannot fix it, the standard attempt-tracker (5 attempts per HEAD SHA) applies -- do **not** un-pause. Autoloop remains paused for human review.

> The same 5-attempts-per-SHA rule applies to Autoloop PRs: Evergreen eventually gives up rather than burning cycles on a hopelessly broken change.

## Missing CI checks (autoloop PRs only -- handled in pre-flight)

If an autoloop PR's HEAD has no CI check runs (or only stuck queued/in-progress runs), the **pre-flight step** detects this (`missing_checks` issue) and dispatches the `ci.yml` workflow directly via the GitHub API using `GH_AW_CI_TRIGGER_TOKEN`. It also posts a comment on the PR. No agent run is needed in that case, so the PR will not appear in `selected` for `missing_checks`-only issues.

You will only see a `missing_checks` issue in `selected.issues` if it is combined with another fixable issue (e.g. `merge_conflict` or `failing_checks`). In that case, focus on the *other* issue: pre-flight will re-trigger CI on the next scheduled run if checks are still missing after your push.

**Security gate:** the pre-flight only ever triggers CI for branches matching `autoloop/*`. Do not bypass this check or trigger CI yourself for PRs from outside contributors.

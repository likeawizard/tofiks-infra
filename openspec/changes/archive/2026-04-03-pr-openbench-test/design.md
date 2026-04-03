## Context

OpenBench is running at bench.likeawizard.dev with a `/scripts/` endpoint that accepts POST requests to create tests. Workers on the developer's laptop pick up and execute tests. The missing piece is automated test creation when a PR is opened on the tofiks repo.

## Goals / Non-Goals

**Goals:**
- Automatically create an OpenBench SPRT test when a PR is opened or updated
- Test compares PR branch against master using STC preset defaults
- Post a comment on the PR with a link to the test
- Keep tofiks repo workflows lightweight — just a dispatch trigger

**Non-Goals:**
- Waiting for test completion or reporting results back to the PR
- Blocking PR merge based on OpenBench results (advisory only)
- Running workers in CI (workers run on the developer's laptop)

## Decisions

### 1. Split workflow: lightweight trigger in tofiks, logic in tofiks-infra

The tofiks repo sends a `repository_dispatch` with PR metadata (branch, PR number). The tofiks-infra repo has the workflow that does the heavy lifting: build both branches, extract bench values, create the OpenBench test, and comment on the PR.

**Why:** Keeps engine and test infrastructure concerns separate. Tofiks workflows stay minimal. Same pattern as deploy-bot.

### 2. Use the /scripts/ endpoint with POST

OpenBench exposes `POST /scripts/` with `action=CREATE_TEST` for programmatic test creation. Auth is username/password in the POST body.

### 3. Trigger on PR open, synchronize, and ready_for_review — skip drafts

The tofiks trigger workflow fires on `pull_request: [opened, synchronize, ready_for_review]` and exits early if the PR is a draft. This means:
- Draft PRs are ignored entirely
- When a draft is marked ready, `ready_for_review` fires and dispatches
- Subsequent pushes to a non-draft PR dispatch via `synchronize`

### 4. Bench value from the engine

The tofiks-infra workflow checks out the tofiks repo, builds both branches, runs `./tofiks bench`, and extracts the node count. Both values are passed to the test creation request.

### 5. PR comment with test link

After creating the test, the workflow uses the GitHub API to post a comment on the tofiks PR with the OpenBench test URL. This requires the `INFRA_DISPATCH_PAT` to have permissions to comment on PRs in the tofiks repo.

### 6. Credentials

- **tofiks repo**: `INFRA_DISPATCH_PAT` (already exists) for dispatching to tofiks-infra
- **tofiks-infra repo**: `OPENBENCH_USER`, `OPENBENCH_PASS` for creating tests; `INFRA_DISPATCH_PAT` or a token with PR comment permissions on the tofiks repo

## Risks / Trade-offs

- **Cross-repo complexity** → Two workflows instead of one. But keeps concerns separated and tofiks repo clean.
- **PAT permissions** → The PAT needs both dispatch and PR comment permissions across repos. Fine-grained token scoped to both repos.
- **No auto-cancel of previous tests** → If a PR is updated, the old test keeps running. Acceptable.
- **CI build time** → Building tofiks twice (dev + base) adds ~30s. Acceptable.

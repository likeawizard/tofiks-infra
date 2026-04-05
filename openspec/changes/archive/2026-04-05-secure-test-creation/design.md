## Context

The `openbench-test.yml` workflow currently uses `OPENBENCH_ADMIN_USER`/`OPENBENCH_ADMIN_PASS` to create tests via the `/scripts/` endpoint. Admin accounts have `approver=True`, which auto-approves tests. The Location header from the API response contains the test ID but it's only used for success validation, not in the PR comment.

## Goals / Non-Goals

**Goals:**
- CI-created tests require manual approval before workers run them
- PR comments link directly to the specific test
- PGNs are saved for debugging context

**Non-Goals:**
- Label/checkbox-based gating on the tofiks repo (requires tofiks repo changes, deferred)
- Changing the OpenBench server config (`use_cross_approval` stays false)

## Decisions

### 1. Non-admin CI account

**Choice**: Create a `github-actions` OpenBench user with `enabled=True`, `approver=False`. Use new secrets `OPENBENCH_CI_USER`/`OPENBENCH_CI_PASS`.

**Rationale**: OpenBench auto-approves tests only when the creating user has `approver=True`. A non-approver account means tests land in "Pending" state. The admin (you) approves them manually from the OpenBench UI. No server config changes needed.

### 2. Extract test ID from Location header

**Choice**: The API redirects to `/index/{test_id}` on success. Parse the test ID with grep and construct the direct URL `https://bench.likeawizard.dev/test/{id}/`.

**Rationale**: The redirect URL uses `/index/{id}` format but the canonical test URL is `/test/{id}/`. Simple sed/grep extraction, no additional API calls.

### 3. PGN uploads

**Choice**: Change `upload_pgns=FALSE` to `upload_pgns=TRUE` in the API call.

**Rationale**: One-line change. OpenBench will store full PGNs per test. Disk usage is minimal for the volume of tests we run.

## Risks / Trade-offs

- **[Manual approval step]** → Every CI test now requires you to approve it. Small friction, but that's the point — prevents unauthorized cost.
- **[PGN disk usage]** → Negligible for our test volume. Can be disabled later if needed.

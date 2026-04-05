## Why

The current CI workflow creates OpenBench tests using admin credentials, which auto-approves them and immediately triggers workers (costing money). Any public PR to the tofiks repo triggers a test — a security and cost risk. Additionally, PR comments link to the base OpenBench URL rather than the specific test, and PGNs aren't saved, losing debugging context.

## What Changes

- Use a dedicated non-admin OpenBench account (`github-actions`) for CI test creation. Tests will require manual approval before workers pick them up.
- Parse the test ID from the API response and include a direct link (e.g., `/test/5/`) in the PR comment.
- Enable PGN uploads (`upload_pgns=TRUE`) so full game records are saved for debugging.

## Capabilities

### New Capabilities
- `ci-test-gating`: CI-created tests require manual approval before workers execute them
- `test-link-and-pgn`: PR comments link directly to the test, and PGNs are saved

### Modified Capabilities

## Impact

- Modified: `.github/workflows/openbench-test.yml` (credentials, link parsing, PGN flag)
- New GitHub secrets: `OPENBENCH_CI_USER`, `OPENBENCH_CI_PASS`
- Requires creating a new OpenBench user account (manual step on bench.likeawizard.dev)

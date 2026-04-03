## Why

Currently OpenBench tests must be created manually via the web UI. When a PR is opened on the tofiks repo, an automated SPRT test against master would catch regressions and validate improvements without manual intervention.

## What Changes

- **GitHub Action in the tofiks repo**: on PR open/sync, POST to bench.likeawizard.dev/scripts/ to create an OpenBench test (PR branch vs master)
- **OpenBench credentials** stored as secrets on the tofiks repo
- **PR comment** with a link to the OpenBench test page so the author can track progress

## Capabilities

### New Capabilities
- `pr-openbench-trigger`: GitHub Actions workflow that creates an OpenBench test when a PR is opened or updated on the tofiks repo

### Modified Capabilities

None.

## Impact

- **tofiks repo**: New workflow file `.github/workflows/openbench-test.yml`
- **tofiks repo secrets**: OPENBENCH_USER, OPENBENCH_PASS for authenticating with the scripts endpoint
- **OpenBench server**: Receives automated test creation requests — no server changes needed

## Why

We need a way to benchmark tofiks changes against the current version. OpenBench is the standard tool for distributed chess engine testing (SPRT tests, Elo estimation). Having our own instance gives us automated regression testing on PRs and a public results page for contributors.

## What Changes

- **OpenBench container** added to docker-compose, running the Django app with gunicorn, exposed on port 80
- **OpenBench fork** added as a git submodule to pin the version
- **SQLite database** stored in a Docker volume on the Hetzner instance (data loss is acceptable — tests can be re-run)
- **Initial setup automation** for Django migrations and admin user creation
- **Deploy pipeline** updated to include the OpenBench service

## Capabilities

### New Capabilities
- `openbench-server`: OpenBench Django application containerized and publicly accessible on the Hetzner instance

### Modified Capabilities

None.

## Impact

- **Hetzner server**: Port 80 now serves OpenBench web UI (publicly accessible)
- **docker-compose.yml**: New service added
- **Deploy workflow**: No changes needed — existing `docker compose up --build` picks up the new service
- **Disk**: SQLite DB grows with test results, minimal footprint

## Context

We already have a working Docker-based deploy pipeline for lichess-bot on Hetzner. OpenBench is a Django application that uses SQLite by default. It needs to be added as a second service in the existing docker-compose stack, publicly accessible for viewing test results.

## Goals / Non-Goals

**Goals:**
- OpenBench web UI accessible at the server's public IP
- Database persisted across container restarts via Docker volume
- Admin user created automatically on first deploy
- Deploys alongside lichess-bot with the existing pipeline

**Non-Goals:**
- HTTPS / domain name (future improvement)
- OpenBench workers (separate phase)
- PR-triggered test creation (separate phase)
- High-availability or database backups

## Decisions

### 1. OpenBench as a git submodule

Same pattern as lichess-bot. Fork OpenBench to `likeawizard/OpenBench`, add as submodule to pin the version.

**Why fork:** Allows customization (e.g., default engine config) and ensures stability if upstream changes.

### 2. Gunicorn for production serving

Django's `runserver` is not production-ready. Use gunicorn with a few workers behind port 8000, exposed as port 80 on the host.

**Why not nginx in front:** For a single Django app on a personal server, gunicorn directly serving HTTP is fine. Add nginx later if needed for HTTPS or static file performance.

### 3. SQLite with Docker named volume

OpenBench defaults to SQLite. A Docker named volume (`openbench-data`) persists the database file across container rebuilds. No PostgreSQL needed at this scale.

**Why not bind mount:** Named volumes are managed by Docker, portable, and don't depend on host directory permissions.

### 4. Admin user via environment variables

Pass `OPENBENCH_ADMIN_USER` and `OPENBENCH_ADMIN_PASS` as secrets. The entrypoint script creates the superuser if it doesn't exist (idempotent).

### 5. Separate Dockerfile for OpenBench

A `Dockerfile.openbench` in the repo root. Installs Python deps, copies the OpenBench source, runs migrations in the entrypoint, then starts gunicorn.

## Risks / Trade-offs

- **SQLite under concurrent writes** → Fine for a small instance. OpenBench's write pattern (test results from workers) is low-frequency. If workers are added later and contention appears, migrate to PostgreSQL then.
- **No HTTPS** → Test results are public and non-sensitive. Credentials go over plain HTTP — acceptable risk for a personal instance, but worth revisiting with a domain + Let's Encrypt.
- **Port 80 conflict** → If we later want to serve both OpenBench and another web service, we'll need a reverse proxy. For now, only OpenBench needs a public port.

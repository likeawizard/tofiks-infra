## Context

The Hetzner server (188.34.201.182) currently runs lichess-bot as a systemd service with tofiks built from source using Go installed directly on the host. Deployment is via a daily cron GitHub Action that SSHs in, `git pull`s, builds, copies the binary, and restarts the service. This setup has caused friction with Go version management and systemd config maintenance.

The tofiks-infra repo is a new repo intended to own all deployment configuration.

## Goals / Non-Goals

**Goals:**
- Reproducible, containerized deployment of lichess-bot + tofiks
- One-button deploy from tofiks-infra (push to main → deploy)
- Auto-redeploy on tofiks merge to main
- Secrets (lichess API token, SSH key) managed in GitHub Secrets, never committed
- Foundation for adding more services (OpenBench) later via docker-compose

**Non-Goals:**
- OpenBench setup (future change)
- HTTPS / reverse proxy (not needed until OpenBench)
- Container orchestration beyond docker-compose (no k8s)
- Automated server provisioning (Hetzner instance already exists)

## Decisions

### 1. Multi-stage Docker build for lichess-bot + tofiks

Build tofiks from source in a Go stage, copy the binary into a Python runtime image that runs lichess-bot. Single container, single process (lichess-bot.py).

**Why over separate containers:** lichess-bot spawns tofiks as a subprocess via UCI protocol. They must share a filesystem. Two containers would require IPC complexity for zero benefit.

**Why over pre-built binary:** Building from source in Docker ensures the binary matches the target architecture and Go version without maintaining build artifacts.

### 2. lichess-bot pinned to a specific commit via Git submodule

Add lichess-bot as a git submodule in tofiks-infra. This gives us version control over which lichess-bot version we run.

**Why submodule over copying files:** lichess-bot is an external project. Submodule tracks the upstream commit and makes upgrades explicit.

**Why not latest:** The current server runs an Oct 2023 version. We pin to that first, upgrade separately.

### 3. Config and opening book stored in tofiks-infra

`config.yml` lives in the repo with the API token replaced by an environment variable placeholder. The polyglot book (`tofiks.bin`) is committed to the repo (it's a small binary, ~few KB).

At deploy time, the lichess API token is injected via Docker environment variable from GitHub Secrets.

### 4. Deploy via SSH + docker compose

The GitHub Action SSHs into Hetzner, pulls the latest tofiks-infra, and runs `docker compose up -d --build`. Simple, no registry needed — images are built on the server.

**Why not a container registry:** Adds complexity (auth, storage, push/pull) for a single-server setup. Building on the server is fast enough for a Go binary + Python app.

### 5. Cross-repo deploy trigger from tofiks

The tofiks repo triggers a `repository_dispatch` event on tofiks-infra when merging to main. The tofiks-infra deploy workflow listens for both `push` (its own changes) and `repository_dispatch` (engine updates).

**Why repository_dispatch over workflow_dispatch:** Can be triggered programmatically from another repo's Action without a PAT for workflow_dispatch API.

### 6. Docker installed on Hetzner as a prerequisite

Docker and Docker Compose will be installed on the server as a one-time manual step (or scripted in the first deploy). Not managed by this change beyond documenting the requirement.

## Risks / Trade-offs

- **Build on server uses server CPU** → Acceptable for infrequent deploys. Go cross-compilation is fast (~10s). If it becomes an issue, switch to a registry later.
- **Single container = coupled lifecycle** → Upgrading lichess-bot requires rebuilding. Fine for this scale; the image builds in seconds.
- **Server disk fills with old Docker images** → Add `docker image prune -f` to the deploy script.
- **Downtime during deploy** → Brief (seconds) while container restarts. Lichess handles reconnection. Acceptable for a bot.
- **repository_dispatch needs a PAT** → Actually needs a fine-grained token with repo scope on tofiks-infra. Store as a secret in the tofiks repo.

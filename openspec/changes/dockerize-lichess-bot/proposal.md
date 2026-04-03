## Why

The current deployment is fragile — a systemd service running lichess-bot directly on the Hetzner host, with Go installed on the server to build tofiks from source. This has caused issues with manual Go version upgrades, systemd config drift, and no reproducible deploy process. Dockerizing the setup gives us reproducible builds, one-command deploys from GitHub Actions, and a foundation for adding OpenBench later.

## What Changes

- **Multi-stage Dockerfile** that builds the tofiks Go binary and packages it into a lichess-bot Python container
- **docker-compose.yml** to orchestrate the lichess-bot service (extensible for future services)
- **GitHub Action in tofiks-infra** that deploys the stack to Hetzner via SSH on push to main
- **GitHub Action in tofiks repo** that triggers a redeploy of the bot container on merge to main
- **Lichess-bot config and polyglot book** stored in this repo (secrets like the API token stay in GitHub Secrets)
- **Replaces** the existing systemd service, manual Go builds, and the daily-cron deploy workflow in the tofiks repo

## Capabilities

### New Capabilities
- `lichess-bot-container`: Multi-stage Docker build combining tofiks (Go) and lichess-bot (Python) into a single container with config and opening book
- `deploy-pipeline`: GitHub Actions workflow for one-button deploy of the full stack from tofiks-infra, plus cross-repo trigger from tofiks on merge to main

### Modified Capabilities

None — no existing specs.

## Impact

- **Hetzner server**: Docker + Docker Compose need to be installed. Existing systemd service and bare-metal Go installation become unused
- **tofiks repo**: The existing `deploy.yml` workflow (SSH + go build + systemctl) gets replaced by a workflow that triggers tofiks-infra deploy
- **GitHub Secrets**: SSH key, host, lichess API token need to be configured on both repos (or use a shared org secret)
- **lichess-bot**: Pinning to a specific version/commit rather than the Oct 2023 clone currently on the server

## 1. Repository Setup

- [ ] 1.1 Add lichess-bot as a git submodule pinned to the Oct 2023 commit currently deployed
- [ ] 1.2 Add lichess-bot `config.yml` to tofiks-infra (with token placeholder, not hardcoded)
- [ ] 1.3 Add `tofiks.bin` polyglot opening book to the repo

## 2. Docker Build

- [ ] 2.1 Create multi-stage Dockerfile: Go builder (GOAMD64=v3) + Python runtime with lichess-bot
- [ ] 2.2 Create `docker-compose.yml` with lichess-bot service, passing `LICHESS_TOKEN` env var
- [ ] 2.3 Test local build and verify container starts correctly

## 3. Deploy Pipeline (tofiks-infra)

- [ ] 3.1 Create GitHub Actions deploy workflow: triggers on push to main, workflow_dispatch, and repository_dispatch
- [ ] 3.2 Deploy script: SSH into Hetzner, git pull (with submodules), docker compose up --build, prune old images
- [ ] 3.3 Write `.env` file on server from GitHub Secrets (LICHESS_TOKEN) during deploy

## 4. Cross-repo Trigger (tofiks)

- [ ] 4.1 Create GitHub Actions workflow in tofiks repo: on merge to main, send repository_dispatch to tofiks-infra
- [ ] 4.2 Remove or disable the old deploy.yml workflow in tofiks repo

## 5. Server Preparation

- [ ] 5.1 Install Docker and Docker Compose on Hetzner server
- [ ] 5.2 Clone tofiks-infra repo on the server (with submodules)
- [ ] 5.3 Stop and disable the existing systemd tofiks.service
- [ ] 5.4 Verify bot connects to lichess and plays games via Docker

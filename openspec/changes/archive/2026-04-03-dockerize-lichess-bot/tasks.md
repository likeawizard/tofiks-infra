## 1. Repository Setup

- [x] 1.1 Add lichess-bot as a git submodule pinned to the Oct 2023 commit currently deployed
- [x] 1.2 Add lichess-bot `config.yml` to tofiks-infra (with token placeholder, not hardcoded)
- [x] 1.3 Add `tofiks.bin` polyglot opening book to the repo (via Git LFS)

## 2. Docker Build

- [x] 2.1 Create multi-stage Dockerfile: Go builder (GOAMD64=v3) + Python runtime with lichess-bot
- [x] 2.2 Create `docker-compose.yml` with lichess-bot service, passing `LICHESS_TOKEN` env var
- [x] 2.3 Test local build and verify container starts correctly

## 3. GitHub Secrets

- [x] 3.1 Set secrets on tofiks-infra repo via `gh secret set`: SSH_KEY, SSH_HOST, SSH_USER, LICHESS_TOKEN
- [x] 3.2 Set secrets on tofiks repo via `gh secret set`: INFRA_DISPATCH_PAT for repository_dispatch to tofiks-infra

## 4. Deploy Pipeline (tofiks-infra)

- [x] 4.1 Create GitHub Actions deploy workflow: triggers on push to main, workflow_dispatch, and repository_dispatch
- [x] 4.2 Deploy script: SSH into Hetzner, git pull (with submodules), docker compose up --build, prune old images
- [x] 4.3 Write `.env` file on server from GitHub Secrets (LICHESS_TOKEN) during deploy

## 5. Cross-repo Trigger (tofiks)

- [x] 5.1 Create GitHub Actions workflow in tofiks repo: on merge to main, send repository_dispatch to tofiks-infra
- [x] 5.2 Replaced old deploy-bot.yml (was SSH+go build+systemctl) with repository_dispatch trigger

## 6. Server Preparation

- [x] 6.1 Install Docker and Docker Compose on Hetzner server
- [x] 6.2 Clone tofiks-infra repo on the server at /opt/tofiks-infra
- [x] 6.3 Stop and disable the existing systemd tofiks.service
- [x] 6.4 Verify bot connects to lichess and plays games via Docker (verified locally; server deploy pending merge)

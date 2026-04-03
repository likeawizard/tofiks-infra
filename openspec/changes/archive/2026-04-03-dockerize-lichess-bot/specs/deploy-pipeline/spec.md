## ADDED Requirements

### Requirement: tofiks-infra deploys on push to main
A GitHub Actions workflow in tofiks-infra SHALL deploy the Docker stack to the Hetzner server when changes are pushed to the main branch.

#### Scenario: Push to main triggers deploy
- **WHEN** a commit is pushed to `main` on tofiks-infra
- **THEN** the workflow SSHs into the Hetzner server, pulls latest code, and runs `docker compose up -d --build`

#### Scenario: Deploy cleans up old images
- **WHEN** the deploy completes successfully
- **THEN** old unused Docker images are pruned to prevent disk bloat

### Requirement: tofiks merge to main triggers bot redeploy
A GitHub Actions workflow in the tofiks repo SHALL trigger a redeploy of the lichess-bot container when a PR is merged to main, so the bot always runs the latest engine.

#### Scenario: Merge to tofiks main triggers tofiks-infra deploy
- **WHEN** a PR is merged to `main` on the tofiks repo
- **THEN** a `repository_dispatch` event is sent to tofiks-infra, triggering its deploy workflow

### Requirement: Secrets are managed in GitHub
All sensitive values (SSH key, server host, lichess API token, dispatch PAT) SHALL be stored as GitHub Secrets, never committed to the repository.

#### Scenario: Deploy uses secrets for SSH
- **WHEN** the deploy workflow runs
- **THEN** it uses `SSH_KEY`, `SSH_HOST`, and `SSH_USER` secrets to connect to the server

#### Scenario: Lichess token is passed to the container
- **WHEN** the deploy workflow runs `docker compose up`
- **THEN** the `LICHESS_TOKEN` secret is written to the server environment and passed to the container

### Requirement: Deploy supports manual trigger
The deploy workflow SHALL support `workflow_dispatch` for manual one-button deploys.

#### Scenario: Manual deploy from GitHub UI
- **WHEN** a user triggers the workflow manually from the Actions tab
- **THEN** the full deploy runs identically to an automated trigger

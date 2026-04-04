## Context

OpenBench workers are currently deployed manually or via the existing docker-compose on the main server. Scaling testing requires launching additional Hetzner servers, each needing manual setup. Hetzner offers a Docker CE image and cloud-config (cloud-init) support, enabling fully automated server provisioning.

The worker already supports `--fleet` mode which exits when no work is available. We leverage this to build ephemeral workers that self-terminate.

## Goals / Non-Goals

**Goals:**
- Zero-touch worker provisioning via Hetzner cloud-config
- Workers auto-shutdown when idle to minimize costs
- Reuse existing Dockerfile.worker without modification
- Simple workflow: paste cloud-config, set variables, create server

**Non-Goals:**
- Hetzner API automation (creating/deleting servers programmatically)
- Pre-built Docker images or container registry
- Monitoring or alerting for worker status
- Auto-scaling (user manually creates servers from Hetzner panel)

## Decisions

### 1. Build Docker image on the server vs. pull from registry

**Choice**: Build on the server by cloning the repo and running `docker compose build`.

**Rationale**: No container registry to maintain. The build takes ~1-2 minutes which is acceptable for ephemeral workers. The Dockerfile.worker is small and deterministic.

**Alternative**: Push images to Docker Hub/GHCR. Faster startup but adds registry management complexity for minimal benefit.

### 2. Cloud-config writes all files inline vs. cloning repo

**Choice**: Clone the repo on the server. Cloud-config only writes the `.env` file and a startup script.

**Rationale**: Keeps the cloud-config small and maintainable. The Dockerfile and docker-compose are already in the repo. Avoids duplicating file contents in cloud-config YAML.

### 3. Standalone docker-compose.worker.yml vs. reusing main compose

**Choice**: Create a dedicated `docker-compose.worker.yml` that only defines the worker service with `--fleet` flag.

**Rationale**: The main docker-compose includes openbench server, caddy, lichess-bot — none of which are needed on worker nodes. A dedicated file is cleaner and avoids accidentally starting other services.

### 4. Shutdown mechanism

**Choice**: A wrapper script that runs `docker compose up` (blocking), then executes `shutdown -h now` when the container exits.

**Rationale**: Simple and reliable. The `--fleet` flag causes the worker process to exit when there's no work. Docker compose exits when the service stops. The wrapper script then shuts down the server. Hetzner stops billing for powered-off servers (though the server still exists — user deletes it manually or could use Hetzner API later).

**Note**: `shutdown -h now` powers off the server. Hetzner bills for powered-off servers at a reduced rate. The user should delete the server from the Hetzner panel to fully stop billing.

## Risks / Trade-offs

- **[Build time on first boot]** → ~2-3 minutes for docker build. Acceptable for batch testing workloads.
- **[Server not deleted after shutdown]** → Hetzner still bills for powered-off servers (disk storage). User must manually delete. Document this clearly.
- **[Credentials in cloud-config]** → Cloud-config is visible in Hetzner console. Use a dedicated non-admin worker account to limit blast radius. Document this.
- **[Fleet mode exit delay]** → Worker waits `TIMEOUT_ERROR` seconds before exiting when no work. This is a short delay, not a significant cost.

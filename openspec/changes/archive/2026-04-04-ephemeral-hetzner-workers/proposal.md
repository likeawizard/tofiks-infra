## Why

Spinning up OpenBench workers currently requires manual server setup — installing Docker, cloning the repo, configuring environment variables, and running docker-compose. This makes it impractical to quickly scale testing by launching multiple Hetzner servers. We need a one-shot cloud-config that turns a fresh Hetzner Docker CE server into a working OpenBench worker with zero manual intervention, and shuts itself down when there's no more work to minimize costs.

## What Changes

- Add a `worker-cloud-config.yml` cloud-init script that can be pasted directly into Hetzner's cloud-config field when creating a server
- The cloud-config will: clone the repo, build the worker Docker image, and run it with `--fleet` mode
- When the worker exits (no more work), the server automatically shuts down via `shutdown -h now`
- Add a `.env.worker.example` documenting required variables (credentials, threads, server URL)
- Add a standalone `docker-compose.worker.yml` for just the worker service (used by cloud-config)
- Document the workflow: create a dedicated OpenBench worker account, configure cloud-config, launch servers

## Capabilities

### New Capabilities
- `ephemeral-worker`: Cloud-init provisioning of disposable OpenBench worker servers on Hetzner, including auto-shutdown when idle

### Modified Capabilities

## Impact

- New files: `worker-cloud-config.yml`, `docker-compose.worker.yml`, `.env.worker.example`
- No changes to existing services or Docker images
- Reuses existing `Dockerfile.worker` as-is
- Requires a dedicated (non-admin) OpenBench user account for workers

## Why

After setting up cloud-config for ephemeral workers, the remaining manual step is creating and deleting Hetzner servers. A Go service can automate this: detect pending OpenBench work, spin up workers, and destroy them when done — fully hands-off scaling with no idle server costs.

## What Changes

- New Go service `worker-manager/` that runs as a Docker container alongside OpenBench
- Polls OpenBench to detect active/pending tests (via HTTP to the local OpenBench instance)
- Creates Hetzner servers via hcloud-go SDK when work is detected
- Destroys Hetzner servers when all work is finished
- Manages worker lifecycle: tracks which servers it created, avoids duplicate creation
- Added to docker-compose.yml as a new service
- Configurable via environment variables: Hetzner API token, server type, max workers, poll interval

## Capabilities

### New Capabilities
- `work-detection`: Detecting whether OpenBench has pending or active tests that need workers
- `server-lifecycle`: Creating and destroying Hetzner servers based on work availability

### Modified Capabilities

## Impact

- New directory: `worker-manager/` (Go module with Dockerfile)
- Modified: `docker-compose.yml` (new worker-manager service)
- Modified: `.env` / `.env.worker.example` (new HCLOUD_TOKEN variable)
- Depends on: existing `worker-cloud-config.yml` (passed as user-data to new servers)
- External dependency: `github.com/hetznercloud/hcloud-go`

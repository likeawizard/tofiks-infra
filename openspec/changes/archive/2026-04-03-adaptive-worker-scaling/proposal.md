## Why

The Hetzner server has 4 vCPUs shared between the lichess-bot (up to 2 concurrent games = 2 tofiks processes) and the OpenBench worker. Currently the worker only runs on the developer's laptop. Running a worker on the server too would mean tests progress 24/7, but it must not starve the lichess-bot of CPU when games are active.

## What Changes

- **Server-side OpenBench worker** running on the Hetzner instance as a docker-compose service (not under the `worker` profile — always on)
- **Adaptive scaling script** that monitors lichess-bot game activity and adjusts the worker thread count:
  - No games active → 2 worker threads (using all spare vCPUs)
  - 1+ games active → 1 worker thread (leave headroom for the bot)
- **Monitoring via process count** — check for running tofiks processes inside the lichess-bot container to determine if games are active

## Capabilities

### New Capabilities
- `server-worker-scaling`: Adaptive OpenBench worker on the Hetzner server that scales thread count based on lichess-bot activity

### Modified Capabilities

None.

## Impact

- **docker-compose.yml**: New server-worker service (always on, not under worker profile)
- **Hetzner server**: CPU shared between bot and worker. Worker yields when games are active
- **Deploy workflow**: Server-worker starts/restarts with deploy
- **Laptop worker**: Unchanged — can still run independently via `make worker-start`

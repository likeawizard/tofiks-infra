## Why

OpenBench tests need workers to run games. Without workers, tests sit idle in the queue. The first worker should run on the developer's laptop (not the server) to keep the Hetzner instance lean for hosting. A Docker container ensures the worker has all build dependencies (Python, Go, g++) without polluting the host system.

## What Changes

- **Dockerfile.worker** that bundles Python, Go, g++, and the OpenBench client — everything needed to build tofiks and run games
- **Worker service** added to docker-compose.yml under a `worker` profile so it doesn't start on server deploys
- **Tofiks Makefile** updated with an `EXE=` variable so OpenBench's build system can control the output binary path

## Capabilities

### New Capabilities
- `openbench-worker`: Docker-based OpenBench worker that connects to bench.likeawizard.dev, builds tofiks from source, and runs SPRT tests

### Modified Capabilities

None.

## Impact

- **docker-compose.yml**: New `worker` service under `worker` profile — does not affect server deploys
- **tofiks repo Makefile**: Needs a default target that accepts `EXE=` for OpenBench compatibility
- **Laptop**: Run `docker compose --profile worker up worker` to start contributing games

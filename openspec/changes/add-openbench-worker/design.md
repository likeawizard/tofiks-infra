## Context

OpenBench is deployed at bench.likeawizard.dev with Tofiks as the only engine. Tests can be created via the web UI but no workers exist to execute them. The developer's laptop is the intended first worker — it has CPU for running games and the Docker setup avoids dependency management.

## Goals / Non-Goals

**Goals:**
- Worker container with Python, Go, and g++ (for fastchess)
- Connects to bench.likeawizard.dev using OpenBench credentials
- Runs as a docker compose profile so it's opt-in
- Configurable thread count via environment variable
- Tofiks Makefile compatible with OpenBench's `make -j EXE=<path>` build system

**Non-Goals:**
- Running workers on the Hetzner server (future, separate change)
- Auto-scaling or multiple worker orchestration
- Worker monitoring or alerting

## Decisions

### 1. Docker compose profile for the worker

The worker uses `profiles: [worker]` so `docker compose up` on the server ignores it. Start it explicitly with `docker compose --profile worker up worker`.

**Why not a separate compose file:** Single file is simpler. Profiles are the idiomatic way to handle optional services.

### 2. Worker uses `--no-client-downloads` flag

The OpenBench client normally auto-updates itself from the server. Since we bake the client into the Docker image, we skip this to avoid writing to the container filesystem. The image is rebuilt when tofiks-infra is updated.

### 3. Tofiks Makefile needs a default target with EXE=

OpenBench workers run `make -j EXE=<output_path>`. The current Makefile's default target is `build-tofiks` which hardcodes the output name. A new default target must accept `EXE=` to be OpenBench-compatible.

## Risks / Trade-offs

- **Worker on laptop = intermittent** → Tests only progress when the laptop is running the worker. Acceptable for personal use.
- **Go download in Dockerfile = larger image** → ~300MB for Go SDK. Unavoidable since we need to compile from source.
- **`--no-client-downloads` may drift from server** → If OpenBench updates its client protocol, rebuild the Docker image. Low risk since we control the server version.

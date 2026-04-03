## ADDED Requirements

### Requirement: Worker container has all build dependencies
The worker Docker image SHALL include Python 3.12, Go 1.26, g++, make, and git — everything needed to build tofiks and fastchess from source.

#### Scenario: Worker can build tofiks
- **WHEN** the worker receives a workload for tofiks
- **THEN** it clones the repo, runs `make -j EXE=<path>`, and produces a working binary

#### Scenario: Worker can build fastchess
- **WHEN** the worker starts and needs to build the match runner
- **THEN** it compiles fastchess using g++

### Requirement: Worker connects to the OpenBench server
The worker SHALL authenticate with bench.likeawizard.dev using credentials from environment variables and poll for workloads.

#### Scenario: Worker authenticates and waits for work
- **WHEN** the worker container starts with valid credentials
- **THEN** it connects to the server and begins polling for workloads

#### Scenario: Worker reports results
- **WHEN** the worker completes a batch of games
- **THEN** results are submitted to the server and visible in the web UI

### Requirement: Worker is an opt-in compose profile
The worker service SHALL use a docker compose profile so it does not start during server deploys.

#### Scenario: Server deploy ignores worker
- **WHEN** `docker compose up -d --build` runs on the server
- **THEN** the worker container is not started

#### Scenario: Developer starts worker explicitly
- **WHEN** `docker compose --profile worker up worker` is run on the laptop
- **THEN** the worker container starts and connects to the server

### Requirement: Thread count is configurable
The number of worker threads SHALL be configurable via the `WORKER_THREADS` environment variable, defaulting to 4.

#### Scenario: Custom thread count
- **WHEN** `WORKER_THREADS=8 docker compose --profile worker up worker` is run
- **THEN** the worker uses 8 threads for running games

### Requirement: Tofiks Makefile supports EXE variable
The tofiks Makefile SHALL have a default target that accepts `EXE=<path>` to control the output binary location, compatible with OpenBench's build system.

#### Scenario: OpenBench build command works
- **WHEN** `make -j EXE=/tmp/tofiks` is run in the tofiks repo
- **THEN** the tofiks binary is produced at `/tmp/tofiks`

## ADDED Requirements

### Requirement: Server runs an OpenBench worker continuously
An OpenBench worker SHALL run on the Hetzner server as a docker-compose service that starts automatically on deploy.

#### Scenario: Worker starts on deploy
- **WHEN** `docker compose up -d --build` runs on the server
- **THEN** the server worker starts and connects to bench.likeawizard.dev

#### Scenario: Worker is independent from laptop worker
- **WHEN** both the server worker and laptop worker are running
- **THEN** both connect to OpenBench and share workloads independently

### Requirement: Worker thread count adapts to lichess-bot activity
A scaling script SHALL monitor the lichess-bot container for active tofiks processes and adjust the worker thread count accordingly.

#### Scenario: No games active — scale up
- **WHEN** no tofiks processes are running in the lichess-bot container
- **THEN** the worker runs with 2 threads

#### Scenario: Games active — scale down
- **WHEN** one or more tofiks processes are running in the lichess-bot container
- **THEN** the worker runs with 1 thread

### Requirement: Scaling is debounced to avoid thrashing
The scaling script SHALL wait for game activity to be stable before changing the worker thread count.

#### Scenario: Brief game start does not immediately scale down
- **WHEN** a game starts and ends within the debounce period
- **THEN** the worker thread count is not changed

#### Scenario: Sustained game activity triggers scale down
- **WHEN** games have been active for longer than the debounce period
- **THEN** the worker is restarted with 1 thread

### Requirement: Scaling script detects games via process count
The script SHALL detect active games by counting tofiks processes inside the lichess-bot container using the Docker API or Docker exec.

#### Scenario: Process detection works
- **WHEN** the lichess-bot is playing a game
- **THEN** `docker exec lichess-bot` shows tofiks processes with count > 0

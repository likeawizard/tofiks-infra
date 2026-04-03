## ADDED Requirements

### Requirement: Multi-stage Docker image builds tofiks and runs lichess-bot
The Dockerfile SHALL use a multi-stage build: a Go builder stage compiles the tofiks binary with `GOAMD64=v3`, and a Python runtime stage runs lichess-bot with the binary copied in.

#### Scenario: Image builds successfully
- **WHEN** `docker build` is run on the Dockerfile
- **THEN** the resulting image contains the tofiks binary at `/usr/local/bin/tofiks` and the lichess-bot Python application ready to run

#### Scenario: Tofiks binary is optimized for modern x86
- **WHEN** the Go builder stage compiles tofiks
- **THEN** it uses `GOAMD64=v3` to enable AVX2 and related instructions

### Requirement: lichess-bot configuration is bundled with secret injection
The container SHALL include the lichess-bot `config.yml` with all static settings (matchmaking, greetings, UCI options). The lichess API token SHALL NOT be in the image and MUST be injected via environment variable at runtime.

#### Scenario: Container starts with API token from environment
- **WHEN** the container starts with `LICHESS_TOKEN` environment variable set
- **THEN** lichess-bot uses that token to authenticate with lichess.org

#### Scenario: Container fails without API token
- **WHEN** the container starts without `LICHESS_TOKEN` environment variable
- **THEN** the container exits with an error indicating the missing token

### Requirement: Polyglot opening book is included in the image
The tofiks polyglot book (`tofiks.bin`) SHALL be included in the Docker image alongside the engine binary.

#### Scenario: Engine uses opening book
- **WHEN** lichess-bot starts a game and the position is in the book
- **THEN** tofiks uses the polyglot book moves as configured in config.yml

### Requirement: lichess-bot version is pinned via git submodule
The lichess-bot source SHALL be tracked as a git submodule in tofiks-infra, pinned to a specific commit.

#### Scenario: Submodule checkout provides lichess-bot source
- **WHEN** the repo is cloned with `--recurse-submodules`
- **THEN** the lichess-bot directory contains the pinned version of the lichess-bot source code

### Requirement: docker-compose orchestrates the service
A `docker-compose.yml` SHALL define the lichess-bot service with the API token passed from an environment variable. It MUST be extensible for future services.

#### Scenario: Stack starts with compose
- **WHEN** `docker compose up -d` is run with a `.env` file or environment containing `LICHESS_TOKEN`
- **THEN** the lichess-bot container starts and connects to lichess.org

#### Scenario: Stack stops cleanly
- **WHEN** `docker compose down` is run
- **THEN** the lichess-bot container stops gracefully

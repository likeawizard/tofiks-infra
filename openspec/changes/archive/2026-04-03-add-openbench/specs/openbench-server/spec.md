## ADDED Requirements

### Requirement: OpenBench runs as a Docker service
The OpenBench Django application SHALL run as a service in docker-compose, using gunicorn as the application server.

#### Scenario: Service starts successfully
- **WHEN** `docker compose up -d` is run
- **THEN** the OpenBench service starts and serves the web UI on port 8000 inside the container

#### Scenario: Web UI is publicly accessible
- **WHEN** a user navigates to `http://<server-ip>/`
- **THEN** the OpenBench web interface is displayed

### Requirement: Database is persisted via Docker volume
The SQLite database SHALL be stored in a Docker named volume so data survives container rebuilds.

#### Scenario: Data survives container rebuild
- **WHEN** `docker compose down && docker compose up -d --build` is run
- **THEN** previously created tests and results are still present in the OpenBench UI

#### Scenario: Fresh start is possible
- **WHEN** `docker compose down -v` is run (removing volumes)
- **THEN** the database is wiped and OpenBench starts fresh on next launch

### Requirement: Admin user is created automatically
The entrypoint SHALL create a Django superuser from environment variables if one does not already exist. The user MUST also be enabled and approved for OpenBench access.

#### Scenario: First deploy creates admin
- **WHEN** the container starts for the first time with `OPENBENCH_ADMIN_USER` and `OPENBENCH_ADMIN_PASS` set
- **THEN** a superuser account is created and enabled

#### Scenario: Subsequent deploys skip creation
- **WHEN** the container starts and the admin user already exists
- **THEN** no duplicate user is created and the existing account is unchanged

### Requirement: Django migrations run automatically
The entrypoint SHALL run `python manage.py migrate` before starting gunicorn to ensure the database schema is up to date.

#### Scenario: Migrations applied on startup
- **WHEN** the container starts
- **THEN** all pending Django migrations are applied before the application serves requests

### Requirement: OpenBench is pinned as a git submodule
The OpenBench source SHALL be tracked as a git submodule in tofiks-infra, forked to the user's GitHub account.

#### Scenario: Submodule provides OpenBench source
- **WHEN** the repo is cloned with `--recurse-submodules`
- **THEN** the OpenBench directory contains the pinned version of the source code

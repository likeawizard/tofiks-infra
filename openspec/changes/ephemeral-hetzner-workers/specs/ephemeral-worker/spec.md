## ADDED Requirements

### Requirement: Cloud-config provisions a working OpenBench worker
The system SHALL provide a cloud-init cloud-config file that, when applied to a Hetzner server with Docker CE image, results in a fully operational OpenBench worker without manual intervention.

#### Scenario: Fresh server boot with cloud-config
- **WHEN** a Hetzner server is created with the Docker CE image and the cloud-config applied
- **THEN** the server SHALL clone the repository, build the worker Docker image, and start the worker connected to the configured OpenBench server

#### Scenario: Worker authenticates with OpenBench
- **WHEN** the worker container starts
- **THEN** it SHALL authenticate using the credentials provided via cloud-config variables (OPENBENCH_USERNAME, OPENBENCH_PASSWORD) against the configured OPENBENCH_SERVER

### Requirement: Worker runs in fleet mode
The worker SHALL run with the `--fleet` flag, causing it to exit when no work is available on the OpenBench server.

#### Scenario: Work available
- **WHEN** the OpenBench server has pending test workloads
- **THEN** the worker SHALL pick up and execute workloads using the configured number of threads

#### Scenario: No work available
- **WHEN** the OpenBench server has no pending workloads
- **THEN** the worker process SHALL exit

### Requirement: Server auto-shutdown on worker exit
The system SHALL shut down the Hetzner server when the worker container exits, to stop incurring compute costs.

#### Scenario: Worker exits after completing work
- **WHEN** the worker container exits (fleet mode, no more work)
- **THEN** the server SHALL execute `shutdown -h now` to power off

#### Scenario: Worker exits due to error
- **WHEN** the worker container exits with a non-zero exit code
- **THEN** the server SHALL still shut down to avoid idle billing

### Requirement: Configurable worker parameters
The cloud-config SHALL support configuration of worker parameters without modifying the cloud-config file itself.

#### Scenario: Custom thread count
- **WHEN** WORKER_THREADS is set in the cloud-config variables
- **THEN** the worker SHALL use that number of threads

#### Scenario: Default thread count
- **WHEN** WORKER_THREADS is not specified
- **THEN** the worker SHALL default to using all available CPU threads (auto)

### Requirement: Dedicated worker account
The documentation SHALL instruct users to create a dedicated non-admin OpenBench account for workers to avoid exposing admin credentials in cloud-config.

#### Scenario: Worker account usage
- **WHEN** setting up ephemeral workers
- **THEN** the user SHALL create a separate OpenBench account and use those credentials in the cloud-config

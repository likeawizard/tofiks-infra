## ADDED Requirements

### Requirement: Create worker servers when work is detected
The manager SHALL create Hetzner servers when active tests are detected and no managed workers currently exist.

#### Scenario: Work detected, no workers running
- **WHEN** active tests are detected AND no servers with the manager label exist
- **THEN** the manager SHALL create the configured number of Hetzner servers using the cloud-config as user-data

#### Scenario: Work detected, workers already running
- **WHEN** active tests are detected AND managed servers already exist
- **THEN** the manager SHALL NOT create additional servers

### Requirement: Destroy worker servers when work is complete
The manager SHALL destroy all managed Hetzner servers when no active tests remain, after a cooldown period.

#### Scenario: No work, cooldown expires
- **WHEN** no active tests are detected for the duration of the cooldown period (default 5 minutes)
- **THEN** the manager SHALL delete all Hetzner servers with the manager label

#### Scenario: Work reappears during cooldown
- **WHEN** no active tests are detected but work reappears before cooldown expires
- **THEN** the manager SHALL cancel the pending destruction

### Requirement: Track servers via labels
The manager SHALL use Hetzner labels to track servers it created, enabling stateless operation.

#### Scenario: Server creation
- **WHEN** creating a new Hetzner server
- **THEN** the server SHALL be labeled with `managed-by=worker-manager`

#### Scenario: Manager restart
- **WHEN** the manager process restarts
- **THEN** it SHALL discover existing managed servers by querying for the label

### Requirement: Configurable server parameters
The manager SHALL support configuration of Hetzner server parameters via environment variables.

#### Scenario: Server type configuration
- **WHEN** HCLOUD_SERVER_TYPE is set (e.g., "cpx31")
- **THEN** created servers SHALL use that server type

#### Scenario: Worker count configuration
- **WHEN** HCLOUD_WORKER_COUNT is set (e.g., "3")
- **THEN** the manager SHALL create that many servers when work is detected

#### Scenario: Default values
- **WHEN** optional parameters are not set
- **THEN** the manager SHALL use defaults: server type "cpx31", worker count 1, location "fsn1"

### Requirement: Cloud-config as user-data
The manager SHALL read the worker-cloud-config.yml file and pass it as user-data when creating servers.

#### Scenario: Cloud-config loaded
- **WHEN** the manager starts
- **THEN** it SHALL read worker-cloud-config.yml from the configured path and use it for all server creation

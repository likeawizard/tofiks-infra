## ADDED Requirements

### Requirement: Detect active OpenBench tests
The manager SHALL poll the OpenBench instance to determine whether there are active tests that need worker threads.

#### Scenario: Active tests present
- **WHEN** the OpenBench index page contains one or more active tests
- **THEN** the manager SHALL report that work is available

#### Scenario: No active tests
- **WHEN** the OpenBench index page contains no active tests (only pending, awaiting, or completed)
- **THEN** the manager SHALL report that no work is available

#### Scenario: OpenBench unreachable
- **WHEN** the manager cannot reach the OpenBench instance
- **THEN** the manager SHALL log the error and retry on the next poll cycle without changing state

### Requirement: Configurable poll interval
The manager SHALL poll at a configurable interval, defaulting to 60 seconds.

#### Scenario: Default poll interval
- **WHEN** no poll interval is configured
- **THEN** the manager SHALL poll every 60 seconds

#### Scenario: Custom poll interval
- **WHEN** POLL_INTERVAL is set to a duration (e.g., "30s", "2m")
- **THEN** the manager SHALL use that interval between polls

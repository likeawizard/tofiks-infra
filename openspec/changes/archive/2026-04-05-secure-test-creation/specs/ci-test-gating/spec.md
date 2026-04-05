## ADDED Requirements

### Requirement: CI uses non-admin OpenBench account
The openbench-test workflow SHALL authenticate with a dedicated non-admin account that does not have approver privileges.

#### Scenario: Test created by CI
- **WHEN** the workflow creates a test via the OpenBench API
- **THEN** it SHALL use `OPENBENCH_CI_USER` and `OPENBENCH_CI_PASS` secrets instead of admin credentials

### Requirement: CI-created tests require approval
Tests created by the CI account SHALL be in pending (unapproved) state, requiring manual approval before workers execute them.

#### Scenario: Test lands in pending state
- **WHEN** a test is created by the non-approver CI account
- **THEN** the test SHALL appear in the "Pending" section of the OpenBench index and NOT be picked up by workers

#### Scenario: Admin approves test
- **WHEN** an admin approves a pending test on the OpenBench UI
- **THEN** the test SHALL move to "Active" and workers SHALL begin processing it

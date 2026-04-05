## ADDED Requirements

### Requirement: PR comment links to specific test
The PR comment posted by the workflow SHALL include a direct link to the created test.

#### Scenario: Test created successfully
- **WHEN** the OpenBench API returns a Location header with the test ID
- **THEN** the PR comment SHALL include a link in the format `https://bench.likeawizard.dev/test/{id}/`

### Requirement: PGN uploads enabled
The test creation API call SHALL request PGN uploads so full game records are preserved.

#### Scenario: Test creation parameters
- **WHEN** the workflow POSTs to the OpenBench `/scripts/` endpoint
- **THEN** the `upload_pgns` parameter SHALL be set to `TRUE`

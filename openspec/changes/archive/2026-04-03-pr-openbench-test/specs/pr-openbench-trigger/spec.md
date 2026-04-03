## ADDED Requirements

### Requirement: OpenBench test is created when a PR is opened
A GitHub Actions workflow SHALL create an OpenBench SPRT test when a pull request is opened on the tofiks repo, comparing the PR branch against master.

#### Scenario: Non-draft PR opened triggers test creation
- **WHEN** a non-draft pull request is opened on the tofiks repo
- **THEN** the workflow builds both branches, extracts bench values, and creates an OpenBench test via the /scripts/ endpoint

#### Scenario: Draft PR is ignored
- **WHEN** a draft pull request is opened or updated
- **THEN** no OpenBench test is created

#### Scenario: Draft marked ready triggers test creation
- **WHEN** a draft pull request is marked as ready for review
- **THEN** the workflow creates an OpenBench test

#### Scenario: Non-draft PR updated triggers new test
- **WHEN** new commits are pushed to a non-draft pull request
- **THEN** a new OpenBench test is created for the updated code

### Requirement: Bench values are extracted from builds
The workflow SHALL build tofiks for both the PR branch and master, run `./tofiks bench`, and parse the node count to use as bench values in the test creation request.

#### Scenario: Bench values extracted successfully
- **WHEN** the workflow builds both branches
- **THEN** the bench output matching `<N> nodes` is parsed and used for dev_bench and base_bench

### Requirement: PR comment with test link
The workflow SHALL post a comment on the pull request with a link to the created OpenBench test.

#### Scenario: Comment posted after test creation
- **WHEN** the OpenBench test is created successfully
- **THEN** a comment appears on the PR with the test URL

### Requirement: Credentials stored as secrets
OpenBench username and password SHALL be stored as GitHub Secrets on the tofiks repo, never hardcoded in the workflow.

#### Scenario: Workflow authenticates with OpenBench
- **WHEN** the workflow sends a POST to /scripts/
- **THEN** it uses OPENBENCH_USER and OPENBENCH_PASS secrets for authentication

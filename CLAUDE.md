# Project Guidelines

## OpenSpec Workflow

When working with OpenSpec changes, follow this git workflow automatically:

### On `/opsx:propose` (after creating the change)
1. Create a new git branch named after the change: `git checkout -b <change-name>`
2. Commit the openspec artifacts with message: `openspec: propose <change-name>`

### On `/opsx:apply` completion (all tasks done)
1. Run `/opsx:archive` to archive the change before the final commit
2. Create a final commit with all implementation changes
3. Push the branch and create a PR using `gh pr create`

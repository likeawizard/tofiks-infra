# Project Guidelines

## Common Pitfalls

Things that have bitten us before — check these when making changes:

### OpenBench: /test/ vs /tune/
OpenBench uses different URL paths for SPRT tests (`/test/<id>/`) and SPSA tunes (`/tune/<id>/`). Any code that parses the index page or builds links must handle both. This includes:
- `worker-manager/openbench.go` — active work detection
- Workflow files — tune ID extraction and PR comment links

### Submodule URLs must be HTTPS
The hcloud workers clone tofiks-infra via cloud-init on fresh machines with no SSH keys. If any submodule uses `git@github.com:` URLs, the clone will silently fail and workers will never connect to OpenBench. Always use `https://github.com/` in `.gitmodules`.

### Deploy workflow: env block vs envs list
The `appleboy/ssh-action` requires secrets in **two** places in `deploy.yml`:
1. The `env:` block (maps GitHub secrets to env vars)
2. The `envs:` field (forwards those env vars to the remote script)

Missing either one means the value is empty on the server. When adding a new secret, update both.

### OpenBench fork (likeawizard/OpenBench)
We maintain a fork of `AndyGrant/OpenBench` with Tofiks-specific customizations baked in (engine config, Django settings, `CREATE_TUNE` API handler). Changes to OpenBench config go in the fork directly — there is no config overlay mechanism.

### Hcloud workers: no recovery from failed cloud-init
If workers fail during boot (e.g., bad submodule URL, missing credentials), they won't self-heal. They must be manually deleted from Hetzner so the worker-manager creates fresh ones.

## OpenSpec Workflow

When working with OpenSpec changes, follow this git workflow automatically:

### On `/opsx:propose` (after creating the change)
1. Create a new git branch named after the change: `git checkout -b <change-name>`
2. Commit the openspec artifacts with message: `openspec: propose <change-name>`

### On `/opsx:apply` completion (all tasks done)
1. Run `/opsx:archive` to archive the change before the final commit
2. Create a final commit with all implementation changes
3. Push the branch and create a PR using `gh pr create`

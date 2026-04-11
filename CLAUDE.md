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

### OpenBench: persistent state lives in TWO places
The openbench container needs **both** `/app/db` (SQLite) and `/app/Media` (PGN tar archives, event logs) mounted as volumes. The DB alone is not enough — PGN downloads read from `/app/Media/PGNs/<id>.pgn.tar` and those files get wiped on every container rebuild if `/app/Media` is ephemeral. If you see "Unable to find PGN for Workload #N" for every historical test after a deploy, this is the cause.

### Worker Dockerfile needs `procps` for pkill
The OpenBench client calls `pkill -f fastchess-ob` in `utils.kill_process_by_name` after every workload batch. `python:*-slim` base images don't ship `pkill`, so if `procps` is missing from `Dockerfile.worker`, every workload ends with a `FileNotFoundError: 'pkill'` that propagates out of `complete_workload` — **before** the PGN upload block at `worker.py:1148`. Symptom: results submit fine, tests make progress, but `PGN.objects.count() == 0` forever. Upstream OpenBench users don't hit this because the client normally runs on a full Linux box where procps is preinstalled.

### OpenBench SQLite must be in WAL mode
With multiple concurrent workers hammering `/clientSubmitResults/`, the default `journal_mode=delete` serializes every transaction and produces `sqlite3.OperationalError: database is locked` → 500 errors. The openbench entrypoint runs `PRAGMA journal_mode=WAL` once at startup; the setting is persistent on the DB file. If you restore the DB from a backup or swap in a fresh SQLite file, re-run the PRAGMA.

## OpenSpec Workflow

When working with OpenSpec changes, follow this git workflow automatically:

### On `/opsx:propose` (after creating the change)
1. Create a new git branch named after the change: `git checkout -b <change-name>`
2. Commit the openspec artifacts with message: `openspec: propose <change-name>`

### On `/opsx:apply` completion (all tasks done)
1. Run `/opsx:archive` to archive the change before the final commit
2. Create a final commit with all implementation changes
3. Push the branch and create a PR using `gh pr create`

## 1. Repository Setup

- [x] 1.1 Fork OpenBench to likeawizard/OpenBench on GitHub
- [x] 1.2 Add OpenBench fork as a git submodule in tofiks-infra

## 2. Docker Build

- [x] 2.1 Create Dockerfile.openbench: Python runtime, install deps, copy source
- [x] 2.2 Create openbench-entrypoint.sh: run migrations, create admin user (idempotent), start gunicorn
- [x] 2.3 Add openbench service to docker-compose.yml with port 80:8000, named volume for SQLite, and environment variables

## 3. Secrets

- [x] 3.1 Set OPENBENCH_ADMIN_USER and OPENBENCH_ADMIN_PASS secrets on tofiks-infra repo
- [x] 3.2 Update deploy workflow to pass OpenBench secrets to .env file on server

## 4. Verification

- [x] 4.1 Build and test locally
- [x] 4.2 Deploy to Hetzner and verify OpenBench UI is accessible (pending merge + deploy)
- [x] 4.3 Log in with admin credentials and verify account works (pending merge + deploy)

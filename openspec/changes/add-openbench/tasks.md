## 1. Repository Setup

- [ ] 1.1 Fork OpenBench to likeawizard/OpenBench on GitHub
- [ ] 1.2 Add OpenBench fork as a git submodule in tofiks-infra

## 2. Docker Build

- [ ] 2.1 Create Dockerfile.openbench: Python runtime, install deps, copy source
- [ ] 2.2 Create openbench-entrypoint.sh: run migrations, create admin user (idempotent), start gunicorn
- [ ] 2.3 Add openbench service to docker-compose.yml with port 80:8000, named volume for SQLite, and environment variables

## 3. Secrets

- [ ] 3.1 Set OPENBENCH_ADMIN_USER and OPENBENCH_ADMIN_PASS secrets on tofiks-infra repo
- [ ] 3.2 Update deploy workflow to pass OpenBench secrets to .env file on server

## 4. Verification

- [ ] 4.1 Build and test locally
- [ ] 4.2 Deploy to Hetzner and verify OpenBench UI is accessible at http://188.34.201.182
- [ ] 4.3 Log in with admin credentials and verify account works

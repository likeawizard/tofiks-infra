## 1. Worker Docker Compose

- [x] 1.1 Create `docker-compose.worker.yml` with a single worker service using Dockerfile.worker, `--fleet` flag, auto thread detection, and environment variables from `.env`

## 2. Cloud-Config Script

- [x] 2.1 Create `worker-cloud-config.yml` cloud-init script that clones the repo, writes `.env` from inline variables, builds and runs the worker via docker-compose, and shuts down the server when the worker exits

## 3. Documentation and Examples

- [x] 3.1 Create `.env.worker.example` documenting all required and optional variables
- [x] 3.2 Add usage instructions to the cloud-config file as comments (create worker account, paste into Hetzner, configure variables, delete server after use)

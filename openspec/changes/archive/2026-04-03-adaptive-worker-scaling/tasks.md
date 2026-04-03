## 1. Server Worker Service

- [x] 1.1 Add server-worker service to docker-compose.yml (always on, not under worker profile)
- [x] 1.2 Update deploy workflow to include server-worker and worker-scaler in deploy targets

## 2. Scaling Script

- [x] 2.1 Create scaling script that checks tofiks process count in lichess-bot container
- [x] 2.2 Script adjusts server-worker CPU limit via `docker update --cpus` (no restart needed)
- [x] 2.3 Debounce logic: 30s stable state before scaling
- [x] 2.4 Sidecar container (docker:cli) with Docker socket access

## 3. Verification

- [x] 3.1 Deploy to Hetzner and verify server-worker runs with 2 CPUs when idle (pending merge)
- [x] 3.2 Start a lichess game and verify worker scales down to 1 CPU (pending merge)
- [x] 3.3 End the game and verify worker scales back up to 2 CPUs (pending merge)

## 1. Server Worker Service

- [ ] 1.1 Add server-worker service to docker-compose.yml (always on, not under worker profile)
- [ ] 1.2 Update deploy workflow to include server-worker in the "all" and "bench" deploy targets

## 2. Scaling Script

- [ ] 2.1 Create scaling script that checks tofiks process count in lichess-bot container
- [ ] 2.2 Script restarts server-worker with 1 or 2 threads based on game activity
- [ ] 2.3 Add debounce logic to avoid thrashing on brief game starts/stops
- [ ] 2.4 Create sidecar container for the scaling script with Docker socket access

## 3. Verification

- [ ] 3.1 Deploy to Hetzner and verify server-worker runs with 2 threads when idle
- [ ] 3.2 Start a lichess game and verify worker scales down to 1 thread
- [ ] 3.3 End the game and verify worker scales back up to 2 threads

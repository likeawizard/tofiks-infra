## 1. Go Module Setup

- [ ] 1.1 Initialize Go module in `worker-manager/` with hcloud-go dependency
- [ ] 1.2 Create `worker-manager/Dockerfile` (Go build + minimal runtime image)

## 2. Work Detection

- [ ] 2.1 Implement OpenBench poller: HTTP GET to index page, parse HTML to detect active tests
- [ ] 2.2 Implement configurable poll loop with interval (default 60s)

## 3. Server Lifecycle

- [ ] 3.1 Implement server creation: read cloud-config file, create labeled Hetzner servers via hcloud-go
- [ ] 3.2 Implement server discovery: list servers by `managed-by=worker-manager` label
- [ ] 3.3 Implement server destruction: delete all managed servers after cooldown period
- [ ] 3.4 Implement state machine: idle ↔ workers-running with cooldown logic

## 4. Configuration and Main

- [ ] 4.1 Implement config from environment variables (HCLOUD_TOKEN, HCLOUD_SERVER_TYPE, HCLOUD_WORKER_COUNT, HCLOUD_LOCATION, POLL_INTERVAL, COOLDOWN, OPENBENCH_URL, CLOUD_CONFIG_PATH)
- [ ] 4.2 Wire up main.go: config → poll loop → state machine → lifecycle actions

## 5. Docker Compose Integration

- [ ] 5.1 Add worker-manager service to docker-compose.yml with required env vars and volume mount for cloud-config
- [ ] 5.2 Update .env and .env.worker.example with new variables

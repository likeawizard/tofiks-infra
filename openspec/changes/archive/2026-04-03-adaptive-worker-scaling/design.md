## Context

The Hetzner server has 4 vCPUs. The lichess-bot runs with concurrency=2, meaning up to 2 tofiks engine processes can be active simultaneously (one per game). Each engine process uses ~1 vCPU. An OpenBench worker thread also uses ~1 vCPU per active game pair.

```
4 vCPUs budget:
- Lichess games active:  2 vCPUs for bot + 1 vCPU for worker = 3 used
- No lichess games:      0 vCPUs for bot + 2 vCPUs for worker = 2 used
- Always 1 vCPU spare for OS, OpenBench server, Caddy
```

## Goals / Non-Goals

**Goals:**
- Run an OpenBench worker on the server 24/7 so tests progress without the laptop
- Scale worker threads based on lichess-bot activity (1 or 2 threads)
- Simple, reliable detection of active games

**Non-Goals:**
- Fine-grained CPU pinning or cgroup limits
- Scaling beyond 2 worker threads
- Replacing the laptop worker (it's complementary)

## Decisions

### 1. Detect games by counting tofiks processes in the lichess-bot container

Run `docker exec lichess-bot pgrep -c tofiks` (or equivalent) to count engine processes. If >0, games are active. This is simpler and more reliable than polling the lichess API — no auth needed, no rate limits, instant response.

### 2. Scaling script as a sidecar container

A small shell script runs in a loop: check game activity, restart the worker with the appropriate thread count if it changed. Runs as its own docker-compose service with access to the Docker socket.

**Why not cron:** The check needs to happen frequently (every 10-30 seconds) and cron's minimum is 1 minute. A simple sleep loop is more appropriate.

**Why sidecar over modifying the worker:** The OpenBench worker doesn't support dynamic thread changes at runtime. The only way to change threads is to restart it with a different `-T` value.

### 3. Worker restart strategy

When scaling up/down, the script stops the current worker and starts a new one with the new thread count. The worker will finish its current game batch before stopping (graceful). To avoid thrashing, the script should only change threads after the game count has been stable for a short period.

### 4. Separate server-worker service (not using worker profile)

The server worker is a distinct service from the laptop worker. It runs without the `worker` profile so it starts automatically on deploy. The laptop worker remains under the `worker` profile.

## Risks / Trade-offs

- **Worker restart interrupts games** → The current batch of OpenBench games is lost on restart. At workload_size=32 with STC, that's a few minutes of work. Acceptable for occasional scaling events.
- **Docker socket access** → The sidecar needs the Docker socket to exec into the bot container and restart the worker. This is a privileged operation but acceptable on a single-user server.
- **Thrashing** → Rapid game start/stop could cause frequent worker restarts. Mitigate with a stability delay before scaling.

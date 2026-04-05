## Context

OpenBench runs on a single Hetzner server. Ephemeral workers can be provisioned via cloud-config, but creating/destroying servers is manual. The worker-manager automates this lifecycle.

OpenBench's index page shows tests in four states: pending (not yet approved), active (running), awaiting (paused), and completed. The manager cares about **active** tests — those are the ones that need worker threads.

The OpenBench index page is rendered server-side with Django templates. There is no REST API that directly returns test counts. However, the manager runs on the same Docker network as OpenBench, so it can hit `http://openbench:8000/` directly.

## Goals / Non-Goals

**Goals:**
- Automatically create Hetzner workers when OpenBench has active tests
- Automatically destroy workers when all tests are finished
- Simple, single-binary Go service
- Configurable server type, worker count, and poll interval
- Run as a docker-compose service alongside OpenBench
- Also runnable locally as a plain Go binary (just set env vars, no Docker required)

**Non-Goals:**
- Dynamic scaling (adjusting worker count based on workload size) — fixed count of workers
- Supporting multiple Hetzner projects or locations
- Web UI or API for the manager itself
- Graceful mid-test worker removal (workers finish their current batch naturally via --fleet)

## Decisions

### 1. Work detection: scrape index page vs query DB vs custom API

**Choice**: Parse the OpenBench index HTML page for active test entries.

**Rationale**: No modifications to OpenBench needed. The index page reliably shows active tests in a table. The OpenBench URL is configurable via `OPENBENCH_URL` — defaults to `http://openbench:8000` when running in docker-compose, but can be set to `https://bench.likeawizard.dev` for local runs. Parsing HTML is simple enough for this — we just need to detect whether there are any active tests, not parse their details.

**Alternative considered**: Query SQLite directly — too coupled, breaks if OpenBench changes schema. Custom API endpoint — requires modifying OpenBench (a submodule we don't control).

### 2. State machine: two states (idle / workers-running)

**Choice**: Simple two-state loop:
- **Idle**: Poll for work. If active tests found → create servers, transition to workers-running.
- **Workers-running**: Poll for work. If no active tests → destroy all managed servers, transition to idle.

**Rationale**: Minimizes complexity. No need for intermediate states. The cloud-config handles all worker setup. The --fleet flag handles worker exit. The manager only needs to create and destroy servers.

### 3. Server tracking: label-based

**Choice**: Tag created servers with a label (e.g., `managed-by=worker-manager`). On each poll cycle, list servers with that label to know what we're managing.

**Rationale**: Stateless — the manager can restart without losing track of servers. Hetzner API supports label filtering natively. No local state file needed.

### 4. Cloud-config: read from file

**Choice**: Read `worker-cloud-config.yml` from the filesystem at startup, use it as user-data when creating servers.

**Rationale**: Reuses the existing cloud-config. Single source of truth. Changes to cloud-config are picked up on manager restart.

### 5. Cooldown before destroying

**Choice**: After detecting no active work, wait for a configurable cooldown period (default: 5 minutes) before destroying servers. If work reappears during cooldown, cancel the destruction.

**Rationale**: Avoids destroying and recreating servers for rapid test submissions. Also gives workers time to finish submitting their last batch of results.

## Risks / Trade-offs

- **[HTML parsing fragility]** → If OpenBench changes its index page layout, parsing breaks. Mitigation: keep the parser minimal (just detect presence of active tests), and pin the OpenBench submodule version.
- **[Race: work finishes during server creation]** → Servers take ~3 min to boot and build. If work finishes in that window, the workers will start, find no work (--fleet), exit, and the manager will destroy them on next poll. Minor cost waste, acceptable.
- **[Hetzner API rate limits]** → Polling every 60s with a few API calls is well within limits. Not a concern.
- **[Manager crash while servers running]** → On restart, the manager discovers existing servers via labels and resumes management. No orphaned servers.

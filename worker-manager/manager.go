package main

import (
	"context"
	"log"
	"time"
)

type state int

const (
	stateIdle state = iota
	stateWorkersRunning
	stateCoolingDown
)

func (s state) String() string {
	switch s {
	case stateIdle:
		return "idle"
	case stateWorkersRunning:
		return "workers-running"
	case stateCoolingDown:
		return "cooling-down"
	default:
		return "unknown"
	}
}

type Manager struct {
	config       Config
	hcloud       *HCloudManager
	state        state
	cooldownFrom time.Time
}

func NewManager(cfg Config, hcloud *HCloudManager) *Manager {
	return &Manager{
		config: cfg,
		hcloud: hcloud,
		state:  stateIdle,
	}
}

func (m *Manager) Run(ctx context.Context) {
	log.Printf("Starting worker manager (poll=%s, cooldown=%s, workers=%d, type=%s, location=%s)",
		m.config.PollInterval, m.config.Cooldown, m.config.WorkerCount, m.config.ServerType, m.config.Location)

	// Check for existing managed servers on startup
	servers, err := m.hcloud.ManagedServers(ctx)
	if err != nil {
		log.Printf("Error checking existing servers: %v", err)
	} else if len(servers) > 0 {
		log.Printf("Found %d existing managed servers, entering workers-running state", len(servers))
		m.state = stateWorkersRunning
	}

	ticker := time.NewTicker(m.config.PollInterval)
	defer ticker.Stop()

	// Run immediately on start, then on ticker
	m.tick(ctx)

	for {
		select {
		case <-ctx.Done():
			log.Println("Shutting down worker manager")
			return
		case <-ticker.C:
			m.tick(ctx)
		}
	}
}

func (m *Manager) tick(ctx context.Context) {
	hasWork, err := HasActiveWork(m.config.OpenBenchURL)
	if err != nil {
		log.Printf("Error polling OpenBench: %v", err)
		return
	}

	prevState := m.state

	switch m.state {
	case stateIdle:
		if hasWork {
			log.Println("Active tests detected, creating workers")
			if err := m.hcloud.CreateWorkers(ctx); err != nil {
				log.Printf("Error creating workers: %v", err)
				return
			}
			m.state = stateWorkersRunning
		}

	case stateWorkersRunning:
		if !hasWork {
			log.Printf("No active tests, starting cooldown (%s)", m.config.Cooldown)
			m.cooldownFrom = time.Now()
			m.state = stateCoolingDown
		} else {
			// Top up if we have fewer workers than desired
			if err := m.hcloud.CreateWorkers(ctx); err != nil {
				log.Printf("Error topping up workers: %v", err)
			}
		}

	case stateCoolingDown:
		if hasWork {
			log.Println("Work reappeared during cooldown, cancelling destruction")
			m.state = stateWorkersRunning
		} else if time.Since(m.cooldownFrom) >= m.config.Cooldown {
			log.Println("Cooldown expired, destroying workers")
			if err := m.hcloud.DestroyWorkers(ctx); err != nil {
				log.Printf("Error destroying workers: %v", err)
				return
			}
			m.state = stateIdle
		} else {
			remaining := m.config.Cooldown - time.Since(m.cooldownFrom)
			log.Printf("Cooling down, %s remaining before destruction", remaining.Round(time.Second))
		}
	}

	if m.state != prevState {
		log.Printf("State: %s → %s", prevState, m.state)
	}
}

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

const managedByLabel = "managed-by"
const managedByValue = "worker-manager"

type HCloudManager struct {
	client    *hcloud.Client
	config    Config
	userData  string
	sshKeyID  int64
}

func NewHCloudManager(cfg Config) (*HCloudManager, error) {
	rawUserData, err := os.ReadFile(cfg.CloudConfigPath)
	if err != nil {
		return nil, fmt.Errorf("reading cloud-config: %w", err)
	}

	// Inject worker credentials into cloud-config
	userData := string(rawUserData)
	userData = strings.ReplaceAll(userData, "OPENBENCH_USERNAME=worker", "OPENBENCH_USERNAME="+cfg.OpenBenchUsername)
	userData = strings.ReplaceAll(userData, "OPENBENCH_PASSWORD=CHANGE_ME", "OPENBENCH_PASSWORD="+cfg.OpenBenchPassword)

	client := hcloud.NewClient(hcloud.WithToken(cfg.HCloudToken))

	m := &HCloudManager{
		client:   client,
		config:   cfg,
		userData: userData,
	}

	if cfg.SSHKeyName != "" {
		key, _, err := client.SSHKey.GetByName(context.Background(), cfg.SSHKeyName)
		if err != nil {
			return nil, fmt.Errorf("looking up SSH key %q: %w", cfg.SSHKeyName, err)
		}
		if key == nil {
			return nil, fmt.Errorf("SSH key %q not found", cfg.SSHKeyName)
		}
		m.sshKeyID = key.ID
	}

	return m, nil
}

func (m *HCloudManager) ManagedServers(ctx context.Context) ([]*hcloud.Server, error) {
	servers, err := m.client.Server.AllWithOpts(ctx, hcloud.ServerListOpts{
		ListOpts: hcloud.ListOpts{
			LabelSelector: fmt.Sprintf("%s=%s", managedByLabel, managedByValue),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("listing managed servers: %w", err)
	}
	return servers, nil
}

func (m *HCloudManager) CreateWorkers(ctx context.Context) error {
	existing, err := m.ManagedServers(ctx)
	if err != nil {
		return err
	}
	needed := m.config.WorkerCount - len(existing)
	if needed <= 0 {
		log.Printf("Workers already running (%d/%d servers)", len(existing), m.config.WorkerCount)
		return nil
	}
	if len(existing) > 0 {
		log.Printf("Have %d/%d workers, creating %d more", len(existing), m.config.WorkerCount, needed)
	}

	image, _, err := m.client.Image.GetByNameAndArchitecture(ctx, "docker-ce", "x86")
	if err != nil {
		return fmt.Errorf("looking up docker-ce image: %w", err)
	}
	if image == nil {
		return fmt.Errorf("docker-ce image not found")
	}

	serverType, _, err := m.client.ServerType.GetByName(ctx, m.config.ServerType)
	if err != nil {
		return fmt.Errorf("looking up server type %q: %w", m.config.ServerType, err)
	}
	if serverType == nil {
		return fmt.Errorf("server type %q not found", m.config.ServerType)
	}

	location, _, err := m.client.Location.GetByName(ctx, m.config.Location)
	if err != nil {
		return fmt.Errorf("looking up location %q: %w", m.config.Location, err)
	}
	if location == nil {
		return fmt.Errorf("location %q not found", m.config.Location)
	}

	batch := time.Now().Format("0102-1504")
	for i := range needed {
		name := fmt.Sprintf("ob-worker-%s-%d", batch, i+1)

		opts := hcloud.ServerCreateOpts{
			Name:       name,
			ServerType: serverType,
			Image:      image,
			Location:   location,
			UserData:   m.userData,
			Labels: map[string]string{
				managedByLabel: managedByValue,
			},
		}

		if m.sshKeyID != 0 {
			opts.SSHKeys = []*hcloud.SSHKey{{ID: m.sshKeyID}}
		}

		result, _, err := m.client.Server.Create(ctx, opts)
		if err != nil {
			return fmt.Errorf("creating server %s: %w", name, err)
		}
		log.Printf("Created server %s (ID: %d)", name, result.Server.ID)
	}

	return nil
}

func (m *HCloudManager) DestroyWorkers(ctx context.Context) error {
	servers, err := m.ManagedServers(ctx)
	if err != nil {
		return err
	}

	var failed int
	for _, server := range servers {
		_, _, err := m.client.Server.DeleteWithResult(ctx, server)
		if err != nil {
			log.Printf("Error deleting server %s (ID: %d): %v", server.Name, server.ID, err)
			failed++
			continue
		}
		log.Printf("Deleted server %s (ID: %d)", server.Name, server.ID)
	}

	if failed > 0 {
		return fmt.Errorf("failed to delete %d of %d servers", failed, len(servers))
	}
	return nil
}

package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	HCloudToken      string
	ServerType       string
	WorkerCount      int
	Location         string
	PollInterval     time.Duration
	Cooldown         time.Duration
	OpenBenchURL     string
	CloudConfigPath  string
	SSHKeyName       string
	OpenBenchUsername string
	OpenBenchPassword string
}

func LoadConfig() (Config, error) {
	token := os.Getenv("HCLOUD_TOKEN")
	if token == "" {
		return Config{}, fmt.Errorf("HCLOUD_TOKEN is required")
	}

	obUser := os.Getenv("OPENBENCH_USERNAME")
	obPass := os.Getenv("OPENBENCH_PASSWORD")
	if obUser == "" || obPass == "" {
		return Config{}, fmt.Errorf("OPENBENCH_USERNAME and OPENBENCH_PASSWORD are required")
	}

	cloudConfigPath := envOrDefault("CLOUD_CONFIG_PATH", "worker-cloud-config.yml")
	if _, err := os.Stat(cloudConfigPath); err != nil {
		return Config{}, fmt.Errorf("cloud-config file not found: %s", cloudConfigPath)
	}

	workerCount := 1
	if s := os.Getenv("HCLOUD_WORKER_COUNT"); s != "" {
		n, err := strconv.Atoi(s)
		if err != nil || n < 1 {
			return Config{}, fmt.Errorf("HCLOUD_WORKER_COUNT must be a positive integer")
		}
		workerCount = n
	}

	pollInterval, err := parseDuration("POLL_INTERVAL", 60*time.Second)
	if err != nil {
		return Config{}, err
	}

	cooldown, err := parseDuration("COOLDOWN", 5*time.Minute)
	if err != nil {
		return Config{}, err
	}

	return Config{
		HCloudToken:     token,
		ServerType:      envOrDefault("HCLOUD_SERVER_TYPE", "cpx32"),
		WorkerCount:     workerCount,
		Location:        envOrDefault("HCLOUD_LOCATION", "fsn1"),
		PollInterval:    pollInterval,
		Cooldown:        cooldown,
		OpenBenchURL:    envOrDefault("OPENBENCH_URL", "http://openbench:8000"),
		CloudConfigPath: cloudConfigPath,
		SSHKeyName:       os.Getenv("HCLOUD_SSH_KEY"),
		OpenBenchUsername: obUser,
		OpenBenchPassword: obPass,
	}, nil
}

func envOrDefault(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func parseDuration(envKey string, defaultVal time.Duration) (time.Duration, error) {
	s := os.Getenv(envKey)
	if s == "" {
		return defaultVal, nil
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return 0, fmt.Errorf("%s: invalid duration %q", envKey, s)
	}
	return d, nil
}

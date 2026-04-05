package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	hcloudMgr, err := NewHCloudManager(cfg)
	if err != nil {
		log.Fatalf("HCloud setup error: %v", err)
	}

	mgr := NewManager(cfg, hcloudMgr)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	mgr.Run(ctx)
}

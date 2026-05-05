// Package main is the entry point for the driftwatch daemon.
// It wires together configuration loading, Docker client setup,
// compose spec parsing, drift detection, scheduling, and notifications.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/user/driftwatch/internal/compose"
	"github.com/user/driftwatch/internal/config"
	"github.com/user/driftwatch/internal/docker"
	"github.com/user/driftwatch/internal/drift"
	"github.com/user/driftwatch/internal/notify"
	"github.com/user/driftwatch/internal/schedule"
)

func main() {
	configPath := flag.String("config", "driftwatch.yaml", "path to driftwatch config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Parse the compose spec declared in config.
	spec, err := compose.LoadSpec(cfg.ComposePath)
	if err != nil {
		log.Fatalf("failed to load compose spec: %v", err)
	}

	// Connect to the Docker daemon.
	dockerClient, err := docker.NewClient(cfg.DockerHost)
	if err != nil {
		log.Fatalf("failed to create docker client: %v", err)
	}
	defer dockerClient.Close()

	// Build notifier from configured level.
	level, err := notify.ParseLevel(cfg.NotifyLevel)
	if err != nil {
		log.Fatalf("invalid notify level %q: %v", cfg.NotifyLevel, err)
	}
	notifier := notify.New(level, os.Stdout)

	// Set up the scheduler.
	sched := schedule.New(schedule.Config{
		Interval: cfg.Interval,
	})

	// Trap OS signals for graceful shutdown.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigCh
		log.Printf("received signal %s, shutting down", sig)
		cancel()
	}()

	log.Printf("driftwatch started (compose=%s, interval=%s, notify=%s)",
		cfg.ComposePath, cfg.Interval, cfg.NotifyLevel)

	// Run the detection loop.
	sched.Run(ctx, func() {
		containers, err := dockerClient.ListContainers(ctx)
		if err != nil {
			log.Printf("error listing containers: %v", err)
			return
		}

		results := drift.Detect(spec, containers)

		if err := notifier.Notify(results); err != nil {
			log.Printf("error sending notification: %v", err)
		}
	})

	fmt.Println("driftwatch stopped")
}

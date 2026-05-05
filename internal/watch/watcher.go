package watch

import (
	"context"
	"log"
	"time"

	"driftwatch/internal/compose"
	"driftwatch/internal/docker"
	"driftwatch/internal/drift"
)

// Config holds configuration for the watcher.
type Config struct {
	ComposePath string
	Interval    time.Duration
	Quiet       bool
}

// Watcher periodically checks for configuration drift.
type Watcher struct {
	cfg    Config
	client *docker.Client
}

// New creates a new Watcher with the given config.
func New(cfg Config, client *docker.Client) *Watcher {
	return &Watcher{cfg: cfg, client: client}
}

// Run starts the watch loop, blocking until ctx is cancelled.
func (w *Watcher) Run(ctx context.Context) error {
	ticker := time.NewTicker(w.cfg.Interval)
	defer ticker.Stop()

	if err := w.check(); err != nil {
		log.Printf("[driftwatch] initial check error: %v", err)
	}

	for {
		select {
		case <-ticker.C:
			if err := w.check(); err != nil {
				log.Printf("[driftwatch] check error: %v", err)
			}
		case <-ctx.Done():
			log.Println("[driftwatch] shutting down")
			return ctx.Err()
		}
	}
}

func (w *Watcher) check() error {
	spec, err := compose.LoadSpec(w.cfg.ComposePath)
	if err != nil {
		return err
	}

	containers, err := w.client.ListContainers()
	if err != nil {
		return err
	}

	results := drift.Detect(spec, containers)
	summary := drift.Summarize(results)

	if !w.cfg.Quiet || summary.DriftCount > 0 {
		log.Printf("[driftwatch] %s", summary.OneLiner())
	}

	return nil
}

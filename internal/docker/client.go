package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

// ContainerInfo holds relevant runtime info for a running container.
type ContainerInfo struct {
	ID      string
	Name    string
	Image   string
	Env     map[string]string
	Labels  map[string]string
	Project string
	Service string
}

// Client wraps the Docker SDK client.
type Client struct {
	dc *client.Client
}

// NewClient creates a new Docker client using environment defaults.
func NewClient() (*Client, error) {
	dc, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("docker: failed to create client: %w", err)
	}
	return &Client{dc: dc}, nil
}

// Close releases resources held by the Docker client.
func (c *Client) Close() error {
	return c.dc.Close()
}

// ListProjectContainers returns running containers belonging to the given
// Compose project (matched via the com.docker.compose.project label).
func (c *Client) ListProjectContainers(ctx context.Context, project string) ([]ContainerInfo, error) {
	f := filters.NewArgs()
	f.Add("label", fmt.Sprintf("com.docker.compose.project=%s", project))
	f.Add("status", "running")

	containers, err := c.dc.ContainerList(ctx, types.ContainerListOptions{Filters: f})
	if err != nil {
		return nil, fmt.Errorf("docker: list containers: %w", err)
	}

	result := make([]ContainerInfo, 0, len(containers))
	for _, ctr := range containers {
		result = append(result, toContainerInfo(ctr))
	}
	return result, nil
}

func toContainerInfo(ctr types.Container) ContainerInfo {
	name := ""
	if len(ctr.Names) > 0 {
		name = ctr.Names[0]
	}

	env := make(map[string]string)
	// Env values are not available from list; populated separately if needed.

	return ContainerInfo{
		ID:      ctr.ID,
		Name:    name,
		Image:   ctr.Image,
		Env:     env,
		Labels:  ctr.Labels,
		Project: ctr.Labels["com.docker.compose.project"],
		Service: ctr.Labels["com.docker.compose.service"],
	}
}

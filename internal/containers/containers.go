package containers

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/Lakelimbo/kaeru/internal/database"
	"github.com/Lakelimbo/kaeru/internal/logger"
	"github.com/docker/compose/v5/pkg/api"
	"github.com/moby/moby/client"
)

// Interface to place each method for container interactions.
type ContainerRepository interface {
	PullImage(ctx context.Context, image string) error

	Up(ctx context.Context, app database.App) error
	Down(ctx context.Context, name string) error

	ListContainers(ctx context.Context, project string) (client.ContainerListResult, error)
	ListServices(ctx context.Context, project string) (client.ServiceListResult, error)

	Logs(ctx context.Context, id string, follow bool) (io.ReadCloser, error)
	Stats(ctx context.Context, id string) (client.ContainerStartResult, error)
}

func (r *DockerRepository) PullImage(ctx context.Context, image string) error {
	reader, err := r.cli.Client().ImagePull(ctx, image, client.ImagePullOptions{})
	if err != nil {
		logger.Errorf("Failed to pull image: %v", err)
		return err
	}
	defer reader.Close()

	io.Copy(io.Discard, reader)
	return nil
}

func (r *DockerRepository) Up(ctx context.Context, app database.App) error {
	tmp, err := os.CreateTemp(".", fmt.Sprintf("*.compose.tmp-%s.yaml", app.ProjectName))
	if err != nil {
		logger.Errorf("Unable to create Compose declaration: %v", err)
		return err
	}

	if _, err := tmp.WriteString(app.ComposeYAML); err != nil {
		logger.Errorf("Unable to write Compose YAML: %v", err)
		os.Remove(tmp.Name())
		return err
	}
	if err := tmp.Sync(); err != nil {
		os.Remove(tmp.Name())
		return err
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmp.Name())
		return err
	}

	project, err := r.compose.LoadProject(ctx, api.ProjectLoadOptions{
		ProjectName: app.ProjectName,
		ConfigPaths: []string{tmp.Name()},
	})
	if err != nil {
		logger.Errorf("Unable to load Compose file: %v", err)
		os.Remove(tmp.Name())
		return err
	}

	defer os.Remove(tmp.Name())

	return r.compose.Up(ctx, project, api.UpOptions{})
}

func (r *DockerRepository) Down(ctx context.Context, name string) error {
	return r.compose.Down(ctx, name, api.DownOptions{})
}

func (r *DockerRepository) ListContainers(ctx context.Context, project string) (client.ContainerListResult, error) {
	filters := make(client.Filters)
	filters.Add("label", fmt.Sprintf("com.docker.compose.project=%s", project))

	return r.cli.Client().ContainerList(ctx, client.ContainerListOptions{
		All:     true,
		Filters: filters,
	})
}

// ListServices list each service ("things") created by a Compose stack.
func (r *DockerRepository) ListServices(ctx context.Context, project string) (client.ServiceListResult, error) {
	filters := make(client.Filters)
	filters.Add("label", fmt.Sprintf("com.docker.compose.project=%s", project))

	return r.cli.Client().ServiceList(ctx, client.ServiceListOptions{
		Filters: filters,
	})
}

func (r *DockerRepository) Logs(ctx context.Context, id string, follow bool) (io.ReadCloser, error) {
	return r.cli.Client().ContainerLogs(ctx, id, client.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     follow,
	})
}

func (r *DockerRepository) Stats(ctx context.Context, id string) (client.ContainerStatsResult, error) {
	return r.cli.Client().ContainerStats(ctx, id, client.ContainerStatsOptions{
		Stream: true,
	})
}

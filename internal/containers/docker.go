package containers

import (
	"github.com/Lakelimbo/kaeru/config"
	"github.com/Lakelimbo/kaeru/internal/logger"
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/flags"
	"github.com/docker/compose/v5/pkg/api"
	"github.com/docker/compose/v5/pkg/compose"
	"github.com/moby/moby/client"
	"github.com/pocketbase/dbx"
)

// A Docker repository instance, which will also spawn
// Compose.
type DockerRepository struct {
	cli     *command.DockerCli
	cfg     *config.Config
	db      *dbx.DB
	compose api.Compose
}

func NewDockerRepository(cfg *config.Config, db *dbx.DB) (*DockerRepository, error) {
	host := client.FromEnv
	if cfg.Server.DockerHost != "" {
		host = client.WithHost(cfg.Server.DockerHost)
	}

	cli, err := command.NewDockerCli(command.WithAPIClientOptions(host))
	if err != nil {
		logger.Fatalf("Unable to create Docker CLI instance: %v", err)
	}

	err = cli.Initialize(&flags.ClientOptions{})
	if err != nil {
		logger.Fatalf("Failed to initialize Docker CLI: %v", err)
	}

	compose := getCompose(cli)

	return &DockerRepository{
		cli:     cli,
		cfg:     cfg,
		db:      db,
		compose: compose,
	}, nil
}

func getCompose(cli *command.DockerCli) api.Compose {
	service, err := compose.NewComposeService(cli)
	if err != nil {
		logger.Fatalf("Failed to create Compose service: %v", err)
	}

	return service
}

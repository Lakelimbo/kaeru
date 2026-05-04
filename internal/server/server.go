package server

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Lakelimbo/kaeru/config"
	"github.com/Lakelimbo/kaeru/internal/containers"
	"github.com/Lakelimbo/kaeru/internal/jobs"
	"github.com/Lakelimbo/kaeru/internal/logger"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
)

// The "core" of Kaeru. Think of it as dependency injection,
// containing things that are required to work, or that is better
// to do a "global" initialization rather than a per-method one.
type App struct {
	Server *echo.Echo
	Config *config.Config
	DB     *dbx.DB
	Docker *containers.DockerRepository
	Jobs   *jobs.Job
}

func New(
	cfg *config.Config,
	db *dbx.DB,
	docker *containers.DockerRepository,
	jobs *jobs.Job,
) *App {
	e := echo.New()
	e.Logger = slog.New(logger.New())

	return &App{e, cfg, db, docker, jobs}
}

// Start starts the Echo server on the designated port.
func (app *App) Start(ctx context.Context) error {
	srv := echo.StartConfig{
		Address:    fmt.Sprintf("%s:%d", app.Config.Server.Host, app.Config.Server.Port),
		HidePort:   true,
		HideBanner: true,
	}

	logger.Infof("Server started on %s", srv.Address)
	return srv.Start(ctx, app.Server)
}

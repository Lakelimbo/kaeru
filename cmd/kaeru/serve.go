package main

import (
	"context"
	"fmt"

	"github.com/Lakelimbo/kaeru/config"
	"github.com/Lakelimbo/kaeru/internal/api"
	"github.com/Lakelimbo/kaeru/internal/auth"
	"github.com/Lakelimbo/kaeru/internal/containers"
	"github.com/Lakelimbo/kaeru/internal/database"
	"github.com/Lakelimbo/kaeru/internal/jobs"
	"github.com/Lakelimbo/kaeru/internal/logger"
	"github.com/Lakelimbo/kaeru/internal/server"
	"github.com/Lakelimbo/kaeru/tools/misc"
	"github.com/spf13/cobra"
)

var (
	serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "An open source, self-hosted server management tool in a single binary",
		PreRun: func(cmd *cobra.Command, args []string) {
			fmt.Println(misc.PrintLogo())
		},
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := config.Load(cfgPath)
			if err != nil {
				loadConfigErr()
			}

			db := database.New(cfg)
			auth.FirstRun(db)

			docker, _ := containers.NewDockerRepository(cfg, db)

			jobPubSub := jobs.NewPubSub()
			jobQueue := jobs.Spawn(db, docker, jobPubSub)

			app := server.New(cfg, db, docker, jobQueue)
			api.RegisterRoutes(app)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			if err := app.Start(ctx); err != nil {
				logger.Fatalf("Failed to start server: %v", err)
			}
		},
	}
)

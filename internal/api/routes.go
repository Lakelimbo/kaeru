package api

import (
	"strings"

	"github.com/Lakelimbo/kaeru/internal/api/handlers"
	"github.com/Lakelimbo/kaeru/internal/auth"
	"github.com/Lakelimbo/kaeru/internal/server"
	"github.com/Lakelimbo/kaeru/web"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func RegisterRoutes(app *server.App) {
	origins := app.Config.Server.CORSOrigins
	if app.Config.IsDev() {
		origins = append(origins, "http://localhost:5173") // vite frontend
	}

	app.Server.Use(frontend())
	app.Server.Use(server.CORS(origins))

	v1 := app.Server.Group("/api/v1")
	v1.Use(auth.JWTMiddleware())

	if app.Config.IsDev() || app.Config.Server.APISpecProduction {
		spec := handlers.NewSpecHandler(&app.Config.Server)
		spec.Routes(v1)
	}

	auth := handlers.NewAuthHandler(app.Config, app.DB)
	auth.Routes(v1)

	apps := handlers.NewAppHandler(app.Config, app.DB, app.Docker, app.Jobs)
	apps.Routes(v1)

	jobs := handlers.NewJobHandler(app.Jobs)
	jobs.Routes(v1)
}

func frontend() echo.MiddlewareFunc {
	return middleware.StaticWithConfig(middleware.StaticConfig{
		HTML5:      true,
		Root:       "build",
		Filesystem: web.Frontend,
		Skipper: func(c *echo.Context) bool {
			return strings.HasPrefix(c.Request().URL.Path, "/api")
		},
	})
}

package handlers

import (
	"net/http"

	"github.com/Lakelimbo/kaeru/config"
	"github.com/Lakelimbo/kaeru/internal/apps"
	"github.com/Lakelimbo/kaeru/internal/containers"
	"github.com/Lakelimbo/kaeru/internal/jobs"
	"github.com/Lakelimbo/kaeru/tools/errors"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
)

type AppHandler struct {
	config *config.Config
	db     *dbx.DB
	docker *containers.DockerRepository
	jobs   *jobs.Job
}

func NewAppHandler(
	cfg *config.Config,
	db *dbx.DB,
	docker *containers.DockerRepository,
	jobs *jobs.Job,
) *AppHandler {
	return &AppHandler{cfg, db, docker, jobs}
}

func (h *AppHandler) Routes(e *echo.Group) {
	e.POST("/apps", h.Create)
	e.GET("/apps", h.List)
	e.GET("/apps/:id", h.View)
}

// Create is the endpoint to create a new application manually with a Compose stack
//
//	@Summary		Create
//	@Description	Creates an app with Compose manually
//	@Tags			apps
//	@Accept			json
//	@Produce		json
//	@Param			request	body		apps.AppRequest	true	"Data to create the app stack"
//	@Success		200		{object}	apps.AppCreateResponse
//	@Failure		400		{object}	errors.Error
//	@Failure		401		{object}	errors.Error
//	@Failure		500		{object}	errors.Error
//	@Router			/api/v1/apps [post]
func (h *AppHandler) Create(c *echo.Context) error {
	var req apps.AppRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errors.Error{
			Identifier: errors.ErrInvalidBody,
			Message:    "Invalid request body",
		})
	}

	app := apps.New(h.db, h.docker, h.jobs)
	res, err := app.NewApp(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errors.Error{
			Identifier: errors.ErrAppCreation,
			Message:    "Failed to create app",
		})
	}

	return c.JSON(http.StatusOK, res)
}

// List is the endpoint to view all created apps
//
//	@Summary		List
//	@Description	Lists all the apps
//	@Tags			apps
//	@Produce		json
//	@Success		200	{array}		apps.AppResponse
//	@Failure		500	{object}	errors.Error
//	@Router			/api/v1/apps [get]
func (h *AppHandler) List(c *echo.Context) error {
	app := apps.New(h.db, h.docker, h.jobs)
	res, err := app.ListApps()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errors.Error{
			Identifier: errors.ErrInternal,
			Message:    "An unknown error has happened",
		})
	}

	return c.JSON(http.StatusOK, res)
}

// View is the endpoint to get an individual app
//
//	@Summary		View
//	@Description	Views an app
//	@Tags			apps
//	@Produce		json
//	@Param			id	path		string	true	"App ID"
//	@Success		200	{object}	apps.AppResponse
//	@Failure		404	{object}	errors.Error
//	@Failure		500	{object}	errors.Error
//	@Router			/api/v1/apps/{id} [get]
func (h *AppHandler) View(c *echo.Context) error {
	id := c.Param("id")

	app := apps.New(h.db, h.docker, h.jobs)
	res, err := app.GetApp(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, errors.Error{
			Identifier: errors.ErrAppNotFound,
			Message:    "App not found",
		})
	}

	return c.JSON(http.StatusOK, res)
}

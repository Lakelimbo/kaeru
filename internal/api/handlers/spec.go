package handlers

import (
	"fmt"
	"net/http"

	"github.com/Lakelimbo/kaeru/config"
	"github.com/Lakelimbo/kaeru/internal/logger"
	"github.com/Lakelimbo/kaeru/internal/server/spec"
	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/labstack/echo/v5"
)

type SpecHandler struct {
	config *config.ServerConfig
}

func NewSpecHandler(cfg *config.ServerConfig) *SpecHandler {
	return &SpecHandler{cfg}
}

func (h *SpecHandler) Routes(e *echo.Group) {
	e.GET("/swagger.json", h.File)
	e.GET("/scalar", h.Scalar)
}

func (h *SpecHandler) File(c *echo.Context) error {
	file, err := spec.Embed.ReadFile("swagger.json")
	if err != nil {
		logger.Debug(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.Blob(http.StatusOK, "application/json", file)
}

func (h *SpecHandler) Scalar(c *echo.Context) error {
	s, err := scalar.ApiReferenceHTML(&scalar.Options{
		SpecURL:       fmt.Sprintf("http://%s:%d/api/v1/swagger.json", h.config.Host, h.config.Port),
		BaseServerURL: "0.0.0.0",
		Theme:         scalar.ThemeDeepSpace,
		Layout:        scalar.LayoutClassic,
	})

	if err != nil {
		logger.Errorf("Failed to initialize Scalar: %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.HTML(http.StatusOK, s)
}

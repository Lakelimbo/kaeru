package server

import (
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

// CORS is a middleware to allow cross-origin HTTP requests
// (such as in a SPA), with specific methods, content-type, and
// headers.
//
// For the frontend on development, you should add Vite's localhost URL
// (or whatever you have set it), because you're likely using `pnpm dev`
// on a separate terminal.
func CORS(origins []string) echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: origins,
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAccessControlAllowOrigin,
			echo.HeaderAuthorization,
		},
		AllowCredentials: true,
	})
}

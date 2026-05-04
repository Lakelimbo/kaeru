package auth

import (
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
)

// Authentication (user) repository
type Auth struct {
	ctx *echo.Context
	db  *dbx.DB
}

func New(e *echo.Context, db *dbx.DB) *Auth {
	return &Auth{e, db}
}

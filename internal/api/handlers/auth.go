package handlers

import (
	"net/http"

	"github.com/Lakelimbo/kaeru/config"
	"github.com/Lakelimbo/kaeru/internal/auth"
	"github.com/Lakelimbo/kaeru/internal/logger"
	"github.com/Lakelimbo/kaeru/tools/errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
)

type AuthHandler struct {
	config *config.Config
	db     *dbx.DB
}

func NewAuthHandler(cfg *config.Config, db *dbx.DB) *AuthHandler {
	return &AuthHandler{cfg, db}
}

func (h *AuthHandler) Routes(e *echo.Group) {
	g := e.Group("/auth")

	g.POST("/login", h.Login)
	g.POST("/logout", h.Logout)
	g.POST("/refresh", h.Refresh)
	g.GET("/me", h.Me)
}

// Login is the endpoint to authenticate an existing user
//
//	@Summary		Login
//	@Description	Authenticates the user
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		auth.UserLoginRequest	true	"Login information"
//	@Success		200		{object}	auth.Token
//	@Failure		400		{object}	errors.Error
//	@Failure		409		{object}	errors.Error
//	@Failure		500		{object}	errors.Error
//	@Router			/api/v1/auth/login [post]
func (h *AuthHandler) Login(c *echo.Context) error {
	var req auth.UserLoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errors.Error{
			Identifier: errors.ErrInvalidBody,
			Message:    "Invalid request body",
		})
	}

	session := auth.New(c, h.db)
	token, err := session.AccessToken(req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, errors.Error{
			Identifier: errors.ErrUserInvalidCredentials,
			Message:    "Incorrect email or password",
		})
	}

	return c.JSON(http.StatusOK, token)
}

// Refresh is the endpoint to refresh the access token to keep the user logged in
//
//	@Summary		Refresh
//	@Description	Refreshes the authentication token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		auth.RefreshRequest	true	"Refresh token to be swapped"
//	@Success		200		{object}	auth.Token
//	@Failure		400		{object}	errors.Error
//	@Failure		401		{object}	errors.Error
//	@Failure		500		{object}	errors.Error
//	@Router			/api/v1/auth/refresh [post]
func (h *AuthHandler) Refresh(c *echo.Context) error {
	var req auth.RefreshRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errors.Error{
			Identifier: errors.ErrInvalidBody,
			Message:    "Invalid request body",
		})
	}

	session := auth.New(c, h.db)
	token, err := session.RefreshToken(req.RefreshToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, errors.Error{
			Identifier: errors.ErrUserInvalidCredentials,
			Message:    "Invalid or expired refresh token",
		})
	}

	return c.JSON(http.StatusOK, token)
}

// Logout removes the session (and tokens) from the current authenticated user
//
//	@Summary		Logout
//	@Description	Logs out the user
//	@Tags			auth
//	@Accept			json
//	@Param			request	body		auth.RefreshRequest	true	"Refresh token to be deleted"
//	@Success		200		{object}	nil
//	@Failure		400		{object}	errors.Error
//	@Failure		500		{object}	errors.Error
//	@Router			/api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *echo.Context) error {
	var req auth.RefreshRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errors.Error{
			Identifier: errors.ErrInvalidBody,
			Message:    "Invalid request body",
		})
	}

	session := auth.New(c, h.db)
	if err := session.DeleteToken(req.RefreshToken); err != nil {
		return c.JSON(http.StatusInternalServerError, errors.Error{
			Identifier: errors.ErrInternal,
			Message:    "An unknown error has happened",
		})
	}

	return c.NoContent(http.StatusOK)
}

// Me gets some basic user information. Used mostly for checking the currented
// authenticated user
//
//	@Summary		Me
//	@Description	Basic user information
//	@Tags			auth
//	@Param			Authorization	header	string	false	"Access token"
//	@Produces		json
//	@Success		200	{object}	auth.UserResponse
//	@Failure		400	{object}	errors.Error
//	@Failure		401	{object}	errors.Error
//	@Failure		404	{object}	errors.Error
//	@Failure		500	{object}	errors.Error
//	@Router			/api/v1/auth/me [get]
func (h *AuthHandler) Me(c *echo.Context) error {
	token, err := echo.ContextGet[*jwt.Token](c, "user")
	if err != nil {
		logger.Error(err)
		return c.JSON(http.StatusUnauthorized, errors.Error{
			Identifier: errors.ErrUserUnauthenticated,
			Message:    "Unable to get authentication token",
		})
	}

	claims, ok := token.Claims.(*auth.JWTClaims)
	if !ok {
		return c.JSON(http.StatusInternalServerError, errors.Error{
			Identifier: errors.ErrInternal,
			Message:    "An unknown error has happened",
		})
	}

	var user auth.UserResponse
	if err := h.db.Select("username", "email", "formatted_name", "avatar").
		From("users").
		Where(dbx.HashExp{"id": claims.ID}).
		One(&user); err != nil {
		return c.JSON(http.StatusNotFound, errors.Error{
			Identifier: errors.ErrUserNotFound,
			Message:    "User not found",
		})
	}

	return c.JSON(http.StatusOK, user)
}

package auth

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"os"
	"time"

	"github.com/Lakelimbo/kaeru/config"
	"github.com/Lakelimbo/kaeru/internal/database"
	"github.com/Lakelimbo/kaeru/internal/logger"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
)

type JWTClaims struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func generateAccessToken(user database.User) (string, error) {
	claims := JWTClaims{
		ID:    user.ID,
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "kaeru",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret()))
}

func generateRefreshToken() (string, time.Time, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", time.Time{}, err
	}

	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	return base64.URLEncoding.EncodeToString(b), expiresAt, err
}

func jwtSecret() string {
	env := os.Getenv(config.JWT_SECRET)
	if env == "" {
		// still unsure if the app should generate a secret if not
		// set.
		// on one hand, generating a token automatically is convenient, but
		// on the other, creating it manually is a safer strategy.

		// logger.Warnf("%s not found. Generating random 256-bit token...", config.JWT_SECRET)
		// b := make([]byte, 32)
		// _, _ = rand.Read(b)
		// os.Setenv(config.JWT_SECRET, hex.EncodeToString(b))
		// return hex.EncodeToString(b)
		logger.Errorf("%s not set. Set one before authenticating.", config.JWT_SECRET)
	}

	return env
}

func JWTMiddleware() echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(jwtSecret()),
		NewClaimsFunc: func(c *echo.Context) jwt.Claims {
			return new(JWTClaims)
		},
		ErrorHandler: func(c *echo.Context, err error) error {
			logger.Errorf("Unauthenticated: %v", err)
			return c.JSON(http.StatusUnauthorized, "You're not authorized to perform this task")
		},
		Skipper: func(c *echo.Context) bool {
			path := c.Request().URL.Path

			return path == "/api/v1/auth/login" ||
				path == "/api/v1/scalar" ||
				path == "/api/v1/swagger.json"
		},
	})
}

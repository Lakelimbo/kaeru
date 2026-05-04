package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/Lakelimbo/kaeru/internal/database"
	"github.com/Lakelimbo/kaeru/internal/logger"
	"github.com/Lakelimbo/kaeru/tools/security"
	"github.com/pocketbase/dbx"
	"golang.org/x/crypto/bcrypt"
)

type Token struct {
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	Expiration   time.Duration `json:"expiration"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *Auth) AccessToken(req UserLoginRequest) (*Token, error) {
	var user database.User
	if err := a.db.Select("id", "email", "password").
		From("users").
		Where(dbx.HashExp{"email": req.Email}).
		One(&user); err != nil {
		logger.Errorf("User by email <%s> not found", req.Email)
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		logger.Errorf("Failed to compare hash and password: %v", err)
		return nil, err
	}

	accessToken, err := generateAccessToken(user)
	if err != nil {
		logger.Errorf("Failed to generate access token for <%s>: %v", user.Email, err)
		return nil, err
	}

	refreshToken, _, err := generateRefreshToken()
	if err != nil {
		logger.Errorf("Failed to generate refresh token for <%s>: %v", user.Email, err)
		return nil, err
	}

	if err := a.NewAccess(user.ID, refreshToken); err != nil {
		logger.Errorf("Failed to add access token and hash for <%s>: %v", user.Email, err)
		return nil, err
	}

	return &Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Expiration:   15 * time.Minute,
	}, nil
}

func (a *Auth) RefreshToken(token string) (*Token, error) {
	hashed := sha256.Sum256([]byte(token))
	hashedStr := hex.EncodeToString(hashed[:])

	var access UserAccess
	if err := a.db.Select("*").
		From("user_access").
		Where(dbx.HashExp{"token_hash": hashedStr}).
		One(&access); err != nil {
		logger.Errorf("Failed to fetch hash: %v", err)
		return nil, err
	}

	var user database.User
	if err := a.db.Select("id", "email").
		From("users").
		Where(dbx.HashExp{"id": access.UserID}).
		One(&user); err != nil {
		logger.Errorf("Failed to fetch user with ID '%s' from hash token: %v", access.UserID, err)
		return nil, err
	}

	if err := a.DeleteToken(hashedStr); err != nil {
		logger.Errorf("Failed to delete old access for hash '%s': %v", hashedStr, err)
		return nil, err
	}

	newAccessToken, err := generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	newRefreshToken, exp, err := generateRefreshToken()
	if err != nil {
		return nil, err
	}

	if err := a.NewAccess(user.ID, newRefreshToken); err != nil {
		logger.Errorf("Failed to refresh token for <%s>: %v", user.Email, err)
		return nil, err
	}

	return &Token{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		Expiration:   time.Until(exp),
	}, nil
}

func (a *Auth) DeleteToken(token string) error {
	hashed := security.HashSHA256([]byte(token))

	if _, err := a.db.Delete(
		"user_access",
		dbx.HashExp{"token_hash": hashed},
	).Execute(); err != nil {
		logger.Errorf("Failed to delete token with hash '%s': %v", hashed, err)
		return err
	}

	return nil
}

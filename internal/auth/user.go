package auth

import (
	"time"

	"github.com/Lakelimbo/kaeru/internal/database"
	"github.com/Lakelimbo/kaeru/internal/logger"
	"github.com/Lakelimbo/kaeru/tools/security"
	"github.com/Lakelimbo/kaeru/tools/validations"
	"github.com/go-ozzo/ozzo-validation/is"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserRegisterRequest struct {
	Username        string `json:"username"`
	FormattedName   string `json:"formatted_name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type UserResponse struct {
	Username      string `json:"username"`
	FormattedName string `json:"formatted_name,omitempty"`
	Email         string `json:"email"`
	Avatar        string `json:"avatar,omitempty"`
}

type UserAccess struct {
	UserID    string    `json:"user_id" db:"pk,user_id"`
	TokenHash string    `json:"token_hash" db:"pk,token_hash"`
	IssuedAt  time.Time `json:"issued_at" db:"issued_at"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
}

func (a *Auth) NewUser(req UserRegisterRequest) error {
	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Username, validation.Required, validation.Length(3, 32)),
		validation.Field(&req.Email, validation.Required, is.Email),
		validation.Field(&req.Password, validation.Length(8, 0)),
		validation.Field(&req.ConfirmPassword, validation.By(validation.RuleFunc(validations.StringEq(req.Password, "passwords do not match")))),
	); err != nil {
		logger.Errorf("Invalid payload on user creation: %v", err)
		return err
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Errorf("Failed to encrypt password: %v", err)
		return err
	}

	user := database.User{
		ID:            uuid.NewString(),
		Username:      req.Username,
		Email:         req.Email,
		FormattedName: req.FormattedName,
		Password:      string(pass),
		Creation: database.Creation{
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		},
	}

	if err := a.db.Model(&user).Insert(); err != nil {
		logger.Errorf("Failed to create user: %v", err)
		return err
	}

	logger.Infof("User '%s' <%s> created", user.Username, user.Email)
	return nil
}

// NewAccess will add the hashed refresh token to the user_access table.
// This will take care of refresh tokens, as well as keeping track of
// user sessions.
func (a *Auth) NewAccess(id, token string) error {
	access := UserAccess{
		UserID:    id,
		TokenHash: security.HashSHA256([]byte(token)),
		IssuedAt:  time.Now().UTC(),
		ExpiresAt: time.Now().UTC().Add(time.Minute * 24 * 7),
	}

	if err := a.db.Model(&access).Insert(); err != nil {
		return err
	}

	return nil
}

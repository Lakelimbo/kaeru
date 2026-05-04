package auth_test

import (
	"crypto/sha256"
	"encoding/hex"
	"path/filepath"
	"testing"
	"time"

	"github.com/Lakelimbo/kaeru/config"
	"github.com/Lakelimbo/kaeru/internal/auth"
	"github.com/Lakelimbo/kaeru/tools/tests"
	"github.com/pocketbase/dbx"
)

func TestAccessToken(t *testing.T) {
	t.Parallel()

	cfg := config.Config{
		Database: config.DatabaseConfig{
			Path: filepath.Join(t.TempDir(), "test.db"),
		},
	}

	app := tests.NewTestApp(&cfg)
	a := auth.New(app.Server.AcquireContext(), app.DB)

	user := auth.UserRegisterRequest{
		Username:        "majin",
		Email:           "majin@boo.com",
		Password:        "abc1234567890",
		ConfirmPassword: "abc1234567890",
	}
	if err := a.NewUser(user); err != nil {
		t.Fatalf("expected to register user, got %q", err)
	}

	scenarios := []struct {
		name      string
		login     auth.UserLoginRequest
		expectErr bool
	}{
		{
			name: "successful_access",
			login: auth.UserLoginRequest{
				Email:    "majin@boo.com",
				Password: "abc1234567890",
			},
			expectErr: false,
		},
		{
			name: "wrong_access",
			login: auth.UserLoginRequest{
				Email:    "major@doo.com",
				Password: "shakalakalaka",
			},
			expectErr: true,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			token, err := a.AccessToken(s.login)
			if s.expectErr {
				if err == nil {
					t.Fatal("expected error but got nil")
				}

				return
			}
			if err != nil {
				t.Fatalf("did not expect error, but got %q", err)
			}

			exp := 15 * time.Minute
			if token.Expiration != exp {
				t.Fatalf("expected expiration of %d but got %d", exp, token.Expiration)
			}

			hashed := sha256.Sum256([]byte(token.RefreshToken))
			hashedStr := hex.EncodeToString(hashed[:])
			var access auth.UserAccess
			if err := app.DB.Select("*").
				From("user_access").
				Where(dbx.HashExp{"token_hash": hashedStr}).
				One(&access); err != nil {
				t.Fatal("expected to get result from query, got none")
			}

			// need to truncate time because DB insertion happens just a few miliseconds
			// before this last comparison
			tokenExp := time.Now().UTC().Truncate(time.Second).Add(time.Minute * 24 * 7)
			accessExp := access.ExpiresAt.Truncate(time.Second)
			if accessExp != tokenExp {
				t.Fatalf("expected expiration on %s, got %s", tokenExp, accessExp)
			}
		})
	}
}

func TestRefreshToken(t *testing.T) {
	t.Parallel()

	cfg := config.Config{
		Database: config.DatabaseConfig{
			Path: filepath.Join(t.TempDir(), "test.db"),
		},
	}

	app := tests.NewTestApp(&cfg)
	a := auth.New(app.Server.AcquireContext(), app.DB)

	user := auth.UserRegisterRequest{
		Username:        "majin",
		Email:           "majin@boo.com",
		Password:        "abc1234567890",
		ConfirmPassword: "abc1234567890",
	}
	if err := a.NewUser(user); err != nil {
		t.Fatalf("expected to register user, got %q", err)
	}

	login := auth.UserLoginRequest{
		Email:    "majin@boo.com",
		Password: "abc1234567890",
	}
	initialToken, err := a.AccessToken(login)
	if err != nil {
		t.Fatalf("expected to get token on login, got %q", err)
	}

	scenarios := []struct {
		name         string
		refreshToken string
		expectErr    bool
	}{
		{
			name:         "successful_refresh",
			refreshToken: initialToken.RefreshToken,
			expectErr:    false,
		},
		{
			name:         "wrong_refresh",
			refreshToken: "lololo-lololo-lolo",
			expectErr:    true,
		},
		{
			name:         "non_existent_refresh",
			refreshToken: "00000000-0000-0000-0000-000000000000",
			expectErr:    true,
		},
		{
			name:         "empty_token",
			refreshToken: "",
			expectErr:    true,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			newToken, err := a.RefreshToken(s.refreshToken)
			if s.expectErr {
				if err == nil {
					t.Fatal("expected error but got nil")
				}

				return
			}
			if err != nil {
				t.Fatalf("did not expect error, but got %q", err)
			}
			if newToken == nil {
				t.Fatal("expected token to not be nil")
			}

			newExp := newToken.Expiration.Truncate(time.Minute).Round(time.Hour)
			weekExp := time.Hour * 24 * 7
			if newExp != weekExp {
				t.Fatalf("expected expiration of %s, got %s", weekExp, newExp)
			}
			if newToken.AccessToken == "" {
				t.Fatal("expected access token, but got empty string")
			}
			if newToken.RefreshToken == "" {
				t.Fatal("expected refresh token, but got empty string")
			}
			if newToken.RefreshToken == s.refreshToken {
				t.Fatal("expected new refresh token to be different from the old one")
			}

			hashedOld := sha256.Sum256([]byte(s.refreshToken))
			hashedOldStr := hex.EncodeToString(hashedOld[:])
			var oldAccess auth.UserAccess
			if err := app.DB.Select("*").
				From("user_access").
				Where(dbx.HashExp{"token_hash": hashedOldStr}).
				One(&oldAccess); err != nil {
				t.Fatal("expected old token to be deleted from database")
			}

			hashedNew := sha256.Sum256([]byte(newToken.RefreshToken))
			hashedNewStr := hex.EncodeToString(hashedNew[:])
			var newAccess auth.UserAccess
			if err := app.DB.Select("*").
				From("user_access").
				Where(dbx.HashExp{"token_hash": hashedNewStr}).
				One(&newAccess); err != nil {
				t.Fatal("expected new token on the database, got none")
			}

			tokenExp := time.Now().UTC().Truncate(time.Second).Add(time.Minute * 24 * 7)
			refreshExp := newAccess.ExpiresAt.Truncate(time.Second)
			if refreshExp != tokenExp {
				t.Fatalf("expected expiration on %s, got %s", tokenExp, refreshExp)
			}
		})
	}
}

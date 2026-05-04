package auth_test

import (
	"path/filepath"
	"testing"

	"github.com/Lakelimbo/kaeru/config"
	"github.com/Lakelimbo/kaeru/internal/auth"
	"github.com/Lakelimbo/kaeru/internal/database"
	"github.com/Lakelimbo/kaeru/tools/tests"
	"github.com/google/uuid"
	"github.com/pocketbase/dbx"
)

func TestNewUser(t *testing.T) {
	t.Parallel()

	cfg := config.Config{
		Database: config.DatabaseConfig{
			Path: filepath.Join(t.TempDir(), "test.db"),
		},
	}

	app := tests.NewTestApp(&cfg)
	a := auth.New(app.Server.AcquireContext(), app.DB)

	scenarios := []struct {
		name      string
		user      auth.UserRegisterRequest
		pre       func(*testing.T, *auth.Auth)
		post      func(*testing.T, *auth.Auth, string)
		expectErr bool
	}{
		{
			name: "first_user",
			user: auth.UserRegisterRequest{
				Username:        "Majin",
				Email:           "majin@boo.com",
				Password:        "abc1234567890",
				ConfirmPassword: "abc1234567890",
			},
			post: func(t *testing.T, a *auth.Auth, s string) {
				var user database.User
				if err := app.DB.
					Select("*").
					From("users").
					Where(dbx.HashExp{"id": s}).
					One(&user); err != nil {
					t.Fatalf("expected user to be inserted into the database, got %q", err)
				}

				if user.Username != "Majin" {
					t.Fatalf("expected username 'Majin', got '%s'", user.Username)
				}
				if user.Email != "majin@boo.com" {
					t.Fatalf("expected email <majin@boo.com>, got <%s>", user.Email)
				}
				if user.Password == "abc1234567890" {
					t.Fatal("expected password to be hashed, got plaintext")
				}
				if user.CreatedAt.IsZero() {
					t.Error("expected created_at to be set")
				}
				if user.UpdatedAt.IsZero() {
					t.Error("expected updated_at to be set")
				}

				if _, err := uuid.Parse(user.ID); err != nil {
					t.Errorf("invalid UUID format ('%s')", user.ID)
				}
			},
			expectErr: false,
		},
		{
			name: "invalid_email",
			user: auth.UserRegisterRequest{
				Username:        "majin",
				Email:           "majin!",
				Password:        "abc1234567890",
				ConfirmPassword: "abc1234567890",
			},
			expectErr: true,
		},
		{
			name: "username_too_short",
			user: auth.UserRegisterRequest{
				Username:        "bu",
				Email:           "majin@boo.com",
				Password:        "abc1234567890",
				ConfirmPassword: "abc1234567890",
			},
			expectErr: true,
		},
		{
			name: "username_too_long",
			user: auth.UserRegisterRequest{
				Username:        "majinLoremIpsumDolorSitAmetIForgotTheRestThisShouldBeSoLongThatItDoesntEvenMakeSenseAnymore",
				Email:           "majin!",
				Password:        "abc1234567890",
				ConfirmPassword: "abc1234567890",
			},
			expectErr: true,
		},
		{
			name: "password_too_short",
			user: auth.UserRegisterRequest{
				Username:        "majin",
				Email:           "majin@boo.com",
				Password:        "123",
				ConfirmPassword: "123",
			},
			expectErr: true,
		},
		{
			name: "passwords_do_not_match",
			user: auth.UserRegisterRequest{
				Username:        "majin",
				Email:           "majin@boo.com",
				Password:        "abc1234567890",
				ConfirmPassword: "def0987654321",
			},
			expectErr: true,
		},
		{
			name: "empty_username",
			user: auth.UserRegisterRequest{
				Username:        "",
				Email:           "majin@boo.com",
				Password:        "abc1234567890",
				ConfirmPassword: "abc1234567890",
			},
			expectErr: true,
		},
		{
			name: "empty_email",
			user: auth.UserRegisterRequest{
				Username:        "majin",
				Email:           "",
				Password:        "abc1234567890",
				ConfirmPassword: "abc1234567890",
			},
			expectErr: true,
		},
		{
			name: "empty_password",
			user: auth.UserRegisterRequest{
				Username:        "majin",
				Email:           "majin@boo.com",
				Password:        "",
				ConfirmPassword: "",
			},
			expectErr: true,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			if s.pre != nil {
				s.pre(t, a)
			}

			err := a.NewUser(s.user)
			if s.expectErr {
				if err == nil {
					t.Fatal("expected error but got nil")
				}

				return
			}
			if err != nil {
				t.Fatalf("did not expect error, but got %q", err)
			}

			var user database.User
			if err := app.DB.Select("*").
				From("users").
				Where(dbx.HashExp{"email": s.user.Email}).
				One(&user); err != nil {
				t.Fatal("expected to find created user, got none")
			}

			if s.post != nil {
				s.post(t, a, user.ID)
			}
		})
	}
}

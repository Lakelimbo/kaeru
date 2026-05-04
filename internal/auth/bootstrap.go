package auth

import (
	"errors"
	"fmt"

	"charm.land/huh/v2"
	"github.com/Lakelimbo/kaeru/internal/database"
	"github.com/Lakelimbo/kaeru/internal/logger"
	"github.com/go-ozzo/ozzo-validation/is"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
)

// At the moment (May 2026), Kaeru does have basic support for
// multiple users, but we will focus on that later (so we don't have
// a /register endpoint yet either).
// Due to this, and to also prevent possible issues if you boot up the
// server in a non-HTTPS environment first, the user creation will be
// done inside the terminal.
//
// Eventually, when we have a /register endpoint, we could either ditch this,
// or have this run if it's the first run AND HTTPS is not available at setup
// moment.
func FirstRun(db *dbx.DB) {
	var user database.User
	if err := db.Select("id").From("users").Limit(1).One(&user); err != nil {
		session := New(echo.New().AcquireContext(), db)
		session.firstUserPrompt()
	}
}

func (a *Auth) firstUserPrompt() error {
	fmt.Print("🐸 Welcome to Kaeru! Let's create the your account.\n")

	var req UserRegisterRequest

	huh.ThemeCatppuccin(true)
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Username").
				Prompt(": ").
				CharLimit(32).
				Value(&req.Username).
				Validate(func(s string) error {
					if len(s) < 3 {
						return errors.New("username needs to have at least 3 characters")
					}

					return nil
				}),

			huh.NewInput().
				Title("Email").
				Prompt(": ").
				Value(&req.Email).
				Validate(func(s string) error {
					if err := validation.Validate(s, is.Email); err != nil {
						return errors.New("must be a valid email address")
					}

					return nil
				}),

			huh.NewInput().
				Title("Password").
				Prompt(": ").
				EchoMode(huh.EchoModePassword).
				Value(&req.Password).
				Validate(func(s string) error {
					if len(s) < 8 {
						return errors.New("password needs at least 8 characters")
					}

					return nil
				}),

			huh.NewInput().
				Title("Confirm password").
				Prompt(": ").
				EchoMode(huh.EchoModePassword).
				Value(&req.ConfirmPassword).
				Validate(func(s string) error {
					if s != req.Password {
						return errors.New("passwords do not match")
					}

					return nil
				}),
		),
	)

	if err := form.Run(); err != nil {
		logger.Fatalf("unable to create first user. Exiting...")
		return err
	}

	return a.NewUser(req)
}

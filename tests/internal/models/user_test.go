package models_test

import (
	"errors"
	"testing"

	"github.com/pedrotunin/go-jwt-auth/internal/models"
	"github.com/pedrotunin/go-jwt-auth/internal/utils"
)

func TestNewUser(t *testing.T) {
	t.Run("should return error when e-mail is invalid", func(t *testing.T) {
		email := "not an email"
		password := "test12345"

		_, err := models.NewUser(email, password)
		if !errors.Is(err, utils.ErrInvalidEmail) {
			t.Errorf("expected ErrInvalidEmail, got other: %s", err.Error())
		}

	})

	t.Run("should return error when password is invalid", func(t *testing.T) {
		email := "email@email.com"
		password := ""

		_, err := models.NewUser(email, password)
		if !errors.Is(err, utils.ErrPasswordTooShort) {
			t.Errorf("expected ErrPasswordTooShort, got other: %s", err.Error())
		}

	})

	t.Run("should return user", func(t *testing.T) {
		email := "email@email.com"
		password := "test12345"

		u, err := models.NewUser(email, password)
		if err != nil {
			t.Errorf("expected no error, got one")
		}

		if u.Email != email {
			t.Errorf("emails do not match")
		}

		if u.Password != password {
			t.Error("passwords do not match")
		}

	})
}

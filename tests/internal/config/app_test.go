package config_test

import (
	"database/sql"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pedrotunin/go-jwt-auth/internal/config"
)

func TestSetup(t *testing.T) {

	t.Run("should fail when JWT_TOKEN_SECRET is not set", func(t *testing.T) {
		app := &config.Application{}

		os.Unsetenv("JWT_TOKEN_SECRET")

		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic, got none")
			}
		}()

		app.Setup()
	})

	t.Run("should fail when JWT_REFRESH_TOKEN_SECRET is not set", func(t *testing.T) {
		app := &config.Application{}

		os.Unsetenv("JWT_REFRESH_TOKEN_SECRET")

		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic, got none")
			}
		}()

		app.Setup()
	})

	t.Run("should run without panicking", func(t *testing.T) {
		app := &config.Application{
			DB:     &sql.DB{},
			Router: gin.Default(),
		}

		os.Setenv("JWT_TOKEN_SECRET", "test")
		os.Setenv("JWT_REFRESH_TOKEN_SECRET", "test")

		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("expected no panics, got one: %+v", r)
			}
		}()

		app.Setup()
	})

}

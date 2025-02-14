package repositories

import "github.com/pedrotunin/go-jwt-auth/internal/models"

type AppRepository interface {
	CreateApp(*models.App) error
}

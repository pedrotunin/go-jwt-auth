package repositories

import "github.com/pedrotunin/go-jwt-auth/internal/models"

type AppRepository interface {
	GetAppByID(appID models.AppID) (*models.App, error)
	CreateApp(*models.App) error
	DeleteAppByID(appID models.AppID) error
}

package repositories

import "github.com/pedrotunin/go-jwt-auth/internal/models"

type AppRepository interface {
	GetAppsByUserID(userID models.UserID) ([]models.App, error)
	GetAppByID(appID models.AppID) (*models.App, error)
	CreateApp(*models.App) error
	UpdateApp(app *models.App) error
	DeleteAppByID(appID models.AppID) error
}

package services

import (
	"github.com/pedrotunin/go-jwt-auth/internal/models"
	"github.com/pedrotunin/go-jwt-auth/internal/repositories"
	"github.com/pedrotunin/go-jwt-auth/internal/utils"
	"github.com/pedrotunin/go-jwt-auth/internal/validators"
)

type IAppService interface {
	GetUserApps(userID models.UserID) ([]models.App, error)
	GetAppByID(appID models.AppID) (*models.App, error)
	CreateApp(app *models.App) error
	UpdateApp(app *models.App) error
	DeleteApp(appID models.AppID) error
}

type AppService struct {
	appRepository repositories.AppRepository
}

func NewAppService(repository repositories.AppRepository) IAppService {
	return &AppService{
		appRepository: repository,
	}
}

func (as *AppService) GetUserApps(userID models.UserID) ([]models.App, error) {
	apps, err := as.appRepository.GetAppsByUserID(userID)
	if err != nil {
		return nil, err
	}
	return apps, nil
}

func (as *AppService) GetAppByID(appID models.AppID) (*models.App, error) {
	app, err := as.appRepository.GetAppByID(appID)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (as *AppService) CreateApp(app *models.App) error {
	err := as.appRepository.CreateApp(app)
	if err != nil {
		return err
	}
	return nil
}

func (as *AppService) UpdateApp(app *models.App) error {
	err := validators.IsValidAppName(app.Name)
	if err != nil {
		return utils.ErrAppNameInvalid
	}

	err = validators.IsValidAppDescription(app.Description)
	if err != nil {
		return utils.ErrAppDescInvalid
	}

	err = as.appRepository.UpdateApp(app)
	if err != nil {
		return err
	}

	return nil
}

func (as *AppService) DeleteApp(appID models.AppID) error {
	err := as.appRepository.DeleteAppByID(appID)
	if err != nil {
		return err
	}
	return nil
}

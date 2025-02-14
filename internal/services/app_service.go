package services

import (
	"github.com/pedrotunin/go-jwt-auth/internal/models"
	"github.com/pedrotunin/go-jwt-auth/internal/repositories"
)

type IAppService interface {
	GetAppByID(appID models.AppID) (*models.App, error)
	CreateApp(app *models.App) error
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

func (as *AppService) DeleteApp(appID models.AppID) error {
	err := as.appRepository.DeleteAppByID(appID)
	if err != nil {
		return err
	}
	return nil
}

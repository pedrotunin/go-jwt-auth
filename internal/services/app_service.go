package services

import (
	"github.com/pedrotunin/go-jwt-auth/internal/models"
	"github.com/pedrotunin/go-jwt-auth/internal/repositories"
)

type IAppService interface {
	CreateApp(app *models.App) error
}

type AppService struct {
	appRepository repositories.AppRepository
}

func NewAppService(repository repositories.AppRepository) IAppService {
	return &AppService{
		appRepository: repository,
	}
}

func (as *AppService) CreateApp(app *models.App) error {
	err := as.appRepository.CreateApp(app)
	if err != nil {
		return err
	}
	return nil
}

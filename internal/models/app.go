package models

import (
	"time"

	"github.com/pedrotunin/go-jwt-auth/internal/validators"
)

type AppID = int
type AppName = string
type AppDescription = string

type App struct {
	ID          AppID
	Name        AppName
	Description AppDescription
	UserID      UserID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
}

func NewApp(name, description string, userID UserID) (*App, error) {

	err := validators.IsValidAppName(name)
	if err != nil {
		return nil, err
	}

	err = validators.IsValidAppDescription(description)
	if err != nil {
		return nil, err
	}

	return &App{
		Name:        name,
		Description: description,
		UserID:      userID,
	}, nil
}

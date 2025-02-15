package models

import (
	"time"

	"github.com/pedrotunin/go-jwt-auth/internal/validators"
)

type AppID = int
type AppName = string
type AppDescription = string

type App struct {
	ID          AppID          `json:"id"`
	Name        AppName        `json:"name"`
	Description AppDescription `json:"description"`
	UserID      UserID         `json:"user_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   time.Time      `json:"deleted_at,omitempty"`
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

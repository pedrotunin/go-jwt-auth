package repositories

import (
	"github.com/pedrotunin/go-jwt-auth/internal/models"
)

type UserRepository interface {
	GetUserByEmail(email models.UserEmail) (*models.User, error)
	CreateUser(u *models.User) (id int, err error)
}

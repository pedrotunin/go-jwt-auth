package repositories

import (
	"errors"

	"github.com/pedrotunin/jwt-auth/internal/models"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
	GetUserByID(id models.UserID) (*models.User, error)
	GetUserByEmail(email models.UserEmail) (*models.User, error)
	CreateUser(u *models.User) (id int, err error)
	UpdateUser(u *models.User) error
	DeleteUser(u *models.User) error
	DeleteUserByID(id models.UserID) error
}

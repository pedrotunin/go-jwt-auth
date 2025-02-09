package repositories

import (
	"errors"

	"github.com/pedrotunin/jwt-auth/internal/models"
)

var ErrUserNotFound = errors.New("user not found")
var ErrUserEmailAlreadyExists = errors.New("user's email already exists in our database")

type UserRepository interface {
	GetUserByEmail(email models.UserEmail) (*models.User, error)
	CreateUser(u *models.User) (id int, err error)
}

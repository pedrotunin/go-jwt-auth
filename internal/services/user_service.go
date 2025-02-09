package services

import (
	"github.com/pedrotunin/jwt-auth/internal/models"
	"github.com/pedrotunin/jwt-auth/internal/repositories"
)

type UserService struct {
	userRepository  repositories.UserRepository
	passwordService *PasswordService
}

func NewUserService(repository repositories.UserRepository, pwdService *PasswordService) *UserService {
	return &UserService{
		userRepository: repository,
	}
}

func (us *UserService) GetUserByEmail(email string) (*models.User, error) {
	user, err := us.userRepository.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserService) CreateUser(u *models.User) error {
	hash, err := us.passwordService.Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = hash

	id, err := us.userRepository.CreateUser(u)
	if err != nil {
		return err
	}

	u.ID = id

	return nil
}

package services

import (
	"log"

	"github.com/pedrotunin/go-jwt-auth/internal/models"
	"github.com/pedrotunin/go-jwt-auth/internal/repositories"
)

type UserService struct {
	userRepository repositories.UserRepository
	hashService    *HashService
}

func NewUserService(repository repositories.UserRepository, pwdService *HashService) *UserService {
	return &UserService{
		userRepository: repository,
	}
}

func (us *UserService) GetUserByEmail(email string) (*models.User, error) {
	user, err := us.userRepository.GetUserByEmail(email)
	if err != nil {
		log.Printf("GetUserByEmail: error getting user in database: %s", err.Error())
		return nil, err
	}

	log.Printf("GetUserByEmail: user found in database")
	return user, nil
}

func (us *UserService) CreateUser(u *models.User) error {
	hash, err := us.hashService.Hash(u.Password)
	if err != nil {
		log.Printf("CreateUser: error hashing password: %s", err.Error())
		return err
	}
	u.Password = hash

	id, err := us.userRepository.CreateUser(u)
	if err != nil {
		log.Printf("CreateUser: error creating user: %s", err.Error())
		return err
	}

	u.ID = id

	log.Printf("CreateUser: user created")
	return nil
}

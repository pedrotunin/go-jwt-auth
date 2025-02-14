package services

import (
	"log"

	"github.com/pedrotunin/go-jwt-auth/internal/models"
	"github.com/pedrotunin/go-jwt-auth/internal/repositories"
	"github.com/pedrotunin/go-jwt-auth/internal/utils"
)

type IUserService interface {
	GetUserByEmail(email string) (*models.User, error)
	CreateUser(u *models.User) error
	VerifyActiveUser(u *models.User) error
	ActivateUser(userID models.UserID) error
}

type UserService struct {
	userRepository repositories.UserRepository
	hashService    IHashService
}

func NewUserService(repository repositories.UserRepository, hashService IHashService) IUserService {
	return &UserService{
		userRepository: repository,
		hashService:    hashService,
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
	hash, err := us.hashService.HashArgon2id(u.Password)
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

func (us *UserService) VerifyActiveUser(u *models.User) error {
	if u.Status == utils.UserStatusActive {
		return nil
	}

	switch u.Status {
	case utils.UserStatusInactive:
		return utils.ErrUserInactive
	case utils.UserStatusPending:
		return utils.ErrUserPending
	default:
		return utils.ErrInvalidUserStatus
	}
}

func (us *UserService) ActivateUser(userID models.UserID) error {
	err := us.userRepository.ActivateUser(userID)
	if err != nil {
		return err
	}

	return nil
}

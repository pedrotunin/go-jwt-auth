package services

import (
	"github.com/alexedwards/argon2id"
	"github.com/pedrotunin/jwt-auth/internal/utils"
)

type PasswordService struct {
}

func NewPasswordService() *PasswordService {
	return &PasswordService{}
}

func (ps *PasswordService) Hash(password string) (hash string, err error) {
	hash, err = argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hash, err
}

func (ps *PasswordService) Compare(password string, hashedPassword string) error {
	match, err := argon2id.ComparePasswordAndHash(password, hashedPassword)
	if err != nil {
		return err
	}
	if !match {
		return utils.ErrPasswordsNotMatch
	}

	return nil
}

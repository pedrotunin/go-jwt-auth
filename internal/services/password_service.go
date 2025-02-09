package services

import (
	"errors"

	"github.com/alexedwards/argon2id"
)

type PasswordService struct {
}

func NewPasswordService() *PasswordService {
	return &PasswordService{}
}

var ErrPasswordsNotMatch = errors.New("passwords don't match")

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
		return ErrPasswordsNotMatch
	}

	return nil
}

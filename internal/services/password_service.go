package services

import (
	"log"

	"github.com/alexedwards/argon2id"
	"github.com/pedrotunin/go-jwt-auth/internal/utils"
)

type PasswordService struct {
}

func NewPasswordService() *PasswordService {
	return &PasswordService{}
}

func (ps *PasswordService) Hash(password string) (hash string, err error) {
	hash, err = argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		log.Printf("Hash: error creating password hash: %s", err.Error())
		return "", err
	}
	log.Printf("Hash: hash created")
	return hash, err
}

func (ps *PasswordService) Compare(password string, hashedPassword string) error {
	match, err := argon2id.ComparePasswordAndHash(password, hashedPassword)
	if err != nil {
		log.Printf("Compare: error comparing password and hash: %s", err.Error())
		return err
	}
	if !match {
		log.Printf("Compare: password do not match with hash")
		return utils.ErrPasswordsNotMatch
	}

	log.Printf("Compare: password and hash match")
	return nil
}

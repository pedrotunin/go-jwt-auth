package services

import (
	"log"

	"github.com/alexedwards/argon2id"
	"github.com/pedrotunin/go-jwt-auth/internal/utils"
)

type HashService struct {
}

func NewHashService() *HashService {
	return &HashService{}
}

func (ps *HashService) Hash(text string) (hash string, err error) {
	hash, err = argon2id.CreateHash(text, argon2id.DefaultParams)
	if err != nil {
		log.Printf("Hash: error creating hash: %s", err.Error())
		return "", err
	}
	log.Printf("Hash: hash created")
	return hash, err
}

func (ps *HashService) Compare(text string, hashedText string) error {
	match, err := argon2id.ComparePasswordAndHash(text, hashedText)
	if err != nil {
		log.Printf("Compare: error comparing text and hash: %s", err.Error())
		return err
	}
	if !match {
		log.Printf("Compare: text do not match with hash")
		return utils.ErrPasswordsNotMatch
	}

	log.Printf("Compare: text and hash match")
	return nil
}

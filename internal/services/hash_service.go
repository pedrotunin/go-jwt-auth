package services

import (
	"crypto/sha256"
	"encoding/hex"
	"log"

	"github.com/alexedwards/argon2id"
	"github.com/pedrotunin/go-jwt-auth/internal/utils"
)

type HashService struct {
}

func NewHashService() *HashService {
	return &HashService{}
}

func (hs *HashService) HashArgon2id(text string) (hash string, err error) {
	hash, err = argon2id.CreateHash(text, argon2id.DefaultParams)
	if err != nil {
		log.Printf("Hash: error creating hash: %s", err.Error())
		return "", err
	}
	log.Printf("Hash: hash created")
	return hash, err
}

func (hs *HashService) CompareArgon2id(text string, hashedText string) error {
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

func (hs *HashService) HashSHA256(text string) (hash string, err error) {
	h := sha256.New()
	h.Write([]byte(text))

	hashBytes := h.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)

	return hashString, nil
}

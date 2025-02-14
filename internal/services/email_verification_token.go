package services

import (
	"time"

	"github.com/pedrotunin/go-jwt-auth/internal/models"
	"github.com/pedrotunin/go-jwt-auth/internal/repositories"
	"github.com/pedrotunin/go-jwt-auth/internal/utils"
)

type IEmailVerificationTokenService interface {
	CreateToken(userID models.UserID) (string, error)
	IsValidToken(token models.EmailVerificationTokenContent, userID models.UserID) (*models.EmailVerificationToken, error)
	UseToken(token *models.EmailVerificationToken) error
}

type EmailVerificationTokenService struct {
	EmailVerificationTokenRepository repositories.EmailVerificationTokenRepository
}

func NewEmailVerificationTokenService(evtRepo repositories.EmailVerificationTokenRepository) *EmailVerificationTokenService {
	return &EmailVerificationTokenService{
		EmailVerificationTokenRepository: evtRepo,
	}
}

func (evts *EmailVerificationTokenService) CreateToken(userID models.UserID) (string, error) {
	token, err := utils.GetRandomString(16)
	if err != nil {
		return "", err
	}

	evToken := models.EmailVerificationToken{
		Content:   token,
		UserID:    userID,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(30 * time.Minute),
	}

	err = evts.EmailVerificationTokenRepository.CreateVerificationToken(&evToken)
	if err != nil {
		return "", err
	}

	return evToken.Content, nil
}

func (evts *EmailVerificationTokenService) IsValidToken(token models.EmailVerificationTokenContent, userID models.UserID) (*models.EmailVerificationToken, error) {
	evToken, err := evts.EmailVerificationTokenRepository.GetVerificationTokenByContent(token)
	if err != nil {
		return nil, err
	}

	if evToken.UserID != userID {
		return nil, utils.ErrUserIDsDoNotMatch
	}

	if evToken.ExpiresAt.After(time.Now()) {
		return nil, utils.ErrVerifyTokenExpired
	}

	return evToken, nil
}

func (evts *EmailVerificationTokenService) UseToken(token *models.EmailVerificationToken) error {
	err := evts.EmailVerificationTokenRepository.SetTokenToUsed(token)
	if err != nil {
		return err
	}

	return nil
}

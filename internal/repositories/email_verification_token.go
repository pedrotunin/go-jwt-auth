package repositories

import "github.com/pedrotunin/go-jwt-auth/internal/models"

type EmailVerificationTokenRepository interface {
	CreateVerificationToken(*models.EmailVerificationToken) error
	GetVerificationTokenByContent(content models.EmailVerificationTokenContent) (*models.EmailVerificationToken, error)
	SetTokenToUsed(*models.EmailVerificationToken) error
}

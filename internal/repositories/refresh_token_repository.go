package repositories

import (
	"github.com/pedrotunin/jwt-auth/internal/models"
)

type RefreshTokenRepository interface {
	CreateRefreshToken(token *models.RefreshToken) error
	GetRefreshTokenByContent(content models.RefreshTokenContent) (*models.RefreshToken, error)
	InvalidateRefreshTokenByContent(content models.RefreshTokenContent) error
}

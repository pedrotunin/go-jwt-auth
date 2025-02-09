package repositories

import (
	"errors"

	"github.com/pedrotunin/jwt-auth/internal/models"
)

var ErrRefreshTokenNotFound = errors.New("refresh token not found")

type RefreshTokenRepository interface {
	CreateRefreshToken(token *models.RefreshToken) error
	GetRefreshTokenByContent(content models.RefreshTokenContent) (*models.RefreshToken, error)
	InvalidateRefreshTokenByContent(content models.RefreshTokenContent) error
}

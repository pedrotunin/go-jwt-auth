package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pedrotunin/jwt-auth/internal/models"
	"github.com/pedrotunin/jwt-auth/internal/repositories"
)

var ErrTokenInvalid = errors.New("invalid token")

type JWTService struct {
	hmacSecret             string
	refreshTokenRepository repositories.RefreshTokenRepository
}

func NewJWTService(hmacSecret string, repo repositories.RefreshTokenRepository) *JWTService {
	return &JWTService{
		hmacSecret:             hmacSecret,
		refreshTokenRepository: repo,
	}
}

type TokenClaims struct {
	UserID models.UserID `json:"uid"`
	jwt.RegisteredClaims
}

func (js *JWTService) GenerateToken(userID models.UserID) (tokenString string, err error) {
	expiration := time.Now().Add(10 * time.Minute)

	claims := &TokenClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "jwt_auth",
			ExpiresAt: jwt.NewNumericDate(expiration),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err = token.SignedString([]byte(js.hmacSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

type RefreshTokenClaims struct {
	UserID models.UserID `json:"uid"`
	jwt.RegisteredClaims
}

func (js *JWTService) GenerateRefreshToken(userID models.UserID) (tokenString string, err error) {
	expiration := time.Now().Add(7 * 24 * time.Hour)

	claims := &RefreshTokenClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "jwt_auth",
			ExpiresAt: jwt.NewNumericDate(expiration),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err = token.SignedString([]byte(js.hmacSecret))
	if err != nil {
		return "", err
	}

	refreshToken := &models.RefreshToken{
		Content: tokenString,
		Status:  models.RefreshTokenStatusActive,
	}

	err = js.refreshTokenRepository.CreateRefreshToken(refreshToken)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (js *JWTService) ValidateToken(tokenString string) (*TokenClaims, error) {
	claims := &TokenClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
		}

		return []byte(js.hmacSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, ErrTokenInvalid
	}

	if claims, ok := token.Claims.(TokenClaims); ok {
		return &claims, nil
	} else {
		return nil, fmt.Errorf("error parsing token")
	}
}

func (js *JWTService) ValidateRefreshToken(tokenString string) (*RefreshTokenClaims, error) {
	claims := &RefreshTokenClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
		}

		return []byte(js.hmacSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, models.ErrRefreshTokenInvalid
	}

	refreshToken, err := js.refreshTokenRepository.GetRefreshTokenByContent(tokenString)
	if err != nil {
		return nil, err
	}

	if refreshToken.Status != models.RefreshTokenStatusActive {
		return nil, models.ErrRefreshTokenInvalid
	}

	return claims, nil
}

func (js *JWTService) InvalidateRefreshToken(tokenString string) error {
	err := js.refreshTokenRepository.InvalidateRefreshTokenByContent(tokenString)
	if err != nil {
		return err
	}

	return nil
}

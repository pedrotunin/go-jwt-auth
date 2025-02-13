package services

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pedrotunin/go-jwt-auth/internal/models"
	"github.com/pedrotunin/go-jwt-auth/internal/repositories"
	"github.com/pedrotunin/go-jwt-auth/internal/utils"
)

type JWTService struct {
	tokenSecret            string
	refreshTokenSecret     string
	refreshTokenRepository repositories.RefreshTokenRepository
	hashService            *HashService
}

func NewJWTService(tokenSecret, refreshTokenSecret string, repo repositories.RefreshTokenRepository, hashService *HashService) *JWTService {
	return &JWTService{
		tokenSecret:            tokenSecret,
		refreshTokenSecret:     refreshTokenSecret,
		refreshTokenRepository: repo,
		hashService:            hashService,
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

	tokenString, err = token.SignedString([]byte(js.tokenSecret))
	if err != nil {
		log.Printf("GenerateToken: error creating token: %s", err.Error())
		return "", err
	}

	log.Print("GenerateToken: token created")
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

	tokenString, err = token.SignedString([]byte(js.refreshTokenSecret))
	if err != nil {
		log.Printf("GenerateRefreshToken: error creating refresh token: %s", err.Error())
		return "", err
	}

	tokenHash, err := js.hashService.Hash(tokenString)
	if err != nil {
		log.Printf("GenerateRefreshToken: error hashing refresh token: %s", err.Error())
		return "", err
	}

	refreshToken := &models.RefreshToken{
		Content: tokenHash,
		Status:  models.RefreshTokenStatusActive,
	}

	err = js.refreshTokenRepository.CreateRefreshToken(refreshToken)
	if err != nil {
		log.Printf("GenerateRefreshToken: error creating refresh token in database: %s", err.Error())
		return "", err
	}

	log.Print("GenerateRefreshToken: refresh token created")
	return tokenString, nil
}

func (js *JWTService) ValidateToken(tokenString string) (*TokenClaims, error) {
	claims := TokenClaims{}

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
		}

		return []byte(js.tokenSecret), nil
	})
	if err != nil {
		log.Printf("ValidateToken: error parsing token: %s", err.Error())
		return nil, err
	}

	if !token.Valid {
		log.Print("ValidateToken: invalid token")
		return nil, utils.ErrTokenInvalid
	}

	log.Print("ValidateToken: token is valid")
	return &claims, nil
}

func (js *JWTService) ValidateRefreshToken(tokenString string) (*RefreshTokenClaims, error) {
	claims := RefreshTokenClaims{}

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
		}

		return []byte(js.refreshTokenSecret), nil
	})
	if err != nil {
		log.Printf("ValidateRefreshToken: error parsing token: %s", err.Error())
		return nil, err
	}

	if !token.Valid {
		log.Print("ValidateRefreshToken: invalid refresh token")
		return nil, utils.ErrRefreshTokenInvalid
	}

	tokenHash, err := js.hashService.Hash(tokenString)
	if err != nil {
		log.Printf("ValidateRefreshToken: error hashing token: %s", err.Error())
		return nil, err
	}

	refreshToken, err := js.refreshTokenRepository.GetRefreshTokenByContent(tokenHash)
	if err != nil {
		log.Printf("ValidateRefreshToken: error getting refresh token in database: %s", err.Error())
		return nil, err
	}

	if refreshToken.Status != models.RefreshTokenStatusActive {
		log.Print("ValidateRefreshToken: refresh token is invalid in the database")
		return nil, utils.ErrRefreshTokenInvalid
	}

	log.Print("ValidateRefreshToken: refresh token is valid")
	return &claims, nil
}

func (js *JWTService) InvalidateRefreshToken(tokenString string) error {
	err := js.refreshTokenRepository.InvalidateRefreshTokenByContent(tokenString)
	if err != nil {
		log.Printf("InvalidateRefreshToken: error invalidating refresh token in database: %s", err.Error())
		return err
	}

	log.Printf("InvalidateRefreshToken: refresh token invalidated")
	return nil
}

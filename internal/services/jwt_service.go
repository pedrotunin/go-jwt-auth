package services

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pedrotunin/jwt-auth/internal/models"
)

type JWTService struct {
	hmacSecret string
}

func NewJWTService(hmacSecret string) *JWTService {
	return &JWTService{
		hmacSecret: hmacSecret,
	}
}

type TokenClaims struct {
	UserID models.UserID
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
	UserID models.UserID
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

	return tokenString, nil
}

func (js *JWTService) ValidateToken(tokenString string) (jwt.Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
		}

		return []byte(js.hmacSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("error parsing token")
	}
}

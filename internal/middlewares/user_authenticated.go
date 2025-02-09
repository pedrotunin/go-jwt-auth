package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pedrotunin/jwt-auth/internal/services"
)

type AuthenticatedUserMiddleware struct {
	jwtService *services.JWTService
}

func NewAuthenticatedUserMiddleware(jwtService *services.JWTService) *AuthenticatedUserMiddleware {
	return &AuthenticatedUserMiddleware{
		jwtService: jwtService,
	}
}

func (aum *AuthenticatedUserMiddleware) IsAuthenticated(c *gin.Context) {
	authorization, ok := c.Request.Header["Authorization"]
	if !ok {
		c.String(http.StatusUnauthorized, "Authorization header not found")
		return
	}

	if len(authorization) > 1 {
		c.String(http.StatusUnauthorized, "multiple Authorization headers not accepted")
		return
	}

	if !strings.HasPrefix(authorization[0], "Bearer ") {
		c.String(http.StatusUnauthorized, "Authorization header malformed")
		return
	}

	tokenString := strings.Split(authorization[0], " ")[1]

	claims, err := aum.jwtService.ValidateToken(tokenString)
	if err != nil {
		c.String(http.StatusUnauthorized, fmt.Errorf("error validating bearer token: %w", err).Error())
		return
	}

	log.Printf("user authenticated: %v", claims)
}

package middlewares

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pedrotunin/go-jwt-auth/internal/services"
	"github.com/pedrotunin/go-jwt-auth/internal/utils"
)

type AuthenticatedUserMiddleware struct {
	jwtService *services.JWTService
}

func NewAuthenticatedUserMiddleware(jwtService *services.JWTService) *AuthenticatedUserMiddleware {
	return &AuthenticatedUserMiddleware{
		jwtService: jwtService,
	}
}

func (aum *AuthenticatedUserMiddleware) IsAuthenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization, ok := c.Request.Header["Authorization"]
		if !ok {
			log.Print("IsAuthenticated: Authorization header not found")
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.GetErrorResponse(utils.ErrAuthorizationHeaderNotFound))
			return
		}

		if len(authorization) > 1 {
			log.Print("IsAuthenticated: multiple Authorization headers not accepted")
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.GetErrorResponse(utils.ErrMultipleAuthorizationHeaders))
			return
		}

		if !strings.HasPrefix(authorization[0], "Bearer ") {
			log.Print("IsAuthenticated: Authorization header malformed")
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.GetErrorResponse(utils.ErrAuthorizationHeaderMalformed))
			return
		}

		tokenString := strings.Split(authorization[0], " ")[1]

		claims, err := aum.jwtService.ValidateToken(tokenString)
		if err != nil {
			log.Printf("IsAuthenticated: error validating token: %s", err.Error())
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.GetErrorResponse(utils.ErrTokenInvalid))
			return
		}

		log.Printf("user %d is authenticated", claims.UserID)
		c.Next()
	}
}

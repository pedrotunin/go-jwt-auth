package controllers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pedrotunin/go-jwt-auth/internal/models"
	"github.com/pedrotunin/go-jwt-auth/internal/services"
	"github.com/pedrotunin/go-jwt-auth/internal/utils"
)

type IAuthController interface {
	Login(c *gin.Context)
	Logout(c *gin.Context)
	Refresh(c *gin.Context)
}

type AuthController struct {
	UserService services.IUserService
	HashService services.IHashService
	JWTService  services.IJWTService
}

type loginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (ac *AuthController) Login(c *gin.Context) {
	var loginDTO loginDTO

	err := c.ShouldBindJSON(&loginDTO)
	if err != nil {
		log.Printf("Login: error during binding loginDTO: %s", err.Error())

		c.JSON(http.StatusBadRequest, utils.GetErrorResponse(
			fmt.Errorf("error parsing request body: %w", err),
		))
		return
	}

	_, err = models.NewUser(loginDTO.Email, loginDTO.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.GetErrorResponse(err))
		return
	}

	user, err := ac.UserService.GetUserByEmail(loginDTO.Email)
	if err != nil {
		log.Printf("Login: error getting user: %s", err.Error())

		if errors.Is(err, utils.ErrUserNotFound) {
			c.JSON(http.StatusUnprocessableEntity, utils.GetErrorResponse(utils.ErrEmailPasswordIncorrect))
			return
		}

		c.JSON(http.StatusInternalServerError, utils.GetErrorResponse(err))
		return
	}

	if err := ac.UserService.VerifyActiveUser(user); err != nil {

		if errors.Is(err, utils.ErrUserPending) {
			c.JSON(http.StatusBadRequest, utils.GetErrorResponse(fmt.Errorf("user is pending activation, check your e-mail to activate.")))
			return
		}

		c.JSON(http.StatusBadRequest, utils.GetErrorResponse(err))
		return
	}

	err = ac.HashService.CompareArgon2id(loginDTO.Password, user.Password)
	if err != nil {
		log.Printf("Login: error comparing password and hash: %s", err.Error())

		c.JSON(http.StatusUnprocessableEntity, utils.GetErrorResponse(utils.ErrEmailPasswordIncorrect))
		return

	}

	accessToken, err := ac.JWTService.GenerateToken(user.ID)
	if err != nil {
		log.Printf("Login: error generating token: %s", err.Error())
		c.JSON(http.StatusInternalServerError, utils.GetErrorResponse(utils.ErrInternalServerError))
		return
	}

	refreshToken, err := ac.JWTService.GenerateRefreshToken(user.ID)
	if err != nil {
		log.Printf("Login: error generating refresh token: %s", err.Error())
		c.JSON(http.StatusInternalServerError, utils.GetErrorResponse(utils.ErrInternalServerError))
		return
	}

	log.Printf("Login: login successful")
	c.JSON(http.StatusOK, map[string]string{
		"messagge":      "login successful",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (ac *AuthController) Logout(c *gin.Context) {
	authorization := c.Request.Header["Authorization"]
	token := strings.Split(authorization[0], " ")[1]

	claims, err := ac.JWTService.ValidateToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.GetErrorResponse(utils.ErrTokenInvalid))
		return
	}

	userID := claims.UserID

	err = ac.JWTService.InvalidateRefreshTokensByUserID(userID)
	if err != nil {
		log.Printf("Logout: error invalidating refresh tokens: %s", err.Error())
		c.JSON(http.StatusInternalServerError, utils.GetErrorResponse(utils.ErrInternalServerError))
		return
	}

	log.Print("Logout: logout successful")
	c.String(http.StatusOK, "")
}

type refreshDTO struct {
	RefreshToken string `json:"refresh_token"`
}

func (ac *AuthController) Refresh(c *gin.Context) {
	var refreshDTO refreshDTO

	err := c.ShouldBindJSON(&refreshDTO)
	if err != nil {
		log.Printf("Refresh: error during binding refreshDTO: %s", err.Error())
		c.JSON(http.StatusBadRequest, utils.GetErrorResponse(
			fmt.Errorf("error parsing request body: %w", err),
		))
		return
	}

	claims, err := ac.JWTService.ValidateRefreshToken(refreshDTO.RefreshToken)
	if err != nil {
		log.Printf("Refresh: error validating refresh token: %s", err.Error())

		if errors.Is(err, utils.ErrRefreshTokenInvalid) {
			c.JSON(http.StatusBadRequest, utils.GetErrorResponse(utils.ErrRefreshTokenInvalid))
			return
		}

		c.JSON(http.StatusInternalServerError, utils.GetErrorResponse(utils.ErrInternalServerError))
		return
	}

	accessToken, err := ac.JWTService.GenerateToken(claims.UserID)
	if err != nil {
		log.Printf("Refresh: error generating token: %s", err.Error())
		c.JSON(http.StatusInternalServerError, utils.GetErrorResponse(utils.ErrInternalServerError))
		return
	}

	refreshToken, err := ac.JWTService.GenerateRefreshToken(claims.UserID)
	if err != nil {
		log.Printf("Refresh: error generating refresh token: %s", err.Error())
		c.JSON(http.StatusInternalServerError, utils.GetErrorResponse(utils.ErrInternalServerError))
		return
	}

	err = ac.JWTService.InvalidateRefreshToken(refreshDTO.RefreshToken)
	if err != nil {
		log.Printf("Refresh: error invalidating refresh token: %s", err.Error())
		c.JSON(http.StatusInternalServerError, utils.GetErrorResponse(utils.ErrInternalServerError))
		return
	}

	log.Printf("Refresh: successfully refreshed tokens")
	c.JSON(http.StatusOK, map[string]string{
		"messagge":      "tokens refreshed",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})

}

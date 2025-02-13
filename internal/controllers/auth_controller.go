package controllers

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pedrotunin/go-jwt-auth/internal/models"
	"github.com/pedrotunin/go-jwt-auth/internal/services"
	"github.com/pedrotunin/go-jwt-auth/internal/utils"
)

type AuthController struct {
	UserService *services.UserService
	HashService *services.HashService
	JWTService  *services.JWTService
}

type loginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (ac *AuthController) Login(c *gin.Context) {
	var loginDTO loginDTO

	err := c.ShouldBind(&loginDTO)
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

type refreshDTO struct {
	RefreshToken string `json:"refresh_token"`
}

func (ac *AuthController) Refresh(c *gin.Context) {
	var refreshDTO refreshDTO

	err := c.ShouldBind(&refreshDTO)
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

package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pedrotunin/jwt-auth/internal/services"
	"github.com/pedrotunin/jwt-auth/internal/utils"
)

type AuthController struct {
	UserService     *services.UserService
	PasswordService *services.PasswordService
	JWTService      *services.JWTService
}

type loginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (ac *AuthController) Login(c *gin.Context) {
	var loginDTO loginDTO

	err := c.ShouldBind(&loginDTO)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	// TODO: implement input validation

	user, err := ac.UserService.GetUserByEmail(loginDTO.Email)
	if err != nil {
		if errors.Is(err, utils.ErrUserNotFound) {
			c.String(http.StatusNotFound, "email or password incorrect")
			return
		}

		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	err = ac.PasswordService.Compare(loginDTO.Password, user.Password)
	if err != nil {
		c.String(http.StatusBadRequest, "email or password incorrect")
		return

	}

	accessToken, err := ac.JWTService.GenerateToken(user.ID)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	refreshToken, err := ac.JWTService.GenerateRefreshToken(user.ID)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

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
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	// TODO: implement input validation

	claims, err := ac.JWTService.ValidateRefreshToken(refreshDTO.RefreshToken)
	if err != nil {
		if errors.Is(err, utils.ErrRefreshTokenInvalid) {
			c.String(http.StatusBadRequest, "invalid refresh token")
			return
		}

		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	accessToken, err := ac.JWTService.GenerateToken(claims.UserID)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	refreshToken, err := ac.JWTService.GenerateRefreshToken(claims.UserID)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	err = ac.JWTService.InvalidateRefreshToken(refreshDTO.RefreshToken)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]string{
		"messagge":      "tokens refreshed",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})

}

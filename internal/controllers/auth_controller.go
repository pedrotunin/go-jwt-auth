package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pedrotunin/jwt-auth/internal/repositories"
	"github.com/pedrotunin/jwt-auth/internal/services"
)

type AuthController struct {
	userService     *services.UserService
	passwordService *services.PasswordService
}

func NewAuthController(userService *services.UserService, pwdService *services.PasswordService) *AuthController {
	return &AuthController{
		userService:     userService,
		passwordService: pwdService,
	}
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

	user, err := ac.userService.GetUserByEmail(loginDTO.Email)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			c.String(http.StatusNotFound, "email or password incorrect")
			return
		}

		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	err = ac.passwordService.Compare(loginDTO.Password, user.Password)
	if err != nil {
		c.String(http.StatusBadRequest, "email or password incorrect")
		return

	}
	c.String(http.StatusOK, "login ok")
}

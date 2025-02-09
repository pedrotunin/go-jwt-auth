package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pedrotunin/jwt-auth/internal/models"
	"github.com/pedrotunin/jwt-auth/internal/services"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (ac *UserController) CreateUser(c *gin.Context) {
	var u models.User

	err := c.ShouldBind(&u)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	// TODO: implement request validation

	err = ac.userService.CreateUser(&u)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.String(http.StatusCreated, fmt.Sprintf("id: %d", u.ID))
}

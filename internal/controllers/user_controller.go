package controllers

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pedrotunin/jwt-auth/internal/models"
	"github.com/pedrotunin/jwt-auth/internal/services"
	"github.com/pedrotunin/jwt-auth/internal/utils"
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
		log.Printf("CreateUser: error during binding user: %s", err.Error())

		c.JSON(http.StatusBadRequest, utils.GetErrorResponse(
			fmt.Errorf("error parsing request body: %w", err),
		))
		return
	}

	// TODO: implement request validation

	err = ac.userService.CreateUser(&u)
	if err != nil {
		log.Printf("CreateUser: error creating user: %s", err.Error())

		if errors.Is(err, utils.ErrUserEmailAlreadyExists) {
			c.JSON(http.StatusUnprocessableEntity, utils.GetErrorResponse(err))
			return
		}

		c.JSON(http.StatusInternalServerError, utils.GetErrorResponse(utils.ErrInternalServerError))
		return
	}

	log.Printf("CreateUser: user created with id %d", u.ID)
	c.String(http.StatusCreated, "")
}

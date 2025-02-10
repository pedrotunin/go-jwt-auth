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

type createUserDTO struct {
	Email    string
	Password string
}

func (ac *UserController) CreateUser(c *gin.Context) {
	var createUserDTO createUserDTO

	err := c.ShouldBind(&createUserDTO)
	if err != nil {
		log.Printf("CreateUser: error during binding user: %s", err.Error())

		c.JSON(http.StatusBadRequest, utils.GetErrorResponse(
			fmt.Errorf("error parsing request body: %w", err),
		))
		return
	}

	user, err := models.NewUser(createUserDTO.Email, createUserDTO.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.GetErrorResponse(err))
		return
	}

	err = ac.userService.CreateUser(user)
	if err != nil {
		log.Printf("CreateUser: error creating user: %s", err.Error())

		if errors.Is(err, utils.ErrUserEmailAlreadyExists) {
			c.JSON(http.StatusUnprocessableEntity, utils.GetErrorResponse(err))
			return
		}

		c.JSON(http.StatusInternalServerError, utils.GetErrorResponse(utils.ErrInternalServerError))
		return
	}

	log.Printf("CreateUser: user created with id %d", user.ID)
	c.String(http.StatusCreated, "")
}

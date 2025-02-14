package controllers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pedrotunin/go-jwt-auth/internal/models"
	"github.com/pedrotunin/go-jwt-auth/internal/services"
	"github.com/pedrotunin/go-jwt-auth/internal/utils"
)

type IUserController interface {
	CreateUser(c *gin.Context)
	VerifyUser(c *gin.Context)
}

type UserController struct {
	userService                   services.IUserService
	emailVerificationTokenService services.IEmailVerificationTokenService
}

func NewUserController(userService services.IUserService, evtService services.IEmailVerificationTokenService) IUserController {
	return &UserController{
		userService:                   userService,
		emailVerificationTokenService: evtService,
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

	evToken, err := ac.emailVerificationTokenService.CreateToken(user.ID)
	if err != nil {
		log.Printf("CreateUser: error creating email verify token: %s", err.Error())
		c.JSON(http.StatusInternalServerError, utils.GetErrorResponse(utils.ErrInternalServerError))
		return
	}

	// TODO: send confirmation email
	log.Printf("sending email with token %s", evToken)

	c.JSON(http.StatusCreated, map[string]string{
		"message": "user created, check e-mail for activation instructions.",
	})
}

func (ac *UserController) VerifyUser(c *gin.Context) {
	id := c.Param("id")
	queryToken := c.Query("token")

	userID, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("VerifyUser: error converting userID to int: %s", err.Error())
		c.JSON(http.StatusBadRequest, utils.GetErrorResponse(utils.ErrInvalidUserID))
		return
	}

	if queryToken == "" {
		log.Print("VerifyUser: verify token not found in URL")
		c.JSON(http.StatusBadRequest, utils.GetErrorResponse(utils.ErrVerifyTokenNotFound))
		return
	}

	evToken, err := ac.emailVerificationTokenService.IsValidToken(queryToken, userID)
	if err != nil {
		log.Printf("VerifyUser: error validing verify token: %s", err.Error())

		if errors.Is(err, utils.ErrVerifyTokenNotFound) {
			c.JSON(http.StatusBadRequest, utils.GetErrorResponse(fmt.Errorf("valid verify token not found")))
			return
		}

		if errors.Is(err, utils.ErrVerifyTokenExpired) {
			c.JSON(http.StatusBadRequest, utils.GetErrorResponse(utils.ErrVerifyTokenExpired))
			return
		}

		if errors.Is(err, utils.ErrUserIDsDoNotMatch) {
			c.JSON(http.StatusForbidden, utils.GetErrorResponse(utils.ErrUserIDsDoNotMatch))
			return
		}

		c.JSON(http.StatusInternalServerError, utils.GetErrorResponse(utils.ErrInternalServerError))
		return
	}

	err = ac.userService.ActivateUser(evToken.UserID)
	if err != nil {
		log.Printf("VerifyUser: error activating user: %s", err.Error())
		c.JSON(http.StatusInternalServerError, utils.GetErrorResponse(utils.ErrInternalServerError))
		return
	}

	err = ac.emailVerificationTokenService.UseToken(evToken)
	if err != nil {
		log.Printf("VerifyUser: error updating token status: %s", err.Error())
		c.JSON(http.StatusInternalServerError, utils.GetErrorResponse(utils.ErrInternalServerError))
		return
	}

	c.JSON(http.StatusOK, map[string]string{
		"message": "user activated",
	})
}

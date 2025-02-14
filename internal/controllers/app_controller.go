package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pedrotunin/go-jwt-auth/internal/models"
	"github.com/pedrotunin/go-jwt-auth/internal/services"
	"github.com/pedrotunin/go-jwt-auth/internal/utils"
)

type IAppController interface {
	Create(c *gin.Context)
}

type AppController struct {
	AppService services.IAppService
}

type createAppDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (ac *AppController) Create(c *gin.Context) {
	var createDTO createAppDTO

	err := c.ShouldBindJSON(&createDTO)
	if err != nil {
		log.Printf("Create: error during binding createAppDTO: %s", err.Error())

		c.JSON(http.StatusBadRequest, utils.GetErrorResponse(
			fmt.Errorf("error parsing request body: %w", err),
		))
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		log.Print("Create: userID value do not exists in context")
		c.JSON(http.StatusInternalServerError, utils.GetErrorResponse(utils.ErrInternalServerError))
		return
	}

	app, err := models.NewApp(createDTO.Name, createDTO.Description, userID.(models.UserID))
	if err != nil {
		log.Printf("Create: error validating app: %s", err.Error())
		c.JSON(http.StatusBadRequest, utils.GetErrorResponse(err))
		return
	}

	err = ac.AppService.CreateApp(app)
	if err != nil {
		log.Printf("Create: error creating app: %s", err.Error())
		c.JSON(http.StatusInternalServerError, utils.GetErrorResponse(utils.ErrInternalServerError))
		return
	}

	c.JSON(http.StatusCreated, map[string]string{
		"message": "app created",
	})
}

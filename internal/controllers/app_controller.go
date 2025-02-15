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

type IAppController interface {
	GetAll(c *gin.Context)
	GetOne(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	DeleteByID(c *gin.Context)
}

type AppController struct {
	AppService services.IAppService
}

func (ac *AppController) GetAll(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		log.Print("GetAll: userID value do not exists in context")
		c.JSON(http.StatusInternalServerError, utils.GetErrorResponse(utils.ErrInternalServerError))
		return
	}

	apps, err := ac.AppService.GetUserApps(userID.(models.UserID))
	if err != nil {

		log.Printf("GetAll: error getting apps: %s", err.Error())

		if errors.Is(err, utils.ErrAppNotFound) {
			c.JSON(http.StatusNotFound, utils.GetErrorResponse(utils.ErrAppNotFound))
			return
		}

		c.JSON(http.StatusInternalServerError, utils.GetErrorResponse(utils.ErrInternalServerError))
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"data": map[string]any{
			"apps": apps,
		},
	})
}

func (ac *AppController) GetOne(c *gin.Context) {
	appID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("GetOne: invalid appID: %s", err.Error())
		c.JSON(http.StatusBadRequest, utils.GetErrorResponse(utils.ErrAppIDInvalid))
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		log.Print("GetOne: userID value do not exists in context")
		c.JSON(http.StatusInternalServerError, utils.GetErrorResponse(utils.ErrInternalServerError))
		return
	}

	app, err := ac.AppService.GetAppByID(appID)
	if err != nil {

		log.Printf("GetOne: error getting app: %s", err.Error())

		if errors.Is(err, utils.ErrAppNotFound) {
			c.JSON(http.StatusNotFound, utils.GetErrorResponse(utils.ErrAppNotFound))
			return
		}

		c.JSON(http.StatusInternalServerError, utils.GetErrorResponse(utils.ErrInternalServerError))
		return
	}

	if app.UserID != userID {
		log.Print("GetOne: userIDs do not match")
		c.JSON(http.StatusForbidden, utils.GetErrorResponse(utils.ErrUserIDsDoNotMatch))
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"data": map[string]any{
			"app": app,
		},
	})
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

type updateAppDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (ac *AppController) Update(c *gin.Context) {
	var updateDTO updateAppDTO

	err := c.ShouldBindJSON(&updateDTO)
	if err != nil {
		log.Printf("Update: error during binding createAppDTO: %s", err.Error())

		c.JSON(http.StatusBadRequest, utils.GetErrorResponse(
			fmt.Errorf("error parsing request body: %w", err),
		))
		return
	}

	appID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("Update: invalid appID: %s", err.Error())
		c.JSON(http.StatusBadRequest, utils.GetErrorResponse(utils.ErrAppIDInvalid))
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		log.Print("Update: userID value do not exists in context")
		c.JSON(http.StatusInternalServerError, utils.GetErrorResponse(utils.ErrInternalServerError))
		return
	}

	app, err := ac.AppService.GetAppByID(appID)
	if err != nil {

		log.Printf("Update: error getting app: %s", err.Error())

		if errors.Is(err, utils.ErrAppNotFound) {
			c.JSON(http.StatusNotFound, utils.GetErrorResponse(utils.ErrAppNotFound))
			return
		}

		c.JSON(http.StatusInternalServerError, utils.GetErrorResponse(utils.ErrInternalServerError))
		return
	}

	if app.UserID != userID {
		log.Print("Update: userIDs do not match")
		c.JSON(http.StatusForbidden, utils.GetErrorResponse(utils.ErrUserIDsDoNotMatch))
		return
	}

	app.Name = updateDTO.Name
	app.Description = updateDTO.Description

	err = ac.AppService.UpdateApp(app)
	if err != nil {
		log.Printf("Update: errors updating app: %s", err.Error())

		if errors.Is(err, utils.ErrAppNameInvalid) || errors.Is(err, utils.ErrAppDescInvalid) {
			c.JSON(http.StatusBadRequest, utils.GetErrorResponse(err))
			return
		}

		c.JSON(http.StatusInternalServerError, utils.GetErrorResponse(utils.ErrInternalServerError))
		return
	}

	c.String(http.StatusOK, "")
}

func (ac *AppController) DeleteByID(c *gin.Context) {
	appID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("DeleteByID: invalid appID: %s", err.Error())
		c.JSON(http.StatusBadRequest, utils.GetErrorResponse(utils.ErrAppIDInvalid))
		return
	}

	app, err := ac.AppService.GetAppByID(appID)
	if err != nil {
		log.Printf("DeleteByID: error getting app: %s", err.Error())

		if errors.Is(err, utils.ErrAppNotFound) {
			c.JSON(http.StatusNotFound, utils.GetErrorResponse(utils.ErrAppNotFound))
			return
		}

		c.JSON(http.StatusInternalServerError, utils.GetErrorResponse(utils.ErrInternalServerError))
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		log.Print("DeleteByID: userID value do not exists in context")
		c.JSON(http.StatusInternalServerError, utils.GetErrorResponse(utils.ErrInternalServerError))
		return
	}

	if app.UserID != userID {
		log.Print("DeleteByID: userIDs do not match")
		c.JSON(http.StatusForbidden, utils.GetErrorResponse(utils.ErrUserIDsDoNotMatch))
		return
	}

	err = ac.AppService.DeleteApp(app.ID)
	if err != nil {
		log.Printf("DeleteByID: error deleting app: %s", err.Error())
		c.JSON(http.StatusInternalServerError, utils.GetErrorResponse(utils.ErrInternalServerError))
		return
	}

	c.String(http.StatusOK, "")
}

package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pedrotunin/jwt-auth/internal/controllers"
	"github.com/pedrotunin/jwt-auth/internal/middlewares"
)

type Routes struct {
	Router      gin.IRouter
	Middlewares *middlewares.Middlewares
	Controllers *controllers.Controllers
}

func (r *Routes) Setup() {
	r.Router.POST("/users", r.Controllers.UserController.CreateUser)
	r.Router.POST("/login", r.Controllers.AuthController.Login)
}

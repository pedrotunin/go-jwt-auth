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
	r.Router.Use(r.Middlewares.LoggerMiddleware.LogRequest())

	usersRoutes := r.Router.Group("/users")
	{
		usersRoutes.POST("", r.Controllers.UserController.CreateUser)
	}

	authRoutes := r.Router.Group("/auth")
	{
		authRoutes.POST("/login", r.Controllers.AuthController.Login)
		authRoutes.POST("/refresh", r.Controllers.AuthController.Refresh)
	}
}

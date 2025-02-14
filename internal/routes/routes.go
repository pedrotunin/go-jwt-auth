package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pedrotunin/go-jwt-auth/internal/controllers"
	"github.com/pedrotunin/go-jwt-auth/internal/middlewares"
)

type Routes struct {
	Router      gin.IRouter
	Middlewares *middlewares.Middlewares
	Controllers *controllers.Controllers
}

func (r *Routes) Setup() {
	r.Router.Use(r.Middlewares.LoggerMiddleware.LogRequest())

	v1 := r.Router.Group("/v1")
	{
		users := v1.Group("/users")
		{
			users.POST("/", r.Controllers.UserController.CreateUser)
			users.GET("/:id/verify", r.Controllers.UserController.VerifyUser)
		}

		auth := v1.Group("/auth")
		{
			auth.POST("/login", r.Controllers.AuthController.Login)
			auth.POST("/logout", r.Middlewares.AuthenticatedUserMiddleware.IsAuthenticated(), r.Controllers.AuthController.Logout)
			auth.POST("/refresh", r.Controllers.AuthController.Refresh)
		}

		apps := v1.Group("/apps", r.Middlewares.AuthenticatedUserMiddleware.IsAuthenticated())
		{
			apps.POST("/", r.Controllers.AppController.Create)
			apps.DELETE("/:id", r.Controllers.AppController.DeleteByID)
		}

	}

}

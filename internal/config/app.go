package config

import (
	"database/sql"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/pedrotunin/jwt-auth/internal/controllers"
	"github.com/pedrotunin/jwt-auth/internal/middlewares"
	"github.com/pedrotunin/jwt-auth/internal/repositories"
	"github.com/pedrotunin/jwt-auth/internal/routes"
	"github.com/pedrotunin/jwt-auth/internal/services"
)

type Application struct {
	Router gin.IRouter
	DB     *sql.DB
}

func (app *Application) Setup() {
	log.Print("starting app setup")

	hmacSecret := os.Getenv("HMAC_SECRET")
	if hmacSecret == "" {
		log.Fatal("HMAC_SECRET env var not found")
	}

	// Setup repositories
	userRepository := repositories.NewPSQLUserRepository(app.DB)

	// Setup services
	jwtService := services.NewJWTService(hmacSecret)
	passwordService := services.NewPasswordService()
	userService := services.NewUserService(userRepository, passwordService)

	// Setup controllers
	authController := controllers.NewAuthController(userService, passwordService)
	userController := controllers.NewUserController(userService)

	// Setup middlewares
	authenticatedUserMiddleware := middlewares.NewAuthenticatedUserMiddleware(jwtService)

	// Setup Routes
	routes := &routes.Routes{
		Router: app.Router,
		Middlewares: &middlewares.Middlewares{
			AuthenticatedUserMiddleware: authenticatedUserMiddleware,
		},
		Controllers: &controllers.Controllers{
			AuthController: authController,
			UserController: userController,
		},
	}
	routes.Setup()

	log.Print("finished app setup")
}

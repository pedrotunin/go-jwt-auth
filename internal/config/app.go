package config

import (
	"database/sql"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/pedrotunin/go-jwt-auth/internal/controllers"
	"github.com/pedrotunin/go-jwt-auth/internal/middlewares"
	"github.com/pedrotunin/go-jwt-auth/internal/repositories"
	"github.com/pedrotunin/go-jwt-auth/internal/routes"
	"github.com/pedrotunin/go-jwt-auth/internal/services"
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
	refreshTokenRepository := repositories.NewPSQLRefreshTokenRepository(app.DB)

	// Setup services
	jwtService := services.NewJWTService(hmacSecret, refreshTokenRepository)
	passwordService := services.NewPasswordService()
	userService := services.NewUserService(userRepository, passwordService)

	// Setup controllers
	authController := &controllers.AuthController{
		UserService:     userService,
		PasswordService: passwordService,
		JWTService:      jwtService,
	}
	userController := controllers.NewUserController(userService)

	// Setup middlewares
	authenticatedUserMiddleware := middlewares.NewAuthenticatedUserMiddleware(jwtService)
	loggerMiddleware := middlewares.NewLoggerMiddleware()

	// Setup Routes
	routes := &routes.Routes{
		Router: app.Router,
		Middlewares: &middlewares.Middlewares{
			AuthenticatedUserMiddleware: authenticatedUserMiddleware,
			LoggerMiddleware:            loggerMiddleware,
		},
		Controllers: &controllers.Controllers{
			AuthController: authController,
			UserController: userController,
		},
	}
	routes.Setup()

	log.Print("finished app setup")
}

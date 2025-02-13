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

	tokenSecret := os.Getenv("JWT_TOKEN_SECRET")
	if tokenSecret == "" {
		log.Fatal("JWT_TOKEN_SECRET env var not found")
	}

	refreshTokenSecret := os.Getenv("JWT_REFRESH_TOKEN_SECRET")
	if refreshTokenSecret == "" {
		log.Fatal("JWT_REFRESH_TOKEN_SECRET env var not found")
	}

	// Setup repositories
	userRepository := repositories.NewPSQLUserRepository(app.DB)
	refreshTokenRepository := repositories.NewPSQLRefreshTokenRepository(app.DB)

	// Setup services
	hashService := services.NewHashService()
	jwtService := services.NewJWTService(tokenSecret, refreshTokenSecret, refreshTokenRepository, hashService)
	userService := services.NewUserService(userRepository, hashService)

	// Setup controllers
	authController := &controllers.AuthController{
		UserService: userService,
		HashService: hashService,
		JWTService:  jwtService,
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

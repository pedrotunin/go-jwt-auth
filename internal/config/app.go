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
		log.Panic("JWT_TOKEN_SECRET env var not found")
	}

	refreshTokenSecret := os.Getenv("JWT_REFRESH_TOKEN_SECRET")
	if refreshTokenSecret == "" {
		log.Panic("JWT_REFRESH_TOKEN_SECRET env var not found")
	}

	// Setup repositories
	userRepository := repositories.NewPSQLUserRepository(app.DB)
	refreshTokenRepository := repositories.NewPSQLRefreshTokenRepository(app.DB)
	evtRepository := repositories.NewPSQLEmailVerificationTokenRepository(app.DB)
	appRepository := repositories.NewPSQLAppRepository(app.DB)

	// Setup services
	sendGridMailerService := services.NewSendGridMailerService(
		os.Getenv("SENDGRID_SENDER_NAME"),
		os.Getenv("SENDGRID_SENDER_EMAIL"),
		os.Getenv("SENDGRID_API_KEY"),
	)
	hashService := services.NewHashService()
	jwtService := services.NewJWTService(tokenSecret, refreshTokenSecret, refreshTokenRepository, hashService)
	userService := services.NewUserService(userRepository, hashService)
	evtService := services.NewEmailVerificationTokenService(evtRepository)
	appService := services.NewAppService(appRepository)

	// Setup controllers
	authController := &controllers.AuthController{
		UserService: userService,
		HashService: hashService,
		JWTService:  jwtService,
	}
	userController := controllers.NewUserController(userService, evtService, sendGridMailerService)
	appController := &controllers.AppController{
		AppService: appService,
	}

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
			AppController:  appController,
		},
	}
	routes.Setup()

	log.Print("finished app setup")
}

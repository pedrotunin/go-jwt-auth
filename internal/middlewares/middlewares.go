package middlewares

type Middlewares struct {
	AuthenticatedUserMiddleware IAuthenticatedUserMiddleware
	LoggerMiddleware            ILoggerMiddleware
}

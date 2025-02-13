package middlewares

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

type ILoggerMiddleware interface {
	LogRequest() gin.HandlerFunc
}

type LoggerMiddleware struct {
}

func NewLoggerMiddleware() ILoggerMiddleware {
	return &LoggerMiddleware{}
}

func (lg *LoggerMiddleware) LogRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		c.Next()

		latency := time.Since(t)

		log.Printf("Method: %s | URL: %s | Status: %d | Duration: %v\n",
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			latency,
		)

	}
}

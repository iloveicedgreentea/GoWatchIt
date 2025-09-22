package webhooks

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/iloveicedgreentea/gowatchit/pkg/logger"
	"github.com/iloveicedgreentea/gowatchit/services/gowatchit/app"
)

// http listener
func NewRouter(ctx context.Context, app *app.App) (*gin.Engine, error) {
	log := logger.GetLoggerFromContext(ctx)

	log.Debug("Making new router")
	router := gin.Default()

	router.Use(addOrigins())

	addRoutes(router, app)

	err := router.SetTrustedProxies(nil)
	if err != nil {
		logger.Fatal("Failed to set trusted proxies: ", err)
	}

	return router, nil
}

func addOrigins() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Define allowed origins
		allowedOrigins := map[string]bool{
			"http://localhost:5173":            true, // bun
			"http://localhost:3000":            true, // nginx
			"http://host.docker.internal:3000": true, // docker
			"http://host.docker.internal:5173": true, // docker
		}

		// Check if origin is allowed and set the header
		if allowedOrigins[origin] {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

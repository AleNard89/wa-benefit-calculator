package middlewares

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func CORSMiddleware() gin.HandlerFunc {
	allowedOrigins := viper.GetStringSlice("api.allowedOrigins")

	cfg := cors.Config{
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"},
		AllowHeaders: []string{
			"Authorization",
			"Content-Type",
			"X-Requested-With",
			"Accept",
			"Origin",
			"Cache-Control",
			"X-Company-Id",
		},
		ExposeHeaders:    []string{},
		AllowCredentials: true,
		MaxAge:           10 * time.Minute,
		AllowOrigins:     allowedOrigins,
	}

	return cors.New(cfg)
}

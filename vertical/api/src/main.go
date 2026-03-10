package main

import (
	"fmt"
	"os"
	"time"

	"benefit-calculator-api/auth"
	"benefit-calculator-api/chat"
	"benefit-calculator-api/core"
	_ "benefit-calculator-api/core/logger"
	"benefit-calculator-api/db"
	"benefit-calculator-api/orchestrator"
	"benefit-calculator-api/orgs"
	"benefit-calculator-api/processes"

	m "benefit-calculator-api/core/middlewares"

	ginzap "github.com/gin-contrib/zap"
	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func init() {
	viper.SetConfigFile(fmt.Sprintf("./%s", os.Getenv("API_SETTINGS")))
	if err := viper.ReadInConfig(); err != nil {
		zap.S().Fatal(fmt.Sprintf("Error reading settings file, %s", err))
	} else {
		zap.S().Info("Successfully read settings file", viper.ConfigFileUsed())
	}
}

func main() {
	defer db.Close()

	zap.S().Info("Reading settings...")
	_ = viper.GetString("environment")

	decimal.MarshalJSONWithoutQuotes = true

	zap.S().Info("Initialising gin router...")
	r := gin.New()
	r.MaxMultipartMemory = 20 * 1024 * 1024 // 20MB

	// Middleware chain
	r.Use(gin.Recovery())
	r.Use(ginzap.Ginzap(zap.L(), time.RFC3339, true))
	r.Use(m.SecurityHeadersMiddleware())
	r.Use(m.CORSMiddleware())
	r.Use(auth.AuthenticationMiddleware)
	r.Use(orgs.CompanyMiddleware)

	// WebSocket
	r.GET("/api/ws", gin.WrapF(core.HandleWs))

	// API routes
	api := r.Group("/api")

	// Health check
	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Module routes
	auth.RoutesRegister(api.Group("/auth"))
	orgs.RoutesRegister(api.Group("/orgs"))
	processes.RoutesRegister(api.Group("/processes"))
	chat.RoutesRegister(api.Group("/chat"))
	orchestrator.RoutesRegister(api.Group("/orchestrator"))

	// Ensure storage folders exist for seeded companies
	orgs.EnsureCompanyFolders()

	zap.S().Info("Starting server...")
	r.Run() // listen and serve on 0.0.0.0:8080
}

package orchestrator

import "github.com/gin-gonic/gin"

func RoutesRegister(router *gin.RouterGroup) {
	// Connectors (admin/superuser only)
	router.GET("/connectors", ListConnectors)
	router.POST("/connectors", CreateConnector)
	router.PUT("/connectors/:id", UpdateConnector)
	router.POST("/connectors/:id/test", TestConnector)
	router.POST("/connectors/:id/sync", SyncConnectorHandler)

	// Dashboard stats (all authenticated users)
	router.GET("/dashboard-stats", GetDashboardStats)

	// Process names (all authenticated users)
	router.GET("/process-names", ListProcessNames)
	router.GET("/bot-names", ListBotNames)

	// Jobs (all authenticated users)
	router.GET("/jobs", ListJobs)
	router.GET("/jobs/:id", GetJob)

	// Queue Items (all authenticated users)
	router.GET("/queue-items", ListQueueItems)

	// Schedules (all authenticated users)
	router.GET("/schedules", ListSchedules)

	// Queue Definitions (all authenticated users)
	router.GET("/queue-definitions", ListQueueDefinitions)

	// Process-Queue Map
	router.GET("/process-queue-map", ListProcessQueueMaps)
	router.POST("/process-queue-map", CreateProcessQueueMap)
	router.DELETE("/process-queue-map/:id", DeleteProcessQueueMap)
}

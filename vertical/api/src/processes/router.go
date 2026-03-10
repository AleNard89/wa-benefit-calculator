package processes

import "github.com/gin-gonic/gin"

func RoutesRegister(router *gin.RouterGroup) {
	router.GET("", GetProcesses)
	router.GET("/stats", GetProcessStats)
	router.GET("/:id", GetProcess)
	router.POST("", CreateProcess)
	router.PUT("/:id", UpdateProcess)
	router.PATCH("/:id/status", UpdateProcessStatus)
	router.POST("/:id/recalculate", RecalculateProcess)
	router.DELETE("/:id", DeleteProcess)
	router.POST("/:id/restore", RestoreProcess)

	// Document upload/download/delete
	router.POST("/:id/document", UploadDocument)
	router.GET("/:id/document", DownloadDocument)
	router.DELETE("/:id/document", DeleteDocument)
}

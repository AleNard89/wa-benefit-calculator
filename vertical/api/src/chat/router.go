package chat

import "github.com/gin-gonic/gin"

func RoutesRegister(router *gin.RouterGroup) {
	router.GET("/conversations", GetConversations)
	router.POST("/conversations", CreateConversation)
	router.DELETE("/conversations/:id", DeleteConversation)
	router.GET("/conversations/:id/messages", GetMessages)
	router.POST("/conversations/:id/messages", SendMessage)
	router.POST("/index", IndexDocuments)
}

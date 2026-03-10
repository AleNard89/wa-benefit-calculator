package auth

import "github.com/gin-gonic/gin"

func RoutesRegister(router *gin.RouterGroup) {
	// authentication
	router.POST("/token/obtain", ObtainToken)
	router.POST("/token/refresh", RefreshToken)
	router.POST("/logout", Logout)
	router.GET("/whoami", WhoAmI)

	// permissions
	router.GET("/permission", GetPermissions)

	// roles
	router.GET("/role", GetRoles)
	router.GET("/role/:id", GetRole)
	router.POST("/role", CreateRole)
	router.PUT("/role/:id", UpdateRole)
	router.DELETE("/role/:id", DeleteRole)

	// users
	router.GET("/user", GetUsers)
	router.POST("/user", CreateUser)
	router.PUT("/user/:id", UpdateUser)
	router.PUT("/user/:id/password", UpdateUserPassword)
	router.DELETE("/user/:id", DeleteUser)
}

package orgs

import "github.com/gin-gonic/gin"

func RoutesRegister(router *gin.RouterGroup) {
	router.GET("/company", GetCompanies)
	router.GET("/company/:id", GetCompany)
	router.POST("/company", CreateCompany)
	router.PUT("/company/:id", UpdateCompany)
	router.DELETE("/company/:id", DeleteCompany)

	router.GET("/area", GetAreas)
	router.POST("/area", CreateArea)
	router.PUT("/area/:id", UpdateArea)
	router.DELETE("/area/:id", DeleteArea)

	router.GET("/company/:id/areas", GetCompanyAreas)
	router.POST("/company/:id/areas", CreateCompanyArea)

	router.GET("/user/:userId/areas", GetUserAreas)
	router.PUT("/user/:userId/areas", SetUserAreas)
}

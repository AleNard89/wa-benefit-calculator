package decorators

import (
	"benefit-calculator-api/core/httpx"
	"benefit-calculator-api/core/shared"
	"benefit-calculator-api/core/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SuperuserRequired(handler func(c *gin.Context)) func(c *gin.Context) {
	return func(c *gin.Context) {
		iuser, exists := c.Get("user")
		user, ok := iuser.(shared.User)

		if !exists || !ok || user.IsAnonymous() {
			c.JSON(http.StatusUnauthorized, httpx.ErrorResponse{Message: "Missing or invalid token"})
			return
		}
		if user.IsSuperUser() {
			handler(c)
			return
		}

		zap.S().Warnw("Superuser access denied", "ip", c.ClientIP())
		c.JSON(http.StatusForbidden, httpx.ErrorResponse{Message: "You don't have the rights to see the requested content"})
	}
}

func CompanyPermissionRequired(perms []string, handler func(c *gin.Context)) func(c *gin.Context) {
	return func(c *gin.Context) {
		iuser, exists := c.Get("user")
		user, ok := iuser.(shared.User)

		if !exists || !ok || user.IsAnonymous() {
			c.JSON(http.StatusUnauthorized, httpx.ErrorResponse{Message: "Missing or invalid token"})
			return
		}

		if user.IsSuperUser() {
			handler(c)
			return
		}

		companyID := httpx.HeaderCompanyID(c)
		if companyID == 0 {
			c.JSON(http.StatusUnauthorized, httpx.ErrorResponse{Message: "Missing company id"})
			return
		}

		if utils.ContainsAtLeastOne(perms, user.AllCompanyPermissionsCodes(companyID)) {
			handler(c)
			return
		}

		c.JSON(http.StatusForbidden, httpx.ErrorResponse{Message: "You don't have the rights to see the requested content"})
	}
}

func CompanyRequiredAndPermissionRequired(perms []string, handler func(c *gin.Context)) func(c *gin.Context) {
	return func(c *gin.Context) {
		companyID := httpx.HeaderCompanyID(c)
		if companyID == 0 {
			c.JSON(http.StatusUnauthorized, httpx.ErrorResponse{Message: "Missing company id"})
			return
		}
		CompanyPermissionRequired(perms, handler)(c)
	}
}

func PermissionRequired(perms []string, handler func(c *gin.Context)) func(c *gin.Context) {
	return func(c *gin.Context) {
		iuser, exists := c.Get("user")
		user, ok := iuser.(shared.User)

		if !exists || !ok || user.IsAnonymous() {
			c.JSON(http.StatusUnauthorized, httpx.ErrorResponse{Message: "Missing or invalid token"})
			return
		}

		if user.IsSuperUser() {
			handler(c)
			return
		}

		if utils.ContainsAtLeastOne(perms, user.AllPermissionsCodes()) {
			handler(c)
			return
		}

		c.JSON(http.StatusForbidden, httpx.ErrorResponse{Message: "You don't have the rights to see the requested content"})
	}
}

func AuthenticationRequired(handler func(c *gin.Context)) func(c *gin.Context) {
	return func(c *gin.Context) {
		iuser, exists := c.Get("user")
		user, ok := iuser.(shared.User)

		if !exists || !ok || user.IsAnonymous() {
			c.JSON(http.StatusUnauthorized, httpx.ErrorResponse{Message: "Missing or invalid token"})
			return
		}

		handler(c)
	}
}

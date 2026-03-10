package orgs

import (
	"benefit-calculator-api/core/shared"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CompanyMiddleware(c *gin.Context) {
	companyHeader := c.GetHeader("X-Company-Id")
	if companyHeader == "" {
		c.Set("companyID", 0)
		return
	}
	companyID, err := strconv.Atoi(companyHeader)
	if err != nil || companyID <= 0 {
		c.Set("companyID", 0)
		return
	}

	iuser, exists := c.Get("user")
	if !exists {
		c.Set("companyID", 0)
		return
	}
	user, ok := iuser.(shared.User)
	if !ok || user.IsAnonymous() {
		c.Set("companyID", 0)
		return
	}

	if !user.IsSuperUser() && !user.BelongsToCompany(companyID) {
		zap.S().Warnw("User tried to access company without membership",
			"companyID", companyID, "ip", c.ClientIP())
		c.Set("companyID", 0)
		return
	}

	c.Set("companyID", companyID)
}

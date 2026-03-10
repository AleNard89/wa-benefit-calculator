package orgs

import "github.com/gin-gonic/gin"

type CompanyBody struct {
	Name     string `json:"name" binding:"required"`
	ParentID *int   `json:"parentId"`
}

func (b *CompanyBody) Bind(c *gin.Context) error {
	return c.ShouldBindJSON(b)
}

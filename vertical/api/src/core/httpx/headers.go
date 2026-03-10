package httpx

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GenerateEtag(obj any) (string, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	hash := sha1.Sum(data)
	return fmt.Sprintf(`"%s"`, hex.EncodeToString(hash[:])), nil
}

func AddEtagHeader(c *gin.Context, obj any) error {
	etag, err := GenerateEtag(obj)
	if err == nil {
		c.Header("ETag", etag)
	}
	return err
}

func CheckEtagMatch(c *gin.Context, obj any) bool {
	etag, err := GenerateEtag(obj)
	if err != nil {
		zap.S().Errorf("Error generating etag: %v", err)
		return false
	}
	return c.GetHeader("If-Match") == etag
}

func HeaderCompanyID(c *gin.Context) int {
	icompanyID, exists := c.Get("companyID")
	if !exists {
		return 0
	}
	companyID, ok := icompanyID.(int)
	if !ok {
		return 0
	}
	return companyID
}

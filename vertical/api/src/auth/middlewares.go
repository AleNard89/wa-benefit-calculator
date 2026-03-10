package auth

import (
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RetrieveAccessToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	extractedToken := strings.Split(authHeader, "Bearer ")
	if len(extractedToken) == 2 {
		return strings.TrimSpace(extractedToken[1])
	}
	return ""
}

func AuthenticationMiddleware(c *gin.Context) {
	anonymousUser := &User{}
	stringToken := RetrieveAccessToken(c)
	claim, err := JWTAuthService().ValidateToken(stringToken)
	if err != nil {
		zap.S().Warnw("Token JWT invalido", "ip", c.ClientIP(), "error", err)
		c.Set("user", anonymousUser)
		return
	}
	userService := new(UserService)
	user, err := userService.GetByEmail(claim.Email)
	if err != nil {
		zap.S().Warn("AuthenticationMiddleware, cannot find user for valid token, email: ", claim.Email)
		c.Set("user", anonymousUser)
		return
	}
	c.Set("user", user)
}

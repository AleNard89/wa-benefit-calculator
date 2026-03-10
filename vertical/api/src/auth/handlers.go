package auth

import (
	"context"
	"benefit-calculator-api/core"
	"benefit-calculator-api/core/httpx"
	"fmt"
	"net/http"
	"strconv"
	"time"

	d "benefit-calculator-api/core/decorators"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

const (
	bruteForceEmailMaxAttempts = 3
	bruteForceIPMaxAttempts    = 10
	bruteForceLockoutDuration  = 15 * time.Minute
	bruteForceLockoutSeconds   = "900"
)

func bruteForceEmailKey(email string) string {
	return fmt.Sprintf("brute_force:email:%s", email)
}

func bruteForceIPKey(ip string) string {
	return fmt.Sprintf("brute_force:ip:%s", ip)
}

func isBruteForceBlocked(ctx context.Context, email, ip string) (bool, string) {
	rdb := core.GetRedisClient()
	emailCount, _ := rdb.Get(ctx, bruteForceEmailKey(email)).Int()
	if emailCount >= bruteForceEmailMaxAttempts {
		return true, "email"
	}
	ipCount, _ := rdb.Get(ctx, bruteForceIPKey(ip)).Int()
	if ipCount >= bruteForceIPMaxAttempts {
		return true, "ip"
	}
	return false, ""
}

func bruteForceIncrement(ctx context.Context, email, ip string) {
	rdb := core.GetRedisClient()
	emailKey := bruteForceEmailKey(email)
	emailCount, _ := rdb.Incr(ctx, emailKey).Result()
	if emailCount == 1 {
		rdb.Expire(ctx, emailKey, bruteForceLockoutDuration)
	}
	ipKey := bruteForceIPKey(ip)
	ipCount, _ := rdb.Incr(ctx, ipKey).Result()
	if ipCount == 1 {
		rdb.Expire(ctx, ipKey, bruteForceLockoutDuration)
	}
}

func bruteForceClear(ctx context.Context, email, ip string) {
	rdb := core.GetRedisClient()
	rdb.Del(ctx, bruteForceEmailKey(email), bruteForceIPKey(ip))
}

func ObtainToken(ctx *gin.Context) {
	var credential SignInCredentialsBody
	err := ctx.ShouldBindJSON(&credential)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpx.ErrorResponse{Message: "Missing e-mail or password"})
		return
	}

	redisCtx := context.Background()
	clientIP := ctx.ClientIP()
	if blocked, reason := isBruteForceBlocked(redisCtx, credential.Email, clientIP); blocked {
		zap.S().Warnw("Brute force bloccato", "email", credential.Email, "ip", clientIP, "motivo", reason)
		ctx.Header("Retry-After", bruteForceLockoutSeconds)
		ctx.JSON(http.StatusTooManyRequests, httpx.ErrorResponse{Message: "Too many failed login attempts. Try again later."})
		return
	}

	authService := AuthenticationService{}
	isUserAuthenticated, _ := authService.Authenticate(credential.Email, credential.Password)

	if isUserAuthenticated {
		bruteForceClear(redisCtx, credential.Email, clientIP)
		jwtService := JWTAuthService()
		tokenDuration := viper.GetInt("jwt.accessTokenValidityMinutes")
		token := jwtService.GenerateToken(credential.Email, tokenDuration)
		refreshTokenDuration := viper.GetInt("jwt.refreshTokenValidityMinutes")
		refreshToken := jwtService.GenerateToken(credential.Email, refreshTokenDuration)

		if refreshToken != "" && refreshTokenDuration > 0 {
			setRefreshTokenCookie(ctx, refreshToken, refreshTokenDuration)
		}
		ctx.JSON(http.StatusCreated, SignInSuccessResponse{AccessToken: token})
	} else {
		bruteForceIncrement(redisCtx, credential.Email, clientIP)
		clearRefreshTokenCookie(ctx)
		zap.S().Warnw("Login fallito", "email", credential.Email, "ip", ctx.ClientIP())
		ctx.JSON(http.StatusUnauthorized, httpx.ErrorResponse{Message: "Wrong authentication credentials"})
	}
}

func RefreshToken(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie("refreshToken")
	if err != nil || refreshToken == "" {
		ctx.JSON(http.StatusBadRequest, httpx.ErrorResponse{Message: "Missing refresh token"})
		return
	}

	claim, err := JWTAuthService().ValidateToken(refreshToken)
	if err != nil {
		clearRefreshTokenCookie(ctx)
		ctx.JSON(http.StatusUnauthorized, httpx.ErrorResponse{Message: "Invalid refresh token"})
		return
	}

	jwtService := JWTAuthService()
	tokenDuration := viper.GetInt("jwt.accessTokenValidityMinutes")
	token := jwtService.GenerateToken(claim.Email, tokenDuration)
	refreshTokenDuration := viper.GetInt("jwt.refreshTokenValidityMinutes")
	newRefreshToken := jwtService.GenerateToken(claim.Email, refreshTokenDuration)

	if newRefreshToken != "" && refreshTokenDuration > 0 {
		setRefreshTokenCookie(ctx, newRefreshToken, refreshTokenDuration)
	}
	ctx.JSON(http.StatusOK, RefreshTokenSuccessResponse{AccessToken: token})
}

func refreshTokenCookieConfig() (string, bool) {
	return viper.GetString("api.refreshTokenCookieDomain"), viper.GetString("scheme") == "https"
}

func setRefreshTokenCookie(ctx *gin.Context, value string, maxAgeMinutes int) {
	if maxAgeMinutes <= 0 {
		return
	}
	domain, secure := refreshTokenCookieConfig()
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("refreshToken", value, maxAgeMinutes*60, "/", domain, secure, true)
}

func clearRefreshTokenCookie(ctx *gin.Context) {
	domain, secure := refreshTokenCookieConfig()
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("refreshToken", "", -1, "/", domain, secure, true)
}

func Logout(ctx *gin.Context) {
	clearRefreshTokenCookie(ctx)
	ctx.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

func WhoAmI(ctx *gin.Context) {
	iuser, exists := ctx.Get("user")
	user, ok := iuser.(*User)
	if !exists || !ok || user.IsAnonymous() {
		ctx.JSON(http.StatusUnauthorized, httpx.ErrorResponse{Message: "User not found"})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

// Users

func getUsers(c *gin.Context) {
	u := MustCurrentUser(c)
	companyID := httpx.HeaderCompanyID(c)
	if u.IsSuperuser {
		companyID = 0
	}
	service := UserService{}
	users, err := service.GetAll(companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}
	c.JSON(http.StatusOK, users)
}

var GetUsers = d.CompanyPermissionRequired([]string{AuthUserRead}, getUsers)

func createUser(c *gin.Context) {
	payload := UserBodyWithPassword{}
	err := payload.Bind(c)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}
	user := User{}
	err = user.ApplyPayload(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}
	c.JSON(http.StatusCreated, user)
}

var CreateUser = d.CompanyPermissionRequired([]string{AuthUserCreate}, createUser)

func updateUser(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Param("id"))
	service := UserService{}
	user, err := service.GetByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, httpx.ErrorResponse{Message: "User not found"})
		return
	}
	currentUser := MustCurrentUser(c)
	if user.IsSuperuser && !currentUser.IsSuperuser {
		c.JSON(http.StatusForbidden, httpx.ErrorResponse{Message: "Only a superuser can modify another superuser"})
		return
	}
	payload := UserBody{}
	err = payload.Bind(c)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}
	err = user.ApplyBasePayload(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}
	c.JSON(http.StatusOK, user)
}

var UpdateUser = d.CompanyPermissionRequired([]string{AuthUserCreate}, updateUser)

func updateUserPassword(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Param("id"))
	service := UserService{}
	user, err := service.GetByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, httpx.ErrorResponse{Message: "User not found"})
		return
	}

	u := MustCurrentUser(c)
	if !u.IsSuperuser && user.ID != u.ID {
		c.JSON(http.StatusForbidden, httpx.ErrorResponse{Message: "Can't change password for another user"})
		return
	}

	payload := ChangePasswordBody{}
	err = payload.Bind(c)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}

	if !u.IsSuperuser {
		passwordHash, hashErr := service.GetPasswordHashByID(user.ID)
		if hashErr != nil {
			c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Message: "error retrieving password"})
			return
		}
		if bcErr := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(payload.CurrentPassword)); bcErr != nil {
			c.JSON(http.StatusBadRequest, httpx.ErrorResponse{Message: "current password is incorrect"})
			return
		}
	}

	err = user.UpdatePassword(payload.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}
	zap.S().Infow("Password aggiornata", "user_id", user.ID, "email", user.Email, "ip", c.ClientIP())
	c.JSON(http.StatusOK, user)
}

var UpdateUserPassword = d.AuthenticationRequired(updateUserPassword)

func deleteUser(c *gin.Context) {
	u := MustCurrentUser(c)
	companyID := httpx.HeaderCompanyID(c)

	userID, _ := strconv.Atoi(c.Param("id"))
	service := UserService{}
	user, err := service.GetByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, httpx.ErrorResponse{Message: "User not found"})
		return
	}

	if u.IsSuperuser || (len(user.CompanyRoles) == 1 && user.CompanyRoles[0].Company.ID == companyID && !user.IsSuperuser) {
		if err := user.Delete(); err != nil {
			c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
			return
		}
	} else if !user.IsSuperuser {
		if err := user.RemoveCompany(companyID); err != nil {
			c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
			return
		}
	} else {
		c.JSON(http.StatusForbidden, httpx.ErrorResponse{Message: "You can't delete a superuser"})
		return
	}
	zap.S().Warnw("Utente eliminato", "deleted_user_id", user.ID, "deleted_email", user.Email, "ip", c.ClientIP())
	c.JSON(http.StatusNoContent, nil)
}

var DeleteUser = d.CompanyPermissionRequired([]string{AuthUserDelete}, deleteUser)

// Roles

func getRoles(c *gin.Context) {
	service := RoleService{}
	roles, err := service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}
	c.JSON(http.StatusOK, roles)
}

var GetRoles = d.PermissionRequired([]string{AuthRoleRead}, getRoles)

func getRole(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}
	service := RoleService{}
	role, err := service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}
	httpx.JSONWithEtag(c, role)
}

var GetRole = d.PermissionRequired([]string{AuthRoleRead}, getRole)

func createRole(c *gin.Context) {
	payload := RoleBody{}
	err := payload.Bind(c)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}
	role := Role{}
	err = role.ApplyPayload(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}
	c.JSON(http.StatusCreated, role)
}

var CreateRole = d.PermissionRequired([]string{AuthRoleCreate}, createRole)

func updateRole(c *gin.Context) {
	roleID, _ := strconv.Atoi(c.Param("id"))
	service := RoleService{}
	role, err := service.GetByID(roleID)
	if err != nil {
		c.JSON(http.StatusNotFound, httpx.ErrorResponse{Message: "Role not found"})
		return
	}
	payload := RoleBody{}
	err = payload.Bind(c)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}
	err = role.ApplyPayload(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}
	zap.S().Infow("Ruolo aggiornato", "role_id", role.ID, "role_name", role.Name, "ip", c.ClientIP())
	c.JSON(http.StatusOK, role)
}

var UpdateRole = d.PermissionRequired([]string{AuthRoleUpdate}, updateRole)

func deleteRole(c *gin.Context) {
	roleID, _ := strconv.Atoi(c.Param("id"))
	service := RoleService{}
	role, err := service.GetByID(roleID)
	if err != nil {
		c.JSON(http.StatusNotFound, httpx.ErrorResponse{Message: "Role not found"})
		return
	}
	if err := role.Delete(); err != nil {
		c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

var DeleteRole = d.PermissionRequired([]string{AuthRoleDelete}, deleteRole)

// Permissions

func getPermissions(c *gin.Context) {
	service := PermissionService{}
	perms, err := service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}
	c.JSON(http.StatusOK, perms)
}

var GetPermissions = d.PermissionRequired([]string{AuthRoleRead}, getPermissions)

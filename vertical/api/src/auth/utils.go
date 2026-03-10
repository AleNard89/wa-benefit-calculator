package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func MustCurrentUser(c *gin.Context) *User {
	return c.MustGet("user").(*User)
}

func HashPassword(password string) (string, error) {
	pwd := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(pwd, 12)
	if err != nil {
		zap.S().Error(err)
		return "", err
	}
	return string(hash), nil
}

func ValidatePassword(password string) (bool, error) {
	if len(password) < 8 {
		return false, fmt.Errorf("password must be at least 8 characters long")
	}
	matched, _ := regexp.Match(`[a-z]`, []byte(password))
	if !matched {
		return false, fmt.Errorf("password must contain at least one lowercase letter")
	}
	matched, _ = regexp.Match(`[A-Z]`, []byte(password))
	if !matched {
		return false, fmt.Errorf("password must contain at least one uppercase letter")
	}
	matched, _ = regexp.Match(`[0-9]`, []byte(password))
	if !matched {
		return false, fmt.Errorf("password must contain at least one digit")
	}
	matched, _ = regexp.Match(`[$!%&]`, []byte(password))
	if !matched {
		return false, fmt.Errorf("password must contain at least one of the following special characters: $!%%&")
	}
	return true, nil
}

func SignedResetPasswordToken(email string, expire string) string {
	data := make(map[string]string)
	data["email"] = email
	data["expire"] = expire
	jsonData, _ := json.Marshal(data)
	mac := hmac.New(sha256.New, []byte(os.Getenv("JWT_SECRET")))
	mac.Write(jsonData)
	sig := base64.URLEncoding.EncodeToString(mac.Sum(nil))
	data["token"] = sig
	jsonParam, _ := json.Marshal(data)
	return base64.URLEncoding.EncodeToString(jsonParam)
}

func CheckResetPasswordToken(email string, expire string, sig string) bool {
	data := make(map[string]string)
	data["email"] = email
	data["expire"] = expire
	jsonData, _ := json.Marshal(data)
	mac := hmac.New(sha256.New, []byte(os.Getenv("JWT_SECRET")))
	mac.Write(jsonData)
	expected := base64.URLEncoding.EncodeToString(mac.Sum(nil))
	return expected == sig
}

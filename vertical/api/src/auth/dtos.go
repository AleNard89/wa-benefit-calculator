package auth

import (
	"errors"

	"github.com/gin-gonic/gin"
)

// Token

type SignInCredentialsBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInSuccessResponse struct {
	AccessToken string `json:"accessToken"`
}

type RefreshTokenSuccessResponse struct {
	AccessToken string `json:"accessToken"`
}

// Role

type RoleBody struct {
	Name          string `json:"name" binding:"required"`
	Description   string `json:"description"`
	PermissionIDs []int  `json:"permissionIds"`
}

func (b *RoleBody) Bind(c *gin.Context) error {
	return c.ShouldBind(b)
}

// User

type UserBody struct {
	Email       string            `json:"email" binding:"required"`
	IsSuperuser bool              `json:"isSuperuser"`
	FirstName   string            `json:"firstName" binding:"required"`
	LastName    string            `json:"lastName" binding:"required"`
	Companies   []UserCompanyBody `json:"companies"`
}

type UserCompanyBody struct {
	CompanyID int   `json:"companyId" binding:"required"`
	RoleIDs   []int `json:"roleIds"`
}

type UserBodyWithPassword struct {
	UserBody
	Password string `json:"password" binding:"required"`
}

func CheckUserBodyCompanies(u *User, companies []UserCompanyBody, requiredPerms []string) bool {
	for _, company := range companies {
		permissions := u.AllCompanyPermissionsCodes(company.CompanyID)
		if len(permissions) == 0 {
			return false
		}
		allowed := false
		for _, required := range requiredPerms {
			for _, perm := range permissions {
				if perm == required {
					allowed = true
					break
				}
			}
			if allowed {
				break
			}
		}
		if !allowed {
			return false
		}
	}
	return true
}

func (b *UserBodyWithPassword) Bind(c *gin.Context) error {
	u := MustCurrentUser(c)
	err := c.ShouldBind(b)
	if err != nil {
		return err
	}
	_, err = ValidatePassword(b.Password)
	if err != nil {
		return err
	}
	if b.IsSuperuser && !u.IsSuperuser {
		return errors.New("only a superuser can grant superuser privileges")
	}
	if !u.IsSuperuser && !CheckUserBodyCompanies(u, b.Companies, []string{AuthUserCreate}) {
		return errors.New("cannot set roles for companies without permission")
	}
	return nil
}

func (b *UserBody) Bind(c *gin.Context) error {
	u := MustCurrentUser(c)
	err := c.ShouldBind(b)
	if err != nil {
		return err
	}
	if b.IsSuperuser && !u.IsSuperuser {
		return errors.New("only a superuser can grant superuser privileges")
	}
	if !u.IsSuperuser && !CheckUserBodyCompanies(u, b.Companies, []string{AuthUserUpdate}) {
		return errors.New("cannot set roles for companies without permission")
	}
	return nil
}

// ChangePassword

type ChangePasswordBody struct {
	CurrentPassword string `json:"currentPassword"`
	Password        string `json:"password" binding:"required"`
}

func (b *ChangePasswordBody) Bind(c *gin.Context) error {
	err := c.ShouldBind(b)
	if err != nil {
		return err
	}
	_, err = ValidatePassword(b.Password)
	return err
}

// Reset password

type ResetPasswordRequestBody struct {
	Email string `json:"email" binding:"required"`
}

type ResetPasswordRequestResponse struct {
	Email   string `json:"email"`
	Result  bool   `json:"result"`
	Message string `json:"message"`
}

type ResetPasswordConfirmBody struct {
	Token    string `json:"token"    binding:"required"`
	Email    string `json:"email"    binding:"required"`
	Expire   string `json:"expire"   binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ResetPasswordConfirmResponse struct {
	Email   string `json:"email"`
	Result  bool   `json:"result"`
	Message string `json:"message"`
}

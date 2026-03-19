package auth

import (
	"strings"
	"time"

	"orbita-api/orgs"

	"go.uber.org/zap"
)

const (
	AuthUserRead   = "auth:user.read"
	AuthUserCreate = "auth:user.create"
	AuthUserUpdate = "auth:user.update"
	AuthUserDelete = "auth:user.delete"
)

type User struct {
	ID           int            `db:"id" json:"id,omitempty"`
	Email        string         `db:"email" json:"email"`
	Password     string         `db:"password" json:"-"`
	IsSuperuser  bool           `db:"is_superuser" json:"isSuperuser,omitempty"`
	FirstName    string         `db:"firstname" json:"firstName"`
	LastName     string         `db:"lastname" json:"lastName"`
	CreatedAt    time.Time      `db:"created_at" json:"createdAt"`
	UpdatedAt    time.Time      `db:"updated_at" json:"updatedAt"`
	CompanyRoles []CompanyRoles `json:"companyRoles"`
}

func (u *User) GetID() int {
	return u.ID
}

func (u *User) IsAnonymous() bool {
	return u.ID == 0
}

func (u *User) IsSuperUser() bool {
	return u.IsSuperuser
}

func (u *User) SetPassword(password string) error {
	hash, err := HashPassword(password)
	if err != nil {
		zap.S().Error(err)
		return err
	}
	u.Password = string(hash)
	return nil
}

func (u *User) BelongsToCompany(companyID int) bool {
	for _, cr := range u.CompanyRoles {
		if cr.Company.ID == companyID {
			return true
		}
	}
	return false
}

func (u *User) AllCompanyPermissionsCodes(companyID int) []string {
	return u.permissionCodesFilter(func(cr CompanyRoles) bool {
		return cr.Company.ID == companyID
	})
}

func (u *User) AllPermissionsCodes() []string {
	return u.permissionCodesFilter(nil)
}

func (u *User) permissionCodesFilter(predicate func(CompanyRoles) bool) []string {
	permissionSet := make(map[string]bool)
	for _, cr := range u.CompanyRoles {
		if predicate != nil && !predicate(cr) {
			continue
		}
		for _, role := range cr.Roles {
			for _, perm := range role.Permissions {
				permissionSet[perm.Code] = true
			}
		}
	}
	permissions := make([]string, 0, len(permissionSet))
	for p := range permissionSet {
		permissions = append(permissions, p)
	}
	return permissions
}

func (u *User) UpdatePassword(password string) error {
	err := u.SetPassword(password)
	if err != nil {
		return err
	}
	s := UserService{}
	return s.UpdatePassword(u)
}

func (u *User) Save() error {
	service := UserService{}
	u.UpdatedAt = time.Now()
	if u.ID > 0 {
		return service.Update(u)
	}
	u.CreatedAt = time.Now()
	return service.Insert(u)
}

func (u *User) SaveCompaniesAndRoles(userCompanyBody []UserCompanyBody) error {
	service := UserService{}
	err := service.UpdateCompaniesAndRoles(u, userCompanyBody)
	if err != nil {
		return err
	}
	companyRoles, err := service.GetCompaniesAndRoles(u)
	if err == nil {
		u.CompanyRoles = companyRoles
	}
	return err
}

func (u *User) RemoveCompany(companyID int) error {
	service := UserService{}
	err := service.RemoveCompany(u, companyID)
	if err != nil {
		return err
	}
	companyRoles, err := service.GetCompaniesAndRoles(u)
	if err == nil {
		u.CompanyRoles = companyRoles
	}
	return err
}

func (u *User) ApplyPayload(payload UserBodyWithPassword) error {
	u.Email = strings.ToLower(payload.Email)
	u.Password = payload.Password
	u.IsSuperuser = payload.IsSuperuser
	u.FirstName = payload.FirstName
	u.LastName = payload.LastName
	u.SetPassword(payload.Password)

	err := u.Save()
	if err != nil {
		return err
	}

	err = u.SaveCompaniesAndRoles(payload.Companies)
	return err
}

func (u *User) ApplyBasePayload(payload UserBody) error {
	u.Email = strings.ToLower(payload.Email)
	u.IsSuperuser = payload.IsSuperuser
	u.FirstName = payload.FirstName
	u.LastName = payload.LastName

	err := u.Save()
	if err != nil {
		return err
	}

	err = u.SaveCompaniesAndRoles(payload.Companies)
	return err
}

func (u *User) Delete() error {
	service := UserService{}
	return service.Delete(u.ID)
}

// Permission model
const (
	AuthPermissionRead = "auth:permission.read"
)

type Permission struct {
	ID          int    `db:"id" json:"id"`
	App         string `db:"app" json:"app"`
	Code        string `db:"code" json:"code"`
	Description string `db:"description" json:"description"`
}

// Role model
const (
	AuthRoleRead   = "auth:role.read"
	AuthRoleCreate = "auth:role.create"
	AuthRoleUpdate = "auth:role.update"
	AuthRoleDelete = "auth:role.delete"
)

type Role struct {
	ID          int          `db:"id" json:"id"`
	Name        string       `db:"name" json:"name"`
	Description string       `db:"description" json:"description"`
	Permissions []Permission `db:"permissions" json:"permissions"`
}

func (r *Role) Save() error {
	service := RoleService{}
	if r.ID > 0 {
		return service.Update(r)
	}
	return service.Insert(r)
}

func (r *Role) SavePermissions(permissionIDs []int) error {
	service := RoleService{}
	err := service.UpdatePermissions(r, permissionIDs)
	if err != nil {
		return err
	}
	permissions, err := service.GetPermissions(r)
	if err == nil {
		r.Permissions = permissions
	}
	return err
}

func (r *Role) Delete() error {
	service := RoleService{}
	return service.Delete(r.ID)
}

func (r *Role) ApplyPayload(payload RoleBody) error {
	r.Name = payload.Name
	r.Description = payload.Description
	err := r.Save()
	if err != nil {
		return err
	}
	return r.SavePermissions(payload.PermissionIDs)
}

type CompanyRoles struct {
	Company orgs.Company `json:"company"`
	Roles   []Role       `json:"roles"`
}

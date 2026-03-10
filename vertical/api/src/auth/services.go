package auth

import (
	"context"
	"errors"
	"benefit-calculator-api/core/utils"
	"benefit-calculator-api/db"
	"fmt"
	"os"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

const (
	UsersTable            = "auth_users"
	RolesTable            = "auth_roles"
	PermissionsTable      = "auth_permissions"
	RolesPermissionsTable = "auth_roles_permissions"
	UsersRolesTable       = "auth_users_roles"
	UsersCompaniesTable   = "auth_users_companies"
)

// AuthenticationService

type AuthenticationService struct{}

func (s *AuthenticationService) Authenticate(email string, password string) (bool, *User) {
	var user User
	db := db.DB()

	stm := db.Builder.Select("*").From(UsersTable).Where(sq.Eq{"email": strings.ToLower(email)})
	sql, args, err := stm.ToSql()
	if err != nil {
		zap.S().Debugw("AuthenticationService, error building query: ", "error", err)
		return false, nil
	}
	err = pgxscan.Get(context.Background(), db.C, &user, sql, args...)
	if err != nil {
		zap.S().Debugw("AuthenticationService, user not found: ", "error", err)
		return false, nil
	}

	byteHash := []byte(user.Password)
	err = bcrypt.CompareHashAndPassword(byteHash, []byte(password))
	if err != nil {
		zap.S().Debugw("AuthenticationService, wrong password: ", "error", err)
		return false, nil
	}
	return true, &user
}

// JWTService

type JWTService interface {
	GenerateToken(email string, validity int) string
	ValidateToken(token string) (*JwtClaim, error)
}

type JwtClaim struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type jwtService struct {
	secretKey string
	issuer    string
}

func JWTAuthService() JWTService {
	return &jwtService{
		secretKey: mustSecretKey(),
		issuer:    "BenefitCalculator",
	}
}

func mustSecretKey() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic(errors.New("JWT_SECRET is not set"))
	}
	return secret
}

func (service *jwtService) GenerateToken(email string, validity int) string {
	expiresAt := time.Now().Local().Add(time.Minute * time.Duration(validity))
	claims := &JwtClaim{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			Issuer:    service.issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(service.secretKey))
	if err != nil {
		panic(err)
	}
	return t
}

func (service *jwtService) ValidateToken(signedToken string) (*JwtClaim, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JwtClaim{},
		func(token *jwt.Token) (any, error) {
			return []byte(service.secretKey), nil
		},
	)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*JwtClaim)
	if !ok {
		return nil, errors.New("failed to parse claims")
	}
	if claims.ExpiresAt.Unix() < time.Now().Local().Unix() {
		return nil, errors.New("JWT expired")
	}
	return claims, nil
}

// UserService

type UserService struct{}

func (s *UserService) GetAll(companyID int) ([]User, error) {
	db := db.DB()
	users := []User{}

	stm := SelectUsersWithCompanyRoles()
	if companyID != 0 {
		subquery, subargs := SelectUserIdsWithCompany(companyID)
		stm = stm.Where(sq.Expr(fmt.Sprintf("u.id IN (%s)", subquery), subargs...))
	}
	sql, args, _ := stm.ToSql()

	err := pgxscan.Select(context.Background(), db.C, &users, sql, args...)
	if err != nil {
		zap.S().Errorw("Error getting users", "error", err)
		return nil, err
	}
	return users, nil
}

func (s *UserService) GetByEmail(email string) (*User, error) {
	db := db.DB()
	user := User{}
	stm := SelectUsersWithCompanyRoles()
	stm = stm.Where(sq.Eq{"u.email": strings.ToLower(email)})
	sql, args, _ := stm.ToSql()
	err := pgxscan.Get(context.Background(), db.C, &user, sql, args...)
	return utils.RetZeroOnError(&user, err)
}

func (s *UserService) GetByID(id int) (*User, error) {
	db := db.DB()
	user := User{}
	stm := SelectUsersWithCompanyRoles()
	stm = stm.Where(sq.Eq{"u.id": id})
	sql, args, _ := stm.ToSql()
	err := pgxscan.Get(context.Background(), db.C, &user, sql, args...)
	return utils.RetZeroOnError(&user, err)
}

func (s *UserService) Insert(user *User) error {
	db := db.DB()
	stm := db.Builder.
		Insert(UsersTable).
		Columns("email", "password", "is_superuser", "firstname", "lastname", "created_at", "updated_at").
		Values(user.Email, user.Password, user.IsSuperuser, user.FirstName, user.LastName, user.CreatedAt, user.UpdatedAt).
		Suffix("RETURNING id")
	sql, args, err := stm.ToSql()
	if err != nil {
		return err
	}
	return pgxscan.Get(context.Background(), db.C, &user.ID, sql, args...)
}

func (s *UserService) Update(user *User) error {
	db := db.DB()
	stm := db.Builder.
		Update(UsersTable).
		Set("email", user.Email).
		Set("is_superuser", user.IsSuperuser).
		Set("firstname", user.FirstName).
		Set("lastname", user.LastName).
		Set("updated_at", user.UpdatedAt).
		Where(sq.Eq{"id": user.ID})
	sql, args, err := stm.ToSql()
	if err != nil {
		return err
	}
	_, err = db.C.Exec(context.Background(), sql, args...)
	return err
}

func (s *UserService) GetPasswordHashByID(id int) (string, error) {
	db := db.DB()
	var hash string
	stm := db.Builder.Select("password").From(UsersTable).Where(sq.Eq{"id": id})
	sql, args, _ := stm.ToSql()
	err := pgxscan.Get(context.Background(), db.C, &hash, sql, args...)
	return hash, err
}

func (s *UserService) UpdatePassword(user *User) error {
	db := db.DB()
	stm := db.Builder.
		Update(UsersTable).
		Set("password", user.Password).
		Where(sq.Eq{"id": user.ID})
	sql, args, err := stm.ToSql()
	if err != nil {
		return err
	}
	_, err = db.C.Exec(context.Background(), sql, args...)
	return err
}

func (s *UserService) UpdateCompaniesAndRoles(user *User, companies []UserCompanyBody) error {
	database := db.DB()
	ctx := context.Background()

	tx, err := database.C.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, fmt.Sprintf("DELETE FROM %s WHERE user_id = $1", UsersCompaniesTable), user.ID)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, fmt.Sprintf("DELETE FROM %s WHERE user_id = $1", UsersRolesTable), user.ID)
	if err != nil {
		return err
	}

	for _, company := range companies {
		_, err = tx.Exec(ctx,
			fmt.Sprintf("INSERT INTO %s (user_id, company_id) VALUES ($1, $2)", UsersCompaniesTable),
			user.ID, company.CompanyID)
		if err != nil {
			return err
		}
		for _, roleID := range company.RoleIDs {
			_, err = tx.Exec(ctx,
				fmt.Sprintf("INSERT INTO %s (user_id, role_id, company_id) VALUES ($1, $2, $3)", UsersRolesTable),
				user.ID, roleID, company.CompanyID)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit(ctx)
}

func (s *UserService) GetCompaniesAndRoles(user *User) ([]CompanyRoles, error) {
	db := db.DB()
	res := make([]CompanyRoles, 0)
	sql, args := SelectCompanyRolesForUser(user.ID)
	err := pgxscan.Select(context.Background(), db.C, &res, sql, args...)
	return utils.RetZeroOnError(res, err)
}

func (s *UserService) Delete(id int) error {
	db := db.DB()
	stm := db.Builder.Delete(UsersTable).Where(sq.Eq{"id": id})
	sql, args, err := stm.ToSql()
	if err != nil {
		return err
	}
	_, err = db.C.Exec(context.Background(), sql, args...)
	return err
}

func (s *UserService) RemoveCompany(user *User, companyID int) error {
	db := db.DB()
	stm := db.Builder.Delete(UsersCompaniesTable).Where(sq.Eq{"user_id": user.ID, "company_id": companyID})
	sql, args, err := stm.ToSql()
	if err != nil {
		return err
	}
	_, err = db.C.Exec(context.Background(), sql, args...)
	return err
}

// RoleService

type RoleService struct{}

func (s *RoleService) GetAll() ([]Role, error) {
	db := db.DB()
	roles := []Role{}
	stm, _ := SelectRolesWithPermissions(0)
	err := pgxscan.Select(context.Background(), db.C, &roles, stm)
	if err != nil {
		zap.S().Errorw("Error getting roles", "error", err)
		return nil, err
	}
	return roles, nil
}

func (s *RoleService) GetByID(id int) (*Role, error) {
	db := db.DB()
	role := Role{}
	stm, args := SelectRolesWithPermissions(id)
	err := pgxscan.Get(context.Background(), db.C, &role, stm, args...)
	return utils.RetZeroOnError(&role, err)
}

func (s *RoleService) Insert(role *Role) error {
	db := db.DB()
	stm := db.Builder.
		Insert(RolesTable).
		Columns("name", "description").
		Values(role.Name, role.Description).
		Suffix("RETURNING id")
	sql, args, err := stm.ToSql()
	if err != nil {
		return err
	}
	return db.C.QueryRow(context.Background(), sql, args...).Scan(&role.ID)
}

func (s *RoleService) Update(role *Role) error {
	db := db.DB()
	stm := db.Builder.
		Update(RolesTable).
		Set("name", role.Name).
		Set("description", role.Description).
		Where(sq.Eq{"id": role.ID})
	sql, args, err := stm.ToSql()
	if err != nil {
		return err
	}
	_, err = db.C.Exec(context.Background(), sql, args...)
	return err
}

func (s *RoleService) Delete(id int) error {
	db := db.DB()
	stm := db.Builder.Delete(RolesTable).Where(sq.Eq{"id": id})
	sql, args, err := stm.ToSql()
	if err != nil {
		return err
	}
	_, err = db.C.Exec(context.Background(), sql, args...)
	return err
}

func (s *RoleService) GetPermissions(role *Role) ([]Permission, error) {
	db := db.DB()
	perms := make([]Permission, 0)
	stm := db.Builder.
		Select("permissions.*").
		From(fmt.Sprintf("%s permissions", PermissionsTable)).
		Join(fmt.Sprintf("%s rp on rp.permission_id = permissions.id", RolesPermissionsTable)).
		Where(sq.Eq{"rp.role_id": role.ID})
	sql, args, _ := stm.ToSql()
	err := pgxscan.Select(context.Background(), db.C, &perms, sql, args...)
	return utils.RetZeroOnError(perms, err)
}

func (s *RoleService) UpdatePermissions(role *Role, permissionIDs []int) error {
	database := db.DB()
	ctx := context.Background()

	tx, err := database.C.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, fmt.Sprintf("DELETE FROM %s WHERE role_id = $1", RolesPermissionsTable), role.ID)
	if err != nil {
		return err
	}

	for _, permissionID := range permissionIDs {
		_, err = tx.Exec(ctx,
			fmt.Sprintf("INSERT INTO %s (role_id, permission_id) VALUES ($1, $2)", RolesPermissionsTable),
			role.ID, permissionID)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

// PermissionService

type PermissionService struct{}

func (s *PermissionService) GetAll() ([]Permission, error) {
	db := db.DB()
	perms := []Permission{}
	stm := db.Builder.Select("*").From(PermissionsTable)
	sql, args, err := stm.ToSql()
	if err != nil {
		return nil, err
	}
	err = pgxscan.Select(context.Background(), db.C, &perms, sql, args...)
	return utils.RetZeroOnError(perms, err)
}

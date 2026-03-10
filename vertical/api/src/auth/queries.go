package auth

import (
	"benefit-calculator-api/db"
	"benefit-calculator-api/orgs"
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

// Roles

func BaseRoleSelectStatement(where string) string {
	return fmt.Sprintf(`
	SELECT
		r.id,
		r.name,
		r.description,
		COALESCE(json_agg(
			json_build_object(
				'id', p.id,
				'app', p.app,
				'code', p.code,
				'description', p.description
			)
		) FILTER (WHERE p.id IS NOT NULL), '[]'::json) AS permissions
	FROM
		%s r
	LEFT JOIN
		%s rp ON r.id = rp.role_id
	LEFT JOIN
		%s p ON rp.permission_id = p.id
	%s
	GROUP BY
		r.id
	ORDER BY
		r.name
	`, RolesTable, RolesPermissionsTable, PermissionsTable, where)
}

func SelectRolesWithPermissions(roleID int) (string, []any) {
	where := ""
	args := []any{}
	if roleID > 0 {
		where = "WHERE r.id = $1"
		args = append(args, roleID)
	}
	return BaseRoleSelectStatement(where), args
}

// Users

func BaseUserRolesSelectStatement() string {
	return fmt.Sprintf(`
		COALESCE(
			(SELECT json_agg(json_build_object(
				'id', r.id,
				'name', r.name,
				'description', r.description,
				'permissions', COALESCE(
					(SELECT json_agg(json_build_object(
						'id', p.id,
						'app', p.app,
						'code', p.code,
						'description', p.description
					))
					FROM %s rp
					JOIN %s p ON rp.permission_id = p.id
					WHERE rp.role_id = r.id
				))
			))
			FROM %s ur
			JOIN %s r ON ur.role_id = r.id
			WHERE ur.user_id = u.id AND ur.company_id = c.id),
			'[]'::json
		)
	`, RolesPermissionsTable, PermissionsTable, UsersRolesTable, RolesTable)
}

func SelectUsersWithCompanyRoles() sq.SelectBuilder {
	companyRolesExpr := fmt.Sprintf(`
		COALESCE(
			(SELECT json_agg(json_build_object(
				'company', json_build_object(
					'id', c.id,
					'name', c.name
				),
				'roles', %s
			))
			FROM %s uc
			JOIN %s c ON uc.company_id = c.id
			WHERE uc.user_id = u.id),
			'[]'::json
		) AS company_roles`, BaseUserRolesSelectStatement(), UsersCompaniesTable, orgs.CompaniesTable)

	db := db.DB()
	stm := db.Builder.
		Select(
			"u.id",
			"u.email",
			"u.is_superuser",
			"u.firstname",
			"u.lastname",
			"u.created_at",
			"u.updated_at").
		Column(sq.Expr(companyRolesExpr)).
		From(UsersTable + " u")

	return stm
}

func SelectCompanyRolesForUser(userID int) (string, []any) {
	args := []any{userID}
	stm := fmt.Sprintf(`
	SELECT
    c.id AS "company.id",
    c.name AS "company.name",
    %s AS roles
	FROM
    %s uc
	INNER JOIN %s u ON u.id = uc.user_id
	INNER JOIN %s c ON c.id = uc.company_id
	WHERE u.id = $1
	`, BaseUserRolesSelectStatement(), UsersCompaniesTable, UsersTable, orgs.CompaniesTable)
	return stm, args
}

func SelectUserIdsWithCompany(companyID int) (string, []any) {
	db := db.DB()
	stm := db.Builder.
		Select("u.id").
		From(UsersTable + " u").
		InnerJoin(UsersCompaniesTable + " uc ON u.id = uc.user_id").
		Where(sq.Eq{"uc.company_id": companyID})
	sql, args, _ := stm.ToSql()
	return sql, args
}

package orgs

import "time"

const AreasTable = "orgs_areas"
const UsersAreasTable = "auth_users_areas"

type Area struct {
	ID        int       `db:"id" json:"id"`
	CompanyID int       `db:"company_id" json:"companyId"`
	Name      string    `db:"name" json:"name"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}

type AreaBody struct {
	Name string `json:"name" binding:"required"`
}

package orgs

import (
	"benefit-calculator-api/core/backends"
	"fmt"
	"regexp"
	"strings"
	"time"
)

const (
	CompaniesTable = "orgs_companies"

	OrgsCompanyRead   = "orgs:company.read"
	OrgsCompanyCreate = "orgs:company.create"
	OrgsCompanyUpdate = "orgs:company.update"
	OrgsCompanyDelete = "orgs:company.delete"
)

type Company struct {
	ID          int       `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description *string   `db:"description" json:"description"`
	ParentID    *int      `db:"parent_id" json:"parentId"`
	StoragePath string    `db:"storage_path" json:"storagePath"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `db:"updated_at" json:"updatedAt"`
}

var slugRegexp = regexp.MustCompile(`[^a-z0-9]+`)

func slugify(name string) string {
	slug := strings.ToLower(strings.TrimSpace(name))
	slug = slugRegexp.ReplaceAllString(slug, "-")
	return strings.Trim(slug, "-")
}

func (c *Company) ApplyPayload(payload CompanyBody) error {
	c.Name = payload.Name
	c.ParentID = payload.ParentID
	return c.Save()
}

func (c *Company) Save() error {
	service := CompanyService{}
	c.UpdatedAt = time.Now()
	if c.ID > 0 {
		return service.Update(c)
	}
	c.CreatedAt = time.Now()

	slug := slugify(c.Name)
	if slug == "" {
		slug = fmt.Sprintf("company-%d", time.Now().UnixMilli())
	}
	storagePath, err := backends.Storage().CreateCompanyFolder(slug)
	if err != nil {
		return fmt.Errorf("cannot create company folder: %w", err)
	}
	c.StoragePath = storagePath

	return service.Insert(c)
}

func (c *Company) Delete() error {
	if c.StoragePath != "" {
		_ = backends.Storage().DeleteCompanyFolder(c.StoragePath)
	}
	service := CompanyService{}
	return service.Delete(c.ID)
}

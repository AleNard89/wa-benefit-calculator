package auth

import (
	"orbita-api/orgs"
	"testing"
)

func makeUserWithCompanies() *User {
	return &User{
		ID:    1,
		Email: "test@example.com",
		CompanyRoles: []CompanyRoles{
			{
				Company: orgs.Company{ID: 10, Name: "Company A"},
				Roles: []Role{
					{ID: 1, Name: "Admin", Permissions: []Permission{
						{ID: 1, Code: "processes:process.read"},
						{ID: 2, Code: "processes:process.create"},
					}},
				},
			},
			{
				Company: orgs.Company{ID: 20, Name: "Company B"},
				Roles: []Role{
					{ID: 2, Name: "Reader", Permissions: []Permission{
						{ID: 1, Code: "processes:process.read"},
					}},
				},
			},
		},
	}
}

func TestUser_IsAnonymous(t *testing.T) {
	anon := &User{}
	if !anon.IsAnonymous() {
		t.Error("zero-ID user should be anonymous")
	}

	u := &User{ID: 1}
	if u.IsAnonymous() {
		t.Error("user with ID should not be anonymous")
	}
}

func TestUser_IsSuperUser(t *testing.T) {
	u := &User{IsSuperuser: true}
	if !u.IsSuperUser() {
		t.Error("expected IsSuperUser to be true")
	}
	u2 := &User{IsSuperuser: false}
	if u2.IsSuperUser() {
		t.Error("expected IsSuperUser to be false")
	}
}

func TestUser_BelongsToCompany(t *testing.T) {
	u := makeUserWithCompanies()

	if !u.BelongsToCompany(10) {
		t.Error("user should belong to company 10")
	}
	if !u.BelongsToCompany(20) {
		t.Error("user should belong to company 20")
	}
	if u.BelongsToCompany(99) {
		t.Error("user should not belong to company 99")
	}
}

func TestUser_AllCompanyPermissionsCodes(t *testing.T) {
	u := makeUserWithCompanies()

	permsA := u.AllCompanyPermissionsCodes(10)
	if len(permsA) != 2 {
		t.Errorf("company 10 should have 2 permissions, got %d", len(permsA))
	}

	permsB := u.AllCompanyPermissionsCodes(20)
	if len(permsB) != 1 {
		t.Errorf("company 20 should have 1 permission, got %d", len(permsB))
	}

	permsNone := u.AllCompanyPermissionsCodes(99)
	if len(permsNone) != 0 {
		t.Errorf("unknown company should have 0 permissions, got %d", len(permsNone))
	}
}

func TestUser_AllPermissionsCodes(t *testing.T) {
	u := makeUserWithCompanies()
	all := u.AllPermissionsCodes()

	// Both companies share "processes:process.read", company 10 also has "processes:process.create"
	if len(all) != 2 {
		t.Errorf("expected 2 unique permissions, got %d: %v", len(all), all)
	}
}

func TestUser_BelongsToCompany_EmptyRoles(t *testing.T) {
	u := &User{ID: 1, CompanyRoles: []CompanyRoles{}}
	if u.BelongsToCompany(1) {
		t.Error("user with no company roles should not belong to any company")
	}
}

func TestUser_GetID(t *testing.T) {
	u := &User{ID: 42}
	if u.GetID() != 42 {
		t.Errorf("expected ID 42, got %d", u.GetID())
	}
}

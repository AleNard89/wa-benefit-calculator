package shared

type User interface {
	GetID() int
	IsSuperUser() bool
	IsAnonymous() bool
	BelongsToCompany(companyID int) bool
	AllCompanyPermissionsCodes(companyID int) []string
	AllPermissionsCodes() []string
}

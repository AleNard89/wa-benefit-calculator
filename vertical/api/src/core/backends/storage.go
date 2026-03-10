package backends

// StorageBackend defines the interface for company document storage.
// Default: local filesystem. Can be swapped for SharePoint, S3, etc.
type StorageBackend interface {
	CreateCompanyFolder(companySlug string) (path string, err error)
	DeleteCompanyFolder(path string) error
	CompanyFolderExists(path string) bool
}

var activeBackend StorageBackend = &LocalStorage{}

func Storage() StorageBackend {
	return activeBackend
}

func SetStorageBackend(b StorageBackend) {
	activeBackend = b
}

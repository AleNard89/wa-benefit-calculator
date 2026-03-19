package orgs

import (
	"orbita-api/core/backends"

	"go.uber.org/zap"
)

// EnsureCompanyFolders creates storage folders for any company that has a
// storage_path set in the DB but no folder on disk yet (e.g. seeded companies).
func EnsureCompanyFolders() {
	service := CompanyService{}
	companies, err := service.GetAll()
	if err != nil {
		zap.S().Warnw("EnsureCompanyFolders: cannot list companies", "error", err)
		return
	}
	storage := backends.Storage()
	for _, c := range companies {
		if c.StoragePath == "" {
			continue
		}
		if storage.CompanyFolderExists(c.StoragePath) {
			continue
		}
		if ls, ok := storage.(*backends.LocalStorage); ok {
			if err := ls.EnsureFolder(c.StoragePath); err != nil {
				zap.S().Warnw("EnsureCompanyFolders: cannot create folder", "company", c.Name, "path", c.StoragePath, "error", err)
			} else {
				zap.S().Infow("EnsureCompanyFolders: created missing folder", "company", c.Name, "path", c.StoragePath)
			}
		}
	}
}

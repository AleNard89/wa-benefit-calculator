package backends

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const defaultBasePath = "/data/companies"

type LocalStorage struct {
	BasePath string
}

func (s *LocalStorage) basePath() string {
	if s.BasePath != "" {
		return s.BasePath
	}
	envPath := os.Getenv("COMPANY_STORAGE_PATH")
	if envPath != "" {
		return envPath
	}
	return defaultBasePath
}

func (s *LocalStorage) CreateCompanyFolder(companySlug string) (string, error) {
	if strings.Contains(companySlug, "..") || strings.ContainsAny(companySlug, "/\\") {
		return "", fmt.Errorf("invalid company slug: %q", companySlug)
	}
	folderPath := filepath.Join(s.basePath(), companySlug)
	if err := os.MkdirAll(folderPath, 0750); err != nil {
		return "", err
	}
	return folderPath, nil
}

func (s *LocalStorage) DeleteCompanyFolder(path string) error {
	if path == "" {
		return nil
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	base, err := filepath.Abs(s.basePath())
	if err != nil {
		return err
	}
	if !strings.HasPrefix(absPath, base+string(filepath.Separator)) {
		return fmt.Errorf("path %q is outside base directory", absPath)
	}
	return os.RemoveAll(absPath)
}

func (s *LocalStorage) CompanyFolderExists(path string) bool {
	if path == "" {
		return false
	}
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// EnsureFolder creates the folder if it doesn't exist yet (idempotent).
func (s *LocalStorage) EnsureFolder(path string) error {
	if path == "" {
		return nil
	}
	if s.CompanyFolderExists(path) {
		return nil
	}
	return os.MkdirAll(path, 0750)
}

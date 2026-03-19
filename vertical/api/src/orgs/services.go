package orgs

import (
	"context"
	"orbita-api/core/utils"
	"orbita-api/db"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
)

type CompanyService struct{}

func (s *CompanyService) GetAll() ([]Company, error) {
	db := db.DB()
	companies := make([]Company, 0)
	stm := db.Builder.Select("*").From(CompaniesTable)
	sql, args, err := stm.ToSql()
	if err != nil {
		return nil, err
	}
	err = pgxscan.Select(context.Background(), db.C, &companies, sql, args...)
	if err != nil {
		return nil, err
	}
	return companies, nil
}

func (s *CompanyService) GetByUserID(userID int) ([]Company, error) {
	db := db.DB()
	companies := make([]Company, 0)
	stm := db.Builder.Select("c.*").From(CompaniesTable + " c").
		Join("auth_users_companies uc ON uc.company_id = c.id").
		Where(sq.Eq{"uc.user_id": userID})
	sql, args, err := stm.ToSql()
	if err != nil {
		return nil, err
	}
	err = pgxscan.Select(context.Background(), db.C, &companies, sql, args...)
	if err != nil {
		return nil, err
	}
	return companies, nil
}

func (s *CompanyService) GetByID(id int) (*Company, error) {
	db := db.DB()
	company := &Company{}
	stm := db.Builder.Select("*").From(CompaniesTable).Where(sq.Eq{"id": id})
	sql, args, err := stm.ToSql()
	if err != nil {
		return nil, err
	}
	err = pgxscan.Get(context.Background(), db.C, company, sql, args...)
	return utils.RetZeroOnError(company, err)
}

func (s *CompanyService) Insert(company *Company) error {
	db := db.DB()
	stm := db.Builder.
		Insert(CompaniesTable).
		Columns("name", "parent_id", "storage_path", "created_at", "updated_at").
		Values(company.Name, company.ParentID, company.StoragePath, company.CreatedAt, company.UpdatedAt).
		Suffix("RETURNING id")
	sql, args, err := stm.ToSql()
	if err != nil {
		return err
	}
	return db.C.QueryRow(context.Background(), sql, args...).Scan(&company.ID)
}

func (s *CompanyService) Update(company *Company) error {
	db := db.DB()
	stm := db.Builder.
		Update(CompaniesTable).
		Set("name", company.Name).
		Set("parent_id", company.ParentID).
		Set("updated_at", company.UpdatedAt).
		Where(sq.Eq{"id": company.ID})
	sql, args, err := stm.ToSql()
	if err != nil {
		return err
	}
	_, err = db.C.Exec(context.Background(), sql, args...)
	return err
}

func (s *CompanyService) Delete(id int) error {
	db := db.DB()
	stm := db.Builder.Delete(CompaniesTable).Where(sq.Eq{"id": id})
	sql, args, err := stm.ToSql()
	if err != nil {
		return err
	}
	_, err = db.C.Exec(context.Background(), sql, args...)
	return err
}

package orgs

import (
	"orbita-api/core/utils"
	"orbita-api/db"
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
)

type AreaService struct{}

func (s *AreaService) GetByCompany(companyID int) ([]Area, error) {
	database := db.DB()
	areas := make([]Area, 0)
	stm := database.Builder.Select("*").From(AreasTable).Where(sq.Eq{"company_id": companyID}).OrderBy("name ASC")
	sql, args, err := stm.ToSql()
	if err != nil {
		return nil, err
	}
	err = pgxscan.Select(context.Background(), database.C, &areas, sql, args...)
	if err != nil {
		return nil, err
	}
	return areas, nil
}

func (s *AreaService) GetByID(id int) (*Area, error) {
	database := db.DB()
	area := &Area{}
	stm := database.Builder.Select("*").From(AreasTable).Where(sq.Eq{"id": id})
	sql, args, err := stm.ToSql()
	if err != nil {
		return nil, err
	}
	err = pgxscan.Get(context.Background(), database.C, area, sql, args...)
	return utils.RetZeroOnError(area, err)
}

func (s *AreaService) Insert(area *Area) error {
	database := db.DB()
	area.CreatedAt = time.Now()
	area.UpdatedAt = time.Now()
	stm := database.Builder.
		Insert(AreasTable).
		Columns("company_id", "name", "created_at", "updated_at").
		Values(area.CompanyID, area.Name, area.CreatedAt, area.UpdatedAt).
		Suffix("RETURNING id")
	sql, args, err := stm.ToSql()
	if err != nil {
		return err
	}
	return database.C.QueryRow(context.Background(), sql, args...).Scan(&area.ID)
}

func (s *AreaService) Update(area *Area) error {
	database := db.DB()
	area.UpdatedAt = time.Now()
	stm := database.Builder.
		Update(AreasTable).
		Set("name", area.Name).
		Set("updated_at", area.UpdatedAt).
		Where(sq.Eq{"id": area.ID})
	sql, args, err := stm.ToSql()
	if err != nil {
		return err
	}
	_, err = database.C.Exec(context.Background(), sql, args...)
	return err
}

func (s *AreaService) Delete(id int) error {
	database := db.DB()
	stm := database.Builder.Delete(AreasTable).Where(sq.Eq{"id": id})
	sql, args, err := stm.ToSql()
	if err != nil {
		return err
	}
	_, err = database.C.Exec(context.Background(), sql, args...)
	return err
}

func (s *AreaService) GetUserAreas(userID int, companyID int) ([]Area, error) {
	database := db.DB()
	areas := make([]Area, 0)
	stm := database.Builder.
		Select("a.*").
		From(AreasTable + " a").
		Join(UsersAreasTable + " ua ON ua.area_id = a.id").
		Where(sq.Eq{"ua.user_id": userID, "a.company_id": companyID}).
		OrderBy("a.name ASC")
	sql, args, err := stm.ToSql()
	if err != nil {
		return nil, err
	}
	err = pgxscan.Select(context.Background(), database.C, &areas, sql, args...)
	if err != nil {
		return nil, err
	}
	return areas, nil
}

func (s *AreaService) SetUserAreas(userID int, areaIDs []int) error {
	database := db.DB()
	ctx := context.Background()

	tx, err := database.C.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, "DELETE FROM "+UsersAreasTable+" WHERE user_id = $1", userID)
	if err != nil {
		return err
	}

	for _, areaID := range areaIDs {
		_, err = tx.Exec(ctx,
			"INSERT INTO "+UsersAreasTable+" (user_id, area_id) VALUES ($1, $2)",
			userID, areaID)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

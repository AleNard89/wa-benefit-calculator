package processes

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"benefit-calculator-api/db"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

type ProcessService struct{}

func (s *ProcessService) setTenant(ctx context.Context, tx pgx.Tx, companyID int) error {
	_, err := tx.Exec(ctx, "SELECT set_config('app.current_tenant', $1, true)", strconv.Itoa(companyID))
	return err
}

func (s *ProcessService) GetAll(companyID int, params ProcessListParams) ([]Process, int, error) {
	database := db.DB()
	ctx := context.Background()
	processes := make([]Process, 0)

	tx, err := database.C.Begin(ctx)
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback(ctx)

	if err := s.setTenant(ctx, tx, companyID); err != nil {
		return nil, 0, err
	}

	companyFilter := sq.Eq{"company_id": companyID}
	var deletedFilter sq.Sqlizer
	if params.Deleted {
		deletedFilter = sq.NotEq{"deleted_at": nil}
	} else {
		deletedFilter = sq.Eq{"deleted_at": nil}
	}

	// Count query
	countStm := database.Builder.Select("COUNT(*)").From(ProcessesTable).Where(companyFilter).Where(deletedFilter)
	if len(params.AreaIDs) > 0 {
		countStm = countStm.Where(sq.Eq{"area_id": params.AreaIDs})
	}
	if params.Status != "" {
		countStm = countStm.Where(sq.Eq{"status": params.Status})
	}
	if params.Search != "" {
		searchPattern := fmt.Sprintf("%%%s%%", params.Search)
		countStm = countStm.Where(
			sq.Or{
				sq.ILike{"process_name": searchPattern},
				sq.Expr("data->>'proposer' ILIKE ?", searchPattern),
				sq.Expr("data->>'area' ILIKE ?", searchPattern),
			},
		)
	}

	countSQL, countArgs, err := countStm.ToSql()
	if err != nil {
		return nil, 0, err
	}
	var total int
	if err := tx.QueryRow(ctx, countSQL, countArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Data query
	stm := database.Builder.Select("*").From(ProcessesTable).Where(companyFilter).Where(deletedFilter)
	if len(params.AreaIDs) > 0 {
		stm = stm.Where(sq.Eq{"area_id": params.AreaIDs})
	}
	if params.Status != "" {
		stm = stm.Where(sq.Eq{"status": params.Status})
	}
	if params.Search != "" {
		searchPattern := fmt.Sprintf("%%%s%%", params.Search)
		stm = stm.Where(
			sq.Or{
				sq.ILike{"process_name": searchPattern},
				sq.Expr("data->>'proposer' ILIKE ?", searchPattern),
				sq.Expr("data->>'area' ILIKE ?", searchPattern),
			},
		)
	}

	allowedSorts := map[string]bool{
		"created_at": true, "updated_at": true, "process_name": true, "status": true,
	}
	sortBy := "created_at"
	if allowedSorts[params.SortBy] {
		sortBy = params.SortBy
	}
	order := "DESC"
	if params.Order == "asc" {
		order = "ASC"
	}
	stm = stm.OrderBy(fmt.Sprintf("%s %s", sortBy, order))

	if params.Limit > 0 {
		stm = stm.Limit(uint64(params.Limit))
	}
	if params.Page > 1 && params.Limit > 0 {
		stm = stm.Offset(uint64((params.Page - 1) * params.Limit))
	}

	sql, args, err := stm.ToSql()
	if err != nil {
		return nil, 0, err
	}

	rows, err := tx.Query(ctx, sql, args...)
	if err != nil {
		return nil, 0, err
	}
	if err := pgxscan.ScanAll(&processes, rows); err != nil {
		return nil, 0, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, 0, err
	}

	return processes, total, nil
}

func (s *ProcessService) GetByID(id, companyID int) (*Process, error) {
	database := db.DB()
	ctx := context.Background()
	process := &Process{}

	tx, err := database.C.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	if err := s.setTenant(ctx, tx, companyID); err != nil {
		return nil, err
	}

	stm := database.Builder.Select("*").From(ProcessesTable).Where(sq.Eq{"id": id, "company_id": companyID, "deleted_at": nil})
	sql, args, err := stm.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	if err := pgxscan.ScanOne(process, rows); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return process, nil
}

func (s *ProcessService) Insert(p *Process) error {
	database := db.DB()
	ctx := context.Background()

	tx, err := database.C.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := s.setTenant(ctx, tx, p.CompanyID); err != nil {
		return err
	}

	stm := database.Builder.
		Insert(ProcessesTable).
		Columns("company_id", "area_id", "process_name", "status", "data", "results", "created_by", "created_at", "updated_at").
		Values(p.CompanyID, p.AreaID, p.ProcessName, p.Status, p.Data, p.Results, p.CreatedBy, p.CreatedAt, p.UpdatedAt).
		Suffix("RETURNING id")

	sql, args, err := stm.ToSql()
	if err != nil {
		return err
	}

	if err := tx.QueryRow(ctx, sql, args...).Scan(&p.ID); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *ProcessService) Update(p *Process) error {
	database := db.DB()
	ctx := context.Background()

	tx, err := database.C.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := s.setTenant(ctx, tx, p.CompanyID); err != nil {
		return err
	}

	stm := database.Builder.
		Update(ProcessesTable).
		Set("area_id", p.AreaID).
		Set("process_name", p.ProcessName).
		Set("data", p.Data).
		Set("results", p.Results).
		Set("document_path", p.DocumentPath).
		Set("document_name", p.DocumentName).
		Set("updated_at", p.UpdatedAt).
		Where(sq.Eq{"id": p.ID, "company_id": p.CompanyID})

	sql, args, err := stm.ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *ProcessService) UpdateStatus(id, companyID int, status string) error {
	database := db.DB()
	ctx := context.Background()

	tx, err := database.C.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := s.setTenant(ctx, tx, companyID); err != nil {
		return err
	}

	stm := database.Builder.
		Update(ProcessesTable).
		Set("status", status).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": id, "company_id": companyID})

	sql, args, err := stm.ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *ProcessService) Delete(id, companyID int) error {
	database := db.DB()
	ctx := context.Background()

	tx, err := database.C.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := s.setTenant(ctx, tx, companyID); err != nil {
		return err
	}

	stm := database.Builder.
		Update(ProcessesTable).
		Set("deleted_at", time.Now()).
		Where(sq.Eq{"id": id, "company_id": companyID})
	sql, args, err := stm.ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *ProcessService) Restore(id, companyID int) error {
	database := db.DB()
	ctx := context.Background()

	tx, err := database.C.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := s.setTenant(ctx, tx, companyID); err != nil {
		return err
	}

	stm := database.Builder.
		Update(ProcessesTable).
		Set("deleted_at", nil).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": id, "company_id": companyID}).
		Where(sq.NotEq{"deleted_at": nil})
	sql, args, err := stm.ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *ProcessService) GetStats(companyID int, areaIDs []int) (*ProcessStats, error) {
	database := db.DB()
	ctx := context.Background()

	tx, err := database.C.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	if err := s.setTenant(ctx, tx, companyID); err != nil {
		return nil, err
	}

	stm := database.Builder.
		Select("COUNT(*)").
		Column(sq.Expr("COUNT(*) FILTER (WHERE status = ?)", StatusToValuate)).
		Column(sq.Expr("COUNT(*) FILTER (WHERE status = ?)", StatusAnalysis)).
		Column(sq.Expr("COUNT(*) FILTER (WHERE status = ?)", StatusOngoing)).
		Column(sq.Expr("COUNT(*) FILTER (WHERE status = ?)", StatusProduction)).
		From(ProcessesTable).
		Where(sq.Eq{"company_id": companyID, "deleted_at": nil})

	if len(areaIDs) > 0 {
		stm = stm.Where(sq.Eq{"area_id": areaIDs})
	}

	sql, args, err := stm.ToSql()
	if err != nil {
		return nil, err
	}

	stats := &ProcessStats{}
	err = tx.QueryRow(ctx, sql, args...).
		Scan(&stats.Total, &stats.ToValuate, &stats.Analysis, &stats.Ongoing, &stats.Production)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return stats, nil
}

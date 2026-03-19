package orchestrator

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"orbita-api/db"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

// GetDistinctProcessCodes returns unique process codes extracted from the first
// underscore-separated token of process_name (jobs), release_name (schedules),
// and queue_name (queue items). Example: "Proc01_Something_PRF" -> "Proc01".
func (s *Service) GetDistinctProcessCodes(companyID int) ([]string, error) {
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

	query := `SELECT DISTINCT code FROM (
		SELECT split_part(process_name, '_', 1) AS code FROM orchestrator_job_executions WHERE company_id = $1 AND process_name IS NOT NULL AND process_name != ''
		UNION
		SELECT split_part(release_name, '_', 1) AS code FROM orchestrator_schedules WHERE company_id = $1 AND release_name IS NOT NULL AND release_name != ''
		UNION
		SELECT split_part(queue_name, '_', 1) AS code FROM orchestrator_queue_items WHERE company_id = $1 AND queue_name IS NOT NULL AND queue_name != ''
	) sub WHERE code != '' ORDER BY code`

	rows, err := tx.Query(ctx, query, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var codes []string
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			return nil, err
		}
		codes = append(codes, code)
	}

	return codes, tx.Commit(ctx)
}

func (s *Service) GetDistinctBotNames(companyID int) ([]string, error) {
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

	query := `SELECT DISTINCT name FROM (
		SELECT process_name AS name FROM orchestrator_job_executions WHERE company_id = $1 AND process_name IS NOT NULL AND process_name != ''
		UNION
		SELECT release_name AS name FROM orchestrator_schedules WHERE company_id = $1 AND release_name IS NOT NULL AND release_name != ''
	) sub ORDER BY name`

	rows, err := tx.Query(ctx, query, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		names = append(names, name)
	}

	return names, tx.Commit(ctx)
}

// processCodeFilter builds an OR clause that matches rows where
// split_part(column, '_', 1) equals any of the given codes.
func processCodeFilter(column string, codes []string) sq.Sqlizer {
	if len(codes) == 1 {
		return sq.Expr("split_part("+column+", '_', 1) = ?", codes[0])
	}
	or := sq.Or{}
	for _, code := range codes {
		or = append(or, sq.Expr("split_part("+column+", '_', 1) = ?", code))
	}
	return or
}

type Service struct{}

func (s *Service) setTenant(ctx context.Context, tx pgx.Tx, companyID int) error {
	_, err := tx.Exec(ctx, "SELECT set_config('app.current_tenant', $1, true)", strconv.Itoa(companyID))
	return err
}

// --- Connectors ---

func (s *Service) ListConnectors(companyID int) ([]Connector, error) {
	database := db.DB()
	ctx := context.Background()
	connectors := make([]Connector, 0)

	tx, err := database.C.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	if err := s.setTenant(ctx, tx, companyID); err != nil {
		return nil, err
	}

	sql, args, err := database.Builder.
		Select("*").From(ConnectorsTable).
		Where(sq.Eq{"company_id": companyID}).
		OrderBy("created_at ASC").
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	if err := pgxscan.ScanAll(&connectors, rows); err != nil {
		return nil, err
	}

	return connectors, tx.Commit(ctx)
}

func (s *Service) GetConnector(id, companyID int) (*Connector, error) {
	database := db.DB()
	ctx := context.Background()
	connector := &Connector{}

	tx, err := database.C.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	if err := s.setTenant(ctx, tx, companyID); err != nil {
		return nil, err
	}

	sql, args, err := database.Builder.
		Select("*").From(ConnectorsTable).
		Where(sq.Eq{"id": id, "company_id": companyID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	if err := pgxscan.ScanOne(connector, rows); err != nil {
		return nil, err
	}

	return connector, tx.Commit(ctx)
}

func (s *Service) CreateConnector(companyID int, body ConnectorBody, encryptFn func(string) (string, error)) (*Connector, error) {
	database := db.DB()
	ctx := context.Background()

	encrypted, err := encryptFn(body.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("encryption failed: %w", err)
	}

	cfg := UiPathConfig{
		OrganizationName:    body.OrganizationName,
		TenantName:          body.TenantName,
		PersonalAccessToken: encrypted,
		FolderID:            body.FolderID,
		FolderName:          body.FolderName,
		Folders:             body.Folders,
	}
	configJSON, err := json.Marshal(cfg)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	connector := &Connector{
		CompanyID: companyID,
		Name:      body.Name,
		Type:      body.Type,
		Config:    configJSON,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	tx, err := database.C.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	if err := s.setTenant(ctx, tx, companyID); err != nil {
		return nil, err
	}

	sql, args, err := database.Builder.
		Insert(ConnectorsTable).
		Columns("company_id", "name", "type", "config", "is_active", "created_at", "updated_at").
		Values(companyID, connector.Name, connector.Type, connector.Config, connector.IsActive, connector.CreatedAt, connector.UpdatedAt).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return nil, err
	}

	if err := tx.QueryRow(ctx, sql, args...).Scan(&connector.ID); err != nil {
		return nil, err
	}

	return connector, tx.Commit(ctx)
}

func (s *Service) UpdateConnector(id, companyID int, body ConnectorBody, encryptFn func(string) (string, error)) (*Connector, error) {
	database := db.DB()
	ctx := context.Background()

	existing, err := s.GetConnector(id, companyID)
	if err != nil {
		return nil, err
	}

	existingCfg, _ := existing.GetUiPathConfig()
	token := existingCfg.PersonalAccessToken
	if body.AccessToken != "" {
		token, err = encryptFn(body.AccessToken)
		if err != nil {
			return nil, fmt.Errorf("encryption failed: %w", err)
		}
	}

	cfg := UiPathConfig{
		OrganizationName:    body.OrganizationName,
		TenantName:          body.TenantName,
		PersonalAccessToken: token,
		FolderID:            body.FolderID,
		FolderName:          body.FolderName,
		Folders:             body.Folders,
	}
	configJSON, err := json.Marshal(cfg)
	if err != nil {
		return nil, err
	}

	isActive := existing.IsActive
	if body.IsActive != nil {
		isActive = *body.IsActive
	}

	tx, err := database.C.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	if err := s.setTenant(ctx, tx, companyID); err != nil {
		return nil, err
	}

	sql, args, err := database.Builder.
		Update(ConnectorsTable).
		Set("name", body.Name).
		Set("type", body.Type).
		Set("config", configJSON).
		Set("is_active", isActive).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return s.GetConnector(id, companyID)
}

// GetActiveConnectors returns active connectors (no RLS, used by sync).
func (s *Service) GetActiveConnectors(connectorType string) ([]Connector, error) {
	database := db.DB()
	ctx := context.Background()
	connectors := make([]Connector, 0)

	sql, args, err := database.Builder.
		Select("*").From(ConnectorsTable).
		Where(sq.Eq{"type": connectorType, "is_active": true}).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := database.C.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	if err := pgxscan.ScanAll(&connectors, rows); err != nil {
		return nil, err
	}

	return connectors, nil
}

// --- Job Executions ---

func (s *Service) ListJobs(companyID int, params JobListParams) ([]JobExecution, int, error) {
	database := db.DB()
	ctx := context.Background()
	jobs := make([]JobExecution, 0)

	tx, err := database.C.Begin(ctx)
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback(ctx)

	if err := s.setTenant(ctx, tx, companyID); err != nil {
		return nil, 0, err
	}

	where := sq.And{sq.Eq{"company_id": companyID}}
	if params.State != "" {
		where = append(where, sq.Eq{"state": params.State})
	}
	if params.ConnectorID > 0 {
		where = append(where, sq.Eq{"connector_id": params.ConnectorID})
	}
	if codes := parseCSV(params.ProcessNames); len(codes) > 0 {
		where = append(where, processCodeFilter("process_name", codes))
	}

	// Count
	countStm := database.Builder.Select("COUNT(*)").From(JobExecutionsTable)
	if len(where) > 0 {
		countStm = countStm.Where(where)
	}
	countSQL, countArgs, err := countStm.ToSql()
	if err != nil {
		return nil, 0, err
	}
	var total int
	if err := tx.QueryRow(ctx, countSQL, countArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Data
	stm := database.Builder.Select("*").From(JobExecutionsTable).OrderBy("start_time DESC NULLS LAST")
	if len(where) > 0 {
		stm = stm.Where(where)
	}
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
	if err := pgxscan.ScanAll(&jobs, rows); err != nil {
		return nil, 0, err
	}

	return jobs, total, tx.Commit(ctx)
}

func (s *Service) GetJob(id, companyID int) (*JobExecution, error) {
	database := db.DB()
	ctx := context.Background()
	job := &JobExecution{}

	tx, err := database.C.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	if err := s.setTenant(ctx, tx, companyID); err != nil {
		return nil, err
	}

	sql, args, err := database.Builder.
		Select("*").From(JobExecutionsTable).
		Where(sq.Eq{"id": id, "company_id": companyID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	if err := pgxscan.ScanOne(job, rows); err != nil {
		return nil, err
	}

	return job, tx.Commit(ctx)
}

// --- Queue Items ---

func (s *Service) ListQueueItems(companyID int, params QueueItemListParams) ([]QueueItem, int, error) {
	database := db.DB()
	ctx := context.Background()
	items := make([]QueueItem, 0)

	tx, err := database.C.Begin(ctx)
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback(ctx)

	if err := s.setTenant(ctx, tx, companyID); err != nil {
		return nil, 0, err
	}

	where := sq.And{sq.Eq{"company_id": companyID}}
	if params.Status != "" {
		where = append(where, sq.Eq{"status": params.Status})
	}
	if params.ConnectorID > 0 {
		where = append(where, sq.Eq{"connector_id": params.ConnectorID})
	}
	if params.QueueName != "" {
		where = append(where, sq.Eq{"queue_name": params.QueueName})
	}
	if codes := parseCSV(params.ProcessNames); len(codes) > 0 {
		where = append(where, processCodeFilter("queue_name", codes))
	}

	// Count
	countStm := database.Builder.Select("COUNT(*)").From(QueueItemsTable)
	if len(where) > 0 {
		countStm = countStm.Where(where)
	}
	countSQL, countArgs, err := countStm.ToSql()
	if err != nil {
		return nil, 0, err
	}
	var total int
	if err := tx.QueryRow(ctx, countSQL, countArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Data
	stm := database.Builder.Select("*").From(QueueItemsTable).OrderBy("created_at DESC")
	if len(where) > 0 {
		stm = stm.Where(where)
	}
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
	if err := pgxscan.ScanAll(&items, rows); err != nil {
		return nil, 0, err
	}

	return items, total, tx.Commit(ctx)
}

// --- Queue Definitions ---

func (s *Service) ListQueueDefinitions(companyID int, connectorID int) ([]QueueDefinition, error) {
	database := db.DB()
	ctx := context.Background()
	defs := make([]QueueDefinition, 0)

	tx, err := database.C.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	if err := s.setTenant(ctx, tx, companyID); err != nil {
		return nil, err
	}

	stm := database.Builder.Select("*").From(QueueDefinitionsTable).Where(sq.Eq{"company_id": companyID}).OrderBy("name ASC")
	if connectorID > 0 {
		stm = stm.Where(sq.Eq{"connector_id": connectorID})
	}

	sql, args, err := stm.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	if err := pgxscan.ScanAll(&defs, rows); err != nil {
		return nil, err
	}

	return defs, tx.Commit(ctx)
}

// --- Schedules ---

func (s *Service) ListSchedules(companyID int, params ScheduleListParams) ([]ProcessSchedule, int, error) {
	database := db.DB()
	ctx := context.Background()
	schedules := make([]ProcessSchedule, 0)

	tx, err := database.C.Begin(ctx)
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback(ctx)

	if err := s.setTenant(ctx, tx, companyID); err != nil {
		return nil, 0, err
	}

	where := sq.And{sq.Eq{"company_id": companyID}}
	if params.Enabled != nil {
		where = append(where, sq.Eq{"enabled": *params.Enabled})
	}
	if params.ConnectorID > 0 {
		where = append(where, sq.Eq{"connector_id": params.ConnectorID})
	}
	if codes := parseCSV(params.ProcessNames); len(codes) > 0 {
		where = append(where, processCodeFilter("release_name", codes))
	}

	countStm := database.Builder.Select("COUNT(*)").From(SchedulesTable)
	if len(where) > 0 {
		countStm = countStm.Where(where)
	}
	countSQL, countArgs, err := countStm.ToSql()
	if err != nil {
		return nil, 0, err
	}
	var total int
	if err := tx.QueryRow(ctx, countSQL, countArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}

	stm := database.Builder.Select("*").From(SchedulesTable).OrderBy("name ASC")
	if len(where) > 0 {
		stm = stm.Where(where)
	}
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
	if err := pgxscan.ScanAll(&schedules, rows); err != nil {
		return nil, 0, err
	}

	return schedules, total, tx.Commit(ctx)
}

// --- Upsert methods for sync ---

func (s *Service) UpsertJob(ctx context.Context, tx pgx.Tx, job *JobExecution) error {
	database := db.DB()
	sql, args, err := database.Builder.
		Insert(JobExecutionsTable).
		Columns("company_id", "connector_id", "external_job_key", "external_job_id", "process_name", "state",
			"source_type", "source", "start_time", "end_time", "host_machine", "folder_name", "info", "details",
			"created_at", "updated_at").
		Values(job.CompanyID, job.ConnectorID, job.ExternalJobKey, job.ExternalJobID, job.ProcessName, job.State,
			job.SourceType, job.Source, job.StartTime, job.EndTime, job.HostMachine, job.FolderName, job.Info, job.Details,
			job.CreatedAt, job.UpdatedAt).
		Suffix(`ON CONFLICT (connector_id, external_job_key) DO UPDATE SET
			state = EXCLUDED.state, end_time = EXCLUDED.end_time, info = EXCLUDED.info,
			details = EXCLUDED.details, updated_at = EXCLUDED.updated_at`).
		ToSql()
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, sql, args...)
	return err
}

func (s *Service) UpsertQueueDefinition(ctx context.Context, tx pgx.Tx, def *QueueDefinition) error {
	database := db.DB()
	sql, args, err := database.Builder.
		Insert(QueueDefinitionsTable).
		Columns("company_id", "connector_id", "external_definition_id", "name", "max_retries", "folder_name",
			"created_at", "updated_at").
		Values(def.CompanyID, def.ConnectorID, def.ExternalDefinitionID, def.Name, def.MaxRetries, def.FolderName,
			def.CreatedAt, def.UpdatedAt).
		Suffix(`ON CONFLICT (connector_id, external_definition_id) DO UPDATE SET
			name = EXCLUDED.name, max_retries = EXCLUDED.max_retries, updated_at = EXCLUDED.updated_at`).
		ToSql()
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, sql, args...)
	return err
}

func (s *Service) UpsertQueueItem(ctx context.Context, tx pgx.Tx, item *QueueItem) error {
	database := db.DB()
	sql, args, err := database.Builder.
		Insert(QueueItemsTable).
		Columns("company_id", "connector_id", "external_item_key", "external_item_id", "queue_definition_id",
			"queue_name", "status", "priority", "reference", "processing_exception_type", "error_message",
			"start_processing", "end_processing", "retry_number", "folder_name", "specific_content",
			"created_at", "updated_at").
		Values(item.CompanyID, item.ConnectorID, item.ExternalItemKey, item.ExternalItemID, item.QueueDefinitionID,
			item.QueueName, item.Status, item.Priority, item.Reference, item.ProcessingExceptionType, item.ErrorMessage,
			item.StartProcessing, item.EndProcessing, item.RetryNumber, item.FolderName, item.SpecificContent,
			item.CreatedAt, item.UpdatedAt).
		Suffix(`ON CONFLICT (connector_id, external_item_key) DO UPDATE SET
			status = EXCLUDED.status, end_processing = EXCLUDED.end_processing, error_message = EXCLUDED.error_message,
			processing_exception_type = EXCLUDED.processing_exception_type, retry_number = EXCLUDED.retry_number,
			specific_content = EXCLUDED.specific_content, updated_at = EXCLUDED.updated_at`).
		ToSql()
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, sql, args...)
	return err
}

// --- Process-Queue Map ---

func (s *Service) ListProcessQueueMaps(companyID int) ([]ProcessQueueMap, error) {
	database := db.DB()
	ctx := context.Background()
	maps := make([]ProcessQueueMap, 0)

	tx, err := database.C.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	if err := s.setTenant(ctx, tx, companyID); err != nil {
		return nil, err
	}

	stm := database.Builder.Select("*").From(ProcessQueueMapTable).
		Where(sq.Eq{"company_id": companyID}).OrderBy("process_name ASC, queue_name ASC")

	sql, args, err := stm.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	if err := pgxscan.ScanAll(&maps, rows); err != nil {
		return nil, err
	}

	return maps, tx.Commit(ctx)
}

func (s *Service) CreateProcessQueueMap(companyID int, m *ProcessQueueMap) error {
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

	sql, args, err := database.Builder.
		Insert(ProcessQueueMapTable).
		Columns("company_id", "connector_id", "process_name", "queue_name", "auto_detected", "created_at", "updated_at").
		Values(m.CompanyID, m.ConnectorID, m.ProcessName, m.QueueName, m.AutoDetected, m.CreatedAt, m.UpdatedAt).
		Suffix("ON CONFLICT (company_id, connector_id, process_name, queue_name) DO NOTHING").
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return err
	}

	err = tx.QueryRow(ctx, sql, args...).Scan(&m.ID)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *Service) DeleteProcessQueueMap(id, companyID int) error {
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

	sql, args, err := database.Builder.Delete(ProcessQueueMapTable).
		Where(sq.Eq{"id": id, "company_id": companyID}).ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *Service) AutoDetectProcessQueueMaps(ctx context.Context, tx pgx.Tx, companyID, connectorID int) error {
	// Get distinct process_name + folder_name from jobs
	type pf struct {
		Name   string
		Folder string
	}
	jobRows, err := tx.Query(ctx, `SELECT DISTINCT process_name, folder_name FROM orchestrator_job_executions
		WHERE company_id = $1 AND connector_id = $2 AND process_name IS NOT NULL AND folder_name IS NOT NULL`, companyID, connectorID)
	if err != nil {
		return err
	}
	defer jobRows.Close()

	folderProcesses := make(map[string][]string)
	for jobRows.Next() {
		var name, folder string
		if err := jobRows.Scan(&name, &folder); err != nil {
			continue
		}
		folderProcesses[folder] = append(folderProcesses[folder], name)
	}
	jobRows.Close()

	// Get distinct queue_name + folder_name from queue items
	queueRows, err := tx.Query(ctx, `SELECT DISTINCT queue_name, folder_name FROM orchestrator_queue_items
		WHERE company_id = $1 AND connector_id = $2 AND folder_name IS NOT NULL`, companyID, connectorID)
	if err != nil {
		return err
	}
	defer queueRows.Close()

	folderQueues := make(map[string][]string)
	for queueRows.Next() {
		var name, folder string
		if err := queueRows.Scan(&name, &folder); err != nil {
			continue
		}
		folderQueues[folder] = append(folderQueues[folder], name)
	}
	queueRows.Close()

	// Delete old auto-detected maps for this connector
	_, err = tx.Exec(ctx, `DELETE FROM orchestrator_process_queue_map WHERE company_id = $1 AND connector_id = $2 AND auto_detected = true`, companyID, connectorID)
	if err != nil {
		return err
	}

	database := db.DB()
	now := time.Now()

	// For each folder, match processes to queues
	for folder, processes := range folderProcesses {
		queues, ok := folderQueues[folder]
		if !ok || len(queues) == 0 {
			continue
		}

		for _, proc := range processes {
			for _, queue := range queues {
				sql, args, err := database.Builder.
					Insert(ProcessQueueMapTable).
					Columns("company_id", "connector_id", "process_name", "queue_name", "auto_detected", "created_at", "updated_at").
					Values(companyID, connectorID, proc, queue, true, now, now).
					Suffix("ON CONFLICT (company_id, connector_id, process_name, queue_name) DO NOTHING").
					ToSql()
				if err != nil {
					continue
				}
				tx.Exec(ctx, sql, args...)
			}
		}
	}
	return nil
}

// GetDashboardStats returns aggregated job stats per bot and recent job executions.
func (s *Service) GetDashboardStats(companyID int) (*DashboardStats, error) {
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

	// Bot stats grouped by process_name
	statsQuery := `SELECT
		COALESCE(process_name, 'N/A') AS process_name,
		COUNT(*) AS total,
		COUNT(*) FILTER (WHERE state = 'Successful') AS successful,
		COUNT(*) FILTER (WHERE state = 'Faulted') AS faulted,
		COUNT(*) FILTER (WHERE state = 'Running') AS running
	FROM orchestrator_job_executions
	WHERE company_id = $1
	GROUP BY process_name
	ORDER BY total DESC`

	rows, err := tx.Query(ctx, statsQuery, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var botStats []BotStat
	var totalJobs, totalSuccessful, totalFaulted, totalRunning int
	for rows.Next() {
		var bs BotStat
		if err := rows.Scan(&bs.ProcessName, &bs.Total, &bs.Successful, &bs.Faulted, &bs.Running); err != nil {
			return nil, err
		}
		if bs.Total > 0 {
			bs.SuccessRate = float64(bs.Successful) / float64(bs.Total) * 100
			bs.ErrorRate = float64(bs.Faulted) / float64(bs.Total) * 100
		}
		totalJobs += bs.Total
		totalSuccessful += bs.Successful
		totalFaulted += bs.Faulted
		totalRunning += bs.Running
		botStats = append(botStats, bs)
	}
	rows.Close()

	if botStats == nil {
		botStats = []BotStat{}
	}

	// Recent 10 job executions
	recentQuery := `SELECT * FROM orchestrator_job_executions
		WHERE company_id = $1
		ORDER BY start_time DESC NULLS LAST
		LIMIT 10`

	recentRows, err := tx.Query(ctx, recentQuery, companyID)
	if err != nil {
		return nil, err
	}

	recentJobs := make([]JobExecution, 0)
	if err := pgxscan.ScanAll(&recentJobs, recentRows); err != nil {
		return nil, err
	}

	return &DashboardStats{
		BotStats:   botStats,
		RecentJobs: recentJobs,
		TotalJobs:  totalJobs,
		Successful: totalSuccessful,
		Faulted:    totalFaulted,
		Running:    totalRunning,
	}, tx.Commit(ctx)
}

func (s *Service) UpsertSchedule(ctx context.Context, tx pgx.Tx, sched *ProcessSchedule) error {
	database := db.DB()
	sql, args, err := database.Builder.
		Insert(SchedulesTable).
		Columns("company_id", "connector_id", "external_schedule_id", "name", "enabled",
			"release_name", "package_name", "cron_expression", "cron_summary", "next_occurrence",
			"timezone_id", "timezone_iana", "start_strategy", "folder_name", "input_arguments",
			"created_at", "updated_at").
		Values(sched.CompanyID, sched.ConnectorID, sched.ExternalScheduleID, sched.Name, sched.Enabled,
			sched.ReleaseName, sched.PackageName, sched.CronExpression, sched.CronSummary, sched.NextOccurrence,
			sched.TimezoneID, sched.TimezoneIANA, sched.StartStrategy, sched.FolderName, sched.InputArguments,
			sched.CreatedAt, sched.UpdatedAt).
		Suffix(`ON CONFLICT (connector_id, external_schedule_id) DO UPDATE SET
			name = EXCLUDED.name, enabled = EXCLUDED.enabled, release_name = EXCLUDED.release_name,
			package_name = EXCLUDED.package_name, cron_expression = EXCLUDED.cron_expression,
			cron_summary = EXCLUDED.cron_summary, next_occurrence = EXCLUDED.next_occurrence,
			timezone_id = EXCLUDED.timezone_id, timezone_iana = EXCLUDED.timezone_iana,
			start_strategy = EXCLUDED.start_strategy, input_arguments = EXCLUDED.input_arguments,
			updated_at = EXCLUDED.updated_at`).
		ToSql()
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, sql, args...)
	return err
}

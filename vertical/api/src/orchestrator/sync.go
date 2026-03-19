package orchestrator

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"orbita-api/core"
	"orbita-api/db"

	"go.uber.org/zap"
)

func parseTime(s *string) *time.Time {
	if s == nil || *s == "" {
		return nil
	}
	for _, layout := range []string{time.RFC3339, "2006-01-02T15:04:05Z", "2006-01-02T15:04:05.999999Z"} {
		if t, err := time.Parse(layout, *s); err == nil {
			return &t
		}
	}
	return nil
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func intPtr(i int) *int { return &i }

// SyncConnector fetches data from UiPath and upserts into the local DB.
func SyncConnector(connector Connector) error {
	cfg, err := connector.GetUiPathConfig()
	if err != nil {
		return fmt.Errorf("failed to parse connector config: %w", err)
	}

	token, err := core.Decrypt(cfg.PersonalAccessToken)
	if err != nil {
		return fmt.Errorf("failed to decrypt token: %w", err)
	}
	cfg.PersonalAccessToken = token

	folders := cfg.GetFolders()
	if len(folders) == 0 {
		return fmt.Errorf("no folders configured")
	}

	client := NewUiPathClient()
	svc := &Service{}
	database := db.DB()
	ctx := context.Background()
	now := time.Now()
	var totalJobs, totalDefs, totalItems, totalSchedules, folderErrors int

	tx, err := database.C.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Global queue name map across all folders
	queueNameMap := make(map[int]string)

	for _, folder := range folders {
		folderCfg := *cfg
		folderCfg.FolderID = folder.ID
		folderCfg.FolderName = folder.Name

		zap.S().Infow("Syncing folder", "connector", connector.Name, "folder", folder.Name, "folderId", folder.ID)

		// Sync Jobs
		jobs, err := client.GetJobs(&folderCfg)
		if err != nil {
			zap.S().Warnw("Failed to sync jobs", "connector", connector.Name, "folder", folder.Name, "error", err)
		} else {
			for _, j := range jobs {
				job := &JobExecution{
					CompanyID:      connector.CompanyID,
					ConnectorID:    connector.ID,
					ExternalJobKey: strPtr(j.Key),
					ExternalJobID:  intPtr(j.ID),
					ProcessName:    strPtr(j.ReleaseName),
					State:          j.State,
					SourceType:     strPtr(j.SourceType),
					Source:         strPtr(j.Source),
					StartTime:      parseTime(j.StartTime),
					EndTime:        parseTime(j.EndTime),
					HostMachine:    strPtr(j.HostMachineName),
					FolderName:     strPtr(folder.Name),
					Info:           j.Info,
					Details:        json.RawMessage("{}"),
					CreatedAt:      now,
					UpdatedAt:      now,
				}
				if err := svc.UpsertJob(ctx, tx, job); err != nil {
					zap.S().Warnw("Failed to upsert job", "jobKey", j.Key, "error", err)
				}
			}
			totalJobs += len(jobs)
		}

		// Sync Queue Definitions
		defs, err := client.GetQueueDefinitions(&folderCfg)
		if err != nil {
			zap.S().Warnw("Failed to sync queue definitions", "connector", connector.Name, "folder", folder.Name, "error", err)
		} else {
			for _, d := range defs {
				queueNameMap[d.ID] = d.Name
				def := &QueueDefinition{
					CompanyID:            connector.CompanyID,
					ConnectorID:          connector.ID,
					ExternalDefinitionID: intPtr(d.ID),
					Name:                 d.Name,
					MaxRetries:           d.MaxNumberOfRetries,
					FolderName:           strPtr(folder.Name),
					CreatedAt:            now,
					UpdatedAt:            now,
				}
				if err := svc.UpsertQueueDefinition(ctx, tx, def); err != nil {
					zap.S().Warnw("Failed to upsert queue definition", "name", d.Name, "error", err)
				}
			}
			totalDefs += len(defs)
		}

		// Sync Queue Items
		items, err := client.GetQueueItems(&folderCfg)
		if err != nil {
			zap.S().Warnw("Failed to sync queue items", "connector", connector.Name, "folder", folder.Name, "error", err)
		} else {
			for _, qi := range items {
				var errorMsg *string
				if qi.Error != nil {
					switch v := qi.Error.(type) {
					case string:
						errorMsg = &v
					case map[string]interface{}:
						if msg, ok := v["Message"].(string); ok && msg != "" {
							errorMsg = &msg
						} else if msg, ok := v["Title"].(string); ok && msg != "" {
							errorMsg = &msg
						} else if b, err := json.Marshal(v); err == nil {
							s := string(b)
							errorMsg = &s
						}
					default:
						s := fmt.Sprintf("%v", v)
						errorMsg = &s
					}
				}

				var specificContent json.RawMessage
				if qi.SpecificContent != nil {
					if sc, err := json.Marshal(qi.SpecificContent); err == nil {
						specificContent = sc
					}
				}
				if specificContent == nil {
					specificContent = json.RawMessage("{}")
				}

				queueName := queueNameMap[qi.QueueDefinitionID]
				if queueName == "" {
					queueName = fmt.Sprintf("Queue_%d", qi.QueueDefinitionID)
				}

				item := &QueueItem{
					CompanyID:               connector.CompanyID,
					ConnectorID:             connector.ID,
					ExternalItemKey:         strPtr(qi.Key),
					ExternalItemID:          intPtr(qi.ID),
					QueueDefinitionID:       intPtr(qi.QueueDefinitionID),
					QueueName:               queueName,
					Status:                  qi.Status,
					Priority:                strPtr(qi.Priority),
					Reference:               qi.Reference,
					ProcessingExceptionType: qi.ProcessingExceptionType,
					ErrorMessage:            errorMsg,
					StartProcessing:         parseTime(qi.StartProcessing),
					EndProcessing:           parseTime(qi.EndProcessing),
					RetryNumber:             qi.RetryNumber,
					FolderName:              strPtr(folder.Name),
					SpecificContent:         specificContent,
					CreatedAt:               now,
					UpdatedAt:               now,
				}
				if err := svc.UpsertQueueItem(ctx, tx, item); err != nil {
					zap.S().Warnw("Failed to upsert queue item", "itemKey", qi.Key, "error", err)
				}
			}
			totalItems += len(items)
		}

		// Sync Schedules
		schedules, err := client.GetSchedules(&folderCfg)
		if err != nil {
			zap.S().Warnw("Failed to sync schedules", "connector", connector.Name, "folder", folder.Name, "error", err)
		} else {
			for _, ps := range schedules {
				var inputArgs json.RawMessage
				if ps.InputArguments != nil && *ps.InputArguments != "" {
					inputArgs = json.RawMessage(*ps.InputArguments)
				}
				if inputArgs == nil {
					inputArgs = json.RawMessage("{}")
				}

				sched := &ProcessSchedule{
					CompanyID:          connector.CompanyID,
					ConnectorID:        connector.ID,
					ExternalScheduleID: intPtr(ps.ID),
					Name:               ps.Name,
					Enabled:            ps.Enabled,
					ReleaseName:        ps.ReleaseName,
					PackageName:        ps.PackageName,
					CronExpression:     strPtr(ps.StartProcessCron),
					CronSummary:        ps.StartProcessCronSummary,
					NextOccurrence:     parseTime(ps.StartProcessNextOccurrence),
					TimezoneID:         strPtr(ps.TimeZoneId),
					TimezoneIANA:       strPtr(ps.TimeZoneIana),
					StartStrategy:      ps.StartStrategy,
					FolderName:         strPtr(folder.Name),
					InputArguments:     inputArgs,
					CreatedAt:          now,
					UpdatedAt:          now,
				}
				if err := svc.UpsertSchedule(ctx, tx, sched); err != nil {
					zap.S().Warnw("Failed to upsert schedule", "name", ps.Name, "error", err)
				}
			}
			totalSchedules += len(schedules)
		}

		// Count folder-level failures
		if jobs == nil && defs == nil && items == nil && schedules == nil {
			folderErrors++
		}
	}

	zap.S().Infow("Sync totals", "connector", connector.Name, "folders", len(folders), "jobs", totalJobs, "definitions", totalDefs, "queueItems", totalItems, "schedules", totalSchedules)

	if folderErrors == len(folders) {
		return fmt.Errorf("all folders failed to sync, check credentials and network")
	}

	// Auto-detect process-queue mappings based on shared folder
	if err := svc.AutoDetectProcessQueueMaps(ctx, tx, connector.CompanyID, connector.ID); err != nil {
		zap.S().Warnw("Failed to auto-detect process-queue maps", "error", err)
	}

	return tx.Commit(ctx)
}

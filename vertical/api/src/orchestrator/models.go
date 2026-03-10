package orchestrator

import (
	"encoding/json"
	"time"
)

const (
	ConnectorsTable       = "orchestrator_connectors"
	JobExecutionsTable    = "orchestrator_job_executions"
	QueueDefinitionsTable = "orchestrator_queue_definitions"
	QueueItemsTable       = "orchestrator_queue_items"
	SchedulesTable        = "orchestrator_schedules"
	ProcessQueueMapTable  = "orchestrator_process_queue_map"

	TypeUiPath      = "UIPATH"
	TypePythonAgent = "PYTHON_AGENT"
)

// Connector represents a connection to an external orchestration platform.
type Connector struct {
	ID        int             `db:"id" json:"id"`
	CompanyID int             `db:"company_id" json:"companyId"`
	Name      string          `db:"name" json:"name"`
	Type      string          `db:"type" json:"type"`
	Config    json.RawMessage `db:"config" json:"config"`
	IsActive  bool            `db:"is_active" json:"isActive"`
	CreatedAt time.Time       `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time       `db:"updated_at" json:"updatedAt"`
}

// UiPathFolder represents a single UiPath folder to sync.
type UiPathFolder struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// UiPathConfig holds the decrypted connector config for UiPath.
type UiPathConfig struct {
	OrganizationName    string         `json:"organizationName"`
	TenantName          string         `json:"tenantName"`
	PersonalAccessToken string         `json:"personalAccessToken"`
	FolderID            string         `json:"folderId"`
	FolderName          string         `json:"folderName"`
	Folders             []UiPathFolder `json:"folders,omitempty"`
}

// GetFolders returns the list of folders to sync. Falls back to single FolderID/FolderName.
func (c *UiPathConfig) GetFolders() []UiPathFolder {
	if len(c.Folders) > 0 {
		return c.Folders
	}
	if c.FolderID != "" {
		return []UiPathFolder{{ID: c.FolderID, Name: c.FolderName}}
	}
	return nil
}

func (c *Connector) GetUiPathConfig() (*UiPathConfig, error) {
	cfg := &UiPathConfig{}
	if len(c.Config) == 0 {
		return cfg, nil
	}
	return cfg, json.Unmarshal(c.Config, cfg)
}

// JobExecution represents a bot execution synced from UiPath.
type JobExecution struct {
	ID             int        `db:"id" json:"id"`
	CompanyID      int        `db:"company_id" json:"companyId"`
	ConnectorID    int        `db:"connector_id" json:"connectorId"`
	ExternalJobKey *string    `db:"external_job_key" json:"externalJobKey,omitempty"`
	ExternalJobID  *int       `db:"external_job_id" json:"externalJobId,omitempty"`
	ProcessName    *string    `db:"process_name" json:"processName,omitempty"`
	State          string     `db:"state" json:"state"`
	SourceType     *string    `db:"source_type" json:"sourceType,omitempty"`
	Source         *string    `db:"source" json:"source,omitempty"`
	StartTime      *time.Time `db:"start_time" json:"startTime,omitempty"`
	EndTime        *time.Time `db:"end_time" json:"endTime,omitempty"`
	HostMachine    *string    `db:"host_machine" json:"hostMachine,omitempty"`
	FolderName     *string    `db:"folder_name" json:"folderName,omitempty"`
	Info           *string    `db:"info" json:"info,omitempty"`
	Details        json.RawMessage `db:"details" json:"details,omitempty"`
	CreatedAt      time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt      time.Time  `db:"updated_at" json:"updatedAt"`
}

// QueueDefinition represents a UiPath queue definition.
type QueueDefinition struct {
	ID                   int       `db:"id" json:"id"`
	CompanyID            int       `db:"company_id" json:"companyId"`
	ConnectorID          int       `db:"connector_id" json:"connectorId"`
	ExternalDefinitionID *int      `db:"external_definition_id" json:"externalDefinitionId,omitempty"`
	Name                 string    `db:"name" json:"name"`
	MaxRetries           int       `db:"max_retries" json:"maxRetries"`
	FolderName           *string   `db:"folder_name" json:"folderName,omitempty"`
	CreatedAt            time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt            time.Time `db:"updated_at" json:"updatedAt"`
}

// ProcessSchedule represents a UiPath process schedule (trigger).
type ProcessSchedule struct {
	ID                 int             `db:"id" json:"id"`
	CompanyID          int             `db:"company_id" json:"companyId"`
	ConnectorID        int             `db:"connector_id" json:"connectorId"`
	ExternalScheduleID *int            `db:"external_schedule_id" json:"externalScheduleId,omitempty"`
	Name               string          `db:"name" json:"name"`
	Enabled            bool            `db:"enabled" json:"enabled"`
	ReleaseName        *string         `db:"release_name" json:"releaseName,omitempty"`
	PackageName        *string         `db:"package_name" json:"packageName,omitempty"`
	CronExpression     *string         `db:"cron_expression" json:"cronExpression,omitempty"`
	CronSummary        *string         `db:"cron_summary" json:"cronSummary,omitempty"`
	NextOccurrence     *time.Time      `db:"next_occurrence" json:"nextOccurrence,omitempty"`
	TimezoneID         *string         `db:"timezone_id" json:"timezoneId,omitempty"`
	TimezoneIANA       *string         `db:"timezone_iana" json:"timezoneIana,omitempty"`
	StartStrategy      int             `db:"start_strategy" json:"startStrategy"`
	FolderName         *string         `db:"folder_name" json:"folderName,omitempty"`
	InputArguments     json.RawMessage `db:"input_arguments" json:"inputArguments,omitempty"`
	CreatedAt          time.Time       `db:"created_at" json:"createdAt"`
	UpdatedAt          time.Time       `db:"updated_at" json:"updatedAt"`
}

// ProcessQueueMap represents a mapping between a process (bot) and a queue.
type ProcessQueueMap struct {
	ID           int       `db:"id" json:"id"`
	CompanyID    int       `db:"company_id" json:"companyId"`
	ConnectorID  int       `db:"connector_id" json:"connectorId"`
	ProcessName  string    `db:"process_name" json:"processName"`
	QueueName    string    `db:"queue_name" json:"queueName"`
	AutoDetected bool      `db:"auto_detected" json:"autoDetected"`
	CreatedAt    time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt    time.Time `db:"updated_at" json:"updatedAt"`
}

// QueueItem represents an element in a UiPath queue.
type QueueItem struct {
	ID                      int             `db:"id" json:"id"`
	CompanyID               int             `db:"company_id" json:"companyId"`
	ConnectorID             int             `db:"connector_id" json:"connectorId"`
	ExternalItemKey         *string         `db:"external_item_key" json:"externalItemKey,omitempty"`
	ExternalItemID          *int            `db:"external_item_id" json:"externalItemId,omitempty"`
	QueueDefinitionID       *int            `db:"queue_definition_id" json:"queueDefinitionId,omitempty"`
	QueueName               string          `db:"queue_name" json:"queueName"`
	Status                  string          `db:"status" json:"status"`
	Priority                *string         `db:"priority" json:"priority,omitempty"`
	Reference               *string         `db:"reference" json:"reference,omitempty"`
	ProcessingExceptionType *string         `db:"processing_exception_type" json:"processingExceptionType,omitempty"`
	ErrorMessage            *string         `db:"error_message" json:"errorMessage,omitempty"`
	StartProcessing         *time.Time      `db:"start_processing" json:"startProcessing,omitempty"`
	EndProcessing           *time.Time      `db:"end_processing" json:"endProcessing,omitempty"`
	RetryNumber             int             `db:"retry_number" json:"retryNumber"`
	FolderName              *string         `db:"folder_name" json:"folderName,omitempty"`
	SpecificContent         json.RawMessage `db:"specific_content" json:"specificContent,omitempty"`
	CreatedAt               time.Time       `db:"created_at" json:"createdAt"`
	UpdatedAt               time.Time       `db:"updated_at" json:"updatedAt"`
}

package orchestrator

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

func parseCSV(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

var validConnectorTypes = []string{TypeUiPath}

// ConnectorBody is the input for creating/updating a connector.
type ConnectorBody struct {
	Name             string         `json:"name" binding:"required"`
	Type             string         `json:"type" binding:"required"`
	OrganizationName string         `json:"organizationName" binding:"required"`
	TenantName       string         `json:"tenantName" binding:"required"`
	AccessToken      string         `json:"accessToken"`
	FolderID         string         `json:"folderId"`
	FolderName       string         `json:"folderName"`
	Folders          []UiPathFolder `json:"folders"`
	IsActive         *bool          `json:"isActive"`
}

func (b *ConnectorBody) Bind(c *gin.Context) error {
	if err := c.ShouldBindJSON(b); err != nil {
		return err
	}
	valid := false
	for _, t := range validConnectorTypes {
		if b.Type == t {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("invalid connector type: must be one of %v", validConnectorTypes)
	}
	return nil
}

// JobListParams holds filters for listing job executions.
type JobListParams struct {
	State        string `form:"state"`
	ConnectorID  int    `form:"connectorId"`
	ProcessNames string `form:"processNames"`
	Page         int    `form:"page,default=1"`
	Limit        int    `form:"limit,default=20"`
}

// QueueItemListParams holds filters for listing queue items.
type QueueItemListParams struct {
	Status       string `form:"status"`
	ConnectorID  int    `form:"connectorId"`
	QueueName    string `form:"queueName"`
	ProcessNames string `form:"processNames"`
	Page         int    `form:"page,default=1"`
	Limit        int    `form:"limit,default=20"`
}

// ScheduleListParams holds filters for listing schedules.
type ScheduleListParams struct {
	Enabled      *bool  `form:"enabled"`
	ConnectorID  int    `form:"connectorId"`
	ProcessNames string `form:"processNames"`
	Page         int    `form:"page,default=1"`
	Limit        int    `form:"limit,default=20"`
}

// PaginatedResponse wraps any list response with pagination info.
type PaginatedResponse struct {
	Data       any `json:"data"`
	Total      int `json:"total"`
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalPages int `json:"totalPages"`
}

// ConnectorResponse omits sensitive fields from the config.
type ConnectorResponse struct {
	ID        int    `json:"id"`
	CompanyID int    `json:"companyId"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	IsActive  bool   `json:"isActive"`
	// Safe config fields (no token)
	OrganizationName string         `json:"organizationName"`
	TenantName       string         `json:"tenantName"`
	FolderID         string         `json:"folderId"`
	FolderName       string         `json:"folderName"`
	Folders          []UiPathFolder `json:"folders"`
	CreatedAt        string         `json:"createdAt"`
	UpdatedAt        string         `json:"updatedAt"`
}

// BotStat holds aggregated job stats for a single bot (process_name).
type BotStat struct {
	ProcessName string  `json:"processName"`
	Total       int     `json:"total"`
	Successful  int     `json:"successful"`
	Faulted     int     `json:"faulted"`
	Running     int     `json:"running"`
	SuccessRate float64 `json:"successRate"`
	ErrorRate   float64 `json:"errorRate"`
}

// DashboardStats holds orchestrator overview data for the home dashboard.
type DashboardStats struct {
	BotStats   []BotStat      `json:"botStats"`
	RecentJobs []JobExecution `json:"recentJobs"`
	TotalJobs  int            `json:"totalJobs"`
	Successful int            `json:"successful"`
	Faulted    int            `json:"faulted"`
	Running    int            `json:"running"`
}

func NewConnectorResponse(c Connector) ConnectorResponse {
	cfg, _ := c.GetUiPathConfig()
	return ConnectorResponse{
		ID:               c.ID,
		CompanyID:        c.CompanyID,
		Name:             c.Name,
		Type:             c.Type,
		IsActive:         c.IsActive,
		OrganizationName: cfg.OrganizationName,
		TenantName:       cfg.TenantName,
		FolderID:         cfg.FolderID,
		FolderName:       cfg.FolderName,
		Folders:          cfg.GetFolders(),
		CreatedAt:        c.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:        c.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

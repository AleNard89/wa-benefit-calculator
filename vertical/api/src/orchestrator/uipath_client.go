package orchestrator

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"
)

const uiPathBaseURL = "https://cloud.uipath.com"

type UiPathClient struct {
	httpClient *http.Client
}

func NewUiPathClient() *UiPathClient {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	if os.Getenv("TLS_SKIP_VERIFY") == "true" {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	return &UiPathClient{
		httpClient: &http.Client{Timeout: 30 * time.Second, Transport: transport},
	}
}

func (c *UiPathClient) baseURL(cfg *UiPathConfig) string {
	return fmt.Sprintf("%s/%s/%s/orchestrator_/odata", uiPathBaseURL, cfg.OrganizationName, cfg.TenantName)
}

func (c *UiPathClient) doRequest(cfg *UiPathConfig, endpoint string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", c.baseURL(cfg), endpoint)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cfg.PersonalAccessToken))
	req.Header.Set("X-UIPATH-OrganizationUnitId", cfg.FolderID)
	req.Header.Set("Content-Type", "application/json")

	zap.S().Debugw("UiPath request", "url", url, "folder_id", cfg.FolderID)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("UiPath API request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		zap.S().Warnw("UiPath API error",
			"status", resp.StatusCode,
			"url", url,
			"body", string(body),
			"response_headers", resp.Header,
			"folder_id", cfg.FolderID,
		)
		return nil, fmt.Errorf("UiPath API returned %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// OData response wrapper
type odataResponse struct {
	Value json.RawMessage `json:"value"`
}

// UiPath API response structs
type UiPathJob struct {
	Key                                  string  `json:"Key"`
	ID                                   int     `json:"Id"`
	State                                string  `json:"State"`
	ReleaseName                          string  `json:"ReleaseName"`
	StartTime                            *string `json:"StartTime"`
	EndTime                              *string `json:"EndTime"`
	Source                               string  `json:"Source"`
	SourceType                           string  `json:"SourceType"`
	Info                                 *string `json:"Info"`
	HostMachineName                      string  `json:"HostMachineName"`
	OrganizationUnitFullyQualifiedName   string  `json:"OrganizationUnitFullyQualifiedName"`
}

type UiPathQueueItem struct {
	Key                                string  `json:"Key"`
	ID                                 int     `json:"Id"`
	Status                             string  `json:"Status"`
	QueueDefinitionID                  int     `json:"QueueDefinitionId"`
	Reference                          *string `json:"Reference"`
	Priority                           string  `json:"Priority"`
	ProcessingExceptionType            *string `json:"ProcessingExceptionType"`
	Error                              any     `json:"Error"`
	StartProcessing                    *string `json:"StartProcessing"`
	EndProcessing                      *string `json:"EndProcessing"`
	RetryNumber                        int     `json:"RetryNumber"`
	CreationTime                       string  `json:"CreationTime"`
	OrganizationUnitFullyQualifiedName string  `json:"OrganizationUnitFullyQualifiedName"`
	SpecificContent                    any     `json:"SpecificContent"`
}

type UiPathQueueDefinition struct {
	Key                                string `json:"Key"`
	ID                                 int    `json:"Id"`
	Name                               string `json:"Name"`
	MaxNumberOfRetries                 int    `json:"MaxNumberOfRetries"`
	OrganizationUnitFullyQualifiedName string `json:"OrganizationUnitFullyQualifiedName"`
}

type UiPathProcessSchedule struct {
	ID                          int     `json:"Id"`
	Name                        string  `json:"Name"`
	Enabled                     bool    `json:"Enabled"`
	ReleaseID                   int     `json:"ReleaseId"`
	ReleaseName                 *string `json:"ReleaseName"`
	PackageName                 *string `json:"PackageName"`
	StartProcessCron            string  `json:"StartProcessCron"`
	StartProcessCronSummary     *string `json:"StartProcessCronSummary"`
	StartProcessNextOccurrence  *string `json:"StartProcessNextOccurrence"`
	StartStrategy               int     `json:"StartStrategy"`
	TimeZoneId                  string  `json:"TimeZoneId"`
	TimeZoneIana                string  `json:"TimeZoneIana"`
	InputArguments              *string `json:"InputArguments"`
	OrganizationUnitFullyQualifiedName string `json:"OrganizationUnitFullyQualifiedName"`
}

func (c *UiPathClient) GetSchedules(cfg *UiPathConfig) ([]UiPathProcessSchedule, error) {
	data, err := c.doRequest(cfg, "ProcessSchedules")
	if err != nil {
		return nil, err
	}
	var odata odataResponse
	if err := json.Unmarshal(data, &odata); err != nil {
		return nil, err
	}
	var schedules []UiPathProcessSchedule
	return schedules, json.Unmarshal(odata.Value, &schedules)
}

func (c *UiPathClient) GetJobs(cfg *UiPathConfig) ([]UiPathJob, error) {
	data, err := c.doRequest(cfg, "Jobs")
	if err != nil {
		return nil, err
	}
	var odata odataResponse
	if err := json.Unmarshal(data, &odata); err != nil {
		return nil, err
	}
	var jobs []UiPathJob
	return jobs, json.Unmarshal(odata.Value, &jobs)
}

func (c *UiPathClient) GetQueueItems(cfg *UiPathConfig) ([]UiPathQueueItem, error) {
	data, err := c.doRequest(cfg, "QueueItems")
	if err != nil {
		return nil, err
	}
	var odata odataResponse
	if err := json.Unmarshal(data, &odata); err != nil {
		return nil, err
	}
	var items []UiPathQueueItem
	return items, json.Unmarshal(odata.Value, &items)
}

func (c *UiPathClient) GetQueueDefinitions(cfg *UiPathConfig) ([]UiPathQueueDefinition, error) {
	data, err := c.doRequest(cfg, "QueueDefinitions")
	if err != nil {
		return nil, err
	}
	var odata odataResponse
	if err := json.Unmarshal(data, &odata); err != nil {
		return nil, err
	}
	var defs []UiPathQueueDefinition
	return defs, json.Unmarshal(odata.Value, &defs)
}

func (c *UiPathClient) TestConnection(cfg *UiPathConfig) error {
	_, err := c.doRequest(cfg, "Settings")
	return err
}

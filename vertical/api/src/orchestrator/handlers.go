package orchestrator

import (
	"math"
	"net/http"
	"strconv"
	"time"

	"benefit-calculator-api/core"
	"benefit-calculator-api/core/httpx"

	d "benefit-calculator-api/core/decorators"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// --- Connectors (admin only) ---

func requireCompanyID(c *gin.Context) int {
	companyID := httpx.HeaderCompanyID(c)
	if companyID == 0 {
		c.JSON(http.StatusBadRequest, httpx.ErrorResponse{Message: "Missing company id"})
	}
	return companyID
}

func listConnectors(c *gin.Context) {
	companyID := requireCompanyID(c)
	if companyID == 0 {
		return
	}
	svc := &Service{}
	connectors, err := svc.ListConnectors(companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}
	resp := make([]ConnectorResponse, len(connectors))
	for i, conn := range connectors {
		resp[i] = NewConnectorResponse(conn)
	}
	c.JSON(http.StatusOK, resp)
}

var ListConnectors = d.SuperuserRequired(listConnectors)

func createConnector(c *gin.Context) {
	companyID := requireCompanyID(c)
	if companyID == 0 {
		return
	}

	body := ConnectorBody{}
	if err := body.Bind(c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, httpx.ErrorResponse{Message: err.Error()})
		return
	}

	if body.AccessToken == "" {
		c.JSON(http.StatusUnprocessableEntity, httpx.ErrorResponse{Message: "Access token is required"})
		return
	}

	svc := &Service{}
	connector, err := svc.CreateConnector(companyID, body, core.Encrypt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}

	zap.S().Infow("Connettore creato", "id", connector.ID, "name", connector.Name, "company_id", companyID)
	c.JSON(http.StatusCreated, NewConnectorResponse(*connector))
}

var CreateConnector = d.SuperuserRequired(createConnector)

func updateConnector(c *gin.Context) {
	companyID := requireCompanyID(c)
	if companyID == 0 {
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, httpx.ErrorResponse{Message: "Invalid connector ID"})
		return
	}

	body := ConnectorBody{}
	if err := body.Bind(c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, httpx.ErrorResponse{Message: err.Error()})
		return
	}

	svc := &Service{}
	connector, err := svc.UpdateConnector(id, companyID, body, core.Encrypt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}

	zap.S().Infow("Connettore aggiornato", "id", connector.ID, "name", connector.Name)
	c.JSON(http.StatusOK, NewConnectorResponse(*connector))
}

var UpdateConnector = d.SuperuserRequired(updateConnector)

func testConnector(c *gin.Context) {
	companyID := requireCompanyID(c)
	if companyID == 0 {
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, httpx.ErrorResponse{Message: "Invalid connector ID"})
		return
	}

	svc := &Service{}
	connector, err := svc.GetConnector(id, companyID)
	if err != nil {
		c.JSON(http.StatusNotFound, httpx.ErrorResponse{Message: "Connector not found"})
		return
	}

	cfg, err := connector.GetUiPathConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Message: "Invalid connector config"})
		return
	}

	token, err := core.Decrypt(cfg.PersonalAccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Message: "Failed to decrypt token"})
		return
	}
	cfg.PersonalAccessToken = token

	client := NewUiPathClient()
	if err := client.TestConnection(cfg); err != nil {
		c.JSON(http.StatusBadGateway, httpx.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Connessione riuscita"})
}

var TestConnector = d.SuperuserRequired(testConnector)

func syncConnector(c *gin.Context) {
	companyID := requireCompanyID(c)
	if companyID == 0 {
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, httpx.ErrorResponse{Message: "Invalid connector ID"})
		return
	}

	svc := &Service{}
	connector, err := svc.GetConnector(id, companyID)
	if err != nil {
		c.JSON(http.StatusNotFound, httpx.ErrorResponse{Message: "Connector not found"})
		return
	}

	if err := SyncConnector(*connector); err != nil {
		zap.S().Errorw("Sync failed", "connector", connector.Name, "error", err)
		c.JSON(http.StatusBadGateway, httpx.ErrorResponse{Message: err.Error()})
		return
	}

	zap.S().Infow("Sync completato", "connector", connector.Name, "company_id", companyID)
	c.JSON(http.StatusOK, gin.H{"message": "Sincronizzazione completata"})
}

var SyncConnectorHandler = d.SuperuserRequired(syncConnector)

// --- Dashboard Stats (all authenticated users) ---

func getDashboardStats(c *gin.Context) {
	companyID := requireCompanyID(c)
	if companyID == 0 {
		return
	}
	svc := &Service{}
	stats, err := svc.GetDashboardStats(companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}
	c.JSON(http.StatusOK, stats)
}

var GetDashboardStats = d.AuthenticationRequired(getDashboardStats)

// --- Process Names (all authenticated users) ---

func listProcessNames(c *gin.Context) {
	companyID := requireCompanyID(c)
	if companyID == 0 {
		return
	}

	svc := &Service{}
	names, err := svc.GetDistinctProcessCodes(companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}
	if names == nil {
		names = []string{}
	}
	c.JSON(http.StatusOK, names)
}

var ListProcessNames = d.AuthenticationRequired(listProcessNames)

func listBotNames(c *gin.Context) {
	companyID := requireCompanyID(c)
	if companyID == 0 {
		return
	}

	svc := &Service{}
	names, err := svc.GetDistinctBotNames(companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}
	if names == nil {
		names = []string{}
	}
	c.JSON(http.StatusOK, names)
}

var ListBotNames = d.AuthenticationRequired(listBotNames)

// --- Jobs (all authenticated users) ---

func listJobs(c *gin.Context) {
	companyID := requireCompanyID(c)
	if companyID == 0 {
		return
	}
	params := JobListParams{}
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, httpx.ErrorResponse{Message: err.Error()})
		return
	}
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 || params.Limit > 100 {
		params.Limit = 20
	}

	svc := &Service{}
	jobs, total, err := svc.ListJobs(companyID, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(params.Limit)))
	c.JSON(http.StatusOK, PaginatedResponse{Data: jobs, Total: total, Page: params.Page, Limit: params.Limit, TotalPages: totalPages})
}

var ListJobs = d.AuthenticationRequired(listJobs)

func getJob(c *gin.Context) {
	companyID := requireCompanyID(c)
	if companyID == 0 {
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, httpx.ErrorResponse{Message: "Invalid job ID"})
		return
	}

	svc := &Service{}
	job, err := svc.GetJob(id, companyID)
	if err != nil {
		c.JSON(http.StatusNotFound, httpx.ErrorResponse{Message: "Job not found"})
		return
	}
	c.JSON(http.StatusOK, job)
}

var GetJob = d.AuthenticationRequired(getJob)

// --- Queue Items (all authenticated users) ---

func listQueueItems(c *gin.Context) {
	companyID := requireCompanyID(c)
	if companyID == 0 {
		return
	}
	params := QueueItemListParams{}
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, httpx.ErrorResponse{Message: err.Error()})
		return
	}
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 || params.Limit > 100 {
		params.Limit = 20
	}

	svc := &Service{}
	items, total, err := svc.ListQueueItems(companyID, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(params.Limit)))
	c.JSON(http.StatusOK, PaginatedResponse{Data: items, Total: total, Page: params.Page, Limit: params.Limit, TotalPages: totalPages})
}

var ListQueueItems = d.AuthenticationRequired(listQueueItems)

// --- Schedules (all authenticated users) ---

func listSchedules(c *gin.Context) {
	companyID := requireCompanyID(c)
	if companyID == 0 {
		return
	}
	params := ScheduleListParams{}
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, httpx.ErrorResponse{Message: err.Error()})
		return
	}
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 || params.Limit > 100 {
		params.Limit = 20
	}

	svc := &Service{}
	schedules, total, err := svc.ListSchedules(companyID, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(params.Limit)))
	c.JSON(http.StatusOK, PaginatedResponse{Data: schedules, Total: total, Page: params.Page, Limit: params.Limit, TotalPages: totalPages})
}

var ListSchedules = d.AuthenticationRequired(listSchedules)

// --- Queue Definitions (all authenticated users) ---

func listQueueDefinitions(c *gin.Context) {
	companyID := requireCompanyID(c)
	if companyID == 0 {
		return
	}
	connectorID := 0
	if id := c.Query("connectorId"); id != "" {
		connectorID, _ = strconv.Atoi(id)
	}

	svc := &Service{}
	defs, err := svc.ListQueueDefinitions(companyID, connectorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}
	c.JSON(http.StatusOK, defs)
}

var ListQueueDefinitions = d.AuthenticationRequired(listQueueDefinitions)

// --- Process-Queue Map ---

func listProcessQueueMaps(c *gin.Context) {
	companyID := requireCompanyID(c)
	if companyID == 0 {
		return
	}
	svc := &Service{}
	maps, err := svc.ListProcessQueueMaps(companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}
	c.JSON(http.StatusOK, maps)
}

var ListProcessQueueMaps = d.AuthenticationRequired(listProcessQueueMaps)

type processQueueMapBody struct {
	ConnectorID int    `json:"connectorId" binding:"required"`
	ProcessName string `json:"processName" binding:"required"`
	QueueName   string `json:"queueName" binding:"required"`
}

func createProcessQueueMap(c *gin.Context) {
	companyID := requireCompanyID(c)
	if companyID == 0 {
		return
	}
	var body processQueueMapBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, httpx.ErrorResponse{Message: err.Error()})
		return
	}

	svc := &Service{}
	now := time.Now()
	m := &ProcessQueueMap{
		CompanyID:    companyID,
		ConnectorID:  body.ConnectorID,
		ProcessName:  body.ProcessName,
		QueueName:    body.QueueName,
		AutoDetected: false,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := svc.CreateProcessQueueMap(companyID, m); err != nil {
		c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}
	zap.S().Infow("Mapping processo-coda creato", "process", body.ProcessName, "queue", body.QueueName)
	c.JSON(http.StatusCreated, m)
}

var CreateProcessQueueMap = d.SuperuserRequired(createProcessQueueMap)

func deleteProcessQueueMap(c *gin.Context) {
	companyID := requireCompanyID(c)
	if companyID == 0 {
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, httpx.ErrorResponse{Message: "Invalid map ID"})
		return
	}

	svc := &Service{}
	if err := svc.DeleteProcessQueueMap(id, companyID); err != nil {
		c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}
	zap.S().Infow("Mapping processo-coda eliminato", "id", id)
	c.JSON(http.StatusOK, gin.H{"message": "Mapping eliminato"})
}

var DeleteProcessQueueMap = d.SuperuserRequired(deleteProcessQueueMap)

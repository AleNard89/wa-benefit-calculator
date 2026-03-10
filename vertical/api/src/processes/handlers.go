package processes

import (
	"math"
	"net/http"
	"strconv"

	"benefit-calculator-api/auth"
	"benefit-calculator-api/core/httpx"
	"benefit-calculator-api/orgs"

	d "benefit-calculator-api/core/decorators"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func getUserAreaIDs(c *gin.Context) []int {
	user := auth.MustCurrentUser(c)
	if user.IsSuperUser() {
		return nil
	}
	companyID := httpx.HeaderCompanyID(c)
	areaService := orgs.AreaService{}
	areas, err := areaService.GetUserAreas(user.ID, companyID)
	if err != nil {
		return []int{}
	}
	ids := make([]int, len(areas))
	for i, a := range areas {
		ids[i] = a.ID
	}
	return ids
}

func getProcesses(c *gin.Context) {
	companyID := httpx.HeaderCompanyID(c)
	params := ProcessListParams{}
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}

	// Only superuser or users with process.delete can view deleted processes
	if params.Deleted {
		user := auth.MustCurrentUser(c)
		if !user.IsSuperUser() {
			perms := user.AllCompanyPermissionsCodes(companyID)
			hasDeletePerm := false
			for _, p := range perms {
				if p == ProcessDelete {
					hasDeletePerm = true
					break
				}
			}
			if !hasDeletePerm {
				c.JSON(http.StatusForbidden, httpx.ErrorResponse{Message: "Non hai i permessi per visualizzare i processi eliminati"})
				return
			}
		}
	}

	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 || params.Limit > 100 {
		params.Limit = 25
	}
	params.AreaIDs = getUserAreaIDs(c)

	service := ProcessService{}
	processes, total, err := service.GetAll(companyID, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(params.Limit)))
	c.JSON(http.StatusOK, ProcessListResponse{
		Data:       processes,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	})
}

var GetProcesses = d.CompanyRequiredAndPermissionRequired([]string{ProcessRead}, getProcesses)

func getProcess(c *gin.Context) {
	companyID := httpx.HeaderCompanyID(c)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, httpx.ErrorResponse{Message: "Invalid process ID"})
		return
	}

	service := ProcessService{}
	process, err := service.GetByID(id, companyID)
	if err != nil {
		c.JSON(http.StatusNotFound, httpx.ErrorResponse{Message: "Process not found"})
		return
	}
	c.JSON(http.StatusOK, process)
}

var GetProcess = d.CompanyRequiredAndPermissionRequired([]string{ProcessRead}, getProcess)

func createProcess(c *gin.Context) {
	companyID := httpx.HeaderCompanyID(c)
	payload := ProcessBody{}
	if err := payload.Bind(c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}

	process := Process{}
	if err := process.ApplyPayload(payload); err != nil {
		c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}

	user := auth.MustCurrentUser(c)
	if !user.IsAnonymous() {
		userID := user.ID
		process.CreatedBy = &userID
	}

	if err := process.Save(companyID); err != nil {
		c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}

	zap.S().Infow("Processo creato", "id", process.ID, "name", process.ProcessName, "company_id", companyID)
	c.JSON(http.StatusCreated, process)
}

var CreateProcess = d.CompanyRequiredAndPermissionRequired([]string{ProcessCreate}, createProcess)

func updateProcess(c *gin.Context) {
	companyID := httpx.HeaderCompanyID(c)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, httpx.ErrorResponse{Message: "Invalid process ID"})
		return
	}

	service := ProcessService{}
	process, err := service.GetByID(id, companyID)
	if err != nil {
		c.JSON(http.StatusNotFound, httpx.ErrorResponse{Message: "Process not found"})
		return
	}

	payload := ProcessBody{}
	if err := payload.Bind(c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}

	if err := process.ApplyPayload(payload); err != nil {
		c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}
	if err := process.Save(companyID); err != nil {
		c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}

	zap.S().Infow("Processo aggiornato", "id", process.ID, "name", process.ProcessName)
	c.JSON(http.StatusOK, process)
}

var UpdateProcess = d.CompanyRequiredAndPermissionRequired([]string{ProcessUpdate}, updateProcess)

func updateProcessStatus(c *gin.Context) {
	companyID := httpx.HeaderCompanyID(c)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, httpx.ErrorResponse{Message: "Invalid process ID"})
		return
	}

	service := ProcessService{}
	process, err := service.GetByID(id, companyID)
	if err != nil {
		c.JSON(http.StatusNotFound, httpx.ErrorResponse{Message: "Process not found"})
		return
	}

	payload := StatusBody{}
	if err := payload.Bind(c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}

	if err := service.UpdateStatus(id, companyID, payload.Status); err != nil {
		c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}

	process.Status = payload.Status
	zap.S().Infow("Stato processo aggiornato", "id", id, "status", payload.Status)
	c.JSON(http.StatusOK, process)
}

var UpdateProcessStatus = d.CompanyRequiredAndPermissionRequired([]string{ProcessUpdate}, updateProcessStatus)

func deleteProcess(c *gin.Context) {
	companyID := httpx.HeaderCompanyID(c)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, httpx.ErrorResponse{Message: "Invalid process ID"})
		return
	}

	service := ProcessService{}
	process, err := service.GetByID(id, companyID)
	if err != nil {
		c.JSON(http.StatusNotFound, httpx.ErrorResponse{Message: "Process not found"})
		return
	}

	if err := process.Delete(companyID); err != nil {
		c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}

	zap.S().Warnw("Processo eliminato", "id", id, "name", process.ProcessName, "company_id", companyID)
	c.JSON(http.StatusNoContent, nil)
}

var DeleteProcess = d.CompanyRequiredAndPermissionRequired([]string{ProcessDelete}, deleteProcess)

func restoreProcess(c *gin.Context) {
	companyID := httpx.HeaderCompanyID(c)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, httpx.ErrorResponse{Message: "Invalid process ID"})
		return
	}

	service := ProcessService{}
	if err := service.Restore(id, companyID); err != nil {
		c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}

	process, err := service.GetByID(id, companyID)
	if err != nil {
		c.JSON(http.StatusNotFound, httpx.ErrorResponse{Message: "Process not found"})
		return
	}

	zap.S().Infow("Processo ripristinato", "id", id, "name", process.ProcessName, "company_id", companyID)
	c.JSON(http.StatusOK, process)
}

var RestoreProcess = d.CompanyRequiredAndPermissionRequired([]string{ProcessDelete}, restoreProcess)

func getProcessStats(c *gin.Context) {
	companyID := httpx.HeaderCompanyID(c)
	areaIDs := getUserAreaIDs(c)
	service := ProcessService{}
	stats, err := service.GetStats(companyID, areaIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}
	c.JSON(http.StatusOK, stats)
}

var GetProcessStats = d.CompanyRequiredAndPermissionRequired([]string{StatsRead}, getProcessStats)

func recalculateProcess(c *gin.Context) {
	companyID := httpx.HeaderCompanyID(c)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, httpx.ErrorResponse{Message: "Invalid process ID"})
		return
	}

	service := ProcessService{}
	process, err := service.GetByID(id, companyID)
	if err != nil {
		c.JSON(http.StatusNotFound, httpx.ErrorResponse{Message: "Process not found"})
		return
	}

	if err := process.Save(companyID); err != nil {
		c.JSON(http.StatusInternalServerError, httpx.ErrorResponse{Message: httpx.GetErrorMessage(err)})
		return
	}

	zap.S().Infow("Processo ricalcolato", "id", process.ID)
	c.JSON(http.StatusOK, process)
}

var RecalculateProcess = d.CompanyRequiredAndPermissionRequired([]string{ProcessUpdate}, recalculateProcess)

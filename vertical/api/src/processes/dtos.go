package processes

import (
	"fmt"
	"slices"

	"github.com/gin-gonic/gin"
)

type ProcessBody struct {
	ProcessName        string `json:"processName" binding:"required"`
	ProcessDescription string `json:"processDescription"`
	Proposer           string `json:"proposer" binding:"required"`
	Area               string `json:"area" binding:"required"`
	AreaID             *int   `json:"areaId"`
	ResponsibleManager string `json:"responsibleManager" binding:"required"`
	Department         string `json:"department"`

	SystemsInvolved int    `json:"systemsInvolved" binding:"required,min=1"`
	ProcessType     string `json:"processType" binding:"required"`
	Periodicity     string `json:"periodicity" binding:"required"`
	FrequentChanges bool     `json:"frequentChanges"`
	Technology      []string `json:"technology" binding:"required,min=1"`
	TechnologyOther string   `json:"technologyOther"`
	LinkedBots      []string `json:"linkedBots"`
	BotNotes        string   `json:"botNotes"`

	ImplementationCost float64 `json:"implementationCost"`
	TrainingCost       float64 `json:"trainingCost"`
	MaintenanceCost    float64 `json:"maintenanceCost"`

	AnnualSalary       float64 `json:"annualSalary"`
	HourlyCost         float64 `json:"hourlyCost"`
	TimePerActivity    int     `json:"timePerActivity"`
	ActivitiesPerDay   int     `json:"activitiesPerDay"`
	WorkingDaysPerYear int     `json:"workingDaysPerYear"`
	HoursPerDay        float64 `json:"hoursPerDay"`
	DaysPerWeek        int     `json:"daysPerWeek"`
	WeeksPerYear       int     `json:"weeksPerYear"`

	CurrentErrorRate float64 `json:"currentErrorRate"`
	PostErrorRate    float64 `json:"postErrorRate"`
	ErrorCost        float64 `json:"errorCost"`

	ProductivityFactor  float64 `json:"productivityFactor"`
	TimeReductionFactor int     `json:"timeReductionFactor"`

	DataQualityScore        int `json:"dataQualityScore" binding:"required,min=1,max=5"`
	AuditScore              int `json:"auditScore" binding:"required,min=1,max=5"`
	CustomerExperienceScore int `json:"customerExperienceScore" binding:"required,min=1,max=5"`
	ErrorReductionScore     int `json:"errorReductionScore" binding:"required,min=1,max=5"`
	StandardizationScore    int `json:"standardizationScore" binding:"required,min=1,max=5"`
	ScalabilityScore        int `json:"scalabilityScore" binding:"required,min=1,max=5"`
}

func (b *ProcessBody) Bind(c *gin.Context) error {
	return c.ShouldBindJSON(b)
}

type StatusBody struct {
	Status string `json:"status" binding:"required"`
}

func (b *StatusBody) Bind(c *gin.Context) error {
	if err := c.ShouldBindJSON(b); err != nil {
		return err
	}
	if !slices.Contains(ValidStatuses, b.Status) {
		return fmt.Errorf("invalid status: must be one of %v", ValidStatuses)
	}
	return nil
}

type ProcessListParams struct {
	Status  string `form:"status"`
	Search  string `form:"search"`
	Deleted bool   `form:"deleted"`
	AreaIDs []int  `form:"-"`
	Page    int    `form:"page,default=1"`
	Limit   int    `form:"limit,default=25"`
	SortBy  string `form:"sortBy,default=created_at"`
	Order   string `form:"order,default=desc"`
}

type ProcessListResponse struct {
	Data       []Process `json:"data"`
	Total      int       `json:"total"`
	Page       int       `json:"page"`
	Limit      int       `json:"limit"`
	TotalPages int       `json:"totalPages"`
}

type ProcessStats struct {
	Total      int `json:"total"`
	ToValuate  int `json:"toValuate"`
	Analysis   int `json:"analysis"`
	Ongoing    int `json:"ongoing"`
	Production int `json:"production"`
}

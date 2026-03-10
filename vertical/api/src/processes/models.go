package processes

import (
	"encoding/json"
	"time"
)

const (
	ProcessesTable = "processes"

	StatusToValuate  = "To Valuate"
	StatusAnalysis   = "Analysis"
	StatusOngoing    = "Ongoing"
	StatusProduction = "Production"

	ProcessRead   = "processes:process.read"
	ProcessCreate = "processes:process.create"
	ProcessUpdate = "processes:process.update"
	ProcessDelete = "processes:process.delete"
	StatsRead     = "processes:stats.read"
)

var ValidStatuses = []string{StatusToValuate, StatusAnalysis, StatusOngoing, StatusProduction}

// Process is the DB row. Data and Results are stored as JSONB.
type Process struct {
	ID          int             `db:"id" json:"id"`
	CompanyID   int             `db:"company_id" json:"companyId"`
	AreaID      *int            `db:"area_id" json:"areaId"`
	ProcessName string          `db:"process_name" json:"processName"`
	Status      string          `db:"status" json:"status"`
	Data        json.RawMessage `db:"data" json:"data"`
	Results     json.RawMessage `db:"results" json:"results"`
	CreatedBy   *int            `db:"created_by" json:"createdBy"`
	CreatedAt    time.Time       `db:"created_at" json:"createdAt"`
	UpdatedAt    time.Time       `db:"updated_at" json:"updatedAt"`
	DeletedAt    *time.Time      `db:"deleted_at" json:"deletedAt,omitempty"`
	DocumentPath *string         `db:"document_path" json:"documentPath,omitempty"`
	DocumentName *string         `db:"document_name" json:"documentName,omitempty"`
}

// ProcessData holds all user-input fields (stored in data JSONB column).
type ProcessData struct {
	ProcessDescription      string  `json:"processDescription"`
	Proposer                string  `json:"proposer"`
	Area                    string  `json:"area"`
	ResponsibleManager      string  `json:"responsibleManager"`
	Department              string  `json:"department"`
	SystemsInvolved         int     `json:"systemsInvolved"`
	ProcessType             string  `json:"processType"`
	Periodicity             string  `json:"periodicity"`
	FrequentChanges         bool     `json:"frequentChanges"`
	Technology              []string `json:"technology"`
	TechnologyOther         string   `json:"technologyOther,omitempty"`
	LinkedBots              []string `json:"linkedBots,omitempty"`
	BotNotes                string   `json:"botNotes,omitempty"`
	ImplementationCost      float64 `json:"implementationCost"`
	TrainingCost            float64 `json:"trainingCost"`
	MaintenanceCost         float64 `json:"maintenanceCost"`
	HourlyCost              float64 `json:"hourlyCost"`
	TimePerActivity         int     `json:"timePerActivity"`
	ActivitiesPerDay        int     `json:"activitiesPerDay"`
	WorkingDaysPerYear      int     `json:"workingDaysPerYear"`
	HoursPerDay             float64 `json:"hoursPerDay,omitempty"`
	DaysPerWeek             int     `json:"daysPerWeek,omitempty"`
	WeeksPerYear            int     `json:"weeksPerYear,omitempty"`
	CurrentErrorRate        float64 `json:"currentErrorRate"`
	PostErrorRate           float64 `json:"postErrorRate"`
	ErrorCost               float64 `json:"errorCost"`
	ProductivityFactor      float64 `json:"productivityFactor"`
	TimeReductionFactor     int     `json:"timeReductionFactor"`
	DataQualityScore        int     `json:"dataQualityScore"`
	AuditScore              int     `json:"auditScore"`
	CustomerExperienceScore int     `json:"customerExperienceScore"`
	ErrorReductionScore     int     `json:"errorReductionScore"`
	StandardizationScore    int     `json:"standardizationScore"`
	ScalabilityScore        int     `json:"scalabilityScore"`
}

// ProcessResults holds all calculated fields (stored in results JSONB column).
type ProcessResults struct {
	OperationalSavings    float64 `json:"operationalSavings"`
	ErrorReductionSavings float64 `json:"errorReductionSavings"`
	ProductivityBenefit   float64 `json:"productivityBenefit"`
	AnnualSavings         float64 `json:"annualSavings"`
	ROI                   float64 `json:"roi"`
	BreakEvenMonths       *int    `json:"breakEvenMonths"`
	HoursSavedMonthly     float64 `json:"hoursSavedMonthly"`
	HoursSavedAnnually    float64 `json:"hoursSavedAnnually"`
	ImpactScore           float64 `json:"impactScore"`
}

func (p *Process) GetData() (*ProcessData, error) {
	d := &ProcessData{}
	if len(p.Data) == 0 {
		return d, nil
	}
	return d, json.Unmarshal(p.Data, d)
}

func (p *Process) SetData(d *ProcessData) error {
	raw, err := json.Marshal(d)
	if err != nil {
		return err
	}
	p.Data = raw
	return nil
}

func (p *Process) SetResults(r *ProcessResults) error {
	raw, err := json.Marshal(r)
	if err != nil {
		return err
	}
	p.Results = raw
	return nil
}

func (p *Process) ApplyPayload(payload ProcessBody) error {
	p.ProcessName = payload.ProcessName
	p.AreaID = payload.AreaID

	d := &ProcessData{
		ProcessDescription:      payload.ProcessDescription,
		Proposer:                payload.Proposer,
		Area:                    payload.Area,
		ResponsibleManager:      payload.ResponsibleManager,
		Department:              payload.Department,
		SystemsInvolved:         payload.SystemsInvolved,
		ProcessType:             payload.ProcessType,
		Periodicity:             payload.Periodicity,
		FrequentChanges:         payload.FrequentChanges,
		Technology:              payload.Technology,
		TechnologyOther:         payload.TechnologyOther,
		LinkedBots:              payload.LinkedBots,
		BotNotes:                payload.BotNotes,
		ImplementationCost:      payload.ImplementationCost,
		TrainingCost:            payload.TrainingCost,
		MaintenanceCost:         payload.MaintenanceCost,
		HourlyCost:              payload.HourlyCost,
		TimePerActivity:         payload.TimePerActivity,
		ActivitiesPerDay:        payload.ActivitiesPerDay,
		WorkingDaysPerYear:      payload.WorkingDaysPerYear,
		HoursPerDay:             payload.HoursPerDay,
		DaysPerWeek:             payload.DaysPerWeek,
		WeeksPerYear:            payload.WeeksPerYear,
		CurrentErrorRate:        payload.CurrentErrorRate,
		PostErrorRate:           payload.PostErrorRate,
		ErrorCost:               payload.ErrorCost,
		ProductivityFactor:      payload.ProductivityFactor,
		TimeReductionFactor:     payload.TimeReductionFactor,
		DataQualityScore:        payload.DataQualityScore,
		AuditScore:              payload.AuditScore,
		CustomerExperienceScore: payload.CustomerExperienceScore,
		ErrorReductionScore:     payload.ErrorReductionScore,
		StandardizationScore:    payload.StandardizationScore,
		ScalabilityScore:        payload.ScalabilityScore,
	}
	return p.SetData(d)
}

func (p *Process) Save(companyID int) error {
	service := ProcessService{}
	p.CompanyID = companyID
	p.UpdatedAt = time.Now()

	if err := CalculateBenefits(p); err != nil {
		return err
	}

	if p.ID > 0 {
		return service.Update(p)
	}
	p.CreatedAt = time.Now()
	if p.Status == "" {
		p.Status = StatusToValuate
	}
	return service.Insert(p)
}

func (p *Process) Delete(companyID int) error {
	service := ProcessService{}
	return service.Delete(p.ID, companyID)
}

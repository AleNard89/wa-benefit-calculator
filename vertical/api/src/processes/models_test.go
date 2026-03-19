package processes

import (
	"encoding/json"
	"testing"
)

func TestProcess_GetData_Empty(t *testing.T) {
	p := &Process{}
	d, err := p.GetData()
	if err != nil {
		t.Fatalf("GetData on empty process should not fail: %v", err)
	}
	if d.ProcessDescription != "" {
		t.Error("empty data should return zero-value struct")
	}
}

func TestProcess_SetAndGetData(t *testing.T) {
	p := &Process{}
	input := &ProcessData{
		ProcessDescription: "Test process",
		Proposer:           "Mario",
		HourlyCost:         25.5,
		Technology:         []string{"UiPath", "Power Automate"},
		LinkedBots:         []string{"Bot1", "Bot2"},
	}
	if err := p.SetData(input); err != nil {
		t.Fatalf("SetData failed: %v", err)
	}

	got, err := p.GetData()
	if err != nil {
		t.Fatalf("GetData failed: %v", err)
	}
	if got.ProcessDescription != "Test process" {
		t.Errorf("ProcessDescription: got %q, want %q", got.ProcessDescription, "Test process")
	}
	if got.HourlyCost != 25.5 {
		t.Errorf("HourlyCost: got %.2f, want 25.50", got.HourlyCost)
	}
	if len(got.Technology) != 2 {
		t.Errorf("Technology: got %d items, want 2", len(got.Technology))
	}
	if len(got.LinkedBots) != 2 {
		t.Errorf("LinkedBots: got %d items, want 2", len(got.LinkedBots))
	}
}

func TestProcess_SetResults(t *testing.T) {
	p := &Process{}
	r := &ProcessResults{
		OperationalSavings: 1000.50,
		ROI:                250.0,
		ImpactScore:        4.5,
	}
	if err := p.SetResults(r); err != nil {
		t.Fatalf("SetResults failed: %v", err)
	}

	var got ProcessResults
	if err := json.Unmarshal(p.Results, &got); err != nil {
		t.Fatalf("unmarshal results failed: %v", err)
	}
	if got.OperationalSavings != 1000.50 {
		t.Errorf("OperationalSavings: got %.2f, want 1000.50", got.OperationalSavings)
	}
	if got.ROI != 250.0 {
		t.Errorf("ROI: got %.2f, want 250.00", got.ROI)
	}
}

func TestProcess_ApplyPayload(t *testing.T) {
	p := &Process{}
	areaID := 5
	payload := ProcessBody{
		ProcessName:         "Automazione Fatture",
		ProcessDescription:  "Automazione processo fatturazione",
		Proposer:            "Mario Rossi",
		Area:                "Finance",
		AreaID:              &areaID,
		ResponsibleManager:  "Luigi Bianchi",
		SystemsInvolved:     3,
		ProcessType:         "Transactional",
		Periodicity:         "daily",
		Technology:          []string{"UiPath"},
		ImplementationCost:  15000,
		HourlyCost:          22.0,
		TimePerActivity:     45,
		ActivitiesPerDay:    20,
		WorkingDaysPerYear:  220,
		TimeReductionFactor: 75,
		DataQualityScore:    4, AuditScore: 3, CustomerExperienceScore: 4,
		ErrorReductionScore: 5, StandardizationScore: 3, ScalabilityScore: 4,
	}

	if err := p.ApplyPayload(payload); err != nil {
		t.Fatalf("ApplyPayload failed: %v", err)
	}

	if p.ProcessName != "Automazione Fatture" {
		t.Errorf("ProcessName: got %q, want %q", p.ProcessName, "Automazione Fatture")
	}
	if p.AreaID == nil || *p.AreaID != 5 {
		t.Errorf("AreaID: got %v, want 5", p.AreaID)
	}

	d, _ := p.GetData()
	if d.Proposer != "Mario Rossi" {
		t.Errorf("Proposer: got %q, want %q", d.Proposer, "Mario Rossi")
	}
	if d.ImplementationCost != 15000 {
		t.Errorf("ImplementationCost: got %.0f, want 15000", d.ImplementationCost)
	}
	if len(d.Technology) != 1 || d.Technology[0] != "UiPath" {
		t.Errorf("Technology: got %v, want [UiPath]", d.Technology)
	}
}

func TestStatusBody_ValidStatuses(t *testing.T) {
	for _, s := range ValidStatuses {
		if s == "" {
			t.Error("valid status should not be empty")
		}
	}
	if len(ValidStatuses) != 4 {
		t.Errorf("expected 4 valid statuses, got %d", len(ValidStatuses))
	}
}

func TestProcess_GetData_InvalidJSON(t *testing.T) {
	p := &Process{Data: json.RawMessage(`{invalid`)}
	_, err := p.GetData()
	if err == nil {
		t.Error("GetData should fail on invalid JSON")
	}
}

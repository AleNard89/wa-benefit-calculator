package processes

import (
	"encoding/json"
	"math"
	"testing"
)

func makeProcess(data ProcessData) *Process {
	raw, _ := json.Marshal(data)
	return &Process{Data: raw}
}

func getResults(t *testing.T, p *Process) *ProcessResults {
	t.Helper()
	r := &ProcessResults{}
	if err := json.Unmarshal(p.Results, r); err != nil {
		t.Fatalf("failed to unmarshal results: %v", err)
	}
	return r
}

func assertClose(t *testing.T, name string, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > 0.02 {
		t.Errorf("%s: got %.2f, want %.2f", name, got, want)
	}
}

func TestCalculateBenefits_BasicScenario(t *testing.T) {
	p := makeProcess(ProcessData{
		TimePerActivity:    30,
		ActivitiesPerDay:   10,
		WorkingDaysPerYear: 220,
		TimeReductionFactor: 80,
		HourlyCost:          25.0,
		CurrentErrorRate:    5.0,
		ErrorCost:           50.0,
		ProductivityFactor:  1.2,
		ImplementationCost:  10000,
		TrainingCost:        2000,
		MaintenanceCost:     1000,
		DataQualityScore:    4, AuditScore: 3, CustomerExperienceScore: 5,
		ErrorReductionScore: 4, StandardizationScore: 3, ScalabilityScore: 4,
	})

	if err := CalculateBenefits(p); err != nil {
		t.Fatalf("CalculateBenefits returned error: %v", err)
	}

	r := getResults(t, p)

	// hours saved: (30/60) * 10 * 220 * 0.8 = 880
	assertClose(t, "HoursSavedAnnually", r.HoursSavedAnnually, 880.0)
	assertClose(t, "HoursSavedMonthly", r.HoursSavedMonthly, 880.0/12.0)

	// operational savings: (0.5 * 10 * 220 * 25) * 0.8 = 22000
	assertClose(t, "OperationalSavings", r.OperationalSavings, 22000.0)

	// error reduction: 10 * 220 * 0.05 * 50 = 5500
	assertClose(t, "ErrorReductionSavings", r.ErrorReductionSavings, 5500.0)

	// productivity: (12 - 10) * 220 * (25 * 0.5) = 5500
	assertClose(t, "ProductivityBenefit", r.ProductivityBenefit, 5500.0)

	// annual savings = 22000 + 5500 + 5500 = 33000
	assertClose(t, "AnnualSavings", r.AnnualSavings, 33000.0)

	// net = 33000 - 1000 = 32000 => ROI = (32000/10000)*100 = 320
	assertClose(t, "ROI", r.ROI, 320.0)

	// break-even: ceil(13000 / (32000/12)) = ceil(4.875) = 5
	if r.BreakEvenMonths == nil || *r.BreakEvenMonths != 5 {
		t.Errorf("BreakEvenMonths: got %v, want 5", r.BreakEvenMonths)
	}

	// impact: (4+3+5+4+3+4)/6 = 3.833...
	assertClose(t, "ImpactScore", r.ImpactScore, 3.83)
}

func TestCalculateBenefits_ZeroImplementationCost(t *testing.T) {
	p := makeProcess(ProcessData{
		TimePerActivity:     60,
		ActivitiesPerDay:    5,
		WorkingDaysPerYear:  200,
		TimeReductionFactor: 50,
		HourlyCost:          30.0,
		ImplementationCost:  0,
		DataQualityScore:    3, AuditScore: 3, CustomerExperienceScore: 3,
		ErrorReductionScore: 3, StandardizationScore: 3, ScalabilityScore: 3,
	})

	if err := CalculateBenefits(p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	r := getResults(t, p)

	if r.ROI != 0 {
		t.Errorf("ROI should be 0 when implementation cost is 0, got %.2f", r.ROI)
	}
}

func TestCalculateBenefits_NegativeNetBenefit_NoBreakEven(t *testing.T) {
	p := makeProcess(ProcessData{
		TimePerActivity:     1,
		ActivitiesPerDay:    1,
		WorkingDaysPerYear:  10,
		TimeReductionFactor: 10,
		HourlyCost:          10.0,
		ImplementationCost:  100000,
		MaintenanceCost:     999999,
		DataQualityScore:    1, AuditScore: 1, CustomerExperienceScore: 1,
		ErrorReductionScore: 1, StandardizationScore: 1, ScalabilityScore: 1,
	})

	if err := CalculateBenefits(p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	r := getResults(t, p)

	if r.BreakEvenMonths != nil {
		t.Errorf("BreakEvenMonths should be nil when net benefit <= 0, got %d", *r.BreakEvenMonths)
	}
}

func TestCalculateBenefits_AllScoresMax(t *testing.T) {
	p := makeProcess(ProcessData{
		TimePerActivity: 10, ActivitiesPerDay: 1, WorkingDaysPerYear: 1,
		TimeReductionFactor: 100, HourlyCost: 1,
		DataQualityScore: 5, AuditScore: 5, CustomerExperienceScore: 5,
		ErrorReductionScore: 5, StandardizationScore: 5, ScalabilityScore: 5,
	})
	CalculateBenefits(p)
	r := getResults(t, p)
	assertClose(t, "ImpactScore", r.ImpactScore, 5.0)
}

func TestCalculateBenefits_AllScoresMin(t *testing.T) {
	p := makeProcess(ProcessData{
		TimePerActivity: 10, ActivitiesPerDay: 1, WorkingDaysPerYear: 1,
		TimeReductionFactor: 100, HourlyCost: 1,
		DataQualityScore: 1, AuditScore: 1, CustomerExperienceScore: 1,
		ErrorReductionScore: 1, StandardizationScore: 1, ScalabilityScore: 1,
	})
	CalculateBenefits(p)
	r := getResults(t, p)
	assertClose(t, "ImpactScore", r.ImpactScore, 1.0)
}

func TestCalculateBenefits_EmptyData(t *testing.T) {
	p := &Process{}
	if err := CalculateBenefits(p); err != nil {
		t.Fatalf("should handle empty data without error, got: %v", err)
	}
}

func TestCalculateBenefits_100PercentReduction(t *testing.T) {
	p := makeProcess(ProcessData{
		TimePerActivity:     60,
		ActivitiesPerDay:    8,
		WorkingDaysPerYear:  250,
		TimeReductionFactor: 100,
		HourlyCost:          20.0,
		DataQualityScore:    3, AuditScore: 3, CustomerExperienceScore: 3,
		ErrorReductionScore: 3, StandardizationScore: 3, ScalabilityScore: 3,
	})
	CalculateBenefits(p)
	r := getResults(t, p)

	// 100% reduction => all hours saved: 1h * 8 * 250 = 2000
	assertClose(t, "HoursSavedAnnually", r.HoursSavedAnnually, 2000.0)
	// operational savings = 1*8*250*20 = 40000
	assertClose(t, "OperationalSavings", r.OperationalSavings, 40000.0)
}

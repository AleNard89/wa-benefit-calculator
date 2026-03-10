package processes

import "math"

func CalculateBenefits(p *Process) error {
	d, err := p.GetData()
	if err != nil {
		return err
	}

	timePerActivityHours := float64(d.TimePerActivity) / 60.0
	activitiesPerDay := float64(d.ActivitiesPerDay)
	workingDays := float64(d.WorkingDaysPerYear)
	reductionFactor := float64(d.TimeReductionFactor) / 100.0

	hoursSavedAnnually := timePerActivityHours * activitiesPerDay * workingDays * reductionFactor
	hoursSavedMonthly := hoursSavedAnnually / 12.0

	annualCostBefore := timePerActivityHours * activitiesPerDay * workingDays * d.HourlyCost
	annualCostAfter := annualCostBefore * (1.0 - reductionFactor)
	operationalSavings := annualCostBefore - annualCostAfter

	errorRateDiff := (d.CurrentErrorRate - d.PostErrorRate) / 100.0
	errorReductionSavings := activitiesPerDay * workingDays * errorRateDiff * d.ErrorCost

	activitiesAfter := activitiesPerDay * d.ProductivityFactor
	activityValue := d.HourlyCost * timePerActivityHours
	productivityBenefit := (activitiesAfter - activitiesPerDay) * workingDays * activityValue

	annualSavings := operationalSavings + errorReductionSavings + productivityBenefit
	netAnnualBenefit := annualSavings - d.MaintenanceCost

	var breakEvenMonths *int
	firstYearCost := d.ImplementationCost + d.TrainingCost + d.MaintenanceCost
	if netAnnualBenefit > 0 {
		m := int(math.Ceil(firstYearCost / (netAnnualBenefit / 12.0)))
		breakEvenMonths = &m
	}

	roi := 0.0
	if d.ImplementationCost > 0 {
		roi = (netAnnualBenefit / d.ImplementationCost) * 100.0
	}

	totalScore := d.DataQualityScore + d.AuditScore + d.CustomerExperienceScore +
		d.ErrorReductionScore + d.StandardizationScore + d.ScalabilityScore
	impactScore := float64(totalScore) / 6.0

	r := &ProcessResults{
		OperationalSavings:    math.Round(operationalSavings*100) / 100,
		ErrorReductionSavings: math.Round(errorReductionSavings*100) / 100,
		ProductivityBenefit:   math.Round(productivityBenefit*100) / 100,
		AnnualSavings:         math.Round(annualSavings*100) / 100,
		ROI:                   math.Round(roi*100) / 100,
		BreakEvenMonths:       breakEvenMonths,
		HoursSavedMonthly:     math.Round(hoursSavedMonthly*100) / 100,
		HoursSavedAnnually:    math.Round(hoursSavedAnnually*100) / 100,
		ImpactScore:           math.Round(impactScore*100) / 100,
	}

	return p.SetResults(r)
}

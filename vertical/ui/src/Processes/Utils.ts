import type { ProcessInput, ProcessResults } from './Types'

export function calculateBenefits(input: ProcessInput): ProcessResults {
  const timePerActivityHours = input.timePerActivity / 60
  const reductionFactor = input.timeReductionFactor / 100

  // Ore risparmiate
  const hoursSavedAnnually =
    timePerActivityHours * input.activitiesPerDay * input.workingDaysPerYear * reductionFactor
  const hoursSavedMonthly = hoursSavedAnnually / 12

  // Risparmio operativo
  const annualCostBefore =
    timePerActivityHours * input.activitiesPerDay * input.workingDaysPerYear * input.hourlyCost
  const annualCostAfter = annualCostBefore * (1 - reductionFactor)
  const operationalSavings = annualCostBefore - annualCostAfter

  // Risparmio riduzione errori
  const errorRateDiff = (input.currentErrorRate - input.postErrorRate) / 100
  const errorReductionSavings =
    input.activitiesPerDay * input.workingDaysPerYear * errorRateDiff * input.errorCost

  // Beneficio produttivita
  const activitiesAfter = input.activitiesPerDay * input.productivityFactor
  const activityValue = input.hourlyCost * timePerActivityHours
  const productivityBenefit =
    (activitiesAfter - input.activitiesPerDay) * input.workingDaysPerYear * activityValue

  // Totali
  const annualSavings = operationalSavings + errorReductionSavings + productivityBenefit
  const netAnnualBenefit = annualSavings - input.maintenanceCost

  // Break-even
  const firstYearCost = input.implementationCost + input.trainingCost + input.maintenanceCost
  let breakEvenMonths: number | null = null
  if (netAnnualBenefit > 0) {
    breakEvenMonths = Math.ceil(firstYearCost / (netAnnualBenefit / 12))
  }

  // ROI
  const roi = input.implementationCost > 0
    ? (netAnnualBenefit / input.implementationCost) * 100
    : 0

  // Impact score
  const impactScore = (
    input.dataQualityScore +
    input.auditScore +
    input.customerExperienceScore +
    input.errorReductionScore +
    input.standardizationScore +
    input.scalabilityScore
  ) / 6

  return {
    operationalSavings: Math.round(operationalSavings * 100) / 100,
    errorReductionSavings: Math.round(errorReductionSavings * 100) / 100,
    productivityBenefit: Math.round(productivityBenefit * 100) / 100,
    annualSavings: Math.round(annualSavings * 100) / 100,
    roi: Math.round(roi * 100) / 100,
    breakEvenMonths,
    hoursSavedMonthly: Math.round(hoursSavedMonthly * 100) / 100,
    hoursSavedAnnually: Math.round(hoursSavedAnnually * 100) / 100,
    impactScore: Math.round(impactScore * 100) / 100,
  }
}

export function formatCurrency(value: number): string {
  return new Intl.NumberFormat('it-IT', { style: 'currency', currency: 'EUR' }).format(value)
}

export function formatNumber(value: number, decimals = 1): string {
  return new Intl.NumberFormat('it-IT', { maximumFractionDigits: decimals }).format(value)
}

export function formatPercent(value: number): string {
  return new Intl.NumberFormat('it-IT', { maximumFractionDigits: 1 }).format(value) + '%'
}

export const statusColorMap: Record<string, string> = {
  'To Valuate': 'gray',
  'Analysis': 'blue',
  'Ongoing': 'orange',
  'Production': 'green',
}

export const statusLabelMap: Record<string, string> = {
  'To Valuate': 'Da Valutare',
  'Analysis': 'In Analisi',
  'Ongoing': 'In Corso',
  'Production': 'Produzione',
}

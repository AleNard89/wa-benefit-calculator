import { describe, expect, it } from 'vitest'
import {
  calculateBenefits,
  formatCurrency,
  formatNumber,
  formatPercent,
  statusColorMap,
  statusLabelMap,
} from './Utils'
import type { ProcessInput } from './Types'

const baseInput: ProcessInput = {
  processName: 'Test',
  areaId: null,
  processDescription: '',
  proposer: '',
  area: '',
  responsibleManager: '',
  department: '',
  systemsInvolved: 1,
  processType: '',
  periodicity: '',
  frequentChanges: false,
  technology: [],
  implementationCost: 10000,
  trainingCost: 2000,
  maintenanceCost: 1000,
  annualSalary: 0,
  hourlyCost: 25,
  timePerActivity: 30,
  activitiesPerDay: 10,
  workingDaysPerYear: 220,
  hoursPerDay: 8,
  daysPerWeek: 5,
  weeksPerYear: 44,
  currentErrorRate: 5,
  postErrorRate: 0,
  errorCost: 50,
  productivityFactor: 1.2,
  timeReductionFactor: 80,
  dataQualityScore: 4,
  auditScore: 3,
  customerExperienceScore: 5,
  errorReductionScore: 4,
  standardizationScore: 3,
  scalabilityScore: 4,
}

describe('calculateBenefits', () => {
  it('calculates hours saved correctly', () => {
    const r = calculateBenefits(baseInput)
    // (30/60) * 10 * 220 * 0.8 = 880
    expect(r.hoursSavedAnnually).toBe(880)
    expect(r.hoursSavedMonthly).toBeCloseTo(880 / 12, 1)
  })

  it('calculates operational savings', () => {
    const r = calculateBenefits(baseInput)
    // (0.5 * 10 * 220 * 25) * 0.8 = 22000
    expect(r.operationalSavings).toBe(22000)
  })

  it('calculates error reduction savings', () => {
    const r = calculateBenefits(baseInput)
    // 10 * 220 * ((5-0)/100) * 50 = 5500
    expect(r.errorReductionSavings).toBe(5500)
  })

  it('calculates productivity benefit', () => {
    const r = calculateBenefits(baseInput)
    // (12 - 10) * 220 * (25 * 0.5) = 5500
    expect(r.productivityBenefit).toBe(5500)
  })

  it('calculates annual savings as sum of components', () => {
    const r = calculateBenefits(baseInput)
    expect(r.annualSavings).toBe(
      r.operationalSavings + r.errorReductionSavings + r.productivityBenefit
    )
  })

  it('calculates ROI correctly', () => {
    const r = calculateBenefits(baseInput)
    // net = 33000 - 1000 = 32000, ROI = (32000/10000)*100 = 320
    expect(r.roi).toBe(320)
  })

  it('calculates break-even months', () => {
    const r = calculateBenefits(baseInput)
    // first year cost = 13000, monthly net = 32000/12 ≈ 2666.67, ceil(13000/2666.67) = 5
    expect(r.breakEvenMonths).toBe(5)
  })

  it('calculates impact score as average of 6 scores', () => {
    const r = calculateBenefits(baseInput)
    // (4+3+5+4+3+4)/6 = 3.83
    expect(r.impactScore).toBeCloseTo(3.83, 1)
  })

  it('returns ROI 0 when implementation cost is 0', () => {
    const r = calculateBenefits({ ...baseInput, implementationCost: 0 })
    expect(r.roi).toBe(0)
  })

  it('returns null break-even when net benefit <= 0', () => {
    const r = calculateBenefits({ ...baseInput, maintenanceCost: 999999 })
    expect(r.breakEvenMonths).toBeNull()
  })

  it('handles 100% time reduction', () => {
    const r = calculateBenefits({ ...baseInput, timeReductionFactor: 100 })
    // (0.5 * 10 * 220) * 1.0 = 1100 hours
    expect(r.hoursSavedAnnually).toBe(1100)
  })

  it('handles 0% time reduction', () => {
    const r = calculateBenefits({ ...baseInput, timeReductionFactor: 0 })
    expect(r.hoursSavedAnnually).toBe(0)
    expect(r.operationalSavings).toBe(0)
  })
})

describe('formatCurrency', () => {
  it('formats positive EUR amount with euro sign', () => {
    const result = formatCurrency(1234.56)
    expect(result).toContain('€')
    expect(result).toContain('1234')
  })

  it('formats zero', () => {
    const result = formatCurrency(0)
    expect(result).toContain('0')
    expect(result).toContain('€')
  })
})

describe('formatNumber', () => {
  it('formats with default decimals', () => {
    const result = formatNumber(1234.567)
    expect(result).toContain('1234')
    expect(result).toContain('6')
  })

  it('formats with custom decimals', () => {
    const result = formatNumber(1234.567, 2)
    expect(result).toContain('1234')
    expect(result).toContain('57')
  })
})

describe('formatPercent', () => {
  it('appends % sign', () => {
    expect(formatPercent(85.5)).toBe('85,5%')
  })
})

describe('statusColorMap', () => {
  it('maps all 4 statuses', () => {
    expect(Object.keys(statusColorMap)).toHaveLength(4)
    expect(statusColorMap['To Valuate']).toBe('gray')
    expect(statusColorMap['Production']).toBe('green')
  })
})

describe('statusLabelMap', () => {
  it('maps statuses to Italian labels', () => {
    expect(statusLabelMap['To Valuate']).toBe('Da Valutare')
    expect(statusLabelMap['Analysis']).toBe('In Analisi')
    expect(statusLabelMap['Ongoing']).toBe('In Corso')
    expect(statusLabelMap['Production']).toBe('Produzione')
  })
})

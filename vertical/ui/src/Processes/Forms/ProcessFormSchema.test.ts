import { describe, expect, it } from 'vitest'
import { processFormSchema, processFormDefaults } from './ProcessFormSchema'

describe('processFormSchema', () => {
  it('accepts valid defaults', () => {
    const valid = {
      ...processFormDefaults,
      processName: 'Test Process',
      proposer: 'Mario Rossi',
      area: 'Finance',
      responsibleManager: 'Luigi Bianchi',
      processType: 'Transactional',
      periodicity: 'daily',
      technology: ['UiPath'],
    }
    const result = processFormSchema.safeParse(valid)
    expect(result.success).toBe(true)
  })

  it('rejects empty processName', () => {
    const invalid = {
      ...processFormDefaults,
      processName: '',
      proposer: 'X',
      area: 'Y',
      responsibleManager: 'Z',
      processType: 'T',
      periodicity: 'daily',
      technology: ['UiPath'],
    }
    const result = processFormSchema.safeParse(invalid)
    expect(result.success).toBe(false)
  })

  it('rejects empty technology array', () => {
    const invalid = {
      ...processFormDefaults,
      processName: 'Test',
      proposer: 'X',
      area: 'Y',
      responsibleManager: 'Z',
      processType: 'T',
      periodicity: 'daily',
      technology: [],
    }
    const result = processFormSchema.safeParse(invalid)
    expect(result.success).toBe(false)
  })

  it('rejects scores outside 1-5 range', () => {
    const tooLow = {
      ...processFormDefaults,
      processName: 'Test',
      proposer: 'X',
      area: 'Y',
      responsibleManager: 'Z',
      processType: 'T',
      periodicity: 'daily',
      technology: ['UiPath'],
      dataQualityScore: 0,
    }
    expect(processFormSchema.safeParse(tooLow).success).toBe(false)

    const tooHigh = { ...tooLow, dataQualityScore: 6 }
    expect(processFormSchema.safeParse(tooHigh).success).toBe(false)
  })

  it('rejects negative costs', () => {
    const invalid = {
      ...processFormDefaults,
      processName: 'Test',
      proposer: 'X',
      area: 'Y',
      responsibleManager: 'Z',
      processType: 'T',
      periodicity: 'daily',
      technology: ['UiPath'],
      implementationCost: -100,
    }
    expect(processFormSchema.safeParse(invalid).success).toBe(false)
  })

  it('rejects currentErrorRate > 100', () => {
    const invalid = {
      ...processFormDefaults,
      processName: 'Test',
      proposer: 'X',
      area: 'Y',
      responsibleManager: 'Z',
      processType: 'T',
      periodicity: 'daily',
      technology: ['UiPath'],
      currentErrorRate: 150,
    }
    expect(processFormSchema.safeParse(invalid).success).toBe(false)
  })

  it('accepts daysPerWeek up to 365', () => {
    const valid = {
      ...processFormDefaults,
      processName: 'Test',
      proposer: 'X',
      area: 'Y',
      responsibleManager: 'Z',
      processType: 'T',
      periodicity: 'annual',
      technology: ['UiPath'],
      daysPerWeek: 365,
    }
    expect(processFormSchema.safeParse(valid).success).toBe(true)
  })

  it('accepts nullable areaId', () => {
    const valid = {
      ...processFormDefaults,
      processName: 'Test',
      proposer: 'X',
      area: 'Y',
      responsibleManager: 'Z',
      processType: 'T',
      periodicity: 'daily',
      technology: ['UiPath'],
      areaId: null,
    }
    expect(processFormSchema.safeParse(valid).success).toBe(true)
  })
})

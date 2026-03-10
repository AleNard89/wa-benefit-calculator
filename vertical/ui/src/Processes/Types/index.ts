export type ProcessStatus = 'To Valuate' | 'Analysis' | 'Ongoing' | 'Production'

export interface ProcessData {
  processDescription: string
  proposer: string
  area: string
  responsibleManager: string
  department: string
  systemsInvolved: number
  processType: string
  periodicity: string
  frequentChanges: boolean
  technology: string[] | string
  technologyOther?: string
  linkedBots?: string[]
  botNotes?: string
  implementationCost: number
  trainingCost: number
  maintenanceCost: number
  hourlyCost: number
  timePerActivity: number
  activitiesPerDay: number
  workingDaysPerYear: number
  hoursPerDay?: number
  daysPerWeek?: number
  weeksPerYear?: number
  currentErrorRate: number
  postErrorRate: number
  errorCost: number
  productivityFactor: number
  timeReductionFactor: number
  dataQualityScore: number
  auditScore: number
  customerExperienceScore: number
  errorReductionScore: number
  standardizationScore: number
  scalabilityScore: number
}

export interface ProcessResults {
  operationalSavings: number
  errorReductionSavings: number
  productivityBenefit: number
  annualSavings: number
  roi: number
  breakEvenMonths: number | null
  hoursSavedMonthly: number
  hoursSavedAnnually: number
  impactScore: number
}

export interface Process {
  id: number
  companyId: number
  areaId: number | null
  processName: string
  status: ProcessStatus
  data: ProcessData
  results: ProcessResults
  createdBy: number | null
  createdAt: string
  updatedAt: string
  documentPath?: string | null
  documentName?: string | null
}

export interface ProcessInput {
  processName: string
  areaId: number | null
  processDescription: string
  proposer: string
  area: string
  responsibleManager: string
  department: string
  systemsInvolved: number
  processType: string
  periodicity: string
  frequentChanges: boolean
  technology: string[]
  technologyOther?: string
  linkedBots?: string[]
  botNotes?: string
  implementationCost: number
  trainingCost: number
  maintenanceCost: number
  hourlyCost: number
  timePerActivity: number
  activitiesPerDay: number
  workingDaysPerYear: number
  hoursPerDay: number
  daysPerWeek: number
  weeksPerYear: number
  currentErrorRate: number
  postErrorRate: number
  errorCost: number
  productivityFactor: number
  timeReductionFactor: number
  dataQualityScore: number
  auditScore: number
  customerExperienceScore: number
  errorReductionScore: number
  standardizationScore: number
  scalabilityScore: number
}

export interface ProcessListResponse {
  data: Process[]
  total: number
  page: number
  limit: number
  totalPages: number
}

export interface ProcessStats {
  total: number
  toValuate: number
  analysis: number
  ongoing: number
  production: number
}

export interface ProcessListParams {
  status?: string
  search?: string
  deleted?: boolean
  page?: number
  limit?: number
  sortBy?: string
  order?: string
}

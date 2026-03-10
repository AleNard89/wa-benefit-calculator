export interface UiPathFolder {
  id: string
  name: string
}

export interface ConnectorResponse {
  id: number
  companyId: number
  name: string
  type: 'UIPATH' | 'PYTHON_AGENT'
  isActive: boolean
  organizationName: string
  tenantName: string
  folderId: string
  folderName: string
  folders: UiPathFolder[]
  createdAt: string
  updatedAt: string
}

export interface ConnectorForm {
  name: string
  type: string
  organizationName: string
  tenantName: string
  accessToken: string
  folderId: string
  folderName: string
  folders: UiPathFolder[]
  isActive?: boolean
}

export interface JobExecution {
  id: number
  companyId: number
  connectorId: number
  externalJobKey?: string
  externalJobId?: number
  processName?: string
  state: string
  sourceType?: string
  source?: string
  startTime?: string
  endTime?: string
  hostMachine?: string
  folderName?: string
  info?: string
  details?: Record<string, unknown>
  createdAt: string
  updatedAt: string
}

export interface QueueItem {
  id: number
  companyId: number
  connectorId: number
  externalItemKey?: string
  externalItemId?: number
  queueDefinitionId?: number
  queueName: string
  status: string
  priority?: string
  reference?: string
  processingExceptionType?: string
  errorMessage?: string
  startProcessing?: string
  endProcessing?: string
  retryNumber: number
  folderName?: string
  specificContent?: Record<string, unknown>
  createdAt: string
  updatedAt: string
}

export interface QueueDefinition {
  id: number
  companyId: number
  connectorId: number
  externalDefinitionId?: number
  name: string
  maxRetries: number
  folderName?: string
  createdAt: string
  updatedAt: string
}

export interface PaginatedResponse<T> {
  data: T[]
  total: number
  page: number
  limit: number
  totalPages: number
}

export interface ProcessSchedule {
  id: number
  companyId: number
  connectorId: number
  externalScheduleId?: number
  name: string
  enabled: boolean
  releaseName?: string
  packageName?: string
  cronExpression?: string
  cronSummary?: string
  nextOccurrence?: string
  timezoneId?: string
  timezoneIana?: string
  startStrategy: number
  folderName?: string
  inputArguments?: Record<string, unknown>
  createdAt: string
  updatedAt: string
}

export interface ScheduleFilters {
  enabled?: boolean
  connectorId?: number
  processNames?: string
  page?: number
  limit?: number
}

export interface JobFilters {
  state?: string
  connectorId?: number
  processNames?: string
  page?: number
  limit?: number
}

export interface ProcessQueueMap {
  id: number
  companyId: number
  connectorId: number
  processName: string
  queueName: string
  autoDetected: boolean
  createdAt: string
  updatedAt: string
}

export interface ProcessQueueMapForm {
  connectorId: number
  processName: string
  queueName: string
}

export interface QueueItemFilters {
  status?: string
  connectorId?: number
  queueName?: string
  processNames?: string
  page?: number
  limit?: number
}

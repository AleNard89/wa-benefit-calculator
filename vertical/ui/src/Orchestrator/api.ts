import { api } from '@/Core/Services/Api'

import type {
  ConnectorForm,
  ConnectorResponse,
  JobExecution,
  JobFilters,
  PaginatedResponse,
  ProcessQueueMap,
  ProcessQueueMapForm,
  ProcessSchedule,
  QueueDefinition,
  QueueItem,
  QueueItemFilters,
  ScheduleFilters,
} from './types'

const prefix = 'orchestrator'

const extendedApi = api.injectEndpoints({
  endpoints: (builder) => ({
    // Connectors
    connectors: builder.query<ConnectorResponse[], void>({
      query: () => `${prefix}/connectors`,
      providesTags: ['Connectors', 'HasCompanyHeader'],
    }),
    createConnector: builder.mutation<ConnectorResponse, ConnectorForm>({
      query: (body) => ({ url: `${prefix}/connectors`, method: 'POST', body }),
      invalidatesTags: ['Connectors'],
    }),
    updateConnector: builder.mutation<ConnectorResponse, { id: number; body: ConnectorForm }>({
      query: ({ id, body }) => ({ url: `${prefix}/connectors/${id}`, method: 'PUT', body }),
      invalidatesTags: ['Connectors'],
    }),
    testConnector: builder.mutation<{ message: string }, number>({
      query: (id) => ({ url: `${prefix}/connectors/${id}/test`, method: 'POST' }),
    }),
    syncConnector: builder.mutation<{ message: string }, number>({
      query: (id) => ({ url: `${prefix}/connectors/${id}/sync`, method: 'POST' }),
      invalidatesTags: ['Jobs', 'QueueItems', 'QueueDefinitions', 'Schedules', 'ProcessNames'],
    }),

    // Process Names (prefixes)
    processNames: builder.query<string[], void>({
      query: () => `${prefix}/process-names`,
      providesTags: ['ProcessNames', 'HasCompanyHeader'],
    }),

    // Bot Names (full names from Orchestrator)
    botNames: builder.query<string[], void>({
      query: () => `${prefix}/bot-names`,
      providesTags: ['ProcessNames', 'HasCompanyHeader'],
    }),

    // Jobs
    jobs: builder.query<PaginatedResponse<JobExecution>, JobFilters>({
      query: (params) => {
        const searchParams = new URLSearchParams()
        if (params.state) searchParams.set('state', params.state)
        if (params.connectorId) searchParams.set('connectorId', params.connectorId.toString())
        if (params.processNames) searchParams.set('processNames', params.processNames)
        if (params.page) searchParams.set('page', params.page.toString())
        if (params.limit) searchParams.set('limit', params.limit.toString())
        const qs = searchParams.toString()
        return `${prefix}/jobs${qs ? `?${qs}` : ''}`
      },
      providesTags: ['Jobs', 'HasCompanyHeader'],
    }),
    job: builder.query<JobExecution, number>({
      query: (id) => `${prefix}/jobs/${id}`,
      providesTags: ['Jobs', 'HasCompanyHeader'],
    }),

    // Queue Items
    queueItems: builder.query<PaginatedResponse<QueueItem>, QueueItemFilters>({
      query: (params) => {
        const searchParams = new URLSearchParams()
        if (params.status) searchParams.set('status', params.status)
        if (params.connectorId) searchParams.set('connectorId', params.connectorId.toString())
        if (params.queueName) searchParams.set('queueName', params.queueName)
        if (params.processNames) searchParams.set('processNames', params.processNames)
        if (params.page) searchParams.set('page', params.page.toString())
        if (params.limit) searchParams.set('limit', params.limit.toString())
        const qs = searchParams.toString()
        return `${prefix}/queue-items${qs ? `?${qs}` : ''}`
      },
      providesTags: ['QueueItems', 'HasCompanyHeader'],
    }),

    // Schedules
    schedules: builder.query<PaginatedResponse<ProcessSchedule>, ScheduleFilters>({
      query: (params) => {
        const searchParams = new URLSearchParams()
        if (params.enabled !== undefined) searchParams.set('enabled', params.enabled.toString())
        if (params.connectorId) searchParams.set('connectorId', params.connectorId.toString())
        if (params.processNames) searchParams.set('processNames', params.processNames)
        if (params.page) searchParams.set('page', params.page.toString())
        if (params.limit) searchParams.set('limit', params.limit.toString())
        const qs = searchParams.toString()
        return `${prefix}/schedules${qs ? `?${qs}` : ''}`
      },
      providesTags: ['Schedules', 'HasCompanyHeader'],
    }),

    // Process-Queue Map
    processQueueMaps: builder.query<ProcessQueueMap[], void>({
      query: () => `${prefix}/process-queue-map`,
      providesTags: ['ProcessQueueMap', 'HasCompanyHeader'],
    }),
    createProcessQueueMap: builder.mutation<ProcessQueueMap, ProcessQueueMapForm>({
      query: (body) => ({ url: `${prefix}/process-queue-map`, method: 'POST', body }),
      invalidatesTags: ['ProcessQueueMap'],
    }),
    deleteProcessQueueMap: builder.mutation<{ message: string }, number>({
      query: (id) => ({ url: `${prefix}/process-queue-map/${id}`, method: 'DELETE' }),
      invalidatesTags: ['ProcessQueueMap'],
    }),

    // Queue Definitions
    queueDefinitions: builder.query<QueueDefinition[], { connectorId?: number }>({
      query: (params) => {
        const qs = params.connectorId ? `?connectorId=${params.connectorId}` : ''
        return `${prefix}/queue-definitions${qs}`
      },
      providesTags: ['QueueDefinitions', 'HasCompanyHeader'],
    }),
  }),
  overrideExisting: false,
})

export const {
  useConnectorsQuery,
  useCreateConnectorMutation,
  useUpdateConnectorMutation,
  useTestConnectorMutation,
  useSyncConnectorMutation,
  useProcessNamesQuery,
  useBotNamesQuery,
  useJobsQuery,
  useJobQuery,
  useSchedulesQuery,
  useQueueItemsQuery,
  useQueueDefinitionsQuery,
  useProcessQueueMapsQuery,
  useCreateProcessQueueMapMutation,
  useDeleteProcessQueueMapMutation,
} = extendedApi

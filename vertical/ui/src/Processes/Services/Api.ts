import { api } from '@/Core/Services/Api'

import type { Process, ProcessInput, ProcessListParams, ProcessListResponse, ProcessStats } from '../Types'

const prefix = 'processes'

function buildQueryString(params: ProcessListParams): string {
  const searchParams = new URLSearchParams()
  if (params.status) searchParams.set('status', params.status)
  if (params.search) searchParams.set('search', params.search)
  if (params.deleted) searchParams.set('deleted', 'true')
  if (params.page) searchParams.set('page', params.page.toString())
  if (params.limit) searchParams.set('limit', params.limit.toString())
  if (params.sortBy) searchParams.set('sortBy', params.sortBy)
  if (params.order) searchParams.set('order', params.order)
  const qs = searchParams.toString()
  return qs ? `?${qs}` : ''
}

const extendedApi = api.injectEndpoints({
  endpoints: (builder) => ({
    processes: builder.query<ProcessListResponse, ProcessListParams>({
      query: (params) => `${prefix}${buildQueryString(params)}`,
      providesTags: ['Processes', 'HasCompanyHeader'],
    }),
    process: builder.query<Process, number>({
      query: (id) => `${prefix}/${id}`,
      providesTags: ['Processes'],
    }),
    processStats: builder.query<ProcessStats, void>({
      query: () => `${prefix}/stats`,
      providesTags: ['Processes', 'HasCompanyHeader'],
    }),
    createProcess: builder.mutation<Process, ProcessInput>({
      query: (body) => ({
        url: prefix,
        method: 'POST',
        body,
      }),
      invalidatesTags: ['Processes'],
    }),
    updateProcess: builder.mutation<Process, { id: number; body: ProcessInput }>({
      query: ({ id, body }) => ({
        url: `${prefix}/${id}`,
        method: 'PUT',
        body,
      }),
      invalidatesTags: ['Processes'],
    }),
    updateProcessStatus: builder.mutation<Process, { id: number; status: string }>({
      query: ({ id, status }) => ({
        url: `${prefix}/${id}/status`,
        method: 'PATCH',
        body: { status },
      }),
      invalidatesTags: ['Processes'],
    }),
    deleteProcess: builder.mutation<void, number>({
      query: (id) => ({
        url: `${prefix}/${id}`,
        method: 'DELETE',
      }),
      invalidatesTags: ['Processes'],
    }),
    recalculateProcess: builder.mutation<Process, number>({
      query: (id) => ({
        url: `${prefix}/${id}/recalculate`,
        method: 'POST',
      }),
      invalidatesTags: ['Processes'],
    }),
    uploadDocument: builder.mutation<Process, { id: number; file: File }>({
      query: ({ id, file }) => {
        const formData = new FormData()
        formData.append('file', file)
        return {
          url: `${prefix}/${id}/document`,
          method: 'POST',
          body: formData,
        }
      },
      invalidatesTags: ['Processes'],
    }),
    deleteDocument: builder.mutation<Process, number>({
      query: (id) => ({
        url: `${prefix}/${id}/document`,
        method: 'DELETE',
      }),
      invalidatesTags: ['Processes'],
    }),
    restoreProcess: builder.mutation<Process, number>({
      query: (id) => ({
        url: `${prefix}/${id}/restore`,
        method: 'POST',
      }),
      invalidatesTags: ['Processes'],
    }),
  }),
  overrideExisting: false,
})

export const {
  useProcessesQuery,
  useProcessQuery,
  useProcessStatsQuery,
  useCreateProcessMutation,
  useUpdateProcessMutation,
  useUpdateProcessStatusMutation,
  useDeleteProcessMutation,
  useRecalculateProcessMutation,
  useUploadDocumentMutation,
  useDeleteDocumentMutation,
  useRestoreProcessMutation,
} = extendedApi

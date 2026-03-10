import { api } from '@/Core/Services/Api'

import type { Company } from '../Types'
import type { Area } from '../Types/Area'

const prefix = 'orgs'
const extendedApi = api.injectEndpoints({
  endpoints: (builder) => ({
    companies: builder.query<Company[], void>({
      query: () => `${prefix}/company`,
      providesTags: ['Companies'],
    }),
    company: builder.query<Company, number>({
      query: (id) => `${prefix}/company/${id}`,
      providesTags: ['Companies'],
    }),
    createCompany: builder.mutation<Company, { name: string; parentId?: number | null }>({
      query: (body) => ({
        url: `${prefix}/company`,
        method: 'POST',
        body,
      }),
      invalidatesTags: ['Companies'],
    }),
    updateCompany: builder.mutation<Company, { id: number; body: { name: string; parentId?: number | null } }>({
      query: ({ id, body }) => ({
        url: `${prefix}/company/${id}`,
        method: 'PUT',
        body,
      }),
      invalidatesTags: ['Companies'],
    }),
    deleteCompany: builder.mutation<void, number>({
      query: (id) => ({
        url: `${prefix}/company/${id}`,
        method: 'DELETE',
      }),
      invalidatesTags: ['Companies'],
    }),
    companyAreas: builder.query<Area[], number>({
      query: (companyId) => `${prefix}/company/${companyId}/areas`,
      providesTags: (_result, _err, companyId) => [{ type: 'CompanyAreas' as const, id: companyId }],
    }),
    createCompanyArea: builder.mutation<Area, { companyId: number; name: string }>({
      query: ({ companyId, name }) => ({
        url: `${prefix}/company/${companyId}/areas`,
        method: 'POST',
        body: { name },
      }),
      invalidatesTags: (_result, _err, { companyId }) => [
        { type: 'CompanyAreas' as const, id: companyId },
        'Areas',
      ],
    }),
  }),
  overrideExisting: false,
})

export const {
  useCompaniesQuery,
  useCompanyQuery,
  useCreateCompanyMutation,
  useUpdateCompanyMutation,
  useDeleteCompanyMutation,
  useCompanyAreasQuery,
  useCreateCompanyAreaMutation,
} = extendedApi

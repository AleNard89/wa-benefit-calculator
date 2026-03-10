import { api } from '@/Core/Services/Api'

import type { Area } from '../Types/Area'

const prefix = 'orgs'

const extendedApi = api.injectEndpoints({
  endpoints: (builder) => ({
    areas: builder.query<Area[], void>({
      query: () => `${prefix}/area`,
      providesTags: ['HasCompanyHeader', 'Areas'],
    }),
    createArea: builder.mutation<Area, { name: string }>({
      query: (body) => ({
        url: `${prefix}/area`,
        method: 'POST',
        body,
      }),
      invalidatesTags: ['Areas'],
    }),
    updateArea: builder.mutation<Area, { id: number; body: { name: string } }>({
      query: ({ id, body }) => ({
        url: `${prefix}/area/${id}`,
        method: 'PUT',
        body,
      }),
      invalidatesTags: ['Areas'],
    }),
    deleteArea: builder.mutation<void, number>({
      query: (id) => ({
        url: `${prefix}/area/${id}`,
        method: 'DELETE',
      }),
      invalidatesTags: ['Areas'],
    }),
    userAreas: builder.query<Area[], number>({
      query: (userId) => `${prefix}/user/${userId}/areas`,
      providesTags: ['Areas'],
    }),
    setUserAreas: builder.mutation<void, { userId: number; areaIds: number[] }>({
      query: ({ userId, areaIds }) => ({
        url: `${prefix}/user/${userId}/areas`,
        method: 'PUT',
        body: { areaIds },
      }),
      invalidatesTags: ['Areas'],
    }),
  }),
  overrideExisting: false,
})

export const {
  useAreasQuery,
  useCreateAreaMutation,
  useUpdateAreaMutation,
  useDeleteAreaMutation,
  useUserAreasQuery,
  useSetUserAreasMutation,
} = extendedApi

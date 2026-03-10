import { api } from '@/Core/Services/Api'

import type { AuthenticatedUser, User, Role, Permission } from '../Types'

export interface SignInBaseResponse {
  accessToken: string
}

export type SignInFormData = { email: string; password: string }

const prefix = 'auth'
const extendedApi = api.injectEndpoints({
  endpoints: (builder) => ({
    currentUser: builder.query<Omit<AuthenticatedUser, 'companyPermissions'>, string | void>({
      query: (token?: string) => ({
        url: `${prefix}/whoami`,
        headers: token ? { Authorization: `Bearer ${token}` } : {},
      }),
      providesTags: ['AuthenticatedUser'],
    }),
    users: builder.query<User[], void>({
      query: () => `${prefix}/user`,
      providesTags: ['HasCompanyHeader', 'Users'],
    }),
    signIn: builder.mutation<SignInBaseResponse, SignInFormData>({
      query: (body) => ({
        url: `${prefix}/token/obtain`,
        method: 'POST',
        body,
      }),
    }),
    createUser: builder.mutation<User, Record<string, unknown>>({
      query: (body) => ({
        url: `${prefix}/user`,
        method: 'POST',
        body,
      }),
      invalidatesTags: ['Users'],
    }),
    updateUser: builder.mutation<User, { id: number; body: Record<string, unknown> }>({
      query: ({ id, body }) => ({
        url: `${prefix}/user/${id}`,
        method: 'PUT',
        body,
      }),
      invalidatesTags: ['Users'],
    }),
    deleteUser: builder.mutation<void, number>({
      query: (id) => ({
        url: `${prefix}/user/${id}`,
        method: 'DELETE',
      }),
      invalidatesTags: ['Users'],
    }),
    updateUserPassword: builder.mutation<User, { id: number; body: { currentPassword: string; password: string } }>({
      query: ({ id, body }) => ({
        url: `${prefix}/user/${id}/password`,
        method: 'PUT',
        body,
      }),
    }),
    permissions: builder.query<Permission[], void>({
      query: () => `${prefix}/permission`,
      providesTags: ['Permission'],
    }),
    roles: builder.query<Role[], void>({
      query: () => `${prefix}/role`,
      providesTags: ['Role'],
    }),
    createRole: builder.mutation<Role, { name: string; description: string; permissionIds: number[] }>({
      query: (body) => ({
        url: `${prefix}/role`,
        method: 'POST',
        body,
      }),
      invalidatesTags: ['Role'],
    }),
    updateRole: builder.mutation<Role, { id: number; body: { name: string; description: string; permissionIds: number[] } }>({
      query: ({ id, body }) => ({
        url: `${prefix}/role/${id}`,
        method: 'PUT',
        body,
      }),
      invalidatesTags: ['Role'],
    }),
    deleteRole: builder.mutation<void, number>({
      query: (id) => ({
        url: `${prefix}/role/${id}`,
        method: 'DELETE',
      }),
      invalidatesTags: ['Role'],
    }),
  }),
  overrideExisting: false,
})

export const {
  useCurrentUserQuery,
  useLazyCurrentUserQuery,
  useSignInMutation,
  useUsersQuery,
  useUpdateUserPasswordMutation,
  useCreateUserMutation,
  useUpdateUserMutation,
  useDeleteUserMutation,
  useRolesQuery,
  usePermissionsQuery,
  useCreateRoleMutation,
  useUpdateRoleMutation,
  useDeleteRoleMutation,
} = extendedApi

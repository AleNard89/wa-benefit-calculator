import { type AuthState, logout, setToken } from '@/Auth/Redux'
import Config from '@/Config'
import type { RootState } from '@/Core/Redux/Store'
import Logger from '@/Core/Services/Logger'
import {
  type BaseQueryFn,
  type FetchArgs,
  type FetchBaseQueryError,
  type QueryReturnValue,
  createApi,
  fetchBaseQuery,
} from '@reduxjs/toolkit/query/react'

export type MessageResponse = { message: string }
export type ResponseError = FetchBaseQueryError & { data?: { message: string } }

const baseQuery = fetchBaseQuery({
  baseUrl: Config.api.basePath,
  credentials: 'include',
  prepareHeaders: (headers, { getState }) => {
    headers.set('Accept-Language', 'it')
    const { auth, orgs } = getState() as RootState
    const companyId = orgs.companyId ?? auth.user?.companyRoles?.[0]?.company.id.toString()
    if (companyId) headers.set('X-Company-Id', companyId)
    if (auth.token) {
      headers.set('Authorization', `Bearer ${auth.token}`)
    }
    return headers
  },
})

let resolveRefreshTokenPromise: () => void
let refreshTokenPromise: Promise<void>
let refreshLocked = false

const isRefreshTokenRequest = (request: FetchArgs | string): boolean =>
  typeof request === 'string' ? request.includes('auth/token/refresh') : request.url.includes('auth/token/refresh')

const baseQueryWithReauth: BaseQueryFn<FetchArgs | string, unknown> = async (args, api, extraOptions) => {
  let result = await baseQuery(args, api, extraOptions)

  const { auth } = api.getState() as { auth: AuthState }
  const hasToken = Boolean(auth.token)
  const isRefreshRequest = isRefreshTokenRequest(args)

  if (result.error && result.error.status === 401 && !isRefreshRequest && hasToken) {
    if (!refreshLocked) {
      refreshTokenPromise = new Promise((resolve) => {
        refreshLocked = true
        resolveRefreshTokenPromise = resolve
      })
      const credentials = (await baseQuery(
        { url: 'auth/token/refresh', method: 'POST' },
        api,
        extraOptions,
      )) as QueryReturnValue<{ accessToken: string }>

      if (credentials.data) {
        api.dispatch(setToken(credentials.data))
        resolveRefreshTokenPromise()
        refreshLocked = false
      } else {
        resolveRefreshTokenPromise()
        refreshLocked = false
        try {
          await fetch(`${Config.api.basePath}/auth/logout`, { method: 'POST', credentials: 'include' })
        } catch {
          // proceed with local logout
        }
        api.dispatch(logout())
        window.location.assign(Config.urls.signIn)
        return result
      }
    }

    if (refreshLocked) {
      Logger.debug('API wait, someone is refreshing token')
      await refreshTokenPromise
    }

    result = await baseQuery(args, api, extraOptions)
  }

  return result
}

export const apiTags = ['AuthenticatedUser', 'Users', 'Roles', 'Permission', 'Companies', 'CompanyAreas', 'Areas', 'Processes', 'Conversations', 'Messages', 'HasCompanyHeader', 'Connectors', 'Jobs', 'QueueItems', 'QueueDefinitions', 'Schedules', 'ProcessNames', 'ProcessQueueMap']

export const api = createApi({
  reducerPath: 'api',
  baseQuery: baseQueryWithReauth,
  tagTypes: apiTags,
  endpoints: () => ({}),
})

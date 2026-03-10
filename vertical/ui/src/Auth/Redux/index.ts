import type { AuthenticatedUser } from '@/Auth/Types'
import type { RootState } from '@/Core/Redux/Store'
import { selectCompanies, selectCompanyId } from '@/Orgs/Redux'
import type { Company } from '@/Orgs/Types'
import { createSelector, createSlice } from '@reduxjs/toolkit'
import type { PayloadAction } from '@reduxjs/toolkit'

export type AuthTokenPayload = { accessToken: string }

export interface AuthState {
  token: string
  user: AuthenticatedUser | null
}

const initialState: AuthState = {
  token: '',
  user: null,
}

const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    setToken: (state, { payload }: PayloadAction<AuthTokenPayload>) => {
      state.token = payload.accessToken
    },
    resetCredentials: (state) => {
      state.token = ''
      state.user = null
    },
    logout(state) {
      state.token = ''
      state.user = null
    },
    setUser: (state, { payload }: PayloadAction<Omit<AuthenticatedUser, 'companyPermissions'>>) => {
      const companyPermissions = payload.companyRoles.reduce(
        (acc, current) => ({
          ...acc,
          [current.company.id]: current.roles.flatMap((r) => r.permissions.map((p) => p.code)),
        }),
        {} as Record<number, string[]>,
      )
      state.user = { ...payload, companyPermissions, currentPermissions: [] }
    },
  },
})

export const { setToken, setUser, resetCredentials, logout } = authSlice.actions
export default authSlice.reducer

export const selectAuth = (state: RootState): AuthState => state.auth
export const selectAccessToken = (state: RootState): string => state.auth.token

const selectAuthUser = (state: RootState): AuthenticatedUser | null => state.auth.user
export const selectCurrentUser = createSelector([selectAuthUser, selectCompanyId], (user, companyId) => {
  if (!user) return null
  const currentCompanyPermissions = companyId ? (user.companyPermissions[parseInt(companyId)] ?? []) : []
  return { ...user, currentPermissions: currentCompanyPermissions }
})

export const selectUserCompanies = createSelector([selectCurrentUser, selectCompanies], (user, companies) => {
  const userCompanies: Company[] = []
  if (!user) return userCompanies
  if (user.isSuperuser) return companies
  return user.companyRoles.map(({ company }) => company)
})

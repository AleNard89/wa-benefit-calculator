import type { RootState } from '@/Core/Redux/Store'
import { readFromStorage } from '@/Core/Services/Storage'
import { createSlice, type PayloadAction } from '@reduxjs/toolkit'

import type { Company } from '../Types'

export type OrgsState = {
  companyId: string
}

const initialState: OrgsState = {
  companyId: readFromStorage<string>('companyId', ''),
}

const orgsSlice = createSlice({
  name: 'orgs',
  initialState,
  reducers: {
    setCompanyId(state, { payload }: PayloadAction<number>) {
      state.companyId = payload.toString()
    },
    unsetCompanyId(state) {
      state.companyId = ''
    },
  },
})

export const { setCompanyId, unsetCompanyId } = orgsSlice.actions
export default orgsSlice.reducer

export const selectCompanyId = (state: RootState): string => state.orgs.companyId.toString()
export const selectCompanies = (state: RootState): Company[] =>
  (state.api.queries['companies(undefined)']?.data as Company[]) ?? []

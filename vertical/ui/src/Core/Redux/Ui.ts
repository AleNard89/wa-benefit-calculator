import { type PayloadAction, createSlice } from '@reduxjs/toolkit'

import type { RootState } from './Store'

export interface Breadcrumb {
  label: string
  path: string
}

interface UiState {
  breadcrumbs: Breadcrumb[]
}

const initialState: UiState = {
  breadcrumbs: [],
}

const uiSlice = createSlice({
  name: 'ui',
  initialState,
  reducers: {
    setBreadcrumbs: (state, { payload }: PayloadAction<Breadcrumb[]>) => {
      state.breadcrumbs = payload
    },
  },
})

export const { setBreadcrumbs } = uiSlice.actions
export default uiSlice.reducer

export const selectBreadcrumbs = (state: RootState): Breadcrumb[] => state.ui.breadcrumbs

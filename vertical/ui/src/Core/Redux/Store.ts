import authMiddleware from '@/Auth/Redux/middlewares'
import { setCompanyId, unsetCompanyId } from '@/Orgs/Redux'
import orgsMiddleware from '@/Orgs/Redux/middlewares'
import { configureStore } from '@reduxjs/toolkit'

import { api } from '../Services/Api'
import { rtkQueryErrorNotifier } from './ApiErrorMiddleware'
import localStorageRemoveMiddleware from './LocalStorageRemoveMiddleware'
import localStorageSetMiddleware from './LocalStorageSetMiddleware'
import { RootReducer } from './Root'

const store = configureStore({
  reducer: RootReducer,
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware({
      immutableCheck: false,
      serializableCheck: false,
    })
      .prepend(authMiddleware.middleware, orgsMiddleware.middleware)
      .concat(
        api.middleware,
        rtkQueryErrorNotifier,
        localStorageSetMiddleware('companyId', setCompanyId.type),
        localStorageRemoveMiddleware('companyId', [unsetCompanyId.type]),
      ),
  devTools: import.meta.env.VITE_DEV,
})

export type RootState = ReturnType<typeof RootReducer>

export default store

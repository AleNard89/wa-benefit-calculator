import authReducer from '@/Auth/Redux'
import orgsReducer from '@/Orgs/Redux'
import { combineReducers } from 'redux'

import { api } from '../Services/Api'
import uiReducer from './Ui'

export const RootReducer = combineReducers({
  api: api.reducer,
  auth: authReducer,
  ui: uiReducer,
  orgs: orgsReducer,
})

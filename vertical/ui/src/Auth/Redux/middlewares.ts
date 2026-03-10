import type { RootState } from '@/Core/Redux/Store'
import { api } from '@/Core/Services/Api'
import { setCompanyId, unsetCompanyId } from '@/Orgs/Redux'
import { createListenerMiddleware } from '@reduxjs/toolkit'

import { logout, setUser } from '.'

const listenerMiddleware = createListenerMiddleware()

listenerMiddleware.startListening({
  actionCreator: logout,
  effect: async (_, listenerApi) => {
    listenerApi.dispatch(unsetCompanyId())
    listenerApi.dispatch(api.util.resetApiState())
  },
})

listenerMiddleware.startListening({
  actionCreator: setUser,
  effect: async ({ payload }, listenerApi) => {
    const userCompanyIds = payload.companyRoles.map((cr) => cr.company.id)
    const currentCompanyId = (listenerApi.getState() as RootState).orgs.companyId

    if (userCompanyIds.length === 0) {
      if (!payload.isSuperuser) {
        listenerApi.dispatch(unsetCompanyId())
      }
      return
    }

    const hasCompanyId = userCompanyIds.includes(+currentCompanyId)
    if (!currentCompanyId || (!hasCompanyId && !payload.isSuperuser)) {
      listenerApi.dispatch(setCompanyId(userCompanyIds[0]))
    }
  },
})

export default listenerMiddleware

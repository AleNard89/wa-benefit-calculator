import type { RootState } from '@/Core/Redux/Store'
import { api } from '@/Core/Services/Api'
import { setCompanyId, unsetCompanyId } from '@/Orgs/Redux'
import type { Company } from '@/Orgs/Types'
import { createListenerMiddleware } from '@reduxjs/toolkit'
import Config from '@/Config'

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
        return
      }
      // Superuser with no companyRoles: auto-select first company if none stored
      if (!currentCompanyId) {
        try {
          const token = (listenerApi.getState() as RootState).auth.token
          const res = await fetch(`${Config.api.basePath}/orgs/company`, {
            headers: { Authorization: `Bearer ${token}` },
          })
          if (res.ok) {
            const companies = (await res.json()) as Company[]
            if (companies.length > 0) {
              listenerApi.dispatch(setCompanyId(companies[0].id))
            }
          }
        } catch {
          // proceed without company, user can select manually
        }
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

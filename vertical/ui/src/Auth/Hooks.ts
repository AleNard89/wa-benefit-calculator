import { useEffect, useState } from 'react'
import { useDispatch, useSelector } from 'react-redux'

import Config from '../Config'
import Logger from '../Core/Services/Logger'
import store from '../Core/Redux/Store'
import type { RootState } from '../Core/Redux/Store'
import { selectCurrentUser, selectUserCompanies, setToken, setUser } from './Redux'
import { useLazyCurrentUserQuery } from './Services/Api'
import { useLazyCompaniesQuery } from '@/Orgs/Services/Api'
import { setCompanyId } from '@/Orgs/Redux'
import type { AuthenticatedUser } from './Types'
import type { Company } from '@/Orgs/Types'

export const useAuthentication = () => {
  const [isComplete, setIsComplete] = useState<boolean>(false)
  const [getCurrentUser] = useLazyCurrentUserQuery()
  const [getCompanies] = useLazyCompaniesQuery()
  const dispatch = useDispatch()

  useEffect(() => {
    const bootstrap = async () => {
      try {
        const res = await fetch(`${Config.api.basePath}/auth/token/refresh`, {
          method: 'POST',
          credentials: 'include',
        })
        if (res.ok) {
          const data = (await res.json()) as { accessToken: string }
          dispatch(setToken(data))
        }
      } catch {
        // refresh failed
      }

      try {
        const user = await getCurrentUser().unwrap()
        dispatch(setUser(user))

        // Superusers with no companyRoles need a company selected before the app renders.
        // The auth middleware skips setCompanyId for them, so we fetch here if needed.
        const currentCompanyId = (store.getState() as RootState).orgs.companyId
        if (user.isSuperuser && user.companyRoles.length === 0 && !currentCompanyId) {
          try {
            const companies = await getCompanies().unwrap()
            if (companies.length > 0) {
              dispatch(setCompanyId(companies[0].id))
            }
          } catch {
            // proceed without company, user can select manually
          }
        }

        setTimeout(() => setIsComplete(true), 200)
      } catch (err) {
        Logger.error('Fetch authenticated user error', err)
        setIsComplete(true)
      }
    }
    void bootstrap()
  }, [getCurrentUser, getCompanies, dispatch])

  return { isComplete }
}

export const useCurrentUser = (): AuthenticatedUser | null => useSelector(selectCurrentUser)

export const useUserCompanies = (): Company[] => useSelector(selectUserCompanies)

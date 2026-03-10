import { useEffect, useState } from 'react'
import { useDispatch, useSelector } from 'react-redux'

import Config from '../Config'
import Logger from '../Core/Services/Logger'
import { selectCurrentUser, selectUserCompanies, setToken, setUser } from './Redux'
import { useLazyCurrentUserQuery } from './Services/Api'
import type { AuthenticatedUser } from './Types'
import type { Company } from '@/Orgs/Types'

export const useAuthentication = () => {
  const [isComplete, setIsComplete] = useState<boolean>(false)
  const [getCurrentUser] = useLazyCurrentUserQuery()
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
        setTimeout(() => setIsComplete(true), 200)
      } catch (err) {
        Logger.error('Fetch authenticated user error', err)
        setIsComplete(true)
      }
    }
    void bootstrap()
  }, [getCurrentUser, dispatch])

  return { isComplete }
}

export const useCurrentUser = (): AuthenticatedUser | null => useSelector(selectCurrentUser)

export const useUserCompanies = (): Company[] => useSelector(selectUserCompanies)

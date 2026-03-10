import { Navigate } from 'react-router-dom'

import Config from '@/Config'
import { useCurrentUser } from '../Hooks'

export default function AdminRoute({ children }: { children: React.ReactNode }) {
  const user = useCurrentUser()

  if (!user) return <Navigate to={Config.urls.signIn} replace />

  const isAdmin = user.isSuperuser || user.companyRoles?.some((cr) => cr.roles.some((r) => r.name === 'Admin'))
  if (!isAdmin) return <Navigate to={Config.urls.home} replace />

  return <>{children}</>
}

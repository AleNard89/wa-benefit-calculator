import { Navigate } from 'react-router-dom'

import Config from '@/Config'
import { useCurrentUser } from '../Hooks'

interface PrivateRouteProps {
  children: React.ReactNode
  permissionCheck?: (permissions: string[]) => boolean
}

export default function PrivateRoute({ children, permissionCheck }: PrivateRouteProps) {
  const user = useCurrentUser()

  if (!user) return <Navigate to={Config.urls.signIn} replace />
  if (permissionCheck && !user.isSuperuser && !permissionCheck(user.currentPermissions)) {
    return <Navigate to={Config.urls.home} replace />
  }

  return <>{children}</>
}

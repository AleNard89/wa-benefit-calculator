import { Flex, Spinner } from '@chakra-ui/react'
import { useSelector } from 'react-redux'
import { Navigate } from 'react-router-dom'

import Config from '@/Config'
import { selectCompanyId } from '@/Orgs/Redux'
import { useCurrentUser } from '../Hooks'

interface PrivateRouteProps {
  children: React.ReactNode
  permissionCheck?: (permissions: string[]) => boolean
}

export default function PrivateRoute({ children, permissionCheck }: PrivateRouteProps) {
  const user = useCurrentUser()
  const companyId = useSelector(selectCompanyId)

  if (!user) return <Navigate to={Config.urls.signIn} replace />
  if (permissionCheck && !user.isSuperuser && !permissionCheck(user.currentPermissions)) {
    return <Navigate to={Config.urls.home} replace />
  }

  // Superusers with no companyRoles need a company fetched asynchronously before rendering.
  // Show a spinner until the company is selected to avoid API calls with missing company ID.
  if (user.isSuperuser && !companyId) {
    return (
      <Flex h="100vh" align="center" justify="center">
        <Spinner size="xl" />
      </Flex>
    )
  }

  return <>{children}</>
}

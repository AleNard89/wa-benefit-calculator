import { Box, Flex, Text } from '@chakra-ui/react'
import { useState } from 'react'
import { LuBot, LuBuilding2, LuKeyRound, LuLink, LuShield, LuUsers } from 'react-icons/lu'

import { useCurrentUser } from '@/Auth/Hooks'
import CompanyView from '@/Orgs/Views/CompanyView'
import UsersView from '@/Auth/Views/UsersView'
import RolesView from '@/Auth/Views/RolesView'
import PasswordView from '@/Auth/Views/PasswordView'
import ConnettoriView from '@/Orchestrator/ConnettoriView'
import MappingTab from '@/Orchestrator/MappingTab'

interface TabItem {
  id: string
  label: string
  icon: React.ReactNode
  permission?: string
  adminOnly?: boolean
}

const tabs: TabItem[] = [
  { id: 'company', label: 'Aziende', icon: <LuBuilding2 size={15} />, permission: 'orgs:company.read' },
  { id: 'users', label: 'Utenti', icon: <LuUsers size={15} />, permission: 'auth:user.read' },
  { id: 'roles', label: 'Ruoli', icon: <LuShield size={15} />, permission: 'auth:role.read' },
  { id: 'connettori', label: 'Connettori', icon: <LuBot size={15} />, adminOnly: true },
  { id: 'mapping', label: 'Mapping Bot-Code', icon: <LuLink size={15} />, adminOnly: true },
  { id: 'password', label: 'Password', icon: <LuKeyRound size={15} /> },
]

export default function SettingsView() {
  const user = useCurrentUser()
  const perms = user?.currentPermissions ?? []
  const isSuperuser = user?.isSuperuser ?? false

  const isAdmin = isSuperuser || user?.companyRoles?.some((cr) => cr.roles.some((r) => r.name === 'Admin'))
  const visibleTabs = tabs.filter((t) => {
    if (t.adminOnly) return isAdmin
    return !t.permission || isSuperuser || perms.includes(t.permission)
  })
  const [activeTab, setActiveTab] = useState(visibleTabs[0]?.id ?? 'password')

  return (
    <Box>
      {/* Header */}
      <Flex justify="space-between" align="center" mb={5}>
        <Box>
          <Text fontSize="20px" fontWeight="700" color="#1d1d1f" letterSpacing="-0.3px">
            Impostazioni
          </Text>
          <Text fontSize="14px" color="#86868b" mt={1}>
            Gestione aziende, utenti, ruoli e account
          </Text>
        </Box>
      </Flex>

      {/* Pill Tabs */}
      <Flex
        bg="#f5f5f7"
        borderRadius="12px"
        p="4px"
        gap="2px"
        mb={5}
        w="fit-content"
      >
        {visibleTabs.map((tab) => (
          <Flex
            key={tab.id}
            as="button"
            type="button"
            align="center"
            gap={1.5}
            px={4}
            py={2}
            borderRadius="10px"
            fontSize="13px"
            fontWeight={activeTab === tab.id ? '600' : '400'}
            color={activeTab === tab.id ? '#1d1d1f' : '#86868b'}
            bg={activeTab === tab.id ? 'white' : 'transparent'}
            boxShadow={activeTab === tab.id ? '0 1px 4px rgba(0,0,0,0.08)' : 'none'}
            cursor="pointer"
            transition="all 0.2s"
            _hover={{ color: '#1d1d1f' }}
            onClick={() => setActiveTab(tab.id)}
          >
            <Box opacity={activeTab === tab.id ? 1 : 0.55}>{tab.icon}</Box>
            {tab.label}
          </Flex>
        ))}
      </Flex>

      {/* Content Card */}
      <Box
        bg="white"
        borderRadius="16px"
        boxShadow="0 2px 20px rgba(0,0,0,0.06)"
        p={7}
      >
        {activeTab === 'company' && <CompanyView />}
        {activeTab === 'users' && <UsersView />}
        {activeTab === 'roles' && <RolesView />}
        {activeTab === 'connettori' && <ConnettoriView />}
        {activeTab === 'mapping' && <MappingTab />}
        {activeTab === 'password' && <PasswordView />}
      </Box>
    </Box>
  )
}

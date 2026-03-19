import { Box, Flex, Text, VStack } from '@chakra-ui/react'
import { useState } from 'react'
import { useDispatch, useSelector } from 'react-redux'
import { useLocation, useNavigate } from 'react-router-dom'
import { LuBot, LuBuilding2, LuCheck, LuChevronDown, LuFilePlus, LuLayoutDashboard, LuList, LuLogOut, LuMessageCircle } from 'react-icons/lu'

import { useCurrentUser } from '@/Auth/Hooks'
import { logout } from '@/Auth/Redux'
import { selectCompanyId, setCompanyId } from '@/Orgs/Redux'
import { useCompaniesQuery } from '@/Orgs/Services/Api'
import type { Company } from '@/Orgs/Types'

interface MenuItem {
  label: string
  icon: React.ReactNode
  path: string
}

interface MenuSection {
  title: string
  items: MenuItem[]
  adminOnly?: boolean
}

const menuSections: MenuSection[] = [
  {
    title: 'Processi',
    items: [
      { label: 'Dashboard', icon: <LuLayoutDashboard size={18} />, path: '/' },
      { label: 'Nuova Proposta', icon: <LuFilePlus size={18} />, path: '/processes/create' },
      { label: 'Lista Processi', icon: <LuList size={18} />, path: '/processes/list' },
    ],
  },
  {
    title: 'Connettori',
    adminOnly: true,
    items: [
      { label: 'Orchestrator', icon: <LuBot size={18} />, path: '/orchestrator' },
    ],
  },
  {
    title: 'AI',
    items: [
      { label: 'Chat', icon: <LuMessageCircle size={18} />, path: '/chat' },
    ],
  },
]

const AVATAR_URL = 'https://api.dicebear.com/9.x/personas/svg?seed=admin&backgroundColor=b6e3f4'

function CompanySwitcher() {
  const dispatch = useDispatch()
  const user = useCurrentUser()
  const { data: allCompanies } = useCompaniesQuery()
  const currentCompanyId = useSelector(selectCompanyId)
  const [open, setOpen] = useState(false)

  const companies = user?.isSuperuser
    ? (allCompanies ?? [])
    : (user?.companyRoles?.map((cr) => cr.company) ?? [])

  const currentCompany = companies.find((c) => c.id.toString() === currentCompanyId) ?? companies[0]

  const handleSelect = (company: Company) => {
    dispatch(setCompanyId(company.id))
    setOpen(false)
    window.location.reload()
  }

  if (companies.length <= 1) {
    return (
      <Box px={5} pt={5} pb={2}>
        <Text fontSize="11px" fontWeight="600" color="#86868b" letterSpacing="0.5px" textTransform="uppercase" mb={1}>
          Orbita
        </Text>
        <Text fontSize="15px" fontWeight="700" color="#1d1d1f" letterSpacing="-0.2px">
          {currentCompany?.name ?? 'Nessuna azienda'}
        </Text>
      </Box>
    )
  }

  return (
    <Box px={3} pt={4} pb={1} position="relative">
      <Text px={2} fontSize="11px" fontWeight="600" color="#86868b" letterSpacing="0.5px" textTransform="uppercase" mb={1}>
        Orbita
      </Text>
      <Flex
        align="center"
        justify="space-between"
        px={3}
        py={2}
        borderRadius="10px"
        cursor="pointer"
        _hover={{ bg: '#f5f5f7' }}
        transition="all 0.15s"
        onClick={() => setOpen(!open)}
      >
        <Flex align="center" gap={2.5}>
          <Flex w="28px" h="28px" align="center" justify="center" borderRadius="8px" bg="#007aff15" color="#007aff" flexShrink={0}>
            <LuBuilding2 size={14} />
          </Flex>
          <Text fontSize="14px" fontWeight="600" color="#1d1d1f" truncate>
            {currentCompany?.name ?? 'Seleziona'}
          </Text>
        </Flex>
        <Box color="#86868b" transition="transform 0.15s" transform={open ? 'rotate(180deg)' : 'none'}>
          <LuChevronDown size={16} />
        </Box>
      </Flex>

      {open && (
        <>
          <Box position="fixed" inset={0} zIndex={999} onClick={() => setOpen(false)} />
          <Box
            position="absolute"
            top="100%"
            left={3}
            right={3}
            mt={1}
            bg="white"
            borderRadius="12px"
            boxShadow="0 8px 30px rgba(0,0,0,0.12)"
            border="1px solid #f0f0f2"
            zIndex={1000}
            overflow="hidden"
            maxH="240px"
            overflowY="auto"
          >
            {companies.map((company) => {
              const isSelected = company.id.toString() === currentCompanyId
              return (
                <Flex
                  key={company.id}
                  align="center"
                  justify="space-between"
                  px={3}
                  py={2.5}
                  cursor="pointer"
                  bg={isSelected ? '#007aff08' : 'transparent'}
                  _hover={{ bg: isSelected ? '#007aff12' : '#f5f5f7' }}
                  transition="background 0.1s"
                  onClick={() => handleSelect(company)}
                >
                  <Flex align="center" gap={2.5}>
                    <Flex w="24px" h="24px" align="center" justify="center" borderRadius="6px" bg={isSelected ? '#007aff15' : '#f5f5f7'} color={isSelected ? '#007aff' : '#86868b'} flexShrink={0}>
                      <LuBuilding2 size={12} />
                    </Flex>
                    <Box>
                      <Text fontSize="13px" fontWeight={isSelected ? '600' : '400'} color="#1d1d1f">
                        {company.name}
                      </Text>
                      {company.parentId && (
                        <Text fontSize="10px" color="#86868b">
                          {companies.find((c) => c.id === company.parentId)?.name ?? ''}
                        </Text>
                      )}
                    </Box>
                  </Flex>
                  {isSelected && <LuCheck size={14} color="#007aff" />}
                </Flex>
              )
            })}
          </Box>
        </>
      )}
    </Box>
  )
}

export default function Sidebar() {
  const dispatch = useDispatch()
  const navigate = useNavigate()
  const location = useLocation()
  const user = useCurrentUser()

  const handleLogout = () => {
    dispatch(logout())
    navigate('/signin')
  }

  return (
    <Flex
      as="nav"
      direction="column"
      w="260px"
      minW="260px"
      h="calc(100vh - 24px)"
      m="12px"
      mr={0}
      bg="white"
      borderRadius="16px"
      boxShadow="0 2px 20px rgba(0,0,0,0.08)"
      overflow="visible"
    >
      <CompanySwitcher />

      <VStack gap={0} px={3} pt={2} align="stretch" flex={1}>
        {menuSections.filter((s) => {
          if (!s.adminOnly) return true
          return user?.isSuperuser || user?.companyRoles?.some((cr) => cr.roles.some((r) => r.name === 'Admin'))
        }).map((section, idx) => (
          <Box key={section.title} mt={idx > 0 ? 3 : 0}>
            <Text
              fontSize="10px"
              fontWeight="700"
              color="#86868b"
              letterSpacing="0.8px"
              textTransform="uppercase"
              px={3}
              mb={1}
            >
              {section.title}
            </Text>
            <VStack gap={0.5} align="stretch">
              {section.items.map((item) => {
                const isActive = location.pathname === item.path ||
                  (item.path !== '/' && location.pathname.startsWith(item.path))

                return (
                  <Flex
                    key={item.path}
                    align="center"
                    gap={3}
                    px={3}
                    py={2}
                    borderRadius="10px"
                    cursor="pointer"
                    bg={isActive ? '#007aff' : 'transparent'}
                    color={isActive ? 'white' : '#1d1d1f'}
                    _hover={{ bg: isActive ? '#007aff' : '#f5f5f7' }}
                    transition="all 0.15s"
                    onClick={() => navigate(item.path)}
                  >
                    <Box opacity={isActive ? 1 : 0.65}>{item.icon}</Box>
                    <Text fontSize="13px" fontWeight={isActive ? '600' : '400'}>{item.label}</Text>
                  </Flex>
                )
              })}
            </VStack>
          </Box>
        ))}
      </VStack>

      <Box px={3} pb={4} pt={2} borderTop="1px solid #f0f0f2">
        <Flex
          align="center"
          gap={3}
          px={3}
          py={2.5}
          borderRadius="10px"
          cursor="pointer"
          bg={location.pathname.startsWith('/settings') ? '#f5f5f7' : 'transparent'}
          _hover={{ bg: '#f5f5f7' }}
          transition="all 0.15s"
          onClick={() => navigate('/settings')}
        >
          <Box w="36px" h="36px" borderRadius="50%" overflow="hidden" bg="#e8e8ed" flexShrink={0}>
            <img src={AVATAR_URL} alt="Avatar" width={36} height={36} style={{ display: 'block' }} />
          </Box>
          <Box flex={1} overflow="hidden">
            <Text fontSize="13px" fontWeight="600" color="#1d1d1f" truncate>
              {user ? `${user.firstName} ${user.lastName}` : 'Utente'}
            </Text>
            <Text fontSize="11px" color="#86868b" truncate>
              {user?.email ?? ''}
            </Text>
          </Box>
          <Box
            as="button"
            p={1.5}
            borderRadius="6px"
            color="#86868b"
            cursor="pointer"
            _hover={{ bg: '#e8e8ed', color: '#ff3b30' }}
            transition="all 0.15s"
            onClick={(e) => { e.stopPropagation(); handleLogout() }}
            title="Esci"
          >
            <LuLogOut size={16} />
          </Box>
        </Flex>
      </Box>
    </Flex>
  )
}

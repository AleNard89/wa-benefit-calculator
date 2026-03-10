import { Box, Button, Flex, Input, Text } from '@chakra-ui/react'
import { useEffect, useState } from 'react'
import { LuBuilding2, LuPencil, LuPlus, LuTrash2, LuX } from 'react-icons/lu'

import { useCurrentUser } from '@/Auth/Hooks'
import {
  useUsersQuery,
  useDeleteUserMutation,
  useCreateUserMutation,
  useUpdateUserMutation,
  useRolesQuery,
} from '../Services/Api'
import { useCompaniesQuery } from '@/Orgs/Services/Api'
import { useAreasQuery, useUserAreasQuery, useSetUserAreasMutation } from '@/Orgs/Services/AreaApi'
import { toaster } from '@/Snippets/toaster'
import type { User, Role } from '../Types'
import type { Company } from '@/Orgs/Types'

const inputStyle = {
  bg: '#f5f5f7',
  border: 'none',
  borderRadius: '10px',
  h: '40px',
  fontSize: '14px',
  _focus: { bg: 'white', boxShadow: '0 0 0 3px rgba(0,122,255,0.15)' },
}

function Chip({ label, selected, onClick }: { label: string; selected: boolean; onClick: () => void }) {
  return (
    <Box
      as="button"
      type="button"
      px={3}
      py={1.5}
      borderRadius="8px"
      bg={selected ? '#007aff' : '#f5f5f7'}
      color={selected ? 'white' : '#1d1d1f'}
      fontSize="12px"
      fontWeight={selected ? '600' : '400'}
      cursor="pointer"
      transition="all 0.15s"
      _hover={{ opacity: 0.85 }}
      onClick={onClick}
    >
      {label}
    </Box>
  )
}

interface CompanyRoleEntry {
  companyId: number
  roleIds: number[]
}

function CompanyRolesSection({
  companies,
  roles,
  entries,
  onChange,
}: {
  companies: Company[]
  roles: Role[]
  entries: CompanyRoleEntry[]
  onChange: (entries: CompanyRoleEntry[]) => void
}) {
  const addCompany = (companyId: number) => {
    if (!entries.find((e) => e.companyId === companyId)) {
      onChange([...entries, { companyId, roleIds: [] }])
    }
  }

  const removeCompany = (companyId: number) => {
    onChange(entries.filter((e) => e.companyId !== companyId))
  }

  const toggleRole = (companyId: number, roleId: number) => {
    onChange(
      entries.map((e) =>
        e.companyId === companyId
          ? { ...e, roleIds: e.roleIds.includes(roleId) ? e.roleIds.filter((r) => r !== roleId) : [...e.roleIds, roleId] }
          : e,
      ),
    )
  }

  const availableCompanies = companies.filter((c) => !entries.find((e) => e.companyId === c.id))

  return (
    <Box>
      <Text fontSize="12px" fontWeight="600" color="#86868b" mb={2}>Aziende e Ruoli</Text>

      {entries.length === 0 && (
        <Text fontSize="12px" color="#86868b" mb={2}>Nessuna azienda assegnata</Text>
      )}

      <Flex direction="column" gap={3}>
        {entries.map((entry) => {
          const company = companies.find((c) => c.id === entry.companyId)
          return (
            <Box key={entry.companyId} bg="#f5f5f7" borderRadius="10px" p={3}>
              <Flex align="center" justify="space-between" mb={2}>
                <Flex align="center" gap={2}>
                  <LuBuilding2 size={14} color="#007aff" />
                  <Text fontSize="13px" fontWeight="600" color="#1d1d1f">{company?.name ?? `ID ${entry.companyId}`}</Text>
                </Flex>
                <Box
                  as="button"
                  type="button"
                  p={1}
                  borderRadius="6px"
                  color="#86868b"
                  cursor="pointer"
                  _hover={{ color: '#ff3b30' }}
                  onClick={() => removeCompany(entry.companyId)}
                >
                  <LuTrash2 size={13} />
                </Box>
              </Flex>
              <Flex gap={2} wrap="wrap">
                {roles.map((r) => (
                  <Chip
                    key={r.id}
                    label={r.name}
                    selected={entry.roleIds.includes(r.id)}
                    onClick={() => toggleRole(entry.companyId, r.id)}
                  />
                ))}
              </Flex>
            </Box>
          )
        })}
      </Flex>

      {availableCompanies.length > 0 && (
        <Flex gap={2} wrap="wrap" mt={2}>
          {availableCompanies.map((c) => (
            <Box
              key={c.id}
              as="button"
              type="button"
              px={3}
              py={1.5}
              borderRadius="8px"
              border="1px dashed #c7c7cc"
              bg="transparent"
              color="#86868b"
              fontSize="12px"
              cursor="pointer"
              _hover={{ borderColor: '#007aff', color: '#007aff' }}
              transition="all 0.15s"
              onClick={() => addCompany(c.id)}
            >
              + {c.name}
            </Box>
          ))}
        </Flex>
      )}
    </Box>
  )
}

function UserDialog({ user, open, onClose }: { user: User | null; open: boolean; onClose: () => void }) {
  const currentUser = useCurrentUser()
  const { data: roles } = useRolesQuery()
  const { data: companies } = useCompaniesQuery()
  const { data: areas } = useAreasQuery()
  const { data: userAreas } = useUserAreasQuery(user?.id ?? 0, { skip: !user?.id })
  const [createUser] = useCreateUserMutation()
  const [updateUser] = useUpdateUserMutation()
  const [setUserAreas] = useSetUserAreasMutation()

  const [email, setEmail] = useState('')
  const [firstName, setFirstName] = useState('')
  const [lastName, setLastName] = useState('')
  const [password, setPassword] = useState('')
  const [companyEntries, setCompanyEntries] = useState<CompanyRoleEntry[]>([])
  const [selectedAreaIds, setSelectedAreaIds] = useState<number[]>([])
  const [loading, setLoading] = useState(false)

  const isEdit = !!user

  useEffect(() => {
    if (open) {
      if (user) {
        setEmail(user.email)
        setFirstName(user.firstName)
        setLastName(user.lastName)
        setPassword('')
        setCompanyEntries(
          user.companyRoles?.map((cr) => ({
            companyId: cr.company.id,
            roleIds: cr.roles.map((r) => r.id),
          })) ?? [],
        )
        setSelectedAreaIds([])
      } else {
        setEmail('')
        setFirstName('')
        setLastName('')
        setPassword('')
        const defaultCompanyId = currentUser?.companyRoles?.[0]?.company.id
        setCompanyEntries(defaultCompanyId ? [{ companyId: defaultCompanyId, roleIds: [] }] : [])
        setSelectedAreaIds([])
      }
    }
  }, [open, user, currentUser])

  useEffect(() => {
    if (userAreas) setSelectedAreaIds(userAreas.map((a) => a.id))
  }, [userAreas])

  const toggleArea = (id: number) =>
    setSelectedAreaIds((prev) => (prev.includes(id) ? prev.filter((a) => a !== id) : [...prev, id]))

  const handleSubmit = async () => {
    if (!email || !firstName || !lastName || (!isEdit && !password)) {
      toaster.error({ title: 'Compila tutti i campi obbligatori' })
      return
    }
    if (companyEntries.length === 0) {
      toaster.error({ title: 'Assegna almeno un\'azienda' })
      return
    }
    setLoading(true)
    try {
      const body: Record<string, unknown> = {
        email,
        firstName,
        lastName,
        companies: companyEntries.map((e) => ({ companyId: e.companyId, roleIds: e.roleIds })),
      }
      let userId = user?.id
      if (isEdit) {
        await updateUser({ id: user!.id, body }).unwrap()
      } else {
        body.password = password
        const created = await createUser(body).unwrap()
        userId = created.id
      }
      if (userId) {
        await setUserAreas({ userId, areaIds: selectedAreaIds }).unwrap()
      }
      toaster.success({ title: isEdit ? 'Utente aggiornato' : 'Utente creato' })
      onClose()
    } catch {
      toaster.error({ title: 'Errore durante il salvataggio' })
    } finally {
      setLoading(false)
    }
  }

  if (!open) return null

  return (
    <Box position="fixed" top={0} left={0} right={0} bottom={0} zIndex={1000}>
      <Box position="absolute" inset={0} bg="blackAlpha.400" onClick={onClose} />
      <Flex position="absolute" inset={0} align="center" justify="center" p={4}>
        <Box
          bg="white"
          borderRadius="16px"
          boxShadow="0 20px 60px rgba(0,0,0,0.15)"
          maxW="520px"
          w="full"
          maxH="85vh"
          overflow="auto"
          onClick={(e) => e.stopPropagation()}
        >
          <Flex justify="space-between" align="center" px={6} pt={5} pb={3}>
            <Text fontSize="17px" fontWeight="700" color="#1d1d1f">
              {isEdit ? 'Modifica Utente' : 'Nuovo Utente'}
            </Text>
            <Box as="button" p={1} color="#86868b" _hover={{ color: '#1d1d1f' }} onClick={onClose}>
              <LuX size={18} />
            </Box>
          </Flex>

          <Box px={6} pb={6}>
            <Flex direction="column" gap={3.5}>
              <Box>
                <Text fontSize="12px" fontWeight="600" color="#86868b" mb={1}>Email *</Text>
                <Input value={email} onChange={(e) => setEmail(e.target.value)} placeholder="email@esempio.it" {...inputStyle} />
              </Box>

              <Flex gap={3}>
                <Box flex={1}>
                  <Text fontSize="12px" fontWeight="600" color="#86868b" mb={1}>Nome *</Text>
                  <Input value={firstName} onChange={(e) => setFirstName(e.target.value)} placeholder="Nome" {...inputStyle} />
                </Box>
                <Box flex={1}>
                  <Text fontSize="12px" fontWeight="600" color="#86868b" mb={1}>Cognome *</Text>
                  <Input value={lastName} onChange={(e) => setLastName(e.target.value)} placeholder="Cognome" {...inputStyle} />
                </Box>
              </Flex>

              {!isEdit && (
                <Box>
                  <Text fontSize="12px" fontWeight="600" color="#86868b" mb={1}>Password *</Text>
                  <Input type="password" value={password} onChange={(e) => setPassword(e.target.value)} placeholder="Min. 8 caratteri" {...inputStyle} />
                </Box>
              )}

              {companies && roles && (
                <CompanyRolesSection
                  companies={companies}
                  roles={roles}
                  entries={companyEntries}
                  onChange={setCompanyEntries}
                />
              )}

              {areas && areas.length > 0 && (
                <Box>
                  <Text fontSize="12px" fontWeight="600" color="#86868b" mb={2}>Aree</Text>
                  <Flex gap={2} wrap="wrap">
                    {areas.map((a) => (
                      <Chip key={a.id} label={a.name} selected={selectedAreaIds.includes(a.id)} onClick={() => toggleArea(a.id)} />
                    ))}
                  </Flex>
                </Box>
              )}

              <Flex gap={2} mt={2}>
                <Button flex={1} bg="#f5f5f7" color="#1d1d1f" borderRadius="10px" h="40px" fontSize="14px" _hover={{ bg: '#e8e8ed' }} onClick={onClose}>
                  Annulla
                </Button>
                <Button flex={1} bg="#007aff" color="white" borderRadius="10px" h="40px" fontSize="14px" fontWeight="600" _hover={{ bg: '#0066d6' }} loading={loading} onClick={handleSubmit}>
                  {isEdit ? 'Salva' : 'Crea Utente'}
                </Button>
              </Flex>
            </Flex>
          </Box>
        </Box>
      </Flex>
    </Box>
  )
}

function UserRow({ user, onEdit, onDelete }: { user: User; onEdit: () => void; onDelete: (id: number) => void }) {
  return (
    <Flex
      align="center"
      justify="space-between"
      py={3}
      px={4}
      borderBottom="1px solid #f0f0f2"
      _hover={{ bg: '#fafafa' }}
      transition="background 0.15s"
      cursor="pointer"
      onClick={onEdit}
    >
      <Box>
        <Text fontSize="14px" fontWeight="500" color="#1d1d1f">
          {user.firstName} {user.lastName}
        </Text>
        <Text fontSize="12px" color="#86868b">{user.email}</Text>
        {user.companyRoles && user.companyRoles.length > 1 && (
          <Text fontSize="11px" color="#86868b" mt={0.5}>
            {user.companyRoles.map((cr) => cr.company.name).join(', ')}
          </Text>
        )}
      </Box>
      <Flex align="center" gap={2}>
        {user.isSuperuser && (
          <Box px={2} py={0.5} borderRadius="6px" bg="#007aff15">
            <Text fontSize="11px" fontWeight="600" color="#007aff">Superuser</Text>
          </Box>
        )}
        {user.companyRoles?.map((cr) =>
          cr.roles.map((r) => (
            <Box key={`${cr.company.id}-${r.id}`} px={2} py={0.5} borderRadius="6px" bg="#f5f5f7">
              <Text fontSize="11px" color="#86868b">{r.name}</Text>
            </Box>
          )),
        )}
        <Box
          as="button"
          p={1.5}
          borderRadius="6px"
          color="#86868b"
          cursor="pointer"
          _hover={{ color: '#007aff', bg: '#007aff10' }}
          transition="all 0.15s"
          onClick={(e) => { e.stopPropagation(); onEdit() }}
        >
          <LuPencil size={15} />
        </Box>
        {!user.isSuperuser && (
          <Box
            as="button"
            p={1.5}
            borderRadius="6px"
            color="#86868b"
            cursor="pointer"
            _hover={{ color: '#ff3b30', bg: '#ff3b3010' }}
            transition="all 0.15s"
            onClick={(e) => { e.stopPropagation(); onDelete(user.id) }}
          >
            <LuTrash2 size={15} />
          </Box>
        )}
      </Flex>
    </Flex>
  )
}

export default function UsersView() {
  const { data: users, isLoading } = useUsersQuery()
  const [deleteUser] = useDeleteUserMutation()
  const [dialogOpen, setDialogOpen] = useState(false)
  const [editingUser, setEditingUser] = useState<User | null>(null)

  const handleDelete = async (id: number) => {
    if (confirm('Sei sicuro di voler eliminare questo utente?')) {
      try {
        await deleteUser(id).unwrap()
        toaster.success({ title: 'Utente eliminato' })
      } catch {
        toaster.error({ title: "Errore durante l'eliminazione" })
      }
    }
  }

  return (
    <Box>
      <Flex justify="space-between" align="center" mb={4}>
        <Text fontSize="15px" fontWeight="700" color="#1d1d1f">Utenti</Text>
        <Button
          bg="#007aff" color="white" borderRadius="10px" h="34px" px={4}
          fontSize="13px" fontWeight="600" _hover={{ bg: '#0066d6' }}
          onClick={() => { setEditingUser(null); setDialogOpen(true) }}
        >
          <Flex align="center" gap={1.5}><LuPlus size={14} /> Nuovo Utente</Flex>
        </Button>
      </Flex>

      <Box bg="#f5f5f7" borderRadius="12px" overflow="hidden">
        {isLoading ? (
          <Flex justify="center" py={10}><Text fontSize="13px" color="#86868b">Caricamento...</Text></Flex>
        ) : users && users.length > 0 ? (
          users.map((user) => (
            <UserRow
              key={user.id}
              user={user}
              onEdit={() => { setEditingUser(user); setDialogOpen(true) }}
              onDelete={handleDelete}
            />
          ))
        ) : (
          <Flex justify="center" py={10}><Text fontSize="13px" color="#86868b">Nessun utente trovato</Text></Flex>
        )}
      </Box>

      <UserDialog user={editingUser} open={dialogOpen} onClose={() => setDialogOpen(false)} />
    </Box>
  )
}

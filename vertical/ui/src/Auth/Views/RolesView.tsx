import { Box, Button, Flex, Input, Text, Textarea } from '@chakra-ui/react'
import { useEffect, useState } from 'react'
import { LuPencil, LuPlus, LuTrash2, LuX } from 'react-icons/lu'

import {
  useRolesQuery,
  useDeleteRoleMutation,
  useCreateRoleMutation,
  useUpdateRoleMutation,
  usePermissionsQuery,
} from '../Services/Api'
import { toaster } from '@/Snippets/toaster'
import type { Role, Permission } from '../Types'

const inputStyle = {
  bg: '#f5f5f7',
  border: 'none',
  borderRadius: '10px',
  fontSize: '14px',
  _focus: { bg: 'white', boxShadow: '0 0 0 3px rgba(0,122,255,0.15)' },
}

const APP_LABELS: Record<string, string> = {
  auth: 'Autenticazione',
  orgs: 'Organizzazione',
  processes: 'Processi',
}

function PermissionChip({ perm, selected, onClick }: { perm: Permission; selected: boolean; onClick: () => void }) {
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
      title={perm.code}
    >
      {perm.description}
    </Box>
  )
}

function RoleDialog({ role, open, onClose }: { role: Role | null; open: boolean; onClose: () => void }) {
  const { data: permissions } = usePermissionsQuery()
  const [createRole] = useCreateRoleMutation()
  const [updateRole] = useUpdateRoleMutation()

  const [name, setName] = useState('')
  const [description, setDescription] = useState('')
  const [selectedPermIds, setSelectedPermIds] = useState<number[]>([])
  const [loading, setLoading] = useState(false)

  const isEdit = !!role

  useEffect(() => {
    if (open) {
      if (role) {
        setName(role.name)
        setDescription(role.description || '')
        setSelectedPermIds(role.permissions?.map((p) => p.id) ?? [])
      } else {
        setName('')
        setDescription('')
        setSelectedPermIds([])
      }
    }
  }, [open, role])

  const togglePerm = (id: number) =>
    setSelectedPermIds((prev) => (prev.includes(id) ? prev.filter((p) => p !== id) : [...prev, id]))

  const selectAllInGroup = (perms: Permission[]) => {
    const ids = perms.map((p) => p.id)
    const allSelected = ids.every((id) => selectedPermIds.includes(id))
    if (allSelected) {
      setSelectedPermIds((prev) => prev.filter((id) => !ids.includes(id)))
    } else {
      setSelectedPermIds((prev) => [...new Set([...prev, ...ids])])
    }
  }

  const handleSubmit = async () => {
    if (!name.trim()) {
      toaster.error({ title: 'Il nome del ruolo e\' obbligatorio' })
      return
    }
    setLoading(true)
    try {
      const body = { name: name.trim(), description: description.trim(), permissionIds: selectedPermIds }
      if (isEdit) {
        await updateRole({ id: role!.id, body }).unwrap()
      } else {
        await createRole(body).unwrap()
      }
      toaster.success({ title: isEdit ? 'Ruolo aggiornato' : 'Ruolo creato' })
      onClose()
    } catch {
      toaster.error({ title: 'Errore durante il salvataggio' })
    } finally {
      setLoading(false)
    }
  }

  const grouped = (permissions ?? []).reduce(
    (acc, perm) => {
      const app = perm.app || 'other'
      if (!acc[app]) acc[app] = []
      acc[app].push(perm)
      return acc
    },
    {} as Record<string, Permission[]>,
  )

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
              {isEdit ? 'Modifica Ruolo' : 'Nuovo Ruolo'}
            </Text>
            <Box as="button" p={1} color="#86868b" _hover={{ color: '#1d1d1f' }} onClick={onClose}>
              <LuX size={18} />
            </Box>
          </Flex>

          <Box px={6} pb={6}>
            <Flex direction="column" gap={3.5}>
              <Box>
                <Text fontSize="12px" fontWeight="600" color="#86868b" mb={1}>Nome *</Text>
                <Input value={name} onChange={(e) => setName(e.target.value)} placeholder="es. Amministratore" h="40px" {...inputStyle} />
              </Box>

              <Box>
                <Text fontSize="12px" fontWeight="600" color="#86868b" mb={1}>Descrizione</Text>
                <Textarea
                  value={description}
                  onChange={(e) => setDescription(e.target.value)}
                  placeholder="Descrizione del ruolo..."
                  rows={2}
                  resize="none"
                  {...inputStyle}
                />
              </Box>

              <Box>
                <Text fontSize="12px" fontWeight="600" color="#86868b" mb={3}>Permessi</Text>
                <Flex direction="column" gap={4}>
                  {Object.entries(grouped).map(([app, perms]) => (
                    <Box key={app}>
                      <Flex align="center" gap={2} mb={2}>
                        <Text fontSize="13px" fontWeight="600" color="#1d1d1f">
                          {APP_LABELS[app] || app}
                        </Text>
                        <Box
                          as="button"
                          type="button"
                          fontSize="11px"
                          color="#007aff"
                          cursor="pointer"
                          _hover={{ textDecoration: 'underline' }}
                          onClick={() => selectAllInGroup(perms)}
                        >
                          {perms.every((p) => selectedPermIds.includes(p.id)) ? 'Deseleziona tutti' : 'Seleziona tutti'}
                        </Box>
                      </Flex>
                      <Flex gap={2} wrap="wrap">
                        {perms.map((perm) => (
                          <PermissionChip
                            key={perm.id}
                            perm={perm}
                            selected={selectedPermIds.includes(perm.id)}
                            onClick={() => togglePerm(perm.id)}
                          />
                        ))}
                      </Flex>
                    </Box>
                  ))}
                </Flex>
              </Box>

              <Flex gap={2} mt={2}>
                <Button flex={1} bg="#f5f5f7" color="#1d1d1f" borderRadius="10px" h="40px" fontSize="14px" _hover={{ bg: '#e8e8ed' }} onClick={onClose}>
                  Annulla
                </Button>
                <Button flex={1} bg="#007aff" color="white" borderRadius="10px" h="40px" fontSize="14px" fontWeight="600" _hover={{ bg: '#0066d6' }} loading={loading} onClick={handleSubmit}>
                  {isEdit ? 'Salva' : 'Crea Ruolo'}
                </Button>
              </Flex>
            </Flex>
          </Box>
        </Box>
      </Flex>
    </Box>
  )
}

function RoleRow({ role, onEdit, onDelete }: { role: Role; onEdit: () => void; onDelete: (id: number) => void }) {
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
        <Text fontSize="14px" fontWeight="500" color="#1d1d1f">{role.name}</Text>
        <Text fontSize="12px" color="#86868b">{role.description || 'Nessuna descrizione'}</Text>
      </Box>
      <Flex align="center" gap={2}>
        <Box px={2} py={0.5} borderRadius="6px" bg="#f5f5f7">
          <Text fontSize="11px" color="#86868b">
            {role.permissions?.length ?? 0} permessi
          </Text>
        </Box>
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
        <Box
          as="button"
          p={1.5}
          borderRadius="6px"
          color="#86868b"
          cursor="pointer"
          _hover={{ color: '#ff3b30', bg: '#ff3b3010' }}
          transition="all 0.15s"
          onClick={(e) => { e.stopPropagation(); onDelete(role.id) }}
        >
          <LuTrash2 size={15} />
        </Box>
      </Flex>
    </Flex>
  )
}

export default function RolesView() {
  const { data: roles, isLoading } = useRolesQuery()
  const [deleteRole] = useDeleteRoleMutation()
  const [dialogOpen, setDialogOpen] = useState(false)
  const [editingRole, setEditingRole] = useState<Role | null>(null)

  const handleDelete = async (id: number) => {
    if (confirm('Sei sicuro di voler eliminare questo ruolo?')) {
      try {
        await deleteRole(id).unwrap()
        toaster.success({ title: 'Ruolo eliminato' })
      } catch {
        toaster.error({ title: "Errore durante l'eliminazione" })
      }
    }
  }

  return (
    <Box>
      <Flex justify="space-between" align="center" mb={4}>
        <Text fontSize="15px" fontWeight="700" color="#1d1d1f">Ruoli</Text>
        <Button
          bg="#007aff" color="white" borderRadius="10px" h="34px" px={4}
          fontSize="13px" fontWeight="600" _hover={{ bg: '#0066d6' }}
          onClick={() => { setEditingRole(null); setDialogOpen(true) }}
        >
          <Flex align="center" gap={1.5}><LuPlus size={14} /> Nuovo Ruolo</Flex>
        </Button>
      </Flex>

      <Box bg="#f5f5f7" borderRadius="12px" overflow="hidden">
        {isLoading ? (
          <Flex justify="center" py={10}><Text fontSize="13px" color="#86868b">Caricamento...</Text></Flex>
        ) : roles && roles.length > 0 ? (
          roles.map((role) => (
            <RoleRow
              key={role.id}
              role={role}
              onEdit={() => { setEditingRole(role); setDialogOpen(true) }}
              onDelete={handleDelete}
            />
          ))
        ) : (
          <Flex justify="center" py={10}><Text fontSize="13px" color="#86868b">Nessun ruolo trovato</Text></Flex>
        )}
      </Box>

      <RoleDialog role={editingRole} open={dialogOpen} onClose={() => setDialogOpen(false)} />
    </Box>
  )
}

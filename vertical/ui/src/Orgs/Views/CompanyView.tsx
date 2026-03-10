import { Box, Button, Flex, Input, Text } from '@chakra-ui/react'
import { useState } from 'react'
import { LuBuilding2, LuChevronRight, LuPencil, LuPlus, LuTrash2, LuX } from 'react-icons/lu'

import { useCurrentUser } from '@/Auth/Hooks'
import {
  useCompaniesQuery,
  useCreateCompanyMutation,
  useUpdateCompanyMutation,
  useDeleteCompanyMutation,
  useCompanyAreasQuery,
  useCreateCompanyAreaMutation,
} from '../Services/Api'
import {
  useUpdateAreaMutation,
  useDeleteAreaMutation,
} from '../Services/AreaApi'
import { toaster } from '@/Snippets/toaster'
import type { Company } from '../Types'
import type { Area } from '../Types/Area'

const inputStyle = {
  bg: '#f5f5f7',
  border: 'none',
  borderRadius: '10px',
  fontSize: '14px',
  _focus: { bg: 'white', boxShadow: '0 0 0 3px rgba(0,122,255,0.15)' },
}

function AreaItem({ area, onUpdate, onDelete }: { area: Area; onUpdate: (name: string) => void; onDelete: () => void }) {
  const [editing, setEditing] = useState(false)
  const [name, setName] = useState(area.name)

  const handleSave = () => {
    if (name.trim() && name.trim() !== area.name) onUpdate(name.trim())
    setEditing(false)
  }

  return (
    <Flex align="center" justify="space-between" py={2} px={3} borderBottom="1px solid #f0f0f2" _hover={{ bg: '#fafafa' }}>
      {editing ? (
        <Input value={name} onChange={(e) => setName(e.target.value)} onKeyDown={(e) => e.key === 'Enter' && handleSave()} onBlur={handleSave} size="sm" bg="#f5f5f7" border="none" borderRadius="8px" fontSize="13px" autoFocus flex={1} mr={2} />
      ) : (
        <Text fontSize="13px" color="#1d1d1f">{area.name}</Text>
      )}
      <Flex gap={1} flexShrink={0}>
        <Box as="button" p={1} borderRadius="6px" color="#86868b" cursor="pointer" _hover={{ color: '#007aff' }} onClick={() => { setName(area.name); setEditing(!editing) }}><LuPencil size={13} /></Box>
        <Box as="button" p={1} borderRadius="6px" color="#86868b" cursor="pointer" _hover={{ color: '#ff3b30' }} onClick={onDelete}><LuTrash2 size={13} /></Box>
      </Flex>
    </Flex>
  )
}

function ParentSelector({ companies, currentId, parentId, onChange }: { companies: Company[]; currentId: number | null; parentId: number | null; onChange: (id: number | null) => void }) {
  const [open, setOpen] = useState(false)
  const available = companies.filter((c) => c.id !== currentId)
  const parentName = available.find((c) => c.id === parentId)?.name

  return (
    <Box position="relative">
      <Text fontSize="12px" fontWeight="600" color="#86868b" mb={1}>Azienda Padre (Holding)</Text>
      <Flex
        align="center"
        justify="space-between"
        px={3}
        h="40px"
        bg="#f5f5f7"
        borderRadius="10px"
        cursor="pointer"
        onClick={() => setOpen(!open)}
      >
        <Text fontSize="14px" color={parentId ? '#1d1d1f' : '#86868b'}>
          {parentName ?? 'Nessuna (top-level)'}
        </Text>
        <LuChevronRight size={14} color="#86868b" style={{ transform: open ? 'rotate(90deg)' : 'none', transition: '0.15s' }} />
      </Flex>
      {open && (
        <>
          <Box position="fixed" inset={0} zIndex={99} onClick={() => setOpen(false)} />
          <Box
            position="absolute"
            top="100%"
            left={0}
            right={0}
            mt={1}
            bg="white"
            borderRadius="10px"
            boxShadow="0 8px 30px rgba(0,0,0,0.12)"
            border="1px solid #f0f0f2"
            zIndex={100}
            overflow="hidden"
            maxH="200px"
            overflowY="auto"
          >
            <Flex
              px={3} py={2} cursor="pointer" fontSize="13px" color="#86868b"
              _hover={{ bg: '#f5f5f7' }}
              fontWeight={!parentId ? '600' : '400'}
              onClick={() => { onChange(null); setOpen(false) }}
            >
              Nessuna (top-level)
            </Flex>
            {available.map((c) => (
              <Flex
                key={c.id}
                px={3} py={2} cursor="pointer" fontSize="13px" color="#1d1d1f"
                _hover={{ bg: '#f5f5f7' }}
                fontWeight={c.id === parentId ? '600' : '400'}
                onClick={() => { onChange(c.id); setOpen(false) }}
              >
                {c.name}
              </Flex>
            ))}
          </Box>
        </>
      )}
    </Box>
  )
}

function CompanyDialog({ company, allCompanies, open, onClose }: { company: Company | null; allCompanies: Company[]; open: boolean; onClose: () => void }) {
  const [createCompany] = useCreateCompanyMutation()
  const [updateCompany] = useUpdateCompanyMutation()
  const [createCompanyArea] = useCreateCompanyAreaMutation()
  const [updateArea] = useUpdateAreaMutation()
  const [deleteArea] = useDeleteAreaMutation()

  const [name, setName] = useState('')
  const [parentId, setParentId] = useState<number | null>(null)
  const [newAreaName, setNewAreaName] = useState('')
  const [loading, setLoading] = useState(false)
  const [createdCompanyId, setCreatedCompanyId] = useState<number | null>(null)

  const isEdit = !!company || !!createdCompanyId
  const activeCompanyId = company?.id ?? createdCompanyId
  const { data: areas } = useCompanyAreasQuery(activeCompanyId!, { skip: !activeCompanyId })

  if (open && name === '' && company?.name) {
    setName(company.name)
    setParentId(company.parentId ?? null)
  }

  const handleSaveCompany = async () => {
    if (!name.trim()) {
      toaster.error({ title: 'Il nome e\' obbligatorio' })
      return
    }
    setLoading(true)
    try {
      if (activeCompanyId) {
        await updateCompany({ id: activeCompanyId, body: { name: name.trim(), parentId } }).unwrap()
        toaster.success({ title: 'Azienda aggiornata' })
      } else {
        const created = await createCompany({ name: name.trim(), parentId }).unwrap()
        setCreatedCompanyId(created.id)
        toaster.success({ title: 'Azienda creata. Ora puoi aggiungere le aree.' })
      }
    } catch {
      toaster.error({ title: 'Errore durante il salvataggio' })
    } finally {
      setLoading(false)
    }
  }

  const handleAddArea = async () => {
    if (!newAreaName.trim() || !activeCompanyId) return
    try {
      await createCompanyArea({ companyId: activeCompanyId, name: newAreaName.trim() }).unwrap()
      setNewAreaName('')
    } catch {
      toaster.error({ title: 'Errore durante la creazione dell\'area' })
    }
  }

  const handleClose = () => {
    setName('')
    setParentId(null)
    setNewAreaName('')
    setCreatedCompanyId(null)
    onClose()
  }

  if (!open) return null

  return (
    <Box position="fixed" top={0} left={0} right={0} bottom={0} zIndex={1000}>
      <Box position="absolute" inset={0} bg="blackAlpha.400" onClick={handleClose} />
      <Flex position="absolute" inset={0} align="center" justify="center" p={4}>
        <Box bg="white" borderRadius="16px" boxShadow="0 20px 60px rgba(0,0,0,0.15)" maxW="480px" w="full" maxH="85vh" overflow="auto" onClick={(e) => e.stopPropagation()}>
          <Flex justify="space-between" align="center" px={6} pt={5} pb={3}>
            <Text fontSize="17px" fontWeight="700" color="#1d1d1f">
              {isEdit ? 'Modifica Azienda' : 'Nuova Azienda'}
            </Text>
            <Box as="button" p={1} color="#86868b" _hover={{ color: '#1d1d1f' }} onClick={handleClose}><LuX size={18} /></Box>
          </Flex>

          <Box px={6} pb={6}>
            <Flex direction="column" gap={4}>
              <Box>
                <Text fontSize="12px" fontWeight="600" color="#86868b" mb={1}>Nome Azienda *</Text>
                <Flex gap={2}>
                  <Input value={name} onChange={(e) => setName(e.target.value)} onKeyDown={(e) => e.key === 'Enter' && handleSaveCompany()} placeholder="es. Nuova S.p.A." h="40px" flex={1} {...inputStyle} autoFocus={!isEdit} />
                  <Button bg="#007aff" color="white" borderRadius="10px" h="40px" px={4} fontSize="13px" fontWeight="600" _hover={{ bg: '#0066d6' }} loading={loading} onClick={handleSaveCompany}>
                    {isEdit ? 'Salva' : 'Crea'}
                  </Button>
                </Flex>
              </Box>

              {allCompanies.length > 0 && (
                <ParentSelector
                  companies={allCompanies}
                  currentId={activeCompanyId}
                  parentId={parentId}
                  onChange={setParentId}
                />
              )}

              {activeCompanyId && (
                <Box>
                  <Text fontSize="12px" fontWeight="600" color="#86868b" mb={2}>Aree Aziendali</Text>
                  <Box bg="#f5f5f7" borderRadius="10px" overflow="hidden">
                    {areas && areas.length > 0 ? (
                      areas.map((area) => (
                        <AreaItem
                          key={area.id}
                          area={area}
                          onUpdate={(n) => updateArea({ id: area.id, body: { name: n } })}
                          onDelete={() => { if (confirm('Eliminare questa area?')) deleteArea(area.id) }}
                        />
                      ))
                    ) : (
                      <Flex justify="center" py={4}><Text fontSize="12px" color="#86868b">Nessuna area</Text></Flex>
                    )}
                    <Flex gap={2} p={3} borderTop={areas && areas.length > 0 ? '1px solid #e8e8ed' : 'none'}>
                      <Input placeholder="Nuova area..." value={newAreaName} onChange={(e) => setNewAreaName(e.target.value)} onKeyDown={(e) => e.key === 'Enter' && handleAddArea()} bg="white" border="none" borderRadius="8px" h="34px" fontSize="13px" flex={1} />
                      <Button onClick={handleAddArea} bg="#007aff" color="white" borderRadius="8px" h="34px" px={3} fontSize="12px" _hover={{ bg: '#0066d6' }} disabled={!newAreaName.trim()}>
                        <LuPlus size={14} />
                      </Button>
                    </Flex>
                  </Box>
                </Box>
              )}

              {activeCompanyId && (
                <Button bg="#f5f5f7" color="#1d1d1f" borderRadius="10px" h="40px" fontSize="14px" _hover={{ bg: '#e8e8ed' }} onClick={handleClose} mt={1}>
                  Chiudi
                </Button>
              )}
            </Flex>
          </Box>
        </Box>
      </Flex>
    </Box>
  )
}

function CompanyTree({ companies, onEdit, onDelete, level = 0 }: { companies: Company[]; onEdit: (c: Company) => void; onDelete: (id: number) => void; level?: number; allCompanies?: Company[] }) {
  const roots = companies.filter((c) => level === 0 ? !c.parentId : false)
  const items = level === 0 ? roots : companies

  if (items.length === 0 && level === 0) {
    return <Flex justify="center" py={10}><Text fontSize="13px" color="#86868b">Nessuna azienda</Text></Flex>
  }

  return (
    <>
      {items.map((company) => {
        const children = level === 0 ? companies.filter((c) => c.parentId === company.id) : []
        return (
          <Box key={company.id}>
            <Flex
              align="center"
              justify="space-between"
              py={3}
              px={4}
              pl={4 + level * 6}
              borderBottom="1px solid #f0f0f2"
              _hover={{ bg: '#fafafa' }}
              transition="background 0.15s"
              cursor="pointer"
              onClick={() => onEdit(company)}
            >
              <Flex align="center" gap={3}>
                <Flex w="32px" h="32px" align="center" justify="center" borderRadius="8px" bg={level === 0 ? '#007aff10' : '#f5f5f7'} color={level === 0 ? '#007aff' : '#86868b'} flexShrink={0}>
                  <LuBuilding2 size={15} />
                </Flex>
                <Box>
                  <Text fontSize="14px" fontWeight={level === 0 ? '600' : '400'} color="#1d1d1f">{company.name}</Text>
                  {children.length > 0 && (
                    <Text fontSize="11px" color="#86868b">{children.length} sotto-aziend{children.length === 1 ? 'a' : 'e'}</Text>
                  )}
                </Box>
              </Flex>
              <Flex align="center" gap={1}>
                <Box as="button" p={1.5} borderRadius="6px" color="#86868b" _hover={{ color: '#007aff' }} onClick={(e) => { e.stopPropagation(); onEdit(company) }}>
                  <LuPencil size={14} />
                </Box>
                <Box as="button" p={1.5} borderRadius="6px" color="#86868b" _hover={{ color: '#ff3b30' }} onClick={(e) => { e.stopPropagation(); onDelete(company.id) }}>
                  <LuTrash2 size={14} />
                </Box>
              </Flex>
            </Flex>
            {children.length > 0 && <CompanyTree companies={children} onEdit={onEdit} onDelete={onDelete} level={level + 1} />}
          </Box>
        )
      })}
    </>
  )
}

export default function CompanyView() {
  const user = useCurrentUser()
  const isSuperuser = user?.isSuperuser ?? false
  const { data: companies, isLoading } = useCompaniesQuery()
  const [deleteCompany] = useDeleteCompanyMutation()
  const [dialogOpen, setDialogOpen] = useState(false)
  const [editingCompany, setEditingCompany] = useState<Company | null>(null)

  const handleDelete = async (id: number) => {
    const children = companies?.filter((c) => c.parentId === id)
    const msg = children && children.length > 0
      ? `Questa azienda ha ${children.length} sotto-aziend${children.length === 1 ? 'a' : 'e'}. Eliminare comunque?`
      : 'Sei sicuro di voler eliminare questa azienda?'
    if (confirm(msg)) {
      try {
        await deleteCompany(id).unwrap()
        toaster.success({ title: 'Azienda eliminata' })
      } catch {
        toaster.error({ title: "Errore durante l'eliminazione" })
      }
    }
  }

  return (
    <Box>
      <Flex justify="space-between" align="center" mb={4}>
        <Text fontSize="15px" fontWeight="700" color="#1d1d1f">Aziende</Text>
        {isSuperuser && (
          <Button bg="#007aff" color="white" borderRadius="10px" h="34px" px={4} fontSize="13px" fontWeight="600" _hover={{ bg: '#0066d6' }} onClick={() => { setEditingCompany(null); setDialogOpen(true) }}>
            <Flex align="center" gap={1.5}><LuPlus size={14} /> Nuova Azienda</Flex>
          </Button>
        )}
      </Flex>

      <Box bg="#f5f5f7" borderRadius="12px" overflow="hidden">
        {isLoading ? (
          <Flex justify="center" py={10}><Text fontSize="13px" color="#86868b">Caricamento...</Text></Flex>
        ) : companies && companies.length > 0 ? (
          <CompanyTree
            companies={companies}
            onEdit={(c) => { setEditingCompany(c); setDialogOpen(true) }}
            onDelete={handleDelete}
          />
        ) : (
          <Flex justify="center" py={10}><Text fontSize="13px" color="#86868b">Nessuna azienda</Text></Flex>
        )}
      </Box>

      <CompanyDialog company={editingCompany} allCompanies={companies ?? []} open={dialogOpen} onClose={() => { setDialogOpen(false); setEditingCompany(null) }} />
    </Box>
  )
}

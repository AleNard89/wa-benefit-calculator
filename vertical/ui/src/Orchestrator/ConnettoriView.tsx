import { Badge, Box, Button, Flex, Input, Spinner, Text } from '@chakra-ui/react'
import { useState } from 'react'
import { useSelector } from 'react-redux'
import { LuBot, LuBuilding2, LuPlus, LuRefreshCw, LuPencil, LuPlugZap, LuTrash2 } from 'react-icons/lu'

import { selectCompanyId } from '@/Orgs/Redux'
import { useCompaniesQuery } from '@/Orgs/Services/Api'
import { toaster } from '@/Snippets/toaster'
import {
  useConnectorsQuery,
  useCreateConnectorMutation,
  useUpdateConnectorMutation,
  useTestConnectorMutation,
  useSyncConnectorMutation,
} from './api'
import type { ConnectorForm, ConnectorResponse, UiPathFolder } from './types'

function EmptyState({ onAdd }: { onAdd: () => void }) {
  return (
    <Flex direction="column" align="center" justify="center" py={12} gap={4}>
      <Flex w="56px" h="56px" align="center" justify="center" borderRadius="14px" bg="#007aff10" color="#007aff">
        <LuBot size={28} />
      </Flex>
      <Text fontSize="15px" fontWeight="600" color="#1d1d1f">Nessun connettore configurato</Text>
      <Text fontSize="13px" color="#86868b" textAlign="center" maxW="320px">
        Aggiungi un connettore per sincronizzare esecuzioni e code dal tuo Orchestrator.
      </Text>
      <Button size="sm" colorPalette="brand" onClick={onAdd}>
        <LuPlus size={14} /> Aggiungi Connettore
      </Button>
    </Flex>
  )
}

function ConnectorFormFields({
  form,
  onChange,
  isEdit,
}: {
  form: ConnectorForm
  onChange: (f: ConnectorForm) => void
  isEdit: boolean
}) {
  return (
    <Flex direction="column" gap={3}>
      <Box>
        <Text fontSize="12px" fontWeight="600" color="#86868b" mb={1}>Nome</Text>
        <Input size="sm" value={form.name} onChange={(e) => onChange({ ...form, name: e.target.value })} placeholder="Es. Connettore Produzione" />
      </Box>
      <Box>
        <Text fontSize="12px" fontWeight="600" color="#86868b" mb={1}>Tipo</Text>
        <Input size="sm" value="UiPath" disabled />
      </Box>
      <Box>
        <Text fontSize="12px" fontWeight="600" color="#86868b" mb={1}>Organization Name</Text>
        <Input size="sm" value={form.organizationName} onChange={(e) => onChange({ ...form, organizationName: e.target.value })} placeholder="nomeorganizzazione" />
      </Box>
      <Box>
        <Text fontSize="12px" fontWeight="600" color="#86868b" mb={1}>Tenant Name</Text>
        <Input size="sm" value={form.tenantName} onChange={(e) => onChange({ ...form, tenantName: e.target.value })} placeholder="NomeTenant" />
      </Box>
      <Box>
        <Text fontSize="12px" fontWeight="600" color="#86868b" mb={1}>
          Personal Access Token {isEdit && <Text as="span" fontSize="11px" color="#86868b">(lascia vuoto per mantenere il precedente)</Text>}
        </Text>
        <Input size="sm" type="password" value={form.accessToken} onChange={(e) => onChange({ ...form, accessToken: e.target.value })} placeholder={isEdit ? '••••••••' : 'rt_...'} />
      </Box>
      <Box>
        <Flex align="center" justify="space-between" mb={2}>
          <Text fontSize="12px" fontWeight="600" color="#86868b">Folders</Text>
          <Button size="xs" variant="outline" onClick={() => onChange({ ...form, folders: [...form.folders, { id: '', name: '' }] })}>
            <LuPlus size={12} /> Aggiungi Folder
          </Button>
        </Flex>
        {form.folders.map((folder: UiPathFolder, idx: number) => (
          <Flex key={idx} gap={2} mb={2} align="center">
            <Input size="sm" flex={1} placeholder="Folder ID (es. 1231441)" value={folder.id}
              onChange={(e) => { const f = [...form.folders]; f[idx] = { ...f[idx], id: e.target.value }; onChange({ ...form, folders: f }) }} />
            <Input size="sm" flex={1} placeholder="Nome (es. Shared)" value={folder.name}
              onChange={(e) => { const f = [...form.folders]; f[idx] = { ...f[idx], name: e.target.value }; onChange({ ...form, folders: f }) }} />
            <Button size="xs" variant="ghost" colorPalette="red" onClick={() => { const f = form.folders.filter((_: UiPathFolder, i: number) => i !== idx); onChange({ ...form, folders: f }) }}>
              <LuTrash2 size={13} />
            </Button>
          </Flex>
        ))}
        {form.folders.length === 0 && <Text fontSize="12px" color="#86868b">Nessun folder configurato. Aggiungi almeno un folder.</Text>}
      </Box>
    </Flex>
  )
}

function ConnectorRow({
  connector,
  onEdit,
  onTest,
  onSync,
  isTesting,
  isSyncing,
}: {
  connector: ConnectorResponse
  onEdit: () => void
  onTest: () => void
  onSync: () => void
  isTesting: boolean
  isSyncing: boolean
}) {
  return (
    <Flex
      align="center"
      justify="space-between"
      p={4}
      borderRadius="12px"
      bg="white"
      boxShadow="0 1px 4px rgba(0,0,0,0.06)"
    >
      <Flex align="center" gap={3}>
        <Flex w="40px" h="40px" align="center" justify="center" borderRadius="10px" bg="#007aff10" color="#007aff" flexShrink={0}>
          <LuBot size={20} />
        </Flex>
        <Box>
          <Flex align="center" gap={2}>
            <Text fontSize="14px" fontWeight="600" color="#1d1d1f">{connector.name}</Text>
            <Badge colorPalette={connector.isActive ? 'green' : 'gray'} fontSize="10px">
              {connector.isActive ? 'Attivo' : 'Inattivo'}
            </Badge>
            <Badge colorPalette="blue" fontSize="10px">{connector.type}</Badge>
          </Flex>
          <Text fontSize="12px" color="#86868b">
            {connector.organizationName} / {connector.tenantName} — {connector.folders?.length > 0
              ? `Folders: ${connector.folders.map((f) => f.name).join(', ')}`
              : `Folder: ${connector.folderName}`}
          </Text>
        </Box>
      </Flex>
      <Flex gap={2}>
        <Button size="xs" variant="outline" onClick={onTest} loading={isTesting}>
          <LuPlugZap size={13} /> Testa
        </Button>
        <Button size="xs" variant="outline" onClick={onSync} loading={isSyncing}>
          <LuRefreshCw size={13} /> Sincronizza
        </Button>
        <Button size="xs" variant="outline" onClick={onEdit}>
          <LuPencil size={13} /> Modifica
        </Button>
      </Flex>
    </Flex>
  )
}

const emptyForm: ConnectorForm = {
  name: '',
  type: 'UIPATH',
  organizationName: '',
  tenantName: '',
  accessToken: '',
  folderId: '',
  folderName: '',
  folders: [{ id: '', name: '' }],
}

export default function ConnettoriView() {
  const currentCompanyId = useSelector(selectCompanyId)
  const { data: allCompanies } = useCompaniesQuery()
  const currentCompany = allCompanies?.find((c) => c.id.toString() === currentCompanyId)

  const { data: connectors, isLoading } = useConnectorsQuery()
  const [createConnector, { isLoading: isCreating }] = useCreateConnectorMutation()
  const [updateConnector, { isLoading: isUpdating }] = useUpdateConnectorMutation()
  const [testConnector] = useTestConnectorMutation()
  const [syncConnector] = useSyncConnectorMutation()

  const [showForm, setShowForm] = useState(false)
  const [editingId, setEditingId] = useState<number | null>(null)
  const [form, setForm] = useState<ConnectorForm>(emptyForm)
  const [testingId, setTestingId] = useState<number | null>(null)
  const [syncingId, setSyncingId] = useState<number | null>(null)

  const handleAdd = () => {
    setForm(emptyForm)
    setEditingId(null)
    setShowForm(true)
  }

  const handleEdit = (c: ConnectorResponse) => {
    setForm({
      name: c.name,
      type: c.type,
      organizationName: c.organizationName,
      tenantName: c.tenantName,
      accessToken: '',
      folderId: c.folderId,
      folderName: c.folderName,
      folders: c.folders?.length > 0 ? c.folders : [{ id: c.folderId, name: c.folderName }],
    })
    setEditingId(c.id)
    setShowForm(true)
  }

  const handleSave = async () => {
    try {
      if (editingId) {
        await updateConnector({ id: editingId, body: form }).unwrap()
        toaster.success({ title: 'Connettore aggiornato' })
      } else {
        await createConnector(form).unwrap()
        toaster.success({ title: 'Connettore creato' })
      }
      setShowForm(false)
      setEditingId(null)
    } catch {
      toaster.error({ title: 'Errore nel salvataggio' })
    }
  }

  const handleTest = async (id: number) => {
    setTestingId(id)
    try {
      await testConnector(id).unwrap()
      toaster.success({ title: 'Connessione riuscita' })
    } catch (err: any) {
      toaster.error({ title: err?.data?.message || 'Connessione fallita' })
    } finally {
      setTestingId(null)
    }
  }

  const handleSync = async (id: number) => {
    setSyncingId(id)
    try {
      await syncConnector(id).unwrap()
      toaster.success({ title: 'Sincronizzazione completata' })
    } catch (err: any) {
      toaster.error({ title: err?.data?.message || 'Sincronizzazione fallita' })
    } finally {
      setSyncingId(null)
    }
  }

  if (isLoading) {
    return <Flex justify="center" py={10}><Spinner size="lg" color="brand.300" /></Flex>
  }

  return (
    <Flex direction="column" gap={4}>
      <Flex justify="space-between" align="center">
        <Box>
          <Text fontWeight="700" fontSize="17px" color="#1d1d1f">Connettori</Text>
          {currentCompany && (
            <Flex align="center" gap={1.5} mt={1}>
              <LuBuilding2 size={12} color="#86868b" />
              <Text fontSize="12px" color="#86868b">
                Stai configurando connettori per <Text as="span" fontWeight="600" color="#1d1d1f">{currentCompany.name}</Text>
              </Text>
            </Flex>
          )}
        </Box>
        {connectors && connectors.length > 0 && !showForm && (
          <Button size="sm" colorPalette="brand" onClick={handleAdd}>
            <LuPlus size={14} /> Aggiungi
          </Button>
        )}
      </Flex>

      {showForm && (
        <Box p={5} borderRadius="12px" bg="#f9f9fb" border="1px solid #e8e8ed">
          <Text fontSize="15px" fontWeight="600" color="#1d1d1f" mb={3}>
            {editingId ? 'Modifica Connettore' : 'Nuovo Connettore'}
          </Text>
          <ConnectorFormFields form={form} onChange={setForm} isEdit={!!editingId} />
          <Flex gap={2} mt={4}>
            <Button size="sm" colorPalette="brand" onClick={handleSave} loading={isCreating || isUpdating}>
              {editingId ? 'Salva' : 'Crea'}
            </Button>
            <Button size="sm" variant="outline" onClick={() => { setShowForm(false); setEditingId(null) }}>
              Annulla
            </Button>
          </Flex>
        </Box>
      )}

      {(!connectors || connectors.length === 0) && !showForm ? (
        <EmptyState onAdd={handleAdd} />
      ) : (
        <Flex direction="column" gap={2}>
          {connectors?.map((c) => (
            <ConnectorRow
              key={c.id}
              connector={c}
              onEdit={() => handleEdit(c)}
              onTest={() => handleTest(c.id)}
              onSync={() => handleSync(c.id)}
              isTesting={testingId === c.id}
              isSyncing={syncingId === c.id}
            />
          ))}
        </Flex>
      )}
    </Flex>
  )
}

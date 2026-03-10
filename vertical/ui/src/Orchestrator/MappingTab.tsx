import { Badge, Box, Button, Flex, Spinner, Text } from '@chakra-ui/react'
import { useState } from 'react'
import { LuPlus, LuTrash2 } from 'react-icons/lu'
import { NativeSelect } from '@chakra-ui/react'

import { toaster } from '@/Snippets/toaster'
import {
  useProcessQueueMapsQuery,
  useCreateProcessQueueMapMutation,
  useDeleteProcessQueueMapMutation,
  useConnectorsQuery,
} from './api'
import type { ProcessQueueMap } from './types'

function MappingRow({ mapping, onDelete }: { mapping: ProcessQueueMap; onDelete: (id: number) => void }) {
  return (
    <Flex
      align="center"
      py={3}
      px={2}
      borderBottom="1px solid #f0f0f2"
      _hover={{ bg: '#fafafa' }}
      transition="background 0.1s"
      gap={3}
      fontSize="13px"
    >
      <Box flex={2} fontWeight="500" color="#1d1d1f">
        {mapping.processName}
      </Box>
      <Box flex={0.3} color="#86868b" textAlign="center">→</Box>
      <Box flex={2} fontWeight="500" color="#1d1d1f">
        {mapping.queueName}
      </Box>
      <Box flex={0.8}>
        <Badge colorPalette={mapping.autoDetected ? 'blue' : 'green'} fontSize="10px">
          {mapping.autoDetected ? 'Auto' : 'Manuale'}
        </Badge>
      </Box>
      <Box flex={0.5} textAlign="right">
        <Button size="xs" colorPalette="red" variant="ghost" onClick={() => onDelete(mapping.id)}>
          <LuTrash2 size={13} />
        </Button>
      </Box>
    </Flex>
  )
}

export default function MappingTab() {
  const { data: maps, isLoading } = useProcessQueueMapsQuery()
  const { data: connectors } = useConnectorsQuery()
  const [createMap] = useCreateProcessQueueMapMutation()
  const [deleteMap] = useDeleteProcessQueueMapMutation()
  const [showForm, setShowForm] = useState(false)
  const [formProcess, setFormProcess] = useState('')
  const [formQueue, setFormQueue] = useState('')
  const [formConnector, setFormConnector] = useState(0)

  const handleCreate = async () => {
    if (!formProcess || !formQueue || !formConnector) {
      toaster.error({ title: 'Compila tutti i campi' })
      return
    }
    try {
      await createMap({ connectorId: formConnector, processName: formProcess, queueName: formQueue }).unwrap()
      toaster.success({ title: 'Mapping creato' })
      setShowForm(false)
      setFormProcess('')
      setFormQueue('')
    } catch {
      toaster.error({ title: 'Errore creazione mapping' })
    }
  }

  const handleDelete = async (id: number) => {
    try {
      await deleteMap(id).unwrap()
      toaster.success({ title: 'Mapping eliminato' })
    } catch {
      toaster.error({ title: 'Errore eliminazione' })
    }
  }

  // Extract distinct process names and queue names from existing maps
  const processNames = [...new Set(maps?.map((m) => m.processName) ?? [])]
  const queueNames = [...new Set(maps?.map((m) => m.queueName) ?? [])]
  const autoCount = maps?.filter((m) => m.autoDetected).length ?? 0
  const manualCount = maps?.filter((m) => !m.autoDetected).length ?? 0

  return (
    <Flex direction="column" gap={4}>
      <Flex gap={3} wrap="wrap">
        <StatCard label="Totale" value={maps?.length ?? 0} color="#1d1d1f" />
        <StatCard label="Auto-rilevati" value={autoCount} color="#007aff" />
        <StatCard label="Manuali" value={manualCount} color="#34c759" />
      </Flex>

      <Flex justify="space-between" align="center">
        <Text fontSize="13px" color="#86868b">
          Associa i bot (processi) alle code che alimentano. Le associazioni automatiche vengono create durante la sincronizzazione.
        </Text>
        <Button size="sm" colorPalette="blue" variant="outline" onClick={() => setShowForm(!showForm)}>
          <LuPlus size={14} /> Nuovo mapping
        </Button>
      </Flex>

      {showForm && (
        <Flex bg="#f9f9fb" borderRadius="12px" p={4} gap={3} align="end" wrap="wrap">
          {connectors && connectors.length > 0 && (
            <Box>
              <Text fontSize="11px" color="#86868b" mb={1}>Connettore</Text>
              <NativeSelect.Root size="sm" maxW="200px">
                <NativeSelect.Field value={formConnector || ''} onChange={(e) => setFormConnector(parseInt(e.target.value) || 0)}>
                  <option value="">Seleziona...</option>
                  {connectors.map((c) => <option key={c.id} value={c.id}>{c.name}</option>)}
                </NativeSelect.Field>
              </NativeSelect.Root>
            </Box>
          )}
          <Box flex={1} minW="200px">
            <Text fontSize="11px" color="#86868b" mb={1}>Nome Processo (Bot)</Text>
            <input
              type="text"
              value={formProcess}
              onChange={(e) => setFormProcess(e.target.value)}
              placeholder="es. NomeProcesso_PRF"
              list="process-names-list"
              style={{
                width: '100%', padding: '6px 10px', fontSize: '13px',
                border: '1px solid #d2d2d7', borderRadius: '8px', outline: 'none',
              }}
            />
            {processNames.length > 0 && (
              <datalist id="process-names-list">
                {processNames.map((n) => <option key={n} value={n} />)}
              </datalist>
            )}
          </Box>
          <Box flex={1} minW="200px">
            <Text fontSize="11px" color="#86868b" mb={1}>Nome Coda</Text>
            <input
              type="text"
              value={formQueue}
              onChange={(e) => setFormQueue(e.target.value)}
              placeholder="es. NomeCoda_Input"
              list="queue-names-list"
              style={{
                width: '100%', padding: '6px 10px', fontSize: '13px',
                border: '1px solid #d2d2d7', borderRadius: '8px', outline: 'none',
              }}
            />
            {queueNames.length > 0 && (
              <datalist id="queue-names-list">
                {queueNames.map((n) => <option key={n} value={n} />)}
              </datalist>
            )}
          </Box>
          <Flex gap={2}>
            <Button size="sm" colorPalette="blue" onClick={handleCreate}>Salva</Button>
            <Button size="sm" variant="outline" onClick={() => setShowForm(false)}>Annulla</Button>
          </Flex>
        </Flex>
      )}

      {isLoading ? (
        <Flex justify="center" py={10}><Spinner size="lg" color="brand.300" /></Flex>
      ) : !maps || maps.length === 0 ? (
        <Text fontSize="13px" color="#86868b" textAlign="center" py={10}>
          Nessun mapping trovato. Sincronizza un connettore per generare mapping automatici, oppure creali manualmente.
        </Text>
      ) : (
        <>
          <Flex px={2} py={2} borderBottom="2px solid #f0f0f2" gap={3} fontSize="11px" fontWeight="700" color="#86868b" textTransform="uppercase" letterSpacing="0.5px">
            <Box flex={2}>Processo (Bot)</Box>
            <Box flex={0.3}></Box>
            <Box flex={2}>Coda</Box>
            <Box flex={0.8}>Tipo</Box>
            <Box flex={0.5}></Box>
          </Flex>

          {maps.map((m) => <MappingRow key={m.id} mapping={m} onDelete={handleDelete} />)}
        </>
      )}
    </Flex>
  )
}

function StatCard({ label, value, color }: { label: string; value: number; color: string }) {
  return (
    <Box bg="#f9f9fb" borderRadius="10px" p={3} flex={1} minW="120px">
      <Text fontSize="11px" color="#86868b">{label}</Text>
      <Text fontSize="22px" fontWeight="700" color={color}>{value}</Text>
    </Box>
  )
}

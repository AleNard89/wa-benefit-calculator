import { Badge, Box, Button, Flex, Input, Spinner, Text } from '@chakra-ui/react'
import { useState } from 'react'
import { LuChevronLeft, LuChevronRight } from 'react-icons/lu'

import { NativeSelect } from '@chakra-ui/react'
import { useQueueItemsQuery, useConnectorsQuery } from './api'
import type { QueueItem, QueueItemFilters } from './types'

const statusColors: Record<string, string> = {
  Successful: 'green',
  Failed: 'red',
  New: 'blue',
  InProgress: 'blue',
  Retried: 'orange',
  Abandoned: 'gray',
  Deleted: 'gray',
}

const statusLabels: Record<string, string> = {
  Successful: 'Completato',
  Failed: 'Fallito',
  New: 'Nuovo',
  InProgress: 'In corso',
  Retried: 'Ritentato',
  Abandoned: 'Abbandonato',
  Deleted: 'Eliminato',
}

function formatDate(d?: string): string {
  if (!d) return '-'
  return new Date(d).toLocaleString('it-IT', { day: '2-digit', month: '2-digit', year: '2-digit', hour: '2-digit', minute: '2-digit' })
}

function QueueItemRow({ item }: { item: QueueItem }) {
  const color = statusColors[item.status] || 'gray'
  return (
    <Flex
      align="start"
      py={3}
      px={2}
      borderBottom="1px solid #f0f0f2"
      _hover={{ bg: '#fafafa' }}
      transition="background 0.1s"
      gap={3}
      fontSize="13px"
    >
      <Box flex={1.5} fontWeight="500" color="#1d1d1f">
        {item.queueName || '-'}
      </Box>
      <Box flex={1}>
        <Badge colorPalette={color} fontSize="10px">
          {statusLabels[item.status] || item.status}
        </Badge>
      </Box>
      <Box flex={0.8} color="#86868b">{item.priority || '-'}</Box>
      <Box flex={1} color="#86868b">{formatDate(item.startProcessing)}</Box>
      <Box flex={1} color="#86868b">{formatDate(item.endProcessing)}</Box>
      <Box flex={1} color="#86868b">
        {item.processingExceptionType ? (
          <Badge colorPalette="red" fontSize="9px" variant="subtle">{item.processingExceptionType}</Badge>
        ) : '-'}
      </Box>
      <Box flex={2} color="#86868b" fontSize="12px" wordBreak="break-word">
        {item.errorMessage || '-'}
      </Box>
    </Flex>
  )
}

import type { ReactNode } from 'react'

export default function QueueItemsTab({ processNames, processFilter }: { processNames?: string; processFilter?: ReactNode }) {
  const { data: connectors } = useConnectorsQuery()
  const [filters, setFilters] = useState<QueueItemFilters>({ page: 1, limit: 20 })
  const { data, isLoading } = useQueueItemsQuery(
    { ...filters, processNames },
    { pollingInterval: 60000 },
  )

  const updateFilter = (key: keyof QueueItemFilters, value: string | number) => {
    setFilters((prev) => ({ ...prev, [key]: value || undefined, page: 1 }))
  }

  const total = data?.total ?? 0
  const successCount = data?.data.filter((i) => i.status === 'Successful').length ?? 0
  const failedCount = data?.data.filter((i) => i.status === 'Failed').length ?? 0
  const inProgressCount = data?.data.filter((i) => i.status === 'New' || i.status === 'InProgress').length ?? 0
  const retriedCount = data?.data.filter((i) => i.status === 'Retried').length ?? 0

  return (
    <Flex direction="column" gap={4}>
      <Flex gap={3} wrap="wrap">
        <StatCard label="Totale" value={total} color="#1d1d1f" />
        <StatCard label="Completati" value={successCount} color="#34c759" />
        <StatCard label="Falliti" value={failedCount} color="#ff3b30" />
        <StatCard label="In lavorazione" value={inProgressCount} color="#007aff" />
        <StatCard label="Ritentati" value={retriedCount} color="#ff9500" />
      </Flex>

      <Flex gap={3} wrap="wrap" align="end">
        <Box>
          <Text fontSize="11px" color="#86868b" mb={1}>Processo</Text>
          {processFilter}
        </Box>
        <Box>
          <Text fontSize="11px" color="#86868b" mb={1}>Stato</Text>
          <NativeSelect.Root size="sm" maxW="160px">
            <NativeSelect.Field value={filters.status || ''} onChange={(e) => updateFilter('status', e.target.value)}>
              <option value="">Tutti</option>
              <option value="Successful">Completato</option>
              <option value="Failed">Fallito</option>
              <option value="New">Nuovo</option>
              <option value="Retried">Ritentato</option>
              <option value="Abandoned">Abbandonato</option>
            </NativeSelect.Field>
          </NativeSelect.Root>
        </Box>
        {connectors && connectors.length > 1 && (
          <Box>
            <Text fontSize="11px" color="#86868b" mb={1}>Connettore</Text>
            <NativeSelect.Root size="sm" maxW="180px">
              <NativeSelect.Field value={filters.connectorId || ''} onChange={(e) => updateFilter('connectorId', parseInt(e.target.value) || 0)}>
                <option value="">Tutti</option>
                {connectors.map((c) => <option key={c.id} value={c.id}>{c.name}</option>)}
              </NativeSelect.Field>
            </NativeSelect.Root>
          </Box>
        )}
        <Box>
          <Text fontSize="11px" color="#86868b" mb={1}>Coda</Text>
          <Input size="sm" maxW="180px" placeholder="Nome coda" value={filters.queueName || ''} onChange={(e) => updateFilter('queueName', e.target.value)} />
        </Box>
      </Flex>

      {isLoading ? (
        <Flex justify="center" py={10}><Spinner size="lg" color="brand.300" /></Flex>
      ) : !data || data.data.length === 0 ? (
        <Text fontSize="13px" color="#86868b" textAlign="center" py={10}>
          Nessun elemento trovato. Sincronizza un connettore dalla sezione Impostazioni.
        </Text>
      ) : (
        <>
          <Flex px={2} py={2} borderBottom="2px solid #f0f0f2" gap={3} fontSize="11px" fontWeight="700" color="#86868b" textTransform="uppercase" letterSpacing="0.5px">
            <Box flex={1.5}>Coda</Box>
            <Box flex={1}>Stato</Box>
            <Box flex={0.8}>Priorita'</Box>
            <Box flex={1}>Inizio</Box>
            <Box flex={1}>Fine</Box>
            <Box flex={1}>Eccezione</Box>
            <Box flex={2}>Errore</Box>
          </Flex>

          {data.data.map((item) => <QueueItemRow key={item.id} item={item} />)}

          {data.totalPages > 1 && (
            <Flex justify="space-between" align="center" pt={3}>
              <Text fontSize="12px" color="#86868b">
                Pagina {data.page} di {data.totalPages} ({data.total} risultati)
              </Text>
              <Flex gap={1}>
                <Button size="xs" variant="outline" disabled={data.page <= 1} onClick={() => setFilters((p) => ({ ...p, page: (p.page || 1) - 1 }))}>
                  <LuChevronLeft size={14} />
                </Button>
                <Button size="xs" variant="outline" disabled={data.page >= data.totalPages} onClick={() => setFilters((p) => ({ ...p, page: (p.page || 1) + 1 }))}>
                  <LuChevronRight size={14} />
                </Button>
              </Flex>
            </Flex>
          )}
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

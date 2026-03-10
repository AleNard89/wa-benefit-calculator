import { Badge, Box, Button, Flex, Spinner, Text } from '@chakra-ui/react'
import { useState } from 'react'
import { LuChevronLeft, LuChevronRight } from 'react-icons/lu'

import { NativeSelect } from '@chakra-ui/react'
import { useJobsQuery, useConnectorsQuery } from './api'
import type { JobExecution, JobFilters } from './types'

const stateColors: Record<string, string> = {
  Successful: 'green',
  Faulted: 'red',
  Stopped: 'orange',
  Running: 'blue',
  Pending: 'gray',
}

const stateLabels: Record<string, string> = {
  Successful: 'Completato',
  Faulted: 'Errore',
  Stopped: 'Fermato',
  Running: 'In esecuzione',
  Pending: 'In attesa',
}

function formatDuration(start?: string, end?: string): string {
  if (!start || !end) return '-'
  const ms = new Date(end).getTime() - new Date(start).getTime()
  if (ms < 0) return '-'
  const secs = Math.floor(ms / 1000)
  if (secs < 60) return `${secs}s`
  const mins = Math.floor(secs / 60)
  const remSecs = secs % 60
  if (mins < 60) return `${mins}m ${remSecs}s`
  const hours = Math.floor(mins / 60)
  const remMins = mins % 60
  return `${hours}h ${remMins}m`
}

function formatDate(d?: string): string {
  if (!d) return '-'
  return new Date(d).toLocaleString('it-IT', { day: '2-digit', month: '2-digit', year: '2-digit', hour: '2-digit', minute: '2-digit' })
}

function JobRow({ job }: { job: JobExecution }) {
  const color = stateColors[job.state] || 'gray'
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
      <Box flex={2} fontWeight="500" color="#1d1d1f">
        {job.processName || '-'}
      </Box>
      <Box flex={1}>
        <Badge colorPalette={color} fontSize="10px">
          {stateLabels[job.state] || job.state}
        </Badge>
      </Box>
      <Box flex={1} color="#86868b">{job.sourceType || '-'}</Box>
      <Box flex={1} color="#86868b">{job.hostMachine || '-'}</Box>
      <Box flex={1} color="#86868b">{formatDate(job.startTime)}</Box>
      <Box flex={1} color="#86868b">{formatDate(job.endTime)}</Box>
      <Box flex={0.8} color="#86868b">{formatDuration(job.startTime, job.endTime)}</Box>
      <Box flex={2} color="#86868b" fontSize="12px" wordBreak="break-word">
        {job.info || '-'}
      </Box>
    </Flex>
  )
}

import type { ReactNode } from 'react'

export default function JobsTab({ processNames, processFilter }: { processNames?: string; processFilter?: ReactNode }) {
  const { data: connectors } = useConnectorsQuery()
  const [filters, setFilters] = useState<JobFilters>({ page: 1, limit: 20 })
  const { data, isLoading } = useJobsQuery(
    { ...filters, processNames },
    { pollingInterval: 60000 },
  )

  const updateFilter = (key: keyof JobFilters, value: string | number) => {
    setFilters((prev) => ({ ...prev, [key]: value || undefined, page: 1 }))
  }

  return (
    <Flex direction="column" gap={4}>
      <Flex gap={3} wrap="wrap" align="end">
        <Box>
          <Text fontSize="11px" color="#86868b" mb={1}>Processo</Text>
          {processFilter}
        </Box>
        <Box>
          <Text fontSize="11px" color="#86868b" mb={1}>Stato</Text>
          <NativeSelect.Root size="sm" maxW="160px">
            <NativeSelect.Field value={filters.state || ''} onChange={(e) => updateFilter('state', e.target.value)}>
              <option value="">Tutti</option>
              <option value="Successful">Completato</option>
              <option value="Faulted">Errore</option>
              <option value="Stopped">Fermato</option>
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
      </Flex>

      {isLoading ? (
        <Flex justify="center" py={10}><Spinner size="lg" color="brand.300" /></Flex>
      ) : !data || data.data.length === 0 ? (
        <Text fontSize="13px" color="#86868b" textAlign="center" py={10}>
          Nessuna esecuzione trovata. Sincronizza un connettore dalla sezione Impostazioni.
        </Text>
      ) : (
        <>
          <Flex px={2} py={2} borderBottom="2px solid #f0f0f2" gap={3} fontSize="11px" fontWeight="700" color="#86868b" textTransform="uppercase" letterSpacing="0.5px">
            <Box flex={2}>Processo</Box>
            <Box flex={1}>Stato</Box>
            <Box flex={1}>Tipo</Box>
            <Box flex={1}>Macchina</Box>
            <Box flex={1}>Inizio</Box>
            <Box flex={1}>Fine</Box>
            <Box flex={0.8}>Durata</Box>
            <Box flex={2}>Info</Box>
          </Flex>

          {data.data.map((job) => <JobRow key={job.id} job={job} />)}

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

import { Badge, Box, Button, Flex, Spinner, Text } from '@chakra-ui/react'
import { useState } from 'react'
import { LuChevronLeft, LuChevronRight } from 'react-icons/lu'

import { NativeSelect } from '@chakra-ui/react'
import { useSchedulesQuery, useConnectorsQuery } from './api'
import type { ProcessSchedule, ScheduleFilters } from './types'

function formatDate(d?: string): string {
  if (!d) return '-'
  return new Date(d).toLocaleString('it-IT', { day: '2-digit', month: '2-digit', year: '2-digit', hour: '2-digit', minute: '2-digit' })
}

function ScheduleRow({ schedule }: { schedule: ProcessSchedule }) {
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
        {schedule.releaseName || schedule.packageName || '-'}
      </Box>
      <Box flex={0.6}>
        <Badge colorPalette={schedule.enabled ? 'green' : 'gray'} fontSize="10px">
          {schedule.enabled ? 'Attivo' : 'Disattivato'}
        </Badge>
      </Box>
      <Box flex={2} color="#86868b" fontSize="12px">
        {schedule.cronSummary || schedule.cronExpression || '-'}
      </Box>
      <Box flex={1} color="#86868b">
        {formatDate(schedule.nextOccurrence)}
      </Box>
      <Box flex={1} color="#86868b">
        {schedule.timezoneIana || schedule.timezoneId || '-'}
      </Box>
      <Box flex={0.8} color="#86868b">
        {schedule.folderName || '-'}
      </Box>
    </Flex>
  )
}

import type { ReactNode } from 'react'

export default function SchedulesTab({ processNames, processFilter }: { processNames?: string; processFilter?: ReactNode }) {
  const { data: connectors } = useConnectorsQuery()
  const [filters, setFilters] = useState<ScheduleFilters>({ page: 1, limit: 20 })
  const { data, isLoading } = useSchedulesQuery(
    { ...filters, processNames },
    { pollingInterval: 60000 },
  )

  const updateFilter = (key: keyof ScheduleFilters, value: string | number | boolean | undefined) => {
    setFilters((prev) => ({ ...prev, [key]: value, page: 1 }))
  }

  const enabledCount = data?.data.filter((s) => s.enabled).length ?? 0
  const disabledCount = data?.data.filter((s) => !s.enabled).length ?? 0

  return (
    <Flex direction="column" gap={4}>
      <Flex gap={3} wrap="wrap">
        <StatCard label="Totale" value={data?.total ?? 0} color="#1d1d1f" />
        <StatCard label="Attivi" value={enabledCount} color="#34c759" />
        <StatCard label="Disattivati" value={disabledCount} color="#86868b" />
      </Flex>

      <Flex gap={3} wrap="wrap" align="end">
        <Box>
          <Text fontSize="11px" color="#86868b" mb={1}>Processo</Text>
          {processFilter}
        </Box>
        <Box>
          <Text fontSize="11px" color="#86868b" mb={1}>Stato</Text>
          <NativeSelect.Root size="sm" maxW="160px">
            <NativeSelect.Field
              value={filters.enabled === undefined ? '' : filters.enabled.toString()}
              onChange={(e) => {
                const val = e.target.value
                updateFilter('enabled', val === '' ? undefined : val === 'true')
              }}
            >
              <option value="">Tutti</option>
              <option value="true">Attivo</option>
              <option value="false">Disattivato</option>
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
          Nessuna schedulazione trovata. Sincronizza un connettore dalla sezione Impostazioni.
        </Text>
      ) : (
        <>
          <Flex px={2} py={2} borderBottom="2px solid #f0f0f2" gap={3} fontSize="11px" fontWeight="700" color="#86868b" textTransform="uppercase" letterSpacing="0.5px">
            <Box flex={2}>Processo</Box>
            <Box flex={0.6}>Stato</Box>
            <Box flex={2}>Schedulazione</Box>
            <Box flex={1}>Prossima Esecuzione</Box>
            <Box flex={1}>Fuso Orario</Box>
            <Box flex={0.8}>Cartella</Box>
          </Flex>

          {data.data.map((schedule) => <ScheduleRow key={schedule.id} schedule={schedule} />)}

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

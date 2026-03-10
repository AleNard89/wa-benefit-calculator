import { Badge, Box, Button, Flex, Heading, Input, Spinner, Table, Text } from '@chakra-ui/react'
import { useState } from 'react'
import { useNavigate } from 'react-router-dom'

import { useCurrentUser } from '@/Auth/Hooks'
import { toaster } from '@/Snippets/toaster'
import KpiCards from '../Components/KpiCards'
import { useProcessesQuery, useProcessStatsQuery, useRestoreProcessMutation } from '../Services/Api'
import { formatCurrency, formatPercent, statusColorMap, statusLabelMap } from '../Utils'
import type { ProcessListParams } from '../Types'

import { NativeSelect } from '@chakra-ui/react'

export default function ProcessListView() {
  const navigate = useNavigate()
  const user = useCurrentUser()
  const isAdmin = user?.isSuperuser || user?.currentPermissions?.includes('processes:process.delete')
  const [showDeleted, setShowDeleted] = useState(false)
  const [params, setParams] = useState<ProcessListParams>({ page: 1, limit: 25 })
  const [search, setSearch] = useState('')
  const [restoreProcess] = useRestoreProcessMutation()

  const { data: listData, isLoading } = useProcessesQuery({ ...params, deleted: showDeleted || undefined })
  const { data: stats } = useProcessStatsQuery()

  const handleSearch = () => {
    setParams((prev) => ({ ...prev, search, page: 1 }))
  }

  const handleStatusFilter = (status: string) => {
    setParams((prev) => ({ ...prev, status: status || undefined, page: 1 }))
  }

  return (
    <Flex direction="column" gap={6}>
      <Flex justify="space-between" align="center">
        <Heading size="lg">Lista Processi</Heading>
        <Button colorPalette="brand" onClick={() => navigate('/processes/create')}>
          + Nuova Proposta
        </Button>
      </Flex>

      {stats && <KpiCards stats={stats} />}

      <Flex gap={3} align="center">
        <Input
          placeholder="Cerca per nome, proponente, area..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          onKeyDown={(e) => e.key === 'Enter' && handleSearch()}
          maxW="400px"
        />
        <Button variant="outline" onClick={handleSearch}>Cerca</Button>
        {!showDeleted && (
          <NativeSelect.Root maxW="200px">
            <NativeSelect.Field
              value={params.status || ''}
              onChange={(e) => handleStatusFilter(e.target.value)}
            >
              <option value="">Tutti gli stati</option>
              <option value="To Valuate">Da Valutare</option>
              <option value="Analysis">In Analisi</option>
              <option value="Ongoing">In Corso</option>
              <option value="Production">Produzione</option>
            </NativeSelect.Field>
          </NativeSelect.Root>
        )}
        {isAdmin && (
          <Button
            size="sm"
            variant={showDeleted ? 'solid' : 'outline'}
            colorPalette={showDeleted ? 'red' : 'gray'}
            onClick={() => { setShowDeleted(!showDeleted); setParams(p => ({ ...p, page: 1, status: undefined })) }}
          >
            {showDeleted ? 'Mostra attivi' : 'Cestino'}
          </Button>
        )}
      </Flex>

      {isLoading ? (
        <Flex justify="center" py={10}><Spinner size="xl" color="brand.300" /></Flex>
      ) : !listData?.data.length ? (
        <Text color="fg.muted" py={10} textAlign="center">Nessun processo trovato</Text>
      ) : (
        <>
          <Box overflowX="auto">
            <Table.Root size="sm">
              <Table.Header>
                <Table.Row>
                  <Table.ColumnHeader>Nome Processo</Table.ColumnHeader>
                  <Table.ColumnHeader>Area</Table.ColumnHeader>
                  <Table.ColumnHeader>Proponente</Table.ColumnHeader>
                  <Table.ColumnHeader>Bot Collegati</Table.ColumnHeader>
                  <Table.ColumnHeader>ROI</Table.ColumnHeader>
                  <Table.ColumnHeader>Risparmio Annuale</Table.ColumnHeader>
                  <Table.ColumnHeader>Stato</Table.ColumnHeader>
                  {showDeleted && <Table.ColumnHeader>Azione</Table.ColumnHeader>}
                </Table.Row>
              </Table.Header>
              <Table.Body>
                {listData.data.map((p) => (
                  <Table.Row
                    key={p.id}
                    cursor={showDeleted ? 'default' : 'pointer'}
                    _hover={{ bg: 'bg.subtle' }}
                    onClick={() => !showDeleted && navigate(`/processes/${p.id}`)}
                    opacity={showDeleted ? 0.7 : 1}
                  >
                    <Table.Cell fontWeight="600">{p.processName}</Table.Cell>
                    <Table.Cell>{p.data?.area}</Table.Cell>
                    <Table.Cell>{p.data?.proposer}</Table.Cell>
                    <Table.Cell>
                      {p.data?.linkedBots && p.data.linkedBots.length > 0 ? (
                        <Flex gap={1} wrap="wrap">
                          {p.data.linkedBots.map((bot: string) => (
                            <Badge key={bot} colorPalette="blue" fontSize="10px">{bot}</Badge>
                          ))}
                        </Flex>
                      ) : (
                        <Text fontSize="xs" color="fg.muted">-</Text>
                      )}
                    </Table.Cell>
                    <Table.Cell>{p.results?.roi != null ? formatPercent(p.results.roi) : '-'}</Table.Cell>
                    <Table.Cell>{p.results?.annualSavings != null ? formatCurrency(p.results.annualSavings) : '-'}</Table.Cell>
                    <Table.Cell>
                      <Badge colorPalette={statusColorMap[p.status] || 'gray'}>
                        {statusLabelMap[p.status] || p.status}
                      </Badge>
                    </Table.Cell>
                    {showDeleted && (
                      <Table.Cell>
                        <Button
                          size="xs"
                          colorPalette="green"
                          variant="outline"
                          onClick={async (e) => {
                            e.stopPropagation()
                            try {
                              await restoreProcess(p.id).unwrap()
                              toaster.success({ title: 'Processo ripristinato' })
                            } catch {
                              toaster.error({ title: 'Errore ripristino' })
                            }
                          }}
                        >
                          Ripristina
                        </Button>
                      </Table.Cell>
                    )}
                  </Table.Row>
                ))}
              </Table.Body>
            </Table.Root>
          </Box>

          {listData.totalPages > 1 && (
            <Flex justify="center" gap={2}>
              <Button
                size="sm"
                variant="outline"
                disabled={params.page === 1}
                onClick={() => setParams((prev) => ({ ...prev, page: (prev.page || 1) - 1 }))}
              >
                Precedente
              </Button>
              <Text fontSize="sm" color="fg.muted" lineHeight="32px">
                Pagina {listData.page} di {listData.totalPages} ({listData.total} risultati)
              </Text>
              <Button
                size="sm"
                variant="outline"
                disabled={params.page === listData.totalPages}
                onClick={() => setParams((prev) => ({ ...prev, page: (prev.page || 1) + 1 }))}
              >
                Successiva
              </Button>
            </Flex>
          )}
        </>
      )}
    </Flex>
  )
}

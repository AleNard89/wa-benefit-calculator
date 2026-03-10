import { Badge, Box, Flex, Text } from '@chakra-ui/react'
import { useNavigate } from 'react-router-dom'
import { Cell, Pie, PieChart, ResponsiveContainer, Tooltip } from 'recharts'

import { useCurrentUser } from '@/Auth/Hooks'
import { useProcessesQuery, useProcessStatsQuery } from '@/Processes/Services/Api'
import { useOrchestratorDashboardStatsQuery } from '@/Orchestrator/api'
import { formatCurrency, statusLabelMap } from '@/Processes/Utils'
import type { ProcessStats } from '@/Processes/Types'
import type { JobExecution } from '@/Orchestrator/types'

const STATUS_COLORS: Record<string, string> = {
  'Da Valutare': '#86868b',
  'In Analisi': '#007aff',
  'In Corso': '#ff9500',
  'Produzione': '#34c759',
}

const JOB_STATE_COLORS: Record<string, string> = {
  Successful: '#34c759',
  Faulted: '#ff3b30',
  Stopped: '#ff9500',
  Running: '#007aff',
  Pending: '#86868b',
}

const JOB_STATE_LABELS: Record<string, string> = {
  Successful: 'Completato',
  Faulted: 'Errore',
  Stopped: 'Fermato',
  Running: 'In esecuzione',
  Pending: 'In attesa',
}

const CARD_STYLE = {
  bg: 'white',
  borderRadius: '12px',
  boxShadow: '0 1px 4px rgba(0,0,0,0.06)',
  p: 4,
} as const

const ROW_HEIGHT = 280

function StatCard({ label, value, color }: { label: string; value: number; color: string }) {
  return (
    <Box {...CARD_STYLE} flex={1} minW="120px" py={2} px={4}>
      <Text fontSize="10px" color="#86868b" mb={0}>{label}</Text>
      <Text fontSize="20px" fontWeight="700" color={color}>{value}</Text>
    </Box>
  )
}

function StatusChart({ stats }: { stats: ProcessStats }) {
  const data = [
    { name: 'Da Valutare', value: stats.toValuate },
    { name: 'In Analisi', value: stats.analysis },
    { name: 'In Corso', value: stats.ongoing },
    { name: 'Produzione', value: stats.production },
  ].filter((d) => d.value > 0)

  if (data.length === 0) {
    return (
      <Flex align="center" justify="center" h="100%">
        <Text fontSize="13px" color="#86868b">Nessun processo presente</Text>
      </Flex>
    )
  }

  return (
    <ResponsiveContainer width="100%" height="100%">
      <PieChart>
        <Pie
          data={data}
          dataKey="value"
          nameKey="name"
          cx="50%"
          cy="50%"
          innerRadius={45}
          outerRadius={75}
          paddingAngle={3}
          label={({ name, value }) => `${name}: ${value}`}
        >
          {data.map((entry) => (
            <Cell key={entry.name} fill={STATUS_COLORS[entry.name] || '#86868b'} />
          ))}
        </Pie>
        <Tooltip />
      </PieChart>
    </ResponsiveContainer>
  )
}

const BOT_CHART_COLORS: Record<string, string> = {
  Completati: '#34c759',
  Errori: '#ff3b30',
  'In esecuzione': '#007aff',
}

function BotExecutionsPieChart({ successful, faulted, running }: { successful: number; faulted: number; running: number }) {
  const data = [
    { name: 'Completati', value: successful },
    { name: 'Errori', value: faulted },
    { name: 'In esecuzione', value: running },
  ].filter((d) => d.value > 0)

  if (data.length === 0) {
    return (
      <Flex align="center" justify="center" h="100%">
        <Text fontSize="13px" color="#86868b">Nessun dato disponibile dall'Orchestrator</Text>
      </Flex>
    )
  }

  return (
    <ResponsiveContainer width="100%" height="100%">
      <PieChart>
        <Pie
          data={data}
          dataKey="value"
          nameKey="name"
          cx="50%"
          cy="50%"
          innerRadius={45}
          outerRadius={75}
          paddingAngle={3}
          label={({ name, value }) => `${name}: ${value}`}
        >
          {data.map((entry) => (
            <Cell key={entry.name} fill={BOT_CHART_COLORS[entry.name] || '#86868b'} />
          ))}
        </Pie>
        <Tooltip />
      </PieChart>
    </ResponsiveContainer>
  )
}

function formatJobDate(d?: string): string {
  if (!d) return '-'
  return new Date(d).toLocaleString('it-IT', { day: '2-digit', month: '2-digit', year: '2-digit', hour: '2-digit', minute: '2-digit' })
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

function RecentJobsTable({ jobs }: { jobs: JobExecution[] }) {
  if (jobs.length === 0) {
    return (
      <Flex align="center" justify="center" h="100%">
        <Text fontSize="13px" color="#86868b">Nessuna esecuzione recente</Text>
      </Flex>
    )
  }

  return (
    <Flex direction="column" gap={0} h="100%">
      <Flex px={2} py={1.5} borderBottom="2px solid #f0f0f2" gap={3} fontSize="10px" fontWeight="700" color="#86868b" textTransform="uppercase" letterSpacing="0.5px" flexShrink={0}>
        <Box flex={2}>Bot</Box>
        <Box flex={1}>Stato</Box>
        <Box flex={1}>Inizio</Box>
        <Box flex={0.8}>Durata</Box>
      </Flex>
      <Box overflowY="auto" flex={1}>
        {jobs.map((job, i) => (
          <Flex
            key={job.id}
            align="center"
            py={2}
            px={2}
            borderTop={i > 0 ? '1px solid #f0f0f2' : 'none'}
            gap={3}
            fontSize="12px"
          >
            <Box flex={2} fontWeight="500" color="#1d1d1f" overflow="hidden" textOverflow="ellipsis" whiteSpace="nowrap">
              {job.processName || '-'}
            </Box>
            <Box flex={1}>
              <Badge
                colorPalette={JOB_STATE_COLORS[job.state] === '#34c759' ? 'green' : JOB_STATE_COLORS[job.state] === '#ff3b30' ? 'red' : JOB_STATE_COLORS[job.state] === '#ff9500' ? 'orange' : JOB_STATE_COLORS[job.state] === '#007aff' ? 'blue' : 'gray'}
                fontSize="10px"
              >
                {JOB_STATE_LABELS[job.state] || job.state}
              </Badge>
            </Box>
            <Box flex={1} color="#86868b" fontSize="11px">{formatJobDate(job.startTime)}</Box>
            <Box flex={0.8} color="#86868b" fontSize="11px">{formatDuration(job.startTime, job.endTime)}</Box>
          </Flex>
        ))}
      </Box>
    </Flex>
  )
}

export default function HomeView() {
  const user = useCurrentUser()
  const navigate = useNavigate()
  const { data: stats } = useProcessStatsQuery()
  const { data: recentProcesses } = useProcessesQuery({ limit: 5, sortBy: 'created_at', order: 'desc' })
  const { data: orchStats } = useOrchestratorDashboardStatsQuery()

  return (
    <Box>
      <Text fontSize="24px" fontWeight="700" color="#1d1d1f">Dashboard</Text>
      <Text fontSize="14px" color="#86868b" mt={1} mb={4}>
        Benvenuto{user ? `, ${user.firstName} ${user.lastName}` : ''}
      </Text>

      {/* Process KPI Cards */}
      {stats && (
        <Flex gap={2} mb={4} wrap="wrap">
          <StatCard label="Totale Processi" value={stats.total} color="#1d1d1f" />
          <StatCard label="Da Valutare" value={stats.toValuate} color="#86868b" />
          <StatCard label="In Analisi" value={stats.analysis} color="#007aff" />
          <StatCard label="In Corso" value={stats.ongoing} color="#ff9500" />
          <StatCard label="Produzione" value={stats.production} color="#34c759" />
        </Flex>
      )}

      <Flex gap={3} mb={4}>
        <Box {...CARD_STYLE} flex={1} display="flex" flexDirection="column" h={`${ROW_HEIGHT}px`}>
          <Text fontSize="14px" fontWeight="600" color="#1d1d1f" mb={2}>
            Distribuzione per Stato
          </Text>
          <Box flex={1}>
            {stats && <StatusChart stats={stats} />}
          </Box>
        </Box>

        <Box {...CARD_STYLE} flex={1.5} display="flex" flexDirection="column" h={`${ROW_HEIGHT}px`}>
          <Text fontSize="14px" fontWeight="600" color="#1d1d1f" mb={2}>
            Ultimi Processi
          </Text>
          <Box flex={1} overflowY="auto">
            {recentProcesses?.data && recentProcesses.data.length > 0 ? (
              <Flex direction="column" gap={0}>
                {recentProcesses.data.map((p, i) => (
                  <Flex
                    key={p.id}
                    align="center"
                    justify="space-between"
                    py={2}
                    px={2}
                    borderTop={i > 0 ? '1px solid #f0f0f2' : 'none'}
                    borderRadius="6px"
                    cursor="pointer"
                    _hover={{ bg: '#f5f5f7' }}
                    transition="background 0.15s"
                    onClick={() => navigate(`/processes/${p.id}`)}
                  >
                    <Box>
                      <Text fontSize="13px" fontWeight="500" color="#1d1d1f">{p.processName}</Text>
                      <Text fontSize="11px" color="#86868b">{p.data?.area} &middot; {p.data?.proposer}</Text>
                    </Box>
                    <Flex align="center" gap={2}>
                      {p.results?.annualSavings != null && (
                        <Text fontSize="12px" fontWeight="600" color="#007aff">
                          {formatCurrency(p.results.annualSavings)}
                        </Text>
                      )}
                      <Box
                        px={2}
                        py={0.5}
                        borderRadius="5px"
                        bg={
                          p.status === 'Production' ? '#34c75915' :
                          p.status === 'Ongoing' ? '#ff950015' :
                          p.status === 'Analysis' ? '#007aff15' : '#86868b15'
                        }
                      >
                        <Text
                          fontSize="10px"
                          fontWeight="600"
                          color={
                            p.status === 'Production' ? '#34c759' :
                            p.status === 'Ongoing' ? '#ff9500' :
                            p.status === 'Analysis' ? '#007aff' : '#86868b'
                          }
                        >
                          {statusLabelMap[p.status] || p.status}
                        </Text>
                      </Box>
                    </Flex>
                  </Flex>
                ))}
              </Flex>
            ) : (
              <Flex align="center" justify="center" h="100%">
                <Text fontSize="13px" color="#86868b">
                  Nessun processo. Inizia creandone uno nuovo.
                </Text>
              </Flex>
            )}
          </Box>
        </Box>
      </Flex>

      {/* Orchestrator Section */}
      {orchStats && orchStats.totalJobs > 0 && (
        <>
          <Text fontSize="18px" fontWeight="700" color="#1d1d1f" mb={3}>Orchestrator</Text>

          <Flex gap={2} mb={4} wrap="wrap">
            <StatCard label="Esecuzioni Totali" value={orchStats.totalJobs} color="#1d1d1f" />
            <StatCard label="Completate" value={orchStats.successful} color="#34c759" />
            <StatCard label="Errori" value={orchStats.faulted} color="#ff3b30" />
            <StatCard label="In Esecuzione" value={orchStats.running} color="#007aff" />
          </Flex>

          <Flex gap={3}>
            <Box {...CARD_STYLE} flex={1} display="flex" flexDirection="column" h={`${ROW_HEIGHT}px`}>
              <Text fontSize="14px" fontWeight="600" color="#1d1d1f" mb={2}>
                Esecuzioni Bot
              </Text>
              <Box flex={1}>
                <BotExecutionsPieChart
                  successful={orchStats.successful}
                  faulted={orchStats.faulted}
                  running={orchStats.running}
                />
              </Box>
            </Box>

            <Box {...CARD_STYLE} flex={1.5} display="flex" flexDirection="column" h={`${ROW_HEIGHT}px`}>
              <Text fontSize="14px" fontWeight="600" color="#1d1d1f" mb={2}>
                Ultime Esecuzioni
              </Text>
              <Box flex={1} overflow="hidden">
                <RecentJobsTable jobs={orchStats.recentJobs} />
              </Box>
            </Box>
          </Flex>
        </>
      )}
    </Box>
  )
}

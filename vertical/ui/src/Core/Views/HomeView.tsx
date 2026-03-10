import { Box, Flex, Text } from '@chakra-ui/react'
import { useNavigate } from 'react-router-dom'
import { Cell, Pie, PieChart, ResponsiveContainer, Tooltip } from 'recharts'

import { useCurrentUser } from '@/Auth/Hooks'
import { useProcessesQuery, useProcessStatsQuery } from '@/Processes/Services/Api'
import { formatCurrency, statusLabelMap } from '@/Processes/Utils'
import type { ProcessStats } from '@/Processes/Types'

const STATUS_COLORS: Record<string, string> = {
  'Da Valutare': '#86868b',
  'In Analisi': '#007aff',
  'In Corso': '#ff9500',
  'Produzione': '#34c759',
}

function StatCard({ label, value, color }: { label: string; value: number; color: string }) {
  return (
    <Box bg="white" borderRadius="12px" boxShadow="0 1px 4px rgba(0,0,0,0.06)" p={5} flex={1} minW="150px">
      <Text fontSize="12px" color="#86868b" mb={1}>{label}</Text>
      <Text fontSize="28px" fontWeight="700" color={color}>{value}</Text>
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
      <Flex align="center" justify="center" h="200px">
        <Text fontSize="13px" color="#86868b">Nessun processo presente</Text>
      </Flex>
    )
  }

  return (
    <ResponsiveContainer width="100%" height={220}>
      <PieChart>
        <Pie
          data={data}
          dataKey="value"
          nameKey="name"
          cx="50%"
          cy="50%"
          innerRadius={50}
          outerRadius={85}
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

export default function HomeView() {
  const user = useCurrentUser()
  const navigate = useNavigate()
  const { data: stats } = useProcessStatsQuery()
  const { data: recentProcesses } = useProcessesQuery({ limit: 5, sortBy: 'created_at', order: 'desc' })

  return (
    <Box>
      <Text fontSize="28px" fontWeight="700" color="#1d1d1f">Dashboard</Text>
      <Text fontSize="15px" color="#86868b" mt={1} mb={6}>
        Benvenuto{user ? `, ${user.firstName} ${user.lastName}` : ''}
      </Text>

      {/* KPI Cards */}
      {stats && (
        <Flex gap={3} mb={6} wrap="wrap">
          <StatCard label="Totale Processi" value={stats.total} color="#1d1d1f" />
          <StatCard label="Da Valutare" value={stats.toValuate} color="#86868b" />
          <StatCard label="In Analisi" value={stats.analysis} color="#007aff" />
          <StatCard label="In Corso" value={stats.ongoing} color="#ff9500" />
          <StatCard label="Produzione" value={stats.production} color="#34c759" />
        </Flex>
      )}

      <Flex gap={4} wrap="wrap">
        {/* Chart */}
        <Box
          bg="white"
          borderRadius="12px"
          boxShadow="0 1px 4px rgba(0,0,0,0.06)"
          p={5}
          flex={1}
          minW="300px"
        >
          <Text fontSize="15px" fontWeight="600" color="#1d1d1f" mb={3}>
            Distribuzione per Stato
          </Text>
          {stats && <StatusChart stats={stats} />}
        </Box>

        {/* Recent Processes */}
        <Box
          bg="white"
          borderRadius="12px"
          boxShadow="0 1px 4px rgba(0,0,0,0.06)"
          p={5}
          flex={1.5}
          minW="400px"
        >
          <Text fontSize="15px" fontWeight="600" color="#1d1d1f" mb={3}>
            Ultimi Processi
          </Text>
          {recentProcesses?.data && recentProcesses.data.length > 0 ? (
            <Flex direction="column" gap={0}>
              {recentProcesses.data.map((p, i) => (
                <Flex
                  key={p.id}
                  align="center"
                  justify="space-between"
                  py={3}
                  px={2}
                  borderTop={i > 0 ? '1px solid #f0f0f2' : 'none'}
                  borderRadius="8px"
                  cursor="pointer"
                  _hover={{ bg: '#f5f5f7' }}
                  transition="background 0.15s"
                  onClick={() => navigate(`/processes/${p.id}`)}
                >
                  <Box>
                    <Text fontSize="14px" fontWeight="500" color="#1d1d1f">{p.processName}</Text>
                    <Text fontSize="12px" color="#86868b">{p.data?.area} &middot; {p.data?.proposer}</Text>
                  </Box>
                  <Flex align="center" gap={3}>
                    {p.results?.annualSavings != null && (
                      <Text fontSize="13px" fontWeight="600" color="#007aff">
                        {formatCurrency(p.results.annualSavings)}
                      </Text>
                    )}
                    <Box
                      px={2.5}
                      py={1}
                      borderRadius="6px"
                      bg={
                        p.status === 'Production' ? '#34c75915' :
                        p.status === 'Ongoing' ? '#ff950015' :
                        p.status === 'Analysis' ? '#007aff15' : '#86868b15'
                      }
                    >
                      <Text
                        fontSize="11px"
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
            <Flex align="center" justify="center" h="180px">
              <Text fontSize="13px" color="#86868b">
                Nessun processo. Inizia creandone uno nuovo.
              </Text>
            </Flex>
          )}
        </Box>
      </Flex>
    </Box>
  )
}

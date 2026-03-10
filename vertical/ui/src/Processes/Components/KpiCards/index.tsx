import { Box, Flex, Text } from '@chakra-ui/react'
import type { ProcessStats } from '../../Types'

interface Props {
  stats: ProcessStats
}

function Card({ label, count, color }: { label: string; count: number; color: string }) {
  return (
    <Box p={4} borderRadius="12px" bg="white" boxShadow="0 1px 4px rgba(0,0,0,0.06)" flex={1} minW="130px">
      <Text fontSize="12px" color="#86868b" mb={1}>{label}</Text>
      <Text fontSize="24px" fontWeight="700" color={color}>{count}</Text>
    </Box>
  )
}

export default function KpiCards({ stats }: Props) {
  return (
    <Flex gap={3} wrap="wrap">
      <Card label="Totale" count={stats.total} color="#1d1d1f" />
      <Card label="Da Valutare" count={stats.toValuate} color="#86868b" />
      <Card label="In Analisi" count={stats.analysis} color="#007aff" />
      <Card label="In Corso" count={stats.ongoing} color="#ff9500" />
      <Card label="Produzione" count={stats.production} color="#34c759" />
    </Flex>
  )
}

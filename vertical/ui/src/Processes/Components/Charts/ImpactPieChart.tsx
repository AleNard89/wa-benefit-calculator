import { Box, Text } from '@chakra-ui/react'
import { Cell, Pie, PieChart, ResponsiveContainer, Tooltip } from 'recharts'

interface Props {
  scores: {
    dataQualityScore: number
    auditScore: number
    customerExperienceScore: number
    errorReductionScore: number
    standardizationScore: number
    scalabilityScore: number
  }
}

const COLORS = ['#FFE600', '#3F7DE8', '#3ABB87', '#F56565', '#ED8936', '#9F7AEA']

const labels: Record<string, string> = {
  dataQualityScore: 'Qualita Dato',
  auditScore: 'Audit',
  customerExperienceScore: 'Esp. Cliente',
  errorReductionScore: 'Rid. Errori',
  standardizationScore: 'Standard.',
  scalabilityScore: 'Scalabilita',
}

export default function ImpactPieChart({ scores }: Props) {
  const data = Object.entries(scores).map(([key, value]) => ({
    name: labels[key] || key,
    value,
  }))

  return (
    <Box>
      <Text fontWeight="700" fontSize="sm" mb={2}>Distribuzione Impatto</Text>
      <ResponsiveContainer width="100%" height={250}>
        <PieChart>
          <Pie data={data} dataKey="value" nameKey="name" cx="50%" cy="50%" outerRadius={90} label>
            {data.map((_entry, i) => (
              <Cell key={i} fill={COLORS[i % COLORS.length]} />
            ))}
          </Pie>
          <Tooltip />
        </PieChart>
      </ResponsiveContainer>
    </Box>
  )
}

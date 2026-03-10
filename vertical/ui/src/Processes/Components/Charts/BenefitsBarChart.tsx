import { Box, Text } from '@chakra-ui/react'
import { Bar, BarChart, CartesianGrid, ResponsiveContainer, Tooltip, XAxis, YAxis } from 'recharts'
import type { ProcessResults as ProcessResult } from '../../Types'

interface Props {
  result: ProcessResult
}

export default function BenefitsBarChart({ result }: Props) {
  const data = [
    { name: 'Operativo', value: result.operationalSavings },
    { name: 'Rid. Errori', value: result.errorReductionSavings },
    { name: 'Produttivita', value: result.productivityBenefit },
  ]

  return (
    <Box>
      <Text fontWeight="700" fontSize="sm" mb={2}>Composizione Benefici (EUR)</Text>
      <ResponsiveContainer width="100%" height={250}>
        <BarChart data={data}>
          <CartesianGrid strokeDasharray="3 3" stroke="#444" />
          <XAxis dataKey="name" tick={{ fill: '#999', fontSize: 12 }} />
          <YAxis tick={{ fill: '#999', fontSize: 12 }} />
          <Tooltip />
          <Bar dataKey="value" fill="#FFE600" radius={[4, 4, 0, 0]} />
        </BarChart>
      </ResponsiveContainer>
    </Box>
  )
}

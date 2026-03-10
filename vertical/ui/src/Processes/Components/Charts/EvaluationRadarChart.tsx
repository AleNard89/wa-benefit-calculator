import { Box, Text } from '@chakra-ui/react'
import { PolarAngleAxis, PolarGrid, PolarRadiusAxis, Radar, RadarChart, ResponsiveContainer, Tooltip } from 'recharts'

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

export default function EvaluationRadarChart({ scores }: Props) {
  const data = [
    { subject: 'Qualita Dato', value: scores.dataQualityScore },
    { subject: 'Audit', value: scores.auditScore },
    { subject: 'Esp. Cliente', value: scores.customerExperienceScore },
    { subject: 'Rid. Errori', value: scores.errorReductionScore },
    { subject: 'Standard.', value: scores.standardizationScore },
    { subject: 'Scalabilita', value: scores.scalabilityScore },
  ]

  return (
    <Box>
      <Text fontWeight="700" fontSize="sm" mb={2}>Valutazione Impatto</Text>
      <ResponsiveContainer width="100%" height={250}>
        <RadarChart data={data}>
          <PolarGrid stroke="#444" />
          <PolarAngleAxis dataKey="subject" tick={{ fill: '#999', fontSize: 11 }} />
          <PolarRadiusAxis angle={30} domain={[0, 5]} tick={{ fill: '#666', fontSize: 10 }} />
          <Radar dataKey="value" stroke="#FFE600" fill="#FFE600" fillOpacity={0.3} />
          <Tooltip />
        </RadarChart>
      </ResponsiveContainer>
    </Box>
  )
}

import { Box, Flex, Text } from '@chakra-ui/react'
import { useState } from 'react'
import { LuInfo } from 'react-icons/lu'
import { PolarAngleAxis, PolarGrid, PolarRadiusAxis, Radar, RadarChart, ResponsiveContainer, Tooltip } from 'recharts'
import type { ProcessResults as ProcessResult } from '../../Types'
import { formatCurrency, formatNumber, formatPercent } from '../../Utils'

interface Props {
  result: ProcessResult
  scores: {
    dataQualityScore: number
    auditScore: number
    customerExperienceScore: number
    errorReductionScore: number
    standardizationScore: number
    scalabilityScore: number
  }
}

const kpiExplanations: Record<string, string> = {
  'ROI': 'Return on Investment = (Beneficio Netto Annuale / Costo Implementazione) x 100. Il beneficio netto e il risparmio annuale meno il costo di manutenzione.',
  'Risparmio Annuale': 'Somma di: Risparmio Operativo + Riduzione Errori + Beneficio Produttivita.',
  'Break-even': 'Mesi necessari per recuperare l\'investimento iniziale (implementazione + formazione + manutenzione) dato il beneficio netto mensile.',
  'Ore Risparmiate / Anno': 'Tempo per attivita (h) x Attivita al giorno x Giorni lavorativi x Fattore riduzione tempo (%).',
  'Risparmio Operativo': 'Costo annuale del lavoro manuale x Fattore riduzione tempo. Rappresenta il risparmio sul costo del personale.',
  'Riduzione Errori': 'Attivita al giorno x Giorni lavorativi x (Tasso errore attuale - Tasso errore post) x Costo per errore.',
  'Beneficio Produttivita': '(Attivita post-automazione - Attivita attuali) x Giorni lavorativi x Valore per attivita. Misura il valore delle attivita aggiuntive rese possibili.',
  'Impact Score': 'Media dei 6 punteggi di valutazione impatto (Qualita Dati, Audit, Customer Experience, Riduzione Errori, Standardizzazione, Scalabilita). Scala 1-5.',
}

function KpiItem({ label, value, accent }: { label: string; value: string; accent?: string }) {
  const [showInfo, setShowInfo] = useState(false)
  const explanation = kpiExplanations[label]

  return (
    <Box py={2.5} px={3} borderRadius="10px" bg="white" boxShadow="0 1px 4px rgba(0,0,0,0.06)" flex={1} minW="140px" position="relative">
      <Flex justify="space-between" align="center" mb={0.5}>
        <Text fontSize="11px" color="#86868b">{label}</Text>
        {explanation && (
          <Box
            as="button"
            type="button"
            display="flex"
            alignItems="center"
            justifyContent="center"
            w="16px"
            h="16px"
            borderRadius="full"
            bg={showInfo ? '#007aff' : '#f0f0f2'}
            color={showInfo ? 'white' : '#86868b'}
            cursor="pointer"
            transition="all 0.15s"
            _hover={{ bg: '#007aff', color: 'white' }}
            onClick={() => setShowInfo(!showInfo)}
          >
            <LuInfo size={10} />
          </Box>
        )}
      </Flex>
      <Text fontSize="16px" fontWeight="700" color={accent || '#1d1d1f'}>{value}</Text>
      {showInfo && explanation && (
        <Box
          mt={2}
          p={2.5}
          bg="#f5f5f7"
          borderRadius="8px"
          fontSize="11px"
          lineHeight="1.5"
          color="#48484a"
        >
          {explanation}
        </Box>
      )}
    </Box>
  )
}

export default function ProcessResults({ result, scores }: Props) {
  const radarData = [
    { subject: 'Qualita Dato', value: scores.dataQualityScore },
    { subject: 'Audit', value: scores.auditScore },
    { subject: 'Esp. Cliente', value: scores.customerExperienceScore },
    { subject: 'Rid. Errori', value: scores.errorReductionScore },
    { subject: 'Standard.', value: scores.standardizationScore },
    { subject: 'Scalabilita', value: scores.scalabilityScore },
  ]

  return (
    <Flex direction="column" gap={3}>
      <Text fontWeight="700" fontSize="17px" color="#1d1d1f">Risultati Calcolo</Text>

      <Flex gap={4}>
        <Flex direction="column" gap={2} flex={1}>
          <Flex gap={2}>
            <KpiItem label="ROI" value={formatPercent(result.roi)} accent={result.roi >= 0 ? '#34c759' : '#ff3b30'} />
            <KpiItem label="Risparmio Annuale" value={formatCurrency(result.annualSavings)} accent="#007aff" />
            <KpiItem label="Break-even" value={result.breakEvenMonths != null ? `${result.breakEvenMonths} mesi` : 'N/A'} />
          </Flex>
          <Flex gap={2}>
            <KpiItem label="Ore Risparmiate / Anno" value={formatNumber(result.hoursSavedAnnually)} />
            <KpiItem label="Risparmio Operativo" value={formatCurrency(result.operationalSavings)} />
            <KpiItem label="Riduzione Errori" value={formatCurrency(result.errorReductionSavings)} />
          </Flex>
          <Flex gap={2}>
            <KpiItem label="Beneficio Produttivita" value={formatCurrency(result.productivityBenefit)} />
            <KpiItem label="Impact Score" value={`${formatNumber(result.impactScore)} / 5`} accent="#ff9500" />
          </Flex>
        </Flex>

        <Box w="320px" flexShrink={0} bg="white" borderRadius="12px" boxShadow="0 1px 4px rgba(0,0,0,0.06)" p={4}>
          <Text fontSize="13px" fontWeight="600" color="#1d1d1f" mb={2}>Valutazione Impatto</Text>
          <ResponsiveContainer width="100%" height={220}>
            <RadarChart data={radarData}>
              <PolarGrid stroke="#e0e0e0" />
              <PolarAngleAxis dataKey="subject" tick={{ fill: '#86868b', fontSize: 11 }} />
              <PolarRadiusAxis angle={30} domain={[0, 5]} tick={{ fill: '#aaa', fontSize: 10 }} />
              <Radar dataKey="value" stroke="#34c759" fill="#34c759" fillOpacity={0.25} />
              <Tooltip />
            </RadarChart>
          </ResponsiveContainer>
        </Box>
      </Flex>
    </Flex>
  )
}

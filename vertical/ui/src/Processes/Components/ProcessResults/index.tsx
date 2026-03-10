import { Box, Flex, Text } from '@chakra-ui/react'
import { useState } from 'react'
import { LuInfo } from 'react-icons/lu'
import type { ProcessResults as ProcessResult } from '../../Types'
import { formatCurrency, formatNumber, formatPercent } from '../../Utils'

interface Props {
  result: ProcessResult
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
    <Box p={4} borderRadius="12px" bg="white" boxShadow="0 1px 4px rgba(0,0,0,0.06)" flex={1} minW="170px" position="relative">
      <Flex justify="space-between" align="center" mb={1}>
        <Text fontSize="12px" color="#86868b">{label}</Text>
        {explanation && (
          <Box
            as="button"
            type="button"
            display="flex"
            alignItems="center"
            justifyContent="center"
            w="18px"
            h="18px"
            borderRadius="full"
            bg={showInfo ? '#007aff' : '#f0f0f2'}
            color={showInfo ? 'white' : '#86868b'}
            cursor="pointer"
            transition="all 0.15s"
            _hover={{ bg: '#007aff', color: 'white' }}
            onClick={() => setShowInfo(!showInfo)}
          >
            <LuInfo size={11} />
          </Box>
        )}
      </Flex>
      <Text fontSize="20px" fontWeight="700" color={accent || '#1d1d1f'}>{value}</Text>
      {showInfo && explanation && (
        <Box
          mt={2}
          p={3}
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

export default function ProcessResults({ result }: Props) {
  return (
    <Flex direction="column" gap={4}>
      <Text fontWeight="700" fontSize="17px" color="#1d1d1f">Risultati Calcolo</Text>

      <Flex gap={3} wrap="wrap">
        <KpiItem label="ROI" value={formatPercent(result.roi)} accent={result.roi >= 0 ? '#34c759' : '#ff3b30'} />
        <KpiItem label="Risparmio Annuale" value={formatCurrency(result.annualSavings)} accent="#007aff" />
        <KpiItem label="Break-even" value={result.breakEvenMonths != null ? `${result.breakEvenMonths} mesi` : 'N/A'} />
        <KpiItem label="Ore Risparmiate / Anno" value={formatNumber(result.hoursSavedAnnually)} />
      </Flex>

      <Flex gap={3} wrap="wrap">
        <KpiItem label="Risparmio Operativo" value={formatCurrency(result.operationalSavings)} />
        <KpiItem label="Riduzione Errori" value={formatCurrency(result.errorReductionSavings)} />
        <KpiItem label="Beneficio Produttivita" value={formatCurrency(result.productivityBenefit)} />
        <KpiItem label="Impact Score" value={`${formatNumber(result.impactScore)} / 5`} accent="#ff9500" />
      </Flex>
    </Flex>
  )
}

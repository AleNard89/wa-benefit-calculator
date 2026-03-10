import { Box, Flex, Input, Text } from '@chakra-ui/react'
import { useEffect, useRef } from 'react'
import type { UseFormReturn } from 'react-hook-form'
import type { ProcessFormValues } from '../../Forms/ProcessFormSchema'

interface Props {
  form: UseFormReturn<ProcessFormValues>
}

const inputStyle = {
  bg: '#f5f5f7',
  border: 'none',
  borderRadius: '10px',
  h: '42px',
  fontSize: '14px',
  _focus: { bg: 'white', boxShadow: '0 0 0 3px rgba(0,122,255,0.15)' },
}

function Field({ label, error, children }: { label: string; error?: string; children: React.ReactNode }) {
  return (
    <Box flex={1}>
      <Text fontSize="12px" fontWeight="600" color="#86868b" mb={1.5}>{label}</Text>
      {children}
      {error && <Text fontSize="11px" color="#ff3b30" mt={1}>{error}</Text>}
    </Box>
  )
}

const PERIODICITY_DEFAULTS: Record<string, { daysPerWeek: number; weeksPerYear: number }> = {
  Giornaliera: { daysPerWeek: 5, weeksPerYear: 44 },
  Settimanale: { daysPerWeek: 1, weeksPerYear: 44 },
  Mensile: { daysPerWeek: 1, weeksPerYear: 12 },
  Trimestrale: { daysPerWeek: 1, weeksPerYear: 4 },
  Annuale: { daysPerWeek: 1, weeksPerYear: 1 },
}

function PeriodicityFields({ form, periodicity }: { form: UseFormReturn<ProcessFormValues>; periodicity: string }) {
  const { register, formState: { errors } } = form

  switch (periodicity) {
    case 'Giornaliera':
      return (
        <Flex gap={4}>
          <Field label="Ore / Minuti al Giorno *" error={errors.hoursPerDay?.message}>
            <Input type="number" step="0.25" {...register('hoursPerDay', { valueAsNumber: true })} min={0} max={24} placeholder="es. 2.5" {...inputStyle} />
          </Field>
          <Field label="Volte al Giorno *" error={errors.activitiesPerDay?.message}>
            <Input type="number" {...register('activitiesPerDay', { valueAsNumber: true })} min={1} placeholder="es. 3" {...inputStyle} />
          </Field>
          <Field label="Giorni / Settimana" error={errors.daysPerWeek?.message}>
            <Input type="number" {...register('daysPerWeek', { valueAsNumber: true })} min={1} max={7} {...inputStyle} />
          </Field>
          <Field label="Settimane / Anno" error={errors.weeksPerYear?.message}>
            <Input type="number" {...register('weeksPerYear', { valueAsNumber: true })} min={1} max={52} {...inputStyle} />
          </Field>
        </Flex>
      )
    case 'Settimanale':
      return (
        <Flex gap={4}>
          <Field label="Ore per Sessione *" error={errors.hoursPerDay?.message}>
            <Input type="number" step="0.25" {...register('hoursPerDay', { valueAsNumber: true })} min={0} max={24} placeholder="es. 3" {...inputStyle} />
          </Field>
          <Field label="Giorni / Settimana *" error={errors.daysPerWeek?.message}>
            <Input type="number" {...register('daysPerWeek', { valueAsNumber: true })} min={1} max={7} {...inputStyle} />
          </Field>
          <Field label="Settimane / Anno" error={errors.weeksPerYear?.message}>
            <Input type="number" {...register('weeksPerYear', { valueAsNumber: true })} min={1} max={52} {...inputStyle} />
          </Field>
        </Flex>
      )
    case 'Mensile':
      return (
        <Flex gap={4}>
          <Field label="Ore per Sessione *" error={errors.hoursPerDay?.message}>
            <Input type="number" step="0.25" {...register('hoursPerDay', { valueAsNumber: true })} min={0} max={24} placeholder="es. 4" {...inputStyle} />
          </Field>
          <Field label="Giorni / Mese *" error={errors.daysPerWeek?.message}>
            <Input type="number" {...register('daysPerWeek', { valueAsNumber: true })} min={1} max={31} {...inputStyle} />
          </Field>
          <Field label="Mesi / Anno" error={errors.weeksPerYear?.message}>
            <Input type="number" {...register('weeksPerYear', { valueAsNumber: true })} min={1} max={12} {...inputStyle} />
          </Field>
        </Flex>
      )
    case 'Trimestrale':
      return (
        <Flex gap={4}>
          <Field label="Ore per Sessione *" error={errors.hoursPerDay?.message}>
            <Input type="number" step="0.25" {...register('hoursPerDay', { valueAsNumber: true })} min={0} max={24} placeholder="es. 8" {...inputStyle} />
          </Field>
          <Field label="Giorni / Trimestre *" error={errors.daysPerWeek?.message}>
            <Input type="number" {...register('daysPerWeek', { valueAsNumber: true })} min={1} max={90} {...inputStyle} />
          </Field>
          <Field label="Trimestri / Anno" error={errors.weeksPerYear?.message}>
            <Input type="number" {...register('weeksPerYear', { valueAsNumber: true })} min={1} max={4} {...inputStyle} />
          </Field>
        </Flex>
      )
    case 'Annuale':
      return (
        <Flex gap={4}>
          <Field label="Ore per Sessione *" error={errors.hoursPerDay?.message}>
            <Input type="number" step="0.25" {...register('hoursPerDay', { valueAsNumber: true })} min={0} max={24} placeholder="es. 8" {...inputStyle} />
          </Field>
          <Field label="Giorni / Anno *" error={errors.daysPerWeek?.message}>
            <Input type="number" {...register('daysPerWeek', { valueAsNumber: true })} min={1} max={365} {...inputStyle} />
          </Field>
        </Flex>
      )
    default:
      return (
        <Flex bg="#f5f5f7" borderRadius="10px" px={4} py={3}>
          <Text fontSize="13px" color="#86868b">
            Seleziona la periodicita nel tab "Informazioni Generali" per configurare i parametri operativi.
          </Text>
        </Flex>
      )
  }
}

export default function CostsTab({ form }: Props) {
  const { register, watch, setValue, formState: { errors } } = form

  const periodicity = watch('periodicity')
  const annualSalary = watch('annualSalary')
  const hoursPerDay = watch('hoursPerDay')
  const daysPerWeek = watch('daysPerWeek')
  const weeksPerYear = watch('weeksPerYear')
  const activitiesPerDay = watch('activitiesPerDay')

  // Auto-compute hourly cost from RAL (1720 ore/anno standard italiano)
  useEffect(() => {
    const ral = annualSalary || 0
    setValue('hourlyCost', ral > 0 ? Math.round((ral / 1720) * 100) / 100 : 0)
  }, [annualSalary, setValue])

  // Set defaults only when periodicity actually changes
  const prevPeriodicity = useRef(periodicity)
  useEffect(() => {
    if (periodicity === prevPeriodicity.current) return
    prevPeriodicity.current = periodicity
    const defaults = PERIODICITY_DEFAULTS[periodicity]
    if (defaults) {
      setValue('daysPerWeek', defaults.daysPerWeek)
      setValue('weeksPerYear', defaults.weeksPerYear)
      if (periodicity !== 'Giornaliera') {
        setValue('activitiesPerDay', 1)
      }
    }
  }, [periodicity, setValue])

  // Auto-compute backend fields
  useEffect(() => {
    const h = hoursPerDay || 0
    const d = daysPerWeek || 0
    const w = periodicity === 'Annuale' ? 1 : (weeksPerYear || 0)
    const acts = periodicity === 'Giornaliera' ? (activitiesPerDay || 1) : 1
    setValue('timePerActivity', Math.round((h / acts) * 60))
    if (periodicity !== 'Giornaliera') setValue('activitiesPerDay', 1)
    setValue('workingDaysPerYear', d * w)
  }, [hoursPerDay, daysPerWeek, weeksPerYear, activitiesPerDay, periodicity, setValue])

  const totalDays = (daysPerWeek || 0) * (periodicity === 'Annuale' ? 1 : (weeksPerYear || 0))
  const totalHoursYear = (hoursPerDay || 0) * totalDays * (periodicity === 'Giornaliera' ? (activitiesPerDay || 1) : 1)

  return (
    <Flex direction="column" gap={5}>
      <Box>
        <Text fontSize="15px" fontWeight="700" color="#1d1d1f" mb={4}>Costi di Implementazione</Text>
        <Flex gap={4}>
          <Field label="Costo Implementazione (EUR) *" error={errors.implementationCost?.message}>
            <Input type="number" step="0.01" {...register('implementationCost', { valueAsNumber: true })} {...inputStyle} />
          </Field>
          <Field label="Costo Formazione (EUR)" error={errors.trainingCost?.message}>
            <Input type="number" step="0.01" {...register('trainingCost', { valueAsNumber: true })} {...inputStyle} />
          </Field>
          <Field label="Costo Manutenzione Annuale (EUR)" error={errors.maintenanceCost?.message}>
            <Input type="number" step="0.01" {...register('maintenanceCost', { valueAsNumber: true })} {...inputStyle} />
          </Field>
        </Flex>
      </Box>

      <Box h="1px" bg="#f0f0f2" />

      <Box>
        <Flex align="center" gap={3} mb={4}>
          <Text fontSize="15px" fontWeight="700" color="#1d1d1f">Parametri Operativi</Text>
          {periodicity && (
            <Text fontSize="12px" color="#007aff" fontWeight="500">
              Periodicita: {periodicity}
            </Text>
          )}
        </Flex>

        <Flex direction="column" gap={4}>
          <Flex gap={4} align="end">
            <Field label="RAL Media Dipendente (EUR) *" error={errors.annualSalary?.message}>
              <Input type="number" step="100" {...register('annualSalary', { valueAsNumber: true })} placeholder="es. 30000" {...inputStyle} />
            </Field>
            {(annualSalary || 0) > 0 && (
              <Box pb={1}>
                <Text fontSize="11px" color="#86868b">Costo orario calcolato</Text>
                <Text fontSize="16px" fontWeight="700" color="#007aff">
                  {(Math.round(((annualSalary || 0) / 1720) * 100) / 100).toFixed(2)} EUR/h
                </Text>
              </Box>
            )}
          </Flex>

          <PeriodicityFields form={form} periodicity={periodicity} />

          {totalHoursYear > 0 && (
            <Flex gap={4} bg="#f5f5f7" borderRadius="10px" px={4} py={2.5}>
              <Box flex={1}>
                <Text fontSize="11px" color="#86868b">Giorni lavorativi / anno</Text>
                <Text fontSize="14px" fontWeight="600" color="#1d1d1f">{totalDays}</Text>
              </Box>
              <Box flex={1}>
                <Text fontSize="11px" color="#86868b">Ore totali / anno</Text>
                <Text fontSize="14px" fontWeight="600" color="#007aff">{Math.round(totalHoursYear)}</Text>
              </Box>
            </Flex>
          )}
        </Flex>
      </Box>

      <Box h="1px" bg="#f0f0f2" />

      <Box>
        <Text fontSize="15px" fontWeight="700" color="#1d1d1f" mb={4}>Produttivita</Text>
        <Flex gap={4}>
          <Field label="Fattore Riduzione Tempo (%)">
            <Input type="number" {...register('timeReductionFactor', { valueAsNumber: true })} min={0} max={100} {...inputStyle} />
          </Field>
          <Field label="Fattore Aumento Produttivita (x)">
            <Input type="number" step="0.1" {...register('productivityFactor', { valueAsNumber: true })} min={1} {...inputStyle} />
          </Field>
        </Flex>
      </Box>
    </Flex>
  )
}

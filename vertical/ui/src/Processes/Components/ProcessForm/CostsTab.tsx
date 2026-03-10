import { Box, Flex, Input, Text } from '@chakra-ui/react'
import { useEffect } from 'react'
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

export default function CostsTab({ form }: Props) {
  const { register, watch, setValue, formState: { errors } } = form

  const hoursPerDay = watch('hoursPerDay')
  const daysPerWeek = watch('daysPerWeek')
  const weeksPerYear = watch('weeksPerYear')

  useEffect(() => {
    const h = hoursPerDay || 0
    const d = daysPerWeek || 0
    const w = weeksPerYear || 0
    setValue('timePerActivity', Math.round(h * 60))
    setValue('activitiesPerDay', 1)
    setValue('workingDaysPerYear', d * w)
  }, [hoursPerDay, daysPerWeek, weeksPerYear, setValue])

  const totalDays = (daysPerWeek || 0) * (weeksPerYear || 0)
  const totalHoursYear = (hoursPerDay || 0) * totalDays

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
        <Text fontSize="15px" fontWeight="700" color="#1d1d1f" mb={4}>Parametri Operativi</Text>
        <Flex direction="column" gap={4}>
          <Flex gap={4}>
            <Field label="Costo Orario Personale (EUR/h) *" error={errors.hourlyCost?.message}>
              <Input type="number" step="0.01" {...register('hourlyCost', { valueAsNumber: true })} {...inputStyle} />
            </Field>
            <Field label="Ore al Giorno dedicate *" error={errors.hoursPerDay?.message}>
              <Input type="number" step="0.5" {...register('hoursPerDay', { valueAsNumber: true })} min={0} max={24} {...inputStyle} />
            </Field>
          </Flex>
          <Flex gap={4}>
            <Field label="Giorni alla Settimana *" error={errors.daysPerWeek?.message}>
              <Input type="number" {...register('daysPerWeek', { valueAsNumber: true })} min={1} max={7} {...inputStyle} />
            </Field>
            <Field label="Settimane all'Anno *" error={errors.weeksPerYear?.message}>
              <Input type="number" {...register('weeksPerYear', { valueAsNumber: true })} min={1} max={52} {...inputStyle} />
            </Field>
          </Flex>
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
    </Flex>
  )
}

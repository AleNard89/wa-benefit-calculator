import { Box, Flex, Input, Text } from '@chakra-ui/react'
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
  const { register, formState: { errors } } = form

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
            <Field label="Tempo per Attivita (min) *" error={errors.timePerActivity?.message}>
              <Input type="number" {...register('timePerActivity', { valueAsNumber: true })} min={1} {...inputStyle} />
            </Field>
          </Flex>
          <Flex gap={4}>
            <Field label="Attivita al Giorno *" error={errors.activitiesPerDay?.message}>
              <Input type="number" {...register('activitiesPerDay', { valueAsNumber: true })} min={1} {...inputStyle} />
            </Field>
            <Field label="Giorni Lavorativi / Anno *" error={errors.workingDaysPerYear?.message}>
              <Input type="number" {...register('workingDaysPerYear', { valueAsNumber: true })} min={1} {...inputStyle} />
            </Field>
          </Flex>
        </Flex>
      </Box>
    </Flex>
  )
}

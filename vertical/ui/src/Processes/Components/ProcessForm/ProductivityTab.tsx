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

function Field({ label, hint, error, children }: { label: string; hint?: string; error?: string; children: React.ReactNode }) {
  return (
    <Box flex={1}>
      <Text fontSize="12px" fontWeight="600" color="#86868b" mb={0.5}>{label}</Text>
      {hint && <Text fontSize="11px" color="#c7c7cc" mb={1}>{hint}</Text>}
      {!hint && <Box mb={1} />}
      {children}
      {error && <Text fontSize="11px" color="#ff3b30" mt={1}>{error}</Text>}
    </Box>
  )
}

export default function ProductivityTab({ form }: Props) {
  const { register, formState: { errors } } = form

  return (
    <Flex direction="column" gap={5}>
      <Box>
        <Text fontSize="15px" fontWeight="700" color="#1d1d1f" mb={4}>Riduzione Errori</Text>
        <Flex gap={4}>
          <Field label="Tasso Errore Attuale (%)" error={errors.currentErrorRate?.message}>
            <Input type="number" step="0.1" {...register('currentErrorRate', { valueAsNumber: true })} min={0} max={100} {...inputStyle} />
          </Field>
          <Field label="Tasso Errore Post-Automazione (%)" error={errors.postErrorRate?.message}>
            <Input type="number" step="0.1" {...register('postErrorRate', { valueAsNumber: true })} min={0} max={100} {...inputStyle} />
          </Field>
          <Field label="Costo per Errore (EUR)" error={errors.errorCost?.message}>
            <Input type="number" step="0.01" {...register('errorCost', { valueAsNumber: true })} min={0} {...inputStyle} />
          </Field>
        </Flex>
      </Box>

      <Box h="1px" bg="#f0f0f2" />

      <Box>
        <Text fontSize="15px" fontWeight="700" color="#1d1d1f" mb={4}>Produttivita</Text>
        <Flex gap={4}>
          <Field label="Fattore Aumento Produttivita (x)" hint="Es. 2.0 = le attivita raddoppiano" error={errors.productivityFactor?.message}>
            <Input type="number" step="0.1" {...register('productivityFactor', { valueAsNumber: true })} min={1} {...inputStyle} />
          </Field>
          <Field label="Fattore Riduzione Tempo (%)" hint="Es. 50 = il tempo si dimezza" error={errors.timeReductionFactor?.message}>
            <Input type="number" {...register('timeReductionFactor', { valueAsNumber: true })} min={0} max={100} {...inputStyle} />
          </Field>
        </Flex>
      </Box>
    </Flex>
  )
}

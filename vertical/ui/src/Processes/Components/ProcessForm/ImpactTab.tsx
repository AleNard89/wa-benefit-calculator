import { Box, Flex, Text } from '@chakra-ui/react'
import { Slider } from '@chakra-ui/react'
import { Controller, type UseFormReturn } from 'react-hook-form'
import type { ProcessFormValues } from '../../Forms/ProcessFormSchema'

interface Props {
  form: UseFormReturn<ProcessFormValues>
}

interface ScoreFieldProps {
  label: string
  name: keyof ProcessFormValues
  form: UseFormReturn<ProcessFormValues>
}

function ScoreField({ label, name, form }: ScoreFieldProps) {
  const value = form.watch(name) as number
  return (
    <Box bg="#f5f5f7" borderRadius="12px" p={4}>
      <Flex justify="space-between" mb={3}>
        <Text fontSize="13px" fontWeight="600" color="#1d1d1f">{label}</Text>
        <Flex
          w="28px"
          h="28px"
          align="center"
          justify="center"
          borderRadius="8px"
          bg={value >= 4 ? '#34c75920' : value >= 3 ? '#ff950020' : '#ff3b3020'}
          flexShrink={0}
        >
          <Text
            fontSize="13px"
            fontWeight="700"
            color={value >= 4 ? '#34c759' : value >= 3 ? '#ff9500' : '#ff3b30'}
          >
            {value}
          </Text>
        </Flex>
      </Flex>
      <Controller
        name={name}
        control={form.control}
        render={({ field }) => (
          <Slider.Root
            min={1}
            max={5}
            step={1}
            value={[field.value as number]}
            onValueChange={({ value }) => field.onChange(value[0])}
          >
            <Slider.Control>
              <Slider.Track h="6px" borderRadius="3px" bg="#e8e8ed">
                <Slider.Range bg="#007aff" borderRadius="3px" />
              </Slider.Track>
              <Slider.Thumb
                index={0}
                w="22px"
                h="22px"
                bg="white"
                borderRadius="50%"
                boxShadow="0 1px 4px rgba(0,0,0,0.15)"
                border="2px solid #007aff"
              />
            </Slider.Control>
          </Slider.Root>
        )}
      />
      <Flex justify="space-between" mt={1.5}>
        <Text fontSize="10px" color="#c7c7cc">Basso</Text>
        <Text fontSize="10px" color="#c7c7cc">Alto</Text>
      </Flex>
    </Box>
  )
}

export default function ImpactTab({ form }: Props) {
  return (
    <Flex direction="column" gap={5}>
      <Text fontSize="15px" fontWeight="700" color="#1d1d1f">Valutazione Impatto</Text>

      <Flex direction="column" gap={3}>
        <Flex gap={3}>
          <Box flex={1}><ScoreField label="Qualita del Dato" name="dataQualityScore" form={form} /></Box>
          <Box flex={1}><ScoreField label="Audit e Compliance" name="auditScore" form={form} /></Box>
        </Flex>
        <Flex gap={3}>
          <Box flex={1}><ScoreField label="Esperienza Cliente" name="customerExperienceScore" form={form} /></Box>
          <Box flex={1}><ScoreField label="Riduzione degli Errori" name="errorReductionScore" form={form} /></Box>
        </Flex>
        <Flex gap={3}>
          <Box flex={1}><ScoreField label="Standardizzazione" name="standardizationScore" form={form} /></Box>
          <Box flex={1}><ScoreField label="Scalabilita" name="scalabilityScore" form={form} /></Box>
        </Flex>
      </Flex>
    </Flex>
  )
}

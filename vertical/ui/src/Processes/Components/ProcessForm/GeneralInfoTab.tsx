import { Box, Flex, Input, Text } from '@chakra-ui/react'
import { useState } from 'react'
import { type UseFormReturn, Controller } from 'react-hook-form'
import { LuX } from 'react-icons/lu'
import type { ProcessFormValues } from '../../Forms/ProcessFormSchema'
import { useAreasQuery } from '@/Orgs/Services/AreaApi'
import { useBotNamesQuery } from '@/Orchestrator/api'

interface Props {
  form: UseFormReturn<ProcessFormValues>
}

const processTypes = ['Manuale', 'Semi-automatico', 'Automatizzabile', 'Complesso']
const periodicities = ['Giornaliera', 'Settimanale', 'Mensile', 'Trimestrale', 'Annuale']
const technologies = ['SAP', 'Excel', 'Web App', 'ERP', 'CRM', 'Altro']

const inputStyle = {
  bg: '#f5f5f7',
  border: 'none',
  borderRadius: '10px',
  h: '42px',
  fontSize: '14px',
  _focus: { bg: 'white', boxShadow: '0 0 0 3px rgba(0,122,255,0.15)' },
}

const selectStyle = {
  bg: '#f5f5f7',
  border: 'none',
  borderRadius: '10px',
  h: '42px',
  fontSize: '14px',
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

function Select({ children, ...props }: React.SelectHTMLAttributes<HTMLSelectElement> & Record<string, unknown>) {
  return (
    <Box
      as="select"
      w="100%"
      px={3}
      cursor="pointer"
      {...selectStyle}
      _focus={{ bg: 'white', boxShadow: '0 0 0 3px rgba(0,122,255,0.15)' }}
      {...props}
    >
      {children}
    </Box>
  )
}

export default function GeneralInfoTab({ form }: Props) {
  const { register, setValue, formState: { errors } } = form
  const { data: areas } = useAreasQuery()
  const { data: botNames } = useBotNamesQuery()
  const [botSearch, setBotSearch] = useState('')

  return (
    <Flex direction="column" gap={5}>
      <Text fontSize="15px" fontWeight="700" color="#1d1d1f">Informazioni Generali</Text>

      <Flex gap={4}>
        <Field label="Nome Processo *" error={errors.processName?.message}>
          <Input {...register('processName')} placeholder="Nome del processo" {...inputStyle} />
        </Field>
        <Field label="Descrizione">
          <Input {...register('processDescription')} placeholder="Descrizione" {...inputStyle} />
        </Field>
      </Flex>

      <Flex gap={4}>
        <Field label="Proponente *" error={errors.proposer?.message}>
          <Input {...register('proposer')} placeholder="Nome proponente" {...inputStyle} />
        </Field>
        <Field label="Area *" error={errors.area?.message}>
          <Select
            {...register('area')}
            onChange={(e) => {
              const name = (e.target as HTMLSelectElement).value
              setValue('area', name)
              const area = areas?.find((a) => a.name === name)
              setValue('areaId', area?.id ?? null)
            }}
          >
            <option value="">Seleziona area...</option>
            {areas?.map((a) => <option key={a.id} value={a.name}>{a.name}</option>)}
          </Select>
        </Field>
      </Flex>

      <Flex gap={4}>
        <Field label="Manager Responsabile *" error={errors.responsibleManager?.message}>
          <Input {...register('responsibleManager')} placeholder="Manager responsabile" {...inputStyle} />
        </Field>
        <Field label="Reparto">
          <Input {...register('department')} placeholder="Reparto" {...inputStyle} />
        </Field>
      </Flex>

      <Flex gap={4}>
        <Field label="Numero Sistemi Coinvolti *" error={errors.systemsInvolved?.message}>
          <Input type="number" {...register('systemsInvolved', { valueAsNumber: true })} min={1} {...inputStyle} />
        </Field>
        <Field label="Tipo Processo *" error={errors.processType?.message}>
          <Select {...register('processType')}>
            <option value="">Seleziona...</option>
            {processTypes.map((t) => <option key={t} value={t}>{t}</option>)}
          </Select>
        </Field>
      </Flex>

      <Flex gap={4}>
        <Field label="Periodicita *" error={errors.periodicity?.message}>
          <Select {...register('periodicity')}>
            <option value="">Seleziona...</option>
            {periodicities.map((p) => <option key={p} value={p}>{p}</option>)}
          </Select>
        </Field>
        <Field label="Tecnologie *" error={errors.technology?.message as string}>
          <Controller
            name="technology"
            control={form.control}
            render={({ field }) => (
              <Flex direction="column" gap={2}>
                <Flex gap={2} wrap="wrap">
                  {technologies.map((t) => {
                    const selected = (field.value ?? []).includes(t)
                    return (
                      <Box
                        key={t}
                        as="button"
                        type="button"
                        px={3}
                        py={1.5}
                        borderRadius="8px"
                        fontSize="13px"
                        fontWeight={selected ? '600' : '400'}
                        bg={selected ? '#007aff' : '#f5f5f7'}
                        color={selected ? 'white' : '#1d1d1f'}
                        cursor="pointer"
                        transition="all 0.15s"
                        _hover={{ opacity: 0.85 }}
                        onClick={() => {
                          const curr = field.value ?? []
                          field.onChange(
                            selected ? curr.filter((v: string) => v !== t) : [...curr, t]
                          )
                        }}
                      >
                        {t}
                      </Box>
                    )
                  })}
                </Flex>
                {(field.value ?? []).includes('Altro') && (
                  <Input
                    placeholder="Specifica tecnologia..."
                    {...inputStyle}
                    h="36px"
                    fontSize="13px"
                    {...register('technologyOther')}
                  />
                )}
              </Flex>
            )}
          />
        </Field>
      </Flex>

      <Box h="1px" bg="#f0f0f2" />
      <Text fontSize="15px" fontWeight="700" color="#1d1d1f">Bot Collegati</Text>
      <Field label="Seleziona i bot Orchestrator associati a questo processo">
        <Controller
          name="linkedBots"
          control={form.control}
          render={({ field }) => {
            const selected = field.value ?? []
            const available = botNames ?? []
            const filtered = available.filter(
              (b) => !selected.includes(b) && b.toLowerCase().includes(botSearch.toLowerCase())
            )
            const showDropdown = botSearch.length > 0 && filtered.length > 0
            return (
              <Flex direction="column" gap={2}>
                {selected.length > 0 && (
                  <Flex gap={2} wrap="wrap">
                    {selected.map((bot: string) => (
                      <Flex
                        key={bot}
                        align="center"
                        gap={1}
                        px={3}
                        py={1.5}
                        borderRadius="8px"
                        fontSize="13px"
                        fontWeight="500"
                        bg="#007aff"
                        color="white"
                      >
                        {bot}
                        <Box
                          as="button"
                          type="button"
                          display="flex"
                          cursor="pointer"
                          onClick={() => field.onChange(selected.filter((v: string) => v !== bot))}
                        >
                          <LuX size={12} />
                        </Box>
                      </Flex>
                    ))}
                  </Flex>
                )}
                <Box position="relative">
                  <Input
                    placeholder={available.length > 0 ? 'Cerca bot per nome...' : 'Nessun bot sincronizzato'}
                    value={botSearch}
                    onChange={(e) => setBotSearch(e.target.value)}
                    disabled={available.length === 0}
                    {...inputStyle}
                    h="36px"
                    fontSize="13px"
                  />
                  {showDropdown && (
                    <Flex
                      direction="column"
                      position="absolute"
                      top="100%"
                      left={0}
                      right={0}
                      mt={1}
                      bg="white"
                      border="1px solid #e8e8ed"
                      borderRadius="10px"
                      maxH="200px"
                      overflowY="auto"
                      boxShadow="0 4px 12px rgba(0,0,0,0.08)"
                      zIndex={10}
                    >
                      {filtered.length > 1 && (
                        <Box
                          as="button"
                          type="button"
                          px={3}
                          py={2}
                          fontSize="13px"
                          fontWeight="600"
                          color="#007aff"
                          textAlign="left"
                          cursor="pointer"
                          borderBottom="1px solid #f0f0f2"
                          _hover={{ bg: '#f0f5ff' }}
                          onClick={() => {
                            field.onChange([...selected, ...filtered])
                            setBotSearch('')
                          }}
                        >
                          Seleziona tutti ({filtered.length})
                        </Box>
                      )}
                      {filtered.map((bot) => (
                        <Box
                          key={bot}
                          as="button"
                          type="button"
                          px={3}
                          py={2}
                          fontSize="13px"
                          textAlign="left"
                          cursor="pointer"
                          _hover={{ bg: '#f5f5f7' }}
                          onClick={() => {
                            field.onChange([...selected, bot])
                            setBotSearch('')
                          }}
                        >
                          {bot}
                        </Box>
                      ))}
                    </Flex>
                  )}
                </Box>
              </Flex>
            )
          }}
        />
      </Field>
      <Field label="Note sui bot (ruolo di ogni bot, visibile al chatbot)">
        <Box
          as="textarea"
          w="100%"
          rows={3}
          px={3}
          py={2}
          {...inputStyle}
          h="auto"
          placeholder="Es: BotA_DSP: lettura e validazione dati. BotB_PRF: inserimento dati nel gestionale"
          {...register('botNotes')}
        />
      </Field>
    </Flex>
  )
}

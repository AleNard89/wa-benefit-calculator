import { Box, Button, Flex, Text } from '@chakra-ui/react'
import { zodResolver } from '@hookform/resolvers/zod'
import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { useNavigate } from 'react-router-dom'
import { LuArrowLeft, LuSave } from 'react-icons/lu'

import { toaster } from '@/Snippets/toaster'
import GeneralInfoTab from '../Components/ProcessForm/GeneralInfoTab'
import CostsTab from '../Components/ProcessForm/CostsTab'
import ProductivityTab from '../Components/ProcessForm/ProductivityTab'
import ImpactTab from '../Components/ProcessForm/ImpactTab'
import { processFormDefaults, processFormSchema, type ProcessFormValues } from '../Forms/ProcessFormSchema'
import { useCreateProcessMutation } from '../Services/Api'

const tabs = [
  { key: 'general', label: 'Informazioni Generali' },
  { key: 'costs', label: 'Costi e Operativita' },
  { key: 'productivity', label: 'Errori e Produttivita' },
  { key: 'impact', label: 'Valutazione Impatto' },
] as const

type TabKey = (typeof tabs)[number]['key']

export default function BenefitCalculationView() {
  const navigate = useNavigate()
  const [createProcess, { isLoading }] = useCreateProcessMutation()
  const [activeTab, setActiveTab] = useState<TabKey>('general')

  const form = useForm<ProcessFormValues>({
    resolver: zodResolver(processFormSchema),
    defaultValues: processFormDefaults,
    mode: 'onChange',
  })

  const onSubmit = async (data: ProcessFormValues) => {
    try {
      const process = await createProcess(data).unwrap()
      toaster.success({ title: 'Processo creato', description: `${process.processName} salvato con successo` })
      navigate(`/processes/${process.id}`)
    } catch {
      toaster.error({ title: 'Errore', description: 'Impossibile salvare il processo' })
    }
  }

  const tabContent: Record<TabKey, React.ReactNode> = {
    general: <GeneralInfoTab form={form} />,
    costs: <CostsTab form={form} />,
    productivity: <ProductivityTab form={form} />,
    impact: <ImpactTab form={form} />,
  }

  return (
    <Box>
      {/* Header */}
      <Flex justify="space-between" align="center" mb={5}>
        <Box>
          <Text fontSize="20px" fontWeight="700" color="#1d1d1f" letterSpacing="-0.3px">
            Nuova Proposta
          </Text>
          <Text fontSize="14px" color="#86868b" mt={1}>
            Compila i dati per calcolare i benefici dell&apos;automazione
          </Text>
        </Box>
        <Flex gap={2}>
          <Button
            bg="#f5f5f7"
            color="#1d1d1f"
            borderRadius="10px"
            h="36px"
            px={4}
            fontSize="13px"
            _hover={{ bg: '#e8e8ed' }}
            onClick={() => navigate('/processes/list')}
          >
            <Flex align="center" gap={1.5}><LuArrowLeft size={15} /> Lista Processi</Flex>
          </Button>
          <Button
            bg="#007aff"
            color="white"
            borderRadius="10px"
            h="36px"
            px={5}
            fontSize="13px"
            fontWeight="600"
            _hover={{ bg: '#0066d6' }}
            loading={isLoading}
            onClick={form.handleSubmit(onSubmit)}
          >
            <Flex align="center" gap={1.5}><LuSave size={15} /> Salva Processo</Flex>
          </Button>
        </Flex>
      </Flex>

      {/* Pill Tabs */}
      <Flex
        bg="#f5f5f7"
        borderRadius="12px"
        p="4px"
        gap="2px"
        mb={5}
        w="fit-content"
      >
        {tabs.map((tab) => (
          <Box
            key={tab.key}
            as="button"
            type="button"
            px={4}
            py={2}
            borderRadius="10px"
            fontSize="13px"
            fontWeight={activeTab === tab.key ? '600' : '400'}
            color={activeTab === tab.key ? '#1d1d1f' : '#86868b'}
            bg={activeTab === tab.key ? 'white' : 'transparent'}
            boxShadow={activeTab === tab.key ? '0 1px 4px rgba(0,0,0,0.08)' : 'none'}
            cursor="pointer"
            transition="all 0.2s"
            _hover={{ color: '#1d1d1f' }}
            onClick={() => setActiveTab(tab.key)}
          >
            {tab.label}
          </Box>
        ))}
      </Flex>

      {/* Form Card */}
      <Box
        as="form"
        onSubmit={form.handleSubmit(onSubmit)}
        bg="white"
        borderRadius="16px"
        boxShadow="0 2px 20px rgba(0,0,0,0.06)"
        p={7}
      >
        {tabContent[activeTab]}
      </Box>
    </Box>
  )
}

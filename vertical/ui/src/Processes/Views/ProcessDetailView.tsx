import { Badge, Box, Button, Dialog, Flex, Heading, Portal, Spinner, Text } from '@chakra-ui/react'
import { zodResolver } from '@hookform/resolvers/zod'
import { useRef, useState, useCallback } from 'react'
import { useForm } from 'react-hook-form'
import { useSelector } from 'react-redux'
import { useNavigate, useParams } from 'react-router-dom'
import { LuPencil, LuX, LuSave } from 'react-icons/lu'
import { selectAccessToken } from '@/Auth/Redux'

import Config from '@/Config'
import { toaster } from '@/Snippets/toaster'
import ProcessResults from '../Components/ProcessResults'
import GeneralInfoTab from '../Components/ProcessForm/GeneralInfoTab'
import CostsTab from '../Components/ProcessForm/CostsTab'
import ProductivityTab from '../Components/ProcessForm/ProductivityTab'
import ImpactTab from '../Components/ProcessForm/ImpactTab'
import { processFormSchema, type ProcessFormValues } from '../Forms/ProcessFormSchema'
import { useDeleteDocumentMutation, useDeleteProcessMutation, useProcessQuery, useUpdateProcessMutation, useUpdateProcessStatusMutation, useUploadDocumentMutation } from '../Services/Api'
import { formatCurrency, statusColorMap, statusLabelMap } from '../Utils'
import type { ProcessResults as ProcessResultsType, ProcessStatus } from '../Types'

import { NativeSelect } from '@chakra-ui/react'

const allStatuses: ProcessStatus[] = ['To Valuate', 'Analysis', 'Ongoing', 'Production']

const editTabs = [
  { key: 'general', label: 'Informazioni Generali' },
  { key: 'costs', label: 'Costi e Operativita' },
  { key: 'productivity', label: 'Errori e Produttivita' },
  { key: 'impact', label: 'Valutazione Impatto' },
] as const

type EditTabKey = (typeof editTabs)[number]['key']

export default function ProcessDetailView() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const { data: process, isLoading } = useProcessQuery(Number(id))
  const [updateStatus] = useUpdateProcessStatusMutation()
  const [deleteProcess] = useDeleteProcessMutation()
  const [uploadDocument, { isLoading: isUploading }] = useUploadDocumentMutation()
  const [deleteDocument] = useDeleteDocumentMutation()
  const token = useSelector(selectAccessToken)
  const [updateProcess, { isLoading: isUpdating }] = useUpdateProcessMutation()
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false)
  const [isEditing, setIsEditing] = useState(false)
  const [editTab, setEditTab] = useState<EditTabKey>('general')
  const cancelRef = useRef<HTMLButtonElement>(null)
  const fileInputRef = useRef<HTMLInputElement>(null)

  const form = useForm<ProcessFormValues>({
    resolver: zodResolver(processFormSchema),
    mode: 'onChange',
  })

  const enterEditMode = () => {
    if (!process) return
    const d = process.data
    form.reset({
      processName: process.processName,
      processDescription: d?.processDescription ?? '',
      proposer: d?.proposer ?? '',
      area: d?.area ?? '',
      areaId: process.areaId ?? null,
      responsibleManager: d?.responsibleManager ?? '',
      department: d?.department ?? '',
      systemsInvolved: d?.systemsInvolved ?? 1,
      processType: d?.processType ?? '',
      periodicity: d?.periodicity ?? '',
      frequentChanges: d?.frequentChanges ?? false,
      technology: Array.isArray(d?.technology) ? d.technology : (d?.technology ? [d.technology] : []),
      technologyOther: d?.technologyOther ?? '',
      linkedBots: d?.linkedBots ?? [],
      botNotes: d?.botNotes ?? '',
      implementationCost: d?.implementationCost ?? 0,
      trainingCost: d?.trainingCost ?? 0,
      maintenanceCost: d?.maintenanceCost ?? 0,
      hourlyCost: d?.hourlyCost ?? 0,
      timePerActivity: d?.timePerActivity ?? 0,
      activitiesPerDay: d?.activitiesPerDay ?? 0,
      workingDaysPerYear: d?.workingDaysPerYear ?? 220,
      currentErrorRate: d?.currentErrorRate ?? 0,
      postErrorRate: d?.postErrorRate ?? 0,
      errorCost: d?.errorCost ?? 0,
      productivityFactor: d?.productivityFactor ?? 2,
      timeReductionFactor: d?.timeReductionFactor ?? 50,
      dataQualityScore: d?.dataQualityScore ?? 3,
      auditScore: d?.auditScore ?? 3,
      customerExperienceScore: d?.customerExperienceScore ?? 3,
      errorReductionScore: d?.errorReductionScore ?? 3,
      standardizationScore: d?.standardizationScore ?? 3,
      scalabilityScore: d?.scalabilityScore ?? 3,
    })
    setEditTab('general')
    setIsEditing(true)
  }

  const handleSave = async (data: ProcessFormValues) => {
    if (!process) return
    try {
      await updateProcess({ id: process.id, body: data }).unwrap()
      toaster.success({ title: 'Processo aggiornato' })
      setIsEditing(false)
    } catch {
      toaster.error({ title: 'Errore aggiornamento' })
    }
  }

  const handleDownloadDocument = useCallback(async () => {
    if (!process) return
    try {
      const res = await fetch(`${Config.api.basePath}/processes/${process.id}/document`, {
        headers: { Authorization: `Bearer ${token}`, 'X-Company-Id': String(process.companyId) },
      })
      if (!res.ok) throw new Error()
      const blob = await res.blob()
      const url = URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = process.documentName || 'documento'
      a.click()
      URL.revokeObjectURL(url)
    } catch {
      toaster.error({ title: 'Errore download documento' })
    }
  }, [process, token])

  if (isLoading) {
    return <Flex justify="center" py={20}><Spinner size="xl" color="brand.300" /></Flex>
  }

  if (!process) {
    return <Text color="fg.muted" py={10}>Processo non trovato</Text>
  }

  const d = process.data
  const r = process.results

  const result: ProcessResultsType = {
    operationalSavings: r?.operationalSavings ?? 0,
    errorReductionSavings: r?.errorReductionSavings ?? 0,
    productivityBenefit: r?.productivityBenefit ?? 0,
    annualSavings: r?.annualSavings ?? 0,
    roi: r?.roi ?? 0,
    breakEvenMonths: r?.breakEvenMonths ?? null,
    hoursSavedMonthly: r?.hoursSavedMonthly ?? 0,
    hoursSavedAnnually: r?.hoursSavedAnnually ?? 0,
    impactScore: r?.impactScore ?? 0,
  }

  const scores = {
    dataQualityScore: d?.dataQualityScore ?? 3,
    auditScore: d?.auditScore ?? 3,
    customerExperienceScore: d?.customerExperienceScore ?? 3,
    errorReductionScore: d?.errorReductionScore ?? 3,
    standardizationScore: d?.standardizationScore ?? 3,
    scalabilityScore: d?.scalabilityScore ?? 3,
  }

  const handleStatusChange = async (status: string) => {
    try {
      await updateStatus({ id: process.id, status }).unwrap()
      toaster.success({ title: 'Stato aggiornato' })
    } catch {
      toaster.error({ title: 'Errore aggiornamento stato' })
    }
  }

  const handleFileUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (!file) return
    try {
      await uploadDocument({ id: process.id, file }).unwrap()
      toaster.success({ title: 'Documento caricato e indicizzato' })
    } catch {
      toaster.error({ title: 'Errore caricamento documento' })
    }
    if (fileInputRef.current) fileInputRef.current.value = ''
  }

  const handleDeleteDocument = async () => {
    try {
      await deleteDocument(process.id).unwrap()
      toaster.success({ title: 'Documento rimosso' })
    } catch {
      toaster.error({ title: 'Errore rimozione documento' })
    }
  }

  const handleDelete = async () => {
    try {
      await deleteProcess(process.id).unwrap()
      setDeleteDialogOpen(false)
      toaster.success({ title: 'Processo eliminato' })
      navigate('/processes/list')
    } catch {
      toaster.error({ title: 'Errore eliminazione' })
    }
  }

  return (
    <Flex direction="column" gap={6}>
      <Flex justify="space-between" align="center">
        <Box>
          <Heading size="lg">{process.processName}</Heading>
          <Flex gap={2} mt={1} align="center">
            <Text fontSize="sm" color="fg.muted">{d?.area} - {d?.proposer}</Text>
            <Badge colorPalette={statusColorMap[process.status] || 'gray'}>
              {statusLabelMap[process.status] || process.status}
            </Badge>
          </Flex>
        </Box>
        <Flex gap={2} align="center">
          {!isEditing && (
            <>
              <NativeSelect.Root maxW="180px" size="sm">
                <NativeSelect.Field
                  value={process.status}
                  onChange={(e) => handleStatusChange(e.target.value)}
                >
                  {allStatuses.map((s) => (
                    <option key={s} value={s}>{statusLabelMap[s]}</option>
                  ))}
                </NativeSelect.Field>
              </NativeSelect.Root>
              <Button size="sm" colorPalette="blue" onClick={enterEditMode}>
                <LuPencil size={14} /> Modifica
              </Button>
            </>
          )}
          {isEditing && (
            <>
              <Button size="sm" colorPalette="blue" loading={isUpdating} onClick={form.handleSubmit(handleSave)}>
                <LuSave size={14} /> Salva
              </Button>
              <Button size="sm" variant="outline" onClick={() => setIsEditing(false)}>
                <LuX size={14} /> Annulla
              </Button>
            </>
          )}
          <Button size="sm" variant="outline" onClick={() => navigate('/processes/list')}>
            Torna alla lista
          </Button>
          {!isEditing && (
            <Button size="sm" colorPalette="red" variant="outline" onClick={() => setDeleteDialogOpen(true)}>
              Elimina
            </Button>
          )}
        </Flex>
      </Flex>

      <Dialog.Root role="alertdialog" open={deleteDialogOpen} onOpenChange={(e) => setDeleteDialogOpen(e.open)} initialFocusEl={() => cancelRef.current}>
        <Portal>
          <Dialog.Backdrop />
          <Dialog.Positioner>
            <Dialog.Content>
              <Dialog.Header>
                <Dialog.Title>Conferma eliminazione</Dialog.Title>
              </Dialog.Header>
              <Dialog.Body>
                Sei sicuro di voler eliminare il processo <strong>{process.processName}</strong>?
              </Dialog.Body>
              <Dialog.Footer>
                <Button ref={cancelRef} variant="outline" onClick={() => setDeleteDialogOpen(false)}>
                  Annulla
                </Button>
                <Button colorPalette="red" onClick={handleDelete} ml={3}>
                  Elimina
                </Button>
              </Dialog.Footer>
            </Dialog.Content>
          </Dialog.Positioner>
        </Portal>
      </Dialog.Root>

      {isEditing && (
        <>
          <Flex bg="#f5f5f7" borderRadius="12px" p="4px" gap="2px" w="fit-content">
            {editTabs.map((tab) => (
              <Box
                key={tab.key}
                as="button"
                type="button"
                px={4}
                py={2}
                borderRadius="10px"
                fontSize="13px"
                fontWeight={editTab === tab.key ? '600' : '400'}
                color={editTab === tab.key ? '#1d1d1f' : '#86868b'}
                bg={editTab === tab.key ? 'white' : 'transparent'}
                boxShadow={editTab === tab.key ? '0 1px 4px rgba(0,0,0,0.08)' : 'none'}
                cursor="pointer"
                transition="all 0.2s"
                _hover={{ color: '#1d1d1f' }}
                onClick={() => setEditTab(tab.key)}
              >
                {tab.label}
              </Box>
            ))}
          </Flex>
          <Box bg="white" borderRadius="16px" boxShadow="0 2px 20px rgba(0,0,0,0.06)" p={7}>
            {editTab === 'general' && <GeneralInfoTab form={form} />}
            {editTab === 'costs' && <CostsTab form={form} />}
            {editTab === 'productivity' && <ProductivityTab form={form} />}
            {editTab === 'impact' && <ImpactTab form={form} />}
          </Box>
        </>
      )}

      {!isEditing && (
        <>
          <Flex gap={4} wrap="wrap">
            <InfoBox label="Manager" value={d?.responsibleManager ?? '-'} />
            <InfoBox label="Reparto" value={d?.department || '-'} />
            <InfoBox label="Tecnologia" value={
              (() => {
                const techs = (Array.isArray(d?.technology) ? d.technology : (d?.technology ? [d.technology] : [])).filter((t: string) => t !== 'Altro')
                const other = d?.technologyOther ? [d.technologyOther] : []
                const all = [...techs, ...other]
                return all.length > 0 ? all.join(', ') : '-'
              })()
            } />
            <InfoBox label="Tipo" value={d?.processType ?? '-'} />
            <InfoBox label="Periodicita" value={d?.periodicity ?? '-'} />
            <InfoBox label="Costo Impl." value={formatCurrency(d?.implementationCost ?? 0)} />
          </Flex>

          <Box p={4} borderRadius="md" bg="bg.subtle" borderWidth={1} borderColor="border.muted">
            <Text fontSize="sm" fontWeight="600" mb={2}>Documento</Text>
            <input
              ref={fileInputRef}
              type="file"
              accept=".pptx,.pdf,.docx"
              style={{ display: 'none' }}
              onChange={handleFileUpload}
            />
            {process.documentName ? (
              <Flex align="center" gap={3}>
                <Text fontSize="sm">{process.documentName}</Text>
                <Button size="xs" variant="outline" onClick={handleDownloadDocument}>
                  Scarica
                </Button>
                <Button size="xs" colorPalette="red" variant="ghost" onClick={handleDeleteDocument}>
                  Rimuovi
                </Button>
              </Flex>
            ) : (
              <Button size="sm" variant="outline" loading={isUploading} onClick={() => fileInputRef.current?.click()}>
                Carica PPTX / DOCX / PDF
              </Button>
            )}
          </Box>

          {d?.linkedBots && d.linkedBots.length > 0 && (
            <Box p={4} borderRadius="md" bg="bg.subtle" borderWidth={1} borderColor="border.muted">
              <Text fontSize="sm" fontWeight="600" mb={2}>Bot Collegati ({d.linkedBots.length})</Text>
              <Flex gap={2} wrap="wrap">
                {d.linkedBots.map((bot: string) => (
                  <Badge key={bot} colorPalette="blue" fontSize="12px" px={2} py={1}>
                    {bot}
                  </Badge>
                ))}
              </Flex>
            </Box>
          )}

          <ProcessResults result={result} scores={scores} />
        </>
      )}
    </Flex>
  )
}

function InfoBox({ label, value }: { label: string; value: string }) {
  return (
    <Box p={3} borderRadius="md" bg="bg.subtle" borderWidth={1} borderColor="border.muted" minW="140px">
      <Text fontSize="xs" color="fg.muted">{label}</Text>
      <Text fontSize="sm" fontWeight="600">{value}</Text>
    </Box>
  )
}

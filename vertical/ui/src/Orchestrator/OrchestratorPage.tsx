import { Box, Flex, Text } from '@chakra-ui/react'
import { useEffect, useRef, useState } from 'react'
import { LuCheck, LuChevronDown, LuX } from 'react-icons/lu'

import { useProcessNamesQuery } from './api'
import JobsTab from './JobsTab'
import QueueItemsTab from './QueueItemsTab'
import SchedulesTab from './SchedulesTab'

const tabs = [
  { id: 'jobs', label: 'Esecuzioni' },
  { id: 'schedules', label: 'Schedulazioni' },
  { id: 'queues', label: 'Code' },
]

function ProcessNameFilter({
  selected,
  onChange,
}: {
  selected: string[]
  onChange: (names: string[]) => void
}) {
  const { data: options = [] } = useProcessNamesQuery()
  const [isOpen, setIsOpen] = useState(false)
  const ref = useRef<HTMLDivElement>(null)

  useEffect(() => {
    function handleClickOutside(e: MouseEvent) {
      if (ref.current && !ref.current.contains(e.target as Node)) {
        setIsOpen(false)
      }
    }
    document.addEventListener('mousedown', handleClickOutside)
    return () => document.removeEventListener('mousedown', handleClickOutside)
  }, [])

  const toggle = (name: string) => {
    onChange(
      selected.includes(name)
        ? selected.filter((n) => n !== name)
        : [...selected, name],
    )
  }

  if (options.length === 0) return null

  return (
    <Box ref={ref} position="relative">
      <Flex
        as="button"
        align="center"
        gap={2}
        px={3}
        py={1.5}
        borderRadius="8px"
        border="1px solid"
        borderColor={selected.length > 0 ? '#007aff' : '#d2d2d7'}
        bg={selected.length > 0 ? '#007aff0a' : 'white'}
        cursor="pointer"
        fontSize="13px"
        color="#1d1d1f"
        transition="all 0.15s"
        _hover={{ borderColor: '#007aff' }}
        onClick={() => setIsOpen(!isOpen)}
      >
        <Text>
          {selected.length === 0
            ? 'Tutti i processi'
            : `${selected.length} process${selected.length === 1 ? 'o' : 'i'}`}
        </Text>
        <LuChevronDown size={14} style={{ opacity: 0.5 }} />
      </Flex>

      {selected.length > 0 && (
        <Flex
          as="button"
          align="center"
          justify="center"
          position="absolute"
          top="-6px"
          right="-6px"
          w="18px"
          h="18px"
          borderRadius="full"
          bg="#007aff"
          cursor="pointer"
          onClick={(e) => {
            e.stopPropagation()
            onChange([])
          }}
        >
          <LuX size={10} color="white" />
        </Flex>
      )}

      {isOpen && (
        <Box
          position="absolute"
          top="calc(100% + 4px)"
          left={0}
          minW="280px"
          maxH="320px"
          overflowY="auto"
          bg="white"
          borderRadius="12px"
          boxShadow="0 4px 24px rgba(0,0,0,0.12)"
          border="1px solid #e5e5ea"
          zIndex={20}
          py={1}
        >
          {options.map((name) => {
            const isSelected = selected.includes(name)
            return (
              <Flex
                key={name}
                align="center"
                gap={2.5}
                px={3}
                py={2}
                cursor="pointer"
                bg={isSelected ? '#007aff08' : 'transparent'}
                _hover={{ bg: isSelected ? '#007aff12' : '#f5f5f7' }}
                transition="background 0.1s"
                onClick={() => toggle(name)}
              >
                <Flex
                  align="center"
                  justify="center"
                  w="18px"
                  h="18px"
                  borderRadius="5px"
                  border="1.5px solid"
                  borderColor={isSelected ? '#007aff' : '#c7c7cc'}
                  bg={isSelected ? '#007aff' : 'transparent'}
                  flexShrink={0}
                  transition="all 0.15s"
                >
                  {isSelected && <LuCheck size={11} color="white" />}
                </Flex>
                <Text fontSize="13px" color="#1d1d1f" lineClamp={1}>
                  {name}
                </Text>
              </Flex>
            )
          })}
        </Box>
      )}
    </Box>
  )
}

export default function OrchestratorPage() {
  const [activeTab, setActiveTab] = useState('jobs')
  const [selectedProcessNames, setSelectedProcessNames] = useState<string[]>([])

  const processNamesParam =
    selectedProcessNames.length > 0 ? selectedProcessNames.join(',') : undefined

  const processFilterNode = (
    <ProcessNameFilter
      selected={selectedProcessNames}
      onChange={setSelectedProcessNames}
    />
  )

  return (
    <Box>
      <Flex justify="space-between" align="center" mb={5}>
        <Box>
          <Text fontSize="20px" fontWeight="700" color="#1d1d1f" letterSpacing="-0.3px">
            Orchestrator
          </Text>
          <Text fontSize="14px" color="#86868b" mt={1}>
            Monitoraggio esecuzioni bot e code
          </Text>
        </Box>
      </Flex>

      <Flex bg="#f5f5f7" borderRadius="12px" p="4px" gap="2px" mb={5} w="fit-content">
        {tabs.map((tab) => (
          <Flex
            key={tab.id}
            as="button"
            align="center"
            px={4}
            py={2}
            borderRadius="10px"
            fontSize="13px"
            fontWeight={activeTab === tab.id ? '600' : '400'}
            color={activeTab === tab.id ? '#1d1d1f' : '#86868b'}
            bg={activeTab === tab.id ? 'white' : 'transparent'}
            boxShadow={activeTab === tab.id ? '0 1px 4px rgba(0,0,0,0.08)' : 'none'}
            cursor="pointer"
            transition="all 0.2s"
            _hover={{ color: '#1d1d1f' }}
            onClick={() => setActiveTab(tab.id)}
          >
            {tab.label}
          </Flex>
        ))}
      </Flex>

      <Box bg="white" borderRadius="16px" boxShadow="0 2px 20px rgba(0,0,0,0.06)" p={6}>
        {activeTab === 'jobs' && (
          <JobsTab
            processNames={processNamesParam}
            processFilter={processFilterNode}
          />
        )}
        {activeTab === 'schedules' && (
          <SchedulesTab
            processNames={processNamesParam}
            processFilter={processFilterNode}
          />
        )}
        {activeTab === 'queues' && (
          <QueueItemsTab
            processNames={processNamesParam}
            processFilter={processFilterNode}
          />
        )}

      </Box>
    </Box>
  )
}

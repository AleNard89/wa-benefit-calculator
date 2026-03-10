import { Box, Flex, Heading, Text } from '@chakra-ui/react'
import { LuBot } from 'react-icons/lu'

export default function OrchestratorView() {
  return (
    <Flex direction="column" align="center" justify="center" py={20} gap={4}>
      <Flex
        w="64px"
        h="64px"
        align="center"
        justify="center"
        borderRadius="16px"
        bg="#007aff10"
        color="#007aff"
      >
        <LuBot size={32} />
      </Flex>
      <Heading size="lg">Connettori - Orchestrator</Heading>
      <Box textAlign="center" maxW="400px">
        <Text color="fg.muted">
          Monitoraggio dei bot: stato dei run, log di esecuzione e gestione delle code.
        </Text>
        <Text color="fg.muted" mt={2} fontSize="sm">
          In arrivo.
        </Text>
      </Box>
    </Flex>
  )
}

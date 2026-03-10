import { Box, Spinner, Text, VStack } from '@chakra-ui/react'

export default function StartupView() {
  return (
    <Box display="flex" alignItems="center" justifyContent="center" minH="100vh">
      <VStack gap={4}>
        <Spinner size="xl" color="brand.300" />
        <Text fontSize="lg" color="fg.muted">Caricamento...</Text>
      </VStack>
    </Box>
  )
}

import { Box, Flex } from '@chakra-ui/react'

import Sidebar from '../Components/Sidebar'

interface BaseLayoutProps {
  children: React.ReactNode
}

export default function BaseLayout({ children }: BaseLayoutProps) {
  return (
    <Flex minH="100vh" bg="#f5f5f7">
      <Sidebar />
      <Box flex={1} p={10} overflowY="auto" maxH="100vh">
        {children}
      </Box>
    </Flex>
  )
}

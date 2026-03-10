'use client'

import themeSystem from '@/Theme'
import { ChakraProvider } from '@chakra-ui/react'

import { ColorModeProvider, type ColorModeProviderProps } from './color-mode'

export function Provider(props: ColorModeProviderProps) {
  return (
    <ChakraProvider value={themeSystem}>
      <ColorModeProvider forcedTheme="light" {...props} />
    </ChakraProvider>
  )
}

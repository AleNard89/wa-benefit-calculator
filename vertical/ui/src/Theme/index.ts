import { createSystem, defaultConfig, defineConfig } from '@chakra-ui/react'

import { semanticTokens } from './semantic-tokens'
import { tokens } from './tokens'

const config = defineConfig({
  globalCss: {
    html: {
      colorPalette: 'brand',
    },
    body: {
      bg: '#f5f5f7',
      color: '#1d1d1f',
      fontFamily: '"Lato", -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif',
    },
  },
  theme: {
    tokens,
    semanticTokens,
  },
})

export default createSystem(defaultConfig, config)

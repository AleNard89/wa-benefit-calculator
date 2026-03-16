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
    '.chat-markdown': {
      '& p': { margin: '0 0 0.5em 0', '&:last-child': { marginBottom: 0 } },
      '& ul, & ol': { paddingLeft: '1.2em', margin: '0.3em 0' },
      '& li': { marginBottom: '0.15em' },
      '& strong': { fontWeight: 700 },
      '& table': {
        width: '100%',
        borderCollapse: 'collapse',
        margin: '0.5em 0',
        fontSize: '13px',
      },
      '& th, & td': {
        border: '1px solid #d2d2d7',
        padding: '6px 10px',
        textAlign: 'left',
      },
      '& th': {
        bg: '#e8e8ed',
        fontWeight: 700,
        fontSize: '12px',
        textTransform: 'uppercase',
        letterSpacing: '0.02em',
      },
      '& tr:nth-of-type(even)': { bg: '#fafafa' },
      '& code': {
        bg: '#e8e8ed',
        padding: '1px 4px',
        borderRadius: '4px',
        fontSize: '13px',
      },
      '& pre': {
        bg: '#1d1d1f',
        color: '#f5f5f7',
        padding: '12px',
        borderRadius: '8px',
        overflow: 'auto',
        margin: '0.5em 0',
        '& code': { bg: 'transparent', padding: 0, color: 'inherit' },
      },
    },
    '@keyframes blink': {
      '0%, 100%': { opacity: 1 },
      '50%': { opacity: 0 },
    },
  },
  theme: {
    tokens,
    semanticTokens,
  },
})

export default createSystem(defaultConfig, config)

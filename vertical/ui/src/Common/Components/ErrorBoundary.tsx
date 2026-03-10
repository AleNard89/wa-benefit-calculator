import { Component, type ErrorInfo, type ReactNode } from 'react'
import { Box, Text } from '@chakra-ui/react'

interface Props {
  children: ReactNode
}

interface State {
  hasError: boolean
  error: Error | null
}

export default class ErrorBoundary extends Component<Props, State> {
  constructor(props: Props) {
    super(props)
    this.state = { hasError: false, error: null }
  }

  static getDerivedStateFromError(error: Error): State {
    return { hasError: true, error }
  }

  componentDidCatch(error: Error, info: ErrorInfo) {
    console.error('[ErrorBoundary]', error, info.componentStack)
  }

  render() {
    if (this.state.hasError) {
      return (
        <Box p={8}>
          <Text fontWeight="700" color="red.500" mb={2}>Errore di rendering</Text>
          <Box as="pre" fontSize="sm" bg="gray.100" p={4} borderRadius="md" overflow="auto" whiteSpace="pre-wrap">
            {this.state.error?.message}
            {'\n\n'}
            {this.state.error?.stack}
          </Box>
        </Box>
      )
    }
    return this.props.children
  }
}

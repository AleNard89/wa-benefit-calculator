import '@fontsource/lato/400.css'
import '@fontsource/lato/700.css'
import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { Provider } from 'react-redux'

import App from './App.tsx'
import store from './Core/Redux/Store.ts'
import { Provider as ChakraProvider } from './Snippets/provider.tsx'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <Provider store={store}>
      <ChakraProvider>
        <App />
      </ChakraProvider>
    </Provider>
  </StrictMode>,
)

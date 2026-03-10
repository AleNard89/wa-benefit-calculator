import { useAuthentication } from './Auth/Hooks'
import AppRouter from './Core/Router'
import StartupView from './Core/Views/StartupView'
import { Toaster } from './Snippets/toaster'

const App: React.FC = () => {
  const { isComplete } = useAuthentication()

  if (!isComplete) return <StartupView />

  return (
    <>
      <Toaster />
      <AppRouter />
    </>
  )
}

export default App

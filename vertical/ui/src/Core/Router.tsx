import { BrowserRouter, Route, Routes, Navigate } from 'react-router-dom'

import Config from '@/Config'
import AdminRoute from '@/Auth/Guards/AdminRoute'
import PrivateRoute from '@/Auth/Guards/PrivateRoute'
import SignInView from '@/Auth/Views/SignInView'
import BaseLayout from '@/Common/Layouts/BaseLayout'
import ErrorBoundary from '@/Common/Components/ErrorBoundary'
import ProcessRouter from '@/Processes/Routers/Router'
import HomeView from './Views/HomeView'
import SettingsView from './Views/SettingsView'
import ChatView from './Views/ChatView'
import OrchestratorPage from '@/Orchestrator/OrchestratorPage'

function PrivateLayout({ children }: { children: React.ReactNode }) {
  return (
    <PrivateRoute>
      <BaseLayout>{children}</BaseLayout>
    </PrivateRoute>
  )
}

export default function AppRouter() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path={Config.urls.signIn} element={<SignInView />} />

        <Route
          path={Config.urls.home}
          element={<PrivateLayout><HomeView /></PrivateLayout>}
        />
        <Route
          path="/processes/*"
          element={<PrivateLayout><ErrorBoundary><ProcessRouter /></ErrorBoundary></PrivateLayout>}
        />
        <Route
          path="/settings"
          element={<PrivateLayout><SettingsView /></PrivateLayout>}
        />
        <Route
          path="/orchestrator"
          element={<PrivateLayout><AdminRoute><OrchestratorPage /></AdminRoute></PrivateLayout>}
        />
        <Route
          path="/chat"
          element={<PrivateLayout><ChatView /></PrivateLayout>}
        />

        <Route path="*" element={<Navigate to={Config.urls.home} replace />} />
      </Routes>
    </BrowserRouter>
  )
}

import { BrowserRouter, Route, Routes } from 'react-router-dom'
import { AuthProvider } from './auth/AuthContext'
import { ProtectedRoute } from './components/ProtectedRoute'
import { AdminRoute } from './components/AdminRoute'
import { AuthCallbackPage } from './pages/AuthCallbackPage'
import { AdminDashboard } from './pages/admin/AdminDashboard'
import { AdminUserDetail } from './pages/admin/AdminUserDetail'
import { GetStartedPage } from './pages/GetStartedPage'
import { LandingPage } from './pages/LandingPage'
import { ProfilePage } from './pages/ProfilePage'
import { TutorialPage } from './pages/TutorialPage'
import { WaitingListPage } from './pages/WaitingListPage'

function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <Routes>
          <Route path="/" element={<LandingPage />} />
          <Route path="/auth/callback" element={<AuthCallbackPage />} />
          <Route path="/get-started" element={<GetStartedPage />} />
          <Route path="/tutorial" element={<TutorialPage />} />
          <Route path="/waiting-list" element={<ProtectedRoute><WaitingListPage /></ProtectedRoute>} />
          <Route path="/profile" element={<ProtectedRoute><ProfilePage /></ProtectedRoute>} />
          <Route path="/admin" element={<AdminRoute><AdminDashboard /></AdminRoute>} />
          <Route path="/admin/users/:id" element={<AdminRoute><AdminUserDetail /></AdminRoute>} />
        </Routes>
      </AuthProvider>
    </BrowserRouter>
  )
}

export default App

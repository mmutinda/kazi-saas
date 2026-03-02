import { Routes, Route, Navigate } from 'react-router-dom'
import { HomePage } from '@/pages/HomePage'
import { JobDetailPage } from '@/pages/JobDetailPage'
import { LoginPage } from '@/pages/admin/LoginPage'
import { DashboardPage } from '@/pages/admin/DashboardPage'
import { JobsPage } from '@/pages/admin/JobsPage'

function RequireAuth({ children }: { children: React.ReactNode }) {
  const token = localStorage.getItem('kazi_token')
  if (!token) return <Navigate to="/admin/login" replace />
  return <>{children}</>
}

export default function App() {
  return (
    <Routes>
      {/* Public */}
      <Route path="/" element={<HomePage />} />
      <Route path="/jobs/:id" element={<JobDetailPage />} />

      {/* Admin */}
      <Route path="/admin/login" element={<LoginPage />} />
      <Route
        path="/admin/dashboard"
        element={<RequireAuth><DashboardPage /></RequireAuth>}
      />
      <Route
        path="/admin/jobs"
        element={<RequireAuth><JobsPage /></RequireAuth>}
      />

      {/* Fallbacks */}
      <Route path="/admin" element={<Navigate to="/admin/dashboard" replace />} />
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  )
}

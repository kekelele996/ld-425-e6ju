import { createBrowserRouter, Navigate } from 'react-router-dom'
import AppLayout from '@/components/layout/AppLayout'
import ErrorBoundary from '@/components/ErrorBoundary'
import { RequireAuth } from './guards'
import Login from '@/pages/Login'
import Dashboard from '@/pages/Dashboard'
import DesignManage from '@/pages/DesignManage'
import MaterialManage from '@/pages/MaterialManage'
import ConstructionProgress from '@/pages/ConstructionProgress'
import BudgetManage from '@/pages/BudgetManage'

export const router = createBrowserRouter([
  { path: '/login', element: <Login /> },
  {
    path: '/',
    element: (
      <ErrorBoundary>
        <RequireAuth>
          <AppLayout />
        </RequireAuth>
      </ErrorBoundary>
    ),
    children: [
      { index: true, element: <Navigate to="/dashboard" replace /> },
      { path: 'dashboard', element: <Dashboard /> },
      { path: 'designs', element: <DesignManage /> },
      { path: 'materials', element: <MaterialManage /> },
      { path: 'construction', element: <ConstructionProgress /> },
      { path: 'budget', element: <BudgetManage /> },
    ],
  },
])

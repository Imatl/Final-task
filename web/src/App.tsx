import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { type ReactNode } from 'react';
import '@/i18n';
import { MainLayout } from '@/components/layout';
import { AgentDashboardPage } from '@/pages/AgentDashboard';
import { AnalyticsPage } from '@/pages/Analytics';
import { IntegrationsPage } from '@/pages/Integrations';
import { SettingsPage } from '@/pages/Settings';
import { AuthPage } from '@/pages/Auth';
import { RegisterPage } from '@/pages/Register';
import { StaffPanelPage } from '@/pages/StaffPanel';
import { CompaniesPage } from '@/pages/Companies';
import { useAuthStore, type UserLevel } from '@/store/auth';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 30_000,
      retry: 1,
      refetchOnWindowFocus: false,
    },
  },
});

function ProtectedRoute({ children, minLevel = 1 }: { children: ReactNode; minLevel?: UserLevel }) {
  const { user } = useAuthStore();
  if (!user) return <Navigate to="/login" replace />;
  if (user.level < minLevel) return <Navigate to="/dashboard" replace />;
  return <>{children}</>;
}

export default function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <Routes>
          <Route
            element={
              <ProtectedRoute>
                <MainLayout />
              </ProtectedRoute>
            }
          >
            <Route path="dashboard" element={<AgentDashboardPage />} />
            <Route
              path="analytics"
              element={
                <ProtectedRoute minLevel={2}>
                  <AnalyticsPage />
                </ProtectedRoute>
              }
            />
            <Route
              path="staff"
              element={
                <ProtectedRoute minLevel={3}>
                  <StaffPanelPage />
                </ProtectedRoute>
              }
            />
            <Route
              path="companies"
              element={
                <ProtectedRoute minLevel={5}>
                  <CompaniesPage />
                </ProtectedRoute>
              }
            />
            <Route
              path="integrations"
              element={
                <ProtectedRoute minLevel={4}>
                  <IntegrationsPage />
                </ProtectedRoute>
              }
            />
            <Route
              path="settings"
              element={
                <ProtectedRoute minLevel={4}>
                  <SettingsPage />
                </ProtectedRoute>
              }
            />
          </Route>

          <Route path="login" element={<AuthPage />} />
          <Route path="register/:token" element={<RegisterPage />} />
          <Route index element={<Navigate to="/dashboard" replace />} />
          <Route path="*" element={<Navigate to="/dashboard" replace />} />
        </Routes>
      </BrowserRouter>
    </QueryClientProvider>
  );
}

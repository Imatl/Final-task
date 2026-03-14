import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import '@/i18n';
import { MainLayout, ChatLayout } from '@/components/layout';
import { CustomerChatPage } from '@/pages/CustomerChat';
import { AgentDashboardPage } from '@/pages/AgentDashboard';
import { AnalyticsPage } from '@/pages/Analytics';
import { SettingsPage } from '@/pages/Settings';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 30_000,
      retry: 1,
      refetchOnWindowFocus: false,
    },
  },
});

export default function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <Routes>
          <Route element={<ChatLayout />}>
            <Route path="chat" element={<CustomerChatPage />} />
          </Route>

          <Route element={<MainLayout />}>
            <Route path="agent/dashboard" element={<AgentDashboardPage />} />
            <Route path="agent/analytics" element={<AnalyticsPage />} />
            <Route path="agent/settings" element={<SettingsPage />} />
          </Route>

          <Route index element={<Navigate to="/chat" replace />} />
          <Route path="*" element={<Navigate to="/chat" replace />} />
        </Routes>
      </BrowserRouter>
    </QueryClientProvider>
  );
}

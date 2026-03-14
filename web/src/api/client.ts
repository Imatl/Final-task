import axios from 'axios';

const api = axios.create({
  baseURL: '/api',
  headers: { 'Content-Type': 'application/json' },
});

api.interceptors.request.use((config) => {
  try {
    const raw = localStorage.getItem('sf-auth-user');
    if (raw) {
      const user = JSON.parse(raw);
      if (user?.company && config.params !== false) {
        config.params = { ...config.params, company: user.company };
      }
    }
  } catch {
    // ignore
  }
  return config;
});

export interface ChatRequest {
  ticket_id?: string;
  customer_id: string;
  channel?: string;
  message: string;
}

export interface ChatResponse {
  ticket_id: string;
  message: string;
  analysis?: AIAnalysis;
  actions?: Action[];
  auto_fixed: boolean;
}

export interface Ticket {
  id: string;
  customer_id: string;
  subject: string;
  channel: string;
  status: string;
  priority: string;
  category: string;
  agent_id?: string;
  ai_summary?: string;
  created_at: string;
  updated_at: string;
  closed_at?: string;
}

export interface Message {
  id: string;
  ticket_id: string;
  role: string;
  content: string;
  created_at: string;
}

export interface Action {
  id: string;
  ticket_id: string;
  type: string;
  params: string;
  status: string;
  result?: string;
  confidence: number;
  created_at: string;
  executed_at?: string;
}

export interface Customer {
  id: string;
  name: string;
  email: string;
  plan: string;
  created_at: string;
}

export interface AIAnalysis {
  ticket_id: string;
  intent: string;
  sentiment: string;
  urgency: string;
  suggested_tools: string[];
  reasoning: string;
  confidence: number;
}

export interface TicketDetail {
  ticket: Ticket;
  customer: Customer;
  messages: Message[];
  actions: Action[];
  analysis?: AIAnalysis;
}

export interface AnalyticsOverview {
  total_tickets: number;
  open_tickets: number;
  avg_resolve_time_minutes: number;
  auto_resolve_rate: number;
  by_category: Record<string, number>;
  by_priority: Record<string, number>;
  by_sentiment: Record<string, number>;
}

export interface AgentPerformance {
  agent_id: string;
  agent_name: string;
  tickets_resolved: number;
  avg_resolve_time_minutes: number;
  quality_score: number;
}

export interface LLMMetrics {
  provider: string;
  model: string;
  latency_ms: number;
  input_tokens: number;
  output_tokens: number;
  tool_calls: number;
  timestamp: string;
  error?: string;
}

export interface ProvidersResponse {
  providers: string[];
  active: string;
}

export interface MetricsResponse {
  metrics: LLMMetrics[];
  stats: Record<string, unknown>;
}

export const chatApi = {
  send: (data: ChatRequest) => api.post<ChatResponse>('/chat', data),
};

export const ticketsApi = {
  list: (params?: Record<string, string>) => api.get<{ tickets: Ticket[]; total: number }>('/tickets', { params }),
  get: (id: string) => api.get<TicketDetail>(`/tickets/${id}`),
  updateStatus: (id: string, status: string) => api.put(`/tickets/${id}/status`, { status }),
  assign: (id: string, agentId: string) => api.put(`/tickets/${id}/assign`, { agent_id: agentId }),
  approveAction: (actionId: string, approved: boolean, agentId: string) =>
    api.post('/tickets/actions/approve', { action_id: actionId, approved, agent_id: agentId }),
  reply: (id: string, message: string, agentId?: string) =>
    api.post(`/tickets/${id}/reply`, { message, agent_id: agentId || 'agent-1' }),
  suggest: (id: string) =>
    api.post<{ suggestion: string }>(`/tickets/${id}/suggest`),
};

export const analyticsApi = {
  overview: () => api.get<AnalyticsOverview>('/analytics/overview'),
  agents: () => api.get<AgentPerformance[]>('/analytics/agents'),
};

export const settingsApi = {
  getProviders: () => api.get<ProvidersResponse>('/settings/providers'),
  setProvider: (provider: string) => api.put('/settings/providers', { provider }),
  getMetrics: (limit?: number) => api.get<MetricsResponse>('/settings/metrics', { params: { limit } }),
};

export interface Company {
  id: string;
  name: string;
  staff_count: number;
  ai_spend_usd: number;
  created_at: string;
}

export const companiesApi = {
  list: () => api.get<{ companies: Company[]; total: number }>('/companies'),
};

export interface InviteValidation {
  valid: boolean;
}

export interface InviteGenerated {
  token: string;
  link: string;
}

export interface RegisterRequest {
  token: string;
  name: string;
  company: string;
  email: string;
  password: string;
}

export interface RegisterResponse {
  id: string;
  email: string;
  name: string;
  level: number;
  role: string;
  company?: string;
}

export const inviteApi = {
  validate: (token: string) => api.get<InviteValidation>(`/invite/${token}`),
  generate: (createdBy: string) =>
    api.post<InviteGenerated>('/auth/invite', { created_by: createdBy }),
  register: (data: RegisterRequest) =>
    api.post<RegisterResponse>('/auth/register', data),
};

export interface IntegrationInfo {
  id: string;
  name: string;
  type: string;
  status: string;
}

export const integrationsApi = {
  list: () => api.get<IntegrationInfo[]>('/integrations'),
  connect: (type: string, config: Record<string, string>) =>
    api.post<{ status: string; id: string }>('/integrations/connect', { type, config }),
  disconnect: (id: string) =>
    api.post('/integrations/disconnect', { id }),
};

export interface KBEntry {
  id: string;
  company: string;
  question: string;
  answer: string;
  created_at: string;
}

export const knowledgeApi = {
  list: () => api.get<KBEntry[]>('/knowledge'),
  create: (company: string, question: string, answer: string) =>
    api.post<KBEntry>('/knowledge', { company, question, answer }),
  remove: (id: string) => api.delete(`/knowledge/${id}`),
};

export default api;

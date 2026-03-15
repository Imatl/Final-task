import { useQuery } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import { Users, TicketCheck, Clock, Star, TrendingUp, Zap } from 'lucide-react';
import { analyticsApi, type AgentPerformance } from '@/api/client';

function AgentRow({ agent }: { agent: AgentPerformance }) {
  const { t } = useTranslation();
  const quality = Math.round(agent.quality_score * 100);

  return (
    <div className="flex items-center gap-4 px-5 py-3.5 hover:bg-cosmic-800/40 transition-colors">
      <div className="w-9 h-9 rounded-full bg-gradient-to-br from-velvet-600 to-neon-violet flex items-center justify-center flex-shrink-0">
        <span className="text-sm font-bold text-white">
          {agent.agent_name
            .split(' ')
            .map((n) => n[0])
            .join('')
            .slice(0, 2)}
        </span>
      </div>

      <div className="flex-1 min-w-0">
        <p className="text-sm font-medium text-white truncate">{agent.agent_name}</p>
        <p className="text-xs text-gray-500 truncate">{agent.tickets_resolved} {t('staffPanel.resolved')}</p>
      </div>

      <div className="hidden sm:flex items-center gap-6 flex-shrink-0 text-xs text-gray-400">
        <div className="flex items-center gap-1.5 w-24 justify-end">
          <TicketCheck className="w-3.5 h-3.5 text-neon-violet" />
          <span>
            {agent.tickets_resolved} {t('staffPanel.resolved')}
          </span>
        </div>
        <div className="flex items-center gap-1.5 w-24 justify-end">
          <Clock className="w-3.5 h-3.5 text-neon-cyan" />
          <span>
            {Math.round(agent.avg_resolve_time_minutes)} {t('analytics.minutes')}
          </span>
        </div>
        <div className="flex items-center gap-1.5 w-20 justify-end">
          <Star className="w-3.5 h-3.5 text-amber-400" />
          <span>{quality}%</span>
        </div>
      </div>
    </div>
  );
}

export function StaffPanelPage() {
  const { t } = useTranslation();

  const { data: rawAgents, isLoading: loadingAgents } = useQuery({
    queryKey: ['analytics-agents'],
    queryFn: () => analyticsApi.agents().then((r) => r.data),
    refetchInterval: 30_000,
  });

  const agents = rawAgents ?? [];

  const { data: overview, isLoading: loadingOverview } = useQuery({
    queryKey: ['analytics-overview'],
    queryFn: () => analyticsApi.overview().then((r) => r.data),
    refetchInterval: 30_000,
  });

  const totalResolved = agents.reduce((sum, a) => sum + a.tickets_resolved, 0);
  const avgQuality =
    agents.length > 0
      ? Math.round((agents.reduce((sum, a) => sum + a.quality_score, 0) / agents.length) * 100)
      : 0;

  return (
    <div className="p-6 space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-white">{t('staffPanel.title')}</h1>
        <p className="text-sm text-gray-400 mt-1">{t('staffPanel.subtitle')}</p>
      </div>

      <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
        <div className="bg-cosmic-900/80 border border-cosmic-700/50 rounded-xl p-6">
          <div className="w-11 h-11 rounded-xl bg-neon-violet/10 border border-neon-violet/20 flex items-center justify-center mb-4">
            <Users className="w-6 h-6 text-neon-violet" />
          </div>
          <p className="text-4xl font-bold text-white mb-1">
            {loadingAgents ? '—' : agents.length}
          </p>
          <span className="text-sm text-gray-400">{t('staffPanel.totalAgents')}</span>
        </div>

        <div className="bg-cosmic-900/80 border border-cosmic-700/50 rounded-xl p-6">
          <div className="w-11 h-11 rounded-xl bg-neon-cyan/10 border border-neon-cyan/20 flex items-center justify-center mb-4">
            <TicketCheck className="w-6 h-6 text-neon-cyan" />
          </div>
          <p className="text-4xl font-bold text-white mb-1">
            {loadingAgents ? '—' : totalResolved}
          </p>
          <span className="text-sm text-gray-400">{t('staffPanel.totalResolved')}</span>
        </div>

        <div className="bg-cosmic-900/80 border border-cosmic-700/50 rounded-xl p-6">
          <div className="w-11 h-11 rounded-xl bg-amber-500/10 border border-amber-500/20 flex items-center justify-center mb-4">
            <Star className="w-6 h-6 text-amber-400" />
          </div>
          <p className="text-4xl font-bold text-white mb-1">
            {loadingAgents ? '—' : `${avgQuality}%`}
          </p>
          <span className="text-sm text-gray-400">{t('staffPanel.avgQuality')}</span>
        </div>

        <div className="bg-cosmic-900/80 border border-cosmic-700/50 rounded-xl p-6">
          <div className="w-11 h-11 rounded-xl bg-neon-green/10 border border-neon-green/20 flex items-center justify-center mb-4">
            <Zap className="w-6 h-6 text-neon-green" />
          </div>
          <p className="text-4xl font-bold text-white mb-1">
            {loadingOverview || !overview
              ? '—'
              : `${Math.round(overview.auto_resolve_rate * 100)}%`}
          </p>
          <span className="text-sm text-gray-400">{t('analytics.autoResolve')}</span>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
        <div className="bg-cosmic-900/80 border border-cosmic-700/50 rounded-xl p-4">
          <div className="flex items-center gap-3 mb-3">
            <TrendingUp className="w-4 h-4 text-neon-violet" />
            <span className="text-sm font-semibold text-white">{t('analytics.totalTickets')}</span>
          </div>
          <p className="text-3xl font-bold text-white">
            {loadingOverview || !overview ? '—' : overview.total_tickets}
          </p>
          <p className="text-xs text-gray-500 mt-1">
            {!overview ? '' : `${overview.open_tickets} ${t('analytics.openTickets').toLowerCase()}`}
          </p>
        </div>

        <div className="bg-cosmic-900/80 border border-cosmic-700/50 rounded-xl p-4">
          <div className="flex items-center gap-3 mb-3">
            <Clock className="w-4 h-4 text-neon-cyan" />
            <span className="text-sm font-semibold text-white">{t('analytics.avgResolve')}</span>
          </div>
          <p className="text-3xl font-bold text-white">
            {loadingOverview || !overview
              ? '—'
              : `${Math.round(overview.avg_resolve_time_minutes)} ${t('analytics.minutes')}`}
          </p>
        </div>
      </div>

      <div className="bg-cosmic-900/80 border border-cosmic-700/50 rounded-xl overflow-hidden">
        <div className="flex items-center justify-between px-5 py-4 border-b border-cosmic-700/50">
          <h2 className="text-sm font-semibold text-white">{t('staffPanel.staffList')}</h2>
          <span className="text-xs text-gray-400">
            {agents.length} {t('staffPanel.agents')}
          </span>
        </div>

        {loadingAgents ? (
          <div className="px-5 py-8 text-center text-sm text-gray-500">{t('dashboard.loading')}</div>
        ) : agents.length === 0 ? (
          <div className="px-5 py-8 text-center text-sm text-gray-500">{t('analytics.noData')}</div>
        ) : (
          <div className="divide-y divide-cosmic-700/30">
            <div className="hidden sm:flex items-center gap-4 px-5 py-2 text-[10px] font-semibold uppercase tracking-wider text-gray-500">
              <div className="w-9 flex-shrink-0" />
              <div className="flex-1">{t('analytics.agent')}</div>
              <div className="flex gap-6 flex-shrink-0 text-right">
                <span className="w-24 text-right">{t('analytics.resolved')}</span>
                <span className="w-24 text-right">{t('analytics.avgTime')}</span>
                <span className="w-20 text-right">{t('analytics.quality')}</span>
              </div>
            </div>
            {agents.map((agent) => (
              <AgentRow key={agent.agent_id} agent={agent} />
            ))}
          </div>
        )}
      </div>
    </div>
  );
}

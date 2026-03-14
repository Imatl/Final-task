import { useQuery } from '@tanstack/react-query';
import { analyticsApi } from '@/api/client';
import { useTranslation } from 'react-i18next';
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  PieChart,
  Pie,
  Cell,
} from 'recharts';
import { TrendingUp, Clock, Zap, Users } from 'lucide-react';
import { Card } from '@/components/ui';

const COLORS = ['#a855f7', '#e879f9', '#22d3ee', '#4ade80', '#60a5fa', '#f59e0b'];

function StatCard({ title, value, subtitle, icon: Icon, color }: { title: string; value: string | number; subtitle?: string; icon: typeof TrendingUp; color: string }) {
  return (
    <Card variant="metric" className="p-5">
      <div className="flex items-center justify-between mb-3">
        <span className="text-sm text-gray-400">{title}</span>
        <div className="w-10 h-10 rounded-lg flex items-center justify-center" style={{ backgroundColor: `${color}20` }}>
          <Icon className="w-5 h-5" style={{ color }} />
        </div>
      </div>
      <div className="text-3xl font-bold text-white">{value}</div>
      {subtitle && <div className="text-xs text-gray-500 mt-1">{subtitle}</div>}
    </Card>
  );
}

const darkTooltipStyle = {
  contentStyle: { backgroundColor: '#110d24', border: '1px solid #1a1435', borderRadius: '8px', color: '#e5e7eb' },
  itemStyle: { color: '#c084fc' },
};

export function AnalyticsPage() {
  const { t } = useTranslation();

  const { data: overview } = useQuery({
    queryKey: ['analytics-overview'],
    queryFn: () => analyticsApi.overview().then((r) => r.data),
    refetchInterval: 10000,
  });

  const { data: agents } = useQuery({
    queryKey: ['analytics-agents'],
    queryFn: () => analyticsApi.agents().then((r) => r.data),
    refetchInterval: 10000,
  });

  const categoryData = overview?.by_category
    ? Object.entries(overview.by_category).map(([name, value]) => ({ name, value }))
    : [];

  const priorityData = overview?.by_priority
    ? Object.entries(overview.by_priority).map(([name, value]) => ({ name, value }))
    : [];

  const sentimentData = overview?.by_sentiment
    ? Object.entries(overview.by_sentiment).map(([name, value]) => ({ name, value }))
    : [];

  return (
    <div className="p-6 space-y-6">
      <h1 className="text-lg font-semibold text-white">{t('analytics.title')}</h1>

      <div className="grid grid-cols-4 gap-4">
        <StatCard title={t('analytics.totalTickets')} value={overview?.total_tickets || 0} icon={TrendingUp} color="#a855f7" />
        <StatCard title={t('analytics.openTickets')} value={overview?.open_tickets || 0} icon={Clock} color="#22d3ee" />
        <StatCard
          title={t('analytics.avgResolve')}
          value={overview ? `${overview.avg_resolve_time_minutes.toFixed(1)}${t('analytics.minutes')}` : `0${t('analytics.minutes')}`}
          icon={Zap}
          color="#e879f9"
        />
        <StatCard
          title={t('analytics.autoResolve')}
          value={overview ? `${(overview.auto_resolve_rate * 100).toFixed(0)}%` : '0%'}
          icon={Users}
          color="#4ade80"
        />
      </div>

      <div className="grid grid-cols-3 gap-4">
        <Card className="p-5">
          <h3 className="text-sm font-medium text-gray-300 mb-4">{t('analytics.byCategory')}</h3>
          <ResponsiveContainer width="100%" height={250}>
            <BarChart data={categoryData}>
              <CartesianGrid strokeDasharray="3 3" stroke="#1a1435" />
              <XAxis dataKey="name" tick={{ fontSize: 11, fill: '#9ca3af' }} />
              <YAxis tick={{ fill: '#9ca3af' }} />
              <Tooltip {...darkTooltipStyle} />
              <Bar dataKey="value" fill="#a855f7" radius={[4, 4, 0, 0]} />
            </BarChart>
          </ResponsiveContainer>
        </Card>

        <Card className="p-5">
          <h3 className="text-sm font-medium text-gray-300 mb-4">{t('analytics.byPriority')}</h3>
          <ResponsiveContainer width="100%" height={250}>
            <PieChart>
              <Pie data={priorityData} dataKey="value" nameKey="name" cx="50%" cy="50%" outerRadius={80} label={{ fill: '#9ca3af', fontSize: 11 }}>
                {priorityData.map((_, i) => (
                  <Cell key={i} fill={COLORS[i % COLORS.length]} />
                ))}
              </Pie>
              <Tooltip {...darkTooltipStyle} />
            </PieChart>
          </ResponsiveContainer>
        </Card>

        <Card className="p-5">
          <h3 className="text-sm font-medium text-gray-300 mb-4">{t('analytics.bySentiment')}</h3>
          <ResponsiveContainer width="100%" height={250}>
            <PieChart>
              <Pie data={sentimentData} dataKey="value" nameKey="name" cx="50%" cy="50%" outerRadius={80} label={{ fill: '#9ca3af', fontSize: 11 }}>
                {sentimentData.map((_, i) => (
                  <Cell key={i} fill={COLORS[i % COLORS.length]} />
                ))}
              </Pie>
              <Tooltip {...darkTooltipStyle} />
            </PieChart>
          </ResponsiveContainer>
        </Card>
      </div>

      {agents && agents.length > 0 && (
        <Card className="p-5">
          <h3 className="text-sm font-medium text-gray-300 mb-4">{t('analytics.agentPerformance')}</h3>
          <table className="w-full text-sm">
            <thead>
              <tr className="text-left text-gray-500 border-b border-cosmic-700/50">
                <th className="pb-3">{t('analytics.agent')}</th>
                <th className="pb-3">{t('analytics.resolved')}</th>
                <th className="pb-3">{t('analytics.avgTime')}</th>
                <th className="pb-3">{t('analytics.quality')}</th>
              </tr>
            </thead>
            <tbody>
              {agents.map((a) => (
                <tr key={a.agent_id} className="border-b border-cosmic-700/30 last:border-0">
                  <td className="py-3 font-medium text-gray-200">{a.agent_name}</td>
                  <td className="py-3 text-gray-400">{a.tickets_resolved}</td>
                  <td className="py-3 text-gray-400">{a.avg_resolve_time_minutes.toFixed(1)}{t('analytics.minutes')}</td>
                  <td className="py-3">
                    <div className="flex items-center gap-2">
                      <div className="w-24 bg-cosmic-700 rounded-full h-2">
                        <div className="bg-gradient-to-r from-velvet-600 to-neon-violet h-2 rounded-full transition-all" style={{ width: `${a.quality_score * 100}%` }} />
                      </div>
                      <span className="text-xs text-gray-500">{(a.quality_score * 100).toFixed(0)}%</span>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </Card>
      )}
    </div>
  );
}

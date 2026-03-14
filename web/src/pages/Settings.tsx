import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { settingsApi } from '@/api/client';
import type { LLMMetrics } from '@/api/client';
import { useTranslation } from 'react-i18next';
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  LineChart,
  Line,
} from 'recharts';
import { cn } from '@/lib/cn';
import { Activity, Cpu, Zap } from 'lucide-react';
import { Card, Badge } from '@/components/ui';

function LatencyChart({ metrics }: { metrics: LLMMetrics[] }) {
  const data = metrics.map((m, i) => ({
    idx: i,
    latency: m.latency_ms,
    provider: m.provider,
    tokens: m.input_tokens + m.output_tokens,
  }));

  return (
    <ResponsiveContainer width="100%" height={250}>
      <LineChart data={data}>
        <CartesianGrid strokeDasharray="3 3" stroke="#1a1435" />
        <XAxis dataKey="idx" tick={{ fontSize: 11, fill: '#9ca3af' }} />
        <YAxis tick={{ fontSize: 11, fill: '#9ca3af' }} />
        <Tooltip
          content={({ active, payload }) => {
            if (!active || !payload?.length) return null;
            const d = payload[0].payload;
            return (
              <div className="bg-cosmic-800 border border-cosmic-700/50 rounded-lg p-3 text-xs shadow-xl">
                <div className="text-gray-400">Provider: <span className="text-neon-purple">{d.provider}</span></div>
                <div className="text-gray-400">Latency: <span className="text-white font-mono">{d.latency}ms</span></div>
                <div className="text-gray-400">Tokens: <span className="text-white font-mono">{d.tokens}</span></div>
              </div>
            );
          }}
        />
        <Line type="monotone" dataKey="latency" stroke="#a855f7" dot={{ r: 3, fill: '#a855f7' }} strokeWidth={2} />
      </LineChart>
    </ResponsiveContainer>
  );
}

function TokensChart({ metrics }: { metrics: LLMMetrics[] }) {
  const { t } = useTranslation();
  const data = metrics.map((m, i) => ({
    idx: i,
    input: m.input_tokens,
    output: m.output_tokens,
  }));

  return (
    <ResponsiveContainer width="100%" height={250}>
      <BarChart data={data}>
        <CartesianGrid strokeDasharray="3 3" stroke="#1a1435" />
        <XAxis dataKey="idx" tick={{ fontSize: 11, fill: '#9ca3af' }} />
        <YAxis tick={{ fontSize: 11, fill: '#9ca3af' }} />
        <Tooltip
          contentStyle={{ backgroundColor: '#110d24', border: '1px solid #1a1435', borderRadius: '8px', color: '#e5e7eb', fontSize: 12 }}
        />
        <Bar dataKey="input" fill="#60a5fa" stackId="tokens" name={t('settings.input')} />
        <Bar dataKey="output" fill="#4ade80" stackId="tokens" name={t('settings.output')} radius={[4, 4, 0, 0]} />
      </BarChart>
    </ResponsiveContainer>
  );
}

export function SettingsPage() {
  const { t } = useTranslation();
  const queryClient = useQueryClient();

  const { data: providersData } = useQuery({
    queryKey: ['providers'],
    queryFn: () => settingsApi.getProviders().then((r) => r.data),
  });

  const { data: metricsData } = useQuery({
    queryKey: ['metrics'],
    queryFn: () => settingsApi.getMetrics(50).then((r) => r.data),
    refetchInterval: 5000,
  });

  const switchProvider = useMutation({
    mutationFn: (provider: string) => settingsApi.setProvider(provider),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['providers'] }),
  });

  const stats = metricsData?.stats as Record<string, unknown> | undefined;
  const byProvider = stats?.by_provider as Record<string, { calls: number; avg_ms: number; min_ms: number; max_ms: number }> | undefined;

  return (
    <div className="p-6 space-y-6">
      <h1 className="text-lg font-semibold text-white">{t('settings.title')}</h1>

      <Card className="p-5">
        <h3 className="text-sm font-medium text-gray-300 mb-4 flex items-center gap-2">
          <Cpu className="w-4 h-4 text-neon-violet" />
          {t('settings.provider')}
        </h3>
        <div className="flex gap-3">
          {providersData?.providers?.map((p) => (
            <button
              key={p}
              onClick={() => switchProvider.mutate(p)}
              className={cn(
                'px-6 py-3 rounded-xl border-2 text-sm font-medium transition-all duration-200',
                p === providersData.active
                  ? 'border-neon-violet bg-velvet-900/50 text-neon-purple neon-glow-sm'
                  : 'border-cosmic-600 text-gray-400 hover:border-cosmic-500 hover:text-gray-200'
              )}
            >
              {p === 'anthropic' ? 'Claude (Anthropic)' : p === 'openai' ? 'GPT (OpenAI)' : p}
              {p === providersData.active && <Badge variant="neon" className="ml-2">{t('settings.active')}</Badge>}
            </button>
          ))}
        </div>
      </Card>

      {byProvider && Object.keys(byProvider).length > 0 && (
        <div className="grid grid-cols-2 gap-4">
          {Object.entries(byProvider).map(([name, s]) => (
            <Card key={name} variant="metric" className="p-5">
              <h4 className="text-sm font-medium text-gray-300 mb-3 capitalize">{name}</h4>
              <div className="grid grid-cols-2 gap-3 text-sm">
                <div>
                  <span className="text-gray-500 text-xs">{t('settings.calls')}</span>
                  <div className="text-2xl font-bold text-white">{s.calls}</div>
                </div>
                <div>
                  <span className="text-gray-500 text-xs">{t('settings.avgLatency')}</span>
                  <div className="text-2xl font-bold text-neon-violet font-mono">{s.avg_ms}ms</div>
                </div>
                <div>
                  <span className="text-gray-500 text-xs">{t('settings.min')}</span>
                  <div className="font-medium text-gray-300 font-mono">{s.min_ms}ms</div>
                </div>
                <div>
                  <span className="text-gray-500 text-xs">{t('settings.max')}</span>
                  <div className="font-medium text-gray-300 font-mono">{s.max_ms}ms</div>
                </div>
              </div>
            </Card>
          ))}
        </div>
      )}

      <div className="grid grid-cols-2 gap-4">
        <Card className="p-5">
          <h3 className="text-sm font-medium text-gray-300 mb-4 flex items-center gap-2">
            <Activity className="w-4 h-4 text-neon-cyan" />
            {t('settings.latency')}
          </h3>
          {metricsData?.metrics?.length ? (
            <LatencyChart metrics={metricsData.metrics} />
          ) : (
            <div className="h-[250px] flex items-center justify-center text-gray-600 text-sm">{t('settings.noData')}</div>
          )}
        </Card>

        <Card className="p-5">
          <h3 className="text-sm font-medium text-gray-300 mb-4 flex items-center gap-2">
            <Zap className="w-4 h-4 text-neon-green" />
            {t('settings.tokenUsage')}
          </h3>
          {metricsData?.metrics?.length ? (
            <TokensChart metrics={metricsData.metrics} />
          ) : (
            <div className="h-[250px] flex items-center justify-center text-gray-600 text-sm">{t('settings.noData')}</div>
          )}
        </Card>
      </div>

      {metricsData?.metrics && metricsData.metrics.length > 0 && (
        <Card className="p-5">
          <h3 className="text-sm font-medium text-gray-300 mb-4">{t('settings.recentCalls')}</h3>
          <div className="overflow-x-auto">
            <table className="w-full text-xs">
              <thead>
                <tr className="text-left text-gray-500 border-b border-cosmic-700/50">
                  <th className="pb-2 pr-4">{t('settings.provider')}</th>
                  <th className="pb-2 pr-4">{t('settings.model')}</th>
                  <th className="pb-2 pr-4">{t('settings.latency')}</th>
                  <th className="pb-2 pr-4">{t('settings.input')}</th>
                  <th className="pb-2 pr-4">{t('settings.output')}</th>
                  <th className="pb-2 pr-4">{t('settings.tools')}</th>
                  <th className="pb-2">{t('settings.errorCol')}</th>
                </tr>
              </thead>
              <tbody>
                {[...metricsData.metrics].reverse().slice(0, 20).map((m, i) => (
                  <tr key={i} className="border-b border-cosmic-700/20 last:border-0">
                    <td className="py-2 pr-4 capitalize text-gray-300">{m.provider}</td>
                    <td className="py-2 pr-4 text-gray-500 font-mono">{m.model}</td>
                    <td className="py-2 pr-4 font-mono text-neon-violet">{m.latency_ms}ms</td>
                    <td className="py-2 pr-4 text-gray-400">{m.input_tokens}</td>
                    <td className="py-2 pr-4 text-gray-400">{m.output_tokens}</td>
                    <td className="py-2 pr-4 text-gray-400">{m.tool_calls}</td>
                    <td className="py-2 text-red-400">{m.error || '-'}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </Card>
      )}
    </div>
  );
}

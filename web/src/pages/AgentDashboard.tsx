import { useQuery } from '@tanstack/react-query';
import { ticketsApi } from '@/api/client';
import type { Ticket } from '@/api/client';
import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import {
  Clock,
  AlertCircle,
  CheckCircle,
  User,
  ChevronRight,
  Globe,
  MessageSquare,
  Mail,
} from 'lucide-react';
import { cn } from '@/lib/cn';
import { Card, Badge } from '@/components/ui';

const PRIORITY_BADGE: Record<string, 'error' | 'warning' | 'success' | 'default'> = {
  critical: 'error',
  high: 'warning',
  medium: 'default',
  low: 'success',
};

const CHANNEL_ICONS: Record<string, typeof Globe> = {
  web: Globe,
  telegram: MessageSquare,
  email: Mail,
};

const STATUS_ICONS: Record<string, typeof Clock> = {
  open: AlertCircle,
  in_progress: Clock,
  waiting: Clock,
  resolved: CheckCircle,
  closed: CheckCircle,
};

function TicketRow({ ticket, isSelected, onSelect }: { ticket: Ticket; isSelected: boolean; onSelect: (id: string) => void }) {
  const { t } = useTranslation();
  const StatusIcon = STATUS_ICONS[ticket.status] || Clock;
  const ChannelIcon = CHANNEL_ICONS[ticket.channel] || Globe;
  const age = Math.round((Date.now() - new Date(ticket.created_at).getTime()) / 60000);
  const ageLabel = age < 60 ? `${age}m` : `${Math.round(age / 60)}h`;

  return (
    <div
      onClick={() => onSelect(ticket.id)}
      className={cn(
        'flex items-center gap-4 px-4 py-3 border-b border-cosmic-700/30 cursor-pointer transition-all duration-200',
        isSelected
          ? 'bg-velvet-900/30 border-l-2 border-l-neon-violet'
          : 'hover:bg-cosmic-800/50'
      )}
    >
      <StatusIcon className={cn('w-5 h-5', ticket.status === 'open' ? 'text-red-400' : ticket.status === 'resolved' ? 'text-neon-green' : 'text-amber-400')} />
      <ChannelIcon className="w-4 h-4 text-gray-500" />
      <div className="flex-1 min-w-0">
        <div className="text-sm font-medium text-gray-200 truncate">{ticket.subject}</div>
        <div className="text-xs text-gray-500">{ticket.category} | {ageLabel} {t('dashboard.ago')}</div>
      </div>
      <Badge variant={PRIORITY_BADGE[ticket.priority] || 'default'}>{ticket.priority}</Badge>
      <Badge variant={ticket.status === 'open' ? 'error' : ticket.status === 'resolved' ? 'success' : 'default'} dot>{ticket.status}</Badge>
      <ChevronRight className="w-4 h-4 text-gray-600" />
    </div>
  );
}

function TicketDetailPanel({ ticketId, onClose }: { ticketId: string; onClose: () => void }) {
  const { t } = useTranslation();
  const { data, isLoading } = useQuery({
    queryKey: ['ticket', ticketId],
    queryFn: () => ticketsApi.get(ticketId).then((r) => r.data),
  });

  if (isLoading) return <div className="p-6 text-gray-500">Loading...</div>;
  if (!data) return null;

  const ChannelIcon = CHANNEL_ICONS[data.ticket.channel] || Globe;

  return (
    <div className="bg-cosmic-900/80 backdrop-blur-sm border-l border-cosmic-700/50 h-full overflow-y-auto">
      <div className="p-4 border-b border-cosmic-700/50 flex items-center justify-between">
        <div>
          <h3 className="font-semibold text-white">{data.ticket.subject}</h3>
          <div className="flex items-center gap-2 mt-1">
            <span className="text-xs text-gray-500 font-mono">{data.ticket.id.slice(0, 8)}</span>
            <Badge variant="neon"><ChannelIcon className="w-3 h-3 mr-1 inline" />{data.ticket.channel}</Badge>
          </div>
        </div>
        <button onClick={onClose} className="text-gray-500 hover:text-gray-300 text-xl transition-colors">&times;</button>
      </div>

      <div className="p-4 border-b border-cosmic-700/50">
        <div className="flex items-center gap-2 mb-2">
          <User className="w-4 h-4 text-gray-500" />
          <span className="text-sm font-medium text-gray-200">{data.customer.name}</span>
          <span className="text-xs text-gray-500">{data.customer.email}</span>
          <Badge variant={PRIORITY_BADGE[data.ticket.priority] || 'default'} className="ml-auto">{data.ticket.priority}</Badge>
        </div>
        <div className="text-xs text-gray-500">Plan: <span className="text-neon-purple">{data.customer.plan}</span> | Status: {data.ticket.status}</div>
      </div>

      {data.ticket.ai_summary && (
        <div className="p-4 border-b border-cosmic-700/50 bg-velvet-900/20">
          <p className="text-xs font-medium text-neon-violet mb-1">{t('dashboard.aiSummary')}</p>
          <p className="text-sm text-gray-300">{data.ticket.ai_summary}</p>
        </div>
      )}

      <div className="p-4 space-y-3">
        <p className="text-xs font-semibold text-gray-500 uppercase tracking-wider">{t('dashboard.messages')}</p>
        {data.messages?.map((msg) => (
          <div key={msg.id} className={cn(
            'text-sm p-3 rounded-lg border',
            msg.role === 'customer' ? 'bg-cosmic-800/50 border-cosmic-700/30' :
            msg.role === 'ai' ? 'bg-velvet-900/20 border-velvet-600/20' :
            'bg-cosmic-800/30 border-cosmic-700/20'
          )}>
            <Badge variant={msg.role === 'ai' ? 'neon' : msg.role === 'customer' ? 'info' : 'success'} className="mb-1">{msg.role}</Badge>
            <p className="mt-1 text-gray-300">{msg.content}</p>
          </div>
        ))}
      </div>

      {data.actions && data.actions.length > 0 && (
        <div className="p-4 border-t border-cosmic-700/50 space-y-2">
          <p className="text-xs font-semibold text-gray-500 uppercase tracking-wider">{t('dashboard.actionsTaken')}</p>
          {data.actions.map((action) => (
            <Card key={action.id} variant="glass" className="p-3">
              <div className="flex items-center gap-2">
                <span className="text-sm font-medium text-gray-200">{action.type}</span>
                <Badge variant={action.status === 'executed' ? 'success' : 'warning'} dot>{action.status}</Badge>
                <span className="text-xs text-gray-500 ml-auto font-mono">{(action.confidence * 100).toFixed(0)}%</span>
              </div>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}

export function AgentDashboardPage() {
  const { t } = useTranslation();
  const [statusFilter, setStatusFilter] = useState('');
  const [selectedTicket, setSelectedTicket] = useState<string | null>(null);

  const { data, isLoading } = useQuery({
    queryKey: ['tickets', statusFilter],
    queryFn: () => ticketsApi.list(statusFilter ? { status: statusFilter } : {}).then((r) => r.data),
    refetchInterval: 5000,
  });

  const statuses = ['', 'open', 'in_progress', 'waiting', 'resolved', 'closed'];

  return (
    <div className="flex h-screen">
      <div className={cn('flex-1 flex flex-col', selectedTicket ? 'max-w-[60%]' : '')}>
        <div className="bg-cosmic-900/80 backdrop-blur-sm border-b border-cosmic-700/50 px-6 py-3 flex items-center justify-between">
          <div>
            <h1 className="text-lg font-semibold text-white">{t('dashboard.title')}</h1>
            <span className="text-xs text-gray-500">{data?.total || 0} {t('dashboard.tickets')}</span>
          </div>
          <div className="flex items-center gap-1 bg-cosmic-800 rounded-lg p-1">
            {statuses.map((s) => (
              <button
                key={s}
                onClick={() => setStatusFilter(s)}
                className={cn(
                  'px-3 py-1.5 rounded-md text-xs font-medium transition-all',
                  statusFilter === s
                    ? 'bg-velvet-600 text-white neon-glow-sm'
                    : 'text-gray-400 hover:text-gray-200'
                )}
              >
                {s || t('dashboard.all')}
              </button>
            ))}
          </div>
        </div>

        <div className="flex-1 overflow-y-auto">
          {isLoading ? (
            <div className="p-6 text-gray-500">{t('dashboard.loading')}</div>
          ) : data?.tickets?.length ? (
            data.tickets.map((ticket) => (
              <TicketRow key={ticket.id} ticket={ticket} isSelected={selectedTicket === ticket.id} onSelect={setSelectedTicket} />
            ))
          ) : (
            <div className="p-12 text-center text-gray-500">
              <AlertCircle className="w-12 h-12 mx-auto mb-3 text-gray-600" />
              <p>{t('dashboard.noTickets')}</p>
            </div>
          )}
        </div>
      </div>

      {selectedTicket && (
        <div className="w-[40%]">
          <TicketDetailPanel ticketId={selectedTicket} onClose={() => setSelectedTicket(null)} />
        </div>
      )}
    </div>
  );
}

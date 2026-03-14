import { useQuery, useQueryClient } from '@tanstack/react-query';
import { ticketsApi } from '@/api/client';
import type { Ticket } from '@/api/client';
import { useState, useRef, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import {
  Clock,
  AlertCircle,
  CheckCircle,
  User,
  Bot,
  ChevronRight,
  Globe,
  MessageSquare,
  Mail,
  Send,
  Sparkles,
  Loader2,
  ArrowLeft,
  Shield,
  ThumbsUp,
  ThumbsDown,
} from 'lucide-react';
import { cn } from '@/lib/cn';
import { Badge } from '@/components/ui';
import { useAuthStore } from '@/store/auth';

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
      <Badge variant={ticket.status === 'open' ? 'error' : ticket.status === 'resolved' ? 'success' : 'default'} dot>{t(`dashboard.statuses.${ticket.status}`)}</Badge>
      <ChevronRight className="w-4 h-4 text-gray-600" />
    </div>
  );
}

function TicketChatPanel({ ticketId, onClose }: { ticketId: string; onClose: () => void }) {
  const { t } = useTranslation();
  const queryClient = useQueryClient();
  const [replyText, setReplyText] = useState('');
  const [sending, setSending] = useState(false);
  const [suggesting, setSuggesting] = useState(false);
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const { data, isLoading } = useQuery({
    queryKey: ['ticket', ticketId],
    queryFn: () => ticketsApi.get(ticketId).then((r) => r.data),
    refetchInterval: 3000,
  });

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [data?.messages]);

  const handleReply = async () => {
    if (!replyText.trim() || sending) return;
    setSending(true);
    try {
      await ticketsApi.reply(ticketId, replyText.trim());
      setReplyText('');
      queryClient.invalidateQueries({ queryKey: ['ticket', ticketId] });
      queryClient.invalidateQueries({ queryKey: ['tickets'] });
    } finally {
      setSending(false);
    }
  };

  const handleSuggest = async () => {
    if (suggesting) return;
    setSuggesting(true);
    try {
      const { data } = await ticketsApi.suggest(ticketId);
      setReplyText(data.suggestion);
    } finally {
      setSuggesting(false);
    }
  };

  if (isLoading) return <div className="flex items-center justify-center h-full text-gray-500"><Loader2 className="w-6 h-6 animate-spin" /></div>;
  if (!data) return null;

  const ChannelIcon = CHANNEL_ICONS[data.ticket.channel] || Globe;

  const roleIcon = (role: string) => {
    if (role === 'customer') return <User className="w-4 h-4 text-white" />;
    if (role === 'agent') return <Shield className="w-4 h-4 text-white" />;
    return <Bot className="w-4 h-4 text-white" />;
  };

  const roleColor = (role: string) => {
    if (role === 'customer') return 'bg-gradient-to-br from-neon-blue to-neon-cyan';
    if (role === 'agent') return 'bg-gradient-to-br from-neon-green to-emerald-500';
    return 'bg-gradient-to-br from-velvet-600 to-neon-violet';
  };

  return (
    <div className="flex flex-col h-full bg-cosmic-900/80 backdrop-blur-sm border-l border-cosmic-700/50">
      <div className="flex-shrink-0 px-4 py-3 border-b border-cosmic-700/50 flex items-center gap-3">
        <button onClick={onClose} className="text-gray-500 hover:text-gray-300 transition-colors">
          <ArrowLeft className="w-5 h-5" />
        </button>
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2">
            <h3 className="font-semibold text-white text-sm truncate">{data.ticket.subject}</h3>
            <Badge variant={data.ticket.status === 'open' ? 'error' : data.ticket.status === 'resolved' ? 'success' : 'default'} dot>{t(`dashboard.statuses.${data.ticket.status}`)}</Badge>
          </div>
          <div className="flex items-center gap-2 mt-0.5">
            <User className="w-3 h-3 text-gray-500" />
            <span className="text-xs text-gray-400">{data.customer.name}</span>
            <span className="text-xs text-gray-600">{data.customer.email}</span>
            <ChannelIcon className="w-3 h-3 text-gray-500 ml-1" />
            <Badge variant={PRIORITY_BADGE[data.ticket.priority] || 'default'} className="ml-auto">{data.ticket.priority}</Badge>
          </div>
        </div>
      </div>

      <div className="flex-1 overflow-y-auto p-4 space-y-4 scrollbar-hide">
        {data.ticket.ai_summary && (
          <div className="mx-auto max-w-md bg-velvet-900/20 border border-velvet-600/20 rounded-lg px-3 py-2 text-center">
            <p className="text-xs text-neon-violet">{t('dashboard.aiSummary')}</p>
            <p className="text-xs text-gray-400 mt-1">{data.ticket.ai_summary}</p>
          </div>
        )}

        {data.actions?.filter((a) => a.status === 'pending').map((a) => (
          <div key={a.id} className="mx-auto max-w-md bg-amber-900/20 border border-amber-500/30 rounded-lg px-4 py-3">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-xs font-semibold text-amber-400">{t(`chat.actionLabels.${a.type}`)}</p>
                <p className="text-xs text-gray-400 mt-0.5">{a.params}</p>
              </div>
              <div className="flex gap-2">
                <button
                  onClick={() => ticketsApi.approveAction(a.id, true, '').then(() => queryClient.invalidateQueries({ queryKey: ['ticket', ticketId] }))}
                  className="flex items-center gap-1 px-2.5 py-1.5 rounded-lg text-xs font-medium bg-neon-green/20 border border-neon-green/30 text-neon-green hover:bg-neon-green/30 transition-all"
                >
                  <ThumbsUp className="w-3 h-3" />
                  Approve
                </button>
                <button
                  onClick={() => ticketsApi.approveAction(a.id, false, '').then(() => queryClient.invalidateQueries({ queryKey: ['ticket', ticketId] }))}
                  className="flex items-center gap-1 px-2.5 py-1.5 rounded-lg text-xs font-medium bg-red-500/20 border border-red-500/30 text-red-400 hover:bg-red-500/30 transition-all"
                >
                  <ThumbsDown className="w-3 h-3" />
                  Reject
                </button>
              </div>
            </div>
          </div>
        ))}

        {data.messages?.map((msg) => {
          const isCustomer = msg.role === 'customer';
          return (
            <div key={msg.id} className={cn('flex gap-2.5', isCustomer ? '' : 'flex-row-reverse')}>
              <div className={cn('w-7 h-7 rounded-full flex items-center justify-center flex-shrink-0', roleColor(msg.role))}>
                {roleIcon(msg.role)}
              </div>
              <div className={cn('max-w-[75%]')}>
                <div className="flex items-center gap-1.5 mb-1">
                  <span className={cn('text-xs font-medium', isCustomer ? 'text-neon-cyan' : msg.role === 'agent' ? 'text-neon-green' : 'text-neon-purple')}>
                    {msg.role === 'customer' ? data.customer.name : msg.role === 'agent' ? t('dashboard.agent') : 'AI'}
                  </span>
                  <span className="text-xs text-gray-600">
                    {new Date(msg.created_at).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
                  </span>
                </div>
                <div className={cn(
                  'rounded-xl px-3.5 py-2.5 text-sm',
                  isCustomer
                    ? 'bg-cosmic-800/60 border border-cosmic-700/30 text-gray-300'
                    : msg.role === 'agent'
                      ? 'bg-emerald-900/20 border border-emerald-600/20 text-gray-300'
                      : 'bg-velvet-900/20 border border-velvet-600/20 text-gray-300'
                )}>
                  {msg.content}
                </div>
              </div>
            </div>
          );
        })}
        <div ref={messagesEndRef} />
      </div>

      <div className="flex-shrink-0 p-3 border-t border-cosmic-700/50 bg-cosmic-900/95 backdrop-blur-sm">
        <div className="flex items-center gap-2">
          <button
            onClick={handleSuggest}
            disabled={suggesting}
            className="flex-shrink-0 flex items-center gap-1.5 px-3 py-2 rounded-lg text-xs font-medium bg-velvet-900/50 border border-velvet-600/30 text-neon-purple hover:bg-velvet-900/80 hover:border-neon-violet/50 disabled:opacity-50 transition-all"
          >
            {suggesting ? <Loader2 className="w-3.5 h-3.5 animate-spin" /> : <Sparkles className="w-3.5 h-3.5" />}
            {t('dashboard.aiSuggest')}
          </button>
          <textarea
            value={replyText}
            onChange={(e) => setReplyText(e.target.value)}
            onKeyDown={(e) => {
              if (e.key === 'Enter' && !e.shiftKey) {
                e.preventDefault();
                handleReply();
              }
            }}
            placeholder={t('dashboard.replyPlaceholder')}
            rows={1}
            className="flex-1 bg-cosmic-800/80 border border-cosmic-600 text-gray-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-neon-violet/50 focus:ring-2 focus:ring-neon-violet/20 placeholder-gray-500 transition-all resize-none"
            disabled={sending}
          />
          <button
            onClick={handleReply}
            disabled={sending || !replyText.trim()}
            className="flex-shrink-0 bg-velvet-600 text-white rounded-lg px-3 py-2 hover:bg-velvet-500 disabled:opacity-50 transition-all shadow-[0_0_20px_rgba(109,40,217,0.3)] hover:shadow-[0_0_30px_rgba(109,40,217,0.5)]"
          >
            <Send className="w-4 h-4" />
          </button>
        </div>
      </div>
    </div>
  );
}

export function AgentDashboardPage() {
  const { t } = useTranslation();
  const user = useAuthStore((s) => s.user);
  const [statusFilter, setStatusFilter] = useState('');
  const [selectedTicket, setSelectedTicket] = useState<string | null>(null);

  const params: Record<string, string> = {};
  if (statusFilter) params.status = statusFilter;
  if (user && user.level === 1) params.agent_id = user.id;

  const { data, isLoading } = useQuery({
    queryKey: ['tickets', statusFilter, user?.level],
    queryFn: () => ticketsApi.list(Object.keys(params).length ? params : undefined).then((r) => r.data),
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
                {s ? t(`dashboard.statuses.${s}`) : t('dashboard.all')}
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
          <TicketChatPanel ticketId={selectedTicket} onClose={() => setSelectedTicket(null)} />
        </div>
      )}
    </div>
  );
}

import { useState } from 'react';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import {
  MessageSquare,
  Mail,
  Globe,
  ShoppingBag,
  Send,
  Loader2,
  CheckCircle,
  Copy,
  ExternalLink,
  Plug,
  Code,
  Zap,
  Power,
  AlertCircle,
} from 'lucide-react';
import { cn } from '@/lib/cn';
import { Badge } from '@/components/ui';
import { chatApi, integrationsApi } from '@/api/client';

interface ChannelDef {
  id: string;
  type: string;
  icon: typeof Globe;
  gradient: string;
  shadow: string;
  fields: { key: string; label: string; placeholder: string; type?: string }[];
}

const CHANNELS: ChannelDef[] = [
  {
    id: 'telegram',
    type: 'telegram',
    icon: MessageSquare,
    gradient: 'from-blue-500 to-cyan-400',
    shadow: 'rgba(34,211,238,0.3)',
    fields: [
      { key: 'bot_token', label: 'Bot Token', placeholder: '123456:ABC-DEF...', type: 'password' },
    ],
  },
  {
    id: 'email',
    type: 'email',
    icon: Mail,
    gradient: 'from-red-500 to-orange-400',
    shadow: 'rgba(251,146,60,0.3)',
    fields: [
      { key: 'smtp_host', label: 'SMTP Host', placeholder: 'smtp.gmail.com' },
      { key: 'smtp_user', label: 'Email', placeholder: 'support@company.com' },
      { key: 'smtp_pass', label: 'Password', placeholder: 'app password', type: 'password' },
    ],
  },
  {
    id: 'web',
    type: 'web',
    icon: Globe,
    gradient: 'from-velvet-600 to-neon-violet',
    shadow: 'rgba(168,85,247,0.3)',
    fields: [],
  },
  {
    id: 'shopify',
    type: 'shopify',
    icon: ShoppingBag,
    gradient: 'from-green-500 to-emerald-400',
    shadow: 'rgba(52,211,153,0.3)',
    fields: [
      { key: 'shop_url', label: 'Shop URL', placeholder: 'mystore.myshopify.com' },
      { key: 'api_key', label: 'API Key', placeholder: 'shpat_...', type: 'password' },
    ],
  },
];

const CODE_EXAMPLE = `curl -X POST /api/chat \\
  -H "Content-Type: application/json" \\
  -d '{
    "customer_id": "cust-123",
    "channel": "telegram",
    "message": "I was charged twice"
  }'`;

const RESPONSE_EXAMPLE = `{
  "ticket_id": "t-abc123",
  "message": "I found your account...",
  "actions": [{
    "type": "refund",
    "status": "executed",
    "confidence": 0.9
  }],
  "auto_fixed": true
}`;

function ChannelCard({ channel, status, onSelect, isSelected }: {
  channel: ChannelDef;
  status: string;
  isSelected: boolean;
  onSelect: () => void;
}) {
  const { t } = useTranslation();
  const Icon = channel.icon;
  const isConnected = status === 'connected';

  return (
    <button
      onClick={onSelect}
      className={cn(
        'relative flex flex-col items-center gap-3 p-5 rounded-xl border transition-all duration-300',
        isSelected
          ? 'bg-cosmic-800/80 border-neon-violet/50'
          : 'bg-cosmic-800/40 border-cosmic-700/30 hover:border-cosmic-600/50 hover:bg-cosmic-800/60'
      )}
      style={isSelected ? { boxShadow: `0 0 30px ${channel.shadow}` } : undefined}
    >
      {isConnected && (
        <div className="absolute top-2 right-2 w-2.5 h-2.5 rounded-full bg-neon-green animate-pulse" />
      )}
      <div className={cn('w-12 h-12 rounded-xl bg-gradient-to-br flex items-center justify-center', channel.gradient)}>
        <Icon className="w-6 h-6 text-white" />
      </div>
      <span className="text-sm font-medium text-gray-200">{t(`integrations.channels.${channel.id}`)}</span>
      <Badge variant={isConnected ? 'success' : 'default'} dot={isConnected}>
        {isConnected ? t('integrations.connected') : t('integrations.available')}
      </Badge>
    </button>
  );
}

function ChannelConfig({ channel, status, onStatusChange }: {
  channel: ChannelDef;
  status: string;
  onStatusChange: () => void;
}) {
  const { t } = useTranslation();
  const [config, setConfig] = useState<Record<string, string>>({});
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const isConnected = status === 'connected';

  const handleConnect = async () => {
    setLoading(true);
    setError('');
    try {
      await integrationsApi.connect(channel.type, { ...config, name: channel.id });
      onStatusChange();
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Connection failed';
      setError(msg);
    } finally {
      setLoading(false);
    }
  };

  const handleDisconnect = async () => {
    setLoading(true);
    try {
      await integrationsApi.disconnect(channel.type);
      onStatusChange();
    } finally {
      setLoading(false);
    }
  };

  if (channel.fields.length === 0) {
    return (
      <div className="bg-cosmic-800/40 border border-cosmic-700/30 rounded-xl p-4">
        <div className="flex items-center gap-2 text-sm text-gray-400">
          <CheckCircle className="w-4 h-4 text-neon-green" />
          {t('integrations.builtIn')}
        </div>
      </div>
    );
  }

  return (
    <div className="bg-cosmic-800/40 border border-cosmic-700/30 rounded-xl p-4 space-y-3">
      {!isConnected && channel.fields.map((field) => (
        <div key={field.key}>
          <label className="text-xs font-medium text-gray-400 mb-1 block">{field.label}</label>
          <input
            type={field.type || 'text'}
            value={config[field.key] || ''}
            onChange={(e) => setConfig({ ...config, [field.key]: e.target.value })}
            placeholder={field.placeholder}
            className="w-full bg-cosmic-900/80 border border-cosmic-600 text-gray-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-neon-violet/50 placeholder-gray-600 transition-all font-mono"
            disabled={loading}
          />
        </div>
      ))}

      {error && (
        <div className="flex items-center gap-2 text-xs text-red-400">
          <AlertCircle className="w-3.5 h-3.5" />
          {error}
        </div>
      )}

      <button
        onClick={isConnected ? handleDisconnect : handleConnect}
        disabled={loading || (!isConnected && channel.fields.some((f) => !config[f.key]))}
        className={cn(
          'w-full flex items-center justify-center gap-2 px-4 py-2.5 rounded-lg text-sm font-medium transition-all',
          isConnected
            ? 'bg-red-900/30 border border-red-600/30 text-red-400 hover:bg-red-900/50'
            : 'bg-velvet-600 text-white hover:bg-velvet-500 disabled:opacity-50',
        )}
      >
        {loading ? (
          <Loader2 className="w-4 h-4 animate-spin" />
        ) : (
          <Power className="w-4 h-4" />
        )}
        {isConnected ? t('integrations.disconnect') : t('integrations.connectBtn')}
      </button>
    </div>
  );
}

function ApiPlayground() {
  const { t } = useTranslation();
  const [channel, setChannel] = useState('telegram');
  const [message, setMessage] = useState('');
  const [customerId] = useState('cust-demo');
  const [response, setResponse] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [copied, setCopied] = useState(false);

  const handleSend = async () => {
    if (!message.trim() || loading) return;
    setLoading(true);
    setResponse(null);
    try {
      const { data } = await chatApi.send({
        customer_id: customerId,
        channel,
        message: message.trim(),
      });
      setResponse(JSON.stringify(data, null, 2));
    } catch (err: unknown) {
      const errorMessage = err instanceof Error ? err.message : 'Unknown error';
      setResponse(JSON.stringify({ error: errorMessage }, null, 2));
    } finally {
      setLoading(false);
    }
  };

  const handleCopy = () => {
    if (response) {
      navigator.clipboard.writeText(response);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    }
  };

  return (
    <div className="bg-cosmic-800/40 border border-cosmic-700/30 rounded-xl overflow-hidden">
      <div className="px-4 py-3 border-b border-cosmic-700/30 flex items-center gap-2">
        <Zap className="w-4 h-4 text-neon-violet" />
        <span className="text-sm font-semibold text-white">{t('integrations.playground')}</span>
        <span className="text-xs text-gray-500 ml-auto">POST /api/chat</span>
      </div>

      <div className="p-4 space-y-3">
        <div className="flex gap-2">
          {['telegram', 'email', 'web', 'shopify'].map((ch) => (
            <button
              key={ch}
              onClick={() => setChannel(ch)}
              className={cn(
                'px-3 py-1.5 rounded-lg text-xs font-medium transition-all',
                channel === ch
                  ? 'bg-velvet-600 text-white'
                  : 'bg-cosmic-800 text-gray-400 hover:text-gray-200'
              )}
            >
              {ch}
            </button>
          ))}
        </div>

        <div className="flex gap-2">
          <input
            type="text"
            value={message}
            onChange={(e) => setMessage(e.target.value)}
            onKeyDown={(e) => e.key === 'Enter' && handleSend()}
            placeholder={t('integrations.messagePlaceholder')}
            className="flex-1 bg-cosmic-900/80 border border-cosmic-600 text-gray-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-neon-violet/50 placeholder-gray-500 transition-all"
            disabled={loading}
          />
          <button
            onClick={handleSend}
            disabled={loading || !message.trim()}
            className="bg-velvet-600 text-white rounded-lg px-4 py-2 hover:bg-velvet-500 disabled:opacity-50 transition-all"
          >
            {loading ? <Loader2 className="w-4 h-4 animate-spin" /> : <Send className="w-4 h-4" />}
          </button>
        </div>

        {response && (
          <div className="relative">
            <button
              onClick={handleCopy}
              className="absolute top-2 right-2 p-1.5 rounded-md bg-cosmic-700/50 text-gray-400 hover:text-white transition-colors"
            >
              {copied ? <CheckCircle className="w-3.5 h-3.5 text-neon-green" /> : <Copy className="w-3.5 h-3.5" />}
            </button>
            <pre className="bg-cosmic-900/80 border border-cosmic-700/30 rounded-lg p-3 text-xs text-neon-green font-mono overflow-x-auto max-h-64 overflow-y-auto">
              {response}
            </pre>
          </div>
        )}
      </div>
    </div>
  );
}

export function IntegrationsPage() {
  const { t } = useTranslation();
  const queryClient = useQueryClient();
  const [selectedChannel, setSelectedChannel] = useState('telegram');
  const [copiedCode, setCopiedCode] = useState(false);

  const { data: activeIntegrations } = useQuery({
    queryKey: ['integrations'],
    queryFn: () => integrationsApi.list().then((r) => r.data),
    refetchInterval: 5000,
  });

  const getStatus = (id: string) => {
    const found = activeIntegrations?.find((i) => i.id === id || i.type === id);
    return found?.status || 'not_configured';
  };

  const handleCopyCode = (text: string) => {
    navigator.clipboard.writeText(text);
    setCopiedCode(true);
    setTimeout(() => setCopiedCode(false), 2000);
  };

  const selectedDef = CHANNELS.find((c) => c.id === selectedChannel)!;

  return (
    <div className="h-full overflow-y-auto scrollbar-hide">
      <div className="max-w-5xl mx-auto p-6 space-y-8">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-velvet-600 to-neon-violet flex items-center justify-center">
            <Plug className="w-5 h-5 text-white" />
          </div>
          <div>
            <h1 className="text-xl font-bold text-white">{t('integrations.title')}</h1>
            <p className="text-sm text-gray-400">{t('integrations.subtitle')}</p>
          </div>
        </div>

        <div>
          <h2 className="text-sm font-semibold text-gray-400 uppercase tracking-wider mb-4">{t('integrations.channelsTitle')}</h2>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            {CHANNELS.map((ch) => (
              <ChannelCard
                key={ch.id}
                channel={ch}
                status={getStatus(ch.id)}
                isSelected={selectedChannel === ch.id}
                onSelect={() => setSelectedChannel(ch.id)}
              />
            ))}
          </div>
        </div>

        <ChannelConfig
          channel={selectedDef}
          status={getStatus(selectedDef.id)}
          onStatusChange={() => queryClient.invalidateQueries({ queryKey: ['integrations'] })}
        />

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <div>
            <div className="flex items-center gap-2 mb-3">
              <Code className="w-4 h-4 text-neon-purple" />
              <h2 className="text-sm font-semibold text-gray-400 uppercase tracking-wider">{t('integrations.requestExample')}</h2>
            </div>
            <div className="relative">
              <button
                onClick={() => handleCopyCode(CODE_EXAMPLE)}
                className="absolute top-2 right-2 p-1.5 rounded-md bg-cosmic-700/50 text-gray-400 hover:text-white transition-colors z-10"
              >
                {copiedCode ? <CheckCircle className="w-3.5 h-3.5 text-neon-green" /> : <Copy className="w-3.5 h-3.5" />}
              </button>
              <pre className="bg-cosmic-800/60 border border-cosmic-700/30 rounded-xl p-4 text-xs text-gray-300 font-mono overflow-x-auto">
                <code>{CODE_EXAMPLE}</code>
              </pre>
            </div>
          </div>

          <div>
            <div className="flex items-center gap-2 mb-3">
              <ExternalLink className="w-4 h-4 text-neon-green" />
              <h2 className="text-sm font-semibold text-gray-400 uppercase tracking-wider">{t('integrations.responseExample')}</h2>
            </div>
            <pre className="bg-cosmic-800/60 border border-cosmic-700/30 rounded-xl p-4 text-xs text-neon-green font-mono overflow-x-auto">
              <code>{RESPONSE_EXAMPLE}</code>
            </pre>
          </div>
        </div>

        <ApiPlayground />

        <div>
          <h2 className="text-sm font-semibold text-gray-400 uppercase tracking-wider mb-4">{t('integrations.howItWorks')}</h2>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            {[1, 2, 3].map((step) => (
              <div key={step} className="bg-cosmic-800/40 border border-cosmic-700/30 rounded-xl p-4">
                <div className="w-8 h-8 rounded-full bg-velvet-600/30 border border-velvet-600/50 flex items-center justify-center mb-3">
                  <span className="text-sm font-bold text-neon-purple">{step}</span>
                </div>
                <h3 className="text-sm font-semibold text-white mb-1">{t(`integrations.step${step}Title`)}</h3>
                <p className="text-xs text-gray-400">{t(`integrations.step${step}Desc`)}</p>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}

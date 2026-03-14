import { useState, useRef, useEffect } from 'react';
import { Send, Bot, User, Zap, CheckCircle, AlertTriangle, Loader2 } from 'lucide-react';
import ReactMarkdown from 'react-markdown';
import { useTranslation } from 'react-i18next';
import { useChatStore } from '@/store/chat';
import { chatApi } from '@/api/client';
import type { Action } from '@/api/client';
import { cn } from '@/lib/cn';
import { Badge } from '@/components/ui';

function ActionCard({ action }: { action: Action }) {
  const { t } = useTranslation();
  const result = action.result ? JSON.parse(action.result) : null;
  const label = t(`chat.actionLabels.${action.type}`, { defaultValue: action.type });

  return (
    <div className="mt-2 bg-velvet-900/30 border border-velvet-600/30 rounded-lg p-3">
      <div className="flex items-center gap-2 text-sm font-medium text-neon-purple">
        <Zap className="w-4 h-4 text-neon-violet" />
        {label}
        {result?.success ? (
          <CheckCircle className="w-4 h-4 text-neon-green ml-auto" />
        ) : (
          <AlertTriangle className="w-4 h-4 text-red-400 ml-auto" />
        )}
      </div>
      {result?.message && (
        <p className="text-xs text-gray-400 mt-1">{result.message}</p>
      )}
      <div className="text-xs text-gray-500 mt-1 font-mono">
        confidence: {(action.confidence * 100).toFixed(0)}%
      </div>
    </div>
  );
}

export function CustomerChatPage() {
  const { t } = useTranslation();
  const [input, setInput] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const { messages, ticketId, customerId, addMessage, setTicketId } = useChatStore();

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  const handleSend = async () => {
    if (!input.trim() || isLoading) return;

    const text = input.trim();
    setInput('');
    addMessage({ role: 'customer', content: text });
    setIsLoading(true);
    addMessage({ role: 'ai', content: '', isLoading: true });

    try {
      const { data } = await chatApi.send({
        ticket_id: ticketId || undefined,
        customer_id: customerId,
        channel: 'web',
        message: text,
      });

      if (data.ticket_id && !ticketId) {
        setTicketId(data.ticket_id);
      }

      const store = useChatStore.getState();
      const msgs = [...store.messages];
      msgs[msgs.length - 1] = {
        ...msgs[msgs.length - 1],
        content: data.message,
        actions: data.actions,
        analysis: data.analysis || undefined,
        isLoading: false,
      };
      useChatStore.setState({ messages: msgs });
    } catch {
      const store = useChatStore.getState();
      const msgs = [...store.messages];
      msgs[msgs.length - 1] = {
        ...msgs[msgs.length - 1],
        content: t('chat.error'),
        isLoading: false,
      };
      useChatStore.setState({ messages: msgs });
    } finally {
      setIsLoading(false);
    }
  };

  const quickReplies = [
    t('chat.quickReplies.charged'),
    t('chat.quickReplies.password'),
    t('chat.quickReplies.plan'),
    t('chat.quickReplies.cancel'),
  ];

  return (
    <div className="flex flex-col h-full">
      <div className="flex-1 overflow-y-auto p-6 space-y-4 scrollbar-hide">
        {messages.length === 0 && (
          <div className="text-center mt-16">
            <div className="w-20 h-20 mx-auto mb-6 rounded-2xl bg-gradient-to-br from-velvet-600 to-neon-violet flex items-center justify-center animate-pulse-neon">
              <Bot className="w-10 h-10 text-white" />
            </div>
            <p className="text-xl text-gray-300">{t('chat.welcome')}</p>
            <p className="text-sm mt-2 text-gray-500">{t('chat.welcomeSub')}</p>
            <div className="flex flex-wrap gap-2 justify-center mt-8">
              {quickReplies.map((q) => (
                <button
                  key={q}
                  onClick={() => setInput(q)}
                  className="px-4 py-2 bg-cosmic-800/80 border border-cosmic-700/50 rounded-full text-sm text-gray-300 hover:border-neon-violet/30 hover:text-neon-purple hover:bg-cosmic-800 transition-all duration-200"
                >
                  {q}
                </button>
              ))}
            </div>
          </div>
        )}

        {messages.map((msg) => (
          <div
            key={msg.id}
            className={cn(
              'flex gap-3 max-w-3xl mx-auto',
              msg.role === 'customer' ? 'ml-auto flex-row-reverse' : ''
            )}
          >
            <div className={cn(
              'w-8 h-8 rounded-full flex items-center justify-center flex-shrink-0',
              msg.role === 'customer'
                ? 'bg-gradient-to-br from-neon-blue to-neon-cyan'
                : 'bg-gradient-to-br from-velvet-600 to-neon-violet'
            )}>
              {msg.role === 'customer' ? <User className="w-4 h-4 text-white" /> : <Bot className="w-4 h-4 text-white" />}
            </div>
            <div className={cn(
              'rounded-xl px-4 py-3 max-w-lg',
              msg.role === 'customer'
                ? 'bg-velvet-600/80 text-white'
                : 'bg-cosmic-800/80 border border-cosmic-700/50 backdrop-blur-sm'
            )}>
              {msg.isLoading ? (
                <div className="flex items-center gap-2 text-gray-400">
                  <Loader2 className="w-4 h-4 animate-spin text-neon-violet" />
                  <span className="text-sm animate-shimmer">{t('chat.thinking')}</span>
                </div>
              ) : (
                <>
                  <div className="text-sm prose prose-sm prose-invert max-w-none">
                    <ReactMarkdown>{msg.content}</ReactMarkdown>
                  </div>
                  {msg.actions && msg.actions.length > 0 && (
                    <div className="mt-2 space-y-2">
                      {msg.actions.map((action) => (
                        <ActionCard key={action.id} action={action} />
                      ))}
                    </div>
                  )}
                  {msg.analysis && (
                    <div className="mt-2 flex flex-wrap gap-1">
                      <Badge variant="neon">{msg.analysis.intent}</Badge>
                      <Badge variant={msg.analysis.sentiment === 'angry' ? 'error' : msg.analysis.sentiment === 'negative' ? 'warning' : 'success'}>{msg.analysis.sentiment}</Badge>
                      <Badge variant={msg.analysis.urgency === 'high' ? 'error' : msg.analysis.urgency === 'low' ? 'success' : 'warning'}>{msg.analysis.urgency}</Badge>
                    </div>
                  )}
                </>
              )}
            </div>
          </div>
        ))}
        <div ref={messagesEndRef} />
      </div>

      <div className="bg-cosmic-900/80 backdrop-blur-sm border-t border-cosmic-700/50 p-4">
        <div className="max-w-3xl mx-auto flex gap-3">
          <input
            type="text"
            value={input}
            onChange={(e) => setInput(e.target.value)}
            onKeyDown={(e) => e.key === 'Enter' && handleSend()}
            placeholder={t('chat.placeholder')}
            className="flex-1 bg-cosmic-800/80 border border-cosmic-600 text-gray-200 rounded-xl px-4 py-2.5 text-sm focus:outline-none focus:border-neon-violet/50 focus:ring-2 focus:ring-neon-violet/20 placeholder-gray-500 transition-all backdrop-blur-sm"
            disabled={isLoading}
          />
          <button
            onClick={handleSend}
            disabled={isLoading || !input.trim()}
            className="bg-velvet-600 text-white rounded-xl px-4 py-2.5 hover:bg-velvet-500 disabled:opacity-50 transition-all shadow-[0_0_20px_rgba(109,40,217,0.3)] hover:shadow-[0_0_30px_rgba(109,40,217,0.5)]"
          >
            <Send className="w-5 h-5" />
          </button>
        </div>
      </div>
    </div>
  );
}

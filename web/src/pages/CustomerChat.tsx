import { useState, useRef, useEffect } from 'react';
import { Send, Bot, User, Zap, CheckCircle, AlertTriangle, Loader2, MessageSquare, Globe, Mail } from 'lucide-react';
import ReactMarkdown from 'react-markdown';
import { useChatStore } from '@/store/chat';
import { chatApi } from '@/api/client';
import type { Action } from '@/api/client';
import { cn } from '@/lib/cn';
import { Badge } from '@/components/ui';

const ACTION_LABELS: Record<string, string> = {
  refund: 'Refund Processed',
  change_plan: 'Plan Changed',
  reset_password: 'Password Reset',
  escalate: 'Escalated',
  send_email: 'Email Sent',
  cancel_subscription: 'Subscription Cancelled',
  lookup_billing: 'Billing Lookup',
  lookup_customer: 'Customer Lookup',
};

const CHANNELS = [
  { id: 'web', label: 'Web', icon: Globe },
  { id: 'telegram', label: 'Telegram', icon: MessageSquare },
  { id: 'email', label: 'Email', icon: Mail },
];

const DEMO_CUSTOMERS = [
  { id: 'a0000000-0000-0000-0000-000000000001', name: 'Ivan Petrov' },
  { id: 'a0000000-0000-0000-0000-000000000002', name: 'Maria Sidorova' },
  { id: 'a0000000-0000-0000-0000-000000000003', name: 'Alex Kozlov' },
  { id: 'a0000000-0000-0000-0000-000000000004', name: 'Elena Novikova' },
  { id: 'a0000000-0000-0000-0000-000000000005', name: 'Dmitry Volkov' },
];

function ActionCard({ action }: { action: Action }) {
  const result = action.result ? JSON.parse(action.result) : null;

  return (
    <div className="mt-2 bg-velvet-900/30 border border-velvet-600/30 rounded-lg p-3">
      <div className="flex items-center gap-2 text-sm font-medium text-neon-purple">
        <Zap className="w-4 h-4 text-neon-violet" />
        {ACTION_LABELS[action.type] || action.type}
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
  const [input, setInput] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [channel, setChannel] = useState('web');
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const { messages, ticketId, customerId, addMessage, setTicketId, setCustomerId, clearChat } = useChatStore();

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
        channel,
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
        content: 'Sorry, something went wrong. Please try again.',
        isLoading: false,
      };
      useChatStore.setState({ messages: msgs });
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="flex flex-col h-screen">
      <div className="bg-cosmic-900 border-b border-cosmic-700/50 px-6 py-3 flex items-center justify-between">
        <div>
          <h1 className="text-lg font-semibold text-white">Customer Support Chat</h1>
          <div className="flex items-center gap-2 mt-0.5">
            {ticketId && (
              <span className="text-xs text-gray-500 font-mono">Ticket: {ticketId.slice(0, 8)}</span>
            )}
            <Badge variant="neon" dot>{channel}</Badge>
          </div>
        </div>
        <div className="flex items-center gap-3">
          <div className="flex bg-cosmic-800 rounded-lg p-1 gap-1">
            {CHANNELS.map((ch) => (
              <button
                key={ch.id}
                onClick={() => { setChannel(ch.id); clearChat(); }}
                className={cn(
                  'flex items-center gap-1.5 px-3 py-1.5 rounded-md text-xs font-medium transition-all',
                  channel === ch.id
                    ? 'bg-velvet-600 text-white neon-glow-sm'
                    : 'text-gray-400 hover:text-gray-200'
                )}
              >
                <ch.icon className="w-3.5 h-3.5" />
                {ch.label}
              </button>
            ))}
          </div>
          <select
            value={customerId}
            onChange={(e) => { setCustomerId(e.target.value); clearChat(); }}
            className="text-sm bg-cosmic-800 border border-cosmic-600 text-gray-300 rounded-lg px-3 py-1.5 focus:border-neon-violet/50 focus:outline-none"
          >
            {DEMO_CUSTOMERS.map((c) => (
              <option key={c.id} value={c.id}>{c.name}</option>
            ))}
          </select>
          <button onClick={clearChat} className="text-xs text-gray-500 hover:text-neon-violet transition-colors">
            New Chat
          </button>
        </div>
      </div>

      <div className="flex-1 overflow-y-auto p-6 space-y-4">
        {messages.length === 0 && (
          <div className="text-center text-gray-500 mt-20">
            <div className="w-20 h-20 mx-auto mb-6 rounded-2xl bg-gradient-to-br from-velvet-600 to-neon-violet flex items-center justify-center animate-pulse-neon">
              <Bot className="w-10 h-10 text-white" />
            </div>
            <p className="text-xl text-gray-300">How can we help you today?</p>
            <p className="text-sm mt-2 text-gray-500">Describe your issue and our AI will assist you</p>
            <div className="flex flex-wrap gap-2 justify-center mt-8">
              {['I was charged twice', 'Reset my password', 'Change my plan to premium', 'Cancel my subscription'].map((q) => (
                <button
                  key={q}
                  onClick={() => setInput(q)}
                  className="px-4 py-2 bg-cosmic-800 border border-cosmic-700/50 rounded-full text-sm text-gray-300 hover:border-neon-violet/30 hover:text-neon-purple transition-all duration-200"
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
              'flex gap-3 max-w-3xl',
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
                ? 'bg-velvet-600 text-white'
                : 'bg-cosmic-800 border border-cosmic-700/50'
            )}>
              {msg.isLoading ? (
                <div className="flex items-center gap-2 text-gray-400">
                  <Loader2 className="w-4 h-4 animate-spin text-neon-violet" />
                  <span className="text-sm animate-shimmer">Thinking...</span>
                </div>
              ) : (
                <>
                  <div className={cn('text-sm prose prose-sm prose-invert max-w-none')}>
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

      <div className="bg-cosmic-900 border-t border-cosmic-700/50 p-4">
        <div className="max-w-3xl mx-auto flex gap-3">
          <input
            type="text"
            value={input}
            onChange={(e) => setInput(e.target.value)}
            onKeyDown={(e) => e.key === 'Enter' && handleSend()}
            placeholder="Type your message..."
            className="flex-1 bg-cosmic-800 border border-cosmic-600 text-gray-200 rounded-xl px-4 py-2.5 text-sm focus:outline-none focus:border-neon-violet/50 focus:ring-2 focus:ring-neon-violet/20 placeholder-gray-500 transition-all"
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

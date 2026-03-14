import { create } from 'zustand';
import type { ChatResponse, AIAnalysis } from '@/api/client';

interface ChatMessage {
  id: string;
  role: 'customer' | 'ai' | 'system';
  content: string;
  timestamp: Date;
  actions?: ChatResponse['actions'];
  analysis?: AIAnalysis;
  isLoading?: boolean;
}

interface ChatState {
  messages: ChatMessage[];
  ticketId: string | null;
  customerId: string;
  isConnected: boolean;
  addMessage: (msg: Omit<ChatMessage, 'id' | 'timestamp'>) => void;
  setTicketId: (id: string) => void;
  setCustomerId: (id: string) => void;
  setConnected: (connected: boolean) => void;
  updateLastMessage: (update: Partial<ChatMessage>) => void;
  clearChat: () => void;
}

export const useChatStore = create<ChatState>((set) => ({
  messages: [],
  ticketId: null,
  customerId: 'a0000000-0000-0000-0000-000000000001',
  isConnected: false,

  addMessage: (msg) =>
    set((state) => ({
      messages: [
        ...state.messages,
        { ...msg, id: crypto.randomUUID(), timestamp: new Date() },
      ],
    })),

  setTicketId: (id) => set({ ticketId: id }),
  setCustomerId: (id) => set({ customerId: id }),
  setConnected: (connected) => set({ isConnected: connected }),

  updateLastMessage: (update) =>
    set((state) => {
      const messages = [...state.messages];
      if (messages.length > 0) {
        messages[messages.length - 1] = { ...messages[messages.length - 1], ...update };
      }
      return { messages };
    }),

  clearChat: () => set({ messages: [], ticketId: null }),
}));

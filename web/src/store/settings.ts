import { create } from 'zustand';

interface SettingsState {
  activeProvider: string;
  providers: string[];
  setActiveProvider: (provider: string) => void;
  setProviders: (providers: string[]) => void;
}

export const useSettingsStore = create<SettingsState>((set) => ({
  activeProvider: 'anthropic',
  providers: [],
  setActiveProvider: (provider) => set({ activeProvider: provider }),
  setProviders: (providers) => set({ providers }),
}));

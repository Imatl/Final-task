import { create } from 'zustand';
import api from '@/api/client';

export type UserLevel = 1 | 2 | 3 | 4 | 5;

export interface AuthUser {
  id: string;
  email: string;
  name: string;
  level: UserLevel;
  role: string;
  company?: string;
}

export const ROLE_LABELS: Record<UserLevel, string> = {
  1: 'Support',
  2: 'Senior Support',
  3: 'Team Lead',
  4: 'Admin',
  5: 'Owner',
};

function loadUser(): AuthUser | null {
  try {
    const raw = localStorage.getItem('sf-auth-user');
    if (!raw) return null;
    return JSON.parse(raw) as AuthUser;
  } catch {
    return null;
  }
}

function saveUser(user: AuthUser | null) {
  if (user) {
    localStorage.setItem('sf-auth-user', JSON.stringify(user));
  } else {
    localStorage.removeItem('sf-auth-user');
  }
}

interface AuthState {
  user: AuthUser | null;
  login: (email: string, password: string) => Promise<boolean>;
  loginWithGoogle: (backendUser: AuthUser) => void;
  logout: () => void;
}

export const useAuthStore = create<AuthState>((set) => ({
  user: loadUser(),
  login: async (email: string, password: string) => {
    console.log('[Auth] Login attempt for:', email);
    try {
      const { data } = await api.post<AuthUser>('/auth/login', { email, password });
      console.log('[Auth] Login successful:', data.email);
      const user: AuthUser = {
        id: data.id,
        email: data.email,
        name: data.name,
        level: (data.level || 1) as UserLevel,
        role: data.role || ROLE_LABELS[1],
        company: data.company,
      };
      saveUser(user);
      set({ user });
      return true;
    } catch (err) {
      console.error('[Auth] Login failed:', err);
      return false;
    }
  },
  loginWithGoogle: (backendUser: AuthUser) => {
    console.log('[Auth] Google login for:', backendUser.email);
    const user: AuthUser = {
      id: backendUser.id,
      email: backendUser.email,
      name: backendUser.name,
      level: (backendUser.level || 1) as UserLevel,
      role: backendUser.role || ROLE_LABELS[1],
      company: backendUser.company,
    };
    saveUser(user);
    set({ user });
    console.log('[Auth] User saved to store:', user.id);
  },
  logout: () => {
    console.log('[Auth] Logging out');
    saveUser(null);
    set({ user: null });
  },
}));

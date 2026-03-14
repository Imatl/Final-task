import { create } from 'zustand';

export type UserLevel = 1 | 2 | 3 | 4 | 5;

export interface AuthUser {
  id: string;
  email: string;
  name: string;
  level: UserLevel;
  role: string;
}

export const ROLE_LABELS: Record<UserLevel, string> = {
  1: 'Support',
  2: 'Senior Support',
  3: 'Team Lead',
  4: 'Admin',
  5: 'Owner',
};

const HARDCODED_USERS: Array<AuthUser & { password: string }> = [
  {
    id: 'admin-001',
    email: 'admin@test.com',
    password: 'password123',
    name: 'Admin User',
    level: 5,
    role: 'Owner',
  },
];

function loadUser(): AuthUser | null {
  try {
    const raw = localStorage.getItem('sf-auth-user');
    if (!raw) return null;
    const cached = JSON.parse(raw) as AuthUser;
    const source = HARDCODED_USERS.find((u) => u.email === cached.email);
    if (source) {
      const synced: AuthUser = { ...cached, level: source.level, role: source.role };
      localStorage.setItem('sf-auth-user', JSON.stringify(synced));
      return synced;
    }
    return cached;
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
  login: (email: string, password: string) => boolean;
  logout: () => void;
}

export const useAuthStore = create<AuthState>((set) => ({
  user: loadUser(),
  login: (email: string, password: string) => {
    const found = HARDCODED_USERS.find(
      (u) => u.email === email && u.password === password
    );
    if (found) {
      const { password: _p, ...user } = found;
      saveUser(user);
      set({ user });
      return true;
    }
    return false;
  },
  logout: () => {
    saveUser(null);
    set({ user: null });
  },
}));

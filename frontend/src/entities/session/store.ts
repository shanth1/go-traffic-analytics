import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { api } from '@/shared/api/base';
import type { User, AuthResponse, LoginReq } from '@/shared/api/types';

interface AuthState {
  token: string | null;
  user: User | null;
  isAuthenticated: boolean;
  login: (creds: LoginReq) => Promise<void>;
  logout: () => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      token: null,
      user: null,
      isAuthenticated: false,
      login: async (creds) => {
        const { data } = await api.post<AuthResponse>('/auth/login', creds);
        set({ token: data.token, user: data.user, isAuthenticated: true });
      },
      logout: () => set({ token: null, user: null, isAuthenticated: false }),
    }),
    {
      name: 'auth-storage',
    }
  )
);

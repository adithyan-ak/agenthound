import { create } from "zustand";
import { login as apiLogin } from "@/api/auth";
import type { UserInfo } from "@/api/types";

const TOKEN_KEY = "agenthound_token";
const USER_KEY = "agenthound_user";

function decodeJWTPayload(token: string): Record<string, unknown> | null {
  try {
    const parts = token.split(".");
    if (parts.length !== 3) return null;
    const payload = atob(parts[1]!.replace(/-/g, "+").replace(/_/g, "/"));
    return JSON.parse(payload) as Record<string, unknown>;
  } catch {
    return null;
  }
}

function isTokenExpired(token: string): boolean {
  const payload = decodeJWTPayload(token);
  if (!payload || typeof payload["exp"] !== "number") return true;
  return Date.now() >= payload["exp"] * 1000;
}

function loadInitialState(): { token: string | null; user: UserInfo | null; isAuthenticated: boolean } {
  const token = localStorage.getItem(TOKEN_KEY);
  if (!token || isTokenExpired(token)) {
    localStorage.removeItem(TOKEN_KEY);
    localStorage.removeItem(USER_KEY);
    return { token: null, user: null, isAuthenticated: false };
  }
  try {
    const raw = localStorage.getItem(USER_KEY);
    const user = raw ? (JSON.parse(raw) as UserInfo) : null;
    return { token, user, isAuthenticated: true };
  } catch {
    localStorage.removeItem(TOKEN_KEY);
    localStorage.removeItem(USER_KEY);
    return { token: null, user: null, isAuthenticated: false };
  }
}

const initialState = loadInitialState();

interface AuthState {
  token: string | null;
  user: UserInfo | null;
  isAuthenticated: boolean;
}

interface AuthActions {
  login: (username: string, password: string) => Promise<void>;
  logout: () => void;
  initialize: () => void;
}

export const useAuthStore = create<AuthState & AuthActions>()((set) => ({
  token: initialState.token,
  user: initialState.user,
  isAuthenticated: initialState.isAuthenticated,

  login: async (username, password) => {
    const resp = await apiLogin(username, password);
    localStorage.setItem(TOKEN_KEY, resp.token);
    localStorage.setItem(USER_KEY, JSON.stringify(resp.user));
    set({ token: resp.token, user: resp.user, isAuthenticated: true });
  },

  logout: () => {
    localStorage.removeItem(TOKEN_KEY);
    localStorage.removeItem(USER_KEY);
    set({ token: null, user: null, isAuthenticated: false });
  },

  initialize: () => {
    const token = localStorage.getItem(TOKEN_KEY);
    if (!token || isTokenExpired(token)) {
      localStorage.removeItem(TOKEN_KEY);
      localStorage.removeItem(USER_KEY);
      set({ token: null, user: null, isAuthenticated: false });
      return;
    }
    try {
      const raw = localStorage.getItem(USER_KEY);
      const user = raw ? (JSON.parse(raw) as UserInfo) : null;
      set({ token, user, isAuthenticated: true });
    } catch {
      localStorage.removeItem(TOKEN_KEY);
      localStorage.removeItem(USER_KEY);
      set({ token: null, user: null, isAuthenticated: false });
    }
  },
}));

// src/auth.ts
const TOKEN_KEY = 'token';

export const auth = {
  getToken(): string | null {
    return localStorage.getItem(TOKEN_KEY);
  },

  isAuthenticated(): boolean {
    return !!this.getToken();
  },

  setSession(token: string) {
    localStorage.setItem(TOKEN_KEY, token);
  },

  clearSession() {
    localStorage.removeItem(TOKEN_KEY);
    // Clear sessionStorage too just in case
    sessionStorage.removeItem(TOKEN_KEY);
  }
};
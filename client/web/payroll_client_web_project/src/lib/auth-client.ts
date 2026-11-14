// src/lib/auth-client.ts
import { createAuthClient } from "better-auth/vue"

export const authClient = createAuthClient({
  baseURL: "http://localhost:4002"
})

export const { signIn, signUp, signOut, useSession } = authClient
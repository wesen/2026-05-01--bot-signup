import { createContext } from 'react'
import type { User } from '../store/api'

export interface AuthContextValue {
  user: User | null
  isLoading: boolean
  loginWithDiscord: (returnTo?: string) => void
  logout: () => Promise<void>
}

export const AuthContext = createContext<AuthContextValue | null>(null)

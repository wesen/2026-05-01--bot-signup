import type { ReactNode } from 'react'
import { useLocation } from 'react-router-dom'
import { apiSlice, useGetMeQuery, useLogoutMutation } from '../store/api'
import { AuthContext } from './context'
import { useAppDispatch } from './hooks'

const publicRoutes = new Set(['/', '/tutorial', '/auth/callback'])

export function AuthProvider({ children }: { children: ReactNode }) {
  const dispatch = useAppDispatch()
  const location = useLocation()
  const shouldCheckSession = !publicRoutes.has(location.pathname)
  const { data: user, isLoading } = useGetMeQuery(undefined, {
    skip: !shouldCheckSession,
    refetchOnMountOrArgChange: true,
  })
  const [logoutMutation] = useLogoutMutation()

  const loginWithDiscord = (returnTo = '/waiting-list') => {
    window.location.href = `/auth/discord/login?return_to=${encodeURIComponent(returnTo)}`
  }

  const logout = async () => {
    await logoutMutation().unwrap()
    dispatch(apiSlice.util.resetApiState())
  }

  return (
    <AuthContext.Provider value={{ user: user ?? null, isLoading: shouldCheckSession && isLoading, loginWithDiscord, logout }}>
      {children}
    </AuthContext.Provider>
  )
}

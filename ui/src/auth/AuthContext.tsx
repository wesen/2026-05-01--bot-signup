import type { ReactNode } from 'react'
import { apiSlice, useGetMeQuery, useLogoutMutation } from '../store/api'
import { AuthContext } from './context'
import { useAppDispatch } from './hooks'

export function AuthProvider({ children }: { children: ReactNode }) {
  const dispatch = useAppDispatch()
  const { data: user, isLoading } = useGetMeQuery(undefined, {
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
    <AuthContext.Provider value={{ user: user ?? null, isLoading, loginWithDiscord, logout }}>
      {children}
    </AuthContext.Provider>
  )
}

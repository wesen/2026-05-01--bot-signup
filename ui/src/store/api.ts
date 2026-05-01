import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react'

export type UserStatus = 'waiting' | 'approved' | 'rejected' | 'suspended'
export type UserRole = 'user' | 'admin'

export interface User {
  id: number
  discord_id: string
  email?: string
  display_name: string
  avatar_url?: string
  status: UserStatus
  role: UserRole
  last_login_at?: string
  created_at: string
  updated_at: string
}

export interface BotCredentials {
  application_id: string
  bot_token: string
  guild_id: string
  public_key: string
  approved_at?: string
}

export interface ProfileResponse {
  user: User
  bot_credentials: BotCredentials | null
  message?: string
}

export interface Stats {
  total_users: number
  approved_users: number
  waiting_users: number
  bots_running: number
}

export const apiSlice = createApi({
  reducerPath: 'api',
  baseQuery: fetchBaseQuery({
    baseUrl: '/api',
    credentials: 'include',
  }),
  tagTypes: ['User', 'Stats'],
  endpoints: (builder) => ({
    getMe: builder.query<User, void>({
      query: () => '/auth/me',
      providesTags: ['User'],
    }),
    logout: builder.mutation<{ message: string }, void>({
      query: () => ({ url: '/auth/logout', method: 'POST' }),
      invalidatesTags: ['User'],
    }),
    getProfile: builder.query<ProfileResponse, void>({
      query: () => '/profile',
      providesTags: ['User'],
    }),
    getStats: builder.query<Stats, void>({
      query: () => '/stats',
      providesTags: ['Stats'],
    }),
  }),
})

export const { useGetMeQuery, useLogoutMutation, useGetProfileQuery, useGetStatsQuery } = apiSlice

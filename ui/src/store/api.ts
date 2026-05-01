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

export interface SessionResponse {
  user: User | null
}

export interface Stats {
  total_users: number
  approved_users: number
  waiting_users: number
  bots_running: number
}

export interface WaitlistResponse {
  users: User[]
  total: number
}

export interface AdminUsersResponse {
  users: User[]
  total: number
  page: number
  per_page: number
}

export interface ApproveUserRequest {
  id: number
  application_id: string
  bot_token: string
  guild_id: string
  public_key: string
}

export interface ApproveUserResponse {
  message: string
  user: User
  bot_credentials: BotCredentials
}

export const apiSlice = createApi({
  reducerPath: 'api',
  baseQuery: fetchBaseQuery({
    baseUrl: '/api',
    credentials: 'include',
  }),
  tagTypes: ['User', 'Stats'],
  endpoints: (builder) => ({
    getSession: builder.query<SessionResponse, void>({
      query: () => '/auth/session',
      providesTags: ['User'],
    }),
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
    getWaitlist: builder.query<WaitlistResponse, void>({
      query: () => '/admin/waitlist',
      providesTags: ['User'],
    }),
    getAdminUsers: builder.query<AdminUsersResponse, { page?: number; per_page?: number; status?: string } | void>({
      query: (params) => ({ url: '/admin/users', params: params ?? undefined }),
      providesTags: ['User'],
    }),
    approveUser: builder.mutation<ApproveUserResponse, ApproveUserRequest>({
      query: ({ id, ...body }) => ({ url: `/admin/users/${id}/approve`, method: 'POST', body }),
      invalidatesTags: ['User', 'Stats'],
    }),
    rejectUser: builder.mutation<{ message: string }, number>({
      query: (id) => ({ url: `/admin/users/${id}/reject`, method: 'POST' }),
      invalidatesTags: ['User', 'Stats'],
    }),
    suspendUser: builder.mutation<{ message: string }, number>({
      query: (id) => ({ url: `/admin/users/${id}/suspend`, method: 'POST' }),
      invalidatesTags: ['User', 'Stats'],
    }),
    deleteUser: builder.mutation<{ message: string }, number>({
      query: (id) => ({ url: `/admin/users/${id}`, method: 'DELETE' }),
      invalidatesTags: ['User', 'Stats'],
    }),
  }),
})

export const {
  useGetSessionQuery,
  useGetMeQuery,
  useLogoutMutation,
  useGetProfileQuery,
  useGetStatsQuery,
  useGetWaitlistQuery,
  useGetAdminUsersQuery,
  useApproveUserMutation,
  useRejectUserMutation,
  useSuspendUserMutation,
  useDeleteUserMutation,
} = apiSlice

import type { Meta, StoryObj } from '@storybook/react-vite'
import { Provider } from 'react-redux'
import { MemoryRouter } from 'react-router-dom'
import { AuthContext } from '../../auth/context'
import { apiSlice, type User } from '../../store/api'
import { store } from '../../store/store'
import { AdminDashboard } from './AdminDashboard'

const admin: User = { id: 99, discord_id: '999', display_name: 'Admin', status: 'waiting', role: 'admin', created_at: new Date().toISOString(), updated_at: new Date().toISOString() }
const waiting: User[] = [
  { id: 1, discord_id: '123', display_name: 'CoolBotDev', email: 'user@example.com', status: 'waiting', role: 'user', created_at: new Date().toISOString(), updated_at: new Date().toISOString() },
]

const meta = {
  title: 'Pages/AdminDashboard',
  component: AdminDashboard,
  tags: ['autodocs'],
  decorators: [
    (Story) => {
      store.dispatch(apiSlice.util.upsertQueryData('getStats', undefined, { total_users: 4, approved_users: 2, waiting_users: 1, bots_running: 2 }))
      store.dispatch(apiSlice.util.upsertQueryData('getWaitlist', undefined, { users: waiting, total: waiting.length }))
      return <Provider store={store}><MemoryRouter><AuthContext.Provider value={{ user: admin, isLoading: false, loginWithDiscord: () => undefined, logout: async () => undefined }}><Story /></AuthContext.Provider></MemoryRouter></Provider>
    },
  ],
} satisfies Meta<typeof AdminDashboard>

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {}

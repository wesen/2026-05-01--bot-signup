import type { Meta, StoryObj } from '@storybook/react-vite'
import { MemoryRouter } from 'react-router-dom'
import type { User } from '../store/api'
import { AdminUserTable } from './AdminUserTable'

const users: User[] = [
  { id: 1, discord_id: '123456789012345678', display_name: 'CoolBotDev', email: 'user@example.com', status: 'waiting', role: 'user', created_at: new Date().toISOString(), updated_at: new Date().toISOString() },
  { id: 2, discord_id: '222456789012345678', display_name: 'BotMaster', email: 'bot@example.com', status: 'approved', role: 'user', created_at: new Date().toISOString(), updated_at: new Date().toISOString() },
  { id: 3, discord_id: '333456789012345678', display_name: 'DisabledDev', email: 'disabled@example.com', status: 'suspended', role: 'user', created_at: new Date().toISOString(), updated_at: new Date().toISOString() },
]

const meta = {
  title: 'Components/AdminUserTable',
  component: AdminUserTable,
  tags: ['autodocs'],
  decorators: [(Story) => <MemoryRouter><div className="bg-slate-50 p-8"><Story /></div></MemoryRouter>],
} satisfies Meta<typeof AdminUserTable>

export default meta
type Story = StoryObj<typeof meta>

export const WithUsers: Story = { args: { users } }
export const Empty: Story = { args: { users: [] } }

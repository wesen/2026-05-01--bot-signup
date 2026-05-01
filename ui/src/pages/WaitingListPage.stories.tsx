import type { Meta, StoryObj } from '@storybook/react-vite'
import { MemoryRouter } from 'react-router-dom'
import { AuthContext } from '../auth/context'
import type { User } from '../store/api'
import { WaitingListPage } from './WaitingListPage'

const baseUser: User = {
  id: 1,
  discord_id: '123456789012345678',
  display_name: 'CoolBotDev',
  email: 'user@example.com',
  status: 'waiting',
  role: 'user',
  created_at: new Date().toISOString(),
  updated_at: new Date().toISOString(),
}

function renderWithUser(user: User) {
  return (
    <MemoryRouter>
      <AuthContext.Provider value={{ user, isLoading: false, loginWithDiscord: () => undefined, logout: async () => undefined }}>
        <WaitingListPage />
      </AuthContext.Provider>
    </MemoryRouter>
  )
}

const meta = {
  title: 'Pages/WaitingListPage',
  component: WaitingListPage,
  tags: ['autodocs'],
} satisfies Meta<typeof WaitingListPage>

export default meta
type Story = StoryObj<typeof meta>

export const Waiting: Story = { render: () => renderWithUser(baseUser) }
export const Approved: Story = { render: () => renderWithUser({ ...baseUser, status: 'approved' }) }
export const Rejected: Story = { render: () => renderWithUser({ ...baseUser, status: 'rejected' }) }

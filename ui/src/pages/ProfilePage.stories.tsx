import type { Meta, StoryObj } from '@storybook/react-vite'
import { Provider } from 'react-redux'
import { MemoryRouter } from 'react-router-dom'
import { AuthContext } from '../auth/context'
import { apiSlice, type ProfileResponse, type User } from '../store/api'
import { store } from '../store/store'
import { ProfilePage } from './ProfilePage'

const user: User = {
  id: 1,
  discord_id: '123456789012345678',
  display_name: 'CoolBotDev',
  email: 'user@example.com',
  status: 'approved',
  role: 'user',
  created_at: new Date().toISOString(),
  updated_at: new Date().toISOString(),
}

const approvedProfile: ProfileResponse = {
  user,
  bot_credentials: {
    application_id: '987654321098765432',
    bot_token: 'MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkw.fake-token',
    guild_id: '111222333444555666',
    public_key: 'abcdef1234567890abcdef1234567890',
    approved_at: new Date().toISOString(),
  },
}

const meta = {
  title: 'Pages/ProfilePage',
  component: ProfilePage,
  tags: ['autodocs'],
  decorators: [
    (Story) => {
      store.dispatch(apiSlice.util.upsertQueryData('getProfile', undefined, approvedProfile))
      return (
        <Provider store={store}>
          <MemoryRouter>
            <AuthContext.Provider value={{ user, isLoading: false, loginWithDiscord: () => undefined, logout: async () => undefined }}>
              <Story />
            </AuthContext.Provider>
          </MemoryRouter>
        </Provider>
      )
    },
  ],
} satisfies Meta<typeof ProfilePage>

export default meta
type Story = StoryObj<typeof meta>

export const Approved: Story = {}

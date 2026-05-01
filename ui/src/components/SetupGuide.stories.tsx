import type { Meta, StoryObj } from '@storybook/react-vite'
import { MemoryRouter } from 'react-router-dom'
import { SetupGuide } from './SetupGuide'

const meta = {
  title: 'Components/SetupGuide',
  component: SetupGuide,
  tags: ['autodocs'],
  args: {
    credentials: {
      application_id: '987654321098765432',
      bot_token: 'MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkw.fake-token',
      guild_id: '111222333444555666',
      public_key: 'abcdef1234567890abcdef1234567890',
    },
  },
  decorators: [
    (Story) => (
      <MemoryRouter>
        <div className="min-h-screen bg-slate-50 p-8">
          <div className="mx-auto max-w-4xl">
            <Story />
          </div>
        </div>
      </MemoryRouter>
    ),
  ],
} satisfies Meta<typeof SetupGuide>

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {}

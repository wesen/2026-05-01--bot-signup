import type { Meta, StoryObj } from '@storybook/react-vite'
import { CredentialTable } from './CredentialTable'

const meta = {
  title: 'Components/CredentialTable',
  component: CredentialTable,
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
      <div className="min-h-screen bg-slate-50 p-8">
        <div className="mx-auto max-w-5xl">
          <Story />
        </div>
      </div>
    ),
  ],
} satisfies Meta<typeof CredentialTable>

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {}

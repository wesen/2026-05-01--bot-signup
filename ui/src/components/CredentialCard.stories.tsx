import type { Meta, StoryObj } from '@storybook/react-vite'
import { CredentialCard } from './CredentialCard'

const meta = {
  title: 'Components/CredentialCard',
  component: CredentialCard,
  tags: ['autodocs'],
  args: {
    label: 'Application ID',
    value: '987654321098765432',
    secret: false,
  },
  decorators: [
    (Story) => (
      <div className="max-w-xl bg-slate-50 p-8">
        <Story />
      </div>
    ),
  ],
} satisfies Meta<typeof CredentialCard>

export default meta
type Story = StoryObj<typeof meta>

export const PublicValue: Story = {}
export const SecretValue: Story = {
  args: {
    label: 'Bot Token',
    value: 'MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkw.fake-token',
    secret: true,
  },
}

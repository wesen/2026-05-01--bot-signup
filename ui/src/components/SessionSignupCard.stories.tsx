import type { Meta, StoryObj } from '@storybook/react-vite'
import { SessionSignupCard } from './SessionSignupCard'

const meta = {
  title: 'Components/SessionSignupCard',
  component: SessionSignupCard,
  tags: ['autodocs'],
  args: {
    onContinueWithDiscord: () => console.log('continue with discord'),
  },
  decorators: [
    (Story) => (
      <div className="min-h-screen bg-slate-50 p-8">
        <div className="mx-auto max-w-md">
          <Story />
        </div>
      </div>
    ),
  ],
} satisfies Meta<typeof SessionSignupCard>

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {}

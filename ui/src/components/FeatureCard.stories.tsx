import type { Meta, StoryObj } from '@storybook/react-vite'
import { KeyRound } from 'lucide-react'
import { FeatureCard } from './FeatureCard'

const meta = {
  title: 'Components/FeatureCard',
  component: FeatureCard,
  tags: ['autodocs'],
  args: {
    icon: KeyRound,
    title: 'Your Bot Credentials',
    description: 'Get your own Discord bot token and client ID to keep.',
  },
  decorators: [
    (Story) => (
      <div className="bg-slate-50 p-8">
        <div className="max-w-sm">
          <Story />
        </div>
      </div>
    ),
  ],
} satisfies Meta<typeof FeatureCard>

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {}

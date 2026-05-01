import type { Meta, StoryObj } from '@storybook/react-vite'
import { DiscordOAuthButton } from './DiscordOAuthButton'

const meta = {
  title: 'Components/DiscordOAuthButton',
  component: DiscordOAuthButton,
  tags: ['autodocs'],
  args: {
    onClick: () => console.log('clicked'),
    children: 'Continue with Discord',
    fullWidth: false,
  },
} satisfies Meta<typeof DiscordOAuthButton>

export default meta
type Story = StoryObj<typeof meta>

export const Inline: Story = {}
export const FullWidth: Story = { args: { fullWidth: true } }

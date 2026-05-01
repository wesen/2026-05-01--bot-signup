import type { Meta, StoryObj } from '@storybook/react-vite'
import { StatusBadge } from './StatusBadge'

const meta = {
  title: 'Components/StatusBadge',
  component: StatusBadge,
  tags: ['autodocs'],
} satisfies Meta<typeof StatusBadge>

export default meta
type Story = StoryObj<typeof meta>

export const Waiting: Story = { args: { status: 'waiting' } }
export const Approved: Story = { args: { status: 'approved' } }
export const Rejected: Story = { args: { status: 'rejected' } }
export const Suspended: Story = { args: { status: 'suspended' } }

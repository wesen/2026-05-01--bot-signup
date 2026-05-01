import type { Meta, StoryObj } from '@storybook/react-vite'
import { MemoryRouter } from 'react-router-dom'
import { GetStartedPage } from './GetStartedPage'

const meta = {
  title: 'Pages/GetStartedPage',
  component: GetStartedPage,
  tags: ['autodocs'],
  decorators: [(Story) => <MemoryRouter><Story /></MemoryRouter>],
} satisfies Meta<typeof GetStartedPage>

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {}

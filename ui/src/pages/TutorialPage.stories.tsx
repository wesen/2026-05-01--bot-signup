import type { Meta, StoryObj } from '@storybook/react-vite'
import { MemoryRouter } from 'react-router-dom'
import { TutorialPage } from './TutorialPage'

const meta = {
  title: 'Pages/TutorialPage',
  component: TutorialPage,
  tags: ['autodocs'],
  decorators: [(Story) => <MemoryRouter><Story /></MemoryRouter>],
} satisfies Meta<typeof TutorialPage>

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {}

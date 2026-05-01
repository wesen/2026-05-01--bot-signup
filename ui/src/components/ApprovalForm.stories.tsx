import type { Meta, StoryObj } from '@storybook/react-vite'
import { ApprovalForm } from './ApprovalForm'

const meta = {
  title: 'Components/ApprovalForm',
  component: ApprovalForm,
  tags: ['autodocs'],
  args: {
    onSubmit: async (values) => console.log(values),
    isSubmitting: false,
  },
  decorators: [(Story) => <div className="max-w-xl bg-slate-50 p-8"><Story /></div>],
} satisfies Meta<typeof ApprovalForm>

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {}
export const Submitting: Story = { args: { isSubmitting: true } }

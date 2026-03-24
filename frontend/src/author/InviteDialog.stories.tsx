import type { Meta, StoryObj } from "@storybook/react-vite";
import { InviteDialog } from "./InviteDialog";

const meta: Meta<typeof InviteDialog> = {
  title: "Author/InviteDialog",
  component: InviteDialog,
  parameters: { layout: "fullscreen" },
};
export default meta;
type Story = StoryObj<typeof InviteDialog>;

export const Default: Story = {
  args: {
    onClose: () => alert("Close"),
    onInvite: (email, note) => alert(`Invite ${email}: ${note}`),
  },
};

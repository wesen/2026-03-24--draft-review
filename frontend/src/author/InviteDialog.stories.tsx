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
    shareUrl: "http://127.0.0.1:8080/r/share-demo-article",
    onGenerateShareLink: async () => {
      return "http://127.0.0.1:8080/r/share-demo-article";
    },
    onInvite: async (email, note) => {
      alert(`Invite ${email}: ${note}`);
      return {
        email,
        inviteUrl: `http://127.0.0.1:8080/r/invite-demo-${email}`,
      };
    },
  },
};

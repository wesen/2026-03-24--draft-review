import type { Meta, StoryObj } from "@storybook/react-vite";
import { DoneDialog } from "./DoneDialog";

const meta: Meta<typeof DoneDialog> = {
  title: "Reader/DoneDialog",
  component: DoneDialog,
  parameters: { layout: "fullscreen" },
};
export default meta;
type Story = StoryObj<typeof DoneDialog>;

export const Default: Story = {
  args: {
    stats: { useful: 3, confused: 2, slow: 1, favorite: 4 },
    onClose: () => alert("Done!"),
  },
};

export const NoReactions: Story = {
  args: {
    stats: { useful: 0, confused: 0, slow: 0, favorite: 0 },
    onClose: () => {},
  },
};

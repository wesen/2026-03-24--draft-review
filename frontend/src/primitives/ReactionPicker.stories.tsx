import type { Meta, StoryObj } from "@storybook/react-vite";
import { ReactionPicker } from "./ReactionPicker";

const meta: Meta<typeof ReactionPicker> = {
  title: "Primitives/ReactionPicker",
  component: ReactionPicker,
  decorators: [(Story) => <div style={{ maxWidth: 400, padding: 16 }}><Story /></div>],
};
export default meta;
type Story = StoryObj<typeof ReactionPicker>;

export const Default: Story = {
  args: {
    onSubmit: (type, text) => alert(`${type}: ${text}`),
    onCancel: () => alert("cancelled"),
  },
};

import type { Meta, StoryObj } from "@storybook/react-vite";
import { MacInput } from "./MacInput";
import { MacTextArea } from "./MacTextArea";

const meta: Meta<typeof MacInput> = {
  title: "Chrome/MacInput",
  component: MacInput,
  decorators: [(Story) => <div style={{ maxWidth: 400, padding: 16 }}><Story /></div>],
};
export default meta;
type Story = StoryObj<typeof MacInput>;

export const Default: Story = {
  args: { label: "Email", placeholder: "you@example.com", type: "email" },
};

export const Password: Story = {
  args: { label: "Password", type: "password", placeholder: "\u2022\u2022\u2022\u2022\u2022\u2022\u2022\u2022" },
};

export const TextArea: StoryObj<typeof MacTextArea> = {
  render: () => (
    <MacTextArea
      label="Personal Note"
      placeholder="Add a note for the reader..."
      rows={4}
    />
  ),
};

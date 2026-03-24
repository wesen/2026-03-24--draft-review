import type { Meta, StoryObj } from "@storybook/react-vite";
import { TitleBar } from "./TitleBar";
import "./TitleBar.css";

const meta: Meta<typeof TitleBar> = {
  title: "Chrome/TitleBar",
  component: TitleBar,
  decorators: [
    (Story) => (
      <div style={{ width: 500, border: "2px solid #000", borderRadius: 3 }}>
        <Story />
      </div>
    ),
  ],
};
export default meta;
type Story = StoryObj<typeof TitleBar>;

export const Default: Story = {
  args: { title: "My Document" },
};

export const WithClose: Story = {
  args: { title: "Closeable Window", onClose: () => alert("close") },
};

export const LongTitle: Story = {
  args: { title: "A Very Long Window Title That Should Still Look Good" },
};

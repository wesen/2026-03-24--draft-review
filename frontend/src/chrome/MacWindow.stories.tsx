import type { Meta, StoryObj } from "@storybook/react-vite";
import { MacWindow } from "./MacWindow";

const meta: Meta<typeof MacWindow> = {
  title: "Chrome/MacWindow",
  component: MacWindow,
  decorators: [
    (Story) => (
      <div style={{ position: "relative", width: "100%", height: 500, background: "#e8e8e8" }}>
        <Story />
      </div>
    ),
  ],
};
export default meta;
type Story = StoryObj<typeof MacWindow>;

export const Default: Story = {
  args: {
    title: "My Window",
    children: (
      <div style={{ padding: 16, fontFamily: "Geneva, Monaco, monospace", fontSize: 12 }}>
        Window content goes here. This window is draggable!
      </div>
    ),
    w: 400,
    h: 300,
  },
};

export const WithClose: Story = {
  args: {
    title: "Closeable Window",
    onClose: () => alert("close"),
    children: <div style={{ padding: 16 }}>Click the close box in the title bar.</div>,
    w: 350,
    h: 200,
  },
};

export const Maximized: Story = {
  args: {
    title: "Maximized Window",
    maximized: true,
    children: <div style={{ padding: 16 }}>This window fills the available space.</div>,
  },
};

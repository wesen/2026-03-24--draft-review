import type { Meta, StoryObj } from "@storybook/react-vite";
import { MenuBar } from "./MenuBar";

const meta: Meta<typeof MenuBar> = {
  title: "Chrome/MenuBar",
  component: MenuBar,
  parameters: { layout: "fullscreen" },
};
export default meta;
type Story = StoryObj<typeof MenuBar>;

export const Default: Story = {
  args: {
    menus: [
      {
        label: "File",
        items: [
          { label: "New Article\u2026", shortcut: "\u2318N" },
          { label: "Import from Docs\u2026", shortcut: "\u2318I" },
          { divider: true, label: "" },
          { label: "Export Feedback\u2026", shortcut: "\u2318E" },
        ],
      },
      {
        label: "Edit",
        items: [
          { label: "Undo", shortcut: "\u2318Z" },
          { label: "Redo", shortcut: "\u2318\u21E7Z" },
        ],
      },
      {
        label: "View",
        items: [{ label: "Dashboard" }, { label: "Articles" }],
      },
      {
        label: "Help",
        items: [{ label: "Beta Reading Guide" }, { label: "Keyboard Shortcuts" }],
      },
    ],
    rightStatus: "3 reactions",
  },
};

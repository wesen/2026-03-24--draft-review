import type { Meta, StoryObj } from "@storybook/react-vite";
import { StatCard } from "./StatCard";

const meta: Meta<typeof StatCard> = {
  title: "Primitives/StatCard",
  component: StatCard,
  decorators: [(Story) => <div style={{ maxWidth: 200, padding: 16 }}><Story /></div>],
};
export default meta;
type Story = StoryObj<typeof StatCard>;

export const Readers: Story = { args: { icon: "\u25C9", value: 5, label: "Readers" } };
export const Reactions: Story = { args: { icon: "\u2726", value: 10, label: "Reactions" } };
export const Sections: Story = { args: { icon: "\u00A7", value: 5, label: "Sections" } };
export const Progress: Story = { args: { icon: "\u25B8", value: "60%", label: "Avg Progress" } };

export const Grid: Story = {
  decorators: [(Story) => <div style={{ maxWidth: 600, padding: 16 }}><Story /></div>],
  render: () => (
    <div style={{ display: "grid", gridTemplateColumns: "repeat(4, 1fr)", gap: 10 }}>
      <StatCard icon="\u25C9" value={5} label="Readers" />
      <StatCard icon="\u2726" value={10} label="Reactions" />
      <StatCard icon="\u00A7" value={5} label="Sections" />
      <StatCard icon="\u25B8" value="60%" label="Avg Progress" />
    </div>
  ),
};

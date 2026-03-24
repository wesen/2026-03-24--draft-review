import type { Meta, StoryObj } from "@storybook/react-vite";
import { ProgressBar } from "./ProgressBar";

const meta: Meta<typeof ProgressBar> = {
  title: "Primitives/ProgressBar",
  component: ProgressBar,
  decorators: [(Story) => <div style={{ maxWidth: 400, padding: 16 }}><Story /></div>],
};
export default meta;
type Story = StoryObj<typeof ProgressBar>;

export const Empty: Story = { args: { percent: 0 } };
export const Quarter: Story = { args: { percent: 25 } };
export const Half: Story = { args: { percent: 50 } };
export const ThreeQuarter: Story = { args: { percent: 75 } };
export const Full: Story = { args: { percent: 100 } };

export const Thin: Story = { args: { percent: 60, height: 8 } };
export const FixedWidth: Story = { args: { percent: 40, width: 120, height: 10 } };

export const AllSizes: Story = {
  render: () => (
    <div style={{ display: "flex", flexDirection: "column", gap: 8 }}>
      <ProgressBar percent={20} height={6} />
      <ProgressBar percent={40} height={8} />
      <ProgressBar percent={60} height={10} />
      <ProgressBar percent={80} height={12} />
      <ProgressBar percent={100} height={16} />
    </div>
  ),
};

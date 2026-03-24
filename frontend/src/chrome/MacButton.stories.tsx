import type { Meta, StoryObj } from "@storybook/react-vite";
import { MacButton } from "./MacButton";

const meta: Meta<typeof MacButton> = {
  title: "Chrome/MacButton",
  component: MacButton,
};
export default meta;
type Story = StoryObj<typeof MacButton>;

export const Default: Story = {
  args: { children: "Click Me" },
};

export const Primary: Story = {
  args: { children: "Primary Action", primary: true },
};

export const Small: Story = {
  args: { children: "Small", small: true },
};

export const SmallPrimary: Story = {
  args: { children: "Submit", small: true, primary: true },
};

export const Disabled: Story = {
  args: { children: "Disabled", disabled: true },
};

export const FullWidth: Story = {
  args: { children: "Full Width Button", fullWidth: true },
  decorators: [(Story) => <div style={{ width: 300 }}><Story /></div>],
};

export const AllVariants: Story = {
  render: () => (
    <div style={{ display: "flex", gap: 8, flexWrap: "wrap", alignItems: "center" }}>
      <MacButton>Default</MacButton>
      <MacButton primary>Primary</MacButton>
      <MacButton small>Small</MacButton>
      <MacButton small primary>Small Primary</MacButton>
      <MacButton disabled>Disabled</MacButton>
    </div>
  ),
};

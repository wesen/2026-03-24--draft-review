import type { Meta, StoryObj } from "@storybook/react-vite";
import { ReactionBadge } from "./ReactionBadge";

const meta: Meta<typeof ReactionBadge> = {
  title: "Primitives/ReactionBadge",
  component: ReactionBadge,
  decorators: [(Story) => <div style={{ padding: 16 }}><Story /></div>],
};
export default meta;
type Story = StoryObj<typeof ReactionBadge>;

export const Useful: Story = {
  args: { type: "useful", text: "Great hook -- immediately relevant to my situation." },
};

export const Confused: Story = {
  args: { type: "confused", text: "What do you mean by 'workarounds'?" },
};

export const Slow: Story = {
  args: { type: "slow", text: "This section drags a bit." },
};

export const Favorite: Story = {
  args: { type: "favorite", text: "YES. This is exactly what happened at my last company." },
};

export const WithRemove: Story = {
  args: { type: "useful", text: "Removable badge", onRemove: () => alert("remove") },
};

export const AllTypes: Story = {
  render: () => (
    <div>
      <ReactionBadge type="useful" text="I can use this" />
      <ReactionBadge type="confused" text="I don't follow" />
      <ReactionBadge type="slow" text="This drags" />
      <ReactionBadge type="favorite" text="Highlight" />
    </div>
  ),
};

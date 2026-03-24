import type { Meta, StoryObj } from "@storybook/react-vite";
import { Dashboard } from "./Dashboard";
import { articles, readers, reactions } from "../mocks/db";

const meta: Meta<typeof Dashboard> = {
  title: "Author/Dashboard",
  component: Dashboard,
  parameters: { layout: "fullscreen" },
  decorators: [
    (Story) => (
      <div style={{ height: "100vh", overflow: "auto", background: "#fff" }}>
        <Story />
      </div>
    ),
  ],
};
export default meta;
type Story = StoryObj<typeof Dashboard>;

export const Default: Story = {
  args: {
    articles,
    readers,
    reactions,
    onSelectArticle: (id, sectionId) =>
      alert(`Select article ${id}${sectionId ? ` section ${sectionId}` : ""}`),
    onInvite: () => alert("Invite"),
  },
};

export const EmptyArticle: Story = {
  args: {
    articles: [articles[2]], // "The Art of Saying No" -- no readers/reactions
    readers: [],
    reactions: [],
    onSelectArticle: () => {},
    onInvite: () => alert("Invite"),
  },
};

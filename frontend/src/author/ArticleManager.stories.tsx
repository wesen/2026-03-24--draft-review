import type { Meta, StoryObj } from "@storybook/react-vite";
import { ArticleManager } from "./ArticleManager";
import { articles, readers, reactions } from "../mocks/db";

const meta: Meta<typeof ArticleManager> = {
  title: "Author/ArticleManager",
  component: ArticleManager,
  parameters: { layout: "fullscreen" },
  decorators: [
    (Story) => (
      <div style={{ height: "100vh" }}>
        <Story />
      </div>
    ),
  ],
};
export default meta;
type Story = StoryObj<typeof ArticleManager>;

export const Default: Story = {
  args: {
    articles,
    readers,
    reactions,
    onEdit: (id) => alert(`Edit ${id}`),
    onSettings: (id) => alert(`Settings ${id}`),
    onReview: (id) => alert(`Review ${id}`),
    onNewArticle: () => alert("New article"),
    onInvite: (id) => alert(`Invite for ${id}`),
  },
};

export const Empty: Story = {
  args: {
    articles: [],
    readers: [],
    reactions: [],
    onEdit: () => {},
    onSettings: () => {},
    onReview: () => {},
    onNewArticle: () => alert("New article"),
    onInvite: () => {},
  },
};

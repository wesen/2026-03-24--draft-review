import type { Meta, StoryObj } from "@storybook/react-vite";
import { ArticleReader } from "./ArticleReader";
import { articles, reactions } from "../mocks/db";

const meta: Meta<typeof ArticleReader> = {
  title: "Author/ArticleReader",
  component: ArticleReader,
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
type Story = StoryObj<typeof ArticleReader>;

export const Default: Story = {
  args: {
    article: articles[0],
    reactions: reactions.filter((r) => r.articleId === "1"),
    onBack: () => alert("Back to dashboard"),
  },
};

export const FocusedSection: Story = {
  args: {
    article: articles[0],
    reactions: reactions.filter((r) => r.articleId === "1"),
    focusSection: "s2",
    onBack: () => alert("Back"),
  },
};

export const NoReactions: Story = {
  args: {
    article: articles[2],
    reactions: [],
    onBack: () => alert("Back"),
  },
};

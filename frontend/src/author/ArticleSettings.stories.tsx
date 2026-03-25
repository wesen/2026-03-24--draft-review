import type { Meta, StoryObj } from "@storybook/react-vite";
import { ArticleSettings } from "./ArticleSettings";
import { articles } from "../mocks/db";

const meta: Meta<typeof ArticleSettings> = {
  title: "Author/ArticleSettings",
  component: ArticleSettings,
  parameters: { layout: "fullscreen" },
  decorators: [
    (Story) => (
      <div style={{ height: "100vh", background: "#fff" }}>
        <Story />
      </div>
    ),
  ],
};
export default meta;
type Story = StoryObj<typeof ArticleSettings>;

export const Default: Story = {
  args: {
    article: articles[0],
    onSave: (updates) => alert(`Saved: ${JSON.stringify(updates)}`),
    onBack: () => alert("Back"),
    onGenerateLink: () => alert("Generate link"),
    shareUrl: "/r/share-1-abc123",
  },
};

export const NoLink: Story = {
  args: {
    article: articles[2],
    onSave: () => {},
    onBack: () => alert("Back"),
    onGenerateLink: () => alert("Generate link"),
  },
};

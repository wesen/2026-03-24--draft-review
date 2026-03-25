import type { Meta, StoryObj } from "@storybook/react-vite";
import { ArticleEditor } from "./ArticleEditor";
import { articles } from "../mocks/db";

const meta: Meta<typeof ArticleEditor> = {
  title: "Author/ArticleEditor",
  component: ArticleEditor,
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
type Story = StoryObj<typeof ArticleEditor>;

export const Default: Story = {
  args: {
    article: articles[0],
    onSave: (a) => alert(`Saved: ${a.title} (${a.sections.length} sections)`),
    onBack: () => alert("Back"),
    onPreview: () => alert("Preview"),
  },
};

export const SingleSection: Story = {
  args: {
    article: articles[2],
    onSave: (a) => alert(`Saved: ${a.title}`),
    onBack: () => alert("Back"),
  },
};

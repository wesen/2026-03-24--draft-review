import type { Meta, StoryObj } from "@storybook/react-vite";
import { ReaderPage } from "./ReaderPage";
import { articles } from "../mocks/db";

const meta: Meta<typeof ReaderPage> = {
  title: "Reader/ReaderPage",
  component: ReaderPage,
  parameters: { layout: "fullscreen" },
  decorators: [
    (Story) => (
      <div style={{ height: "100vh", display: "flex", flexDirection: "column" }}>
        <Story />
      </div>
    ),
  ],
};
export default meta;
type Story = StoryObj<typeof ReaderPage>;

const article1 = articles[0];

export const Default: Story = {
  args: {
    article: {
      id: article1.id,
      title: article1.title,
      author: article1.author,
      version: article1.version,
      intro: article1.intro,
      sections: article1.sections,
    },
  },
};

export const ShortArticle: Story = {
  args: {
    article: {
      id: "3",
      title: "The Art of Saying No",
      author: "Alex Chen",
      version: "Draft 1",
      intro: "A very early draft. Be brutal.",
      sections: [
        {
          id: "s1",
          title: "Why We Say Yes",
          bodyMarkdown:
            "Saying yes feels safe. It keeps the peace, maintains relationships, and avoids the discomfort of disappointing someone. But every yes to something unimportant is a no to something that matters.",
        },
      ],
    },
  },
};

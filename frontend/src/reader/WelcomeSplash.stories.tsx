import type { Meta, StoryObj } from "@storybook/react-vite";
import { WelcomeSplash } from "./WelcomeSplash";

const sampleArticle = {
  id: "1",
  title: "Why Design Systems Fail",
  author: "Alex Chen",
  version: "Draft 3 — March 2026",
  intro:
    "Thanks for reading an early draft! This is a work-in-progress that needs your honest reactions. As you read, highlight anything that feels useful, confusing, or slow — your feedback will directly shape the next revision.",
  sections: [
    { id: "s1", title: "Introduction", bodyMarkdown: "p1\n\np2" },
    { id: "s2", title: "The Adoption Cliff", bodyMarkdown: "p1\n\np2\n\np3" },
    { id: "s3", title: "Governance Gaps", bodyMarkdown: "p1\n\np2\n\np3" },
    { id: "s4", title: "Documentation", bodyMarkdown: "p1\n\np2" },
    { id: "s5", title: "Making It Work", bodyMarkdown: "p1\n\np2\n\np3" },
  ],
};

const meta: Meta<typeof WelcomeSplash> = {
  title: "Reader/WelcomeSplash",
  component: WelcomeSplash,
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
type Story = StoryObj<typeof WelcomeSplash>;

export const Default: Story = {
  args: {
    article: sampleArticle,
    onStart: () => alert("Start reading!"),
  },
};

export const ShortArticle: Story = {
  args: {
    article: {
      ...sampleArticle,
      title: "The Art of Saying No",
      author: "Alex Chen",
      version: "Draft 1",
      intro: "A very early draft. Be brutal.",
      sections: [{ id: "s1", title: "Why We Say Yes", bodyMarkdown: "p1" }],
    },
    onStart: () => {},
  },
};

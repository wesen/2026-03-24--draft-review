import type { Meta, StoryObj } from "@storybook/react-vite";
import { SectionNav } from "./SectionNav";

const sections = [
  { id: "s1", title: "Introduction", paragraphs: [] },
  { id: "s2", title: "The Adoption Cliff", paragraphs: [] },
  { id: "s3", title: "Governance Gaps", paragraphs: [] },
  { id: "s4", title: "Documentation", paragraphs: [] },
  { id: "s5", title: "Making It Work", paragraphs: [] },
];

const meta: Meta<typeof SectionNav> = {
  title: "Primitives/SectionNav",
  component: SectionNav,
  decorators: [(Story) => <div style={{ padding: 16 }}><Story /></div>],
};
export default meta;
type Story = StoryObj<typeof SectionNav>;

export const NoProgress: Story = {
  args: { sections, currentId: "s1", readIds: [], onPick: () => {} },
};

export const PartialProgress: Story = {
  args: { sections, currentId: "s3", readIds: ["s1", "s2"], onPick: () => {} },
};

export const AllRead: Story = {
  args: {
    sections,
    currentId: "s5",
    readIds: ["s1", "s2", "s3", "s4", "s5"],
    onPick: () => {},
  },
};

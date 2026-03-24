import type { Article, Reaction, Reader } from "../types";

/** In-memory mock database seeded with sample data from the prototypes. */

export const articles: Article[] = [
  {
    id: "1",
    title: "Why Design Systems Fail",
    author: "Alex Chen",
    version: "Draft 3 — March 2026",
    status: "in_review",
    intro:
      "Thanks for reading an early draft! This is a work-in-progress that needs your honest reactions. As you read, highlight anything that feels useful, confusing, or slow — your feedback will directly shape the next revision.",
    sections: [
      {
        id: "s1",
        title: "Introduction",
        paragraphs: [
          "Design systems promise consistency and speed, but many teams abandon them within a year. This article explores why that happens and what you can do about it.",
          "The problem isn't usually technical — it's organizational. A design system is a product, and like any product, it fails when nobody owns it, nobody documents it, and nobody listens to its users.",
        ],
      },
      {
        id: "s2",
        title: "The Adoption Cliff",
        paragraphs: [
          "Most design systems see strong initial adoption. Engineers are excited, designers feel heard, and leadership is optimistic about the efficiency gains on the horizon.",
          "But around month six, usage plateaus. Teams start creating workarounds, one-off components multiply, and the system starts to feel like a constraint rather than an enabler. The honeymoon is over.",
          "This is the adoption cliff, and it kills more design systems than bad component APIs ever will. The root cause is almost always a gap between what the system offers and what teams actually need day-to-day.",
        ],
      },
      {
        id: "s3",
        title: "Governance Gaps",
        paragraphs: [
          "Without clear ownership, design systems drift. Who decides when to add a new component? Who approves breaking changes? Who triages the 47 open GitHub issues?",
          "When these questions go unanswered, teams lose trust in the system and start going rogue. They fork components, inline styles, and build parallel mini-systems inside their own repos.",
          "The fix isn't more rules — it's more clarity. A lightweight decision-making framework, published and linked from your docs, can prevent months of confusion.",
        ],
      },
      {
        id: "s4",
        title: "The Documentation Problem",
        paragraphs: [
          "Even well-built systems fail without good docs. If a developer can't figure out how to use a component in under two minutes, they'll build their own. Documentation isn't a nice-to-have — it's the product.",
          "The best design system docs don't just list props. They show real examples, explain when to use (and when not to use) each component, and link to the Figma source so designers and engineers are always looking at the same thing.",
        ],
      },
      {
        id: "s5",
        title: "Making It Work",
        paragraphs: [
          "The teams that succeed treat their design system like a product, not a project. They have dedicated maintainers, regular release cycles, and feedback loops with the teams who consume the system.",
          "Executive sponsorship matters too — not as a rubber stamp, but as air cover. When a product team wants to skip the system to hit a deadline, someone senior needs to weigh the long-term cost.",
          "If you're starting a design system today, start small: one token file, one button, one clear owner. Grow it only when real demand pulls you forward. That patience is the hardest part — and the most important.",
        ],
      },
    ],
    createdAt: "2026-03-01T00:00:00Z",
    updatedAt: "2026-03-10T00:00:00Z",
  },
  {
    id: "2",
    title: "Remote Work Isn't Working",
    author: "Alex Chen",
    version: "Draft 1",
    status: "draft",
    intro:
      "I'd love your thoughts on this early draft about remote work. Focus on what resonates and what feels off.",
    sections: [
      {
        id: "s1",
        title: "The Broken Promise",
        paragraphs: [
          "We were told remote work would free us. In many ways it has. But the cracks are showing: loneliness, miscommunication, and a creeping sense that we've traded one set of problems for another.",
        ],
      },
      {
        id: "s2",
        title: "What We Lost",
        paragraphs: [
          "Serendipity. The hallway conversation that sparks an idea. The lunch that becomes a brainstorm. These aren't just perks — they're how organizations actually think.",
        ],
      },
      {
        id: "s3",
        title: "A Middle Path",
        paragraphs: [
          "The answer isn't going back to the office five days a week. It's being intentional about when and why we gather. Async by default, sync with purpose.",
        ],
      },
    ],
    createdAt: "2026-03-15T00:00:00Z",
    updatedAt: "2026-03-15T00:00:00Z",
  },
  {
    id: "3",
    title: "The Art of Saying No",
    author: "Alex Chen",
    version: "Draft 1",
    status: "draft",
    intro: "A very early draft. Be brutal.",
    sections: [
      {
        id: "s1",
        title: "Why We Say Yes",
        paragraphs: [
          "Saying yes feels safe. It keeps the peace, maintains relationships, and avoids the discomfort of disappointing someone. But every yes to something unimportant is a no to something that matters.",
        ],
      },
    ],
    createdAt: "2026-03-20T00:00:00Z",
    updatedAt: "2026-03-20T00:00:00Z",
  },
];

export const readers: Reader[] = [
  { id: "r1", name: "Sarah K.", email: "sarah@example.com", avatar: "SK", articleId: "1", progress: 100, token: "tok-sarah-1", invitedAt: "2026-03-05T00:00:00Z", lastActiveAt: "2026-03-12T00:00:00Z" },
  { id: "r2", name: "Marcus T.", email: "marcus@example.com", avatar: "MT", articleId: "1", progress: 60, token: "tok-marcus-1", invitedAt: "2026-03-06T00:00:00Z", lastActiveAt: "2026-03-11T00:00:00Z" },
  { id: "r3", name: "Priya R.", email: "priya@example.com", avatar: "PR", articleId: "1", progress: 80, token: "tok-priya-1", invitedAt: "2026-03-05T00:00:00Z", lastActiveAt: "2026-03-12T00:00:00Z" },
  { id: "r4", name: "James L.", email: "james@example.com", avatar: "JL", articleId: "1", progress: 40, token: "tok-james-1", invitedAt: "2026-03-07T00:00:00Z", lastActiveAt: "2026-03-09T00:00:00Z" },
  { id: "r5", name: "Chen W.", email: "chen@example.com", avatar: "CW", articleId: "1", progress: 20, token: "tok-chen-1", invitedAt: "2026-03-08T00:00:00Z", lastActiveAt: "2026-03-09T00:00:00Z" },
  { id: "r6", name: "Sarah K.", email: "sarah@example.com", avatar: "SK", articleId: "2", progress: 100, token: "tok-sarah-2", invitedAt: "2026-03-16T00:00:00Z", lastActiveAt: "2026-03-18T00:00:00Z" },
  { id: "r7", name: "James L.", email: "james@example.com", avatar: "JL", articleId: "2", progress: 33, token: "tok-james-2", invitedAt: "2026-03-16T00:00:00Z" },
];

export const reactions: Reaction[] = [
  { id: "rx1", articleId: "1", sectionId: "s1", paragraphId: "s1-p0", readerId: "r1", readerName: "Sarah K.", type: "useful", text: "Great hook — immediately relevant to my situation.", createdAt: "2026-03-10T10:00:00Z" },
  { id: "rx2", articleId: "1", sectionId: "s2", paragraphId: "s2-p1", readerId: "r2", readerName: "Marcus T.", type: "confused", text: "What do you mean by 'workarounds'? Can you give an example?", createdAt: "2026-03-10T11:00:00Z" },
  { id: "rx3", articleId: "1", sectionId: "s2", paragraphId: "s2-p2", readerId: "r3", readerName: "Priya R.", type: "slow", text: "This section drags a bit. Maybe tighten the middle paragraph?", createdAt: "2026-03-10T12:00:00Z" },
  { id: "rx4", articleId: "1", sectionId: "s2", paragraphId: "s2-p1", readerId: "r1", readerName: "Sarah K.", type: "favorite", text: "YES. This is exactly what happened at my last company.", createdAt: "2026-03-10T13:00:00Z" },
  { id: "rx5", articleId: "1", sectionId: "s3", paragraphId: "s3-p0", readerId: "r4", readerName: "James L.", type: "confused", text: "Who typically owns this? Would help to see a real org chart.", createdAt: "2026-03-10T14:00:00Z" },
  { id: "rx6", articleId: "1", sectionId: "s3", paragraphId: "s3-p2", readerId: "r3", readerName: "Priya R.", type: "useful", text: "The governance framework idea is gold. Worth expanding.", createdAt: "2026-03-10T15:00:00Z" },
  { id: "rx7", articleId: "1", sectionId: "s4", paragraphId: "s4-p0", readerId: "r1", readerName: "Sarah K.", type: "useful", text: "The 2-minute rule is a great benchmark.", createdAt: "2026-03-10T16:00:00Z" },
  { id: "rx8", articleId: "1", sectionId: "s4", paragraphId: "s4-p0", readerId: "r2", readerName: "Marcus T.", type: "favorite", text: "Documentation IS the product. Quotable.", createdAt: "2026-03-10T17:00:00Z" },
  { id: "rx9", articleId: "1", sectionId: "s5", paragraphId: "s5-p0", readerId: "r3", readerName: "Priya R.", type: "slow", text: "Feels rushed compared to earlier sections. Needs more depth.", createdAt: "2026-03-10T18:00:00Z" },
  { id: "rx10", articleId: "1", sectionId: "s1", paragraphId: "s1-p0", readerId: "r5", readerName: "Chen W.", type: "useful", text: "Strong opening. Sets the right expectations.", createdAt: "2026-03-10T19:00:00Z" },
  { id: "rx11", articleId: "2", sectionId: "s1", paragraphId: "s1-p0", readerId: "r6", readerName: "Sarah K.", type: "favorite", text: "This opening paragraph is perfect.", createdAt: "2026-03-18T10:00:00Z" },
  { id: "rx12", articleId: "2", sectionId: "s1", paragraphId: "s1-p0", readerId: "r7", readerName: "James L.", type: "confused", text: "What cracks specifically? Feels vague.", createdAt: "2026-03-18T11:00:00Z" },
];

let nextReactionId = 13;

export function addReaction(reaction: Omit<Reaction, "id" | "createdAt">): Reaction {
  const newReaction: Reaction = {
    ...reaction,
    id: `rx${nextReactionId++}`,
    createdAt: new Date().toISOString(),
  };
  reactions.push(newReaction);
  return newReaction;
}

import { http, HttpResponse } from "msw";
import { articles, readers, reactions, addReaction, updateArticle, createArticle, generateShareToken } from "./db";

export const handlers = [
  // Articles
  http.get("/api/articles", () => {
    return HttpResponse.json(articles);
  }),

  http.get("/api/articles/:id", ({ params }) => {
    const article = articles.find((a) => a.id === params.id);
    if (!article) return new HttpResponse(null, { status: 404 });
    return HttpResponse.json(article);
  }),

  // Create article
  http.post("/api/articles", async ({ request }) => {
    const body = (await request.json()) as Record<string, unknown>;
    const article = createArticle({
      title: (body.title as string) || "Untitled Article",
      author: (body.author as string) || "You",
      intro: (body.intro as string) || "",
    });
    return HttpResponse.json(article, { status: 201 });
  }),

  // Update article
  http.patch("/api/articles/:id", async ({ params, request }) => {
    const body = (await request.json()) as Record<string, unknown>;
    const updated = updateArticle(params.id as string, body);
    if (!updated) return new HttpResponse(null, { status: 404 });
    return HttpResponse.json(updated);
  }),

  // Delete article
  http.delete("/api/articles/:id", ({ params }) => {
    const index = articles.findIndex((a) => a.id === params.id);
    if (index === -1) return new HttpResponse(null, { status: 404 });
    articles.splice(index, 1);
    return new HttpResponse(null, { status: 204 });
  }),

  // Generate share token
  http.post("/api/articles/:id/share-token", ({ params }) => {
    const token = generateShareToken(params.id as string);
    return HttpResponse.json({ token, url: `/r/${token}` });
  }),

  // Readers for an article
  http.get("/api/articles/:id/readers", ({ params }) => {
    const articleReaders = readers.filter((r) => r.articleId === params.id);
    return HttpResponse.json(articleReaders);
  }),

  // Reactions for an article
  http.get("/api/articles/:id/reactions", ({ params }) => {
    const articleReactions = reactions.filter((r) => r.articleId === params.id);
    return HttpResponse.json(articleReactions);
  }),

  // Add a reaction
  http.post("/api/articles/:id/reactions", async ({ params, request }) => {
    const body = (await request.json()) as Record<string, unknown>;
    const newReaction = addReaction({
      articleId: params.id as string,
      sectionId: body.sectionId as string,
      paragraphId: body.paragraphId as string,
      readerId: (body.readerId as string) || "anonymous",
      readerName: (body.readerName as string) || "Anonymous",
      type: body.type as "useful" | "confused" | "slow" | "favorite",
      text: body.text as string,
    });
    return HttpResponse.json(newReaction, { status: 201 });
  }),

  // Reader link resolution
  http.get("/api/r/:token", ({ params }) => {
    const reader = readers.find((r) => r.token === params.token);
    if (!reader) return new HttpResponse(null, { status: 404 });
    const article = articles.find((a) => a.id === reader.articleId);
    if (!article) return new HttpResponse(null, { status: 404 });
    return HttpResponse.json({
      reader: { id: reader.id, name: reader.name },
      article: {
        id: article.id,
        title: article.title,
        author: article.author,
        version: article.version,
        intro: article.intro,
        sections: article.sections,
      },
    });
  }),

  // Invite a reader
  http.post("/api/articles/:id/invite", async ({ params, request }) => {
    const body = (await request.json()) as Record<string, unknown>;
    const newReader = {
      id: `r${readers.length + 1}`,
      name: (body.email as string).split("@")[0],
      email: body.email as string,
      avatar: (body.email as string).slice(0, 2).toUpperCase(),
      articleId: params.id as string,
      progress: 0,
      token: `tok-${Date.now()}`,
      invitedAt: new Date().toISOString(),
    };
    readers.push(newReader);
    return HttpResponse.json(newReader, { status: 201 });
  }),
];

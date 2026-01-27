import { defineCollection, z } from 'astro:content';

const projects = defineCollection({
  schema: z.object({
    title: z.string(),
    description: z.string(),
    category: z.enum(['Web Dev', 'Data Science', 'AI Engineering']),
    pubDate: z.date(),
    image: z.string().optional(),
    tags: z.array(z.string()),
    link: z.string().url(),
  }),
});

export const collections = { projects };